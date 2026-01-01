//go:build reference && !cloud

package buildtags

// GetStorageBackend returns the local storage backend
func GetStorageBackend() string {
	return "Local Filesystem"
}
