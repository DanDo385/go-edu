//go:build !debug

package main

// IsDebugEnabled returns whether debug mode is enabled
func IsDebugEnabled() bool {
	return false
}

// LogDebug is a no-op in production builds
func LogDebug(component, message string) {
	// No-op: debug logging disabled in production
}

// EnableVerboseLogging is a no-op in production builds
func EnableVerboseLogging() {
	// No-op: verbose logging disabled in production
}

// DumpMemoryStats is a no-op in production builds
func DumpMemoryStats() {
	// No-op: memory stats disabled in production
}

// GetDebugInfo returns production build information
func GetDebugInfo() string {
	return "PRODUCTION BUILD - optimized for performance, minimal logging"
}
