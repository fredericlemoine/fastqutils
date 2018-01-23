package io

import (
	"bufio"
	"compress/gzip"
	"errors"
	errorp "github.com/fredericlemoine/fastqutils/error"
	"github.com/fredericlemoine/fastqutils/fastq"
	"os"
	"strings"
)

var ErrBufferFull = errors.New("bufio: buffer full")

type FastQParser struct {
	reader1 *bufio.Reader // First read file
	reader2 *bufio.Reader // paired read file (if any, nil otherwise)
}

func NewSingleEndParser(file string) *FastQParser {
	reader := getReader(file)
	return &FastQParser{
		reader,
		nil,
	}
}

func NewPairedEndParser(read1 string, read2 string) *FastQParser {
	reader1 := getReader(read1)
	reader2 := getReader(read2)
	return &FastQParser{
		reader1,
		reader2,
	}
}

func getReader(file string) (reader *bufio.Reader) {
	var fi *os.File
	var err error
	if file == "stdin" || file == "-" {
		fi = os.Stdin
	} else {
		fi, err = os.Open(file)
		if err != nil {
			errorp.ExitWithMessage(err)
		}
	}

	if strings.HasSuffix(file, ".gz") {
		if gr, err := gzip.NewReader(fi); err != nil {
			errorp.ExitWithMessage(err)
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
func readln(r *bufio.Reader) (ln []byte, err error) {
	if ln, err = r.ReadBytes('\n'); err == nil && ln[len(ln)-1] == '\n' {
		ln = ln[:len(ln)-1]
	}
	return
}

// Returns the next entries:
// If paired end returns 2 fastq entries
// If single end: returns 1 entry and nil
func (p *FastQParser) NextEntry() (*fastq.FastqEntry, *fastq.FastqEntry, error) {
	var name1, name2 []byte
	var seq1, seq2 []byte
	var qual1, qual2 []byte
	var err error
	var entry1, entry2 *fastq.FastqEntry

	if name1, err = readln(p.reader1); err != nil {
		return nil, nil, err
	}
	if seq1, err = readln(p.reader1); err != nil {
		return nil, nil, err
	}
	if _, err = readln(p.reader1); err != nil {
		return nil, nil, err
	}
	if qual1, err = readln(p.reader1); err != nil {
		return nil, nil, err
	}
	entry1 = &fastq.FastqEntry{name1, seq1, qual1}

	if p.reader2 != nil {
		if name2, err = readln(p.reader2); err != nil {
			return nil, nil, err
		}
		if seq2, err = readln(p.reader2); err != nil {
			return nil, nil, err
		}
		if _, err = readln(p.reader2); err != nil {
			return nil, nil, err
		}
		if qual2, err = readln(p.reader2); err != nil {
			return nil, nil, err
		}
		if len(seq2) != len(qual2) {
			errorp.ExitWithMessage(errors.New("Length of sequence is different from length of quality"))
		}
		entry2 = &fastq.FastqEntry{name2, seq2, qual2}
	}
	return entry1, entry2, nil
}
