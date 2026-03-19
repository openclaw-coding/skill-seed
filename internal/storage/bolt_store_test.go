package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/openclaw-coding/grow-check/pkg/models"
)

func TestNew(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// 测试创建存储
	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// 验证文件是否创建
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file not created")
	}
}

func TestPatternCRUD(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// 创建测试模式
	pattern := &models.CodePattern{
		ID:          "test-pattern-1",
		Type:        models.PatternNaming,
		Description: "Test pattern",
		Examples:    []string{"example1"},
		Frequency:   1,
		Confidence:  0.8,
		AutoFixable: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存
	if err := store.SavePattern(pattern); err != nil {
		t.Fatalf("Failed to save pattern: %v", err)
	}

	// 读取
	retrieved, err := store.GetPattern(pattern.ID)
	if err != nil {
		t.Fatalf("Failed to get pattern: %v", err)
	}

	if retrieved.ID != pattern.ID {
		t.Errorf("Expected ID %s, got %s", pattern.ID, retrieved.ID)
	}

	if retrieved.Description != pattern.Description {
		t.Errorf("Expected description %s, got %s", pattern.Description, retrieved.Description)
	}

	// 获取所有模式
	patterns, err := store.GetAllPatterns()
	if err != nil {
		t.Fatalf("Failed to get all patterns: %v", err)
	}

	if len(patterns) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(patterns))
	}
}

func TestRuleCRUD(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// 创建测试规则
	rule := &models.Rule{
		ID:         "test-rule-1",
		Name:       "Test Rule",
		Type:       models.PatternErrorHandling,
		Condition:  "error without log",
		Severity:   "warning",
		AutoFix:    true,
		Confidence: 0.9,
		Source:     "learned",
	}

	// 保存
	if err := store.SaveRule(rule); err != nil {
		t.Fatalf("Failed to save rule: %v", err)
	}

	// 获取所有规则
	rules, err := store.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get all rules: %v", err)
	}

	if len(rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(rules))
	}

	if rules[0].ID != rule.ID {
		t.Errorf("Expected ID %s, got %s", rule.ID, rules[0].ID)
	}
}

func TestMetadata(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	key := "test-key"
	value := []byte("test-value")

	// 保存元数据
	if err := store.SaveMetadata(key, value); err != nil {
		t.Fatalf("Failed to save metadata: %v", err)
	}

	// 读取元数据
	retrieved, err := store.GetMetadata(key)
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}

	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}
}

func TestLearnHistory(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	commitHash := "abc123"

	// 检查未学习
	learned, err := store.HasLearned(commitHash)
	if err != nil {
		t.Fatalf("Failed to check learned: %v", err)
	}
	if learned {
		t.Error("Should not be learned yet")
	}

	// 标记为已学习
	if err := store.SaveLearnRecord(commitHash, time.Now()); err != nil {
		t.Fatalf("Failed to save learn record: %v", err)
	}

	// 检查已学习
	learned, err = store.HasLearned(commitHash)
	if err != nil {
		t.Fatalf("Failed to check learned: %v", err)
	}
	if !learned {
		t.Error("Should be learned")
	}
}

func TestLastLearnTime(t *testing.T) {
	store, cleanup := createTestStore(t)
	defer cleanup()

	// 初始应该为零值
	lastTime, err := store.GetLastLearnTime()
	if err != nil {
		t.Fatalf("Failed to get last learn time: %v", err)
	}
	if !lastTime.IsZero() {
		t.Error("Initial last learn time should be zero")
	}

	// 更新学习时间
	now := time.Now()
	if err := store.UpdateLastLearnTime(now); err != nil {
		t.Fatalf("Failed to update last learn time: %v", err)
	}

	// 验证更新
	retrieved, err := store.GetLastLearnTime()
	if err != nil {
		t.Fatalf("Failed to get last learn time: %v", err)
	}

	// 允许 1 秒误差
	if retrieved.Sub(now) > time.Second {
		t.Errorf("Time mismatch: expected %v, got %v", now, retrieved)
	}
}

// 辅助函数：创建测试存储
func createTestStore(t *testing.T) (*Store, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cleanup := func() {
		store.Close()
	}

	return store, cleanup
}
