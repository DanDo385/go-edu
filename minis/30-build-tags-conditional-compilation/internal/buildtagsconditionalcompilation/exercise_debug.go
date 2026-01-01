//go:build !reference && debug

package buildtagsconditionalcompilation

import "log"

func LogMessage(level, message string) {
	log.Printf("[%s] %s", level, message)
}

func IsLoggingEnabled() bool { return true }
