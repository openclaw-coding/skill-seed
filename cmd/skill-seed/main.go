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
	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/utils"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/spf13/cobra"
)

func main() {
	// 1. 查找 .skill-seed 目录
	seedPath, err := utils.GetSeedPath()

	// 2. 初始化 i18n（在创建任何命令之前）
	locale := "zh-CN" // 默认中文
	if err == nil && seedPath != "" {
		// 尝试从配置文件读取语言设置
		if configData, err := utils.LoadConfig(seedPath); err == nil && configData != nil && configData.Project.Locale != "" {
			locale = configData.Project.Locale
		}
	}

	if err := i18n.Init(locale); err != nil {
		// i18n 初始化失败不应该阻止程序运行，只输出警告
		fmt.Fprintf(os.Stderr, "Warning: failed to initialize i18n: %v\n", err)
	}

	// 3. 创建 rootCmd（使用 i18n）
	rootCmd := &cobra.Command{
		Use:   "skill-seed",
		Short: i18n.Get("RootShort"),
		Long:  i18n.Get("RootLong"),
	}

	// 4. 注册命令
	// 创建应用容器（如果 .skill-seed 存在）
	var cont *container.Container
	if err == nil && seedPath != "" {
		ctx := context.Background()
		cont, err = container.NewContainer(ctx, seedPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer cont.Close()
	}

	// 注册所有命令
	rootCmd.AddCommand(initcmd.Cmd())
	rootCmd.AddCommand(learn.Cmd(cont))
	rootCmd.AddCommand(check.Cmd(cont))
	rootCmd.AddCommand(generate.Cmd(cont))
	rootCmd.AddCommand(hook.Cmd())
	rootCmd.AddCommand(analyze.Cmd(cont))
	rootCmd.AddCommand(view.Cmd(cont))
	rootCmd.AddCommand(scan.Cmd(cont))

	// 5. 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
