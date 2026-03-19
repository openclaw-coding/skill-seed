//go:build !cn
// +build !cn

package i18n

// Messages contains all internationalized messages
var Messages = map[string]string{
	"check_init_failed":        "grow-check not initialized",
	"check_init_hint":          "Hint: Please run the following command in the project root directory to initialize:",
	"check_init_command":       "   grow-check init",
	"check_init_more_info":     "For more information: https://github.com/openclaw-coding/grow-check",
	"init_creating_dirs":       "Creating directory structure...",
	"init_generating_config":   "Generating configuration...",
	"init_initializing_db":     "Initializing memory database...",
	"init_installing_hook":     "Installing Git pre-commit hook...",
	"init_creating_hook":       "Creating hook script...",
	"init_creating_readme":     "Creating README...",
	"init_success":             "grow-check initialized successfully!",
	"init_skill_location":      "Skill location: %s",
	"init_next_steps":          "Next steps:",
	"init_step_learn":          "  1. Learn from history: grow-check learn --since=30d",
	"init_step_watch":          "  2. Make commits and watch it learn!",
	"init_step_patterns":       "  3. View patterns: grow-check patterns",
	"init_step_rules":          "  4. View rules: grow-check rules",
	"learn_from_history":       "Learning from Git history (last %d days)...",
	"check_failed":             "Check failed: %v",
	"check_checking_files":     "Checking %d files...",
	"check_analyzing_claude":   "Analyzing with Claude...",
	"check_claude_failed":      "Claude analysis failed: %v",
	"check_no_issues":          "No issues found",
	"check_found_issues":       "Found %d issues:",
	"check_interactive_options": "Options:",
	"check_option_autofix":     "1. Auto-fix (recommended)",
	"check_option_details":     "2. View details",
	"check_option_ignore":      "3. Ignore (with reason)",
	"check_option_abort":       "4. Abort commit",
	"check_choice_prompt":      "Your choice [1-4]: ",
	"check_autofix_disabled":   "Auto-fix is disabled in config",
	"check_ignore_reason":      "Please provide a reason: ",
	"check_ignored":            "Issues ignored. Reason: %s",
	"check_aborted":            "commit aborted by user",
	"check_invalid_choice":     "Invalid choice",
	"init_failed":              "Init failed: %v",
	"learn_failed":             "Learn failed: %v",
	"check_error":              "Error: %v",
}

// Get returns the message for the given key
func Get(key string) string {
	if msg, ok := Messages[key]; ok {
		return msg
	}
	return key
}
