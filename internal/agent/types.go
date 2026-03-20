package agent

import (
	"context"
	"time"
)

// Agent AI Agent 接口，支持多种后端
type Agent interface {
	// Name 返回 Agent 名称
	Name() string

	// IsAvailable 检查 Agent 是否可用
	IsAvailable() bool

	// AnalyzeCode 分析代码，返回问题和建议
	AnalyzeCode(ctx context.Context, req *AnalyzeRequest) (*AnalyzeResult, error)

	// LearnFromCommit 从提交中学习模式
	LearnFromCommit(ctx context.Context, req *LearnRequest) (*LearnResult, error)
}

// AnalyzeRequest 分析请求
type AnalyzeRequest struct {
	Files         []FileInfo      // 待分析文件
	Context       ProjectContext  // 项目上下文
	Patterns      []Pattern       // 已知模式
	RecentCommits []CommitInfo    // 最近提交
}

// AnalyzeResult 分析结果
type AnalyzeResult struct {
	Issues      []Issue   // 发现的问题
	Suggestions []string  // 改进建议
	Confidence  float64   // 置信度
	AnalyzedAt  time.Time // 分析时间
}

// LearnRequest 学习请求
type LearnRequest struct {
	Commit        CommitInfo // 提交信息
	Diff          string     // 代码变更
	KnownPatterns []Pattern  // 已知模式
}

// LearnResult 学习结果
type LearnResult struct {
	Patterns     []Pattern // 新学习的模式
	UpdatedRules []Rule    // 更新的规则
	LearnedAt    time.Time // 学习时间
}

// FileInfo 文件信息
type FileInfo struct {
	Path     string // 文件路径
	Content  string // 文件内容
	Language string // 语言类型
	Status   string // 状态 (added/modified/deleted)
}

// ProjectContext 项目上下文
type ProjectContext struct {
	Name         string   // 项目名称
	Language     string   // 主要语言
	Frameworks   []string // 使用的框架
	Dependencies []string // 依赖项
}

// CommitInfo 提交信息
type CommitInfo struct {
	Hash    string    // 提交哈希
	Author  string    // 作者
	Date    time.Time // 提交时间
	Message string    // 提交消息
}

// Pattern 代码模式
type Pattern struct {
	ID          string    // 模式ID
	Name        string    // 模式名称
	Category    string    // 分类 (naming/error/structure/concurrency/testing)
	Description string    // 描述
	GoodExample string    // 好的示例
	BadExample  string    // 坏的示例
	Rule        string    // 规则说明
	Confidence  float64   // 置信度 (0.0-1.0)
	Frequency   int       // 出现频率
	Source      string    // 来源 (learned/default)
	CreatedAt   time.Time // 创建时间
}

// Issue 问题
type Issue struct {
	File       string // 文件路径
	Line       int    // 行号
	Severity   string // 严重程度 (error/warning/info)
	Message    string // 问题描述
	Suggestion string // 修复建议
	PatternID  string // 关联的模式ID
}

// Rule 规则
type Rule struct {
	ID          string   // 规则ID
	Name        string   // 规则名称
	Category    string   // 分类
	Description string   // 描述
	PatternIDs  []string // 关联的模式ID
	Enabled     bool     // 是否启用
}
