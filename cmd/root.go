package cmd

import (
	"fmt"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var cfgFile string

var Version string = "Unknown"

var input1 string
var input2 string
var parser *io.FastQParser
var seed int64

var RootCmd = &cobra.Command{
	Use:   "fastqutils",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rand.Seed(seed)
		if input2 != "none" {
			parser = io.NewPairedEndParser(input1, input2)
		} else {
			parser = io.NewSingleEndParser(input1)
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&input1, "input-1", "1", "stdin", "First read fastq file")
	RootCmd.PersistentFlags().StringVarP(&input2, "input-2", "2", "none", "Second read fastq file")
	RootCmd.PersistentFlags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
