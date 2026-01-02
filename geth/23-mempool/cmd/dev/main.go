package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/23-mempool/internal/mempool"
)

/*
Debug Harness for Mempool Module
*/
func main() {
	fmt.Println("=== Mempool Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer client.Close()

	// BREAKPOINT: Step into Run()
	result, err := mempool.Run(ctx, client, mempool.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("Next: Proceed to geth/24-monitor")
}
