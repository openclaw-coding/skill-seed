# Skill-Seed

> 🌱 让项目 Skills 持续成长 - 为 AI Agents 学习项目编码模式的智能工具

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25%2B-brightgreen.svg)](https://go.dev/)

## 什么是 Skill-Seed？

Skill-Seed 是一个**项目级技能成长工具**，它从 Git 历史中学习你的团队编码模式，自动生成和更新 Claude Code（及其他 AI Agents）的 Skills，帮助 AI 更好地理解你的项目。

### 核心特性

- 🌱 **持续成长**: 从 Git 提交历史中学习编码模式，越来越聪明
- 🤖 **AI 深度分析**: 集成 Claude 进行深度代码分析
- 📊 **模式提取**: 自动识别命名、错误处理、架构等模式
- 🔄 **增量学习**: 只学习新的提交，不重复工作
- 🎯 **多 AI 支持**: 生成的 Skills 可被 Claude Code、Cursor 等 AI 工具使用
- 🌍 **国际化**: 完整的中英文支持

## 快速开始

### 安装

```bash
go install github.com/openclaw-coding/skill-seed@latest
```

### 初始化

在你的 Git 项目根目录运行：

```bash
skill-seed init
```

这会创建：
```
.seed/skill-seed/
├── SKILL.md              # 主技能文件
├── config.yaml           # 配置
├── memory/               # 学习数据库
└── references/           # 模式分类
    ├── naming-patterns/
    ├── error-handling-patterns/
    └── ...
```

### 学习模式

#### 1. 从 Git 历史学习

```bash
# 学习最近 30 天的提交
skill-seed learn

# 学习最近 100 个提交
skill-seed learn --max=100

# 学习最近 7 天的提交
skill-seed learn --since=7

# 强制重新学习
skill-seed learn --force
```

#### 2. 扫描当前项目

```bash
# 扫描当前项目状态
skill-seed scan
```

#### 3. 分析特定文件

```bash
skill-seed analyze main.go
skill-seed analyze src/
```

### 生成 Skills

```bash
# 生成 Claude Code Skills
skill-seed generate-skills
```

生成的 Skills 会被输出到 `~/.claude/skills/skill-seed-skills/`，Claude Code 会自动发现。

### 查看

```bash
# 查看学到的模式
skill-seed view patterns

# 查看生成的规则
skill-seed view rules
```

### Git Hooks

可选：安装 Git pre-commit hook，每次提交前自动运行检查：

```bash
# 安装 hook
skill-seed hook install

# 卸载 hook
skill-seed hook uninstall
```

## 工作流程

```bash
# 1. 初始化项目
skill-seed init

# 2. 学习历史模式
skill-seed learn --max=50

# 3. 生成 Skills
skill-seed generate-skills

# 4. 在 Claude Code 中使用
#（Skills 会自动被发现）

# 5. 继续学习...
skill-seed learn --since=7
skill-seed generate-skills
```

## 配置

编辑 `.seed/skill-seed/config.yaml`:

```yaml
project:
  name: my-project
  git_remote: git@github.com:user/repo.git
  initialized_at: 2026-03-19T00:00:00Z

claude:
  enabled: true
  timeout_seconds: 30
  fallback_to_basic: true

learning:
  max_history_analyze: 100
  min_samples_for_rule: 3
```

## 环境变量

```bash
# Claude API Key（用于深度分析）
export ANTHROPIC_API_KEY=your-api-key-here
```

## 生成的 Skills 结构

```
~/.claude/skills/skill-seed-skills/
├── SKILL.md                      # 主入口
└── references/                   # 分类知识
    ├── naming-patterns/          # 命名模式
    ├── error-handling-patterns/  # 错误处理
    ├── structure-patterns/       # 代码结构
    ├── concurrency-patterns/     # 并发模式
    └── testing-patterns/         # 测试模式
```

## 为什么选择 Skill-Seed？

### vs 静态文档

| 特性 | 静态文档 | Skill-Seed |
|------|---------|-----------|
| 更新方式 | 手动维护 | 自动从 Git 学习 |
| 准确性 | 可能过时 | 始终反映最新代码 |
| 代码示例 | 需要手动编写 | 自动提取 |
| 覆盖率 | 有限 | 全面分析 |

### vs LLM Context

| 特性 | LLM Context | Skill-Seed |
|------|------------|-----------|
| Token 消耗 | 每次都消耗 | 一次学习，多次使用 |
| 上下文限制 | 受限 | 不受限 |
| 成本 | 高 | 低 |
| 准确性 | 可能幻觉 | 基于实际代码 |

## 命令参考

```bash
skill-seed init              # 初始化
skill-seed learn [flags]      # 学习 Git 历史
skill-seed scan [flags]       # 扫描当前项目
skill-seed analyze <files>    # 分析特定文件
skill-seed check              # 手动运行检查
skill-seed generate-skills    # 生成 Skills
skill-seed view patterns      # 查看模式
skill-seed view rules         # 查看规则
skill-seed hook install       # 安装 hook
skill-seed hook uninstall     # 卸载 hook
```

## 项目状态

- ✅ 基础学习功能
- ✅ Claude 深度分析
- ✅ Skills 生成
- ✅ Git Hooks 集成
- ✅ 国际化支持（中英文）
- 🔄 更多 AI Agents 支持（开发中）

## 贡献

欢迎贡献！请随时提交 Issue 或 Pull Request。

## 许可证

Apache License 2.0

## 致谢

- [Claude Code](https://code.claude.com) - AI 编程助手
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [BoltDB](https://github.com/etcd-io/bbolt) - 嵌入式数据库

---

**让 AI 更懂你的代码！** 🌱
