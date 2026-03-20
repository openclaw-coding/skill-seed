package generate

import (
	"context"
	"fmt"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/container"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/infra/output"
	"github.com/spf13/cobra"
)

var (
	outputPath string
)

// Cmd 返回 generate 命令
func Cmd(cont *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-skills",
		Short: i18n.Get("GenerateShort"),
		Long:  i18n.Get("GenerateLongDesc"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 检查 container 是否初始化
			if cont == nil {
				output.Error("%s", i18n.Get("GenerateNotInitialized"))
				output.Dim("%s", i18n.Get("GenerateRunInitFirst")+"\n")
				return fmt.Errorf("skill-seed not initialized")
			}
			return runGenerate(cont, cmd)
		},
	}

	// 添加 flags
	defaultOutputPath := "~/.claude/skills/skill-seed-skills"
	if cont != nil {
		defaultOutputPath = cont.ConfigRepo.GetOutputConfig().SkillsPath
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", defaultOutputPath, "Output path for generated skills")

	return cmd
}

func runGenerate(cont *container.Container, cmd *cobra.Command) error {
	ctx := context.Background()

	output.Info("%s", i18n.Get("GenerateStarting"))

	// 获取模式数量
	count, err := cont.PatternRepo.Count(ctx)
	if err != nil {
		output.Error("%s", i18n.GetWithParams("GenerateCountFailed", map[string]interface{}{"Error": err.Error()}))
		return err
	}

	if count == 0 {
		output.Warning("%s", i18n.Get("GenerateNoPatterns"))
		return nil
	}

	output.Dim("%s", i18n.GetWithParams("GenerateFoundPatterns", map[string]interface{}{"Count": count})+"\n")

	// 生成 Skills
	if err := cont.GeneratorSvc.GenerateSkills(ctx, outputPath); err != nil {
		output.Error("%s", i18n.GetWithParams("GenerateFailed", map[string]interface{}{"Error": err.Error()}))
		return err
	}

	output.Success("%s", i18n.Get("GenerateSuccessMsg"))
	output.Info("%s", i18n.GetWithParams("GenerateOutputPath", map[string]interface{}{"Path": outputPath}))

	return nil
}
