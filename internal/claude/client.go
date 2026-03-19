package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/openclaw-coding/grow-check/pkg/models"
)

// Client Claude 客户端
type Client struct {
	CommandPath string        // claude 命令路径
	Timeout     time.Duration // 超时时间
	Fallback    bool          // 是否降级到基础检查
}

// NewClient 创建 Claude 客户端
func NewClient(commandPath string, timeout time.Duration, fallback bool) *Client {
	if commandPath == "" {
		commandPath = "claude" // 默认命令
	}
	return &Client{
		CommandPath: commandPath,
		Timeout:     timeout,
		Fallback:    fallback,
	}
}

// IsAvailable 检查 Claude 是否可用
func (c *Client) IsAvailable() bool {
	_, err := exec.LookPath(c.CommandPath)
	return err == nil
}

// AnalyzeCode 分析代码
func (c *Client) AnalyzeCode(files []models.FileChange, analysisContext *models.AnalysisContext) (*models.AnalysisResult, error) {
	if !c.IsAvailable() {
		if c.Fallback {
			return nil, nil // 降级到基础检查
		}
		return nil, fmt.Errorf("claude not available")
	}

	// 构建提示词
	prompt := c.buildPrompt(files, analysisContext)

	// 调用 Claude
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.CommandPath, "--print", prompt)
	output, err := cmd.Output()
	if err != nil {
		if c.Fallback {
			return nil, nil // 降级
		}
		return nil, fmt.Errorf("claude analysis failed: %w", err)
	}

	// 解析结果
	result, err := c.parseResult(string(output))
	if err != nil {
		// 如果解析失败，返回基础结果
		return &models.AnalysisResult{
			Issues:     []models.Issue{},
			Summary:    "Claude analysis completed but result parsing failed",
			Confidence: 0.5,
		}, nil
	}

	return result, nil
}

// buildPrompt 构建分析提示词
func (c *Client) buildPrompt(files []models.FileChange, analysisContext *models.AnalysisContext) string {
	var sb strings.Builder

	sb.WriteString("你是一位资深的代码审查专家。请分析以下代码变更：\n\n")

	// 项目上下文
	sb.WriteString("## 项目上下文\n")
	sb.WriteString(fmt.Sprintf("- 项目类型: %s\n", analysisContext.ProjectType))
	
	if len(analysisContext.HistoricalBugs) > 0 {
		sb.WriteString("- 历史 bug 模式:\n")
		for _, bug := range analysisContext.HistoricalBugs {
			sb.WriteString(fmt.Sprintf("  - %s\n", bug))
		}
	}

	if analysisContext.TeamConventions != "" {
		sb.WriteString(fmt.Sprintf("- 团队规范: %s\n", analysisContext.TeamConventions))
	}

	// 学习到的模式
	if len(analysisContext.LearnedPatterns) > 0 {
		sb.WriteString("\n## 学习到的代码模式\n")
		for _, pattern := range analysisContext.LearnedPatterns {
			sb.WriteString(fmt.Sprintf("- [%s] %s (置信度: %.2f)\n", 
				pattern.Type, pattern.Description, pattern.Confidence))
		}
	}

	// 变更文件
	sb.WriteString("\n## 变更文件\n")
	for i, file := range files {
		sb.WriteString(fmt.Sprintf("\n### 文件 %d: %s\n", i+1, file.Path))
		if file.Diff != "" {
			sb.WriteString("```diff\n")
			sb.WriteString(file.Diff)
			sb.WriteString("\n```\n")
		} else {
			sb.WriteString("```\n")
			sb.WriteString(file.Content)
			sb.WriteString("\n```\n")
		}
	}

	// 分析要求
	sb.WriteString(`
## 请分析以下方面：

1. **潜在问题**：是否存在 bug、逻辑错误、性能问题
2. **模式匹配**：是否符合项目的历史编码模式
3. **回归风险**：是否可能引入回归问题
4. **改进建议**：具体的修复或优化建议

## 返回格式

请以 JSON 格式返回，格式如下：
` + "```json\n" + `
{
  "issues": [
    {
      "file": "文件路径",
      "line": 行号,
      "severity": "error|warning|info",
      "message": "问题描述",
      "suggestion": "修复建议"
    }
  ],
  "summary": "整体分析摘要",
  "confidence": 0.85
}
` + "```\n")

	return sb.String()
}

// parseResult 解析 Claude 返回的结果
func (c *Client) parseResult(output string) (*models.AnalysisResult, error) {
	// 尝试提取 JSON 部分
	jsonStart := strings.Index(output, "{")
	jsonEnd := strings.LastIndex(output, "}")
	
	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in output")
	}

	jsonStr := output[jsonStart : jsonEnd+1]

	var result models.AnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &result, nil
}

// LearnFromCommit 从提交学习
func (c *Client) LearnFromCommit(commit *models.CommitInfo, diff string, existingPatterns []models.CodePattern) ([]models.CodePattern, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("claude not available")
	}

	prompt := c.buildLearnPrompt(commit, diff, existingPatterns)

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.CommandPath, "--print", prompt)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("claude learning failed: %w", err)
	}

	// 解析学习到的模式
	patterns, err := c.parsePatterns(string(output))
	if err != nil {
		return nil, err
	}

	return patterns, nil
}

// buildLearnPrompt 构建学习提示词
func (c *Client) buildLearnPrompt(commit *models.CommitInfo, diff string, existingPatterns []models.CodePattern) string {
	var sb strings.Builder

	sb.WriteString("分析以下 Git 提交，识别代码模式和规范：\n\n")
	sb.WriteString(fmt.Sprintf("## 提交信息\n"))
	sb.WriteString(fmt.Sprintf("- Hash: %s\n", commit.Hash))
	sb.WriteString(fmt.Sprintf("- Message: %s\n", commit.Message))
	sb.WriteString(fmt.Sprintf("- Author: %s\n", commit.Author))
	sb.WriteString(fmt.Sprintf("- Time: %s\n", commit.Timestamp.Format(time.RFC3339)))

	sb.WriteString("\n## 代码变更\n")
	sb.WriteString("```diff\n")
	sb.WriteString(diff)
	sb.WriteString("\n```\n")

	if len(existingPatterns) > 0 {
		sb.WriteString("\n## 已知的代码模式\n")
		for _, p := range existingPatterns {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", p.Type, p.Description))
		}
	}

	sb.WriteString(`
## 分析任务

1. 识别新的代码模式（命名、结构、错误处理等）
2. 判断是否是 bug fix，如果是，记录错误模式
3. 提取团队编码习惯

## 返回格式

以 JSON 数组返回发现的新模式：
` + "```json\n" + `
[
  {
    "type": "naming|structure|error_handling|concurrency|testing|comment",
    "description": "模式描述",
    "example": "代码示例",
    "auto_fixable": true
  }
]
` + "```\n")

	return sb.String()
}

// parsePatterns 解析学习到的模式
func (c *Client) parsePatterns(output string) ([]models.CodePattern, error) {
	jsonStart := strings.Index(output, "[")
	jsonEnd := strings.LastIndex(output, "]")
	
	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return []models.CodePattern{}, nil
	}

	jsonStr := output[jsonStart : jsonEnd+1]

	var rawPatterns []struct {
		Type         string `json:"type"`
		Description  string `json:"description"`
		Example      string `json:"example"`
		AutoFixable  bool   `json:"auto_fixable"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &rawPatterns); err != nil {
		return nil, fmt.Errorf("failed to parse patterns: %w", err)
	}

	patterns := make([]models.CodePattern, 0, len(rawPatterns))
	for _, rp := range rawPatterns {
		patterns = append(patterns, models.CodePattern{
			Type:        models.PatternType(rp.Type),
			Description: rp.Description,
			Examples:    []string{rp.Example},
			Frequency:   1,
			Confidence:  0.7, // 初始置信度
			AutoFixable: rp.AutoFixable,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return patterns, nil
}
