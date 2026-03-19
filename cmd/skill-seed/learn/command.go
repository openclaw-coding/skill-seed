package learn

import (
	"fmt"
	"os"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/utils"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/learner"
	"github.com/spf13/cobra"
)

var (
	sinceDays   int
	maxCommits  int
	forceLearn  bool
)

// Cmd 返回学习命令
func Cmd() *cobra.Command {
	learnCmd := &cobra.Command{
		Use:   "learn",
		Short: i18n.Get("cmd_learn_short"),
		Long:  i18n.Get("cmd_learn_long") + `

This will:
  - Analyze recent commits (default: last 30 days)
  - Use Claude to identify patterns
  - Generate rules from learned patterns
  - Store everything in the memory database

Examples:
  # Learn from last 30 days
  skill-seed learn

  # Learn from last 7 days
  skill-seed learn --since=7

  # Learn from last 100 commits
  skill-seed learn --max=100

  # Force re-learn all commits (ignore learned status)
  skill-seed learn --max=100 --force`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := learnFromHistory(); err != nil {
				fmt.Println(i18n.Get("learn_failed"), err)
				os.Exit(1)
			}
		},
	}

	learnCmd.Flags().IntVarP(&sinceDays, "since", "s", 30, "Days to look back")
	learnCmd.Flags().IntVarP(&maxCommits, "max", "m", 0, "Maximum commits to analyze (0 = use config)")
	learnCmd.Flags().BoolVarP(&forceLearn, "force", "f", false, "Force re-learn all commits (ignore learned status)")

	return learnCmd
}

func learnFromHistory() error {
	// Find skill path
	skillPath, err := utils.RequireSkillPath()
	if err != nil {
		return err
	}

	// Create learner
	learn, err := learner.New(skillPath)
	if err != nil {
		return fmt.Errorf("failed to create learner: %w", err)
	}
	defer learn.Close()

	// Execute learning (output message is handled by LearnFromHistory)
	if err := learn.LearnFromHistory(sinceDays, maxCommits, forceLearn); err != nil {
		return err
	}

	return nil
}
