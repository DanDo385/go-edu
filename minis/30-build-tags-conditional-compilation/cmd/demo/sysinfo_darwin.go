//go:build darwin

package main

import (
	"os"
	"runtime"
	"strings"
)

// GetPlatformName returns the platform-specific name
func GetPlatformName() string {
	return "macOS (Apple Silicon & Intel)"
}

// GetUsername returns the current username using macOS-specific method
func GetUsername() string {
	// On macOS, try USER environment variable
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	return "unknown"
}

// GetSystemInfo returns macOS-specific system information
func GetSystemInfo() string {
	var info strings.Builder

	info.WriteString("Operating System: macOS\n")
	info.WriteString("Architecture: " + runtime.GOARCH + "\n")

	if runtime.GOARCH == "arm64" {
		info.WriteString("Processor: Apple Silicon (M1/M2/M3)\n")
	} else {
		info.WriteString("Processor: Intel\n")
	}

	info.WriteString("CPU Count: " + string(rune(runtime.NumCPU()+'0')) + "\n")
	info.WriteString("Package Manager: Homebrew (probably)\n")
	info.WriteString("Shell: $SHELL = " + os.Getenv("SHELL") + "\n")

	return info.String()
}

// GetSpecialFeature returns a macOS-specific feature description
func GetSpecialFeature() string {
	return "macOS features: Unix foundation, polished UI, great developer tools (Xcode)"
}
