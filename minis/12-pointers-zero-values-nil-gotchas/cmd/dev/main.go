package main

import (
	"fmt"
)

/*
minis/12-pointers-zero-values-nil-gotchas: cmd/dev (Debug Harness)

Fixed test inputs for debugging. No CLI arguments needed.

Usage:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5, select "Debug: cmd/dev (Debug Harness)"
  3. Step through with F10/F11

BREAKPOINT: Start here
*/

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  minis/12-pointers-zero-values-nil-gotchas: Debug Harness")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// BREAKPOINT: Fixed test inputs
	fmt.Println("Running with fixed test inputs...")
	fmt.Println()
	
	// TODO: Define test inputs and call exercise functions
	// Example:
	// testInput := "hello world"
	// result := pointerszerovaluesnilgotchas.SomeFunction(testInput)
	// fmt.Printf("Result: %v\n", result)
	
	fmt.Println("✓ Complete")
	fmt.Println()
	fmt.Println("Tips:")
	fmt.Println("  • F10 = Step Over")
	fmt.Println("  • F11 = Step Into") 
	fmt.Println("  • Watch Variables panel")
	fmt.Println("  • See internal/pointerszerovaluesnilgotchas/exercise.go for implementation")
}
