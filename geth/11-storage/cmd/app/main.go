package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"geth/11-storage/internal/storage"
)

/*
Contract Storage Demo

This application demonstrates reading raw contract storage slots directly.
Useful for understanding how Solidity stores state variables.

Usage:

	go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>

Examples:

	go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
*/
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	contractAddr := os.Args[2]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	contract := common.HexToAddress(contractAddr)

	// BREAKPOINT: Step into Run()
	result, err := storage.Run(ctx, client, storage.Config{Contract: contract})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Storage Results ===")
	fmt.Printf("Result: %+v\n", result)
}
