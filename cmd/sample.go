package cmd

import (
	"bufio"
	"compress/gzip"
	"github.com/fredericlemoine/fastqutils/error"
	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
)

var sampleNumber int

// sampleCmd represents the sample command
var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nbrecords := 0

		sampled1 := make([]*fastq.FastqEntry, sampleNumber)
		sampled2 := make([]*fastq.FastqEntry, sampleNumber)

		for {
			entry1, entry2, err := parser.NextEntry()
			if err != nil {
				if err.Error() != "EOF" {
					error.WarnMessage(err)
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

		var w1, w2 *bufio.Writer
		var f1, f2 *os.File
		var g1, g2 *gzip.Writer

		w1, g1, f1 = io.GetWriter(output1, gziped)
		if input2 != "none" {
			w2, g2, f2 = io.GetWriter(output2, gziped)
		}

		for i := 0; i < sampleNumber; i++ {
			entry1 := sampled1[i]
			entry2 := sampled2[i]

			io.WriteEntry(w1, entry1)
			if w2 != nil {
				io.WriteEntry(w2, entry2)
			}
		}
		w1.Flush()
		if g1 != nil {
			g1.Flush()
			g1.Close()
		}
		f1.Close()
		if input2 != "none" {
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
	RootCmd.AddCommand(sampleCmd)

	sampleCmd.PersistentFlags().IntVarP(&sampleNumber, "number", "n", 1, "Number of reads to sample from the FastQ file")
	sampleCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s)")
	sampleCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	sampleCmd.PersistentFlags().StringVar(&output2, "output2", "stdout", "Output file 2 (if paired)")

}
