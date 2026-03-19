//go:build cn
// +build cn

package i18n

// Messages contains all internationalized messages (Chinese)
var Messages = map[string]string{
	"check_init_failed":        "grow-check 未初始化",
	"check_init_hint":          "💡 提示: 请先在项目根目录运行以下命令初始化:",
	"check_init_command":       "   grow-check init",
	"check_init_more_info":     "📚 更多信息请查看: https://github.com/openclaw-coding/grow-check",
	"init_creating_dirs":       "  创建目录结构...",
	"init_generating_config":   "  生成配置...",
	"init_initializing_db":     "  初始化内存数据库...",
	"init_installing_hook":     "  安装 Git pre-commit 钩子...",
	"init_creating_hook":       "  创建钩子脚本...",
	"init_creating_readme":     "  创建 README...",
	"init_success":             "✅ grow-check 初始化成功!",
	"init_skill_location":      "📁 Skill 位置: %s",
	"init_next_steps":          "后续步骤:",
	"init_step_learn":          "  1. 从历史学习: grow-check learn --since=30d",
	"init_step_watch":          "  2. 提交代码并观察它学习!",
	"init_step_patterns":       "  3. 查看模式: grow-check patterns",
	"init_step_rules":          "  4. 查看规则: grow-check rules",
	"learn_from_history":       "🤖 从 Git 历史学习中 (最近 %d 天)...\n\n",
	"check_failed":             "❌ Check failed: %v\n",
	"check_checking_files":     "🔍 Checking %d files...\n",
	"check_analyzing_claude":   "🤖 Analyzing with Claude...\n",
	"check_claude_failed":      "⚠ Claude analysis failed: %v\n",
	"check_no_issues":          "✅ No issues found\n",
	"check_found_issues":       "\n⚠ Found %d issues:\n\n",
	"check_interactive_options": "Options:",
	"check_option_autofix":     "1. Auto-fix (recommended)",
	"check_option_details":     "2. View details",
	"check_option_ignore":      "3. Ignore (with reason)",
	"check_option_abort":       "4. Abort commit",
	"check_choice_prompt":      "Your choice [1-4]: ",
	"check_autofix_disabled":   "⚠ Auto-fix is disabled in config",
	"check_ignore_reason":      "Please provide a reason: ",
	"check_ignored":            "✅ Issues ignored. Reason: %s\n",
	"check_aborted":            "commit aborted by user",
	"check_invalid_choice":     "Invalid choice",
	"init_failed":              "❌ Init failed: %v\n",
	"learn_failed":             "❌ Learn failed: %v\n",
	"check_error":              "Error: %v\n",
}

// Get returns the message for the given key
func Get(key string) string {
	if msg, ok := Messages[key]; ok {
		return msg
	}
	return key
}
