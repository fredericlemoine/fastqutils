package cmd

import (
	"bufio"
	stdio "io"
	"log"
	"math/rand"

	"github.com/spf13/cobra"

	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/io"
)

var sampleNumber int

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// sampleCmd represents the sample command
var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Subsample a FastQ File",
	Long:  `Subsample a FastQ File`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var parser *io.FastQParser

		nbrecords := 0

		sampled1 := make([]*fastq.FastqEntry, sampleNumber)
		sampled2 := make([]*fastq.FastqEntry, sampleNumber)

		if parser, err = openFastqParser(input1, input2); err != nil {
			log.Fatal(err)
		}
		defer parser.Close()

		for {
			entry1, entry2, err := parser.NextEntry()
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				}
				break
			}

			if nbrecords < sampleNumber {
				sampled1[nbrecords] = entry1
				if entry2 != nil {
					sampled2[nbrecords] = entry2
				}
			} else {
				random := rand.Intn(nbrecords)
				if random < sampleNumber {
					sampled1[random] = entry1
					if entry2 != nil {
						sampled2[random] = entry2
					}
				}
			}
			nbrecords++
		}

		if nbrecords < sampleNumber {
			log.Printf("fastq file length (%d) is < sampling number (%d) , will write only %d reads", nbrecords, sampleNumber, nbrecords)
		}

		var w1, w2 *bufio.Writer
		var closer1, closer2 stdio.Closer

		if w1, closer1, err = io.GetWriter(output1, gziped, dsrcOut); err != nil {
			log.Fatal(err)
		}
		if input2 != "none" && output2 != "none" {
			if w2, closer2, err = io.GetWriter(output2, gziped, dsrcOut); err != nil {
				log.Fatal(err)
			}
		}

		for i := 0; i < min(sampleNumber, nbrecords); i++ {
			entry1 := sampled1[i]
			entry2 := sampled2[i]

			io.WriteEntry(w1, entry1)
			if w2 != nil {
				io.WriteEntry(w2, entry2)
			}
		}
		if err = closer1.Close(); err != nil {
			log.Fatal(err)
		}
		if input2 != "none" && output2 != "none" {
			if err = closer2.Close(); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(sampleCmd)
	sampleCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	sampleCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")
	sampleCmd.PersistentFlags().IntVarP(&sampleNumber, "number", "n", 1, "Number of reads to sample from the FastQ file")
	sampleCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
	sampleCmd.PersistentFlags().BoolVar(&dsrcOut, "dsrc", false, "If true, will generate dsrc-compressed file(s) : .dsrc extension is added automatically (requires the 'dsrc' executable, see dsrc/README.md)")
	sampleCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	sampleCmd.PersistentFlags().StringVar(&output2, "output2", "none", "Output file 2 (if paired)")
}
