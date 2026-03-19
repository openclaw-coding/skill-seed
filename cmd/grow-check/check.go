package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/grow-check/internal/checker"
	"github.com/openclaw-coding/grow-check/internal/i18n"
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
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck() error {
	// Find skill path
	skillPath, err := findSkillPath()
	if err != nil {
		// Provide friendly error message
		println("")
		print(i18n.Get("check_init_failed"))
		println("")
		print(i18n.Get("check_init_hint"))
		println("")
		print(i18n.Get("check_init_command"))
		println("")
		print(i18n.Get("check_init_more_info"))
		println("")
		return fmt.Errorf("grow-check not initialized")
	}

	// Find project root directory
	projectRoot := filepath.Dir(filepath.Dir(skillPath))

	// Create checker
	check, err := checker.New(skillPath, projectRoot)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	defer check.Close()

	// Run check
	return check.Run()
}
