//go:build !reference && !debug

package buildtagsconditionalcompilation

func LogMessage(level, message string) {}

func IsLoggingEnabled() bool { return false }
