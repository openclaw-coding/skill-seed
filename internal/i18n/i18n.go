//go:build !en
// +build !en

package i18n

import "fmt"

// Messages contains all internationalized messages (Chinese - Default)
var Messages = map[string]string{
	// ========== 初始化相关 ==========
	"init_short": "初始化 grow-check 项目",
	"init_long": "在当前项目中初始化 grow-check。\n\n此命令将：\n  • 创建 .skills/grow-check/ 目录\n  • 生成默认配置\n  • 安装 Git pre-commit 钩子\n  • 初始化内存数据库\n\n请在 Git 仓库根目录运行此命令。",

	"init_creating_dirs":       "  [创建目录结构]",
	"init_generating_config":   "  [生成配置文件]",
	"init_initializing_db":     "  [初始化内存数据库]",
	"init_installing_hook":     "  [安装 Git pre-commit 钩子]",
	"init_creating_hook_script": "  [创建钩子脚本]",
	"init_creating_readme":     "  [创建说明文档]",
	"init_success":             "grow-check 初始化成功",
	"init_already_initialized": "grow-check 已经在此目录中初始化过",

	"init_skill_location": "Skill 位置: %s",
	"init_next_steps":     "后续步骤:",
	"init_step_learn":     "  1. 从历史学习: grow-check learn --since=30d",
	"init_step_watch":     "  2. 正常提交并观察它学习",
	"init_step_patterns":  "  3. 查看学到的模式: grow-check patterns",
	"init_step_rules":     "  4. 查看生成的规则: grow-check rules",

	// ========== 学习相关 ==========
	"learn_from_history": "从 Git 历史学习中 (最近 %d 天)...",
	"learn_failed":       "学习失败: %v",

	// ========== 检查相关 ==========
	"check_checking_files":   "正在检查 %d 个文件...",
	"check_analyzing_claude": "正在使用 Claude 进行深度分析...",
	"check_claude_failed":    "Claude 分析失败: %v",
	"check_no_issues":       "未发现问题",
	"check_found_issues":    "发现 %d 个问题:",

	// ========== 交互式选项 ==========
	"check_interactive_options": "请选择操作:",
	"check_option_autofix":     "  1. 自动修复 (推荐)",
	"check_option_details":     "  2. 查看详细信息",
	"check_option_ignore":      "  3. 忽略 (需提供原因)",
	"check_option_abort":       "  4. 终止提交",
	"check_choice_prompt":      "请选择 [1-4]: ",
	"check_autofix_disabled":   "配置中禁用了自动修复功能",
	"check_ignore_reason":      "请提供忽略原因: ",
	"check_ignored":           "问题已忽略。原因: %s",
	"check_aborted":           "用户已终止提交",
	"check_invalid_choice":    "无效的选择",

	// ========== 错误提示 ==========
	"check_init_failed":    "grow-check 未初始化",
	"check_init_hint":      "提示: 请先在项目根目录运行以下命令进行初始化:",
	"check_init_command":   "   grow-check init",
	"check_init_more_info": "更多信息: https://github.com/openclaw-coding/grow-check",

	"init_failed":  "初始化失败: %v",
	"check_failed": "检查失败: %v",
	"check_error":   "错误: %v",

	// ========== 命令描述 ==========
	"cmd_init_short":     "初始化 grow-check",
	"cmd_init_long":      "在当前项目中初始化 grow-check 作为项目级 skill",
	"cmd_learn_short":    "从 Git 历史学习",
	"cmd_learn_long":     "分析 Git 提交历史并学习编码模式",
	"cmd_check_short":    "手动运行 pre-commit 检查",
	"cmd_check_long":     "运行与 pre-commit 钩子相同的检查",
	"cmd_analyze_short":  "分析特定文件或目录",
	"cmd_analyze_long":   "在不检查 Git 状态的情况下分析特定文件",
	"cmd_view_short":     "查看学习内容",
	"cmd_view_long":      "查看学习到的模式和生成的规则",
	"cmd_generate_short": "生成 Claude Code skills",
	"cmd_generate_long":  "从学习到的模式生成 Claude Code skills",

	// ========== 其他消息 ==========
	"msg_analyzing":  "正在分析...",
	"msg_generating": "正在生成...",
	"msg_loading":    "正在加载...",
	"msg_saving":     "正在保存...",
	"msg_done":       "完成",
	"msg_warning":    "警告",
	"msg_error":      "错误",
	"msg_info":       "信息",

	// ========== 初始化命令详情 ==========
	"init_start":                   "初始化 grow-check",
	"init_hook_install_failed":     "  钩子安装失败: %v",
	"init_hook_install_manual":     "  可以稍后手动安装",
	"init_hook_installed":          "  钩子已安装",

	// ========== 学习命令详情 ==========
	"learn_no_commits":           "没有新的提交需要学习",
	"learn_last_learn_time":      "   上次学习时间: %s",
	"learn_since_time":           "   起始时间: %s",
	"learn_tip_force":            "提示: 使用 --force 参数重新学习所有提交",
	"learn_analyzing_commits":    "正在分析 %d 个提交...",
	"learn_analyzing_commit":     "  [%d/%d] 正在分析 %s...",
	"learn_get_diff_failed":      "    获取差异失败: %v",
	"learn_learning_failed":      "    学习失败: %v",
	"learn_save_pattern_failed":  "    保存模式失败: %v",
	"learn_summary":              "\n摘要:",
	"learn_total_commits":        "  总提交数: %d",
	"learn_skipped_commits":      "  已跳过（已学习）: %d",
	"learn_analyzed_count":       "  已分析: %d",
	"learn_new_patterns":         "  新学到的模式: %d",
	"learn_generating_rules":     "\n从模式生成规则...",
	"learn_create_rule_failed":   "  创建规则失败: %v",
	"learn_rules_created":        "已创建 %d 个新规则",
	"learn_no_rules_created":     "未创建规则（模式需要更多样本）",
	"learn_no_patterns":          "尚未学习到任何模式",
	"learn_patterns_header":      "学到的模式 (共 %d 个):",
	"learn_pattern_item":         "%d. [%s] %s",
	"learn_pattern_details":      "   置信度: %.2f | 频率: %d | 可自动修复: %v",
	"learn_pattern_example":      "   示例: %s",
	"learn_no_rules":             "尚未生成任何规则",
	"learn_rules_header":         "有效规则 (共 %d 个):",
	"learn_rule_item":            "%d. %s [%s]",
	"learn_rule_details":         "   来源: %s | 置信度: %.2f",

	// ========== 检查命令详情 ==========
	"check_read_file_failed":     "警告: 无法读取 %s: %v",
	"check_no_valid_files":       "没有有效的文件可分析",
	"check_autofix_not_impl":     "自动修复未实现: %s:%d",
	"check_fixed_count":          "已修复 %d 个问题",
	"check_autofix_not_ready":    "自动修复功能尚未完全实现",

	// ========== 生成命令详情 ==========
	"generate_exists":            "Skills 已存在于: %s",
	"generate_use_force":         "使用 --force 覆盖",
	"generate_generating":        "生成 Claude Code skills...",
	"generate_success":           "Skills 生成成功",
	"generate_next_steps":        "后续步骤:",
	"generate_step1":             "  1. 查看生成的 skills: %s",
	"generate_step2":             "  2. 在 Claude Code 中测试: /grow-check-skills",
	"generate_step3":             "  3. 提交到版本控制（可选）",
	"generate_output_location":   "已生成到: %s",
	"generate_patterns_count":    "学到的模式数: %d",
	"generate_avg_confidence":    "平均置信度: %.1f%%",
	"generate_file_created":      "  已生成: %s",
}

// Get returns the message for the given key
func Get(key string) string {
	if msg, ok := Messages[key]; ok {
		return msg
	}
	return key
}

// Getf returns the formatted message for the given key
func Getf(key string, args ...interface{}) string {
	msg := Get(key)
	return fmt.Sprintf(msg, args...)
}
