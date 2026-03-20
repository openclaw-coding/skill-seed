package container

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/openclaw-coding/skill-seed/internal/agent"
	"github.com/openclaw-coding/skill-seed/internal/agent/claude"
	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/infra/config"
	"github.com/openclaw-coding/skill-seed/internal/infra/git"
	"github.com/openclaw-coding/skill-seed/internal/infra/storage/boltdb"
	"github.com/openclaw-coding/skill-seed/internal/service/checker"
	"github.com/openclaw-coding/skill-seed/internal/service/generator"
	"github.com/openclaw-coding/skill-seed/internal/service/learner"
	"github.com/openclaw-coding/skill-seed/internal/templates/prompts"
	"github.com/openclaw-coding/skill-seed/internal/templates/skills"
)

// Container 应用容器
type Container struct {
	Config       *config.Config
	ConfigRepo   *config.Repository
	GitRepo      *git.Repository
	PatternRepo  *boltdb.PatternRepository
	RuleRepo     *boltdb.RuleRepository
	Agent        agent.Agent
	LearnerSvc   *learner.Service
	CheckerSvc   *checker.Service
	GeneratorSvc *generator.Service
	PromptLoader *prompts.Loader
	SkillsLoader *skills.Loader
}

// NewContainer 创建应用容器
func NewContainer(ctx context.Context, seedPath string) (*Container, error) {
	// 1. 加载配置
	configRepo, err := config.NewRepository(seedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	cfg := configRepo.Get()

	// 2. 初始化 i18n
	locale := cfg.Project.Locale
	if locale == "" {
		locale = "zh-CN" // 默认中文
	}
	if err := i18n.Init(locale); err != nil {
		return nil, fmt.Errorf("failed to init i18n: %w", err)
	}

	// 3. 创建 Git 仓储
	projectRoot := filepath.Dir(filepath.Dir(seedPath))
	gitRepo := git.NewRepository(projectRoot)

	// 4. 创建 BoltDB 仓储
	dbPath := filepath.Join(seedPath, "memory", "project.db")
	patternRepo, err := boltdb.NewPatternRepository(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create pattern repository: %w", err)
	}

	ruleRepo := boltdb.NewRuleRepository(patternRepo.GetDB())

	// 5. 创建加载器
	promptLoader := prompts.NewLoader("claude")
	skillsLoader := skills.NewLoader()

	// 6. 创建 Agent
	claudeAgent := claude.New(
		cfg.Agent.Command,
		time.Duration(cfg.Agent.Timeout)*time.Second,
		promptLoader,
	)

	// 7. 创建服务
	learnerSvc := learner.NewService(claudeAgent, gitRepo, patternRepo)
	checkerSvc := checker.NewService(claudeAgent, gitRepo, patternRepo)
	generatorSvc := generator.NewService(patternRepo, skillsLoader)

	return &Container{
		Config:       cfg,
		ConfigRepo:   configRepo,
		GitRepo:      gitRepo,
		PatternRepo:  patternRepo,
		RuleRepo:     ruleRepo,
		Agent:        claudeAgent,
		LearnerSvc:   learnerSvc,
		CheckerSvc:   checkerSvc,
		GeneratorSvc: generatorSvc,
		PromptLoader: promptLoader,
		SkillsLoader: skillsLoader,
	}, nil
}

// Close 关闭容器
func (c *Container) Close() error {
	if c.PatternRepo != nil {
		return c.PatternRepo.Close()
	}
	return nil
}

// GetGitRepository 获取 Git 仓储
func (c *Container) GetGitRepository() *git.Repository {
	return c.GitRepo
}

// GetPatternRepository 获取 Pattern 仓储
func (c *Container) GetPatternRepository() *boltdb.PatternRepository {
	return c.PatternRepo
}

// GetLearnerService 获取学习服务
func (c *Container) GetLearnerService() *learner.Service {
	return c.LearnerSvc
}

// GetCheckerService 获取检查服务
func (c *Container) GetCheckerService() *checker.Service {
	return c.CheckerSvc
}

// GetGeneratorService 获取生成服务
func (c *Container) GetGeneratorService() *generator.Service {
	return c.GeneratorSvc
}
