//go:build reference && amd64

package buildtagsconditionalcompilation

// Using filename convention - automatically applies to amd64 architecture

// GetWordSize returns 64 for amd64 architecture
func GetWordSize() int {
	return 64
}
