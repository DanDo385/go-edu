//go:build !reference && !(linux && (amd64 || arm64))

package buildtagsconditionalcompilation

func IsLinux64Bit() bool { return false }
