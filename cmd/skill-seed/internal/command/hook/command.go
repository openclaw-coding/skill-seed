package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/spf13/cobra"
)

// Cmd 返回 hook 命令
func Cmd() *cobra.Command {
	var install bool
	var uninstall bool

	hookCmd := &cobra.Command{
		Use:   "hook",
		Short: i18n.Get("HookShort"),
		Long:  i18n.Get("HookLongDesc"),
		Run: func(cmd *cobra.Command, args []string) {
			if install && uninstall {
				fmt.Println(i18n.Get("HookBothFlagsError"))
				os.Exit(1)
			}

			if install {
				if err := installHook(); err != nil {
					fmt.Println(i18n.GetWithParams("HookInstallFailed", map[string]interface{}{"Error": err.Error()}))
					os.Exit(1)
				}
				fmt.Println(i18n.Get("HookInstallSuccess"))
			} else if uninstall {
				if err := uninstallHook(); err != nil {
					fmt.Println(i18n.GetWithParams("HookUninstallFailed", map[string]interface{}{"Error": err.Error()}))
					os.Exit(1)
				}
				fmt.Println(i18n.Get("HookUninstallSuccess"))
			} else {
				// 默认执行 pre-commit hook
				if err := runPreCommitHook(); err != nil {
					fmt.Println(i18n.GetWithParams("HookRunFailed", map[string]interface{}{"Error": err.Error()}))
					os.Exit(1)
				}
			}
		},
	}

	hookCmd.Flags().BoolVarP(&install, "install", "i", false, i18n.Get("HookFlagInstall"))
	hookCmd.Flags().BoolVarP(&uninstall, "uninstall", "u", false, i18n.Get("HookFlagUninstall"))

	return hookCmd
}

func installHook() error {
	// 获取项目根目录
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 检查是否是 Git 仓库
	gitDir := filepath.Join(projectRoot, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository")
	}

	// 检查 .skill-seed 是否已初始化
	seedPath := filepath.Join(projectRoot, ".skill-seed")
	if _, err := os.Stat(seedPath); os.IsNotExist(err) {
		return fmt.Errorf("skill-seed not initialized, run 'skill-seed init' first")
	}

	// 创建 hook 脚本
	hookPath := filepath.Join(gitDir, "hooks", "pre-commit")
	hookContent := `#!/bin/bash
# skill-seed pre-commit hook

# 获取暂存的文件
STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')

if [ -z "$STAGED_FILES" ]; then
    exit 0
fi

echo "Running skill-seed check..."

# 运行 skill-seed check
if ! skill-seed check; then
    echo "skill-seed check found issues. Please fix them before committing."
    exit 1
fi

exit 0
`

	// 确保 hooks 目录存在
	hooksDir := filepath.Join(gitDir, "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// 写入 hook 脚本
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	return nil
}

func uninstallHook() error {
	// 获取项目根目录
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 检查 hook 文件是否存在
	hookPath := filepath.Join(projectRoot, ".git", "hooks", "pre-commit")
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return fmt.Errorf("pre-commit hook not found")
	}

	// 删除 hook 文件
	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("failed to remove hook file: %w", err)
	}

	return nil
}

func runPreCommitHook() error {
	// 获取暂存的 Go 文件
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get staged files: %w", err)
	}

	files := strings.Split(string(output), "\n")
	goFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file, ".go") && file != "" {
			goFiles = append(goFiles, file)
		}
	}

	if len(goFiles) == 0 {
		fmt.Println(i18n.Get("HookNoStagedFiles"))
		return nil
	}

	fmt.Println(i18n.GetWithParams("HookCheckingFiles", map[string]interface{}{"Count": len(goFiles)}))

	// 运行 skill-seed check
	checkCmd := exec.Command("skill-seed", "check")
	checkCmd.Stdout = os.Stdout
	checkCmd.Stderr = os.Stderr

	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("check failed: %w", err)
	}

	fmt.Println(i18n.Get("HookCheckPassed"))
	return nil
}
