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
	"time"
)

func main() {
	fmt.Println("=== Debugging Goroutines ===")

	// Fixed default values - modify these directly if you want to test different inputs
	numGoroutines := 10

	fmt.Printf("Creating %d goroutines...\n", numGoroutines)

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("Goroutine %d completed\n", id)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	fmt.Println("All goroutines completed!")
}

