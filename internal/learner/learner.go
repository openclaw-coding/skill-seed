package learner

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/openclaw-coding/grow-check/internal/claude"
	"github.com/openclaw-coding/grow-check/internal/config"
	"github.com/openclaw-coding/grow-check/internal/git"
	"github.com/openclaw-coding/grow-check/internal/i18n"
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

// LearnFromHistory learn from history commits
func (l *Learner) LearnFromHistory(sinceDays int, maxCommits int, force bool) error {
	if maxCommits == 0 {
		maxCommits = l.config.Learning.MaxHistoryAnalyze
	}

	// Get last learn time
	lastLearn, _ := l.store.GetLastLearnTime()

	// Calculate start time
	var since time.Time
	if sinceDays > 0 {
		since = time.Now().AddDate(0, 0, -sinceDays)
	} else if !lastLearn.IsZero() && !force {
		since = lastLearn
	}

	// Get new commits
	commits, err := l.git.GetRecentCommits(maxCommits, since)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Println(i18n.Get("learn_no_commits"))
		if !lastLearn.IsZero() {
			fmt.Printf(i18n.Get("learn_last_learn_time")+"\n", lastLearn.Format("2006-01-02 15:04:05"))
		}
		if !since.IsZero() {
			fmt.Printf(i18n.Get("learn_since_time")+"\n", since.Format("2006-01-02 15:04:05"))
		}
		fmt.Println(i18n.Get("learn_tip_force"))
		return nil
	}

	fmt.Printf(i18n.Get("learn_analyzing_commits")+"\n", len(commits))

	// Load existing patterns
	existingPatterns, err := l.store.GetAllPatterns()
	if err != nil {
		existingPatterns = []models.CodePattern{}
	}

	// Learn count
	learnedCount := 0
	skippedCount := 0
	newPatterns := make([]models.CodePattern, 0)

	// Analyze each commit
	for i, commit := range commits {
		// Check if already learned (skip if force mode)
		if !force {
			if learned, _ := l.store.HasLearned(commit.Hash); learned {
				skippedCount++
				continue
			}
		}

		fmt.Printf(i18n.Get("learn_analyzing_commit")+"\n", i+1, len(commits), commit.Hash[:8])

		// Get commit diff
		diff, err := l.git.GetCommitDiff(commit.Hash)
		if err != nil {
			fmt.Printf(i18n.Get("learn_get_diff_failed")+"\n", err)
			continue
		}

		// Use Claude to learn
		if l.config.Claude.Enabled && l.claude.IsAvailable() {
			patterns, err := l.claude.LearnFromCommit(&commit, diff, existingPatterns)
			if err != nil {
				fmt.Printf(i18n.Get("learn_learning_failed")+"\n", err)
				continue
			}

			// Save new patterns
			for _, pattern := range patterns {
				pattern.ID = uuid.New().String()
				pattern.CreatedAt = time.Now()
				pattern.UpdatedAt = time.Now()

				if err := l.store.SavePattern(&pattern); err != nil {
					fmt.Printf(i18n.Get("learn_save_pattern_failed")+"\n", err)
					continue
				}

				newPatterns = append(newPatterns, pattern)
				existingPatterns = append(existingPatterns, pattern)
				learnedCount++
			}
		}

		// Mark as learned
		l.store.SaveLearnRecord(commit.Hash, time.Now())
	}

	// Update learn time
	l.store.UpdateLastLearnTime(time.Now())

	// Show summary
	fmt.Printf(i18n.Get("learn_summary")+"\n")
	fmt.Printf(i18n.Get("learn_total_commits")+"\n", len(commits))
	if skippedCount > 0 {
		fmt.Printf(i18n.Get("learn_skipped_commits")+"\n", skippedCount)
	}
	fmt.Printf(i18n.Get("learn_analyzed_count")+"\n", len(commits)-skippedCount)
	fmt.Printf(i18n.Get("learn_new_patterns")+"\n", learnedCount)

	// Generate rules
	if len(newPatterns) > 0 {
		l.generateRules(newPatterns)
	}

	return nil
}

// generateRules generate rules from patterns
func (l *Learner) generateRules(patterns []models.CodePattern) error {
	fmt.Println(i18n.Get("learn_generating_rules"))

	rulesCreated := 0
	for _, pattern := range patterns {
		// Only create rules for high-frequency patterns
		if pattern.Frequency < l.config.Learning.MinSamplesForRule {
			continue
		}

		// Create rule
		rule := models.Rule{
			ID:         uuid.New().String(),
			Name:       fmt.Sprintf("Learned: %s", pattern.Description),
			Type:       pattern.Type,
			Condition:  pattern.Description, // Simplified version
			Severity:   "warning",
			AutoFix:    pattern.AutoFixable,
			Confidence: pattern.Confidence,
			Source:     "learned",
		}

		if err := l.store.SaveRule(&rule); err != nil {
			fmt.Printf(i18n.Get("learn_create_rule_failed")+"\n", err)
			continue
		}

		rulesCreated++
	}

	if rulesCreated > 0 {
		fmt.Printf(i18n.Get("learn_rules_created")+"\n", rulesCreated)
	} else {
		fmt.Println(i18n.Get("learn_no_rules_created"))
	}

	return nil
}

// ListPatterns list learned patterns
func (l *Learner) ListPatterns() error {
	patterns, err := l.store.GetAllPatterns()
	if err != nil {
		return err
	}

	if len(patterns) == 0 {
		fmt.Println(i18n.Get("learn_no_patterns"))
		return nil
	}

	fmt.Printf(i18n.Get("learn_patterns_header")+"\n\n", len(patterns))
	for i, pattern := range patterns {
		fmt.Printf(i18n.Get("learn_pattern_item")+"\n", i+1, pattern.Type, pattern.Description)
		fmt.Printf(i18n.Get("learn_pattern_details")+"\n",
			pattern.Confidence, pattern.Frequency, pattern.AutoFixable)
		if len(pattern.Examples) > 0 {
			fmt.Printf(i18n.Get("learn_pattern_example")+"\n", pattern.Examples[0])
		}
		fmt.Println()
	}

	return nil
}

// ListRules list rules
func (l *Learner) ListRules() error {
	rules, err := l.store.GetAllRules()
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		fmt.Println(i18n.Get("learn_no_rules"))
		return nil
	}

	fmt.Printf(i18n.Get("learn_rules_header")+"\n\n", len(rules))
	for i, rule := range rules {
		fmt.Printf(i18n.Get("learn_rule_item")+"\n", i+1, rule.Name, rule.Severity)
		fmt.Printf(i18n.Get("learn_rule_details")+"\n", rule.Source, rule.Confidence)
		fmt.Println()
	}

	return nil
}

// Close close learner
func (l *Learner) Close() error {
	return l.store.Close()
}
