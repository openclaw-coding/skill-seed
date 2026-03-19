package output

import (
	"fmt"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	Bold   = "\033[1m"
)

// Colors for different severity levels
const (
	ErrorColor   = Red
	WarningColor = Yellow
	InfoColor    = Blue
	SuccessColor = Green
)

// Colorize adds color to text
func Colorize(text string, color string) string {
	return color + text + Reset
}

// Error prints error message in red
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(Colorize(msg, ErrorColor))
}

// Warning prints warning message in yellow
func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(Colorize(msg, WarningColor))
}

// Info prints info message in blue
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(Colorize(msg, InfoColor))
}

// Success prints success message in green
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(Colorize(msg, SuccessColor))
}

// Dim prints dimmed message in gray
func Dim(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print(Colorize(msg, Gray))
}

// Print prints message without color
func Print(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println prints message with newline
func Println(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// SeverityLabel returns colored severity label
func SeverityLabel(severity string) string {
	switch severity {
	case "error":
		return Colorize("[ERROR]", ErrorColor)
	case "warning":
		return Colorize("[WARNING]", WarningColor)
	case "info":
		return Colorize("[INFO]", InfoColor)
	default:
		return fmt.Sprintf("[%s]", severity)
	}
}

