package git

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/openclaw-coding/skill-seed/internal/domain"
)

// Repository Git 仓储实现
type Repository struct {
	projectRoot string
}

// NewRepository 创建 Git 仓储
func NewRepository(projectRoot string) *Repository {
	return &Repository{
		projectRoot: projectRoot,
	}
}

// GetCommits 获取提交历史
func (r *Repository) GetCommits(ctx context.Context, limit int) ([]domain.CommitInfo, error) {
	cmd := exec.CommandContext(ctx, "git", "log",
		fmt.Sprintf("--max-count=%d", limit),
		"--pretty=format:%H|%an|%ad|%s",
		"--date=iso",
	)
	cmd.Dir = r.projectRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	commits := make([]domain.CommitInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		date, err := time.Parse("2006-01-02 15:04:05 -0700", parts[2])
		if err != nil {
			continue
		}

		commits = append(commits, domain.NewCommitInfo(
			parts[0],  // hash
			parts[1],  // author
			parts[3],  // message
			date,      // date
		))
	}

	return commits, nil
}

// GetDiff 获取指定提交的差异
func (r *Repository) GetDiff(ctx context.Context, hash string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "show", hash, "--format=")
	cmd.Dir = r.projectRoot

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git show failed: %w", err)
	}

	return string(output), nil
}

// GetStagedFiles 获取暂存文件
func (r *Repository) GetStagedFiles(ctx context.Context) ([]domain.FileInfo, error) {
	// 获取暂存文件列表
	cmd := exec.CommandContext(ctx, "git", "diff", "--cached", "--name-status")
	cmd.Dir = r.projectRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff --cached failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	files := make([]domain.FileInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		// 格式: M\tpath/to/file.go
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}

		status := r.parseStatus(parts[0])
		filePath := parts[1]

		// 获取文件内容
		content, err := r.getFileContent(ctx, filePath)
		if err != nil {
			continue
		}

		files = append(files, domain.FileInfo{
			Path:     filePath,
			Content:  content,
			Language: domain.NewFileInfo(filePath, "").Language,
			Status:   status,
		})
	}

	return files, nil
}

// GetAllFiles 获取所有文件
func (r *Repository) GetAllFiles(ctx context.Context) ([]domain.FileInfo, error) {
	cmd := exec.CommandContext(ctx, "git", "ls-files")
	cmd.Dir = r.projectRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	files := make([]domain.FileInfo, 0, len(lines))

	for _, filePath := range lines {
		if filePath == "" {
			continue
		}

		// 获取文件内容
		content, err := r.getFileContent(ctx, filePath)
		if err != nil {
			continue
		}

		files = append(files, domain.FileInfo{
			Path:     filePath,
			Content:  content,
			Language: domain.NewFileInfo(filePath, "").Language,
			Status:   domain.StatusModified,
		})
	}

	return files, nil
}

// GetCurrentBranch 获取当前分支
func (r *Repository) GetCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = r.projectRoot

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetProjectRoot 获取项目根目录
func (r *Repository) GetProjectRoot(ctx context.Context) (string, error) {
	return r.projectRoot, nil
}

// getFileContent 获取文件内容
func (r *Repository) getFileContent(ctx context.Context, filePath string) (string, error) {
	fullPath := filepath.Join(r.projectRoot, filePath)

	cmd := exec.CommandContext(ctx, "cat", fullPath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// parseStatus 解析文件状态
func (r *Repository) parseStatus(status string) domain.Status {
	switch strings.TrimSpace(status) {
	case "A":
		return domain.StatusAdded
	case "D":
		return domain.StatusDeleted
	default:
		return domain.StatusModified
	}
}
