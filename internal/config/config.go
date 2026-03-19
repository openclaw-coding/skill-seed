package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	Project struct {
		Name          string    `yaml:"name"`
		InitializedAt time.Time `yaml:"initialized_at"`
		GitRemote     string    `yaml:"git_remote"`
	} `yaml:"project"`

	Claude struct {
		Enabled         bool `yaml:"enabled"`
		TimeoutSeconds  int  `yaml:"timeout_seconds"`
		FallbackToBasic bool `yaml:"fallback_to_basic"`
	} `yaml:"claude"`

	Learning struct {
		MinSamplesForRule     int `yaml:"min_samples_for_rule"`
		AutoLearnIntervalDays int `yaml:"auto_learn_interval_days"`
		MaxHistoryAnalyze     int `yaml:"max_history_analyze"`
	} `yaml:"learning"`

	Checking struct {
		Interactive     bool     `yaml:"interactive"`
		AutoFix         bool     `yaml:"auto_fix"`
		SeverityLevels  []string `yaml:"severity_levels"`
		ExcludePatterns []string `yaml:"exclude_patterns"`
	} `yaml:"checking"`
}

// Load 加载配置
func Load(skillPath string) (*Config, error) {
	configPath := filepath.Join(skillPath, "config.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save 保存配置
func (c *Config) Save(skillPath string) error {
	configPath := filepath.Join(skillPath, "config.yaml")
	
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// DefaultConfig 生成默认配置
func DefaultConfig(projectName, gitRemote string) *Config {
	cfg := &Config{}

	cfg.Project.Name = projectName
	cfg.Project.InitializedAt = time.Now()
	cfg.Project.GitRemote = gitRemote

	cfg.Claude.Enabled = true
	cfg.Claude.TimeoutSeconds = 30
	cfg.Claude.FallbackToBasic = true

	cfg.Learning.MinSamplesForRule = 3
	cfg.Learning.AutoLearnIntervalDays = 7
	cfg.Learning.MaxHistoryAnalyze = 1000

	cfg.Checking.Interactive = true
	cfg.Checking.AutoFix = true
	cfg.Checking.SeverityLevels = []string{"error", "warning", "info"}
	cfg.Checking.ExcludePatterns = []string{
		"vendor/*",
		"node_modules/*",
		"*.pb.go",
		"*.gen.go",
	}

	return cfg
}
