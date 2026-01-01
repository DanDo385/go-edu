package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"geth/25-toolbox/internal/toolbox"
)

/*
Debug Harness for Toolbox Module

This is the final module - combining all skills learned throughout the course.
*/
func main() {
	fmt.Println("=== Toolbox Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer client.Close()

	// BREAKPOINT: Step into Run()
	result, err := toolbox.Run(ctx, client, toolbox.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("=== Congratulations! ===")
	fmt.Println("You've completed all geth modules!")
	fmt.Println()
	fmt.Println("Skills learned:")
	fmt.Println("- Ethereum connectivity and RPC basics")
	fmt.Println("- Key management and addresses")
	fmt.Println("- Transaction nonces and EIP-1559")
	fmt.Println("- Smart contract interaction (console and Go)")
	fmt.Println("- Events, filters, and logs")
	fmt.Println("- Storage, proofs, and tracing")
	fmt.Println("- Indexing and chain monitoring")
	fmt.Println("- Node operations and networking")
}
