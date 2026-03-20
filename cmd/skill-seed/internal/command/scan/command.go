package scan

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/container"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/spf13/cobra"
)

var (
	scanAll bool
)

// Cmd 返回 scan 命令
func Cmd(cont *container.Container) *cobra.Command {
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: i18n.Get("ScanShort"),
		Long:  i18n.Get("ScanLongDesc"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScan(cont, cmd)
		},
	}

	scanCmd.Flags().BoolVarP(&scanAll, "all", "a", false, "Scan all files (not just Go files)")

	return scanCmd
}

func runScan(cont *container.Container, cmd *cobra.Command) error {
	ctx := context.Background()

	// 获取项目根目录
	projectRoot, err := cont.GetGitRepository().GetProjectRoot(ctx)
	if err != nil {
		return fmt.Errorf("failed to get project root: %w", err)
	}

	fmt.Println(i18n.Get("ScanAnalyzingProject"))
	fmt.Printf("  Project root: %s\n\n", projectRoot)

	// 获取所有文件
	var files []string
	if scanAll {
		files, err = findAllFiles(projectRoot)
	} else {
		files, err = findGoFiles(projectRoot)
	}

	if err != nil {
		return fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println(i18n.Get("ScanNoFiles"))
		return nil
	}

	fmt.Println(i18n.GetWithParams("ScanFoundFiles", map[string]interface{}{"Count": len(files)}))
	fmt.Println()

	// 转换为绝对路径
	absFiles := make([]string, len(files))
	for i, file := range files {
		absFiles[i] = filepath.Join(projectRoot, file)
	}

	// 使用 checker 服务分析文件
	if err := cont.GetCheckerService().AnalyzeFiles(ctx, absFiles); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(i18n.Get("ScanCompleted"))

	return nil
}

func findGoFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过隐藏目录和常见排除目录
		if info.IsDir() {
			name := filepath.Base(path)
			if name == ".git" || name == ".skill-seed" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// 只包含 Go 文件
		if filepath.Ext(path) == ".go" {
			relPath, err := filepath.Rel(root, path)
			if err == nil {
				files = append(files, relPath)
			}
		}

		return nil
	})

	return files, err
}

func findAllFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过隐藏目录和常见排除目录
		if info.IsDir() {
			name := filepath.Base(path)
			if name == ".git" || name == ".skill-seed" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过常见的二进制/生成文件
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext != ".exe" && ext != ".bin" && ext != ".o" && ext != ".a" {
				relPath, err := filepath.Rel(root, path)
				if err == nil {
					files = append(files, relPath)
				}
			}
		}

		return nil
	})

	return files, err
}
