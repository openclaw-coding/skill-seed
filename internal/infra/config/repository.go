package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/openclaw-coding/skill-seed/embedfs"
	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Project  ProjectConfig  `yaml:"project"`
	Claude   ClaudeConfig   `yaml:"claude"`
	Agent    AgentConfig    `yaml:"agent"`
	Learning LearningConfig `yaml:"learning"`
	Checking CheckingConfig `yaml:"checking"`
	Output   OutputConfig   `yaml:"output"`
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Name          string    `yaml:"name"`
	Language      string    `yaml:"language"`
	InitializedAt time.Time `yaml:"initialized_at"`
	GitRemote     string    `yaml:"git_remote"`
	RootPath      string    `yaml:"root_path"`
	Locale        string    `yaml:"locale"` // 语言设置：zh-CN, en-US
}

// ClaudeConfig Claude 配置
type ClaudeConfig struct {
	Enabled         bool `yaml:"enabled"`
	TimeoutSeconds  int  `yaml:"timeout_seconds"`
	FallbackToBasic bool `yaml:"fallback_to_basic"`
}

// AgentConfig Agent 配置
type AgentConfig struct {
	Type    string `yaml:"type"` // claude, gpt, local
	Command string `yaml:"command"`
	Timeout int    `yaml:"timeout"`
}

// LearningConfig 学习配置
type LearningConfig struct {
	MaxCommits         int `yaml:"max_commits"`
	MinSamplesForRule  int `yaml:"min_samples_for_rule"`
}

// CheckingConfig 检查配置
type CheckingConfig struct {
	ExcludePatterns []string `yaml:"exclude_patterns"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	SkillsPath       string `yaml:"skills_path"`
	DefaultLanguage  string `yaml:"default_language"`
}

// Repository 配置仓储
type Repository struct {
	configPath string
	config     *Config
}

// NewRepository 创建配置仓储
func NewRepository(seedPath string) (*Repository, error) {
	configPath := filepath.Join(seedPath, "config.yaml")

	repo := &Repository{
		configPath: configPath,
	}

	// 加载配置
	cfg, err := repo.load()
	if err != nil {
		// 如果配置文件不存在，创建默认配置
		var pathErr *os.PathError
		if errors.As(err, &pathErr) || errors.Is(err, os.ErrNotExist) {
			cfg = repo.defaultConfig()
			if err := repo.save(cfg); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	repo.config = cfg
	return repo, nil
}

// Get 获取配置
func (r *Repository) Get() *Config {
	return r.config
}

// Update 更新配置
func (r *Repository) Update(cfg *Config) error {
	if err := r.save(cfg); err != nil {
		return err
	}
	r.config = cfg
	return nil
}

// load 加载配置
func (r *Repository) load() (*Config, error) {
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// save 保存配置
func (r *Repository) save(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(r.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// defaultConfig 默认配置
func (r *Repository) defaultConfig() *Config {
	// 从 embedfs 读取默认配置模板
	data, err := embedfs.FS.ReadFile("templates/config/config.yaml.tmpl")
	if err != nil {
		// 如果读取失败，使用硬编码的默认值
		return r.hardcodedDefaultConfig()
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		// 如果解析失败，使用硬编码的默认值
		return r.hardcodedDefaultConfig()
	}

	// 设置初始化时间
	cfg.Project.InitializedAt = time.Now()

	return &cfg
}

// hardcodedDefaultConfig 硬编码的默认配置（作为后备）
func (r *Repository) hardcodedDefaultConfig() *Config {
	return &Config{
		Project: ProjectConfig{
			Name:          "project",
			Language:      "go",
			InitializedAt: time.Now(),
			Locale:        "zh-CN",
		},
		Claude: ClaudeConfig{
			Enabled:         true,
			TimeoutSeconds:  60,
			FallbackToBasic: true,
		},
		Agent: AgentConfig{
			Type:    "claude",
			Command: "claude",
			Timeout: 60,
		},
		Learning: LearningConfig{
			MaxCommits:        50,
			MinSamplesForRule: 3,
		},
		Checking: CheckingConfig{
			ExcludePatterns: []string{
				"vendor/*",
				"node_modules/*",
				"*.pb.go",
				"*.gen.go",
				"*/mocks/*",
				"**/testdata/*",
				"**/test/*",
			},
		},
		Output: OutputConfig{
			SkillsPath:      ".claude/skills/skill-seed-skills",
			DefaultLanguage: "go",
		},
	}
}

// GetProjectConfig 获取项目配置
func (r *Repository) GetProjectConfig() ProjectConfig {
	return r.config.Project
}

// GetClaudeConfig 获取 Claude 配置
func (r *Repository) GetClaudeConfig() ClaudeConfig {
	return r.config.Claude
}

// GetAgentConfig 获取 Agent 配置
func (r *Repository) GetAgentConfig() AgentConfig {
	return r.config.Agent
}

// GetLearningConfig 获取学习配置
func (r *Repository) GetLearningConfig() LearningConfig {
	return r.config.Learning
}

// GetOutputConfig 获取输出配置
func (r *Repository) GetOutputConfig() OutputConfig {
	return r.config.Output
}

// GetCheckingConfig 获取检查配置
func (r *Repository) GetCheckingConfig() CheckingConfig {
	return r.config.Checking
}

// SetProjectName 设置项目名称
func (r *Repository) SetProjectName(name string) error {
	r.config.Project.Name = name
	return r.Update(r.config)
}

// SetProjectLanguage 设置项目语言
func (r *Repository) SetProjectLanguage(language string) error {
	r.config.Project.Language = language
	return r.Update(r.config)
}

// SetGitRemote 设置 Git Remote
func (r *Repository) SetGitRemote(gitRemote string) error {
	r.config.Project.GitRemote = gitRemote
	return r.Update(r.config)
}

// SetRootPath 设置根路径
func (r *Repository) SetRootPath(rootPath string) error {
	r.config.Project.RootPath = rootPath
	return r.Update(r.config)
}

// SetLocale 设置语言
func (r *Repository) SetLocale(locale string) error {
	r.config.Project.Locale = locale
	return r.Update(r.config)
}
