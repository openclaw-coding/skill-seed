package boltdb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/openclaw-coding/skill-seed/internal/domain"
	bolt "go.etcd.io/bbolt"
)

// PatternRepository Pattern 仓储实现
type PatternRepository struct {
	db *bolt.DB
}

var bucketPatterns = []byte("patterns")

// NewPatternRepository 创建 Pattern 仓储
func NewPatternRepository(dbPath string) (*PatternRepository, error) {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// 创建 bucket
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketPatterns); err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketPatterns, err)
		}
		return nil
	})

	if err != nil {
		db.Close()
		return nil, err
	}

	return &PatternRepository{db: db}, nil
}

// Get 根据ID获取模式
func (r *PatternRepository) Get(ctx context.Context, id string) (*domain.Pattern, error) {
	var p domain.Pattern

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("pattern not found: %s", id)
		}
		return json.Unmarshal(data, &p)
	})

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetAll 获取所有模式
func (r *PatternRepository) GetAll(ctx context.Context) ([]domain.Pattern, error) {
	var patterns []domain.Pattern

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		return b.ForEach(func(k, v []byte) error {
			var p domain.Pattern
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			patterns = append(patterns, p)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return patterns, nil
}

// GetByCategory 根据分类获取模式
func (r *PatternRepository) GetByCategory(ctx context.Context, category domain.Category) ([]domain.Pattern, error) {
	all, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []domain.Pattern
	for _, p := range all {
		if p.Category == category {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

// GetHighConfidence 获取高置信度模式
func (r *PatternRepository) GetHighConfidence(ctx context.Context, threshold float64) ([]domain.Pattern, error) {
	all, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []domain.Pattern
	for _, p := range all {
		if p.Confidence >= threshold {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

// Save 保存模式
func (r *PatternRepository) Save(ctx context.Context, p *domain.Pattern) error {
	if !p.IsValid() {
		return fmt.Errorf("invalid pattern")
	}

	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(p.ID), data)
	})
}

// Delete 删除模式
func (r *PatternRepository) Delete(ctx context.Context, id string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		return b.Delete([]byte(id))
	})
}

// Count 统计模式数量
func (r *PatternRepository) Count(ctx context.Context) (int, error) {
	count := 0

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		return b.ForEach(func(k, v []byte) error {
			count++
			return nil
		})
	})

	return count, err
}

// Close 关闭数据库
func (r *PatternRepository) Close() error {
	return r.db.Close()
}

// GetDB 获取底层数据库实例
func (r *PatternRepository) GetDB() *bolt.DB {
	return r.db
}

// RuleRepository Rule 仓储实现
type RuleRepository struct {
	db *bolt.DB
}

var bucketRules = []byte("rules")

// NewRuleRepository 创建 Rule 仓储
func NewRuleRepository(db *bolt.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

// Get 根据ID获取规则
func (r *RuleRepository) Get(ctx context.Context, id string) (*domain.Rule, error) {
	var rl domain.Rule

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRules)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("rule not found: %s", id)
		}
		return json.Unmarshal(data, &rl)
	})

	if err != nil {
		return nil, err
	}

	return &rl, nil
}

// GetAll 获取所有规则
func (r *RuleRepository) GetAll(ctx context.Context) ([]domain.Rule, error) {
	var rules []domain.Rule

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRules)
		return b.ForEach(func(k, v []byte) error {
			var rl domain.Rule
			if err := json.Unmarshal(v, &rl); err != nil {
				return err
			}
			rules = append(rules, rl)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return rules, nil
}

// GetByCategory 根据分类获取规则
func (r *RuleRepository) GetByCategory(ctx context.Context, category string) ([]domain.Rule, error) {
	all, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []domain.Rule
	for _, rl := range all {
		if rl.Category == category {
			filtered = append(filtered, rl)
		}
	}

	return filtered, nil
}

// Save 保存规则
func (r *RuleRepository) Save(ctx context.Context, rl *domain.Rule) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRules)
		data, err := json.Marshal(rl)
		if err != nil {
			return err
		}
		return b.Put([]byte(rl.ID), data)
	})
}

// Delete 删除规则
func (r *RuleRepository) Delete(ctx context.Context, id string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRules)
		return b.Delete([]byte(id))
	})
}

// initBuckets 初始化所有 buckets
func initBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{bucketPatterns, bucketRules}
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})
}
