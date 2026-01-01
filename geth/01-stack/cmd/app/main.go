package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"geth/01-stack/internal/stack"
)

/*
geth/01-stack: Prove RPC connectivity by reading network identifiers and latest header.

This is the foundational module that teaches:
- Connecting to an Ethereum RPC endpoint
- Retrieving chain ID and network ID
- Fetching block headers
- Understanding the Ethereum stack

Usage:
  go run ./cmd/app/main.go <RPC_URL> [block_number]

Examples:
  # Get latest block info
  go run ./cmd/app/main.go https://eth.llamarpc.com

  # Get specific block info
  go run ./cmd/app/main.go https://eth.llamarpc.com 12345

BREAKPOINT: Set breakpoints throughout this file to understand RPC connection flow.
*/

func main() {
	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [block_number]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 12345\n", os.Args[0])
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	var blockNumber *big.Int

	// Parse optional block number
	if len(os.Args) >= 3 {
		blockNum, ok := new(big.Int).SetString(os.Args[2], 10)
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid block number: %s\n", os.Args[2])
			os.Exit(1)
		}
		blockNumber = blockNum
	}

	// BREAKPOINT: Set here to inspect parsed arguments
	fmt.Printf("Connecting to RPC endpoint: %s\n", rpcURL)
	if blockNumber != nil {
		fmt.Printf("Block number: %s\n", blockNumber.String())
	} else {
		fmt.Println("Block number: latest")
	}
	fmt.Println()

	// Create RPC client with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set here before RPC connection
	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to RPC endpoint: %v\n", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	// Create ethclient wrapper
	client := ethclient.NewClient(rpcClient)

	// BREAKPOINT: Set here after client creation
	fmt.Println("Connected successfully!")
	fmt.Println()

	// Call the stack.Run function
	cfg := stack.Config{
		BlockNumber: blockNumber,
	}

	// BREAKPOINT: Set here before calling stack.Run
	result, err := stack.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// BREAKPOINT: Set here to inspect result
	fmt.Println("=== Ethereum Stack Information ===")
	fmt.Printf("Chain ID:   %s\n", result.ChainID.String())
	fmt.Printf("Network ID: %s\n", result.NetworkID.String())
	fmt.Printf("Block #:    %s\n", result.Header.Number.String())
	fmt.Printf("Block Hash: %s\n", result.Header.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", result.Header.ParentHash.Hex())
	fmt.Printf("Gas Used:   %d\n", result.Header.GasUsed)
	fmt.Printf("Gas Limit:  %d\n", result.Header.GasLimit)
	fmt.Printf("Timestamp:  %s\n", time.Unix(int64(result.Header.Time), 0).UTC().Format(time.RFC3339))
}
