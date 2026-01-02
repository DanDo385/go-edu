package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/01-hello-strings/internal/hellostrings"
)

func main() {
	fmt.Println("Dev Harness: 01-hello-strings")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// hellostrings.Run("dev-input-value")
}
