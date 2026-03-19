package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/grow-check/internal/config"
	"github.com/openclaw-coding/grow-check/internal/git"
	"github.com/openclaw-coding/grow-check/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize grow-check in current project",
	Long: `Initialize grow-check as a project-level skill.

This will:
  - Create .skills/grow-check/ directory
  - Generate default configuration
  - Install Git pre-commit hook
  - Initialize memory database

Run this command in the root directory of your Git repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initializeSkill(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Init failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initializeSkill() error {
	// 获取项目根目录
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 检查是否是 Git 仓库
	gitOp := git.NewGitOperator(projectRoot)
	if !gitOp.IsGitRepo() {
		return fmt.Errorf("not a git repository: %s", projectRoot)
	}

	// 检查是否已初始化
	skillPath := filepath.Join(projectRoot, ".skills", "grow-check")
	if _, err := os.Stat(skillPath); err == nil {
		return fmt.Errorf("grow-check already initialized at %s", skillPath)
	}

	fmt.Println("📦 Initializing grow-check...")

	// 1. 创建目录结构
	fmt.Println("  Creating directory structure...")
	dirs := []string{
		skillPath,
		filepath.Join(skillPath, "memory"),
		filepath.Join(skillPath, "hooks"),
		filepath.Join(skillPath, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 2. 生成配置
	fmt.Println("  Generating configuration...")
	projectName := gitOp.GetProjectName()
	gitRemote, _ := gitOp.GetRemoteURL()

	cfg := config.DefaultConfig(projectName, gitRemote)
	if err := cfg.Save(skillPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 3. 初始化数据库
	fmt.Println("  Initializing memory database...")
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	store, err := storage.New(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	store.Close()

	// 4. 安装 Git 钩子
	fmt.Println("  Installing Git pre-commit hook...")
	if err := gitOp.InstallPreCommitHook(skillPath); err != nil {
		return fmt.Errorf("failed to install hook: %w", err)
	}

	// 5. 创建 pre-commit 钩子脚本
	fmt.Println("  Creating hook script...")
	if err := createHookScript(skillPath); err != nil {
		return fmt.Errorf("failed to create hook script: %w", err)
	}

	// 6. 创建 README
	fmt.Println("  Creating README...")
	if err := createReadme(skillPath, projectName); err != nil {
		return fmt.Errorf("failed to create readme: %w", err)
	}

	fmt.Println("\n✅ grow-check initialized successfully!")
	fmt.Printf("\n📁 Skill location: %s\n", skillPath)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Learn from history: grow-check learn --since=30d")
	fmt.Println("  2. Make commits and watch it learn!")
	fmt.Println("  3. View patterns: grow-check patterns")
	fmt.Println("  4. View rules: grow-check rules")

	return nil
}

func createHookScript(skillPath string) error {
	hookPath := filepath.Join(skillPath, "hooks", "pre-commit")
	
	// 获取二进制文件路径
	binPath := filepath.Join(skillPath, "bin", "grow-check")

	content := fmt.Sprintf(`#!/bin/sh
# grow-check pre-commit hook
# This hook learns from your commits and checks code quality

# Check if grow-check is initialized
if [ ! -d "%s" ]; then
    echo "grow-check not initialized, skipping..."
    exit 0
fi

# Run the check
%s check
exit $?
`, skillPath, binPath)

	if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
		return err
	}

	return nil
}

func createReadme(skillPath, projectName string) error {
	readmePath := filepath.Join(skillPath, "README.md")
	
	content := fmt.Sprintf(`# grow-check for %s

This is your project's growing code checker that learns from Git history.

## What it does

- ✅ Checks code before every commit
- 🤖 Uses Claude for deep analysis
- 📚 Learns from your team's patterns
- 🔧 Supports auto-fix

## Configuration

Edit ` + "`config.yaml`" + ` to customize:
- Enable/disable Claude analysis
- Adjust learning parameters
- Configure auto-fix behavior

## Commands

- ` + "`grow-check learn`" + ` - Learn from Git history
- ` + "`grow-check check`" + ` - Run manual check
- ` + "`grow-check patterns`" + ` - View learned patterns
- ` + "`grow-check rules`" + ` - View generated rules

## Files

- ` + "`config.yaml`" + ` - Configuration
- ` + "`memory/project.db`" + ` - Learned patterns (BoltDB)
- ` + "`memory/history.log`" + ` - Learning history

---

*This skill was initialized on your project and will grow smarter over time.*
`, projectName)

	return os.WriteFile(readmePath, []byte(content), 0644)
}
