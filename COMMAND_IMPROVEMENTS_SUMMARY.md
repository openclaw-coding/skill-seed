# Grow-Check 命令改进总结

## 完成时间
2026-03-19

## 改进内容

### 1. 子命令支持 Help ✅

**状态**: Cobra 框架默认支持

所有子命令都支持 `--help` 参数：

```bash
$ skill-seed learn --help
$ skill-seed analyze --help
$ skill-seed scan --help
$ skill-seed hook --help
```

### 2. 不同命令显示不同的提示信息 ✅

**改进前**：
```bash
$ skill-seed learn --max=5
从 Git 历史学习中 (最近 30 天)...  # 固定显示，不准确
```

**改进后**：根据参数动态显示不同的提示

```bash
# 使用 --max 参数（指定提交数）
$ skill-seed learn --max=5
正在分析最近 5 个提交...

# 使用 --since 参数（指定天数）
$ skill-seed learn --since=7
正在分析最近 7 天的提交...

# 使用 --force 参数（强制重新学习）
$ skill-seed learn --force --max=3
强制重新学习最近 3 个提交...

# 默认参数
$ skill-seed learn
正在分析最近 30 天的提交...
```

### 3. 新增 Scan 命令 - 分析当前项目 ✅

**新命令**: `skill-seed scan`

**功能**：
- 扫描当前项目的所有 Go 文件
- 使用 Claude 分析代码模式
- 标记当前 commit 为已学习
- 往前推 commit 时会推到最近的已学习 commit

**使用示例**：

```bash
# 扫描 Go 文件
$ skill-seed scan
正在分析当前项目...
  项目根目录: /Users/Apple/GolandProjects/examples

找到 227 个文件

正在使用 Claude 进行深度分析...
未发现问题
已标记当前提交为已学习: eecefee1
项目扫描完成

# 扫描所有文件
$ skill-seed scan --all
```

**实现细节**：

1. **文件扫描**：
   - 默认只扫描 `.go` 文件
   - 使用 `--all` 扫描所有文件
   - 自动跳过 `.git/`, `.skills/`, `vendor/`, `node_modules/` 等目录

2. **代码分析**：
   - 使用现有的 `checker.AnalyzeFiles()` 进行分析
   - 支持 Claude 深度分析
   - 支持模式学习

3. **标记已学习**：
   - 获取当前 commit hash
   - 调用 `storage.SaveLearnRecord()` 标记为已学习
   - 下次 `learn` 命令会从这里继续

4. **与 Learn 命令配合**：
   ```
   场景1: 直接使用 learn
   $ skill-seed learn --max=10
   # 从最近 10 个未学习的 commit 开始学习

   场景2: 先 scan，再 learn
   $ skill-seed scan
   # 标记当前 commit 为已学习

   $ skill-seed learn --max=10
   # 从当前 commit 往前推 10 个 commit（包括当前已标记的）
   ```

## 新增的文件

1. `cmd/skill-seed/scan/command.go` - Scan 命令实现
2. `internal/git/operations.go` - 添加 `GetCurrentCommitHash()` 方法

## 新增的 i18n 消息

### 中文
```go
"cmd_scan_short":            "扫描并分析当前项目",
"cmd_scan_long":             "分析当前项目状态并学习代码模式",
"scan_analyzing_project":     "正在分析当前项目...",
"scan_project_root":         "  项目根目录: %s",
"scan_no_files":             "未找到可分析的文件",
"scan_found_files":          "找到 %d 个文件",
"scan_marked_learned":       "已标记当前提交为已学习: %s",
"scan_completed":            "项目扫描完成",

"learn_msg_recent_days":     "正在分析最近 %d 天的提交...",
"learn_msg_recent_max":      "正在分析最近 %d 个提交...",
"learn_msg_force_max":       "强制重新学习最近 %d 个提交...",
"learn_msg_force_days":      "强制重新学习最近 %d 天的提交...",
"learn_msg_force_all":       "强制重新学习所有提交...",
"learn_no_commits_force_hint": "所有提交都已学习完成",
"learn_no_commits_days_hint": "最近 %d 天内没有新的提交",
"learn_no_commits_max_hint":  "最近 %d 个提交都已学习过",
```

### 英文
```go
"cmd_scan_short":            "Scan and analyze current project",
"cmd_scan_long":             "Analyze current project state and learn code patterns",
"scan_analyzing_project":     "Analyzing current project...",
"scan_project_root":         "  Project root: %s",
"scan_no_files":             "No files found to analyze",
"scan_found_files":          "Found %d files",
"scan_marked_learned":       "Marked current commit as learned: %s",
"scan_completed":            "Project scan completed",

"learn_msg_recent_days":     "Analyzing commits from last %d days...",
"learn_msg_recent_max":      "Analyzing last %d commits...",
"learn_msg_force_max":       "Force re-learning last %d commits...",
"learn_msg_force_days":      "Force re-learning commits from last %d days...",
"learn_msg_force_all":       "Force re-learning all commits...",
"learn_no_commits_force_hint": "All commits have been learned",
"learn_no_commits_days_hint": "No new commits in last %d days",
"learn_no_commits_max_hint":  "Last %d commits are already learned",
```

## 完整命令列表

```bash
$ skill-seed --help
Available Commands:
  analyze         分析特定文件或目录
  check           手动运行 pre-commit 检查
  generate-skills 生成 Claude Code skills
  help            Help about any command
  hook            管理 Git hooks
  init            初始化 skill-seed
  learn           从 Git 历史学习
  scan            扫描并分析当前项目  ← 新增
  view            查看学习内容
```

## 使用场景

### 场景1: 学习 Git 历史
```bash
# 学习最近 30 天的提交
skill-seed learn

# 学习最近 100 个提交
skill-seed learn --max=100

# 强制重新学习
skill-seed learn --force
```

### 场景2: 分析当前项目
```bash
# 扫描当前项目（标记为已学习）
skill-seed scan

# 继续往前学习
skill-seed learn --max=50
```

### 场景3: 分析特定文件
```bash
# 分析特定文件
skill-seed analyze main.go

# 分析目录
skill-seed analyze src/
```

## 技术改进

1. **动态消息生成**：
   - `printLearningStartMessage()` 根据参数生成不同消息
   - `printNoCommitsMessage()` 根据情况显示不同的提示

2. **Commit 追踪**：
   - `GetCurrentCommitHash()` 获取当前 commit
   - `SaveLearnRecord()` 标记已学习
   - `GetLastLearnTime()` 获取上次学习时间

3. **智能文件扫描**：
   - 自动过滤不需要的目录
   - 支持只扫描 Go 文件或所有文件
   - 生成相对路径保持可读性

---

**总结**: 所有三个需求都已完成，命令系统更加完善和用户友好。
