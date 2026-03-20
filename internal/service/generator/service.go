package generator

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/openclaw-coding/skill-seed/internal/domain"
	"github.com/openclaw-coding/skill-seed/internal/templates/skills"
)

// Service 生成服务
type Service struct {
	patternRepo  domain.PatternRepository
	skillsLoader *skills.Loader
}

// NewService 创建生成服务
func NewService(
	patternRepo domain.PatternRepository,
	skillsLoader *skills.Loader,
) *Service {
	return &Service{
		patternRepo:  patternRepo,
		skillsLoader: skillsLoader,
	}
}

// GenerateSkills 生成 Skills 文件夹
func (s *Service) GenerateSkills(ctx context.Context, outputPath string) error {
	// 1. 获取所有模式
	patterns, err := s.patternRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// 2. 计算统计信息
	stats := s.calculateStats(patterns)

	// 3. 准备模板数据（使用大写键名，与模板变量名对应）
	data := map[string]interface{}{
		"TIMESTAMP":               time.Now(),
		"PROJECT_NAME":            "project",
		"PATTERN_COUNT":           len(patterns),
		"AVG_CONFIDENCE":          stats.AvgConfidence,
		"HIGH_CONFIDENCE_PATTERMS": stats.HighConfidence,
		"FREQUENT_PATTERNS":       stats.Frequent,
		"RECENT_PATTERNS":         stats.HighConfidence, // 暂时用高置信度模式
		"ALWAYS_FOLLOW":           stats.Frequent,
		"NEVER_DO":                []domain.Pattern{},
		"TOTAL_COMMITS":           0,
		"TOTAL_PATTERNS":          len(patterns),
		"PATTERN_CATEGORIES":      len(stats.ByCategory),
		"ACTIVE_FILES":            "",
		"GIT_REMOTE":              "",
		"LAST_LEARN_TIME":         time.Now(),
		"FILE_NAMING_PATTERN":     stats.FileNamingPattern,
		"ERROR_CHECK_PATTERN":     stats.ErrorCheckPattern,
		"DIRECTORY_STRUCTURE":     stats.DirectoryStructure,
		"GOROUTINE_PATTERN":       stats.GoroutinePattern,
		"TEST_NAMING_PATTERN":     stats.TestNamingPattern,
	}

	// 4. 确保输出目录存在
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	// 5. 生成主 SKILL.md 文件
	mainContent, err := s.skillsLoader.Render("skill", data)
	if err != nil {
		return err
	}

	mainPath := filepath.Join(outputPath, "SKILL.md")
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return err
	}

	// 6. 生成 references 文件夹中的所有文件
	if err := s.generateAllReferences(outputPath, data); err != nil {
		return err
	}

	return nil
}

// Stats 统计信息
type Stats struct {
	Total             int
	AvgConfidence     float64
	HighConfidence    []domain.Pattern
	Frequent          []domain.Pattern
	ByCategory        map[string][]domain.Pattern
	FileNamingPattern string
	ErrorCheckPattern string
	DirectoryStructure string
	GoroutinePattern  string
	TestNamingPattern string
}

// calculateStats 计算统计信息
func (s *Service) calculateStats(patterns []domain.Pattern) *Stats {
	stats := &Stats{
		Total:          len(patterns),
		ByCategory:     make(map[string][]domain.Pattern),
	}

	if len(patterns) == 0 {
		return stats
	}

	// 计算平均置信度
	var totalConfidence float64
	for _, p := range patterns {
		totalConfidence += p.Confidence
	}
	stats.AvgConfidence = totalConfidence / float64(len(patterns))

	// 按分类统计
	for _, p := range patterns {
		category := string(p.Category)
		stats.ByCategory[category] = append(stats.ByCategory[category], p)
	}

	// 筛选高置信度模式（>0.8）
	for _, p := range patterns {
		if p.Confidence > 0.8 {
			stats.HighConfidence = append(stats.HighConfidence, p)
		}
	}

	// 筛选频繁模式（>3次）
	for _, p := range patterns {
		if p.Frequency > 3 {
			stats.Frequent = append(stats.Frequent, p)
		}
	}

	// 提取特定模式
	stats.FileNamingPattern = s.extractPattern(patterns, "naming", "file")
	stats.ErrorCheckPattern = s.extractPattern(patterns, "error", "check")
	stats.DirectoryStructure = s.extractPattern(patterns, "structure", "directory")
	stats.GoroutinePattern = s.extractPattern(patterns, "concurrency", "goroutine")
	stats.TestNamingPattern = s.extractPattern(patterns, "testing", "naming")

	return stats
}

// extractPattern 提取特定模式
func (s *Service) extractPattern(patterns []domain.Pattern, category, keyword string) string {
	for _, p := range patterns {
		if string(p.Category) == category {
			return p.Rule
		}
	}
	return ""
}

// generateAllReferences 生成所有 references 文件
func (s *Service) generateAllReferences(outputPath string, data map[string]interface{}) error {
	// 定义所有需要生成的分类和文件
	// 注意：category 名称不带 "-patterns" 后缀，LoadReference 会自动添加
	categories := map[string][]string{
		"naming": {
			"overview",
			"file-naming",
			"variable-naming",
			"function-naming",
			"interface-naming",
			"package-naming",
		},
		"error-handling": {
			"overview",
			"error-checking",
			"error-wrapping",
			"error-types",
			"error-logging",
			"error-recovery",
		},
		"structure": {
			"overview",
			"project-layout",
			"package-organization",
			"file-structure",
			"layer-architecture",
		},
		"concurrency": {
			"overview",
			"goroutine-usage",
			"channel-patterns",
			"synchronization",
			"context-usage",
		},
		"testing": {
			"overview",
			"test-organization",
			"test-structure",
			"assertions",
			"mocking",
		},
	}

	// 为每个分类生成文件
	for category, files := range categories {
		categoryPath := filepath.Join(outputPath, "references", category+"-patterns")
		if err := os.MkdirAll(categoryPath, 0755); err != nil {
			return err
		}

		for _, file := range files {
			// 尝试加载并渲染模板
			content, err := s.skillsLoader.RenderReference(category, file, data)
			if err != nil {
				// 如果模板不存在，跳过
				continue
			}

			// 写入文件
			outputFile := filepath.Join(categoryPath, file+".md")
			if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
