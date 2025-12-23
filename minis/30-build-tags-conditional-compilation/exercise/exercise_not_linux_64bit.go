//go:build !linux || (!amd64 && !arm64)

package exercise

// IsLinux64Bit returns true if running on 64-bit Linux.
func IsLinux64Bit() bool {
	// TODO: Implement this function to return false as the fallback.
	return false
}
