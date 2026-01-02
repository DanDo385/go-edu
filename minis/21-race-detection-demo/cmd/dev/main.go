package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/21-race-detection-demo/internal/racedetectiondemo"
)

func main() {
	fmt.Println("Dev Harness: 21-race-detection-demo")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// racedetectiondemo.Run("dev-input-value")
}
