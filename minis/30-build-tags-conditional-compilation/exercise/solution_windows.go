//go:build windows

package exercise

import "os"

// GetPathSeparator returns the Windows path separator
func GetPathSeparator() string {
	return "\\"
}

// GetHomeDirectory returns the home directory on Windows
func GetHomeDirectory() string {
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home
	}
	// Fallback to HOMEDRIVE + HOMEPATH
	if homeDrive := os.Getenv("HOMEDRIVE"); homeDrive != "" {
		if homePath := os.Getenv("HOMEPATH"); homePath != "" {
			return homeDrive + homePath
		}
	}
	return "C:\\Users\\Default"
}
