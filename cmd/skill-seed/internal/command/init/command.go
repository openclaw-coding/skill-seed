package initcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openclaw-coding/skill-seed/internal/i18n"
	"github.com/openclaw-coding/skill-seed/internal/infra/config"
	"github.com/openclaw-coding/skill-seed/internal/infra/storage/boltdb"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: i18n.Get("InitShort"),
		Long:  i18n.Get("InitLongDesc"),
		Run: func(cmd *cobra.Command, args []string) {
			if err := initializeSkill(); err != nil {
				fmt.Println(i18n.GetWithParams("InitFailed", map[string]interface{}{"Error": err.Error()}))
				os.Exit(1)
			}
		},
	}

	return initCmd
}

func initializeSkill() error {
	// 获取项目根目录
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 检查是否是 Git 仓库
	gitDir := filepath.Join(projectRoot, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", projectRoot)
	}

	// 检查是否已经初始化
	skillPath := filepath.Join(projectRoot, ".skill-seed")
	if _, err := os.Stat(skillPath); err == nil {
		return fmt.Errorf("skill-seed already initialized")
	}

	fmt.Println(i18n.Get("InitStart"))

	// 1. 创建目录结构
	dirs := []string{
		skillPath,
		filepath.Join(skillPath, "memory"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 2. 生成配置
	configRepo, err := config.NewRepository(skillPath)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// 获取项目名称（从目录名）
	projectName := filepath.Base(projectRoot)
	_ = configRepo.SetProjectName(projectName)

	// 3. 初始化数据库
	dbPath := filepath.Join(skillPath, "memory", "project.db")
	patternRepo, err := boltdb.NewPatternRepository(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	_ = patternRepo.Close()

	fmt.Println(i18n.Get("InitSuccess"))
	fmt.Println(i18n.GetWithParams("InitSkillLocation", map[string]interface{}{"Path": skillPath}))
	fmt.Println(i18n.Get("InitNextSteps"))
	fmt.Println(i18n.Get("InitStepLearn"))
	fmt.Println(i18n.Get("InitStepWatch"))
	fmt.Println(i18n.Get("InitStepPatterns"))
	fmt.Println(i18n.Get("InitStepRules"))

	return nil
}
