//go:build reference && local

package buildtagsconditionalcompilation

// GetStorageBackend returns the local storage backend
func GetStorageBackend() string {
	return "Local Filesystem"
}
