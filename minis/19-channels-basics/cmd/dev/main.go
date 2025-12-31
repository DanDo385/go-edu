// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/channelsbasics.go
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

	"github.com/example/go-10x-minis/minis/19-channels-basics/internal/channelsbasics"
)

func main() {
	fmt.Println("=== Debugging Channels ===")

	// Fixed default values - modify these directly if you want to test different inputs
	value := 42

	fmt.Printf("Testing Ping with value: %d\n", value)

	// Set breakpoint in channelsbasics.go at Ping function
	ch := channelsbasics.Ping(value)

	result := <-ch
	fmt.Printf("Received: %d\n", result)

	// Test PingPong
	in, out := channelsbasics.PingPong(3)
	in <- 1
	in <- 2
	in <- 3
	close(in)

	fmt.Println("PingPong results:")
	for v := range out {
		fmt.Printf("  %d\n", v)
	}
}

