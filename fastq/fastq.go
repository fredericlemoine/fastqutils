package fastq

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
