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
		var strenc string

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
			fmt.Printf("%.2f\n", v)
		}
		if strenc, err = stats.EncodingToString(stat.Encoding); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Encoding\t%s\n", strenc)
		fmt.Printf("AvgQual\t%.3f\n", stat.MeanQual)
		fmt.Printf("MinQual\t%d\n", stat.MinQual)
		fmt.Printf("MaxQual\t%d\n", stat.MaxQual)
		fmt.Printf("Quality Histogram\n%s\n", stat.QualHistogram.Draw(100))
		fmt.Printf("Length Histogram\n%s\n", stat.LenHistogram.Draw(100))
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	statsCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")

}
