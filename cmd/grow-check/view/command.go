package view

import (
	"fmt"
	"os"

	"github.com/openclaw-coding/grow-check/cmd/grow-check/utils"
	"github.com/openclaw-coding/grow-check/internal/i18n"
	"github.com/openclaw-coding/grow-check/internal/learner"
	"github.com/spf13/cobra"
)

// PatternsCmd 返回查看模式命令
func PatternsCmd() *cobra.Command {
	patternsCmd := &cobra.Command{
		Use:   "patterns",
		Short: "List learned patterns",
		Long: i18n.Get("cmd_view_long") + `

Patterns are extracted from commit history using Claude
and represent common coding practices in your project.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := listPatterns(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return patternsCmd
}

// RulesCmd 返回查看规则命令
func RulesCmd() *cobra.Command {
	rulesCmd := &cobra.Command{
		Use:   "rules",
		Short: "List generated rules",
		Long: i18n.Get("cmd_view_long") + `

Rules are created when patterns reach a certain frequency
and confidence threshold.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := listRules(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return rulesCmd
}

// Cmd 返回查看命令组
func Cmd() *cobra.Command {
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: i18n.Get("cmd_view_short"),
		Long:  i18n.Get("cmd_view_long"),
	}

	viewCmd.AddCommand(PatternsCmd())
	viewCmd.AddCommand(RulesCmd())

	return viewCmd
}

func listPatterns() error {
	skillPath, err := utils.RequireSkillPath()
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
	skillPath, err := utils.RequireSkillPath()
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
