package analyze

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/grow-check/cmd/grow-check/utils"
	"github.com/openclaw-coding/grow-check/internal/checker"
	"github.com/openclaw-coding/grow-check/internal/i18n"
	"github.com/spf13/cobra"
)

// Cmd 返回分析命令
func Cmd() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:   "analyze [files...]",
		Short: i18n.Get("cmd_analyze_short"),
		Long:  i18n.Get("cmd_analyze_long") + `

This is useful for:
  - Analyzing files before staging them
  - Checking files in a different branch
  - Quick analysis without Git context

Examples:
  grow-check analyze main.go
  grow-check analyze src/
  grow-check analyze *.go`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Error: no files provided")
				os.Exit(1)
			}
			if err := analyzeFiles(args); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	return analyzeCmd
}

func analyzeFiles(files []string) error {
	// Find skill path
	skillPath, err := utils.RequireSkillPath()
	if err != nil {
		return err
	}

	// Find project root
	projectRoot := utils.GetProjectRoot(skillPath)

	// Create checker
	chk, err := checker.New(skillPath, projectRoot)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	defer chk.Close()

	// Analyze files
	fmt.Printf("Analyzing %d files...\n", len(files))
	for _, file := range files {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Println("")

	// Convert files to absolute paths
	absFiles := make([]string, 0, len(files))
	for _, file := range files {
		if !filepath.IsAbs(file) {
			abs := filepath.Join(projectRoot, file)
			absFiles = append(absFiles, abs)
		} else {
			absFiles = append(absFiles, file)
		}
	}

	// Run analysis
	return chk.AnalyzeFiles(absFiles)
}
