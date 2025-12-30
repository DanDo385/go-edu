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
	fmt.Println("=== Debugging Bounded Channels and Semaphores ===")

	// Fixed default values - modify these directly if you want to test different inputs
	// Simple bounded channel example
	ch := make(chan int, 3) // Capacity 3

	fmt.Println("Testing bounded channel (capacity 3)...")

	// Fill channel
	for i := 1; i <= 3; i++ {
		ch <- i
		fmt.Printf("Sent: %d\n", i)
	}

	// Receive
	for i := 0; i < 3; i++ {
		v := <-ch
		fmt.Printf("Received: %d\n", v)
	}

	fmt.Println("Done!")
	fmt.Println("\nNote: Use cmd/app/main.go for comprehensive semaphore demos.")
}

