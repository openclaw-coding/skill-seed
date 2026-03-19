package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grow-check",
	Short: "Growing Git pre-commit checker",
	Long: `A project-level skill that learns from Git history 
to improve code checking over time.

This tool integrates with Claude for deep code analysis
and learns your team's coding patterns automatically.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
