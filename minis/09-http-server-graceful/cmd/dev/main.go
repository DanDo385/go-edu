package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/09-http-server-graceful/internal/httpservergraceful"
)

func main() {
	fmt.Println("Dev Harness: 09-http-server-graceful")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// httpservergraceful.Run("dev-input-value")
}
