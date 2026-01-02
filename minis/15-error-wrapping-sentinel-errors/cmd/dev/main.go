package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/15-error-wrapping-sentinel-errors/internal/errorwrappingsentinelerrors"
)

func main() {
	fmt.Println("Dev Harness: 15-error-wrapping-sentinel-errors")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// errorwrappingsentinelerrors.Run("dev-input-value")
}
