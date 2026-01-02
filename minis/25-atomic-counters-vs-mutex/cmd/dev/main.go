package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/25-atomic-counters-vs-mutex/internal/atomiccountersvsmutex"
)

func main() {
	fmt.Println("Dev Harness: 25-atomic-counters-vs-mutex")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// atomiccountersvsmutex.Run("dev-input-value")
}
