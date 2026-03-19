package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/openclaw-coding/grow-check/pkg/models"
)

// GitOperator Git 操作器
type GitOperator struct {
	repoPath string
}

// NewGitOperator 创建 Git 操作器
func NewGitOperator(repoPath string) *GitOperator {
	return &GitOperator{repoPath: repoPath}
}

// GetStagedFiles get staged files
func (g *GitOperator) GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	// Parse output: M filename, MM filename, A filename, etc.
	// Git status output format: XY filename
	// X = staged status, Y = work tree status
	// We only care about staged files (X is not space, ?, !)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	files := make([]string, 0)

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		status := line[0]
		// Check if file is staged: status is M, M, A, R, C, etc. (not space, ?, !)
		if status != ' ' && status != '?' && status != '!' {
			// Extract filename (starts at position 3)
			filename := strings.TrimSpace(line[3:])
			files = append(files, filename)
		}
	}

	return files, nil
}

// GetFileContent 获取文件内容
func (g *GitOperator) GetFileContent(path string) (string, error) {
	fullPath := filepath.Join(g.repoPath, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return string(data), nil
}

// GetStagedFileDiff 获取暂存文件的 diff
func (g *GitOperator) GetStagedFileDiff(path string) (string, error) {
	cmd := exec.Command("git", "diff", "--cached", path)
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff for %s: %w", path, err)
	}
	return string(output), nil
}

// GetRecentCommits 获取最近的提交
func (g *GitOperator) GetRecentCommits(limit int, since time.Time) ([]models.CommitInfo, error) {
	args := []string{
		"log",
		fmt.Sprintf("--max-count=%d", limit),
		"--pretty=format:%H|%s|%an|%aI",
	}
	
	if !since.IsZero() {
		args = append(args, fmt.Sprintf("--since=%s", since.Format(time.RFC3339)))
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]models.CommitInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, parts[3])
		if err != nil {
			timestamp = time.Now()
		}

		commits = append(commits, models.CommitInfo{
			Hash:      parts[0],
			Message:   parts[1],
			Author:    parts[2],
			Timestamp: timestamp,
		})
	}

	return commits, nil
}

// GetCommitFiles 获取提交修改的文件
func (g *GitOperator) GetCommitFiles(hash string) ([]string, error) {
	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", hash)
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}

// GetCommitDiff 获取提交的 diff
func (g *GitOperator) GetCommitDiff(hash string) (string, error) {
	cmd := exec.Command("git", "show", hash)
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit diff: %w", err)
	}
	return string(output), nil
}

// GetRemoteURL 获取远程仓库 URL
func (g *GitOperator) GetRemoteURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote url: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// IsGitRepo 检查是否是 Git 仓库
func (g *GitOperator) IsGitRepo() bool {
	gitDir := filepath.Join(g.repoPath, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// InstallPreCommitHook installs pre-commit hook
func (g *GitOperator) InstallPreCommitHook(skillPath string) error {
	hookPath := filepath.Join(g.repoPath, ".git", "hooks", "pre-commit")

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		// Read existing content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return fmt.Errorf("failed to read existing hook: %w", err)
		}

		// If it already contains our hook, skip
		if strings.Contains(string(content), "grow-check") {
			return nil
		}

		// Create chain call
		newContent := fmt.Sprintf(`#!/bin/sh
# Chain loading for multiple hooks (preserving existing hooks)

# Original hook content:
%s

# grow-check hook
grow-check check || exit $?
`, string(content))

		if err := os.WriteFile(hookPath, []byte(newContent), 0755); err != nil {
			return fmt.Errorf("failed to update hook: %w", err)
		}
		return nil
	}

	// Create new hook script
	content := `#!/bin/sh
# grow-check pre-commit hook
grow-check check || exit $?
`

	if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("failed to create hook: %w", err)
	}

	return nil
}

// GetProjectName 获取项目名称
func (g *GitOperator) GetProjectName() string {
	remote, err := g.GetRemoteURL()
	if err != nil || remote == "" {
		// 如果没有远程仓库，使用目录名
		return filepath.Base(g.repoPath)
	}

	// 从 URL 提取项目名
	// 支持: https://github.com/user/repo.git 或 git@github.com:user/repo.git
	parts := strings.Split(remote, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		return strings.TrimSuffix(name, ".git")
	}

	return filepath.Base(g.repoPath)
}
