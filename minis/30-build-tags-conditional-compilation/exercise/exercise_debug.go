//go:build debug

package exercise

import "log"

// LogMessage logs a message if debug mode is enabled.
func LogMessage(level, message string) {
	// TODO: Implement this function to print a log message.
	// Example: log.Printf("[%s] %s", level, message)
}

// IsLoggingEnabled returns whether logging is enabled.
func IsLoggingEnabled() bool {
	// TODO: Implement this function to return true for debug builds.
	return false
}
