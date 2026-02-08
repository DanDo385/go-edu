package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/16-concurrency/internal/concurrency"
)

// ethProber wraps ethclient to satisfy concurrency.Prober.
// Students will implement the real probe logic in the exercise.
type ethProber struct{ c *ethclient.Client }

func (p *ethProber) Probe(ctx context.Context, endpoint string) error {
	return fmt.Errorf("not yet implemented - complete the exercise")
}

/*
Concurrent RPC Calls Demo

This application demonstrates making concurrent RPC calls to improve
performance when querying multiple pieces of data from Ethereum.

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
	result, err := concurrency.Run(ctx, &ethProber{client}, concurrency.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Concurrency Results ===")
	fmt.Printf("Result: %+v\n", result)
}
