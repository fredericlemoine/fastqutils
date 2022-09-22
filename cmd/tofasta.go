package cmd

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/fredericlemoine/fastqutils/io"
)

// tofastaCmd represents the tofasta command
var tofastaCmd = &cobra.Command{
	Use:   "tofasta",
	Short: "Converts input fastq file into fasta",
	Long:  `Converts input fastq file into fasta.`,
	Run: func(cmd *cobra.Command, args []string) {
		var w1, w2 *bufio.Writer
		var f1, f2 *os.File
		var g1, g2 *gzip.Writer
		var parser *io.FastQParser
		var err error

		if parser, err = openFastqParser(input1, input2); err != nil {
			return
		}

		if w1, g1, f1, err = io.GetWriter(output1, gziped); err != nil {
			log.Fatal(err)
		}
		if input2 != "none" && output2 != "none" {
			if w2, g2, f2, err = io.GetWriter(output2, gziped); err != nil {
				log.Fatal(err)
			}
		}

		for {
			entry1, entry2, err := parser.NextEntry()
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				}
				break
			}
			io.WriteEntryFasta(w1, entry1)
			if w2 != nil {
				io.WriteEntryFasta(w2, entry2)
			}
		}

		w1.Flush()
		if g1 != nil {
			g1.Flush()
			g1.Close()
		}
		f1.Close()
		if input2 != "none" && output2 != "none" {
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
	RootCmd.AddCommand(tofastaCmd)

	tofastaCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
	tofastaCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	tofastaCmd.PersistentFlags().StringVar(&output2, "output2", "none", "Output file 2 (if paired)")
}
