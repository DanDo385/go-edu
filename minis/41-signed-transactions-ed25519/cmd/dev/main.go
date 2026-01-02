package main

import (
	"fmt"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/41-signed-transactions-ed25519/internal/signedtransactionsed25519"
)

func main() {
	fmt.Println("Dev Harness: 41-signed-transactions-ed25519")
	fmt.Println("---------------------------------------------------")
	fmt.Println("This harness runs the code with predefined arguments.")
	
	// TODO: Call your internal package function here with hardcoded args
	// Example:
	// signedtransactionsed25519.Run("dev-input-value")
}
