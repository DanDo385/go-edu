//go:build reference && linux && (amd64 || arm64)

package buildtags

// IsLinux64Bit returns true on 64-bit Linux
func IsLinux64Bit() bool {
	return true
}
