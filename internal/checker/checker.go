package checker

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/openclaw-coding/grow-check/internal/claude"
	"github.com/openclaw-coding/grow-check/internal/config"
	"github.com/openclaw-coding/grow-check/internal/git"
	"github.com/openclaw-coding/grow-check/internal/storage"
	"github.com/openclaw-coding/grow-check/pkg/models"
)

// Checker 代码检查器
type Checker struct {
	config     *config.Config
	store      *storage.Store
	git        *git.GitOperator
	claude     *claude.Client
	skillPath  string
	projectRoot string
}

// New 创建检查器
func New(skillPath, projectRoot string) (*Checker, error) {
	// 加载配置
	cfg, err := config.Load(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 打开数据库
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	store, err := storage.New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}

	// 创建 Git 操作器
	gitOp := git.NewGitOperator(projectRoot)

	// 创建 Claude 客户端
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

// Run 运行检查
func (c *Checker) Run() error {
	// 1. 获取暂存文件
	files, err := c.git.GetStagedFiles()
	if err != nil {
		return fmt.Errorf("failed to get staged files: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	fmt.Printf("🔍 Checking %d files...\n", len(files))

	// 2. 过滤文件
	files = c.filterFiles(files)
	if len(files) == 0 {
		return nil
	}

	// 3. 准备文件数据
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

	// 4. 基础检查
	basicIssues := c.runBasicChecks(fileChanges)

	// 5. Claude 深度分析
	var claudeIssues []models.Issue
	if c.config.Claude.Enabled && c.claude.IsAvailable() {
		fmt.Println("🤖 Analyzing with Claude...")
		
		context := c.buildAnalysisContext()
		result, err := c.claude.AnalyzeCode(fileChanges, context)
		if err != nil {
			fmt.Printf("⚠ Claude analysis failed: %v\n", err)
		} else if result != nil {
			claudeIssues = result.Issues
		}
	}

	// 6. 合并结果
	allIssues := append(basicIssues, claudeIssues...)

	// 7. 处理结果
	if len(allIssues) == 0 {
		fmt.Println("✅ No issues found")
		return nil
	}

	return c.handleIssues(allIssues)
}

// filterFiles 过滤文件
func (c *Checker) filterFiles(files []string) []string {
	filtered := make([]string, 0, len(files))
	
	for _, file := range files {
		// 只检查 Go 文件
		if !strings.HasSuffix(file, ".go") {
			continue
		}

		// 检查排除模式
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

// runBasicChecks 运行基础检查
func (c *Checker) runBasicChecks(files []models.FileChange) []models.Issue {
	var issues []models.Issue

	// 加载规则
	rules, err := c.store.GetAllRules()
	if err != nil {
		return issues
	}

	// 对每个文件应用规则
	for _, file := range files {
		for _, rule := range rules {
			if rule.Type == models.PatternNaming {
				// 简单的命名检查
				issues = append(issues, c.checkNaming(file, rule)...)
			}
			// 可以添加更多规则类型
		}
	}

	return issues
}

// checkNaming 检查命名规范
func (c *Checker) checkNaming(file models.FileChange, rule models.Rule) []models.Issue {
	var issues []models.Issue
	
	// 简单示例：检查是否包含下划线（Go 推荐驼峰命名）
	if strings.Contains(file.Path, "_") && !strings.HasSuffix(file.Path, "_test.go") {
		issues = append(issues, models.Issue{
			File:       file.Path,
			Line:       1,
			Severity:   "warning",
			Message:    "文件名包含下划线，建议使用小写字母和连字符",
			Suggestion: "重命名文件，移除下划线",
			PatternID:  rule.ID,
		})
	}

	return issues
}

// buildAnalysisContext 构建分析上下文
func (c *Checker) buildAnalysisContext() *models.AnalysisContext {
	context := &models.AnalysisContext{
		ProjectType:     "go",
		TeamConventions: "遵循 Go 官方代码规范",
	}

	// 加载学习到的模式
	patterns, err := c.store.GetAllPatterns()
	if err == nil {
		context.LearnedPatterns = patterns
	}

	// 获取最近的提交
	commits, err := c.git.GetRecentCommits(10, time.Time{})
	if err == nil {
		context.RecentCommits = commits
	}

	// 获取历史 bug 模式
	bugs, err := c.store.GetMetadata("historical_bugs")
	if err == nil && len(bugs) > 0 {
		context.HistoricalBugs = strings.Split(string(bugs), "\n")
	}

	return context
}

// handleIssues 处理发现的问题
func (c *Checker) handleIssues(issues []models.Issue) error {
	fmt.Printf("\n⚠ Found %d issues:\n\n", len(issues))
	
	for i, issue := range issues {
		severity := "❌"
		if issue.Severity == "warning" {
			severity = "⚠"
		} else if issue.Severity == "info" {
			severity = "ℹ"
		}

		fmt.Printf("%d. %s %s:%d\n", i+1, severity, issue.File, issue.Line)
		fmt.Printf("   %s\n", issue.Message)
		if issue.Suggestion != "" {
			fmt.Printf("   💡 %s\n", issue.Suggestion)
		}
		fmt.Println()
	}

	// 交互式处理
	if c.config.Checking.Interactive {
		return c.interactiveHandler(issues)
	}

	// 非交互模式：有 error 级别问题就失败
	for _, issue := range issues {
		if issue.Severity == "error" {
			return fmt.Errorf("found %d error-level issues", len(issues))
		}
	}

	return nil
}

// interactiveHandler 交互式处理
func (c *Checker) interactiveHandler(issues []models.Issue) error {
	fmt.Println("Options:")
	fmt.Println("1. Auto-fix (recommended)")
	fmt.Println("2. View details")
	fmt.Println("3. Ignore (with reason)")
	fmt.Println("4. Abort commit")
	fmt.Print("\nYour choice [1-4]: ")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		// 自动修复
		if c.config.Checking.AutoFix {
			return c.autoFix(issues)
		}
		fmt.Println("⚠ Auto-fix is disabled in config")
		return fmt.Errorf("auto-fix disabled")

	case 2:
		// 查看详情
		c.showDetails(issues)
		return c.interactiveHandler(issues) // 重新选择

	case 3:
		// 忽略
		fmt.Print("Please provide a reason: ")
		var reason string
		fmt.Scanln(&reason)
		fmt.Printf("✅ Issues ignored. Reason: %s\n", reason)
		return nil

	case 4:
		// 终止提交
		return fmt.Errorf("commit aborted by user")

	default:
		fmt.Println("Invalid choice")
		return c.interactiveHandler(issues)
	}
}

// autoFix 自动修复
func (c *Checker) autoFix(issues []models.Issue) error {
	fixed := 0
	for _, issue := range issues {
		// TODO: 实现自动修复逻辑
		// 这里可以根据 pattern 进行自动修复
		fmt.Printf("⚠ Auto-fix not implemented for %s:%d\n", issue.File, issue.Line)
	}

	if fixed > 0 {
		fmt.Printf("✅ Fixed %d issues\n", fixed)
	}
	
	return fmt.Errorf("auto-fix not fully implemented yet")
}

// showDetails 显示详细信息
func (c *Checker) showDetails(issues []models.Issue) {
	for _, issue := range issues {
		fmt.Printf("\n--- %s:%d ---\n", issue.File, issue.Line)
		fmt.Printf("Severity: %s\n", issue.Severity)
		fmt.Printf("Message: %s\n", issue.Message)
		fmt.Printf("Suggestion: %s\n", issue.Suggestion)
	}
}

// Close 关闭检查器
func (c *Checker) Close() error {
	return c.store.Close()
}
