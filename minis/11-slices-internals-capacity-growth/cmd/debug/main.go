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

	"github.com/example/go-10x-minis/minis/11-slices-internals-capacity-growth/internal/exercise"
)

func main() {
	fmt.Println("=== Debugging Slice Internals ===")

	// Fixed default values - modify these directly if you want to test different inputs
	s := []int{1, 2, 3}
	elem := 4

	fmt.Printf("Initial slice: %v (len=%d, cap=%d)\n", s, len(s), cap(s))
	fmt.Printf("Appending: %d\n\n", elem)

	// Set breakpoint in exercise.go at GrowSlice function
	newSlice, oldCap, newCap := exercise.GrowSlice(s, elem)

	fmt.Printf("Result: %v (len=%d, cap=%d)\n", newSlice, len(newSlice), cap(newSlice))
	fmt.Printf("Old capacity: %d\n", oldCap)
	fmt.Printf("New capacity: %d\n", newCap)
	fmt.Printf("Reallocation occurred: %v\n", newCap > oldCap)
}

