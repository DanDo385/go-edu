package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/13-trace/internal/trace"
)

// traceAdapter wraps ethclient to satisfy trace.TraceClient.
// Students will implement the real trace logic in the exercise.
type traceAdapter struct{ c *ethclient.Client }

func (a *traceAdapter) TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error) {
	return nil, fmt.Errorf("not yet implemented - complete the exercise")
}

/*
Transaction Tracing Demo

This application demonstrates tracing transaction execution to understand
internal calls, state changes, and debugging contract interactions.

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
	result, err := trace.Run(ctx, &traceAdapter{client}, trace.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Trace Results ===")
	fmt.Printf("Result: %+v\n", result)
}
