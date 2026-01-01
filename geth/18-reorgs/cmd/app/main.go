package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"geth/18-reorgs/internal/reorgs"
)

/*
Chain Reorganization Demo

This application demonstrates understanding and handling chain reorganizations
(reorgs) in Ethereum applications.

Usage:

	go run ./cmd/app/main.go <RPC_URL>

Examples:

	go run ./cmd/app/main.go https://eth.llamarpc.com
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL>")
		os.Exit(1)
	}

	rpcURL := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// BREAKPOINT: Step into Run()
	result, err := reorgs.Run(ctx, client, reorgs.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Reorgs Results ===")
	fmt.Printf("Result: %+v\n", result)
}
