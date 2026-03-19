package initialize

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

// Cmd 返回初始化命令
func Cmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: i18n.Get("cmd_init_short"),
		Long:  i18n.Get("init_long"),
		Run: func(cmd *cobra.Command, args []string) {
			if err := initializeSkill(); err != nil {
				fmt.Println(i18n.Get("init_failed"), err)
				os.Exit(1)
			}
		},
	}

	return initCmd
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
		return fmt.Errorf("grow-check already initialized")
	}

	fmt.Println(i18n.Get("init_start"))

	// 1. Create directory structure
	fmt.Print(i18n.Get("init_creating_dirs"))
	fmt.Println("")
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
	fmt.Print(i18n.Get("init_generating_config"))
	fmt.Println("")
	projectName := gitOp.GetProjectName()
	gitRemote, _ := gitOp.GetRemoteURL()

	cfg := config.DefaultConfig(projectName, gitRemote)
	if err := cfg.Save(skillPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 3. Initialize database
	fmt.Print(i18n.Get("init_initializing_db"))
	fmt.Println("")
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	store, err := storage.New(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	store.Close()

	// 4. Install Git hook
	fmt.Print(i18n.Get("init_installing_hook"))
	fmt.Println("")
	if err := gitOp.InstallPreCommitHook(skillPath); err != nil {
		fmt.Printf(i18n.Get("init_hook_install_failed"), err)
		fmt.Println("")
		fmt.Println(i18n.Get("init_hook_install_manual"))
	} else {
		fmt.Println(i18n.Get("init_hook_installed"))
	}

	// 5. Create hook script
	fmt.Print(i18n.Get("init_creating_hook_script"))
	fmt.Println("")
	hookScriptPath := filepath.Join(skillPath, "hooks", "pre-commit")
	if err := createHookScript(hookScriptPath, skillPath); err != nil {
		return fmt.Errorf("failed to create hook script: %w", err)
	}

	// 6. Create README
	fmt.Print(i18n.Get("init_creating_readme"))
	fmt.Println("")
	if err := createReadme(skillPath); err != nil {
		return fmt.Errorf("failed to create readme: %w", err)
	}

	fmt.Println("")
	fmt.Println(i18n.Get("init_success"))
	fmt.Println("")
	fmt.Printf(i18n.Get("init_skill_location")+"\n", skillPath)
	fmt.Println("")
	fmt.Println(i18n.Get("init_next_steps"))
	fmt.Println("")
	fmt.Println(i18n.Get("init_step_learn"))
	fmt.Println("")
	fmt.Println(i18n.Get("init_step_watch"))
	fmt.Println("")
	fmt.Println(i18n.Get("init_step_patterns"))
	fmt.Println("")
	fmt.Println(i18n.Get("init_step_rules"))
	fmt.Println("")

	return nil
}

func createHookScript(hookPath, skillPath string) error {
	script := fmt.Sprintf(`#!/bin/sh
# grow-check pre-commit hook

# Run grow-check check
cd "$(git rev-parse --show-toplevel)"
grow-check check
if [ $? -ne 0 ]; then
    echo ""
    echo "%s"
    exit 1
fi
`, i18n.Get("learn_tip_force"))
	return os.WriteFile(hookPath, []byte(script), 0755)
}

func createReadme(skillPath string) error {
	readme := fmt.Sprintf(`# grow-check for %s

This is your project's growing code checker that learns from Git history.

## What it does

- Checks code before every commit
- Uses Claude for deep analysis
- Learns from your team's patterns
- Supports auto-fix

## Configuration

Edit %%s/config.yaml to customize:
- Enable/disable Claude analysis
- Adjust learning parameters
- Configure auto-fix behavior

## Commands

- grow-check learn - Learn from Git history
- grow-check check - Run manual check
- grow-check patterns - View learned patterns
- grow-check rules - View generated rules

## Files

- config.yaml - Configuration
- memory/project.db - Learned patterns (BoltDB)
- hooks/pre-commit - Git pre-commit hook

---

*This skill was initialized on your project and will grow smarter over time.*
`, filepath.Base(skillPath), skillPath)

	return os.WriteFile(filepath.Join(skillPath, "README.md"), []byte(readme), 0644)
}
