package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/grow-check/internal/learner"
	"github.com/spf13/cobra"
)

var (
	sinceDays   int
	maxCommits  int
)

var learnCmd = &cobra.Command{
	Use:   "learn",
	Short: "Learn from Git history",
	Long: `Analyze Git commit history and learn code patterns.

This will:
  - Analyze recent commits (default: last 30 days)
  - Use Claude to identify patterns
  - Generate rules from learned patterns
  - Store everything in the memory database

Examples:
  # Learn from last 30 days
  grow-check learn

  # Learn from last 7 days
  grow-check learn --since=7

  # Learn from last 100 commits
  grow-check learn --max=100`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := learnFromHistory(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Learn failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	learnCmd.Flags().IntVarP(&sinceDays, "since", "s", 30, "Days to look back")
	learnCmd.Flags().IntVarP(&maxCommits, "max", "m", 0, "Maximum commits to analyze (0 = use config)")
	rootCmd.AddCommand(learnCmd)
}

func learnFromHistory() error {
	// 查找 skill 路径
	skillPath, err := findSkillPath()
	if err != nil {
		return err
	}

	// 创建学习器
	learn, err := learner.New(skillPath)
	if err != nil {
		return fmt.Errorf("failed to create learner: %w", err)
	}
	defer learn.Close()

	// 执行学习
	fmt.Printf("🤖 Learning from Git history (last %d days)...\n\n", sinceDays)
	
	if err := learn.LearnFromHistory(sinceDays, maxCommits); err != nil {
		return err
	}

	return nil
}

func findSkillPath() (string, error) {
	// 从当前目录开始查找
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上查找 .skills/grow-check/
	for {
		skillPath := filepath.Join(dir, ".skills", "grow-check")
		if _, err := os.Stat(skillPath); err == nil {
			return skillPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("grow-check not initialized. Run 'grow-check init' first")
}
