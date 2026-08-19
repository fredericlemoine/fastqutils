// Package dsrc integrates fastqutils with DSRC (DNA Sequence Reads
// Compressor, https://github.com/refresh-bio/DSRC), a specialized FASTQ
// compressor.
//
// DSRC is a separate C++ program distributed under the GPLv2 license.
// fastqutils does not link against it (fastqutils itself is Apache-2.0
// licensed, and statically embedding GPLv2 code would be a real license
// entanglement, not just an engineering choice). Instead this package
// shells out to the `dsrc` executable found on $PATH (or pointed to by
// the DSRC_BIN environment variable), using DSRC's own streaming mode
// (`-s`) so that raw FASTQ data flows through a pipe rather than a
// temporary file. This keeps fastqutils itself dependency- and
// cgo-free, so normal `go build`/cross-compilation keeps working
// exactly as before for anyone who doesn't touch .dsrc files.
//
// Known caveat: DSRC's `-s` streaming mode has an open upstream
// reliability report (https://github.com/refresh-bio/DSRC/issues/1),
// and building DSRC from source surfaces real mismatched new[]/delete
// bugs in its buffer-handling code (Buffer.h, QualityPositionModeler.cpp,
// huffman.cpp) — undefined behavior that different platforms' allocators
// tolerate differently. It has been observed to segfault on macOS while
// working fine on Linux. If that bites you, the fix is to stop passing
// `-s` and instead round-trip through real temp files using DSRC's
// plain file-to-file interface (`dsrc c in out.dsrc` / `dsrc d in.dsrc
// out>`), which is DSRC's primary, best-tested code path; ask for that
// version of this file if `-s` proves unreliable in your environment.
//
// DSRC must be installed separately; see
// https://github.com/refresh-bio/DSRC for build instructions.
package dsrc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Extension is the file extension fastqutils uses to recognize and name
// DSRC archives.
const Extension = ".dsrc"

// IsDsrcFile reports whether file looks like a DSRC archive, based on its
// extension.
func IsDsrcFile(file string) bool {
	return strings.HasSuffix(file, Extension)
}

// binary resolves the path to the dsrc executable, honoring the DSRC_BIN
// environment variable override, and falling back to $PATH.
func binary() (string, error) {
	if bin := os.Getenv("DSRC_BIN"); bin != "" {
		return bin, nil
	}
	path, err := exec.LookPath("dsrc")
	if err != nil {
		return "", fmt.Errorf(
			"dsrc: the 'dsrc' executable was not found on $PATH; " +
				"install DSRC (https://github.com/refresh-bio/DSRC) to read/write .dsrc files, " +
				"or set the DSRC_BIN environment variable to its full path",
		)
	}
	return path, nil
}

// procCloser releases a running dsrc subprocess. Close must be called
// exactly once when the caller is done reading/writing, or the
// subprocess (and, for writers, the archive it is finalizing) will leak.
type procCloser struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser // set for writers (GetWriter), nil for readers
	stderr *strings.Builder
}

func (c *procCloser) Close() error {
	var closeErr error
	if c.stdin != nil {
		// Closing stdin tells dsrc the raw FASTQ stream is finished; this
		// is what makes it flush and finalize the archive on disk.
		closeErr = c.stdin.Close()
	}

	waitErr := c.cmd.Wait()
	if waitErr != nil {
		msg := strings.TrimSpace(c.stderr.String())
		if msg != "" {
			return fmt.Errorf("dsrc: %w: %s", waitErr, msg)
		}
		return fmt.Errorf("dsrc: %w", waitErr)
	}
	if closeErr != nil {
		return fmt.Errorf("dsrc: closing input pipe: %w", closeErr)
	}
	return nil
}

// GetReader starts `dsrc d -s <file>` and returns a bufio.Reader that
// streams the decompressed FASTQ data from the subprocess's stdout,
// along with a Closer. The Closer must be called once reading is
// finished (typically via defer right after a successful call) — it
// waits for the dsrc subprocess to exit and surfaces anything it wrote
// to stderr as an error.
func GetReader(file string) (*bufio.Reader, io.Closer, error) {
	bin, err := binary()
	if err != nil {
		return nil, nil, err
	}

	// -s: stream decompressed FASTQ data to stdout instead of a file.
	cmd := exec.Command(bin, "d", "-s", file)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("dsrc: %w", err)
	}
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("dsrc: failed to start %s: %w", bin, err)
	}

	return bufio.NewReader(stdout), &procCloser{cmd: cmd, stderr: &stderr}, nil
}

// GetWriter starts `dsrc c -s <file>` and returns a bufio.Writer that
// streams raw FASTQ data to the subprocess's stdin; dsrc compresses it
// and writes the archive directly to `file`. Unlike decompression,
// compressed output cannot be streamed to stdout — DSRC archives are
// written as a real, seekable file — so `file` must be an actual path.
//
// The returned Closer MUST be called after the caller has flushed the
// bufio.Writer: it closes the input pipe (which signals dsrc that the
// FASTQ stream is complete) and waits for dsrc to finish writing the
// archive. Until Close returns without error, the archive on disk is
// not guaranteed to be complete or valid.
func GetWriter(file string) (*bufio.Writer, io.Closer, error) {
	bin, err := binary()
	if err != nil {
		return nil, nil, err
	}

	// -s: read raw FASTQ data from stdin instead of a file.
	cmd := exec.Command(bin, "c", "-s", file)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("dsrc: %w", err)
	}
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("dsrc: failed to start %s: %w", bin, err)
	}

	return bufio.NewWriter(stdin), &procCloser{cmd: cmd, stdin: stdin, stderr: &stderr}, nil
}
