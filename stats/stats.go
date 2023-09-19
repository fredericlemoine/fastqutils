package stats

import (
	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/hist"
	"github.com/fredericlemoine/fastqutils/io"
)

type Stats struct {
	NSeq          int       // Number of sequences
	Paired        bool      // If the Fastq are paired end
	TotalNt       []float64 // global % of A / C / G / T
	MeanQual      float64   // Average base quality
	MinQual       int       // Min quality score
	MaxQual       int       // Max quality score
	Encoding      int       // Quality encoding
	QualHistogram *hist.IntHistogram
	LenHistogram  *hist.IntHistogram
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

func ComputeStats(parser *io.FastQParser, histos bool) (s Stats, err error) {
	var nbrecords int = 0
	var paired bool = true
	var totalNt []int64 = make([]int64, 5)
	var freqNt []float64
	var total int64 = 0
	var meanQual float64
	var minqual, maxqual int = 1000, 0
	var nt int
	var entry1, entry2 *fastq.FastqEntry
	var qualHistogram, lenHistogram *hist.IntHistogram

	if histos {
		qualHistogram = hist.NewIntHistogram(30)
		lenHistogram = hist.NewIntHistogram(20)
	}

	for {
		entry1, entry2, err = parser.NextEntry()
		if err != nil {
			if err.Error() != "EOF" {
				return
			}
			err = nil
			break
		}

		if histos {
			lenHistogram.AddPoint(int(len(entry1.Sequence)))
		}
		for i := 0; i < len(entry1.Sequence); i++ {
			nt, _ = fastq.Index(entry1.Sequence[i])
			totalNt[nt]++
			meanQual += float64(int(entry1.Quality[i]))
			minqual = min(minqual, int(entry1.Quality[i]))
			maxqual = max(maxqual, int(entry1.Quality[i]))
			if histos {
				qualHistogram.AddPoint(int(entry1.Quality[i]))
			}
			total++
		}
		if entry2 != nil {
			if histos {
				lenHistogram.AddPoint(int(len(entry2.Sequence)))
			}
			for i := 0; i < len(entry2.Sequence); i++ {
				nt, _ = fastq.Index(entry2.Sequence[i])
				totalNt[nt]++
				meanQual += float64(int(entry2.Quality[i]))
				minqual = min(minqual, int(entry2.Quality[i]))
				maxqual = max(maxqual, int(entry2.Quality[i]))
				if histos {
					qualHistogram.AddPoint(int(entry2.Quality[i]))
				}
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

	off, _ := EncodingOffset(encoding)

	s = Stats{
		nbrecords,
		paired,
		freqNt,
		meanQual/float64(total) - float64(off),
		minqual - off,
		maxqual - off,
		encoding,
		qualHistogram,
		lenHistogram,
	}

	return
}
