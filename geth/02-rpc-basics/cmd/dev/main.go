package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/02-rpc-basics/internal/rpcbasics"
)

/*
Debug Harness for RPC Basics Module

Fixed inputs for deterministic debugging sessions.
*/
func main() {
	fmt.Println("=== RPC Basics Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Println()

	// BREAKPOINT: Connect to Ethereum
	fmt.Println("Step 1: Connecting...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("         Connected!")
	fmt.Println()

	// BREAKPOINT: Step into Run()
	fmt.Println("Step 2: Running RPC basics...")

	result, err := rpcbasics.Run(ctx, client, rpcbasics.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("Step 3: Results")
	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("Next: Proceed to geth/03-keys-addresses")
}
