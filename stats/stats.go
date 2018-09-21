package stats

import (
	"github.com/fredericlemoine/fastqutils/error"
	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/io"
)

type Stats struct {
	NSeq     int       // Number of sequences
	Paired   bool      // If the Fastq are paired end
	TotalNt  []float64 // global % of A / C / G / T
	MeanQual float64   // Average base quality
	MinQual  int       // Min quality score
	MaxQual  int       // Max quality score
	Encoding int       // Quality encoding
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ComputeStats(parser *io.FastQParser) Stats {
	nbrecords := 0
	paired := true
	totalNt := make([]int64, 5)
	freqNt := make([]float64, 5)
	var total int64 = 0
	var meanQual float64
	minqual, maxqual := 1000, 0

	for {
		entry1, entry2, err := parser.NextEntry()
		if err != nil {
			if err.Error() != "EOF" {
				error.WarnMessage(err)
			}
			break
		}

		for i := 0; i < len(entry1.Sequence); i++ {
			totalNt[fastq.Index(entry1.Sequence[i])]++
			meanQual += float64(int(entry1.Quality[i]))
			minqual = min(minqual, int(entry1.Quality[i]))
			maxqual = max(maxqual, int(entry1.Quality[i]))
			total++
		}
		if entry2 != nil {
			for i := 0; i < len(entry2.Sequence); i++ {
				totalNt[fastq.Index(entry2.Sequence[i])]++
				meanQual += float64(int(entry2.Quality[i]))
				minqual = min(minqual, int(entry2.Quality[i]))
				maxqual = max(maxqual, int(entry2.Quality[i]))
				total++
			}
		}

		if entry2 == nil {
			paired = false
		}
		nbrecords++
	}
	freqNt = make([]float64, len(totalNt))
	for i, v := range totalNt {
		freqNt[i] = float64(v) / float64(total)
	}

	encoding := DetectEncoding(minqual, maxqual)

	return Stats{
		nbrecords,
		paired,
		freqNt,
		meanQual/float64(total) - float64(EncodingOffset(encoding)),
		minqual - EncodingOffset(encoding),
		maxqual - EncodingOffset(encoding),
		encoding,
	}
}
