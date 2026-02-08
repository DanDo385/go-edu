//go:build reference && arm

package buildtagsconditionalcompilation

// Using filename convention - automatically applies to arm architecture

// GetWordSize returns 32 for arm architecture
func GetWordSize() int {
	return 32
}
