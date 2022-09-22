/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// maskCmd represents the mask command
var maskCmd = &cobra.Command{
	Use:   "mask",
	Short: "Commands to mask nucleotides from bam or fastq files",
	Long:  `Commands to mask nucleotides from bam or fastq files.`,
}

func init() {
	RootCmd.AddCommand(maskCmd)
}
