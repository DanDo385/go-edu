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
	"context"
	"fmt"
	"time"

	"github.com/example/go-10x-minis/minis/16-context-cancellation-timeouts/internal/exercise"
)

func main() {
	fmt.Println("=== Debugging Context Cancellation ===")

	// Fixed default values - modify these directly if you want to test different inputs
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	urls := []string{"https://example.com", "https://golang.org"}

	fmt.Printf("Testing FetchAll with timeout: %v\n", 2*time.Second)
	fmt.Printf("URLs: %v\n\n", urls)

	// Set breakpoint in exercise.go at FetchAll function
	results, err := exercise.FetchAll(ctx, urls, 1*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Results: %+v\n", results)
}

