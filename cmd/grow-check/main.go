package main

import (
	"fmt"
	"os"

	"github.com/openclaw-coding/grow-check/cmd/grow-check/initialize"
	"github.com/openclaw-coding/grow-check/cmd/grow-check/learn"
	"github.com/openclaw-coding/grow-check/cmd/grow-check/check"
	"github.com/openclaw-coding/grow-check/cmd/grow-check/analyze"
	"github.com/openclaw-coding/grow-check/cmd/grow-check/view"
	"github.com/openclaw-coding/grow-check/cmd/grow-check/generate"
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
	// 添加子命令
	rootCmd.AddCommand(initialize.Cmd())
	rootCmd.AddCommand(learn.Cmd())
	rootCmd.AddCommand(check.Cmd())
	rootCmd.AddCommand(analyze.Cmd())
	rootCmd.AddCommand(view.Cmd())
	rootCmd.AddCommand(generate.Cmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
