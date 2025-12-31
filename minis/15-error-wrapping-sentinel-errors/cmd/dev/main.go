// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/errorwrappingsentinelerrors.go
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

	"github.com/example/go-10x-minis/minis/15-error-wrapping-sentinel-errors/internal/errorwrappingsentinelerrors"
)

func main() {
	fmt.Println("=== Debugging Error Wrapping and Sentinel Errors ===")

	// Fixed default values - modify these directly if you want to test different inputs
	userID := 1

	fmt.Printf("Testing with user ID: %d\n\n", userID)

	// Set breakpoint in errorwrappingsentinelerrors.go at FindUser function
	username, err := errorwrappingsentinelerrors.FindUser(userID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Is NotFoundError: %v\n", errorwrappingsentinelerrors.IsNotFoundError(err))
	} else {
		fmt.Printf("Username: %s\n", username)
	}

	// Test error wrapping
	config, err := errorwrappingsentinelerrors.ReadConfig(userID)
	if err != nil {
		fmt.Printf("ReadConfig error: %v\n", err)
	}
	if config != "" {
		fmt.Printf("Config: %s\n", config)
	}
}

