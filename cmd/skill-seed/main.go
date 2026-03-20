package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/analyze"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/check"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/generate"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/hook"
	initcmd "github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/init"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/learn"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/scan"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/command/view"
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/container"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skill-seed",
	Short: "Growing project skills for AI agents",
	Long: `A project-level skill that learns from Git history
to help AI agents understand your codebase better.

This tool integrates with Claude for deep code analysis
and learns your team's coding patterns automatically.`,
}

func main() {
	// 1. 查找 .skill-seed 目录
	skillPath, err := container.GetSkillPath()
	if err != nil {
		// 如果没有找到，只能运行 init 命令
		rootCmd.AddCommand(initcmd.Cmd())
	} else {
		// 2. 创建应用容器
		ctx := context.Background()
		cont, err := container.NewContainer(ctx, skillPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer cont.Close()

		// 3. 注册子命令
		rootCmd.AddCommand(initcmd.Cmd())
		rootCmd.AddCommand(learn.Cmd(cont))
		rootCmd.AddCommand(check.Cmd(cont))
		rootCmd.AddCommand(generate.Cmd(cont))
		rootCmd.AddCommand(hook.Cmd())
		rootCmd.AddCommand(analyze.Cmd(cont))
		rootCmd.AddCommand(view.Cmd(cont))
		rootCmd.AddCommand(scan.Cmd(cont))
	}

	// 4. 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
