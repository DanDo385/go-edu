package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/internal/workerpoolwordcount"
)

func main() {
	fmt.Println("Dev Harness: 06-worker-pool-wordcount")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// workerpoolwordcount.Run("dev-input-value")
}
