//go:build !linux && !darwin && !windows

package main

import (
	"runtime"
)

// GetPlatformName returns the platform-specific name
func GetPlatformName() string {
	return "Other OS (" + runtime.GOOS + ")"
}

// GetUsername returns a fallback username
func GetUsername() string {
	return "unknown"
}

// GetSystemInfo returns generic system information
func GetSystemInfo() string {
	return "Operating System: " + runtime.GOOS + "\n" +
		"Architecture: " + runtime.GOARCH + "\n" +
		"This OS is not specifically supported, but Go runs here!\n"
}

// GetSpecialFeature returns a generic feature description
func GetSpecialFeature() string {
	return "Generic fallback for BSD, Solaris, or other Unix-like systems"
}
