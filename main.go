package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vx416/smart-git/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "smart-git",
		Short: "Smart Git is a tool to help you with git operations using AI.",
	}
	rootCmd.AddCommand(cmd.SummaryGitDiff)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
