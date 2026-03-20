package generate

import (
	"context"

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
	// 从配置获取默认值
	defaultOutputPath := cont.ConfigRepo.GetOutputConfig().SkillsPath

	cmd := &cobra.Command{
		Use:   "generate-skills",
		Short: i18n.Get("GenerateShort"),
		Long:  i18n.Get("GenerateLongDesc"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cont, cmd)
		},
	}

	// 添加 flags
	cmd.Flags().StringVarP(&outputPath, "output", "o", defaultOutputPath, "Output path for generated skills")

	return cmd
}

func runGenerate(cont *container.Container, cmd *cobra.Command) error {
	ctx := context.Background()

	output.Info("Generating Claude Code skills...")

	// 获取模式数量
	count, err := cont.PatternRepo.Count(ctx)
	if err != nil {
		output.Error("Failed to count patterns: %v", err)
		return err
	}

	if count == 0 {
		output.Warning("No patterns learned yet. Run 'skill-seed learn' first.")
		return nil
	}

	output.Dim("Found %d patterns\n", count)

	// 生成 Skills
	if err := cont.GeneratorSvc.GenerateSkills(ctx, outputPath); err != nil {
		output.Error("Failed to generate skills: %v", err)
		return err
	}

	output.Success("✓ Skills generated successfully!")
	output.Info("Output: %s", outputPath)

	return nil
}
