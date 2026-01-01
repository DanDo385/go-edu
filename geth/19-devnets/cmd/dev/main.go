package main

import (
	"context"
	"fmt"
	"time"

	"geth/19-devnets/internal/devnets"
)

/*
Debug Harness for Devnets Module
*/
func main() {
	fmt.Println("=== Devnets Debug Harness ===")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// BREAKPOINT: Step into Run()
	result, err := devnets.Run(ctx, devnets.Config{
		RPCURL: "http://localhost:8545",
	})
	if err != nil {
		fmt.Printf("Note: %v\n", err)
		fmt.Println("(This is expected if no local devnet is running)")
		fmt.Println()
		fmt.Println("To start a devnet: geth --dev --http --http.api eth,net,web3,personal")
		return
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("Next: Proceed to geth/20-node")
}
