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
		stats := stats.ComputeStats(parser)
		fmt.Print("NSeq\t")
		fmt.Println(stats.NSeq)
		fmt.Print("Paired\t")
		fmt.Println(stats.Paired)
		for i, v := range stats.TotalNt {
			fmt.Print(fmt.Sprintf("%c", fastq.Nt(i)))
			fmt.Print("\t")
			fmt.Println(v)
		}
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
}
