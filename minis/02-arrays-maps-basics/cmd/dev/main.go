// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/arraysmapsbasics.go
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

	"github.com/example/go-10x-minis/minis/02-arrays-maps-basics/internal/arraysmapsbasics"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	input := `hello
world
hello
go`

	fmt.Println("=== Debugging FreqFromReader ===")
	fmt.Printf("Input:\n%s\n\n", input)

	reader := strings.NewReader(input)

	// Set breakpoint in arraysmapsbasics.go at FreqFromReader function
	freq, mostCommon, err := arraysmapsbasics.FreqFromReader(reader)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Frequency map: %v\n", freq)
	fmt.Printf("Most common word: %q\n", mostCommon)
}

