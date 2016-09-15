package fastq

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fredericlemoine/fastqutils/error"
	"math/rand"
)

type FastqEntry struct {
	Name     string
	Sequence []byte
	Quality  []byte
}

func NewFastQEntry(name string, seq []byte, qual []byte) *FastqEntry {
	return &FastqEntry{
		name,
		seq,
		qual,
	}
}

/* Generates a Fastq Entry */
func GenFastQEntry(length int, id int) *FastqEntry {
	name := fmt.Sprintf("read%d", id)
	seq := genseq(length)
	qual := genqual(length)
	return &FastqEntry{
		name,
		seq,
		qual,
	}
}

func genseq(length int) []byte {
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		buf.WriteByte(nt(rand.Intn(4)))
	}
	return buf.Bytes()
}

// Returns the nt
func nt(n int) byte {
	switch n {
	case 0:
		return 'A'
	case 1:
		return 'C'
	case 2:
		return 'G'
	case 3:
		return 'T'
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("No nucleotide with code %d", n)))
	}
	return '\n'
}

func genqual(length int) []byte {
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		buf.WriteByte(byte(rand.Intn(56) + 33))
	}
	return buf.Bytes()
}
