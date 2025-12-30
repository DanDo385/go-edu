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

	"github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/internal/exercise"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	ctx := context.Background()
	workers := 2
	urls := []string{"https://example.com", "https://golang.org"}

	fmt.Println("=== Debugging WordCount ===")
	fmt.Printf("Workers: %d\n", workers)
	fmt.Printf("URLs: %v\n\n", urls)

	// Set breakpoint in exercise.go at WordCount function
	// Note: This will make actual HTTP requests, so ensure network access
	counts, err := exercise.WordCount(ctx, urls, workers)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Word counts: %+v\n", counts)
}

