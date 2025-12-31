// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/hellostrings.go
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

	"github.com/example/go-10x-minis/minis/01-hello-strings/internal/hellostrings"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	input := "hello world"

	fmt.Println("=== Debugging String Functions ===")
	fmt.Printf("Test input: %q\n\n", input)

	// Test TitleCase
	// Set breakpoint in hellostrings.go at TitleCase function
	// Expected: "hello world" → "Hello World"
	fmt.Println("--- Testing TitleCase ---")
	result1 := hellostrings.TitleCase(input)
	fmt.Printf("TitleCase(%q) = %q\n\n", input, result1)

	// Test Reverse
	// Set breakpoint in hellostrings.go at Reverse function
	// Expected: "hello world" → "dlrow olleh"
	fmt.Println("--- Testing Reverse ---")
	result2 := hellostrings.Reverse(input)
	fmt.Printf("Reverse(%q) = %q\n\n", input, result2)

	// Test RuneLen
	// Set breakpoint in hellostrings.go at RuneLen function
	// Expected: "hello world" → 11 runes
	fmt.Println("--- Testing RuneLen ---")
	result3 := hellostrings.RuneLen(input)
	fmt.Printf("RuneLen(%q) = %d\n", input, result3)
}
