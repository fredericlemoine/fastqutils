package io

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/fredericlemoine/fastqutils/fastq"
)

func WriteEntry(w *bufio.Writer, entry *fastq.FastqEntry) {
	w.WriteString(
		fmt.Sprintf("%s\n%s\n+\n%s\n",
			entry.Name,
			entry.Sequence,
			entry.Quality,
		),
	)
}

func WriteEntryFasta(w *bufio.Writer, entry *fastq.FastqEntry) {
	w.WriteString(fmt.Sprintf(">%s\n%s\n", entry.Name, entry.Sequence))
}

func GetWriter(file string, gz bool) (w *bufio.Writer, gw *gzip.Writer, fi *os.File, err error) {
	ext := ""
	if gz {
		ext = ".gz"
	}

	if file == "stdout" || file == "-" {
		fi = os.Stdout
	} else {
		if fi, err = os.Create(file + ext); err != nil {
			return
		}
	}

	if gz {
		gw = gzip.NewWriter(fi)
		w = bufio.NewWriter(gw)
	} else {
		w = bufio.NewWriter(fi)
	}

	return
}
