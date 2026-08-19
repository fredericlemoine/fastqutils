package io

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/fredericlemoine/fastqutils/dsrc"
	"github.com/fredericlemoine/fastqutils/fastq"
)

func WriteEntry(w *bufio.Writer, entry *fastq.FastqEntry) {
	w.WriteString(
		fmt.Sprintf("%s\n%s\n+\n%s\n",
			entry.Name,
			entry.Sequence,
			entry.Quality,
		),
	)
}

func WriteEntryFasta(w *bufio.Writer, entry *fastq.FastqEntry) {
	w.WriteString(fmt.Sprintf(">%s\n%s\n", entry.Name, entry.Sequence))
}

// multiCloser flushes the buffered writer and then closes, in order,
// every underlying resource that needs it: a gzip stream and its file,
// a plain file, or a dsrc compression subprocess. Callers should call
// Close exactly once when they are done writing.
type multiCloser struct {
	w       *bufio.Writer
	closers []io.Closer
}

func (m *multiCloser) Close() error {
	if err := m.w.Flush(); err != nil {
		return fmt.Errorf("flushing output: %w", err)
	}
	for _, c := range m.closers {
		if c == nil {
			continue
		}
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// GetWriter opens file for writing FASTQ (or FASTA) data and returns a
// buffered writer plus a Closer. The Closer must be called once writing
// is finished (e.g. via defer): it flushes any buffered data and then
// finalizes/closes whatever is underneath — a plain file, a gzip
// stream, or a dsrc compression subprocess.
//
// At most one of gziped/dsrced may be true; requesting both is an
// error. When gziped is true, ".gz" is appended to file, mirroring the
// previous behavior. When dsrced is true, ".dsrc" is appended to file;
// because DSRC archives are written directly to a real file (not a
// stream), file may not be "stdout" or "-" in that case.
func GetWriter(file string, gziped, dsrced bool) (w *bufio.Writer, closer io.Closer, err error) {
	if gziped && dsrced {
		return nil, nil, fmt.Errorf("--gz and --dsrc are mutually exclusive: choose at most one output compression")
	}

	if dsrced {
		if file == "stdout" || file == "-" {
			return nil, nil, fmt.Errorf("dsrc output cannot be streamed to stdout: DSRC archives are written directly to a file, please provide a real --output file name")
		}
		var dc io.Closer
		if w, dc, err = dsrc.GetWriter(file + dsrc.Extension); err != nil {
			return nil, nil, err
		}
		return w, &multiCloser{w: w, closers: []io.Closer{dc}}, nil
	}

	ext := ""
	if gziped {
		ext = ".gz"
	}

	var fi *os.File
	if file == "stdout" || file == "-" {
		fi = os.Stdout
	} else {
		if fi, err = os.Create(file + ext); err != nil {
			return nil, nil, err
		}
	}

	if gziped {
		gw := gzip.NewWriter(fi)
		w = bufio.NewWriter(gw)
		return w, &multiCloser{w: w, closers: []io.Closer{gw, fi}}, nil
	}

	w = bufio.NewWriter(fi)
	return w, &multiCloser{w: w, closers: []io.Closer{fi}}, nil
}
