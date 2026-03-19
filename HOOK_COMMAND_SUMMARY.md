# Hook 命令实现总结

## 完成时间
2026-03-19

## 需求
- `init` 命令默认不安装 git hooks
- 通过独立的 `skill-seed hook` 命令管理 hooks

## 实现内容

### 1. 新增 Hook 命令

**文件**: `cmd/skill-seed/hook/command.go`

```go
// 子命令
hook install   - 安装 Git pre-commit hook
hook uninstall - 卸载 Git pre-commit hook
```

### 2. Git 操作扩展

**文件**: `internal/git/operations.go`

新增 `UninstallPreCommitHook()` 方法：
- 智能卸载：如果 hook 只包含 skill-seed，直接删除
- 链式处理：如果 hook 与其他 hook 链接，只移除 skill-seed 部分
- 安全检查：卸载前检查是否存在 skill-seed hook

### 3. 国际化支持

**新增 i18n 键**：
- `cmd_hook_short` - "管理 Git hooks"
- `cmd_hook_long` - "安装或卸载 Git pre-commit hook"
- `cmd_hook_install_short` - "安装 Git pre-commit hook"
- `cmd_hook_install_long` - "在 .git/hooks/ 中安装 pre-commit hook"
- `cmd_hook_uninstall_short` - "卸载 Git pre-commit hook"
- `cmd_hook_uninstall_long` - "从 .git/hooks/ 中移除 pre-commit hook"
- `hook_installing` - "正在安装 Git pre-commit hook..."
- `hook_install_failed` - "Hook 安装失败: %v"
- `hook_installed` - "Hook 已安装"
- `hook_installed_success` - "每次提交前将自动运行 skill-seed check"
- `hook_uninstalling` - "正在卸载 Git pre-commit hook..."
- `hook_uninstall_failed` - "Hook 卸载失败: %v"
- `hook_uninstalled` - "Hook 已卸载"

### 4. 主命令注册

**文件**: `cmd/skill-seed/main.go`

添加了 hook 命令到根命令。

## 使用示例

### 初始化（不安装 hook）

```bash
$ skill-seed init
初始化 skill-seed
  [创建目录结构]
  [生成配置文件]
  [初始化内存数据库]
  [创建 SKILL.md]
  [创建模式分类]

skill-seed 初始化成功
```

**创建的目录结构**：
```
.skills/skill-seed/
├── SKILL.md
├── config.yaml
├── memory/
│   └── project.db
└── references/
    ├── naming-patterns/overview.md
    ├── error-handling-patterns/overview.md
    ├── structure-patterns/overview.md
    ├── concurrency-patterns/overview.md
    └── testing-patterns/overview.md
```

**注意**: 没有创建 `.git/hooks/pre-commit` 文件

### 安装 Hook

```bash
$ skill-seed hook install
正在安装 Git pre-commit hook...
Hook 已安装

每次提交前将自动运行 skill-seed check
```

**安装的 hook 文件** (`.git/hooks/pre-commit`):
```bash
#!/bin/sh
# skill-seed pre-commit hook
skill-seed check || exit $?
```

### 卸载 Hook

```bash
$ skill-seed hook uninstall
正在卸载 Git pre-commit hook...
Hook 已卸载
```

## 智能处理

### 场景 1: 新安装（没有现有 hook）

**安装前**:
```bash
$ ls .git/hooks/pre-commit
ls: No such file or directory
```

**执行**:
```bash
$ skill-seed hook install
```

**安装后**:
```bash
$ cat .git/hooks/pre-commit
#!/bin/sh
# skill-seed pre-commit hook
skill-seed check || exit $?
```

### 场景 2: 链式安装（已有其他 hook）

**安装前**:
```bash
$ cat .git/hooks/pre-commit
#!/bin/sh
# Existing linter
eslint src/
```

**执行**:
```bash
$ skill-seed hook install
```

**安装后**:
```bash
$ cat .git/hooks/pre-commit
#!/bin/sh
# Chain loading for multiple hooks (preserving existing hooks)

# Original hook content:
#!/bin/sh
# Existing linter
eslint src/

# skill-seed hook
skill-seed check || exit $?
```

### 场景 3: 重复安装

**执行**:
```bash
$ skill-seed hook install
# 第一次：安装成功
$ skill-seed hook install
# 第二次：检测到已存在，跳过
```

### 场景 4: 卸载链式 hook

**卸载前**:
```bash
$ cat .git/hooks/pre-commit
#!/bin/sh
# Chain loading for multiple hooks (preserving existing hooks)

# Original hook content:
#!/bin/sh
# Existing linter
eslint src/

# skill-seed hook
skill-seed check || exit $?
```

**执行**:
```bash
$ skill-seed hook uninstall
```

**卸载后**:
```bash
$ cat .git/hooks/pre-commit
#!/bin/sh
# Existing linter
eslint src/
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
  view            查看学习内容

$ skill-seed hook --help
Available Commands:
  install     安装 Git pre-commit hook
  uninstall   卸载 Git pre-commit hook
```

## 优势

1. ✅ **默认安全** - init 不修改 git hooks，用户需要显式安装
2. ✅ **独立管理** - hook 可以独立安装/卸载，不影响其他配置
3. ✅ **智能处理** - 支持链式 hook，不会覆盖现有 hooks
4. ✅ **完整支持** - 提供安装、卸载的完整功能
5. ✅ **国际化** - 中英文完整支持

## 测试验证

```bash
# 测试安装
$ skill-seed hook install
# ✓ 输出: Hook 已安装

# 验证文件存在
$ cat .git/hooks/pre-commit
# ✓ 内容正确

# 测试卸载
$ skill-seed hook uninstall
# ✓ 输出: Hook 已卸载

# 验证文件删除
$ ls .git/hooks/pre-commit
# ✓ 文件不存在
```

---

**总结**: 已实现独立的 hook 命令，init 默认不安装 hooks，通过 `skill-seed hook install/uninstall` 管理 hooks。
