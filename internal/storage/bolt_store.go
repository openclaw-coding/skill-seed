package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/openclaw-coding/grow-check/pkg/models"
	bolt "go.etcd.io/bbolt"
)

// Store 存储接口
type Store struct {
	db *bolt.DB
}

var (
	bucketPatterns    = []byte("patterns")
	bucketRules       = []byte("rules")
	bucketMetadata    = []byte("metadata")
	bucketLearnHistory = []byte("learn_history")
)

// New 创建存储实例
func New(dbPath string) (*Store, error) {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// 创建 buckets
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range [][]byte{bucketPatterns, bucketRules, bucketMetadata, bucketLearnHistory} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})

	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

// Close 关闭数据库
func (s *Store) Close() error {
	return s.db.Close()
}

// SavePattern 保存代码模式
func (s *Store) SavePattern(pattern *models.CodePattern) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		data, err := json.Marshal(pattern)
		if err != nil {
			return err
		}
		return b.Put([]byte(pattern.ID), data)
	})
}

// GetPattern 获取代码模式
func (s *Store) GetPattern(id string) (*models.CodePattern, error) {
	var pattern models.CodePattern
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("pattern not found: %s", id)
		}
		return json.Unmarshal(data, &pattern)
	})
	if err != nil {
		return nil, err
	}
	return &pattern, nil
}

// GetAllPatterns 获取所有模式
func (s *Store) GetAllPatterns() ([]models.CodePattern, error) {
	var patterns []models.CodePattern
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPatterns)
		return b.ForEach(func(k, v []byte) error {
			var pattern models.CodePattern
			if err := json.Unmarshal(v, &pattern); err != nil {
				return err
			}
			patterns = append(patterns, pattern)
			return nil
		})
	})
	return patterns, err
}

// SaveRule 保存规则
func (s *Store) SaveRule(rule *models.Rule) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRules)
		data, err := json.Marshal(rule)
		if err != nil {
			return err
		}
		return b.Put([]byte(rule.ID), data)
	})
}

// GetAllRules 获取所有规则
func (s *Store) GetAllRules() ([]models.Rule, error) {
	var rules []models.Rule
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRules)
		return b.ForEach(func(k, v []byte) error {
			var rule models.Rule
			if err := json.Unmarshal(v, &rule); err != nil {
				return err
			}
			rules = append(rules, rule)
			return nil
		})
	})
	return rules, err
}

// SaveMetadata 保存元数据
func (s *Store) SaveMetadata(key string, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMetadata)
		return b.Put([]byte(key), value)
	})
}

// GetMetadata 获取元数据
func (s *Store) GetMetadata(key string) ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMetadata)
		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("metadata not found: %s", key)
		}
		value = make([]byte, len(data))
		copy(value, data)
		return nil
	})
	return value, err
}

// SaveLearnRecord 保存学习记录
func (s *Store) SaveLearnRecord(commitHash string, timestamp time.Time) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketLearnHistory)
		data, err := timestamp.MarshalBinary()
		if err != nil {
			return err
		}
		return b.Put([]byte(commitHash), data)
	})
}

// HasLearned 检查是否已学习过该提交
func (s *Store) HasLearned(commitHash string) (bool, error) {
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketLearnHistory)
		data := b.Get([]byte(commitHash))
		if data == nil {
			return fmt.Errorf("not found")
		}
		return nil
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GetLastLearnTime 获取上次学习时间
func (s *Store) GetLastLearnTime() (time.Time, error) {
	var lastTime time.Time
	data, err := s.GetMetadata("last_learn_time")
	if err != nil {
		return time.Time{}, nil // 没有记录返回零值
	}
	err = lastTime.UnmarshalBinary(data)
	return lastTime, err
}

// UpdateLastLearnTime 更新上次学习时间
func (s *Store) UpdateLastLearnTime(t time.Time) error {
	data, err := t.MarshalBinary()
	if err != nil {
		return err
	}
	return s.SaveMetadata("last_learn_time", data)
}
