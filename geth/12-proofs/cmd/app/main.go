package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"

	"github.com/example/go-10x-minis/geth/12-proofs/internal/proofs"
)

/*
Merkle Proofs Demo

This application demonstrates understanding and verifying Merkle proofs
for accounts and storage in Ethereum.

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

	// gethclient wraps the underlying RPC connection and provides GetProof,
	// which satisfies the proofs.ProofClient interface.
	gc := gethclient.New(client.Client())

	// BREAKPOINT: Step into Run()
	result, err := proofs.Run(ctx, gc, proofs.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Proofs Results ===")
	fmt.Printf("Result: %+v\n", result)
}
