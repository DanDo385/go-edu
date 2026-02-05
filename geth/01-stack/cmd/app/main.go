package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/01-stack/internal/stack"
)

/*
Ethereum Stack Connectivity Demo

This application demonstrates basic Ethereum connectivity by querying chain ID,
network ID, and the latest block header. This is the first thing any Ethereum
application should do - verify it can talk to the network.

Usage:

	go run ./cmd/app/main.go <RPC_URL> [BLOCK_NUMBER]

Examples:

	# Query latest block on mainnet
	go run ./cmd/app/main.go https://eth.llamarpc.com

	# Query specific block number
	go run ./cmd/app/main.go https://eth.llamarpc.com 12345678

Arguments:

	RPC_URL      - Ethereum RPC endpoint (required)
	BLOCK_NUMBER - Specific block number to query (optional, default: latest)
*/
func main() {
	// BREAKPOINT: Set a breakpoint here to inspect command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL> [BLOCK_NUMBER]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run ./cmd/app/main.go https://eth.llamarpc.com")
		fmt.Println("  go run ./cmd/app/main.go https://eth.llamarpc.com 12345678")
		os.Exit(1)
	}

	rpcURL := os.Args[1]

	// Connect to Ethereum node
	// BREAKPOINT: Step into DialContext to see connection establishment
	fmt.Printf("Connecting to %s...\n", rpcURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // background context with a timeout of 30 seconds
	defer cancel() // cancel the context if the function returns

	client, err := ethclient.DialContext(ctx, rpcURL) // dial the Ethereum node
	if err != nil { // if the connection fails, print the error and exit the program
		fmt.Printf("Failed to connect: %v\n", err) // print the error
		os.Exit(1) // exit the program with a status code of 1
	}
	defer client.Close() // close the connection when the function returns

	fmt.Println("Connected!")
	fmt.Println()

	// BREAKPOINT: Step into Run() to see the core logic
	result, err := stack.Run(ctx, client, stack.Config{})
	if err != nil {
		fmt.Printf("Failed to query stack: %v\n", err)
		os.Exit(1)
	}

	// Display results
	fmt.Println("=== Ethereum Stack Info ===")
	fmt.Printf("Chain ID:    %s\n", result.ChainID.String())
	fmt.Printf("Network ID:  %s\n", result.NetworkID.String())
	fmt.Println()
	fmt.Println("=== Latest Block Header ===")
	fmt.Printf("Number:      %d\n", result.Header.Number.Uint64())
	fmt.Printf("Hash:        %s\n", result.Header.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", result.Header.ParentHash.Hex())
	fmt.Printf("Timestamp:   %d (%s)\n", result.Header.Time, time.Unix(int64(result.Header.Time), 0).UTC())
	fmt.Printf("Gas Limit:   %d\n", result.Header.GasLimit)
	fmt.Printf("Gas Used:    %d\n", result.Header.GasUsed)
}
