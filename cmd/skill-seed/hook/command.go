package hook

import (
	"fmt"
	"os"

	"github.com/openclaw-coding/skill-seed/cmd/skill-seed/utils"
	"github.com/openclaw-coding/skill-seed/internal/git"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/spf13/cobra"
)

// Cmd 返回 hook 命令组
func Cmd() *cobra.Command {
	hookCmd := &cobra.Command{
		Use:   "hook",
		Short: i18n.Get("cmd_hook_short"),
		Long:  i18n.Get("cmd_hook_long"),
	}

	hookCmd.AddCommand(installCmd())
	hookCmd.AddCommand(uninstallCmd())

	return hookCmd
}

func installCmd() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: i18n.Get("cmd_hook_install_short"),
		Long:  i18n.Get("cmd_hook_install_long"),
		Run: func(cmd *cobra.Command, args []string) {
			if err := installHook(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		},
	}

	return installCmd
}

func uninstallCmd() *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: i18n.Get("cmd_hook_uninstall_short"),
		Long:  i18n.Get("cmd_hook_uninstall_long"),
		Run: func(cmd *cobra.Command, args []string) {
			if err := uninstallHook(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		},
	}

	return uninstallCmd
}

func installHook() error {
	// Find skill path
	skillPath, err := utils.RequireSkillPath()
	if err != nil {
		return err
	}

	// Get project root
	projectRoot := utils.GetProjectRoot(skillPath)

	// Create git operator
	gitOp := git.NewGitOperator(projectRoot)

	fmt.Println(i18n.Get("hook_installing"))

	// Install pre-commit hook
	if err := gitOp.InstallPreCommitHook(skillPath); err != nil {
		fmt.Printf(i18n.Get("hook_install_failed")+"\n", err)
		return fmt.Errorf("hook installation failed")
	}

	fmt.Println(i18n.Get("hook_installed"))
	fmt.Println("")
	fmt.Println(i18n.Get("hook_installed_success"))

	return nil
}

func uninstallHook() error {
	// Get current directory
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create git operator
	gitOp := git.NewGitOperator(projectRoot)

	fmt.Println(i18n.Get("hook_uninstalling"))

	// Uninstall pre-commit hook
	if err := gitOp.UninstallPreCommitHook(); err != nil {
		fmt.Printf(i18n.Get("hook_uninstall_failed")+"\n", err)
		return fmt.Errorf("hook uninstallation failed")
	}

	fmt.Println(i18n.Get("hook_uninstalled"))

	return nil
}
