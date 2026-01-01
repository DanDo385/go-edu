//go:build !reference && !cloud

package buildtagsconditionalcompilation

func GetStorageBackend() string { return "Local Filesystem" }
