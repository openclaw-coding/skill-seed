package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/grow-check/internal/i18n"
)

// FindSkillPath 查找 .skills/grow-check 目录
// 从当前目录向上查找，直到找到包含 .skills/grow-check 的目录
func FindSkillPath() (string, error) {
	// 从当前目录开始
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// 向上查找 .skills/grow-check 目录
	for {
		skillPath := filepath.Join(dir, ".skills", "grow-check")
		if _, err := os.Stat(skillPath); err == nil {
			return skillPath, nil
		}

		// 到达根目录
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("grow-check not initialized in current directory or any parent directory")
}

// PrintInitHelp 打印初始化提示信息
func PrintInitHelp() {
	fmt.Println("")
	fmt.Print(i18n.Get("check_init_failed"))
	fmt.Println("")
	fmt.Print(i18n.Get("check_init_hint"))
	fmt.Println("")
	fmt.Print(i18n.Get("check_init_command"))
	fmt.Println("")
	fmt.Print(i18n.Get("check_init_more_info"))
	fmt.Println("")
}

// RequireSkillPath 查找 skill path，如果失败则打印帮助信息并返回错误
func RequireSkillPath() (string, error) {
	skillPath, err := FindSkillPath()
	if err != nil {
		PrintInitHelp()
		return "", fmt.Errorf("grow-check not initialized")
	}
	return skillPath, nil
}

// GetProjectRoot 从 skillPath 获取项目根目录
func GetProjectRoot(skillPath string) string {
	return filepath.Dir(filepath.Dir(skillPath))
}
