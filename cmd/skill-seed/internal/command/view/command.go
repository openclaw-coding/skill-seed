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
		Short: i18n.Get("ViewPatternsShort"),
		Long:  i18n.Get("ViewPatternsLong"),
		Run: func(cmd *cobra.Command, args []string) {
			// 检查 container 是否初始化
			if cont == nil {
				fmt.Fprintln(os.Stderr, i18n.Get("ViewNotInitialized"))
				fmt.Fprintln(os.Stderr, i18n.Get("ViewRunInitFirst"))
				os.Exit(1)
			}
			if err := listPatterns(cont); err != nil {
				fmt.Fprintf(os.Stderr, "%s", i18n.GetWithParams("ViewFailed", map[string]interface{}{"Error": err.Error()}))
				os.Exit(1)
			}
		},
	}
}

func rulesCmd(cont *container.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "rules",
		Short: i18n.Get("ViewRulesShort"),
		Long:  i18n.Get("ViewRulesLong"),
		Run: func(cmd *cobra.Command, args []string) {
			// 检查 container 是否初始化
			if cont == nil {
				fmt.Fprintln(os.Stderr, i18n.Get("ViewNotInitialized"))
				fmt.Fprintln(os.Stderr, i18n.Get("ViewRunInitFirst"))
				os.Exit(1)
			}
			if err := listRules(cont); err != nil {
				fmt.Fprintf(os.Stderr, "%s", i18n.GetWithParams("ViewFailed", map[string]interface{}{"Error": err.Error()}))
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
		fmt.Println(i18n.Get("ViewNoPatterns"))
		return nil
	}

	fmt.Println(i18n.GetWithParams("ViewPatternsTotal", map[string]interface{}{"Count": len(patterns)}))
	fmt.Println()
	for i, p := range patterns {
		fmt.Printf("%d. [%s] %s\n", i+1, p.Category, p.Name)
		fmt.Println(i18n.GetWithParams("ViewPatternConfidence", map[string]interface{}{
			"Confidence": fmt.Sprintf("%.2f", p.Confidence),
			"Frequency":  p.Frequency,
		}))
		if p.Description != "" {
			fmt.Println(i18n.GetWithParams("ViewPatternDescription", map[string]interface{}{"Description": p.Description}))
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
		fmt.Println(i18n.Get("ViewNoRules"))
		return nil
	}

	fmt.Println(i18n.GetWithParams("ViewRulesTotal", map[string]interface{}{"Count": len(rules)}))
	fmt.Println()
	for i, r := range rules {
		fmt.Printf("%d. %s [%s]\n", i+1, r.Name, r.Category)
		if r.Description != "" {
			fmt.Println(i18n.GetWithParams("ViewRuleDescription", map[string]interface{}{"Description": r.Description}))
		}
		fmt.Println(i18n.GetWithParams("ViewRulePatterns", map[string]interface{}{"Count": len(r.PatternIDs)}))
		fmt.Println()
	}

	return nil
}

