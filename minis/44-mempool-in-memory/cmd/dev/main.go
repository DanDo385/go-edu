package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/44-mempool-in-memory/internal/mempoolinmemory"
)

func main() {
	fmt.Println("Dev Harness: 44-mempool-in-memory")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// mempoolinmemory.Run("dev-input-value")
}
