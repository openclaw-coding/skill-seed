package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/grow-check/internal/checker"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run pre-commit check manually",
	Long: `Run the same checks that the pre-commit hook would run.

This is useful for:
  - Testing your changes before committing
  - Seeing what the hook would find
  - Debugging issues

The check will:
  1. Scan staged files
  2. Run basic pattern checks
  3. Use Claude for deep analysis (if available)
  4. Show results and offer fixes`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCheck(); err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ Check failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck() error {
	// 查找 skill 路径
	skillPath, err := findSkillPath()
	if err != nil {
		return err
	}

	// 查找项目根目录
	projectRoot := filepath.Dir(filepath.Dir(skillPath))

	// 创建检查器
	check, err := checker.New(skillPath, projectRoot)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	defer check.Close()

	// 运行检查
	return check.Run()
}
