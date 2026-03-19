//go:build en
// +build en

package i18n

import "fmt"

// Messages contains all internationalized messages (English)
var Messages = map[string]string{
	// ========== Initialization ==========
	"init_short": "Initialize grow-check",
	"init_long": "Initialize grow-check in the current project.\n\nThis will:\n  • Create .skills/grow-check/ directory\n  • Generate default configuration\n  • Install Git pre-commit hook\n  • Initialize memory database\n\nRun this command in the root directory of your Git repository.",

	"init_creating_dirs":       "  [Creating directory structure]",
	"init_generating_config":   "  [Generating configuration]",
	"init_initializing_db":     "  [Initializing memory database]",
	"init_installing_hook":     "  [Installing Git pre-commit hook]",
	"init_creating_hook_script": "  [Creating hook script]",
	"init_creating_readme":     "  [Creating README]",
	"init_success":             "grow-check initialized successfully",
	"init_already_initialized": "grow-check already initialized in this directory",

	"init_skill_location": "Skill location: %s",
	"init_next_steps":     "Next steps:",
	"init_step_learn":     "  1. Learn from history: grow-check learn --since=30d",
	"init_step_watch":     "  2. Make commits and watch it learn",
	"init_step_patterns":  "  3. View patterns: grow-check patterns",
	"init_step_rules":     "  4. View rules: grow-check rules",

	// ========== Learning ==========
	"learn_from_history": "Learning from Git history (last %d days)...",
	"learn_failed":       "Learning failed: %v",

	// ========== Checking ==========
	"check_checking_files":   "Checking %d files...",
	"check_analyzing_claude": "Analyzing with Claude...",
	"check_claude_failed":    "Claude analysis failed: %v",
	"check_no_issues":       "No issues found",
	"check_found_issues":    "Found %d issues:",

	// ========== Interactive Options ==========
	"check_interactive_options": "Please choose:",
	"check_option_autofix":     "  1. Auto-fix (recommended)",
	"check_option_details":     "  2. View details",
	"check_option_ignore":      "  3. Ignore (with reason)",
	"check_option_abort":       "  4. Abort commit",
	"check_choice_prompt":      "Your choice [1-4]: ",
	"check_autofix_disabled":   "Auto-fix is disabled in config",
	"check_ignore_reason":      "Please provide a reason: ",
	"check_ignored":           "Issues ignored. Reason: %s",
	"check_aborted":           "Commit aborted by user",
	"check_invalid_choice":    "Invalid choice",

	// ========== Error Messages ==========
	"check_init_failed":    "grow-check not initialized",
	"check_init_hint":      "Hint: Please run the following command in the project root directory to initialize:",
	"check_init_command":   "   grow-check init",
	"check_init_more_info": "For more information: https://github.com/openclaw-coding/grow-check",

	"init_failed":  "Init failed: %v",
	"check_failed": "Check failed: %v",
	"check_error":   "Error: %v",

	// ========== Command Descriptions ==========
	"cmd_init_short":     "Initialize grow-check",
	"cmd_init_long":      "Initialize grow-check as a project-level skill",
	"cmd_learn_short":    "Learn from Git history",
	"cmd_learn_long":     "Analyze Git commit history and learn code patterns",
	"cmd_check_short":    "Run pre-commit check manually",
	"cmd_check_long":     "Run the same checks that the pre-commit hook would run",
	"cmd_analyze_short":  "Analyze specific files or directories",
	"cmd_analyze_long":   "Analyze specific files without checking Git status",
	"cmd_view_short":     "View learned content",
	"cmd_view_long":      "View learned patterns and generated rules",
	"cmd_generate_short": "Generate Claude Code skills",
	"cmd_generate_long":  "Generate Claude Code skills from learned patterns",

	// ========== Other Messages ==========
	"msg_analyzing":  "Analyzing...",
	"msg_generating": "Generating...",
	"msg_loading":    "Loading...",
	"msg_saving":     "Saving...",
	"msg_done":       "Done",
	"msg_warning":    "Warning",
	"msg_error":      "Error",
	"msg_info":       "Info",

	// ========== Init Command Details ==========
	"init_start":                   "Initializing grow-check",
	"init_hook_install_failed":     "  Hook installation failed: %v",
	"init_hook_install_manual":     "  You can install it manually later",
	"init_hook_installed":          "  Hook installed",

	// ========== Learn Command Details ==========
	"learn_no_commits":           "No new commits to learn from",
	"learn_last_learn_time":      "   Last learn time: %s",
	"learn_since_time":           "   Since: %s",
	"learn_tip_force":            "Tip: Use --force flag to re-learn all commits",
	"learn_analyzing_commits":    "Analyzing %d commits...",
	"learn_analyzing_commit":     "  [%d/%d] Analyzing %s...",
	"learn_get_diff_failed":      "    Failed to get diff: %v",
	"learn_learning_failed":      "    Learning failed: %v",
	"learn_save_pattern_failed":  "    Failed to save pattern: %v",
	"learn_summary":              "\nSummary:",
	"learn_total_commits":        "  Total commits: %d",
	"learn_skipped_commits":      "  Skipped (already learned): %d",
	"learn_analyzed_count":       "  Analyzed: %d",
	"learn_new_patterns":         "  New patterns learned: %d",
	"learn_generating_rules":     "\nGenerating rules from patterns...",
	"learn_create_rule_failed":   "  Failed to create rule: %v",
	"learn_rules_created":        "Created %d new rules",
	"learn_no_rules_created":     "No rules created (patterns need more samples)",
	"learn_no_patterns":          "No patterns learned yet",
	"learn_patterns_header":      "Learned patterns (%d total):",
	"learn_pattern_item":         "%d. [%s] %s",
	"learn_pattern_details":      "   Confidence: %.2f | Frequency: %d | Auto-fixable: %v",
	"learn_pattern_example":      "   Example: %s",
	"learn_no_rules":             "No rules generated yet",
	"learn_rules_header":         "Active rules (%d total):",
	"learn_rule_item":            "%d. %s [%s]",
	"learn_rule_details":         "   Source: %s | Confidence: %.2f",

	// ========== Check Command Details ==========
	"check_read_file_failed":     "Warning: failed to read %s: %v",
	"check_no_valid_files":       "No valid files to analyze",
	"check_autofix_not_impl":     "Auto-fix not implemented for %s:%d",
	"check_fixed_count":          "Fixed %d issues",
	"check_autofix_not_ready":    "Auto-fix not fully implemented yet",

	// ========== Generate Command Details ==========
	"generate_exists":            "Skills already exist at: %s",
	"generate_use_force":         "Use --force to overwrite",
	"generate_generating":        "Generating Claude Code skills...",
	"generate_success":           "Skills generated successfully",
	"generate_next_steps":        "Next steps:",
	"generate_step1":             "  1. Review generated skills: %s",
	"generate_step2":             "  2. Test with Claude Code: /grow-check-skills",
	"generate_step3":             "  3. Commit to version control (optional)",
	"generate_output_location":   "Generated to: %s",
	"generate_patterns_count":    "Patterns learned: %d",
	"generate_avg_confidence":    "Average confidence: %.1f%%",
	"generate_file_created":      "  Generated: %s",
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
