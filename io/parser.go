package io

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fredericlemoine/fastqutils/dsrc"
	"github.com/fredericlemoine/fastqutils/fastq"
)

var ErrBufferFull = errors.New("bufio: buffer full")

type FastQParser struct {
	reader1 *bufio.Reader // First read file
	reader2 *bufio.Reader // paired read file (if any, nil otherwise)
	closer1 io.Closer     // non-nil when reader1 is backed by a resource that must be released (e.g. a dsrc subprocess)
	closer2 io.Closer     // same, for reader2
}

func NewSingleEndParser(file string) (fp *FastQParser, err error) {
	var reader *bufio.Reader
	var closer io.Closer
	if reader, closer, err = GetReader(file); err != nil {
		return
	}

	fp = &FastQParser{
		reader1: reader,
		closer1: closer,
	}

	return
}

func NewPairedEndParser(read1 string, read2 string) (fp *FastQParser, err error) {
	var reader1, reader2 *bufio.Reader
	var closer1, closer2 io.Closer

	if reader1, closer1, err = GetReader(read1); err != nil {
		return
	}
	if reader2, closer2, err = GetReader(read2); err != nil {
		return
	}
	fp = &FastQParser{
		reader1: reader1,
		reader2: reader2,
		closer1: closer1,
		closer2: closer2,
	}
	return
}

// Close releases any resources backing the parser's underlying readers
// (for example, waits for and releases a dsrc decompression
// subprocess). It is a no-op for plain-text input and a simple
// file/gzip close for .gz input. Callers should call it (e.g. via
// defer) once they are done reading, right after a successful call to
// NewSingleEndParser or NewPairedEndParser.
func (p *FastQParser) Close() (err error) {
	if p.closer1 != nil {
		if cerr := p.closer1.Close(); cerr != nil {
			err = cerr
		}
	}
	if p.closer2 != nil {
		if cerr := p.closer2.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}
	return
}

// multiReadCloser closes several underlying resources, in order, when
// the reader they back is done being used.
type multiReadCloser struct {
	closers []io.Closer
}

func (m *multiReadCloser) Close() error {
	var err error
	for _, c := range m.closers {
		if c == nil {
			continue
		}
		if cerr := c.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}
	return err
}

// GetReader opens file for reading FASTQ data, auto-detecting gzip
// (.gz) and DSRC (.dsrc) compression from the file extension. It
// returns a buffered reader plus a Closer for any backing resource (an
// open file, a gzip stream, or — for .dsrc — a decompression
// subprocess) that should be released once reading is finished; closer
// may be nil when there is nothing to release (e.g. reading from
// stdin).
func GetReader(file string) (reader *bufio.Reader, closer io.Closer, err error) {
	if dsrc.IsDsrcFile(file) {
		return dsrc.GetReader(file)
	}

	var fi *os.File
	var gr *gzip.Reader

	if file == "stdin" || file == "-" {
		fi = os.Stdin
	} else {
		if fi, err = os.Open(file); err != nil {
			return
		}
	}

	if strings.HasSuffix(file, ".gz") {
		if gr, err = gzip.NewReader(fi); err != nil {
			return
		}
		reader = bufio.NewReader(gr)
	} else {
		reader = bufio.NewReader(fi)
	}

	// Close the gzip reader before the underlying file, since it may
	// need to read trailing bytes from it to verify the checksum.
	var closers []io.Closer
	if gr != nil {
		closers = append(closers, gr)
	}
	if fi != nil && fi != os.Stdin {
		closers = append(closers, fi)
	}
	if len(closers) > 0 {
		closer = &multiReadCloser{closers}
	}
	return
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (name, seq, qual []byte, err error) {
	if name, err = r.ReadBytes('\n'); err == nil && name[len(name)-1] == '\n' {
		name = name[:len(name)-1]
		if seq, err = r.ReadBytes('\n'); err == nil && seq[len(seq)-1] == '\n' {
			seq = seq[:len(seq)-1]
			r.ReadBytes('\n') // skip one line
			if qual, err = r.ReadBytes('\n'); err == nil && qual[len(qual)-1] == '\n' {
				qual = qual[:len(qual)-1]
			}
		}
	}
	return
}

// Returns the next entries:
// If paired end returns 2 fastq entries
// If single end: returns 1 entry and nil
func (p *FastQParser) NextEntry() (entry1 *fastq.FastqEntry, entry2 *fastq.FastqEntry, err error) {
	var name1, name2 []byte
	var seq1, seq2 []byte
	var qual1, qual2 []byte
	if name1, seq1, qual1, err = Readln(p.reader1); err != nil {
		return
	}
	entry1 = &fastq.FastqEntry{
		Name:     name1,
		Sequence: seq1,
		Quality:  qual1,
	}

	if p.reader2 != nil {
		if name2, seq2, qual2, err = Readln(p.reader2); err != nil {
			return
		}

		if len(seq2) != len(qual2) {
			err = errors.New("length of sequence is different from length of quality")
			return
		}
		entry2 = &fastq.FastqEntry{
			Name:     name2,
			Sequence: seq2,
			Quality:  qual2,
		}
	}
	return
}
