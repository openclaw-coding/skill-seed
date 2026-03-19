# skill-seed Skills Template

这是一个能够**持续成长**的 Claude Code Skills 模板，通过 skill-seed 工具从 Git 历史学习项目编码模式后自动生成。

## 🎯 核心理念

不同于静态的编码规范文档，这个 skills 模板是**活的**：

- 📚 **从历史学习**：分析 Git 提交历史，提取真实的编码模式
- 🔄 **持续进化**：每次提交都可能学习到新模式
- 📊 **置信度系统**：模式随着出现频率增加而变得更可信
- 🤖 **AI 驱动**：使用 Claude 深度分析代码，理解模式背后的意图

## 📁 目录结构

```
.template/skills/skill-seed-skills/
├── SKILL.md                                    # 主入口文件
└── references/                                 # 详细模式文档
    ├── naming-patterns.md                      # 命名规范
    ├── error-handling-patterns.md              # 错误处理
    ├── structure-patterns.md                   # 代码结构
    ├── concurrency-patterns.md                 # 并发模式
    └── testing-patterns.md                     # 测试模式
```

## 🚀 工作流程

### 1. 初始化阶段

```bash
# 在你的项目中
skill-seed init

# 学习现有代码历史
skill-seed learn --max=100

# 生成初始 skills
skill-seed generate-skills
```

### 2. 持续成长

```bash
# 每周学习新模式
skill-seed learn --since=7d

# 更新 skills
skill-seed generate-skills

# 查看 learning 进度
skill-seed status
```

### 3. AI 使用

当 AI 在此项目工作时：

1. **自动加载** skills 中的模式
2. **遵循**学习到的编码规范
3. **检测**违反模式的代码
4. **建议**符合项目风格的改进

## 📝 模板变量说明

所有模板文件使用 `{{VARIABLE_NAME}}` 格式的占位符，由 skill-seed 填充：

### 全局变量

- `{{TIMESTAMP}}` - 生成时间
- `{{PROJECT_NAME}}` - 项目名称
- `{{GIT_REMOTE}}` - Git 仓库地址
- `{{VERSION}}` - skill-seed 版本
- `{{COMMITS_ANALYZED}}` - 分析的提交数

### 统计变量

- `{{PATTERN_COUNT}}` - 学习到的模式数量
- `{{EXAMPLE_COUNT}}` - 示例数量
- `{{CONFIDENCE}}` - 置信度
- `{{CONSISTENCY_RATE}}` - 一致性率
- `{{AVG_CONFIDENCE}}` - 平均置信度

### 模式变量

每个模式类型都有特定的变量，例如：

#### 命名模式
- `{{FILE_NAMING_PATTERN}}` - 文件命名模式
- `{{FILE_NAMING_GOOD_EXAMPLES}}` - 好的示例
- `{{FILE_NAMING_BAD_EXAMPLES}}` - 坏的示例
- `{{FILE_NAMING_RULE}}` - 命名规则

#### 错误处理模式
- `{{ERROR_CHECK_PATTERN}}` - 错误检查模式
- `{{ERROR_WRAPPING_PATTERN}}` - 错误包装模式
- `{{ERROR_MESSAGE_PATTERN}}` - 错误消息模式

（完整变量列表见各模板文件）

## 🎨 模式分类

根据 `pkg/models/patterns.go` 中的定义，支持以下模式类型：

1. **naming** - 命名规范
   - 文件和目录命名
   - 变量和常量命名
   - 函数和方法命名
   - 接口和类型命名

2. **error_handling** - 错误处理
   - 错误检查
   - 错误包装
   - 错误消息
   - 错误恢复

3. **structure** - 代码结构
   - 项目布局
   - 包组织
   - 文件结构
   - 层次架构

4. **concurrency** - 并发模式
   - Goroutine 使用
   - 通道模式
   - 同步原语
   - Context 使用

5. **testing** - 测试模式
   - 测试文件组织
   - 测试命名
   - 测试结构
   - Mock 和 Stub

6. **comment** - 注释规范
   - 函数注释
   - 包注释
   - 注释风格

## 📊 置信度系统

模式的置信度分为几个阶段：

```
0.0 - 0.3: 新模式 (New Pattern)
  ↓
0.3 - 0.6: 新兴模式 (Emerging Pattern)
  ↓
0.6 - 0.8: 确立模式 (Established Pattern)
  ↓
0.8 - 1.0: 团队规范 (Team Convention) → 成为规则
```

只有高置信度的模式才会被 AI 严格遵循。

## 🔧 生成器实现

`skill-seed generate-skills` 命令应该：

1. **读取模板**：加载 `.template/skills/skill-seed-skills/` 中的模板文件
2. **获取模式**：从 BoltDB 读取学习到的模式
3. **分析统计**：计算各种统计指标
4. **填充变量**：替换模板中的占位符
5. **输出结果**：生成到 `~/.claude/skills/skill-seed-skills/`

### 伪代码

```go
func GenerateSkills(skillPath string) error {
    // 1. 加载模板
    templates := loadTemplates(".template/skills/skill-seed-skills/")

    // 2. 获取学习数据
    patterns := getAllPatterns()
    stats := calculateStats(patterns)

    // 3. 准备变量
    vars := prepareTemplateVars(patterns, stats)

    // 4. 填充模板
    for _, template := range templates {
        content := renderTemplate(template, vars)
        outputPath := getOutputPath(template)
        writeFile(outputPath, content)
    }

    return nil
}
```

## 📖 使用示例

### 示例 1：AI 遵循命名模式

**没有 skill-seed-skills**：
```go
// AI 可能生成
var user_data string
func Get_User_Info() {}
```

**有 skill-seed-skills**：
```go
// AI 会遵循学习到的模式
var userData string
func getUserInfo() {}
```

### 示例 2：AI 遵循错误处理模式

**没有 skill-seed-skills**：
```go
if err != nil {
    return err
}
```

**有 skill-seed-skills**：
```go
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

## 🎯 与 jzero-skills 的对比

| 特性 | jzero-skills | skill-seed-skills |
|------|--------------|-------------------|
| 内容来源 | 框架规范文档 | 项目实际代码 |
| 更新方式 | 手动维护 | 自动学习 |
| 适用范围 | 框架级别 | 项目级别 |
| 内容类型 | 最佳实践指导 | 真实模式提取 |
| 置信度 | 无 | 有（基于频率） |

两者是**互补**的：
- `jzero-skills` 提供框架层面的通用规范
- `skill-seed-skills` 提供项目层面的特定模式

## 🚀 下一步

1. **实现生成器**：在 `internal/generator/` 中实现 `generate-skills` 命令
2. **集成 Claude**：改进模式识别的准确性
3. **添加更多模式**：支持更多模式类型
4. **可视化**：生成模式可视化报告
5. **交互式学习**：允许用户确认或拒绝学到的模式

## 🤝 贡献

欢迎贡献！请查看：
- [CONTRIBUTING.md](../../CONTRIBUTING.md)
- [项目 README](../../README.md)

## 📄 许可证

MIT License - 详见 [LICENSE](../../LICENSE)

---

**Made with ❤️ by the skill-seed team**
