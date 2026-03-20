package generator

import (
	"context"
	"fmt"
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

// GenerateSkills 生成 Skills 文件
func (s *Service) GenerateSkills(ctx context.Context, outputPath string) error {
	// 1. 获取所有模式
	patterns, err := s.patternRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// 2. 计算统计信息
	stats := s.calculateStats(patterns)

	// 3. 准备模板数据
	data := map[string]interface{}{
		"Timestamp":               time.Now(),
		"ProjectName":             "project",
		"PatternCount":            len(patterns),
		"AvgConfidence":           stats.AvgConfidence,
		"HighConfidencePatterns":  stats.HighConfidence,
		"FrequentPatterns":        stats.Frequent,
		"NamingPatterns":          stats.ByCategory["naming"],
		"ErrorPatterns":           stats.ByCategory["error"],
		"StructurePatterns":       stats.ByCategory["structure"],
		"ConcurrencyPatterns":     stats.ByCategory["concurrency"],
		"TestingPatterns":         stats.ByCategory["testing"],
		"FILE_NAMING_PATTERN":     stats.FileNamingPattern,
		"ERROR_CHECK_PATTERN":     stats.ErrorCheckPattern,
		"DIRECTORY_STRUCTURE":     stats.DirectoryStructure,
		"GOROUTINE_PATTERN":       stats.GoroutinePattern,
		"TEST_NAMING_PATTERN":     stats.TestNamingPattern,
	}

	// 4. 渲染主模板
	content, err := s.skillsLoader.Render("skill", data)
	if err != nil {
		return err
	}

	// 5. 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// 6. 写入文件
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return err
	}

	// 7. 生成 references 文件
	if err := s.generateReferences(ctx, filepath.Dir(outputPath), stats); err != nil {
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

// generateReferences 生成 references 文件
func (s *Service) generateReferences(ctx context.Context, outputDir string, stats *Stats) error {
	// 为每个分类生成 reference 文件
	categories := []string{"naming", "error", "structure", "concurrency", "testing"}

	for _, category := range categories {
		patterns, ok := stats.ByCategory[category]
		if !ok || len(patterns) == 0 {
			continue
		}

		data := map[string]interface{}{
			"Category":  category,
			"Patterns":  patterns,
			"Timestamp": time.Now(),
		}

		// 尝试加载并渲染模板
		content, err := s.skillsLoader.Render(fmt.Sprintf("references/%s/overview", category), data)
		if err != nil {
			// 如果模板不存在，跳过
			continue
		}

		// 写入文件
		refDir := filepath.Join(outputDir, "references", category)
		if err := os.MkdirAll(refDir, 0755); err != nil {
			continue
		}

		outputPath := filepath.Join(refDir, "overview.md")
		_ = os.WriteFile(outputPath, []byte(content), 0644)
	}

	return nil
}
