package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/skill-seed/internal/infra/config"
	"gopkg.in/yaml.v3"
)

// GetSeedPath 获取 .skill-seed 目录路径
func GetSeedPath() (string, error) {
	// 从当前目录开始向上查找
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// 检查当前目录是否有 .skill-seed
	seedPath := filepath.Join(currentDir, ".skill-seed")
	if _, err := os.Stat(seedPath); err == nil {
		return seedPath, nil
	}

	// 向上查找父目录
	parentDir := filepath.Dir(currentDir)
	for {
		seedPath = filepath.Join(parentDir, ".skill-seed")
		if _, err := os.Stat(seedPath); err == nil {
			return seedPath, nil
		}

		// 检查是否到达根目录
		if parentDir == "/" || parentDir == currentDir {
			break
		}

		parentDir = filepath.Dir(parentDir)
	}

	return "", fmt.Errorf(".skill-seed directory not found, please run 'skill-seed init' first")
}

// LoadConfig 加载配置文件（不创建 Container）
func LoadConfig(seedPath string) (*config.Config, error) {
	configPath := filepath.Join(seedPath, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
