package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"geth/19-devnets/internal/devnets"
)

/*
Development Networks Demo

This application demonstrates setting up and working with local development
networks (devnets) for testing Ethereum applications.

Usage:

	go run ./cmd/app/main.go [RPC_URL]

Examples:

	go run ./cmd/app/main.go
	go run ./cmd/app/main.go http://localhost:8545
*/
func main() {
	rpcURL := "http://localhost:8545"
	if len(os.Args) > 1 {
		rpcURL = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// BREAKPOINT: Step into Run()
	result, err := devnets.Run(ctx, devnets.Config{RPCURL: rpcURL})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Devnets Results ===")
	fmt.Printf("Result: %+v\n", result)
}
