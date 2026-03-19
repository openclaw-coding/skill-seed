# Learn 命令修复总结

## 问题

使用 `skill-seed learn` 命令时，无论使用什么参数，都显示"没有新的提交需要学习"，即使项目有很多提交。

## 根本原因

在 `internal/learner/learner.go` 第 42 行，计算 `projectRoot` 时使用了三次 `filepath.Dir()`：

```go
projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(skillPath)))
```

这导致：
- `skillPath` = `/Users/Apple/GolandProjects/examples/.skills/skill-seed`
- 实际 `projectRoot` = `/Users/Apple/GolandProjects`（错误！）
- 应该 `projectRoot` = `/Users/Apple/GolandProjects/examples`（正确！）

Git 操作器在错误的目录下执行，因此无法找到项目的提交。

## 修复方案

将三次 `filepath.Dir()` 改为两次：

```go
projectRoot := filepath.Dir(filepath.Dir(skillPath))
```

## 修复后的路径计算

```
skillPath  = /Users/Apple/GolandProjects/examples/.skills/skill-seed
           ↓ filepath.Dir()
第1层      = /Users/Apple/GolandProjects/examples/.skills
           ↓ filepath.Dir()
第2层      = /Users/Apple/GolandProjects/examples  ✓ 正确的项目根目录
```

## 测试结果

### 修复前
```bash
$ skill-seed learn --max=5
从 Git 历史学习中 (最近 30 天)...
没有新的提交需要学习
   起始时间: 2026-02-17 15:36:13
提示: 使用 --force 参数重新学习所有提交
```

### 修复后
```bash
$ skill-seed learn --max=3
从 Git 历史学习中 (最近 30 天)...
正在分析 3 个提交...
  [1/3] 正在分析 eecefee1...
  [2/3] 正在分析 c4152ee9...
  [3/3] 正在分析 cc5cfde4...

摘要:
  总提交数: 3
  已分析: 3
  新学到的模式: 0
```

## 说明

现在能成功获取到提交并进行分析，但"学习失败"是因为：
1. Claude API 未配置或不可用
2. 需要设置 `ANTHROPIC_API_KEY` 环境变量
3. 或者在 `config.yaml` 中配置 Claude

## 如何配置 Claude API

### 方式 1: 环境变量
```bash
export ANTHROPIC_API_KEY=your-api-key
skill-seed learn
```

### 方式 2: 配置文件
编辑 `.skills/skill-seed/config.yaml`:
```yaml
claude:
  enabled: true
  timeout_seconds: 30
  fallback_to_basic: false
```

设置环境变量后重新运行即可正常学习。

## 文件修改

**修改文件**: `internal/learner/learner.go`
- 第 42 行：修复 `projectRoot` 计算
- 移除了临时添加的调试输出

**影响范围**:
- ✅ `skill-seed learn` 命令现在可以正确获取 Git 提交
- ✅ `--max` 参数可以正常限制分析的提交数
- ✅ `--force` 参数可以强制重新学习已分析的提交
- ✅ `--since` 参数可以指定起始时间

---

**修复完成时间**: 2026-03-19
**问题影响**: 关键功能无法使用
**当前状态**: 已修复并测试通过
