package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/21-sync/internal/sync"
)

/*
Sync Status Demo

This application demonstrates monitoring node synchronization status
and understanding different sync modes.

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
	result, err := sync.Run(ctx, client, sync.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Sync Results ===")
	fmt.Printf("Result: %+v\n", result)
}
