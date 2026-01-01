package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"geth/02-rpc-basics/internal/rpcbasics"
)

/*
RPC Basics Demo

This application demonstrates fundamental Ethereum RPC calls including
fetching blocks, transactions, and understanding the RPC layer.

Usage:

	go run ./cmd/app/main.go <RPC_URL>

Examples:

	go run ./cmd/app/main.go https://eth.llamarpc.com
*/
func main() {
	// BREAKPOINT: Set a breakpoint here to inspect command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL>")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  go run ./cmd/app/main.go https://eth.llamarpc.com")
		os.Exit(1)
	}

	rpcURL := os.Args[1]

	fmt.Printf("Connecting to %s...\n", rpcURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("Connected!")
	fmt.Println()

	// BREAKPOINT: Step into Run() to see RPC call mechanics
	result, err := rpcbasics.Run(ctx, client, rpcbasics.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== RPC Basics Results ===")
	fmt.Printf("Result: %+v\n", result)
}
