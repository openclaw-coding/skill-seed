# Skill-Seed

> 🌱 让项目 Skills 持续成长 - 为 AI Agents 学习项目编码模式的智能工具

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25%2B-brightgreen.svg)](https://go.dev/)

## 📖 目录

- [什么是 Skill-Seed？](#什么是-skill-seed)
- [架构设计](#架构设计)
- [快速开始](#快速开始)
- [使用指南](#使用指南)
- [命令参考](#命令参考)
- [配置](#配置)
- [项目结构](#项目结构)
- [常见问题](#常见问题)

---

## 什么是 Skill-Seed？

Skill-Seed 是一个**项目级技能成长工具**，它从 Git 历史中学习你的团队编码模式，自动生成和更新 Claude Code（及其他 AI Agents）的 Skills，帮助 AI 更好地理解你的项目。

### 核心特性

- 🌱 **持续成长**: 从 Git 提交历史中学习编码模式，越来越聪明
- 🤖 **AI 深度分析**: 集成 Claude 进行深度代码分析
- 📊 **模式提取**: 自动识别命名、错误处理、架构等模式
- 🔄 **增量学习**: 只学习新的提交，不重复工作
- 🎯 **多 AI 支持**: 生成的 Skills 可被 Claude Code、Cursor 等 AI 工具使用
- 🌍 **国际化**: 完整的中英文支持

---

## 架构设计

Skill-Seed 采用清晰的分层架构设计，遵循 Go 项目标准布局和最佳实践。

### 设计原则

1. **依赖倒置** - 高层模块不依赖低层模块，都依赖抽象
2. **单一职责** - 每个包只负责一件事
3. **接口隔离** - 使用小接口而非大接口
4. **开闭原则** - 对扩展开放，对修改关闭

### 分层架构

```
┌─────────────────────────────────────────────────────────┐
│                   CLI Layer (cmd)                        │
│          命令解析、参数处理、用户交互                     │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Service Layer (internal/service)            │
│         业务流程编排、事务管理、跨领域协调               │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Domain Layer (internal/domain)              │
│         核心业务模型、业务规则、领域逻辑                 │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│          Infrastructure Layer (internal/infra)           │
│         技术实现：Git、Storage、Config、Output           │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Resources Layer (embedfs)                   │
│           内嵌资源：Templates、Assets（go:embed）        │
└─────────────────────────────────────────────────────────┘
```

### 关键优势

1. **清晰的职责分离**: 每层只负责自己的职责
2. **易于扩展**: 添加新的 Agent 或命令很简单
3. **可测试性**: 每层可以独立测试
4. **依赖注入**: 通过 Container 管理依赖

---

## 快速开始

### 5 分钟上手

#### 1. 安装

```bash
# 下载最新版本
curl -sL https://github.com/openclaw-coding/skill-seed/releases/latest/download/skill-seed-$(uname -s)-$(uname -m) -o skill-seed
chmod +x skill-seed
sudo mv skill-seed /usr/local/bin/

# 验证安装
skill-seed --version
```

或从源码安装：
```bash
go install github.com/openclaw-coding/skill-seed@latest
```

#### 2. 初始化项目

```bash
cd /path/to/your/go-project
skill-seed init
```

这会创建 skill-seed 工具的工作目录：
```
.skill-seed/               # skill-seed 工具工作目录
├── config.yaml           # 配置文件
└── memory/               # 学习数据库
    └── project.db        # BoltDB 数据库（存储学到的模式）
```

生成的 Claude Code Skills 会输出到 `~/.claude/skills/skill-seed-skills/`。

#### 3. 首次学习

```bash
# 学习最近 30 天的提交历史
skill-seed learn --since=30d
```

输出示例：
```
从 Git 历史学习中

正在分析 150 个提交...
  [1/150] 正在分析 a1b2c3d4...
  [2/150] 正在分析 e5f6g7h8...

摘要:
  总提交数: 150
  已分析: 150
  新学到的模式: 12
已创建 5 个新规则
```

#### 4. 生成 Skills

```bash
skill-seed generate-skills
```

生成的 Skills 会被输出到 `~/.claude/skills/skill-seed-skills/`，Claude Code 会自动发现。

#### 5. 查看学习成果

```bash
# 查看学习到的模式
skill-seed view patterns

# 查看生成的规则
skill-seed view rules
```

---

## 使用指南

### 常见使用场景

#### 场景 1：首次在新项目使用

```bash
# 1. 初始化
cd new-project
skill-seed init

# 2. 从历史学习（如果项目已有历史）
skill-seed learn --max=200

# 3. 生成 Skills
skill-seed generate-skills

# 4. 开始正常开发
git add .
git commit -m "feat: add feature"
```

#### 场景 2：定期学习新提交

```bash
# 每周运行一次，学习新的提交
skill-seed learn --since=7
skill-seed generate-skills
```

#### 场景 3：手动检查（不提交）

```bash
# 想在提交前看看检查结果
skill-seed check
```

#### 场景 4：扫描整个项目

```bash
# 扫描当前项目状态
skill-seed scan

# 扫描所有文件类型
skill-seed scan --all
```

#### 场景 5：分析特定文件

```bash
# 分析单个文件
skill-seed analyze main.go

# 分析目录
skill-seed analyze src/

# 分析多个文件
skill-seed analyze *.go
```

### Git Hooks 集成

安装 Git pre-commit hook，每次提交前自动运行检查：

```bash
# 安装 hook
skill-seed hook --install

# 卸载 hook
skill-seed hook --uninstall
```

现在，每次 `git commit` 时都会自动检查：

```bash
git add .
git commit -m "Add new feature"
```

检查流程：
```
正在检查 3 个文件...

⚠ 发现 2 个问题:
1. ⚠ auth/login.go:42 - 错误处理未包含日志
2. ℹ database/connection.go:15 - 数据库连接缺少超时
```

### 性能优化建议

#### 1. 减少分析范围

编辑 `.skill-seed/config.yaml`：
```yaml
checking:
  exclude_patterns:
    - "vendor/*"
    - "*.pb.go"
    - "**/generated/*"
```

#### 2. 调整 Claude 超时

```yaml
agent:
  timeout: 15  # 减少到 15 秒
```

#### 3. 限制学习数量

```bash
# 只学习最近 50 个提交
skill-seed learn --max=50
```

---

## 命令参考

### 核心命令

```bash
skill-seed init              # 初始化项目
skill-seed learn [flags]      # 学习 Git 历史
  --limit=50                 # 学习最近 50 个提交
  --since=30d                # 学习最近 30 天的提交
  --all                      # 学习所有提交
  --force                    # 强制重新学习

skill-seed check [flags]      # 检查暂存文件
  --all                      # 检查所有文件

skill-seed generate-skills    # 生成 Claude Skills
  --output=DIR               # 指定输出目录
  --force                    # 强制覆盖

skill-seed view patterns      # 查看学习的模式
skill-seed view rules         # 查看生成的规则
```

### 辅助命令

```bash
skill-seed scan [flags]       # 扫描项目
  --all                      # 扫描所有文件类型

skill-seed analyze <files>    # 分析特定文件

skill-seed hook               # 管理 Git hooks
  --install                  # 安装 pre-commit hook
  --uninstall                # 卸载 pre-commit hook
```

---

## 配置

编辑 `.skill-seed/config.yaml`:

```yaml
project:
  name: my-project
  language: go
  locale: zh-CN            # 语言设置：zh-CN（中文）或 en-US（英文）
  git_remote: git@github.com:user/repo.git
  initialized_at: 2026-03-20T00:00:00Z

agent:
  type: claude              # Agent 类型：claude, gpt, local
  command: claude           # Claude 命令路径
  timeout: 60               # 超时时间（秒）

claude:
  enabled: true
  timeout_seconds: 60
  fallback_to_basic: true   # Claude 不可用时降级到基础模式

learning:
  max_commits: 50           # 默认分析的提交数量
  min_samples_for_rule: 3   # 生成规则所需最小样本数

checking:
  exclude_patterns:         # 排除的文件模式（不进行检查和分析）
    - vendor/*              # 依赖目录
    - node_modules/*        # Node.js 依赖
    - "*.pb.go"             # Protobuf 生成文件
    - "*.gen.go"            # 代码生成文件
    - "*/mocks/*"           # Mock 文件
    - "**/testdata/*"       # 测试数据

output:
  skills_path: ~/.claude/skills/skill-seed-skills  # Skills 输出路径
  default_language: go     # 默认分析的语言
```

### 配置字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `project.locale` | string | zh-CN | CLI 输出语言：zh-CN（中文）或 en-US（英文） |
| `learning.max_commits` | int | 50 | learn 命令的 `--limit` 默认值 |
| `learning.min_samples_for_rule` | int | 3 | 生成规则所需的最小样本数 |
| `checking.exclude_patterns` | []string | [vendor, node_modules, *.pb.go, etc.] | 检查和分析时排除的文件模式 |
| `output.skills_path` | string | ~/.claude/skills/skill-seed-skills | generate-skills 命令的 `--output` 默认值 |
| `output.default_language` | string | go | 默认分析的语言 |

### 环境变量

```bash
# Claude API Key（用于深度分析）
export ANTHROPIC_API_KEY=your-api-key-here
```

---

## 项目结构

### 目录结构

```
skill-seed/
├── cmd/skill-seed/           # CLI 命令层
│   ├── main.go               # 入口
│   └── internal/
│       ├── command/          # 所有命令
│       │   ├── init/
│       │   ├── learn/
│       │   ├── check/
│       │   ├── generate/
│       │   ├── hook/
│       │   ├── analyze/
│       │   ├── view/
│       │   └── scan/
│       ├── container/        # 依赖注入容器
│       └── utils/           # 工具函数
│
├── internal/
│   ├── agent/               # Agent 抽象层
│   │   ├── agent.go         # 接口定义
│   │   ├── types.go         # 类型定义
│   │   └── claude/          # Claude 实现
│   │
│   ├── service/             # 服务层
│   │   ├── learner.go       # 学习服务
│   │   ├── checker.go       # 检查服务
│   │   └── generator.go     # 生成服务
│   │
│   ├── domain/              # 领域层
│   │   ├── models.go        # 所有领域模型
│   │   └── repository.go    # 仓储接口
│   │
│   ├── infra/               # 基础设施层
│   │   ├── git/            # Git 操作
│   │   ├── storage/        # 存储
│   │   │   └── boltdb/     # BoltDB 实现
│   │   ├── config/         # 配置管理
│   │   └── output/         # 输出格式化
│   │
│   ├── templates/           # 模板加载器
│   │   ├── prompts/        # 提示词模板
│   │   └── skills/         # Skills 模板
│   │
│   └── i18n/                # 国际化
│       ├── i18n.go
│       └── locales/         # 翻译文件
│           ├── active.zh-CN.toml
│           └── active.en-US.toml
│
├── embedfs/                 # 内嵌资源
│   ├── embed.go
│   └── templates/           # 模板文件
│       ├── prompts/         # 提示词
│       └── skills/          # Skills
│
└── docs/                    # 文档
```

### 依赖关系

```
main.go
    ↓
container.go (依赖注入)
    ↓
├── agent/claude      (Agent 实现)
├── service/*         (业务服务)
├── infra/*           (基础设施)
└── templates/*       (模板加载器)
```

---

## 用户项目结构

在你的项目中使用 skill-seed 后，会产生以下目录结构：

### skill-seed 工作目录

```
your-project/
├── .skill-seed/                  # skill-seed 工具的工作目录
│   ├── config.yaml              # 配置文件
│   └── memory/                  # 学习数据库
│       └── project.db           # BoltDB（存储学到的模式）
└── ... (你的项目文件)
```

### 生成的 Claude Code Skills

```
~/.claude/skills/skill-seed-skills/    # Claude Code 使用的 Skills
├── SKILL.md                           # 主入口
├── README.md                          # 使用说明
└── references/                        # 分类知识
    ├── naming-patterns/               # 命名模式
    ├── error-handling-patterns/       # 错误处理
    ├── structure-patterns/            # 代码结构
    ├── concurrency-patterns/          # 并发模式
    └── testing-patterns/              # 测试模式
```

**说明**：
- `.skill-seed/` 是工具的工作目录，存储配置和学到的数据
- `~/.claude/skills/` 是生成的 Skills 目录，Claude Code 会自动加载
- 主要功能是从 Git 历史学习编码模式，"养成"高质量的 Skills

---

## 生成的 Skills 结构

```
~/.claude/skills/skill-seed-skills/
├── SKILL.md                      # 主入口
├── README.md                     # 使用说明
└── references/                   # 分类知识
    ├── naming-patterns/          # 命名模式
    ├── error-handling-patterns/  # 错误处理
    ├── structure-patterns/       # 代码结构
    ├── concurrency-patterns/     # 并发模式
    └── testing-patterns/         # 测试模式
```

---

## 常见问题

### 安装和使用

**Q: skill-seed 会影响提交速度吗？**
A: 会有轻微影响（通常 < 2 秒），但能避免潜在问题。

**Q: 可以在 CI/CD 中使用吗？**
A: 可以！在 CI 中运行 `skill-seed check`。

**Q: 支持非 Go 项目吗？**
A: 目前主要针对 Go，未来会支持更多语言。

**Q: 数据存储在哪里？**
A: 所有数据存储在项目的 `.skill-seed/memory/` 目录。

### 故障排除

#### 问题 1：Claude 不可用

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

#### 问题 2：Hook 未触发

**症状**：
提交时没有运行检查

**解决方案**：
1. 检查 hook 是否存在：
   ```bash
   cat .git/hooks/pre-commit
   ```

2. 重新安装：
   ```bash
   skill-seed hook --install
   ```

#### 问题 3：学习速度慢

**症状**：
`skill-seed learn` 运行很慢

**解决方案**：
1. 减少提交数量：
   ```bash
   skill-seed learn --max=50
   ```

2. 禁用 Claude 学习：
   ```yaml
   claude:
     enabled: false
   ```

---

## 最佳实践

1. **定期学习**: 每周运行 `skill-seed learn --since=7`
2. **审查规则**: 定期运行 `skill-seed view rules` 查看生成的规则
3. **调整配置**: 根据项目特点调整 `exclude_patterns`
4. **团队共享**: 将 `.skill-seed/config.yaml` 提交到 Git
5. **生成 Skills**: 每次学习后运行 `skill-seed generate-skills`

---

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

---

## 项目状态

- ✅ 基础学习功能
- ✅ Claude 深度分析
- ✅ Skills 生成
- ✅ Git Hooks 集成
- ✅ 国际化支持（中英文）
- ✅ 清晰的分层架构
- 🔄 更多 AI Agents 支持（开发中）

---

## 贡献

欢迎贡献！请随时提交 Issue 或 Pull Request。

---

## 许可证

Apache License 2.0

---

## 致谢

- [Claude Code](https://code.claude.com) - AI 编程助手
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [BoltDB](https://github.com/etcd-io/bbolt) - 嵌入式数据库
- [go-i18n](https://github.com/nicksnyder/go-i18n) - 国际化库

---

**让 AI 更懂你的代码！** 🌱
