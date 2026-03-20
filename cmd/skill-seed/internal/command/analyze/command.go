package analyze

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/internal/container"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/spf13/cobra"
)

// Cmd 返回分析命令
func Cmd(cont *container.Container) *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:   "analyze [files...]",
		Short: i18n.Get("AnalyzeShort"),
		Long:  i18n.Get("AnalyzeLongDesc"),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, i18n.Get("AnalyzeNoFiles"))
				os.Exit(1)
			}
			if err := analyzeFiles(cont, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	return analyzeCmd
}

func analyzeFiles(cont *container.Container, files []string) error {
	ctx := context.Background()

	// 获取项目根目录
	projectRoot, err := cont.GetGitRepository().GetProjectRoot(ctx)
	if err != nil {
		return fmt.Errorf("%s", i18n.GetWithParams("AnalyzeGetRootFailed", map[string]interface{}{"Error": err.Error()}))
	}

	// 转换为绝对路径
	absFiles := make([]string, 0, len(files))
	for _, file := range files {
		if !filepath.IsAbs(file) {
			abs := filepath.Join(projectRoot, file)
			absFiles = append(absFiles, abs)
		} else {
			absFiles = append(absFiles, file)
		}
	}

	fmt.Println(i18n.GetWithParams("AnalyzeAnalyzing", map[string]interface{}{"Count": len(files)}))
	for _, file := range files {
		fmt.Println(i18n.GetWithParams("AnalyzeFile", map[string]interface{}{"File": file}))
	}
	fmt.Println("")

	// 使用 checker 服务分析文件
	return cont.GetCheckerService().AnalyzeFiles(ctx, absFiles)
}
