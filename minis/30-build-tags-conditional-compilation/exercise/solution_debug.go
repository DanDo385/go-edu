//go:build debug

package exercise

import (
	"fmt"
	"time"
)

// LogMessage logs a message with timestamp when debug is enabled
func LogMessage(level, message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s: %s\n", timestamp, level, message)
}

// IsLoggingEnabled returns true when debug build tag is active
func IsLoggingEnabled() bool {
	return true
}
