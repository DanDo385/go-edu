package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"geth/01-stack/internal/stack"
)

/*
Debug Harness for Ethereum Stack Module

This file provides a deterministic debug environment with fixed inputs.
Use this for stepping through the code with VS Code debugger.

How to use:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5 in VS Code
  3. Select "Debug: cmd/dev (Debug Harness)"
  4. Step through with F10 (Step Over) and F11 (Step Into)
*/
func main() {
	fmt.Println("=== Stack Debug Harness ===")
	fmt.Println()

	// Fixed test configuration
	// BREAKPOINT: Inspect these values to understand the test setup
	rpcURL := "https://eth.llamarpc.com"
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Println()

	// Connect to Ethereum
	// BREAKPOINT: Step into DialContext to see connection establishment
	fmt.Println("Step 1: Connecting to Ethereum node...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect: %v\n", err)
		fmt.Println()
		fmt.Println("Troubleshooting:")
		fmt.Println("  - Check your internet connection")
		fmt.Println("  - Try a different RPC URL")
		return
	}
	defer client.Close()

	fmt.Println("         Connected successfully!")
	fmt.Println()

	// Query stack information
	// BREAKPOINT: This is the main entry point - step into Run()
	fmt.Println("Step 2: Querying stack information...")

	result, err := stack.Run(ctx, client, stack.Config{})
	if err != nil {
		fmt.Printf("ERROR: Query failed: %v\n", err)
		return
	}

	// Display results
	// BREAKPOINT: Inspect the result struct
	fmt.Println("Step 3: Results")
	fmt.Println()
	fmt.Println("=== Ethereum Stack Info ===")
	fmt.Printf("Chain ID:    %s\n", result.ChainID.String())
	fmt.Printf("Network ID:  %s\n", result.NetworkID.String())
	fmt.Println()
	fmt.Println("=== Latest Block Header ===")
	fmt.Printf("Number:      %d\n", result.Header.Number.Uint64())
	fmt.Printf("Hash:        %s\n", result.Header.Hash().Hex())
	fmt.Printf("Timestamp:   %s\n", time.Unix(int64(result.Header.Time), 0).UTC())
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Chain ID is used for replay protection (EIP-155)")
	fmt.Println("2. Network ID is a legacy identifier for P2P networking")
	fmt.Println("3. Block headers contain cryptographic commitments to state")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/02-rpc-basics")
}
