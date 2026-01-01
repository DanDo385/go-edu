//go:build reference

package buildtagsconditionalcompilation

// GetStorageBackend returns the cloud storage backend
func GetStorageBackend() string {
	return "S3"
}
