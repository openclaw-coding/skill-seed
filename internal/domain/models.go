package domain

import "time"

// ==================== Pattern ====================

// Category 模式分类
type Category string

const (
	CategoryNaming      Category = "naming"
	CategoryError       Category = "error"
	CategoryStructure   Category = "structure"
	CategoryConcurrency Category = "concurrency"
	CategoryTesting     Category = "testing"
)

// Source 模式来源
type Source string

const (
	SourceLearned Source = "learned"
	SourceDefault Source = "default"
)

// Pattern 代码模式聚合根
type Pattern struct {
	ID          string
	Name        string
	Category    Category
	Description string
	GoodExample string
	BadExample  string
	Rule        string
	Confidence  float64
	Frequency   int
	Source      Source
	CreatedAt   time.Time
}

// NewPattern 创建新的模式
func NewPattern(id, name string, category Category) *Pattern {
	return &Pattern{
		ID:         id,
		Name:       name,
		Category:   category,
		Confidence: 0.0,
		Frequency:  0,
		Source:     SourceLearned,
		CreatedAt:  time.Now(),
	}
}

// IsValid 验证模式是否有效
func (p *Pattern) IsValid() bool {
	return p.ID != "" &&
		p.Name != "" &&
		p.Category != "" &&
		p.Confidence >= 0.0 &&
		p.Confidence <= 1.0
}

// UpdateConfidence 更新置信度（基于频率加权平均）
func (p *Pattern) UpdateConfidence(newConfidence float64) {
	p.Confidence = (p.Confidence*float64(p.Frequency) + newConfidence) / float64(p.Frequency+1)
	p.Frequency++
}

// SetExamples 设置示例
func (p *Pattern) SetExamples(good, bad string) {
	p.GoodExample = good
	p.BadExample = bad
}

// SetDescription 设置描述
func (p *Pattern) SetDescription(desc string) {
	p.Description = desc
}

// SetRule 设置规则
func (p *Pattern) SetRule(rule string) {
	p.Rule = rule
}

// ==================== Issue ====================

// Severity 问题严重程度
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Issue 问题实体
type Issue struct {
	File       string   // 文件路径
	Line       int      // 行号
	Severity   Severity // 严重程度
	Message    string   // 问题描述
	Suggestion string   // 修复建议
	PatternID  string   // 关联的模式ID
}

// NewIssue 创建新问题
func NewIssue(file string, line int, severity Severity, message string) *Issue {
	return &Issue{
		File:     file,
		Line:     line,
		Severity: severity,
		Message:  message,
	}
}

// IsError 是否是错误级别
func (i *Issue) IsError() bool {
	return i.Severity == SeverityError
}

// IsWarning 是否是警告级别
func (i *Issue) IsWarning() bool {
	return i.Severity == SeverityWarning
}

// SetSuggestion 设置修复建议
func (i *Issue) SetSuggestion(suggestion string) {
	i.Suggestion = suggestion
}

// SetPatternID 设置关联的模式ID
func (i *Issue) SetPatternID(patternID string) {
	i.PatternID = patternID
}

// ==================== Rule ====================

// Rule 规则实体
type Rule struct {
	ID          string   // 规则ID
	Name        string   // 规则名称
	Category    string   // 分类
	Description string   // 描述
	PatternIDs  []string // 关联的模式ID
	Enabled     bool     // 是否启用
}

// NewRule 创建新规则
func NewRule(id, name, category string) *Rule {
	return &Rule{
		ID:         id,
		Name:       name,
		Category:   category,
		PatternIDs: []string{},
		Enabled:    true,
	}
}

// AddPattern 添加关联的模式
func (r *Rule) AddPattern(patternID string) {
	for _, id := range r.PatternIDs {
		if id == patternID {
			return
		}
	}
	r.PatternIDs = append(r.PatternIDs, patternID)
}

// RemovePattern 移除关联的模式
func (r *Rule) RemovePattern(patternID string) {
	for i, id := range r.PatternIDs {
		if id == patternID {
			r.PatternIDs = append(r.PatternIDs[:i], r.PatternIDs[i+1:]...)
			return
		}
	}
}

// Enable 启用规则
func (r *Rule) Enable() {
	r.Enabled = true
}

// Disable 禁用规则
func (r *Rule) Disable() {
	r.Enabled = false
}

// SetDescription 设置描述
func (r *Rule) SetDescription(desc string) {
	r.Description = desc
}

// ==================== Commit ====================

// CommitInfo 提交值对象
type CommitInfo struct {
	Hash    string    // 提交哈希
	Author  string    // 作者
	Date    time.Time // 提交时间
	Message string    // 提交消息
}

// NewCommitInfo 创建提交信息
func NewCommitInfo(hash, author, message string, date time.Time) CommitInfo {
	return CommitInfo{
		Hash:    hash,
		Author:  author,
		Date:    date,
		Message: message,
	}
}

// IsEmpty 是否为空
func (c CommitInfo) IsEmpty() bool {
	return c.Hash == ""
}

// ShortHash 获取短哈希（前7位）
func (c CommitInfo) ShortHash() string {
	if len(c.Hash) <= 7 {
		return c.Hash
	}
	return c.Hash[:7]
}

// Summary 获取提交摘要（第一行消息）
func (c CommitInfo) Summary() string {
	for i, ch := range c.Message {
		if ch == '\n' {
			return c.Message[:i]
		}
	}
	return c.Message
}

// ==================== File ====================

// Status 文件状态
type Status string

const (
	StatusAdded    Status = "added"
	StatusModified Status = "modified"
	StatusDeleted  Status = "deleted"
)

// FileInfo 文件值对象
type FileInfo struct {
	Path     string // 文件路径
	Content  string // 文件内容
	Language string // 语言类型
	Status   Status // 状态
}

// NewFileInfo 创建文件信息
func NewFileInfo(path, content string) FileInfo {
	return FileInfo{
		Path:     path,
		Content:  content,
		Language: detectLanguage(path),
		Status:   StatusModified,
	}
}

// IsGoFile 是否是 Go 文件
func (f FileInfo) IsGoFile() bool {
	return f.Language == "go"
}

// IsTestFile 是否是测试文件
func (f FileInfo) IsTestFile() bool {
	return len(f.Path) > 8 && f.Path[len(f.Path)-8:] == "_test.go"
}

// IsEmpty 是否为空
func (f FileInfo) IsEmpty() bool {
	return f.Content == ""
}

// LineCount 获取行数
func (f FileInfo) LineCount() int {
	count := 0
	for _, ch := range f.Content {
		if ch == '\n' {
			count++
		}
	}
	return count + 1
}

// detectLanguage 根据文件扩展名检测语言
func detectLanguage(path string) string {
	if len(path) == 0 {
		return ""
	}

	// 获取扩展名
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext := path[i+1:]
			switch ext {
			case "go":
				return "go"
			case "js", "jsx":
				return "javascript"
			case "ts", "tsx":
				return "typescript"
			case "py":
				return "python"
			case "java":
				return "java"
			case "cpp", "cc", "cxx":
				return "cpp"
			case "c":
				return "c"
			case "rs":
				return "rust"
			case "rb":
				return "ruby"
			case "php":
				return "php"
			case "swift":
				return "swift"
			case "kt":
				return "kotlin"
			case "scala":
				return "scala"
			case "md":
				return "markdown"
			case "yaml", "yml":
				return "yaml"
			case "json":
				return "json"
			case "xml":
				return "xml"
			case "sql":
				return "sql"
			case "sh":
				return "shell"
			case "dockerfile":
				return "dockerfile"
			case "makefile":
				return "makefile"
			default:
				return ext
			}
		}
		if path[i] == '/' {
			break
		}
	}

	// 检查特殊文件名
	switch path {
	case "Dockerfile":
		return "dockerfile"
	case "Makefile":
		return "makefile"
	case "go.mod", "go.sum":
		return "go"
	}

	return ""
}
