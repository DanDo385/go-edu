//go:build (linux || darwin) && amd64 && !cgo

package main

// This file demonstrates complex build constraints using boolean logic:
// - Compiled on Linux OR macOS
// - Only on AMD64 architecture
// - Only when CGO is disabled

import "fmt"

// ShowAdvancedConstraints demonstrates the active build constraints
func ShowAdvancedConstraints() {
	fmt.Println("\n=== Advanced Build Constraints Demo ===")
	fmt.Println("This code was compiled with complex constraints:")
	fmt.Println("  ✓ OS: Linux or macOS")
	fmt.Println("  ✓ Architecture: AMD64")
	fmt.Println("  ✓ CGO: Disabled")
	fmt.Println("\nThis pattern is useful for:")
	fmt.Println("  - Pure Go applications (no C dependencies)")
	fmt.Println("  - Unix-like systems only")
	fmt.Println("  - 64-bit optimized code")
	fmt.Println("=======================================\n")
}
