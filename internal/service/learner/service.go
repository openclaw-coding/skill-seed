package learner

import (
	"context"

	"github.com/openclaw-coding/skill-seed/internal/agent"
	"github.com/openclaw-coding/skill-seed/internal/domain"
)

// Service 学习服务
type Service struct {
	agent       agent.Agent
	gitRepo     domain.GitRepository
	patternRepo domain.PatternRepository
}

// NewService 创建学习服务
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

// Learn 从 Git 历史学习模式
func (s *Service) Learn(ctx context.Context, limit int) error {
	// 1. 获取 Git 提交历史
	commits, err := s.gitRepo.GetCommits(ctx, limit)
	if err != nil {
		return err
	}

	// 2. 获取已知模式
	knownPatterns, err := s.patternRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// 转换为 agent.Pattern
	agentPatterns := make([]agent.Pattern, len(knownPatterns))
	for i, p := range knownPatterns {
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

	// 3. 遍历每个提交进行学习
	for _, c := range commits {
		diff, err := s.gitRepo.GetDiff(ctx, c.Hash)
		if err != nil {
			continue
		}

		req := &agent.LearnRequest{
			Commit: agent.CommitInfo{
				Hash:    c.Hash,
				Author:  c.Author,
				Date:    c.Date,
				Message: c.Message,
			},
			Diff:          diff,
			KnownPatterns: agentPatterns,
		}

		result, err := s.agent.LearnFromCommit(ctx, req)
		if err != nil {
			continue
		}

		// 4. 保存新模式
		for _, p := range result.Patterns {
			newPattern := domain.NewPattern(p.ID, p.Name, domain.Category(p.Category))
			newPattern.SetDescription(p.Description)
			newPattern.SetExamples(p.GoodExample, p.BadExample)
			newPattern.SetRule(p.Rule)
			newPattern.Confidence = p.Confidence
			newPattern.Frequency = p.Frequency

			if err := s.patternRepo.Save(ctx, newPattern); err != nil {
				continue
			}
		}
	}

	return nil
}

// LearnFromCommit 从单个提交学习
func (s *Service) LearnFromCommit(ctx context.Context, c domain.CommitInfo) error {
	diff, err := s.gitRepo.GetDiff(ctx, c.Hash)
	if err != nil {
		return err
	}

	knownPatterns, err := s.patternRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	agentPatterns := make([]agent.Pattern, len(knownPatterns))
	for i, p := range knownPatterns {
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

	req := &agent.LearnRequest{
		Commit: agent.CommitInfo{
			Hash:    c.Hash,
			Author:  c.Author,
			Date:    c.Date,
			Message: c.Message,
		},
		Diff:          diff,
		KnownPatterns: agentPatterns,
	}

	result, err := s.agent.LearnFromCommit(ctx, req)
	if err != nil {
		return err
	}

	// 保存新模式
	for _, p := range result.Patterns {
		newPattern := domain.NewPattern(p.ID, p.Name, domain.Category(p.Category))
		newPattern.SetDescription(p.Description)
		newPattern.SetExamples(p.GoodExample, p.BadExample)
		newPattern.SetRule(p.Rule)
		newPattern.Confidence = p.Confidence
		newPattern.Frequency = p.Frequency

		if err := s.patternRepo.Save(ctx, newPattern); err != nil {
			continue
		}
	}

	return nil
}
