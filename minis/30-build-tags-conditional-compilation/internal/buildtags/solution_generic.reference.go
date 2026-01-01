//go:build reference && !amd64 && !386 && !arm64 && !arm

package buildtags

// GetWordSize returns a fallback word size for unknown architectures
func GetWordSize() int {
	// Conservative fallback
	return 32
}
