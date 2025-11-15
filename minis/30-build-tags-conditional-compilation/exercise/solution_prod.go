//go:build !debug

package exercise

// LogMessage is a no-op in production builds
func LogMessage(level, message string) {
	// No-op: no logging overhead in production
}

// IsLoggingEnabled returns false in production builds
func IsLoggingEnabled() bool {
	return false
}
