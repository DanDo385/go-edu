// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/pointerszerovaluesnilgotchas.go
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

	"github.com/example/go-10x-minis/minis/12-pointers-zero-values-nil-gotchas/internal/pointerszerovaluesnilgotchas"
)

func main() {
	fmt.Println("=== Debugging Pointers and Zero Values ===")

	// Fixed default values - modify these directly if you want to test different inputs
	x := 42

	fmt.Printf("Value: %d\n", x)
	fmt.Printf("Address: %p\n", &x)

	// Set breakpoint in pointerszerovaluesnilgotchas.go to debug pointer operations
	// This is a simple test - actual exercises may have different function signatures
	fmt.Println("\nNote: Set breakpoints in internal/exercise/pointerszerovaluesnilgotchas.go")
	fmt.Println("to debug specific pointer operations and zero value behaviors.")
}

