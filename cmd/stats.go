package cmd

import (
	"fmt"
	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/stats"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
