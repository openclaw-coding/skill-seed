package models

import (
	"time"
)

// PatternType 定义模式类型
type PatternType string

const (
	PatternNaming       PatternType = "naming"        // 命名规范
	PatternStructure    PatternType = "structure"     // 代码结构
	PatternErrorHandling PatternType = "error_handling" // 错误处理
	PatternConcurrency  PatternType = "concurrency"   // 并发模式
	PatternTesting      PatternType = "testing"       // 测试模式
	PatternComment      PatternType = "comment"       // 注释规范
)

// CodePattern 代码模式
type CodePattern struct {
	ID          string      `json:"id"`
	Type        PatternType `json:"type"`
	Description string      `json:"description"`
	Examples    []string    `json:"examples"`     // 从历史提交提取的示例
	Frequency   int         `json:"frequency"`    // 出现频率
	Confidence  float64     `json:"confidence"`   // 置信度 (0-1)
	AutoFixable bool        `json:"auto_fixable"` // 是否可自动修复
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Issue 检查发现的问题
type Issue struct {
	File       string `json:"file"`       // 文件路径
	Line       int    `json:"line"`       // 行号
	Column     int    `json:"column"`     // 列号
	Severity   string `json:"severity"`   // error/warning/info
	Message    string `json:"message"`    // 问题描述
	Suggestion string `json:"suggestion"` // 修复建议
	PatternID  string `json:"pattern_id"` // 关联的模式 ID
}

// CommitInfo 提交信息
type CommitInfo struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Files     []string  `json:"files"` // 修改的文件
}

// FileChange 文件变更
type FileChange struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Diff    string `json:"diff"`
}

// AnalysisContext Claude 分析的上下文
type AnalysisContext struct {
	ProjectType       string        `json:"project_type"`        // 项目类型
	HistoricalBugs    []string      `json:"historical_bugs"`     // 历史 bug 模式
	TeamConventions   string        `json:"team_conventions"`    // 团队规范
	LearnedPatterns   []CodePattern `json:"learned_patterns"`    // 学习到的模式
	RecentCommits     []CommitInfo  `json:"recent_commits"`      // 最近提交
}

// AnalysisResult Claude 分析结果
type AnalysisResult struct {
	Issues     []Issue `json:"issues"`
	Summary    string  `json:"summary"`
	Confidence float64 `json:"confidence"`
}

// Rule 检查规则
type Rule struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        PatternType `json:"type"`
	Condition   string      `json:"condition"`   // 规则条件（简化版）
	Severity    string      `json:"severity"`    // error/warning/info
	AutoFix     bool        `json:"auto_fix"`    // 是否自动修复
	FixTemplate string      `json:"fix_template"` // 修复模板
	Confidence  float64     `json:"confidence"`
	Source      string      `json:"source"`      // "builtin" / "learned"
}
