//go:build unix

package exercise

import "os"

// GetPathSeparator returns the Unix path separator
func GetPathSeparator() string {
	return "/"
}

// GetHomeDirectory returns the home directory on Unix systems
func GetHomeDirectory() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return "/tmp" // Fallback
}
