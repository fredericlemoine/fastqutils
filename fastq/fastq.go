package fastq

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/fredericlemoine/gostats"
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
	var nt byte

	for i := 0; i < length; i++ {
		nt, _ = Nt(rand.Intn(4))
		buf.WriteByte(nt)
	}
	return buf.Bytes()
}

// Returns the nt
func Nt(n int) (nt byte, err error) {
	switch n {
	case 0:
		nt = 'A'
	case 1:
		nt = 'C'
	case 2:
		nt = 'G'
	case 3:
		nt = 'T'
	case 4:
		nt = 'N'
	default:
		err = fmt.Errorf("No nucleotide with code %d", n)
	}
	return
}

// Returns the nt
func Index(b byte) (nt int, err error) {
	switch b {
	case 'A':
		nt = 0
	case 'C':
		nt = 1
	case 'G':
		nt = 2
	case 'T':
		nt = 3
	case 'N':
		nt = 4
	default:
		err = fmt.Errorf("No nucleotide %c", b)
	}
	return
}

func genqual(length int, minqual, maxqual int) []byte {
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		buf.WriteByte(byte(gostats.Binomial(float64(length-i)/float64(length)*0.99, (maxqual-minqual)) + minqual))
	}
	return buf.Bytes()
}
