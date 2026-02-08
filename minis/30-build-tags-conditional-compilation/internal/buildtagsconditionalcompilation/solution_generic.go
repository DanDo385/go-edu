//go:build reference && !amd64 && !arm && !arm64 && !386

package buildtagsconditionalcompilation

// GetWordSize returns a fallback word size for unknown architectures
func GetWordSize() int {
	// Conservative fallback
	return 32
}
