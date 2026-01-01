//go:build reference

package buildtagsconditionalcompilation

// IsLinux64Bit returns false on non-Linux or non-64-bit platforms
func IsLinux64Bit() bool {
	return false
}
