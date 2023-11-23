package io

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"strings"

	"github.com/fredericlemoine/fastqutils/fastq"
)

var ErrBufferFull = errors.New("bufio: buffer full")

type FastQParser struct {
	reader1 *bufio.Reader // First read file
	reader2 *bufio.Reader // paired read file (if any, nil otherwise)
}

func NewSingleEndParser(file string) (fp *FastQParser, err error) {
	var reader *bufio.Reader
	if reader, err = GetReader(file); err != nil {
		return
	}

	fp = &FastQParser{
		reader,
		nil,
	}

	return
}

func NewPairedEndParser(read1 string, read2 string) (fp *FastQParser, err error) {
	var reader1, reader2 *bufio.Reader

	if reader1, err = GetReader(read1); err != nil {
		return
	}
	if reader2, err = GetReader(read2); err != nil {
		return
	}
	fp = &FastQParser{
		reader1,
		reader2,
	}
	return
}

func GetReader(file string) (reader *bufio.Reader, err error) {
	var fi *os.File
	var gr *gzip.Reader

	if file == "stdin" || file == "-" {
		fi = os.Stdin
	} else {
		fi, err = os.Open(file)
		if err != nil {
			return
		}
	}

	if strings.HasSuffix(file, ".gz") {
		if gr, err = gzip.NewReader(fi); err != nil {
			return
		} else {
			reader = bufio.NewReader(gr)
		}
	} else {
		reader = bufio.NewReader(fi)
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
