package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/04-jsonl-log-filter/internal/jsonllogfilter"
)

func main() {
	fmt.Println("Dev Harness: 04-jsonl-log-filter")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// jsonllogfilter.Run("dev-input-value")
}
