package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/openclaw-coding/skill-seed/internal/agent"
	"github.com/openclaw-coding/skill-seed/internal/templates/prompts"
)

// ClaudeAgent Claude Agent 实现
type ClaudeAgent struct {
	commandPath  string
	timeout      time.Duration
	fallback     bool
	promptLoader *prompts.Loader
}

// New 创建 Claude Agent
func New(commandPath string, timeout time.Duration, loader *prompts.Loader) *ClaudeAgent {
	if commandPath == "" {
		commandPath = "claude"
	}
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	return &ClaudeAgent{
		commandPath:  commandPath,
		timeout:      timeout,
		fallback:     true,
		promptLoader: loader,
	}
}

// Name 返回 Agent 名称
func (c *ClaudeAgent) Name() string {
	return "claude"
}

// IsAvailable 检查 Agent 是否可用
func (c *ClaudeAgent) IsAvailable() bool {
	_, err := exec.LookPath(c.commandPath)
	return err == nil
}

// AnalyzeCode 分析代码
func (c *ClaudeAgent) AnalyzeCode(ctx context.Context, req *agent.AnalyzeRequest) (*agent.AnalyzeResult, error) {
	// 1. 构建提示词（从模板加载）
	prompt := c.promptLoader.Render("analyze", req)
	if prompt == "" {
		return nil, fmt.Errorf("failed to render analyze prompt")
	}

	// 2. 调用 Claude CLI
	output, err := c.callClaude(ctx, prompt)
	if err != nil {
		if c.fallback {
			return &agent.AnalyzeResult{
				Issues:     []agent.Issue{},
				Confidence: 0.0,
			}, nil
		}
		return nil, fmt.Errorf("claude analyze failed: %w", err)
	}

	// 3. 解析 JSON 结果
	result, err := c.parseAnalyzeResult(output)
	if err != nil {
		return nil, fmt.Errorf("parse result failed: %w", err)
	}

	result.AnalyzedAt = time.Now()
	return result, nil
}

// LearnFromCommit 从提交中学习
func (c *ClaudeAgent) LearnFromCommit(ctx context.Context, req *agent.LearnRequest) (*agent.LearnResult, error) {
	// 1. 构建提示词（从模板加载）
	prompt := c.promptLoader.Render("learn", req)
	if prompt == "" {
		return nil, fmt.Errorf("failed to render learn prompt")
	}

	// 2. 调用 Claude CLI
	output, err := c.callClaude(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("claude learn failed: %w", err)
	}

	// 3. 解析 JSON 结果
	result, err := c.parseLearnResult(output)
	if err != nil {
		return nil, fmt.Errorf("parse result failed: %w", err)
	}

	result.LearnedAt = time.Now()
	return result, nil
}

// callClaude 调用 Claude CLI
func (c *ClaudeAgent) callClaude(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.commandPath, "--print", prompt)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// parseAnalyzeResult 解析分析结果
func (c *ClaudeAgent) parseAnalyzeResult(output string) (*agent.AnalyzeResult, error) {
	// 尝试提取 JSON 部分
	jsonStart := findJSONStart(output, '{')
	jsonEnd := findJSONEnd(output, '}')

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return &agent.AnalyzeResult{
			Issues:     []agent.Issue{},
			Confidence: 0.0,
		}, nil
	}

	jsonStr := output[jsonStart : jsonEnd+1]

	var result struct {
		Issues []struct {
			File       string  `json:"file"`
			Line       int     `json:"line"`
			Severity   string  `json:"severity"`
			Message    string  `json:"message"`
			Suggestion string  `json:"suggestion"`
			PatternID  string  `json:"pattern_id"`
		} `json:"issues"`
		Suggestions []string `json:"suggestions"`
		Confidence  float64  `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	// 转换为 agent.Issue
	issues := make([]agent.Issue, len(result.Issues))
	for i, issue := range result.Issues {
		issues[i] = agent.Issue{
			File:       issue.File,
			Line:       issue.Line,
			Severity:   issue.Severity,
			Message:    issue.Message,
			Suggestion: issue.Suggestion,
			PatternID:  issue.PatternID,
		}
	}

	return &agent.AnalyzeResult{
		Issues:      issues,
		Suggestions: result.Suggestions,
		Confidence:  result.Confidence,
	}, nil
}

// parseLearnResult 解析学习结果
func (c *ClaudeAgent) parseLearnResult(output string) (*agent.LearnResult, error) {
	// 尝试提取 JSON 部分
	jsonStart := findJSONStart(output, '{')
	jsonEnd := findJSONEnd(output, '}')

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return &agent.LearnResult{
			Patterns: []agent.Pattern{},
		}, nil
	}

	jsonStr := output[jsonStart : jsonEnd+1]

	var result struct {
		Patterns []struct {
			ID          string  `json:"id"`
			Name        string  `json:"name"`
			Category    string  `json:"category"`
			Description string  `json:"description"`
			GoodExample string  `json:"good_example"`
			BadExample  string  `json:"bad_example"`
			Rule        string  `json:"rule"`
			Confidence  float64 `json:"confidence"`
			Frequency   int     `json:"frequency"`
		} `json:"patterns"`
		UpdatedRules []struct {
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			PatternIDs []string `json:"pattern_ids"`
			Enabled    bool     `json:"enabled"`
		} `json:"updated_rules"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	// 转换为 agent.Pattern
	patterns := make([]agent.Pattern, len(result.Patterns))
	for i, p := range result.Patterns {
		patterns[i] = agent.Pattern{
			ID:          p.ID,
			Name:        p.Name,
			Category:    p.Category,
			Description: p.Description,
			GoodExample: p.GoodExample,
			BadExample:  p.BadExample,
			Rule:        p.Rule,
			Confidence:  p.Confidence,
			Frequency:   p.Frequency,
			Source:      "learned",
			CreatedAt:   time.Now(),
		}
	}

	// 转换为 agent.Rule
	rules := make([]agent.Rule, len(result.UpdatedRules))
	for i, r := range result.UpdatedRules {
		rules[i] = agent.Rule{
			ID:         r.ID,
			Name:       r.Name,
			PatternIDs: r.PatternIDs,
			Enabled:    r.Enabled,
		}
	}

	return &agent.LearnResult{
		Patterns:     patterns,
		UpdatedRules: rules,
	}, nil
}

// findJSONStart 查找 JSON 开始位置
func findJSONStart(s string, startChar byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == startChar {
			return i
		}
	}
	return -1
}

// findJSONEnd 查找 JSON 结束位置
func findJSONEnd(s string, endChar byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == endChar {
			return i
		}
	}
	return -1
}
