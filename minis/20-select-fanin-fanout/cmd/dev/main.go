// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/selectfaninfanout.go
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
	fmt.Println("=== Debugging Select, Fan-In, and Fan-Out ===")

	// Fixed default values - modify these directly if you want to test different inputs
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Send values
	go func() {
		ch1 <- 1
		ch2 <- 2
		close(ch1)
		close(ch2)
	}()

	// Use select to receive from multiple channels
	fmt.Println("Select demonstration:")
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1:
			fmt.Printf("  Received from ch1: %d\n", v)
		case v := <-ch2:
			fmt.Printf("  Received from ch2: %d\n", v)
		case <-time.After(1 * time.Second):
			fmt.Println("  Timeout")
		}
	}

	fmt.Println("Done!")
}

