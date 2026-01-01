package main

import (
	"fmt"
	"os"
)

/*
minis/08-http-client-retries: cmd/app

CLI application demonstrating the concepts from this module.

Usage:
  go run ./cmd/app/main.go [args...]

See README.md for specific usage examples and argument details.

BREAKPOINT: Set breakpoints at "// BREAKPOINT:" comments for debugging.
*/

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  minis/08-http-client-retries")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// BREAKPOINT: Inspect command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: See README.md for command-line argument details")
		fmt.Printf("Example: %s [args...]\n", os.Args[0])
		fmt.Println()
		fmt.Println("This module demonstrates: httpclientretries concepts")
		fmt.Println("See internal/httpclientretries/ for implementation details")
		os.Exit(1)
	}

	// Parse arguments (project-specific)
	args := os.Args[1:]
	
	// BREAKPOINT: Exercise execution point
	fmt.Println("Running exercise...")
	fmt.Printf("Arguments: %v\n", args)
	fmt.Println()
	
	// TODO: Add exercise-specific logic here
	// Import and call functions from internal/httpclientretries/exercise.go
	
	fmt.Println("✓ Complete")
	fmt.Println()
	fmt.Println("See internal/httpclientretries/exercise.go for implementation")
}
