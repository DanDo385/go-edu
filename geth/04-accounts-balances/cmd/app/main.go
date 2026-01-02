package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/04-accounts-balances/internal/accountsbalances"
)

/*
Accounts and Balances Demo

This application demonstrates querying Ethereum account balances, nonces,
and understanding account state.

Usage:

	go run ./cmd/app/main.go <RPC_URL> <ADDRESS>

Examples:

	# Query Vitalik's address
	go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
*/
func main() {
	// BREAKPOINT: Set a breakpoint here to inspect arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL> <ADDRESS>")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	addressHex := os.Args[2]

	fmt.Printf("Connecting to %s...\n", rpcURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	address := common.HexToAddress(addressHex)
	fmt.Printf("Querying account: %s\n", address.Hex())
	fmt.Println()

	// BREAKPOINT: Step into Run() to see balance queries
	result, err := accountsbalances.Run(ctx, client, accountsbalances.Config{
		Address: address,
	})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Account State ===")
	fmt.Printf("Address: %s\n", address.Hex())
	fmt.Printf("Balance: %s wei\n", result.Balance.String())
	fmt.Printf("Nonce:   %d\n", result.Nonce)
}
