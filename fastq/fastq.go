package fastq

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fredericlemoine/fastqutils/error"
	"github.com/fredericlemoine/gostats"
	"math/rand"
)

type FastqEntry struct {
	Name     []byte
	Sequence []byte
	Quality  []byte
}

/* Generates a Fastq Entry */
func GenFastQEntry(length int, id int, minqual, maxqual int) *FastqEntry {
	name := []byte(fmt.Sprintf("@read%d", id))
	seq := genseq(length)
	qual := genqual(length, minqual, maxqual)
	return &FastqEntry{
		name,
		seq,
		qual,
	}
}

func genseq(length int) []byte {
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		buf.WriteByte(Nt(rand.Intn(4)))
	}
	return buf.Bytes()
}

// Returns the nt
func Nt(n int) byte {
	switch n {
	case 0:
		return 'A'
	case 1:
		return 'C'
	case 2:
		return 'G'
	case 3:
		return 'T'
	case 4:
		return 'N'
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("No nucleotide with code %d", n)))
	}
	return '\n'
}

// Returns the nt
func Index(b byte) int {
	switch b {
	case 'A':
		return 0
	case 'C':
		return 1
	case 'G':
		return 2
	case 'T':
		return 3
	case 'N':
		return 4
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("No nucleotide %c", b)))
	}
	return '\n'
}

func genqual(length int, minqual, maxqual int) []byte {
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		buf.WriteByte(byte(gostats.Binomial(float64(length-i)/float64(length)*0.99, (maxqual-minqual)) + minqual))
	}
	return buf.Bytes()
}
