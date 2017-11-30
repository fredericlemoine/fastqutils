package io

import (
	"bufio"
	"compress/gzip"
	"fmt"
	errorp "github.com/fredericlemoine/fastqutils/error"
	"github.com/fredericlemoine/fastqutils/fastq"
	"os"
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

func GetWriter(file string, gz bool) (*bufio.Writer, *gzip.Writer, *os.File) {
	ext := ""
	var fi *os.File
	var w *bufio.Writer
	var err error
	if gz {
		ext = ".gz"
	}

	if file == "stdout" || file == "-" {
		fi = os.Stdout
	} else {
		if fi, err = os.Create(file + ext); err != nil {
			errorp.ExitWithMessage(err)
		}
	}

	var gw *gzip.Writer
	if gz {
		gw = gzip.NewWriter(fi)
		w = bufio.NewWriter(gw)
	} else {
		w = bufio.NewWriter(fi)
	}

	return w, gw, fi
}
