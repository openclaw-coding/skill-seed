# 快速开始指南

## 🎯 5 分钟上手 grow-check

### 第一步：安装

```bash
# 下载最新版本
curl -sL https://github.com/openclaw-coding/grow-check/releases/latest/download/grow-check-$(uname -s)-$(uname -m) -o grow-check
chmod +x grow-check
sudo mv grow-check /usr/local/bin/

# 验证安装
grow-check --version
```

### 第二步：初始化项目

```bash
cd /path/to/your/go-project
grow-check init
```

这会：
- ✅ 创建 `.skills/grow-check/` 目录
- ✅ 生成配置文件 `config.yaml`
- ✅ 安装 Git pre-commit 钩子
- ✅ 初始化 BoltDB 数据库

### 第三步：首次学习

```bash
# 学习最近 30 天的提交历史
grow-check learn --since=30d
```

输出示例：
```
🤖 Learning from Git history (last 30 days)...

📚 Analyzing 150 commits...
  [1/150] Analyzing a1b2c3d4...
  [2/150] Analyzing e5f6g7h8...
  
✨ Learned 12 new patterns
📏 Generated 5 rules
```

### 第四步：正常提交

现在，每次 `git commit` 时都会自动检查：

```bash
git add .
git commit -m "Add new feature"
```

检查流程：
```
🔍 Checking 3 files...
🤖 Analyzing with Claude...

⚠ Found 2 issues:
1. ⚠ auth/login.go:42 - 错误处理未包含日志
2. ℹ database/connection.go:15 - 数据库连接缺少超时

Options:
1. Auto-fix (recommended)
2. View details  
3. Ignore (with reason)
4. Abort commit

Your choice [1-4]: 1
```

### 第五步：查看学习成果

```bash
# 查看学习到的模式
grow-check patterns

# 查看生成的规则
grow-check rules
```

## 🎨 常见使用场景

### 场景 1：首次在新项目使用

```bash
# 1. 初始化
cd new-project
grow-check init

# 2. 从历史学习（如果项目已有历史）
grow-check learn --max=200

# 3. 开始正常开发
git add .
git commit -m "feat: add feature"
```

### 场景 2：定期学习新提交

```bash
# 每周运行一次，学习新的提交
grow-check learn --since=7
```

### 场景 3：手动检查（不提交）

```bash
# 想在提交前看看检查结果
grow-check check
```

### 场景 4：禁用交互式确认

编辑 `.skills/grow-check/config.yaml`：

```yaml
checking:
  interactive: false  # 改为 false
```

这样发现问题会直接报错，不再询问。

## ⚡ 性能优化建议

### 1. 减少分析范围

```yaml
checking:
  exclude_patterns:
    - "vendor/*"
    - "*.pb.go"
    - "**/generated/*"
```

### 2. 调整 Claude 超时

```yaml
claude:
  timeout_seconds: 15  # 减少到 15 秒
```

### 3. 限制学习数量

```yaml
learning:
  max_history_analyze: 500  # 减少到 500 个提交
```

## 🐛 故障排除

### 问题 1：Claude 不可用

**症状**：
```
⚠ Claude analysis failed: claude not available
```

**解决方案**：
1. 安装 Claude CLI：
   ```bash
   brew install claude
   ```

2. 或启用降级模式：
   ```yaml
   claude:
     fallback_to_basic: true
   ```

### 问题 2：Hook 未触发

**症状**：
提交时没有运行检查

**解决方案**：
1. 检查 hook 是否存在：
   ```bash
   cat .git/hooks/pre-commit
   ```

2. 重新安装：
   ```bash
   grow-check init  # 会重新安装 hook
   ```

### 问题 3：学习速度慢

**症状**：
`grow-check learn` 运行很慢

**解决方案**：
1. 减少提交数量：
   ```bash
   grow-check learn --max=50
   ```

2. 禁用 Claude 学习（只用基础规则）：
   ```yaml
   claude:
     enabled: false
   ```

## 📚 进阶用法

### 自定义规则

编辑 `.skills/grow-check/memory/rules.json`（未来版本支持）

### 导出/导入模式

```bash
# 导出模式（未来版本）
grow-check export > patterns.json

# 导入模式（未来版本）
grow-check import < patterns.json
```

## 🎯 最佳实践

1. **定期学习**：每周运行 `grow-check learn --since=7`
2. **审查规则**：定期运行 `grow-check rules` 查看生成的规则
3. **调整配置**：根据项目特点调整 `exclude_patterns`
4. **团队共享**：将 `.skills/grow-check/config.yaml` 提交到 Git

## ❓ 常见问题

**Q: grow-check 会影响提交速度吗？**
A: 会有轻微影响（通常 < 2 秒），但能避免潜在问题。

**Q: 可以在 CI/CD 中使用吗？**
A: 可以！在 CI 中运行 `grow-check check`。

**Q: 支持非 Go 项目吗？**
A: 目前主要针对 Go，未来会支持更多语言。

**Q: 数据存储在哪里？**
A: 所有数据存储在项目的 `.skills/grow-check/memory/` 目录。

---

**需要帮助？** 查看 [完整文档](README.md) 或 [提交 Issue](https://github.com/openclaw-coding/grow-check/issues)
