//go:build !linux || (!amd64 && !arm64)

package exercise

// IsLinux64Bit returns false on non-Linux or non-64-bit platforms
func IsLinux64Bit() bool {
	return false
}
