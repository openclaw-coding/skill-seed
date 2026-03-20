package check

import (
	"context"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/container"
	"github.com/openclaw-coding/skill-seed/internal/domain"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/infra/output"
	"github.com/spf13/cobra"
)

var (
	interactive bool
	checkAll    bool
	verbose     bool
)

// Cmd 返回 check 命令
func Cmd(cont *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: i18n.Get("CheckShort"),
		Long:  i18n.Get("CheckLongDesc"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(cont, cmd)
		},
	}

	// 添加 flags
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "Interactive mode")
	cmd.Flags().BoolVarP(&checkAll, "all", "a", false, "Check all files (not just staged)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func runCheck(cont *container.Container, cmd *cobra.Command) error {
	ctx := context.Background()

	output.Info("%s", i18n.Get("CheckStarting"))

	var issues []domain.Issue
	var err error

	// 检查所有文件还是只检查暂存文件
	if checkAll {
		output.Dim("%s", i18n.Get("CheckAllFiles")+"\n")
		issues, err = cont.CheckerSvc.CheckAll(ctx)
	} else {
		output.Dim("%s", i18n.Get("CheckStagedFiles")+"\n")
		issues, err = cont.CheckerSvc.Check(ctx)
	}

	if err != nil {
		output.Error("%s", i18n.GetWithParams("CheckFailed", map[string]interface{}{"Error": err.Error()}))
		return err
	}

	// 显示检查结果
	if len(issues) == 0 {
		output.Success("%s", i18n.Get("CheckNoIssues"))
		return nil
	}

	output.Warning("%s", i18n.GetWithParams("CheckFoundIssues", map[string]interface{}{"Count": len(issues)})+"\n")
	for i, iss := range issues {
		output.Print("\n%d. ", i+1)
		output.Print("%s ", output.SeverityLabel(string(iss.Severity)))
		output.Println("%s:%d", iss.File, iss.Line)
		output.Println("   %s", iss.Message)

		if iss.Suggestion != "" {
			output.Dim("%s", i18n.GetWithParams("CheckSuggestion", map[string]interface{}{"Suggestion": iss.Suggestion})+"\n")
		}
	}

	// 如果是交互模式，处理问题
	if interactive {
		return handleIssuesInteractively(issues)
	}

	return nil
}

// handleIssuesInteractively 交互式处理问题
func handleIssuesInteractively(issues []domain.Issue) error {
	// TODO: 实现交互式处理逻辑
	// - 显示每个问题的详情
	// - 提供选项：自动修复、查看详情、忽略、终止
	// - 根据用户选择执行相应操作
	return nil
}
