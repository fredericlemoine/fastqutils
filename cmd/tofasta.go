package cmd

import (
	"bufio"
	stdio "io"
	"log"

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
		var closer1, closer2 stdio.Closer
		var parser *io.FastQParser
		var err error

		if parser, err = openFastqParser(input1, input2); err != nil {
			return
		}
		defer parser.Close()

		// DSRC compresses raw FASTQ streams; it does not apply to FASTA
		// output, so this command only offers --gz (no --dsrc flag).
		if w1, closer1, err = io.GetWriter(output1, gziped, false); err != nil {
			log.Fatal(err)
		}
		if input2 != "none" && output2 != "none" {
			if w2, closer2, err = io.GetWriter(output2, gziped, false); err != nil {
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
	RootCmd.AddCommand(tofastaCmd)

	tofastaCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
	tofastaCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	tofastaCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")
	tofastaCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	tofastaCmd.PersistentFlags().StringVar(&output2, "output2", "none", "Output file 2 (if paired)")
}
