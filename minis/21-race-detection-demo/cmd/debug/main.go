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
	"sync"
)

func main() {
	fmt.Println("=== Debugging Race Detection ===")

	// Fixed default values - modify these directly if you want to test different inputs
	// Simple counter example
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	fmt.Println("Testing concurrent counter increment...")

	// Launch 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d (expected: 10)\n", counter)
	fmt.Println("\nNote: Use cmd/app/main.go for comprehensive race detection demos.")
	fmt.Println("Run with: go run -race ./cmd/app/main.go")
}

