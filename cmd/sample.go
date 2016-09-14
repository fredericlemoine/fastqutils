package cmd

import (
	"fmt"
	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/spf13/cobra"
	"math/rand"
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
					io.WarnMessage(err)
				}
				break
			}

			if nbrecords < sampleNumber {
				sampled1[nbrecords] = entry1
				if entry2 != nil {
					sampled2[nbrecords] = entry2
				}
			} else {
				random := rand.Intn(nbrecords - 1)
				if random < sampleNumber {
					sampled1[random] = entry1
					if entry2 != nil {
						sampled2[random] = entry2
					}
				}
			}
			nbrecords++
		}

		for i := 0; i < sampleNumber; i++ {
			entry1 := sampled1[i]
			entry2 := sampled2[i]

			fmt.Print(entry1.Name)
			if entry2 != nil {
				fmt.Print("\t" + entry2.Name)
			}
			fmt.Println()

			fmt.Print(string(entry1.Sequence))
			if entry2 != nil {
				fmt.Print("\t" + string(entry2.Sequence))
			}
			fmt.Println()

			fmt.Print(string(entry1.Quality))
			if entry2 != nil {
				fmt.Print("\t" + string(entry2.Quality))
			}
			fmt.Println()
		}
	},
}

func init() {
	RootCmd.AddCommand(sampleCmd)

	sampleCmd.PersistentFlags().IntVarP(&sampleNumber, "number", "n", 1, "Number of reads to sample from the FastQ file")

}
