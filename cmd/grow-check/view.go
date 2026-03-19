package main

import (
	"fmt"
	"os"

	"github.com/openclaw-coding/grow-check/internal/learner"
	"github.com/spf13/cobra"
)

var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "List learned patterns",
	Long: `Show all code patterns learned from Git history.

Patterns are extracted from commit history using Claude
and represent common coding practices in your project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listPatterns(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List generated rules",
	Long: `Show all rules generated from learned patterns.

Rules are created when patterns reach a certain frequency
and confidence threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listRules(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(patternsCmd)
	rootCmd.AddCommand(rulesCmd)
}

func listPatterns() error {
	skillPath, err := findSkillPath()
	if err != nil {
		return err
	}

	learn, err := learner.New(skillPath)
	if err != nil {
		return err
	}
	defer learn.Close()

	return learn.ListPatterns()
}

func listRules() error {
	skillPath, err := findSkillPath()
	if err != nil {
		return err
	}

	learn, err := learner.New(skillPath)
	if err != nil {
		return err
	}
	defer learn.Close()

	return learn.ListRules()
}
