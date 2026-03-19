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

// buildPrompt build analysis prompt
func (c *Client) buildPrompt(files []models.FileChange, analysisContext *models.AnalysisContext) string {
	var sb strings.Builder

	sb.WriteString("You are a senior code review expert. Please analyze the following code changes:\n\n")

	// Project context
	sb.WriteString("## Project Context\n")
	sb.WriteString(fmt.Sprintf("- Project Type: %s\n", analysisContext.ProjectType))

	if len(analysisContext.HistoricalBugs) > 0 {
		sb.WriteString("- Historical Bug Patterns:\n")
		for _, bug := range analysisContext.HistoricalBugs {
			sb.WriteString(fmt.Sprintf("  - %s\n", bug))
		}
	}

	if analysisContext.TeamConventions != "" {
		sb.WriteString(fmt.Sprintf("- Team Conventions: %s\n", analysisContext.TeamConventions))
	}

	// Learned patterns
	if len(analysisContext.LearnedPatterns) > 0 {
		sb.WriteString("\n## Learned Code Patterns\n")
		for _, pattern := range analysisContext.LearnedPatterns {
			sb.WriteString(fmt.Sprintf("- [%s] %s (confidence: %.2f)\n",
				pattern.Type, pattern.Description, pattern.Confidence))
		}
	}

	// Changed files - only analyze diff
	sb.WriteString("\n## Changed Files\n")
	for i, file := range files {
		sb.WriteString(fmt.Sprintf("\n### File %d: %s\n", i+1, file.Path))
		if file.Diff != "" {
			sb.WriteString("```diff\n")
			sb.WriteString(file.Diff)
			sb.WriteString("\n```\n")
		} else {
			// If no diff, don't send content
			sb.WriteString("(No changes detected)\n")
		}
	}

	// Analysis requirements
	sb.WriteString(`
## Please analyze the following aspects:

1. **Potential Issues**: Are there bugs, logic errors, or performance issues?
2. **Pattern Matching**: Does it match the project's historical coding patterns?
3. **Regression Risk**: Could it introduce regression issues?
4. **Improvement Suggestions**: Specific fixes or optimization recommendations

## Return Format

Please return in JSON format:
` + "```json\n" + `
{
  "issues": [
    {
      "file": "file path",
      "line": line number,
      "severity": "error|warning|info",
      "message": "issue description",
      "suggestion": "fix suggestion"
    }
  ],
  "summary": "overall analysis summary",
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

// buildLearnPrompt build learning prompt
func (c *Client) buildLearnPrompt(commit *models.CommitInfo, diff string, existingPatterns []models.CodePattern) string {
	var sb strings.Builder

	sb.WriteString("Analyze the following Git commit to identify code patterns and conventions:\n\n")
	sb.WriteString(fmt.Sprintf("## Commit Info\n"))
	sb.WriteString(fmt.Sprintf("- Hash: %s\n", commit.Hash))
	sb.WriteString(fmt.Sprintf("- Message: %s\n", commit.Message))
	sb.WriteString(fmt.Sprintf("- Author: %s\n", commit.Author))
	sb.WriteString(fmt.Sprintf("- Time: %s\n", commit.Timestamp.Format(time.RFC3339)))

	sb.WriteString("\n## Code Changes\n")
	sb.WriteString("```diff\n")
	sb.WriteString(diff)
	sb.WriteString("\n```\n")

	if len(existingPatterns) > 0 {
		sb.WriteString("\n## Known Code Patterns\n")
		for _, p := range existingPatterns {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", p.Type, p.Description))
		}
	}

	sb.WriteString(`
## Analysis Tasks

1. Identify new code patterns (naming, structure, error handling, etc.)
2. Determine if it's a bug fix, if so, record the error pattern
3. Extract team coding conventions

## Return Format

Return discovered patterns in JSON array:
` + "```json\n" + `
[
  {
    "type": "naming|structure|error_handling|concurrency|testing|comment",
    "description": "pattern description",
    "example": "code example",
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
