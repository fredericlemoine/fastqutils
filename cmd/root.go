package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fredericlemoine/fastqutils/io"
	"github.com/spf13/cobra"
)

// Version stores tool version
var Version string = "Unknown"

var input1 string
var input2 string
var seed int64

// RootCmd represents the root Command
var RootCmd = &cobra.Command{
	Use:   "fastqutils",
	Short: "Some tools to handle fastq files",
	Long: `Some tools to handle fastqfiles.

For now:
sample: to take a subset of a whole fastqfile in one pass
stats : a few statistics about the fastq file

Works for single and paired end files.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rand.Seed(seed)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute executes the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func openFastqParser(input1, input2 string) (fp *io.FastQParser, err error) {
	if input2 != "none" {
		fp, err = io.NewPairedEndParser(input1, input2)
	} else {
		fp, err = io.NewSingleEndParser(input1)
	}
	return
}

func init() {
	RootCmd.PersistentFlags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
