package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/08-abigen/internal/abigen"
)

/*
abigen Demo - Typed Contract Bindings

This application demonstrates using abigen-generated typed bindings to interact
with smart contracts. This provides type safety and a much cleaner API compared
to manual ABI encoding (geth/07-eth-call).

Usage:

	go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>

Examples:

	go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
*/
func main() {
	// BREAKPOINT: Set a breakpoint here
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	contractAddr := os.Args[2]

	fmt.Printf("Connecting to %s...\n", rpcURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	contract := common.HexToAddress(contractAddr)
	fmt.Printf("Contract: %s\n", contract.Hex())
	fmt.Println()

	// BREAKPOINT: Step into Run() to see typed bindings
	result, err := abigen.Run(ctx, client, abigen.Config{
		Contract: contract,
	})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Contract Data (via abigen) ===")
	fmt.Printf("Result: %+v\n", result)
}
