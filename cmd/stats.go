package cmd

import (
	"fmt"
	"log"

	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/fredericlemoine/fastqutils/stats"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Displays different statistics about fastq file(s)",
	Long:  `Displays different statistics about fastq file(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		var parser *io.FastQParser
		var err error
		var nt byte
		var stat stats.Stats

		if parser, err = openFastqParser(input1, input2); err != nil {
			return
		}
		if stat, err = stats.ComputeStats(parser); err != nil {
			log.Fatal(err)
		}
		fmt.Print("NSeq\t")
		fmt.Println(stat.NSeq)
		fmt.Print("Paired\t")
		fmt.Println(stat.Paired)
		for i, v := range stat.TotalNt {
			nt, _ = fastq.Nt(i)
			fmt.Printf("%c", nt)
			fmt.Print("\t")
			fmt.Println(v)
		}
		fmt.Print("Encoding\t")
		fmt.Println(stats.EncodingToString(stat.Encoding))
		fmt.Print("AvgQual\t")
		fmt.Println(stat.MeanQual)
		fmt.Print("MinQual\t")
		fmt.Println(stat.MinQual)
		fmt.Print("MaxQual\t")
		fmt.Println(stat.MaxQual)
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	statsCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")

}
