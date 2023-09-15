package cmd

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"

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
		var f1, f2 *os.File
		var g1, g2 *gzip.Writer

		if w1, g1, f1, err = io.GetWriter(output1, gziped); err != nil {
			log.Fatal(err)
		}
		if w2, g2, f2, err = io.GetWriter(output2, gziped); err != nil {
			log.Fatal(err)
		}

		if parser, err = openFastqParser(input1, "none"); err != nil {
			log.Fatal(err)
		}

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

		w1.Flush()
		if g1 != nil {
			g1.Flush()
			g1.Close()
		}
		f1.Close()
		w2.Flush()
		if g2 != nil {
			g2.Flush()
			g2.Close()
		}
		f2.Close()
	},
}

func init() {
	RootCmd.AddCommand(deinterlaceCmd)
	deinterlaceCmd.PersistentFlags().StringVarP(&input1, "input", "i", "stdin", "First read fastq file")
	deinterlaceCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Deinterlaced Output file R1")
	deinterlaceCmd.PersistentFlags().StringVar(&output2, "output2", "stdout", "Deinterlaced Output file R2")
	deinterlaceCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
}
