package check

import (
	"fmt"
	"os"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/utils"
	"github.com/openclaw-coding/skill-seed/internal/checker"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/spf13/cobra"
)

// Cmd 返回检查命令
func Cmd() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: i18n.Get("cmd_check_short"),
		Long:  i18n.Get("cmd_check_long") + `

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

	return checkCmd
}

func runCheck() error {
	// Find skill path
	skillPath, err := utils.RequireSkillPath()
	if err != nil {
		return err
	}

	// Find project root directory
	projectRoot := utils.GetProjectRoot(skillPath)

	// Create checker
	chk, err := checker.New(skillPath, projectRoot)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	defer chk.Close()

	// Run check
	return chk.Run()
}
