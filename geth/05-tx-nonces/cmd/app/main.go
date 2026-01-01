package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"geth/05-tx-nonces/internal/txnonces"
)

/*
Transaction Nonces Demo

This application demonstrates understanding transaction nonces, their role
in transaction ordering, and querying pending vs confirmed nonces.

Usage:

	go run ./cmd/app/main.go <RPC_URL> <ADDRESS>

Examples:

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
	fmt.Printf("Querying nonce for: %s\n", address.Hex())
	fmt.Println()

	// BREAKPOINT: Step into Run()
	result, err := txnonces.Run(ctx, client, txnonces.Config{
		Address: address,
	})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Nonce Information ===")
	fmt.Printf("Address:       %s\n", address.Hex())
	fmt.Printf("Nonce:         %d\n", result.Nonce)
	fmt.Printf("Pending Nonce: %d\n", result.PendingNonce)
}
