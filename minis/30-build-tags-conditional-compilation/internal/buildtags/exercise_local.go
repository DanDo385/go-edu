//go:build !cloud

package buildtags

// GetStorageBackend returns the configured storage backend for local builds.
func GetStorageBackend() string {
	// TODO: Implement this function to return "Local Filesystem".
	return ""
}
