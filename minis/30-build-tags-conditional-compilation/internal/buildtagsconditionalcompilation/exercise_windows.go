//go:build !reference && windows

package buildtagsconditionalcompilation

import "os"

func GetPathSeparator() string { return "\\" }

func GetHomeDirectory() string { return os.Getenv("USERPROFILE") }
