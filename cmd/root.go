package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-wrap",
	Short: "git-wrap is an EC-compliant Git and Submodule manager",
	Long:  `A fast CLI tool built in Go to automate git workflows, enforce European Commission commit guidelines, and track submodules.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
