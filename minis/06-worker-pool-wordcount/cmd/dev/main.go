package main

import "fmt"

/*
06-worker-pool-wordcount Debug Harness

This file demonstrates the module with fixed inputs.
Perfect for debugging with breakpoints.

How to debug:
  1. Set breakpoints in internal/workerpoolwordcount/exercise.go
  2. Press F5 in VS Code
  3. Select "Debug cmd/dev/main.go"
  4. Step through with F10/F11
*/
func main() {
	fmt.Println("=== 06-worker-pool-wordcount Debug Harness ===")
	fmt.Println()
	
	// BREAKPOINT: Set a breakpoint here
	fmt.Println("This debug harness demonstrates the module with fixed inputs.")
	fmt.Println()
	
	// TODO: Add example function calls here
	// Example:
	// import "minis/06-worker-pool-wordcount/internal/workerpoolwordcount"
	// 
	// input := "example"
	// result := workerpoolwordcount.YourFunction(input)
	// fmt.Printf("Result: %v\n", result)
	
	fmt.Println("To use this debug harness:")
	fmt.Println("1. Import the module: import \"minis/06-worker-pool-wordcount/internal/workerpoolwordcount\"")
	fmt.Println("2. Add function calls with example inputs")
	fmt.Println("3. Set breakpoints and debug with F5")
	fmt.Println()
	fmt.Println("See README.md for what this module demonstrates")
}
