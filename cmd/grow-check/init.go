package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/grow-check/internal/config"
	"github.com/openclaw-coding/grow-check/internal/git"
	"github.com/openclaw-coding/grow-check/internal/i18n"
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
			println(i18n.Get("init_failed"), err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initializeSkill() error {
	// Get project root directory
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if it's a Git repository
	gitOp := git.NewGitOperator(projectRoot)
	if !gitOp.IsGitRepo() {
		return fmt.Errorf("not a git repository: %s", projectRoot)
	}

	// Check if already initialized
	skillPath := filepath.Join(projectRoot, ".skills", "grow-check")
	if _, err := os.Stat(skillPath); err == nil {
		return fmt.Errorf("grow-check already initialized at %s", skillPath)
	}

	println("📦 Initializing grow-check...")

	// 1. Create directory structure
	print(i18n.Get("init_creating_dirs"))
	println("")
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

	// 2. Generate configuration
	print(i18n.Get("init_generating_config"))
	println("")
	projectName := gitOp.GetProjectName()
	gitRemote, _ := gitOp.GetRemoteURL()

	cfg := config.DefaultConfig(projectName, gitRemote)
	if err := cfg.Save(skillPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 3. Initialize database
	print(i18n.Get("init_initializing_db"))
	println("")
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	store, err := storage.New(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	store.Close()

	// 4. Install Git hook
	print(i18n.Get("init_installing_hook"))
	println("")
	if err := gitOp.InstallPreCommitHook(skillPath); err != nil {
		return fmt.Errorf("failed to install hook: %w", err)
	}

	// 5. Create pre-commit hook script
	print(i18n.Get("init_creating_hook"))
	println("")
	if err := createHookScript(skillPath); err != nil {
		return fmt.Errorf("failed to create hook script: %w", err)
	}

	// 6. Create README
	print(i18n.Get("init_creating_readme"))
	println("")
	if err := createReadme(skillPath, projectName); err != nil {
		return fmt.Errorf("failed to create readme: %w", err)
	}

	println("")
	print(i18n.Get("init_success"))
	println("")
	printf(i18n.Get("init_skill_location"), skillPath)
	println("")
	println("")
	print(i18n.Get("init_next_steps"))
	println("")
	print(i18n.Get("init_step_learn"))
	println("")
	print(i18n.Get("init_step_watch"))
	println("")
	print(i18n.Get("init_step_patterns"))
	println("")
	print(i18n.Get("init_step_rules"))
	println("")

	return nil
}

func createHookScript(skillPath string) error {
	hookPath := filepath.Join(skillPath, "hooks", "pre-commit")

	content := `#!/bin/sh
# grow-check pre-commit hook
# This hook learns from your commits and checks code quality

# Run the check (assumes grow-check is in PATH)
grow-check check
exit $?
`

	if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
		return err
	}

	return nil
}

func printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func print(s string) {
	fmt.Print(s)
}

func println(args ...interface{}) {
	fmt.Println(args...)
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
