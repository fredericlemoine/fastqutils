package cmd

import (
	"fmt"
	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/stats"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Displays different statistics about fastq file(s)",
	Long:  `Displays different statistics about fastq file(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		stat := stats.ComputeStats(parser)
		fmt.Print("NSeq\t")
		fmt.Println(stat.NSeq)
		fmt.Print("Paired\t")
		fmt.Println(stat.Paired)
		for i, v := range stat.TotalNt {
			fmt.Print(fmt.Sprintf("%c", fastq.Nt(i)))
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
}
