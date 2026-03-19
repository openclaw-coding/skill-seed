package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	skillPath := filepath.Join(projectRoot, ".seed", "skill-seed")
	if _, err := os.Stat(skillPath); err == nil {
		return fmt.Errorf("skill-seed already initialized")
	}

	fmt.Println(i18n.Get("init_start"))

	// 1. Create directory structure
	fmt.Print(i18n.Get("init_creating_dirs"))
	fmt.Println("")
	dirs := []string{
		skillPath,
		filepath.Join(skillPath, "memory"),
		filepath.Join(skillPath, "references"),
		filepath.Join(skillPath, "references", "naming-patterns"),
		filepath.Join(skillPath, "references", "error-handling-patterns"),
		filepath.Join(skillPath, "references", "structure-patterns"),
		filepath.Join(skillPath, "references", "concurrency-patterns"),
		filepath.Join(skillPath, "references", "testing-patterns"),
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

	// 4. Create SKILL.md
	fmt.Print(i18n.Get("init_creating_skill"))
	fmt.Println("")
	if err := createSkillFile(skillPath, projectName, gitRemote); err != nil {
		return fmt.Errorf("failed to create SKILL.md: %w", err)
	}

	// 5. Create pattern category overviews
	fmt.Print(i18n.Get("init_creating_patterns"))
	fmt.Println("")
	if err := createPatternOverviews(skillPath); err != nil {
		return fmt.Errorf("failed to create pattern overviews: %w", err)
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

func createSkillFile(skillPath, projectName, gitRemote string) error {
	backtick := "`"
	skillContent := fmt.Sprintf(`---
name: %s-skill-seed
description: Project-specific coding patterns learned from Git history. Use this skill when working on this project to follow team conventions and best practices.
version: 1.0.0
---

# %s Coding Standards

**Auto-generated knowledge base** from Git history analysis. This skill represents the living coding patterns and conventions of this project, learned from actual commits and code changes.

## Overview

This skill captures the **evolving coding standards** of this project by:
- Learning patterns from Git commit history
- Identifying team conventions over time
- Extracting best practices from real code
- Detecting common anti-patterns to avoid
- Growing smarter with each commit

**Project**: %s
**Repository**: %s
**Last Updated**: %s

## Quick Start

When working on this project:

1. **For naming conventions**: Check [Naming Patterns](references/naming-patterns/overview.md)
2. **For error handling**: Review [Error Handling Patterns](references/error-handling-patterns/overview.md)
3. **For code structure**: See [Structure Patterns](references/structure-patterns/overview.md)
4. **For concurrent code**: Study [Concurrency Patterns](references/concurrency-patterns/overview.md)
5. **For tests**: Follow [Testing Patterns](references/testing-patterns/overview.md)

## How This Skill Grows

This skill is **living documentation** that evolves with your codebase:

1. **Learning Phase**: Run `+backtick+`skill-seed learn`+backtick+` to analyze recent commits
2. **Pattern Extraction**: Claude identifies recurring patterns in your code
3. **Knowledge Base**: Patterns are stored in the memory database
4. **Skill Updates**: Run `+backtick+`skill-seed generate-skills`+backtick+` to update this skill

## Commands

- `+backtick+`skill-seed learn`+backtick+` - Learn from Git history
- `+backtick+`skill-seed check`+backtick+` - Run manual check
- `+backtick+`skill-seed analyze <files>`+backtick+` - Analyze specific files
- `+backtick+`skill-seed view patterns`+backtick+` - View learned patterns
- `+backtick+`skill-seed view rules`+backtick+` - View generated rules
- `+backtick+`skill-seed generate-skills`+backtick+` - Update this skill from learned patterns

## Pattern Categories

### Naming Patterns
File, function, variable, and naming conventions that evolved in this project.

### Error Handling Patterns
How errors are checked, wrapped, and propagated throughout the codebase.

### Structure Patterns
Code organization, directory layout, and architectural patterns.

### Concurrency Patterns
Goroutine usage, channel patterns, and synchronization approaches.

### Testing Patterns
Test organization, naming, and testing conventions used in the project.

---

*This skill was initialized on %s and will grow smarter with every commit.*
`,
		projectName,
		projectName,
		projectName,
		gitRemote,
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02"),
	)

	// Replace the placeholders
	skillContent = strings.ReplaceAll(skillContent, "+backtick+", "`")

	return os.WriteFile(filepath.Join(skillPath, "SKILL.md"), []byte(skillContent), 0644)
}

func createPatternOverviews(skillPath string) error {
	bt := "`"

	// Create overview.md for each category
	overviews := map[string]string{
		"naming-patterns": `# Naming Patterns

Naming conventions learned from this project's codebase.

## Overview

This section captures the naming patterns that have evolved in this project.

*No patterns learned yet. Run ` + bt + `skill-seed learn` + bt + ` to start learning from Git history.*

## Files in this category

- [File Naming](file-naming.md) - How files are named
- [Function Naming](function-naming.md) - Function naming conventions
- [Variable Naming](variable-naming.md) - Variable naming patterns
- [Interface Naming](interface-naming.md) - Interface naming standards
`,
		"error-handling-patterns": `# Error Handling Patterns

Error handling conventions learned from this project's codebase.

## Overview

This section captures error handling patterns that have evolved in this project.

*No patterns learned yet. Run ` + bt + `skill-seed learn` + bt + ` to start learning from Git history.*

## Files in this category

- [Error Checking](error-checking.md) - How errors are checked
- [Error Wrapping](error-wrapping.md) - Error wrapping patterns
- [Error Messages](error-messages.md) - Error message formats
`,
		"structure-patterns": `# Structure Patterns

Code structure and organization patterns learned from this project.

## Overview

This section captures architectural patterns that have evolved in this project.

*No patterns learned yet. Run ` + bt + `skill-seed learn` + bt + ` to start learning from Git history.*

## Files in this category

- [Directory Structure](directory-structure.md) - Project layout
- [File Organization](file-organization.md) - How code is organized
- [Package Structure](package-structure.md) - Package organization
`,
		"concurrency-patterns": `# Concurrency Patterns

Concurrency patterns learned from this project's codebase.

## Overview

This section captures goroutine and channel patterns that have evolved in this project.

*No patterns learned yet. Run ` + bt + `skill-seed learn` + bt + ` to start learning from Git history.*

## Files in this category

- [Goroutine Usage](goroutine-usage.md) - How goroutines are used
- [Channel Patterns](channel-patterns.md) - Channel usage patterns
- [Synchronization](synchronization.md) - Mutex and sync patterns
`,
		"testing-patterns": `# Testing Patterns

Testing conventions learned from this project's codebase.

## Overview

This section captures testing patterns that have evolved in this project.

*No patterns learned yet. Run ` + bt + `skill-seed learn` + bt + ` to start learning from Git history.*

## Files in this category

- [Test Organization](test-organization.md) - How tests are organized
- [Test Naming](test-naming.md) - Test naming conventions
- [Test Patterns](test-patterns.md) - Common test patterns
`,
	}

	for category, content := range overviews {
		if err := os.WriteFile(
			filepath.Join(skillPath, "references", category, "overview.md"),
			[]byte(content),
			0644,
		); err != nil {
			return err
		}
	}

	return nil
}
