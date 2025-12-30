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
	fmt.Println("=== Debugging Worker Pool with Backpressure ===")

	// Fixed default values - modify these directly if you want to test different inputs
	jobs := make(chan int, 3) // Small buffer
	done := make(chan bool)

	// Producer
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("Sending job %d...\n", i)
			jobs <- i
		}
		close(jobs)
		done <- true
	}()

	// Consumer
	go func() {
		for job := range jobs {
			fmt.Printf("Processing job %d...\n", job)
			time.Sleep(100 * time.Millisecond) // Simulate work
		}
	}()

	<-done
	fmt.Println("All jobs processed!")
	fmt.Println("\nNote: Use cmd/app/main.go for comprehensive backpressure demos.")
}

