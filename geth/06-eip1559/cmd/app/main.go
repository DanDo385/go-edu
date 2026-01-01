package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"geth/06-eip1559/internal/eip1559"
)

/*
EIP-1559 Fee Market Demo

This application demonstrates understanding EIP-1559 transaction fee mechanics
including base fee, priority fee (tip), and max fee calculations.

Usage:

	go run ./cmd/app/main.go <RPC_URL>

Examples:

	go run ./cmd/app/main.go https://eth.llamarpc.com
*/
func main() {
	// BREAKPOINT: Set a breakpoint here
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

	// BREAKPOINT: Step into Run() to see EIP-1559 calculations
	result, err := eip1559.Run(ctx, client, eip1559.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== EIP-1559 Fee Information ===")
	fmt.Printf("Base Fee:      %s gwei\n", result.BaseFee.String())
	fmt.Printf("Max Fee:       %s gwei\n", result.MaxFee.String())
	fmt.Printf("Priority Fee:  %s gwei\n", result.PriorityFee.String())
}
