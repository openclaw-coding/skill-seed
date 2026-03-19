package learner

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/openclaw-coding/grow-check/internal/claude"
	"github.com/openclaw-coding/grow-check/internal/config"
	"github.com/openclaw-coding/grow-check/internal/git"
	"github.com/openclaw-coding/grow-check/internal/storage"
	"github.com/openclaw-coding/grow-check/pkg/models"
)

// Learner 学习器
type Learner struct {
	config    *config.Config
	store     *storage.Store
	git       *git.GitOperator
	claude    *claude.Client
	skillPath string
}

// New 创建学习器
func New(skillPath string) (*Learner, error) {
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
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(skillPath)))
	gitOp := git.NewGitOperator(projectRoot)

	// 创建 Claude 客户端
	claudeClient := claude.NewClient(
		"claude",
		time.Duration(cfg.Claude.TimeoutSeconds)*time.Second,
		cfg.Claude.FallbackToBasic,
	)

	return &Learner{
		config:    cfg,
		store:     store,
		git:       gitOp,
		claude:    claudeClient,
		skillPath: skillPath,
	}, nil
}

// LearnFromHistory 从历史提交学习
func (l *Learner) LearnFromHistory(sinceDays int, maxCommits int) error {
	if maxCommits == 0 {
		maxCommits = l.config.Learning.MaxHistoryAnalyze
	}

	// 获取上次学习时间
	lastLearn, _ := l.store.GetLastLearnTime()

	// 计算起始时间
	var since time.Time
	if sinceDays > 0 {
		since = time.Now().AddDate(0, 0, -sinceDays)
	} else if !lastLearn.IsZero() {
		since = lastLearn
	}

	// 获取新提交
	commits, err := l.git.GetRecentCommits(maxCommits, since)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Println("📚 No new commits to learn from")
		return nil
	}

	fmt.Printf("📚 Analyzing %d commits...\n", len(commits))

	// 加载现有模式
	existingPatterns, err := l.store.GetAllPatterns()
	if err != nil {
		existingPatterns = []models.CodePattern{}
	}

	// 学习计数
	learnedCount := 0
	newPatterns := make([]models.CodePattern, 0)

	// 分析每个提交
	for i, commit := range commits {
		// 检查是否已学习过
		if learned, _ := l.store.HasLearned(commit.Hash); learned {
			continue
		}

		fmt.Printf("  [%d/%d] Analyzing %s...\n", i+1, len(commits), commit.Hash[:8])

		// 获取提交的 diff
		diff, err := l.git.GetCommitDiff(commit.Hash)
		if err != nil {
			continue
		}

		// 使用 Claude 学习
		if l.config.Claude.Enabled && l.claude.IsAvailable() {
			patterns, err := l.claude.LearnFromCommit(&commit, diff, existingPatterns)
			if err != nil {
				fmt.Printf("    ⚠ Learning failed: %v\n", err)
				continue
			}

			// 保存新模式
			for _, pattern := range patterns {
				pattern.ID = uuid.New().String()
				pattern.CreatedAt = time.Now()
				pattern.UpdatedAt = time.Now()

				if err := l.store.SavePattern(&pattern); err != nil {
					fmt.Printf("    ⚠ Failed to save pattern: %v\n", err)
					continue
				}

				newPatterns = append(newPatterns, pattern)
				existingPatterns = append(existingPatterns, pattern)
				learnedCount++
			}
		}

		// 标记为已学习
		l.store.SaveLearnRecord(commit.Hash, time.Now())
	}

	// 更新学习时间
	l.store.UpdateLastLearnTime(time.Now())

	// 生成规则
	if len(newPatterns) > 0 {
		fmt.Printf("\n✨ Learned %d new patterns\n", learnedCount)
		l.generateRules(newPatterns)
	}

	return nil
}

// generateRules 从模式生成规则
func (l *Learner) generateRules(patterns []models.CodePattern) error {
	fmt.Println("\n📏 Generating rules from patterns...")

	rulesCreated := 0
	for _, pattern := range patterns {
		// 只为高频模式创建规则
		if pattern.Frequency < l.config.Learning.MinSamplesForRule {
			continue
		}

		// 创建规则
		rule := models.Rule{
			ID:         uuid.New().String(),
			Name:       fmt.Sprintf("Learned: %s", pattern.Description),
			Type:       pattern.Type,
			Condition:  pattern.Description, // 简化版
			Severity:   "warning",
			AutoFix:    pattern.AutoFixable,
			Confidence: pattern.Confidence,
			Source:     "learned",
		}

		if err := l.store.SaveRule(&rule); err != nil {
			fmt.Printf("  ⚠ Failed to create rule: %v\n", err)
			continue
		}

		rulesCreated++
	}

	if rulesCreated > 0 {
		fmt.Printf("✅ Created %d new rules\n", rulesCreated)
	}

	return nil
}

// ListPatterns 列出学习到的模式
func (l *Learner) ListPatterns() error {
	patterns, err := l.store.GetAllPatterns()
	if err != nil {
		return err
	}

	if len(patterns) == 0 {
		fmt.Println("📋 No patterns learned yet")
		return nil
	}

	fmt.Printf("📋 Learned patterns (%d total):\n\n", len(patterns))
	for i, pattern := range patterns {
		fmt.Printf("%d. [%s] %s\n", i+1, pattern.Type, pattern.Description)
		fmt.Printf("   Confidence: %.2f | Frequency: %d | Auto-fixable: %v\n",
			pattern.Confidence, pattern.Frequency, pattern.AutoFixable)
		if len(pattern.Examples) > 0 {
			fmt.Printf("   Example: %s\n", pattern.Examples[0])
		}
		fmt.Println()
	}

	return nil
}

// ListRules 列出规则
func (l *Learner) ListRules() error {
	rules, err := l.store.GetAllRules()
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		fmt.Println("📏 No rules generated yet")
		return nil
	}

	fmt.Printf("📏 Active rules (%d total):\n\n", len(rules))
	for i, rule := range rules {
		fmt.Printf("%d. %s [%s]\n", i+1, rule.Name, rule.Severity)
		fmt.Printf("   Source: %s | Confidence: %.2f\n", rule.Source, rule.Confidence)
		fmt.Println()
	}

	return nil
}

// Close 关闭学习器
func (l *Learner) Close() error {
	return l.store.Close()
}
