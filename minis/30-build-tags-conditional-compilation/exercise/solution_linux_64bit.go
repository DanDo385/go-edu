//go:build linux && (amd64 || arm64)

package exercise

// IsLinux64Bit returns true on 64-bit Linux
func IsLinux64Bit() bool {
	return true
}
