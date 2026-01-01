package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"geth/07-eth-call/internal/ethcall"
)

/*
eth_call Demo - Manual Contract Interaction

This application demonstrates querying ERC20 token metadata using manual
ABI encoding/decoding. This is the Go equivalent of what you learned in
the geth/06-smart-contracts console tutorial.

Usage:

	go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>

Examples:

	# Query USDC metadata
	go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

	# Query DAI metadata
	go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F
*/
func main() {
	// BREAKPOINT: Set a breakpoint here
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  # Query USDC")
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
	fmt.Printf("Querying contract: %s\n", contract.Hex())
	fmt.Println()

	// BREAKPOINT: Step into Run() to see manual ABI encoding/decoding
	result, err := ethcall.Run(ctx, client, ethcall.Config{
		Contract: contract,
	})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== ERC20 Token Metadata ===")
	fmt.Printf("Name:         %s\n", result.Name)
	fmt.Printf("Symbol:       %s\n", result.Symbol)
	fmt.Printf("Decimals:     %d\n", result.Decimals)
	fmt.Printf("Total Supply: %s\n", result.TotalSupply.String())
}
