package domain

import "context"

// PatternRepository 模式仓储接口
type PatternRepository interface {
	// Get 根据ID获取模式
	Get(ctx context.Context, id string) (*Pattern, error)

	// GetAll 获取所有模式
	GetAll(ctx context.Context) ([]Pattern, error)

	// GetByCategory 根据分类获取模式
	GetByCategory(ctx context.Context, category Category) ([]Pattern, error)

	// GetHighConfidence 获取高置信度模式
	GetHighConfidence(ctx context.Context, threshold float64) ([]Pattern, error)

	// Save 保存模式
	Save(ctx context.Context, p *Pattern) error

	// Delete 删除模式
	Delete(ctx context.Context, id string) error

	// Count 统计模式数量
	Count(ctx context.Context) (int, error)
}

// RuleRepository 规则仓储接口
type RuleRepository interface {
	// Get 根据ID获取规则
	Get(ctx context.Context, id string) (*Rule, error)

	// GetAll 获取所有规则
	GetAll(ctx context.Context) ([]Rule, error)

	// GetByCategory 根据分类获取规则
	GetByCategory(ctx context.Context, category string) ([]Rule, error)

	// Save 保存规则
	Save(ctx context.Context, r *Rule) error

	// Delete 删除规则
	Delete(ctx context.Context, id string) error
}

// GitRepository Git 操作接口
type GitRepository interface {
	// GetCommits 获取提交历史
	GetCommits(ctx context.Context, limit int) ([]CommitInfo, error)

	// GetDiff 获取指定提交的差异
	GetDiff(ctx context.Context, hash string) (string, error)

	// GetStagedFiles 获取暂存文件
	GetStagedFiles(ctx context.Context) ([]FileInfo, error)

	// GetAllFiles 获取所有文件
	GetAllFiles(ctx context.Context) ([]FileInfo, error)

	// GetCurrentBranch 获取当前分支
	GetCurrentBranch(ctx context.Context) (string, error)

	// GetProjectRoot 获取项目根目录
	GetProjectRoot(ctx context.Context) (string, error)
}
