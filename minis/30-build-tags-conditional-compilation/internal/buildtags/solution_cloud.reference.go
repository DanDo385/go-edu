//go:build reference && cloud

package buildtags

// GetStorageBackend returns the cloud storage backend
func GetStorageBackend() string {
	return "S3"
}
