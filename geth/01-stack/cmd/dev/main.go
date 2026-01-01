package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"geth/01-stack/internal/stack"
)

/*
Debug harness for geth/01-stack.

This uses fixed, deterministic inputs for easy debugging:
- Test RPC URL: https://eth.llamarpc.com
- Block number: nil (latest)

BREAKPOINT: Set breakpoints throughout this file to step through the RPC connection flow.
*/

func main() {
	fmt.Println("geth/01-stack: Debug Harness")
	fmt.Println("==============================")
	fmt.Println()

	// Fixed test inputs
	rpcURL := "https://eth.llamarpc.com"
	blockNumber := (*big.Int)(nil) // nil means "latest"

	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Block: latest\n")
	fmt.Println()

	// BREAKPOINT: Set here before RPC connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)

	// BREAKPOINT: Set here after client creation
	fmt.Println("Connected successfully!")
	fmt.Println()

	// Call stack.Run with test config
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
	fmt.Println("=== Result ===")
	fmt.Printf("Chain ID:   %s\n", result.ChainID.String())
	fmt.Printf("Network ID: %s\n", result.NetworkID.String())
	fmt.Printf("Block #:    %s\n", result.Header.Number.String())
	fmt.Printf("Block Hash: %s\n", result.Header.Hash().Hex())
}
