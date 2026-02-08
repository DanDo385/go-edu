package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/example/go-10x-minis/geth/20-node/internal/node"
)

/*
Node Configuration Demo

This application demonstrates understanding and configuring Ethereum node
options, sync modes, and operational parameters.

Usage:

	go run ./cmd/app/main.go [OPTIONS]

Examples:

	go run ./cmd/app/main.go
*/
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// BREAKPOINT: Step into Run()
	result, err := node.Run(ctx, node.Config{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Node Results ===")
	fmt.Printf("Result: %+v\n", result)
}
