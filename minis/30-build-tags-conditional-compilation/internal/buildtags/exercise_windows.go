//go:build windows

package buildtags

// GetPathSeparator returns the OS-specific path separator for Windows.
func GetPathSeparator() string {
	// TODO: Implement this function to return "\".
	return ""
}

// GetHomeDirectory returns the user's home directory path on Windows.
func GetHomeDirectory() string {
	// TODO: Implement this function to get the home directory from the %USERPROFILE% environment variable.
	// You can use the `os` package for this.
	return ""
}
