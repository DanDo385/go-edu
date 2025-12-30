// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/exercise.go
// 2. Open this file (cmd/debug/main.go)
// 3. Use "Debug Main Program (Current Package)" configuration
// 4. Press F5 - that's it! The debugger will stop at your breakpoints
//
// Usage:
//   go run ./cmd/debug
//   # Or just press F5 in VS Code

package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/09-http-server-graceful/internal/exercise"
)

func main() {
	fmt.Println("=== Debugging HTTP Server Store ===")

	// Set breakpoint in exercise.go at NewMemStore function
	store := exercise.NewMemStore()

	// Default test operations
	store.Set("key1", "value1")
	fmt.Printf("Set(%q, %q)\n", "key1", "value1")

	store.Set("key2", "value2")
	fmt.Printf("Set(%q, %q)\n", "key2", "value2")

	// Test Get operation
	val, ok := store.Get("key1")
	fmt.Printf("Get('key1'): val=%q, ok=%v\n", val, ok)

	// Test Delete operation
	store.Delete("key1")
	val, ok = store.Get("key1")
	fmt.Printf("Get('key1') after delete: val=%q, ok=%v\n", val, ok)
}

