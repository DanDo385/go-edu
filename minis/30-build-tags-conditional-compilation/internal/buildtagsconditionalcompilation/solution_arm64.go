// Using filename convention - automatically applies to arm64 architecture

package buildtagsconditionalcompilation

// GetWordSize returns 64 for arm64 architecture
func GetWordSize() int {
	return 64
}
