package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"geth/02-rpc-basics/internal/rpcbasics"
)

/*
geth/02-rpc-basics: RPC connection basics with retry logic.

This module teaches:
- RPC connection handling
- Retry logic for unreliable connections
- Fetching block information

Usage:
  go run ./cmd/app/main.go <RPC_URL> [retries]

Examples:
  go run ./cmd/app/main.go https://eth.llamarpc.com
  go run ./cmd/app/main.go https://eth.llamarpc.com 3
*/

func main() {
	var retries int
	flag.IntVar(&retries, "retries", 3, "Number of retry attempts")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [-retries=N]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com -retries=5\n", os.Args[0])
		os.Exit(1)
	}

	rpcURL := args[0]

	fmt.Printf("Connecting to: %s\n", rpcURL)
	fmt.Printf("Retries: %d\n", retries)
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)

	cfg := rpcbasics.Config{
		Retries: retries,
	}

	result, err := rpcbasics.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== RPC Basics Result ===")
	fmt.Printf("Network ID:   %s\n", result.NetworkID.String())
	fmt.Printf("Block Number: %d\n", result.BlockNumber)
	fmt.Printf("Block Hash:   %s\n", result.Block.Hash().Hex())
	fmt.Printf("Gas Used:     %d\n", result.Block.GasUsed())
}
