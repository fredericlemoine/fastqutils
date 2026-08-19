package cmd

import (
	"bufio"
	stdio "io"
	"log"

	"github.com/fredericlemoine/fastqutils/io"
	"github.com/spf13/cobra"
)

// tobamCmd represents the tobam command
var deinterlaceCmd = &cobra.Command{
	Use:   "deinterlace",
	Short: "Place the first reads on file 1 and second reads on file 2",
	Long: `Place the first reads on file 1 and second reads on file 2
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var parser *io.FastQParser
		var w1, w2 *bufio.Writer
		var closer1, closer2 stdio.Closer

		if w1, closer1, err = io.GetWriter(output1, gziped, dsrcOut); err != nil {
			log.Fatal(err)
		}
		if w2, closer2, err = io.GetWriter(output2, gziped, dsrcOut); err != nil {
			log.Fatal(err)
		}

		if parser, err = openFastqParser(input1, "none"); err != nil {
			log.Fatal(err)
		}
		defer parser.Close()

		reads := 0
		for {
			entry1, _, err := parser.NextEntry()
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				}
				break
			}

			if reads%2 == 0 {
				io.WriteEntry(w1, entry1)
			} else {
				io.WriteEntry(w2, entry1)
			}
			reads++
		}

		if err = closer1.Close(); err != nil {
			log.Fatal(err)
		}
		if err = closer2.Close(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deinterlaceCmd)
	deinterlaceCmd.PersistentFlags().StringVarP(&input1, "input", "i", "stdin", "First read fastq file")
	deinterlaceCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Deinterlaced Output file R1")
	deinterlaceCmd.PersistentFlags().StringVar(&output2, "output2", "stdout", "Deinterlaced Output file R2")
	deinterlaceCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
	deinterlaceCmd.PersistentFlags().BoolVar(&dsrcOut, "dsrc", false, "If true, will generate dsrc-compressed file(s) : .dsrc extension is added automatically (requires the 'dsrc' executable, see dsrc/README.md)")
}
