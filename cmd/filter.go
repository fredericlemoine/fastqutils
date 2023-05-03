/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// maskCmd represents the mask command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Commands to filter reads",
	Long:  `Commends to filter reads.`,
}

func init() {
	RootCmd.AddCommand(filterCmd)
}
