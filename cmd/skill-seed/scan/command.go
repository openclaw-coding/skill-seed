package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/utils"
	"github.com/openclaw-coding/skill-seed/internal/checker"
	"github.com/openclaw-coding/skill-seed/internal/generator"
	"github.com/openclaw-coding/skill-seed/internal/git"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/storage"
	"github.com/spf13/cobra"
)

// Cmd 返回扫描命令
func Cmd() *cobra.Command {
	var scanAll bool
	var generateSkills bool
	var outputPath string

	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: i18n.Get("cmd_scan_short"),
		Long:  i18n.Get("cmd_scan_long") + `

This will:
  - Scan all Go files in the current project
  - Use Claude to analyze code patterns
  - Learn from the current codebase state
  - Mark current commit as learned
  - Optionally generate Claude skills

Examples:
  skill-seed scan
  skill-seed scan --all
  skill-seed scan --generate-skills
  skill-seed scan --generate-skills --output .claude/skills`,
		Run: func(cmd *cobra.Command, args []string) {
			scanAll, _ = cmd.Flags().GetBool("all")
			generateSkills, _ = cmd.Flags().GetBool("generate-skills")
			outputPath, _ = cmd.Flags().GetString("output")
			if err := runScan(scanAll, generateSkills, outputPath); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		},
	}

	scanCmd.Flags().BoolP("all", "a", false, "Scan all files (not just Go files)")
	scanCmd.Flags().BoolP("generate-skills", "g", false, "Generate Claude skills after scanning")
	scanCmd.Flags().StringP("output", "o", ".claude/skills", "Output path for generated skills")

	return scanCmd
}

func runScan(scanAll bool, generateSkills bool, outputPath string) error {
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

	fmt.Println(i18n.Get("scan_analyzing_project"))
	fmt.Printf(i18n.Get("scan_project_root")+"\n", projectRoot)
	fmt.Println("")

	// Get all Go files
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
		fmt.Println(i18n.Get("scan_no_files"))
		return nil
	}

	fmt.Printf(i18n.Get("scan_found_files")+"\n", len(files))
	fmt.Println("")

	// Convert to absolute paths
	absFiles := make([]string, len(files))
	for i, file := range files {
		absFiles[i] = filepath.Join(projectRoot, file)
	}

	// Run analysis
	if err := chk.AnalyzeFiles(absFiles); err != nil {
		return err
	}

	// Close checker before accessing database again
	chk.Close()

	// Mark current commit as learned
	gitOp := git.NewGitOperator(projectRoot)
	currentHash, err := gitOp.GetCurrentCommitHash()
	if err == nil && currentHash != "" {
		if err := markCommitAsLearned(skillPath, currentHash); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to mark commit as learned: %v\n", err)
		} else {
			fmt.Printf(i18n.Get("scan_marked_learned")+"\n", currentHash[:8])
		}
	}

	fmt.Println("")
	fmt.Println(i18n.Get("scan_completed"))

	// Generate skills if requested
	if generateSkills {
		fmt.Println("")
		fmt.Println("Generating Claude skills...")

		gen, err := generator.New(skillPath, projectRoot)
		if err != nil {
			return fmt.Errorf("failed to create generator: %w", err)
		}
		defer gen.Close()

		outputDir := filepath.Join(outputPath, "skill-seed-skills")
		if err := gen.Generate(outputDir); err != nil {
			return fmt.Errorf("failed to generate skills: %w", err)
		}
	}

	return nil
}

func findGoFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and common exclude dirs
		if info.IsDir() {
			name := filepath.Base(path)
			if name == ".git" || name == ".skill-seed" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only include Go files
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

		// Skip hidden directories and common exclude dirs
		if info.IsDir() {
			name := filepath.Base(path)
			if name == ".git" || name == ".skill-seed" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip common binary/generated files
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

func markCommitAsLearned(skillPath, commitHash string) error {
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	store, err := storage.New(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	return store.SaveLearnRecord(commitHash, time.Now())
}
