//go:build linux

package main

import (
	"os"
	"runtime"
	"strings"
)

// GetPlatformName returns the platform-specific name
func GetPlatformName() string {
	return "Linux (Penguin Power!)"
}

// GetUsername returns the current username using Linux-specific method
func GetUsername() string {
	// On Linux, try USER environment variable
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	return "unknown"
}

// GetSystemInfo returns Linux-specific system information
func GetSystemInfo() string {
	var info strings.Builder

	info.WriteString("Operating System: Linux\n")
	info.WriteString("Architecture: " + runtime.GOARCH + "\n")
	info.WriteString("CPU Count: " + string(rune(runtime.NumCPU()+'0')) + "\n")
	info.WriteString("Package Manager: apt/yum/pacman (probably)\n")
	info.WriteString("Shell: $SHELL = " + os.Getenv("SHELL") + "\n")

	// Linux-specific: check if running in a container
	if _, err := os.Stat("/.dockerenv"); err == nil {
		info.WriteString("Container: Docker detected\n")
	}

	return info.String()
}

// GetSpecialFeature returns a Linux-specific feature description
func GetSpecialFeature() string {
	return "Linux features: open source kernel, extensive package management, strong server presence"
}
