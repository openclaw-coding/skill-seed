package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/skill-seed/internal/config"
	"github.com/openclaw-coding/skill-seed/internal/git"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/storage"
	"github.com/spf13/cobra"
)

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
	skillPath := filepath.Join(projectRoot, ".skill-seed")
	if _, err := os.Stat(skillPath); err == nil {
		return fmt.Errorf("skill-seed already initialized")
	}

	fmt.Println(i18n.Get("init_start"))

	// 1. Create directory structure
	dirs := []string{
		skillPath,
		filepath.Join(skillPath, "memory"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 2. Generate configuration
	projectName := gitOp.GetProjectName()
	gitRemote, _ := gitOp.GetRemoteURL()

	cfg := config.DefaultConfig(projectName, gitRemote)
	if err := cfg.Save(skillPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 3. Initialize database
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	store, err := storage.New(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	store.Close()

	fmt.Println(i18n.Get("init_success"))
	fmt.Printf(i18n.Get("init_skill_location")+"\n", skillPath)
	fmt.Println(i18n.Get("init_next_steps"))
	fmt.Println(i18n.Get("init_step_learn"))
	fmt.Println(i18n.Get("init_step_watch"))
	fmt.Println(i18n.Get("init_step_patterns"))
	fmt.Println(i18n.Get("init_step_rules"))

	return nil
}
