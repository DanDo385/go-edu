//go:build !solution
// +build !solution

package exercise

// This file contains exercises for learning build tags and conditional compilation.
// Implement the TODOs, then create platform-specific files with appropriate build tags.

// Exercise 1: Platform-Specific Path Separator
// Create separate files for different operating systems that implement GetPathSeparator:
// - exercise_unix.go (for Linux and macOS) - should return "/"
// - exercise_windows.go (for Windows) - should return "\\"
// Use the filename convention or explicit build tags.

// GetPathSeparator returns the OS-specific path separator
// TODO: Implement this function in platform-specific files
// func GetPathSeparator() string

// Exercise 2: Platform-Specific Home Directory
// Create platform-specific implementations that get the user's home directory:
// - On Unix: Use $HOME environment variable
// - On Windows: Use %USERPROFILE% environment variable

// GetHomeDirectory returns the user's home directory path
// TODO: Implement this function in platform-specific files
// func GetHomeDirectory() string

// Exercise 3: Feature Flags - Storage Backend
// Create two versions based on a "cloud" build tag:
// - exercise_cloud.go: Returns "S3" as the storage backend
// - exercise_local.go (no cloud tag): Returns "Local Filesystem"

// GetStorageBackend returns the configured storage backend
// TODO: Implement this function in feature-specific files
// func GetStorageBackend() string

// Exercise 4: Architecture-Specific Word Size
// Create architecture-specific implementations:
// - exercise_amd64.go: Returns 64
// - exercise_386.go: Returns 32
// - exercise_arm64.go: Returns 64
// - exercise_arm.go: Returns 32

// GetWordSize returns the architecture word size in bits
// TODO: Implement this function in architecture-specific files
// func GetWordSize() int

// Exercise 5: Complex Build Constraints
// Create a file that only compiles on Linux AND (amd64 OR arm64):
// - exercise_linux_64bit.go with build tag: //go:build linux && (amd64 || arm64)
// - It should implement IsLinux64Bit() returning true
// - Create a fallback file for other platforms returning false

// IsLinux64Bit returns true if running on 64-bit Linux
// TODO: Implement this function with complex build constraints
// func IsLinux64Bit() bool

// Exercise 6: Debug Logging
// Create debug-enabled and production versions:
// - exercise_debug.go (with debug tag): Actual logging implementation
// - exercise_prod.go (without debug tag): No-op implementation

// LogMessage logs a message if debug mode is enabled
// TODO: Implement this function in debug/production-specific files
// func LogMessage(level, message string)

// IsLoggingEnabled returns whether logging is enabled
// TODO: Implement this function in debug/production-specific files
// func IsLoggingEnabled() bool

// Instructions:
// 1. Read the TODOs above
// 2. Create separate .go files with appropriate build tags
// 3. Each file should implement one or more of the functions above
// 4. Run: go test (for default build)
// 5. Run: go test -tags=debug (to test debug build)
// 6. Run: go test -tags=cloud (to test cloud build)
// 7. Run: GOOS=windows go test (to test Windows build - may fail on Linux)
//
// Example file structure:
//   exercise_unix.go      - //go:build unix
//   exercise_windows.go   - //go:build windows
//   exercise_cloud.go     - //go:build cloud
//   exercise_local.go     - //go:build !cloud
//   exercise_debug.go     - //go:build debug
//   exercise_prod.go      - //go:build !debug
//   exercise_amd64.go     - Uses filename convention
//   exercise_386.go       - Uses filename convention
//   etc.
