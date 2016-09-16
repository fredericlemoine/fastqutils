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
}

func ComputeStats(parser *io.FastQParser) Stats {
	nbrecords := 0
	paired := true
	totalNt := make([]int64, 4)
	freqNt := make([]float64, 4)
	var total int64 = 0
	var meanQual float64
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
			total++
			if entry2 != nil {
				totalNt[fastq.Index(entry2.Sequence[i])]++
				meanQual += float64(int(entry2.Quality[i]))
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

	return Stats{
		nbrecords,
		paired,
		freqNt,
		meanQual / float64(total),
	}
}
