//go:build debug

package main

import (
	"fmt"
	"runtime"
	"time"
)

// IsDebugEnabled returns whether debug mode is enabled
func IsDebugEnabled() bool {
	return true
}

// LogDebug logs detailed debug information
func LogDebug(component, message string) {
	timestamp := time.Now().Format("15:04:05.000")
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("[DEBUG] %s [%s:%d] %s: %s\n", timestamp, file, line, component, message)
}

// EnableVerboseLogging enables verbose output
func EnableVerboseLogging() {
	fmt.Println("🔍 Verbose logging enabled (debug build)")
	fmt.Println("   This provides detailed execution traces for debugging")
}

// DumpMemoryStats prints memory statistics (debug only)
func DumpMemoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\n=== Memory Statistics ===\n")
	fmt.Printf("Alloc = %v MB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
	fmt.Printf("========================\n\n")
}

// GetDebugInfo returns debug build information
func GetDebugInfo() string {
	return "DEBUG BUILD - includes verbose logging, memory stats, and performance tracing"
}
