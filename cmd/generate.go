package cmd

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"

	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/fredericlemoine/fastqutils/stats"
	"github.com/spf13/cobra"
)

var paired bool
var gziped bool
var length int
var nbseqs int
var output1, output2 string
var encoding string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates A random Fastq file",
	Long: `Generates a random Fastq file / single or paired end

It does not follow any specific model. 

It draws uniformly nucleotides from A,C,G,T, and qualities depending on the encoding.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var w1, w2 *bufio.Writer
		var f1, f2 *os.File
		var g1, g2 *gzip.Writer
		var qualenc int
		var minqual, maxqual int
		var err error

		if qualenc, err = stats.EncodingFromString(encoding); err != nil {
			log.Fatal(err)
		}
		if minqual, err = stats.MinQual(qualenc); err != nil {
			log.Fatal(err)
		}
		if maxqual, err = stats.MaxQual(qualenc); err != nil {
			log.Fatal(err)
		}
		if w1, g1, f1, err = io.GetWriter(output1, gziped); err != nil {
			log.Fatal(err)
		}
		if paired && output2 != "none" {
			if w2, g2, f2, err = io.GetWriter(output2, gziped); err != nil {
				log.Fatal(err)
			}
		}

		for i := 0; i < nbseqs; i++ {
			entry1 := fastq.GenFastQEntry(length, i, minqual, maxqual)
			io.WriteEntry(w1, entry1)
			if paired && w2 != nil {
				entry2 := fastq.GenFastQEntry(length, i, minqual, maxqual)
				io.WriteEntry(w2, entry2)
			}
		}
		w1.Flush()
		if g1 != nil {
			g1.Flush()
			g1.Close()
		}
		f1.Close()
		if paired && w2 != nil {
			w2.Flush()
			if g2 != nil {
				g2.Flush()
				g2.Close()
			}
			f2.Close()
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().BoolVarP(&paired, "paired", "p", false, "If true, will generate two files")
	generateCmd.PersistentFlags().IntVarP(&length, "length", "l", 100, "Defines the length of generated sequences")
	generateCmd.PersistentFlags().IntVarP(&nbseqs, "nbseqs", "n", 1000, "Defines the number of sequences to generate")
	generateCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s)")
	generateCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	generateCmd.PersistentFlags().StringVar(&output2, "output2", "stdout", "Output file 2 (if paired)")
	generateCmd.PersistentFlags().StringVar(&encoding, "encoding", "illumina1.8", "Base quality encoding")
}
