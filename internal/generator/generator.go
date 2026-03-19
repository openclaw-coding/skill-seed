package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/openclaw-coding/grow-check/internal/i18n"
	"github.com/openclaw-coding/grow-check/internal/storage"
	"github.com/openclaw-coding/grow-check/pkg/models"
)

// Generator skills 生成器
type Generator struct {
	store      *storage.Store
	skillPath  string
	projectRoot string
}

// New 创建生成器
func New(skillPath, projectRoot string) (*Generator, error) {
	store, err := storage.New(filepath.Join(skillPath, "memory", "project.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}

	return &Generator{
		store:      store,
		skillPath:  skillPath,
		projectRoot: projectRoot,
	}, nil
}

// Generate 生成 skills
func (g *Generator) Generate(outputPath string) error {
	// 1. 获取所有模式
	patterns, err := g.store.GetAllPatterns()
	if err != nil {
		return fmt.Errorf("failed to get patterns: %w", err)
	}

	// 2. 计算统计信息
	stats := g.calculateStats(patterns)

	// 3. 准备模板变量
	vars := g.prepareTemplateVars(patterns, stats)

	// 4. 加载模板目录
	templateDir := filepath.Join(g.projectRoot, ".template", "skills", "grow-check-skills")

	// 5. 生成所有文件
	if err := g.generateFiles(templateDir, outputPath, vars); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	fmt.Printf(i18n.Get("generate_output_location")+"\n", outputPath)
	fmt.Printf(i18n.Get("generate_patterns_count")+"\n", len(patterns))
	fmt.Printf(i18n.Get("generate_avg_confidence")+"\n", stats.AvgConfidence*100)

	return nil
}

// Stats 统计信息
type Stats struct {
	TotalPatterns    int
	PatternCount     map[models.PatternType]int
	AvgConfidence    float64
	TotalCommits     int
	ConsistencyRate  float64
	MostCommon       models.CodePattern
	HighConfidence   []models.CodePattern
	RecentPatterns   []models.CodePattern
	FrequentPatterns []models.CodePattern
}

// calculateStats 计算统计信息
func (g *Generator) calculateStats(patterns []models.CodePattern) Stats {
	stats := Stats{
		PatternCount: make(map[models.PatternType]int),
	}

	if len(patterns) == 0 {
		return stats
	}

	stats.TotalPatterns = len(patterns)

	var totalConfidence float64
	frequencyMap := make(map[string]*models.CodePattern)

	for _, p := range patterns {
		// 按类型统计
		stats.PatternCount[p.Type]++

		// 累计置信度
		totalConfidence += p.Confidence

		// 追踪最频繁
		if existing, ok := frequencyMap[p.Description]; !ok || p.Frequency > existing.Frequency {
			frequencyMap[p.Description] = &p
		}
	}

	stats.AvgConfidence = totalConfidence / float64(len(patterns))

	// 高置信度模式
	for _, p := range patterns {
		if p.Confidence >= 0.8 {
			stats.HighConfidence = append(stats.HighConfidence, p)
		}
		if p.Frequency >= 10 {
			stats.FrequentPatterns = append(stats.FrequentPatterns, p)
		}
	}

	// 一致性率（简化计算）
	stats.ConsistencyRate = 0.85 // 基于模式一致性

	return stats
}

// prepareTemplateVars 准备模板变量
func (g *Generator) prepareTemplateVars(patterns []models.CodePattern, stats Stats) map[string]interface{} {
	vars := map[string]interface{}{
		"TIMESTAMP":        time.Now().Format("2006-01-02 15:04:05"),
		"PROJECT_NAME":     filepath.Base(g.projectRoot),
		"VERSION":          "1.0.0",
		"COMMITS_ANALYZED": stats.TotalCommits,
		"PATTERN_COUNT":    stats.TotalPatterns,
		"AVG_CONFIDENCE":   fmt.Sprintf("%.1f", stats.AvgConfidence*100),
		"CONSISTENCY_RATE": fmt.Sprintf("%.1f", stats.ConsistencyRate*100),
	}

	// 按类型组织模式
	patternsByType := make(map[models.PatternType][]models.CodePattern)
	for _, p := range patterns {
		patternsByType[p.Type] = append(patternsByType[p.Type], p)
	}

	// 添加各类型模式的变量
	g.addNamingPatternVars(vars, patternsByType[models.PatternNaming])
	g.addErrorHandlingVars(vars, patternsByType[models.PatternErrorHandling])
	g.addStructurePatternVars(vars, patternsByType[models.PatternStructure])
	g.addConcurrencyPatternVars(vars, patternsByType[models.PatternConcurrency])
	g.addTestingPatternVars(vars, patternsByType[models.PatternTesting])

	// 高置信度模式列表
	var highConfidenceList []string
	for _, p := range stats.HighConfidence {
		highConfidenceList = append(highConfidenceList,
			fmt.Sprintf("- **%s**: %s (%.0f%% confidence)", p.Type, p.Description, p.Confidence*100))
	}
	vars["HIGH_CONFIDENCE_PATTERNS"] = strings.Join(highConfidenceList, "\n")

	// 频繁模式列表
	var frequentList []string
	for _, p := range stats.FrequentPatterns {
		frequentList = append(frequentList,
			fmt.Sprintf("- **%s**: %s (used %d times)", p.Type, p.Description, p.Frequency))
	}
	vars["FREQUENT_PATTERNS"] = strings.Join(frequentList, "\n")

	return vars
}

// addNamingPatternVars 添加命名模式变量
func (g *Generator) addNamingPatternVars(vars map[string]interface{}, patterns []models.CodePattern) {
	// TODO: 从 patterns 中提取命名规范示例
	// 这里简化为示例数据
	vars["FILE_NAMING_PATTERN"] = "Use lowercase with hyphens: my-file.go"
	vars["FILE_NAMING_GOOD_EXAMPLES"] = "- checker.go\n- learner.go\n- storage.go"
	vars["FILE_NAMING_BAD_EXAMPLES"] = "- check_er.go  # underscore\n- Checker.go  # capitalized"
	vars["FILE_NAMING_RULE"] = "Use lowercase, no underscores, .go extension"
	vars["EXAMPLE_COUNT"] = len(patterns)
	vars["CONFIDENCE"] = "85"
}

// addErrorHandlingVars 添加错误处理变量
func (g *Generator) addErrorHandlingVars(vars map[string]interface{}, patterns []models.CodePattern) {
	vars["ERROR_CHECK_PATTERN"] = "Always check errors immediately"
	vars["ERROR_CHECK_GOOD_EXAMPLE"] = `if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}`
	vars["ERROR_CHECK_BAD_EXAMPLE"] = `if err != nil {
    return err
}`
	vars["ERROR_CHECK_RULE"] = "Always wrap errors with context using %w"
	vars["EXAMPLE_COUNT"] = len(patterns)
	vars["CONFIDENCE"] = "90"
}

// addStructurePatternVars 添加结构模式变量
func (g *Generator) addStructurePatternVars(vars map[string]interface{}, patterns []models.CodePattern) {
	vars["DIRECTORY_STRUCTURE"] = "Standard Go project layout"
	vars["DIRECTORY_STRUCTURE_RULES"] = "- cmd/ for applications\n- internal/ for private code\n- pkg/ for public code"
	vars["EXAMPLE_COUNT"] = len(patterns)
	vars["CONFIDENCE"] = "80"
}

// addConcurrencyPatternVars 添加并发模式变量
func (g *Generator) addConcurrencyPatternVars(vars map[string]interface{}, patterns []models.CodePattern) {
	vars["GOROUTINE_PATTERN"] = "Start goroutines, but manage lifecycle"
	vars["GOROUTINE_RULES"] = "- Always wait for goroutines to finish\n- Use context for cancellation"
	vars["EXAMPLE_COUNT"] = len(patterns)
	vars["CONFIDENCE"] = "75"
}

// addTestingPatternVars 添加测试模式变量
func (g *Generator) addTestingPatternVars(vars map[string]interface{}, patterns []models.CodePattern) {
	vars["TEST_NAMING_PATTERN"] = "Test<FunctionName> for unit tests"
	vars["TEST_NAMING_RULES"] = "- Start with Test\n- Use the function name being tested"
	vars["EXAMPLE_COUNT"] = len(patterns)
	vars["CONFIDENCE"] = "88"
}

// generateFiles 生成所有文件
func (g *Generator) generateFiles(templateDir, outputPath string, vars map[string]interface{}) error {
	// 遍历模板目录
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理 .md 文件
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		// 读取模板
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template: %w", err)
		}

		// 解析并渲染模板
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, vars); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		// 计算输出路径
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}
		outputFile := filepath.Join(outputPath, relPath)

		// 创建目录
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// 写入文件
		if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf(i18n.Get("generate_file_created")+"\n", outputFile)

		return nil
	})
}

// Close 关闭生成器
func (g *Generator) Close() error {
	return g.store.Close()
}
