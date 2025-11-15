//go:build cloud

package exercise

// GetStorageBackend returns the cloud storage backend
func GetStorageBackend() string {
	return "S3"
}
