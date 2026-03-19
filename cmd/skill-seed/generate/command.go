package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/skill-seed/internal/git"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/generator"
	"github.com/spf13/cobra"
)

var (
	outputPath string
	force      bool
)

// Cmd 返回生成命令
func Cmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate-skills",
		Short: i18n.Get("cmd_generate_short"),
		Long:  i18n.Get("cmd_generate_long") + `

This command reads the patterns learned from Git history and generates
a skill that can be used by Claude Code to understand your project's
coding conventions and best practices.

The generated skill will be output to ~/.claude/skills/skill-seed-skills/
by default, where Claude Code can automatically discover it.`,
		RunE: runGenerate,
	}

	generateCmd.Flags().StringVarP(&outputPath, "output", "o",
		filepath.Join(os.Getenv("HOME"), ".claude/skills"), "Output path for generated skills")
	generateCmd.Flags().BoolVarP(&force, "force", "f", false,
		"Overwrite existing skills without asking")

	return generateCmd
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// 确定项目根目录（查找 .git 目录）
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 检查是否在 Git 仓库中
	gitOp := git.NewGitOperator(projectRoot)
	if !gitOp.IsGitRepo() {
		return fmt.Errorf("not a git repository: %s", projectRoot)
	}

	// 确定 skill 路径
	skillPath := filepath.Join(projectRoot, ".seed", "skill-seed")

	// 检查 skill 是否初始化
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return fmt.Errorf("skill-seed not initialized. Run 'skill-seed init' first")
	}

	// 创建生成器
	gen, err := generator.New(skillPath, projectRoot)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}
	defer gen.Close()

	// 检查输出路径
	outputDir := filepath.Join(outputPath, "skill-seed-skills")
	if _, err := os.Stat(outputDir); err == nil && !force {
		fmt.Printf(i18n.Get("generate_exists")+"\n", outputDir)
		fmt.Println(i18n.Get("generate_use_force"))
		return fmt.Errorf("skills already exist")
	}

	// 生成 skills
	fmt.Println(i18n.Get("generate_generating"))
	fmt.Println("")

	if err := gen.Generate(outputDir); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(i18n.Get("generate_success"))
	fmt.Println()
	fmt.Println(i18n.Get("generate_next_steps"))
	fmt.Println("")
	fmt.Printf(i18n.Get("generate_step1")+"\n", outputDir)
	fmt.Println("")
	fmt.Println(i18n.Get("generate_step2"))
	fmt.Println("")
	fmt.Println(i18n.Get("generate_step3"))

	return nil
}
