package main

import (
	"fmt"
	"os"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/initialize"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/learn"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/check"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/analyze"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/view"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/generate"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/hook"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/scan"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skill-seed",
	Short: "Growing project skills for AI agents",
	Long: `A project-level skill that learns from Git history
to help AI agents understand your codebase better.

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
	rootCmd.AddCommand(hook.Cmd())
	rootCmd.AddCommand(scan.Cmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
