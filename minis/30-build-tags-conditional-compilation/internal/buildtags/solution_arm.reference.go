// Using filename convention - automatically applies to arm architecture

package buildtags

// GetWordSize returns 32 for arm architecture
func GetWordSize() int {
	return 32
}
