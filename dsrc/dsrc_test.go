package dsrc

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

// requireDsrc skips the test if the dsrc CLI is not available, so this
// suite doesn't fail CI/dev environments that haven't installed DSRC
// (see README.md).
func requireDsrc(t *testing.T) {
	t.Helper()
	if _, err := binary(); err != nil {
		t.Skip("dsrc executable not found on $PATH (or DSRC_BIN): " + err.Error())
	}
}

func TestRoundTrip(t *testing.T) {
	requireDsrc(t)

	dir := t.TempDir()
	// GetWriter writes the archive to exactly the path it's given (unlike
	// io.GetWriter, which appends the .dsrc extension before calling
	// here), so pass the full archive path.
	archiveFile := filepath.Join(dir, "reads"+Extension)

	const fastq = "@read1\nACGTACGTAC\n+\nIIIIIIIIII\n@read2\nTTTTGGGGCC\n+\n!!!!!!!!!!\n"

	w, wc, err := GetWriter(archiveFile)
	if err != nil {
		t.Fatalf("GetWriter: %v", err)
	}
	if _, err := w.WriteString(fastq); err != nil {
		t.Fatalf("writing fastq data: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("closing writer: %v", err)
	}

	if _, err := os.Stat(archiveFile); err != nil {
		t.Fatalf("expected archive at %s: %v", archiveFile, err)
	}

	r, rc, err := GetReader(archiveFile)
	if err != nil {
		t.Fatalf("GetReader: %v", err)
	}
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading back: %v", err)
	}
	if err := rc.Close(); err != nil {
		t.Fatalf("closing reader: %v", err)
	}

	if string(got) != fastq {
		t.Fatalf("round trip mismatch:\n got: %q\nwant: %q", got, fastq)
	}
}

func TestGetReaderMissingBinary(t *testing.T) {
	t.Setenv("DSRC_BIN", "")
	t.Setenv("PATH", t.TempDir()) // an empty directory: dsrc cannot be found
	if _, _, err := GetReader("whatever.dsrc"); err == nil {
		t.Fatal("expected an error when the dsrc binary cannot be found")
	}
}

func TestIsDsrcFile(t *testing.T) {
	cases := map[string]bool{
		"reads.dsrc":     true,
		"reads.fastq":    false,
		"reads.fastq.gz": false,
		"reads":          false,
	}
	for file, want := range cases {
		if got := IsDsrcFile(file); got != want {
			t.Errorf("IsDsrcFile(%q) = %v, want %v", file, got, want)
		}
	}
}
