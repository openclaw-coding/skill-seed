package learn

import (
	"context"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/container"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/infra/output"
	"github.com/spf13/cobra"
)

var (
	limit   int
	all     bool
	verbose bool
)

// Cmd 返回 learn 命令
func Cmd(cont *container.Container) *cobra.Command {
	// 从配置获取默认值
	defaultLimit := cont.ConfigRepo.GetLearningConfig().MaxCommits

	cmd := &cobra.Command{
		Use:   "learn",
		Short: i18n.Get("LearnShort"),
		Long:  i18n.Get("LearnLongDesc"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLearn(cont, cmd)
		},
	}

	// 添加 flags
	cmd.Flags().IntVarP(&limit, "limit", "l", defaultLimit, "Number of commits to analyze")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Analyze all commits")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func runLearn(cont *container.Container, cmd *cobra.Command) error {
	ctx := context.Background()

	// 检查是否分析所有提交
	if all {
		limit = 1000 // 设置一个较大的值
	}

	output.Info("%s", i18n.Get("LearnStarting"))
	output.Dim("%s", i18n.GetWithParams("LearnAnalyzingCommitsCount", map[string]interface{}{"Count": limit})+"\n")

	// 调用学习服务
	err := cont.LearnerSvc.Learn(ctx, limit)
	if err != nil {
		output.Error("%s", i18n.GetWithParams("LearnFailed", map[string]interface{}{"Error": err.Error()}))
		return err
	}

	output.Success("%s", i18n.Get("LearnCompleted"))

	// 显示统计信息
	count, err := cont.PatternRepo.Count(ctx)
	if err == nil {
		output.Info("%s", i18n.GetWithParams("LearnTotalPatterns", map[string]interface{}{"Count": count}))
	}

	return nil
}
