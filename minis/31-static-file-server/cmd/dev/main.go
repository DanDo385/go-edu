package main

import "fmt"

/*
31-static-file-server Debug Harness

This file demonstrates the module with fixed inputs.
Perfect for debugging with breakpoints.

How to debug:
  1. Set breakpoints in internal/staticfileserver/exercise.go
  2. Press F5 in VS Code
  3. Select "Debug cmd/dev/main.go"
  4. Step through with F10/F11
*/
func main() {
	fmt.Println("=== 31-static-file-server Debug Harness ===")
	fmt.Println()
	
	// BREAKPOINT: Set a breakpoint here
	fmt.Println("This debug harness demonstrates the module with fixed inputs.")
	fmt.Println()
	
	// TODO: Add example function calls here
	// Example:
	// import "minis/31-static-file-server/internal/staticfileserver"
	// 
	// input := "example"
	// result := staticfileserver.YourFunction(input)
	// fmt.Printf("Result: %v\n", result)
	
	fmt.Println("To use this debug harness:")
	fmt.Println("1. Import the module: import \"minis/31-static-file-server/internal/staticfileserver\"")
	fmt.Println("2. Add function calls with example inputs")
	fmt.Println("3. Set breakpoints and debug with F5")
	fmt.Println()
	fmt.Println("See README.md for what this module demonstrates")
}
