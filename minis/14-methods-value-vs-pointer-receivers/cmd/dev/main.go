// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/methodsvaluevspointerreceivers.go
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

	"github.com/example/go-10x-minis/minis/14-methods-value-vs-pointer-receivers/internal/methodsvaluevspointerreceivers"
)

func main() {
	fmt.Println("=== Debugging Value vs Pointer Receivers ===")

	// Fixed default values - modify these directly if you want to test different inputs
	// Test TotalArea function
	shapes := []methodsvaluevspointerreceivers.Shape{
		&methodsvaluevspointerreceivers.Rectangle{Width: 5, Height: 3},
		&methodsvaluevspointerreceivers.Circle{Radius: 2},
	}

	fmt.Printf("Shapes: %+v\n", shapes)

	// Set breakpoint in methodsvaluevspointerreceivers.go at TotalArea function
	total := methodsvaluevspointerreceivers.TotalArea(shapes)
	fmt.Printf("Total area: %.2f\n", total)

	// Test SafeCounterMap
	counter := methodsvaluevspointerreceivers.NewSafeCounterMap()
	counter.Inc("test")
	count := counter.Get("test")
	fmt.Printf("Counter value: %d\n", count)
}

