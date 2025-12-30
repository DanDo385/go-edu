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
	"strings"

	"github.com/example/go-10x-minis/minis/17-file-streaming-bufio/internal/exercise"
)

func main() {
	fmt.Println("=== Debugging File Streaming ===")

	// Fixed default values - modify these directly if you want to test different inputs
	input := `line one
line two
line three
line four
line five`

	reader := strings.NewReader(input)

	fmt.Printf("Input:\n%s\n\n", input)

	// Set breakpoint in exercise.go at CountLines function
	count, err := exercise.CountLines(reader)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Line count: %d\n", count)

	// Test WordFrequency
	reader2 := strings.NewReader(input)
	freq, err := exercise.WordFrequency(reader2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Word frequency: %+v\n", freq)
}

