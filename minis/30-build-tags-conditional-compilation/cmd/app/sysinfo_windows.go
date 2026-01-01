//go:build windows

package main

import (
	"os"
	"runtime"
	"strings"
)

// GetPlatformName returns the platform-specific name
func GetPlatformName() string {
	return "Windows (The Dominant Desktop OS)"
}

// GetUsername returns the current username using Windows-specific method
func GetUsername() string {
	// On Windows, USERNAME is the standard environment variable
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	// Fallback to USERPROFILE parsing
	if userprofile := os.Getenv("USERPROFILE"); userprofile != "" {
		parts := strings.Split(userprofile, "\\")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return "unknown"
}

// GetSystemInfo returns Windows-specific system information
func GetSystemInfo() string {
	var info strings.Builder

	info.WriteString("Operating System: Windows\n")
	info.WriteString("Architecture: " + runtime.GOARCH + "\n")
	info.WriteString("CPU Count: " + string(rune(runtime.NumCPU()+'0')) + "\n")
	info.WriteString("Computer Name: " + os.Getenv("COMPUTERNAME") + "\n")
	info.WriteString("User Domain: " + os.Getenv("USERDOMAIN") + "\n")
	info.WriteString("System Root: " + os.Getenv("SystemRoot") + "\n")

	// Windows-specific: check for WSL
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		info.WriteString("WSL: Running in Windows Subsystem for Linux\n")
	}

	return info.String()
}

// GetSpecialFeature returns a Windows-specific feature description
func GetSpecialFeature() string {
	return "Windows features: DirectX gaming, MS Office integration, WSL for Linux compatibility"
}
