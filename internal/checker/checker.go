package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openclaw-coding/grow-check/internal/claude"
	"github.com/openclaw-coding/grow-check/internal/config"
	"github.com/openclaw-coding/grow-check/internal/git"
	"github.com/openclaw-coding/grow-check/internal/i18n"
	"github.com/openclaw-coding/grow-check/internal/output"
	"github.com/openclaw-coding/grow-check/internal/storage"
	"github.com/openclaw-coding/grow-check/pkg/models"
)

// Checker code checker
type Checker struct {
	config      *config.Config
	store       *storage.Store
	git         *git.GitOperator
	claude      *claude.Client
	skillPath   string
	projectRoot string
}

// New create checker
func New(skillPath, projectRoot string) (*Checker, error) {
	// Load configuration
	cfg, err := config.Load(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Open database
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	store, err := storage.New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}

	// Create Git operator
	gitOp := git.NewGitOperator(projectRoot)

	// Create Claude client
	claudeClient := claude.NewClient(
		"claude",
		time.Duration(cfg.Claude.TimeoutSeconds)*time.Second,
		cfg.Claude.FallbackToBasic,
	)

	return &Checker{
		config:      cfg,
		store:       store,
		git:         gitOp,
		claude:      claudeClient,
		skillPath:   skillPath,
		projectRoot: projectRoot,
	}, nil
}

// Run run check
func (c *Checker) Run() error {
	// 1. Get staged files
	files, err := c.git.GetStagedFiles()
	if err != nil {
		return fmt.Errorf("failed to get staged files: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	fmt.Printf(i18n.Get("check_checking_files"), len(files))
	fmt.Println("")

	// 2. Filter files
	files = c.filterFiles(files)
	if len(files) == 0 {
		return nil
	}

	// 3. Prepare file data
	fileChanges := make([]models.FileChange, 0, len(files))
	for _, path := range files {
		content, err := c.git.GetFileContent(path)
		if err != nil {
			continue
		}

		diff, err := c.git.GetStagedFileDiff(path)
		if err != nil {
			diff = ""
		}

		fileChanges = append(fileChanges, models.FileChange{
			Path:    path,
			Content: content,
			Diff:    diff,
		})
	}

	// 4. Basic checks
	basicIssues := c.runBasicChecks(fileChanges)

	// 5. Claude deep analysis
	var claudeIssues []models.Issue
	if c.config.Claude.Enabled && c.claude.IsAvailable() {
		fmt.Print(i18n.Get("check_analyzing_claude"))
		fmt.Println("")

		context := c.buildAnalysisContext()
		result, err := c.claude.AnalyzeCode(fileChanges, context)
		if err != nil {
			fmt.Printf(i18n.Get("check_claude_failed"), err)
			fmt.Println("")
		} else if result != nil {
			claudeIssues = result.Issues
		}
	}

	// 6. Merge results
	allIssues := append(basicIssues, claudeIssues...)

	// 7. Handle results
	if len(allIssues) == 0 {
		fmt.Print(i18n.Get("check_no_issues"))
		fmt.Println("")
		return nil
	}

	return c.handleIssues(allIssues)
}

// filterFiles filter files
func (c *Checker) filterFiles(files []string) []string {
	filtered := make([]string, 0, len(files))

	for _, file := range files {
		// Only check Go files
		if !strings.HasSuffix(file, ".go") {
			continue
		}

		// Check exclusion patterns
		excluded := false
		for _, pattern := range c.config.Checking.ExcludePatterns {
			matched, err := filepath.Match(pattern, file)
			if err == nil && matched {
				excluded = true
				break
			}
		}

		if !excluded {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// runBasicChecks run basic checks
func (c *Checker) runBasicChecks(files []models.FileChange) []models.Issue {
	var issues []models.Issue

	// Load rules
	rules, err := c.store.GetAllRules()
	if err != nil {
		return issues
	}

	// Apply rules to each file
	for _, file := range files {
		for _, rule := range rules {
			if rule.Type == models.PatternNaming {
				// Simple naming check
				issues = append(issues, c.checkNaming(file, rule)...)
			}
			// Can add more rule types
		}
	}

	return issues
}

// checkNaming check naming convention
func (c *Checker) checkNaming(file models.FileChange, rule models.Rule) []models.Issue {
	var issues []models.Issue

	// Simple example: check if underscore is included (Go recommends camelCase)
	if strings.Contains(file.Path, "_") && !strings.HasSuffix(file.Path, "_test.go") {
		issues = append(issues, models.Issue{
			File:       file.Path,
			Line:       1,
			Severity:   "warning",
			Message:    "File name contains underscore, recommend using lowercase letters and hyphens",
			Suggestion: "Rename file, remove underscores",
			PatternID:  rule.ID,
		})
	}

	return issues
}

// buildAnalysisContext build analysis context
func (c *Checker) buildAnalysisContext() *models.AnalysisContext {
	context := &models.AnalysisContext{
		ProjectType:     "go",
		TeamConventions: "Follow Go official code standards",
	}

	// Load learned patterns
	patterns, err := c.store.GetAllPatterns()
	if err == nil {
		context.LearnedPatterns = patterns
	}

	// Get recent commits
	commits, err := c.git.GetRecentCommits(10, time.Time{})
	if err == nil {
		context.RecentCommits = commits
	}

	// Get historical bug patterns
	bugs, err := c.store.GetMetadata("historical_bugs")
	if err == nil && len(bugs) > 0 {
		context.HistoricalBugs = strings.Split(string(bugs), "\n")
	}

	return context
}

// handleIssues handle discovered issues
func (c *Checker) handleIssues(issues []models.Issue) error {
	fmt.Printf(i18n.Get("check_found_issues"), len(issues))
	fmt.Println("")

	for i, issue := range issues {
		label := output.SeverityLabel(issue.Severity)

		fmt.Printf("%d. %s %s:%d\n", i+1, label, issue.File, issue.Line)
		fmt.Printf("   %s %s\n", label, issue.Message)
		if issue.Suggestion != "" {
			output.Dim("   Suggestion: %s\n", issue.Suggestion)
		}
		fmt.Println("")
	}

	// Interactive handling
	if c.config.Checking.Interactive {
		return c.interactiveHandler(issues)
	}

	// Non-interactive mode: fail if there are error-level issues
	for _, issue := range issues {
		if issue.Severity == "error" {
			return fmt.Errorf("found %d error-level issues", len(issues))
		}
	}

	return nil
}

// interactiveHandler interactive handling
func (c *Checker) interactiveHandler(issues []models.Issue) error {
	fmt.Println(i18n.Get("check_interactive_options"))
	fmt.Println(i18n.Get("check_option_autofix"))
	fmt.Println(i18n.Get("check_option_details"))
	fmt.Println(i18n.Get("check_option_ignore"))
	fmt.Println(i18n.Get("check_option_abort"))
	fmt.Print(i18n.Get("check_choice_prompt"))

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		// Auto-fix
		if c.config.Checking.AutoFix {
			return c.autoFix(issues)
		}
		output.Warning(i18n.Get("check_autofix_disabled"))
		return fmt.Errorf("auto-fix disabled")

	case 2:
		// View details
		c.showDetails(issues)
		return c.interactiveHandler(issues) // Choose again

	case 3:
		// Ignore
		fmt.Print(i18n.Get("check_ignore_reason"))
		var reason string
		fmt.Scanln(&reason)
		output.Success(i18n.Get("check_ignored"), reason)
		return nil

	case 4:
		// Abort commit
		return fmt.Errorf(i18n.Get("check_aborted"))

	default:
		output.Error(i18n.Get("check_invalid_choice"))
		return c.interactiveHandler(issues)
	}
}

// autoFix auto fix
func (c *Checker) autoFix(issues []models.Issue) error {
	fixed := 0
	for _, issue := range issues {
		// TODO: Implement auto-fix logic
		// Can perform auto-fix based on pattern here
		fmt.Printf(i18n.Get("check_autofix_not_impl")+"\n", issue.File, issue.Line)
	}

	if fixed > 0 {
		fmt.Printf(i18n.Get("check_fixed_count")+"\n", fixed)
	}

	return fmt.Errorf(i18n.Get("check_autofix_not_ready"))
}

// showDetails show details
func (c *Checker) showDetails(issues []models.Issue) {
	for _, issue := range issues {
		fmt.Printf("\n--- %s:%d ---\n", issue.File, issue.Line)
		fmt.Printf("Severity: %s\n", issue.Severity)
		fmt.Printf("Message: %s\n", issue.Message)
		fmt.Printf("Suggestion: %s\n", issue.Suggestion)
	}
}

// Close close checker
func (c *Checker) Close() error {
	return c.store.Close()
}

// AnalyzeFiles analyze specific files
func (c *Checker) AnalyzeFiles(filePaths []string) error {
	// Prepare file changes
	fileChanges := make([]models.FileChange, 0, len(filePaths))

	for _, filePath := range filePaths {
		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf(i18n.Get("check_read_file_failed")+"\n", filePath, err)
			continue
		}

		// Get relative path for display
		relPath, err := filepath.Rel(c.projectRoot, filePath)
		if err != nil {
			relPath = filePath
		}

		fileChanges = append(fileChanges, models.FileChange{
			Path:    relPath,
			Content: string(content),
			Diff:    "", // No diff for direct file analysis
		})
	}

	if len(fileChanges) == 0 {
		fmt.Println(i18n.Get("check_no_valid_files"))
		return nil
	}

	// Run basic checks
	basicIssues := c.runBasicChecks(fileChanges)

	// Claude deep analysis
	var claudeIssues []models.Issue
	if c.config.Claude.Enabled && c.claude.IsAvailable() {
		fmt.Print(i18n.Get("check_analyzing_claude"))
		fmt.Println("")

		context := c.buildAnalysisContext()
		result, err := c.claude.AnalyzeCode(fileChanges, context)
		if err != nil {
			fmt.Printf(i18n.Get("check_claude_failed"), err)
			fmt.Println("")
		} else if result != nil {
			claudeIssues = result.Issues
		}
	}

	// Merge results
	allIssues := append(basicIssues, claudeIssues...)

	// Handle results
	if len(allIssues) == 0 {
		fmt.Print(i18n.Get("check_no_issues"))
		fmt.Println("")
		return nil
	}

	return c.handleIssues(allIssues)
}
