//go:build reference && 386

package buildtagsconditionalcompilation

// Using filename convention - automatically applies to 386 architecture

// GetWordSize returns 32 for 386 architecture
func GetWordSize() int {
	return 32
}
