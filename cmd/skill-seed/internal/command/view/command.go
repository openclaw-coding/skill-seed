package view

import (
	"context"
	"fmt"
	"os"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/container"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/spf13/cobra"
)

// Cmd 返回 view 命令
func Cmd(cont *container.Container) *cobra.Command {
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: i18n.Get("ViewShort"),
		Long:  i18n.Get("ViewLongDesc"),
	}

	viewCmd.AddCommand(patternsCmd(cont))
	viewCmd.AddCommand(rulesCmd(cont))

	return viewCmd
}

func patternsCmd(cont *container.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "patterns",
		Short: "View all learned patterns",
		Long:  "View all learned patterns from Git history.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := listPatterns(cont); err != nil {
				fmt.Fprintf(os.Stderr, "Failed: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func rulesCmd(cont *container.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "rules",
		Short: "View all generated rules",
		Long:  "View all generated rules from patterns.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := listRules(cont); err != nil {
				fmt.Fprintf(os.Stderr, "Failed: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func listPatterns(cont *container.Container) error {
	ctx := context.Background()

	patterns, err := cont.GetPatternRepository().GetAll(ctx)
	if err != nil {
		return err
	}

	if len(patterns) == 0 {
		fmt.Println("No patterns learned yet.")
		return nil
	}

	fmt.Printf("Learned patterns (%d total):\n\n", len(patterns))
	for i, p := range patterns {
		fmt.Printf("%d. [%s] %s\n", i+1, p.Category, p.Name)
		fmt.Printf("   Confidence: %.2f | Frequency: %d\n", p.Confidence, p.Frequency)
		if p.Description != "" {
			fmt.Printf("   Description: %s\n", p.Description)
		}
		fmt.Println()
	}

	return nil
}

func listRules(cont *container.Container) error {
	ctx := context.Background()

	rules, err := cont.RuleRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		fmt.Println("No rules generated yet.")
		return nil
	}

	fmt.Printf("Generated rules (%d total):\n\n", len(rules))
	for i, r := range rules {
		fmt.Printf("%d. %s [%s]\n", i+1, r.Name, r.Category)
		if r.Description != "" {
			fmt.Printf("   Description: %s\n", r.Description)
		}
		fmt.Printf("   Patterns: %d\n", len(r.PatternIDs))
		fmt.Println()
	}

	return nil
}

