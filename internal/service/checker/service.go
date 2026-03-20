package checker

import (
	"context"
	"os"

	"github.com/openclaw-coding/skill-seed/internal/agent"
	"github.com/openclaw-coding/skill-seed/internal/domain"
)

// Service 检查服务
type Service struct {
	agent       agent.Agent
	gitRepo     domain.GitRepository
	patternRepo domain.PatternRepository
}

// NewService 创建检查服务
func NewService(
	ag agent.Agent,
	gitRepo domain.GitRepository,
	patternRepo domain.PatternRepository,
) *Service {
	return &Service{
		agent:       ag,
		gitRepo:     gitRepo,
		patternRepo: patternRepo,
	}
}

// Check 检查暂存文件
func (s *Service) Check(ctx context.Context) ([]domain.Issue, error) {
	// 1. 获取暂存文件
	files, err := s.gitRepo.GetStagedFiles(ctx)
	if err != nil {
		return nil, err
	}

	return s.CheckFiles(ctx, files)
}

// CheckAll 检查所有文件
func (s *Service) CheckAll(ctx context.Context) ([]domain.Issue, error) {
	// 1. 获取所有文件
	files, err := s.gitRepo.GetAllFiles(ctx)
	if err != nil {
		return nil, err
	}

	return s.CheckFiles(ctx, files)
}

// CheckFiles 检查指定文件
func (s *Service) CheckFiles(ctx context.Context, files []domain.FileInfo) ([]domain.Issue, error) {
	// 2. 获取项目上下文
	context := s.getProjectContext()

	// 3. 获取已知模式
	patterns, err := s.patternRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 agent.Pattern
	agentPatterns := make([]agent.Pattern, len(patterns))
	for i, p := range patterns {
		agentPatterns[i] = agent.Pattern{
			ID:          p.ID,
			Name:        p.Name,
			Category:    string(p.Category),
			Description: p.Description,
			GoodExample: p.GoodExample,
			BadExample:  p.BadExample,
			Rule:        p.Rule,
			Confidence:  p.Confidence,
			Frequency:   p.Frequency,
			Source:      string(p.Source),
			CreatedAt:   p.CreatedAt,
		}
	}

	// 4. 获取最近提交
	recentCommits, err := s.gitRepo.GetCommits(ctx, 10)
	if err != nil {
		return nil, err
	}

	// 转换为 agent.CommitInfo
	agentCommits := make([]agent.CommitInfo, len(recentCommits))
	for i, c := range recentCommits {
		agentCommits[i] = agent.CommitInfo{
			Hash:    c.Hash,
			Author:  c.Author,
			Date:    c.Date,
			Message: c.Message,
		}
	}

	// 5. 转换文件为 agent.FileInfo
	agentFiles := make([]agent.FileInfo, len(files))
	for i, f := range files {
		agentFiles[i] = agent.FileInfo{
			Path:     f.Path,
			Content:  f.Content,
			Language: f.Language,
			Status:   string(f.Status),
		}
	}

	// 6. 调用 Agent 分析
	req := &agent.AnalyzeRequest{
		Files:         agentFiles,
		Context:       context,
		Patterns:      agentPatterns,
		RecentCommits: agentCommits,
	}

	result, err := s.agent.AnalyzeCode(ctx, req)
	if err != nil {
		return nil, err
	}

	// 7. 转换结果
	issues := make([]domain.Issue, len(result.Issues))
	for i, iss := range result.Issues {
		issues[i] = domain.Issue{
			File:       iss.File,
			Line:       iss.Line,
			Severity:   domain.Severity(iss.Severity),
			Message:    iss.Message,
			Suggestion: iss.Suggestion,
			PatternID:  iss.PatternID,
		}
	}

	return issues, nil
}

// getProjectContext 获取项目上下文
func (s *Service) getProjectContext() agent.ProjectContext {
	// TODO: 从配置或其他来源获取项目上下文
	return agent.ProjectContext{
		Name:         "project",
		Language:     "go",
		Frameworks:   []string{},
		Dependencies: []string{},
	}
}

// GetPatterns 获取检查使用的模式
func (s *Service) GetPatterns(ctx context.Context) ([]domain.Pattern, error) {
	return s.patternRepo.GetAll(ctx)
}

// GetHighConfidencePatterns 获取高置信度模式
func (s *Service) GetHighConfidencePatterns(ctx context.Context, threshold float64) ([]domain.Pattern, error) {
	return s.patternRepo.GetHighConfidence(ctx, threshold)
}

// AnalyzeFiles 分析指定文件（用于 analyze 命令）
func (s *Service) AnalyzeFiles(ctx context.Context, absPaths []string) error {
	// 读取文件内容并转换为 FileInfo
	files := make([]domain.FileInfo, 0, len(absPaths))
	for _, path := range absPaths {
		content, err := s.readFileContent(path)
		if err != nil {
		continue // 跳过无法读取的文件
		}
		files = append(files, domain.FileInfo{
			Path:     path,
			Content:  content,
			Language: domain.NewFileInfo(path, "").Language,
			Status:   domain.StatusModified,
		})
	}

	// 使用 CheckFiles 进行检查
	_, err := s.CheckFiles(ctx, files)
	return err
}

// readFileContent 读取文件内容
func (s *Service) readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
