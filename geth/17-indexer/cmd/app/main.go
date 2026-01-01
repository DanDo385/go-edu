package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

/*
geth/17-indexer: cmd/app

Usage:
  go run ./cmd/app/main.go <RPC_URL> [args...]

Example:
  go run ./cmd/app/main.go https://eth.llamarpc.com

BREAKPOINT: Set breakpoints at "// BREAKPOINT:" comments for debugging.
*/

func main() {
	// BREAKPOINT: Inspect command-line arguments
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s https://eth.llamarpc.com\n", os.Args[0])
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	fmt.Printf("Connecting to: %s\n", rpcURL)

	// BREAKPOINT: Watch RPC connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("✓ Connected")
	fmt.Println()

	// BREAKPOINT: Exercise execution point
	fmt.Println("Running exercise...")
	// TODO: Add exercise-specific logic here
	// See internal/*/exercise.go for implementation
	
	fmt.Println("✓ Complete")
}
