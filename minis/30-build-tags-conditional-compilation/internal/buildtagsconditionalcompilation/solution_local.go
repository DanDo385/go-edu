//go:build reference

package buildtagsconditionalcompilation

// GetStorageBackend returns the local storage backend
func GetStorageBackend() string {
	return "Local Filesystem"
}
