//go:build !cloud

package exercise

// GetStorageBackend returns the local storage backend
func GetStorageBackend() string {
	return "Local Filesystem"
}
