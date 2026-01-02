package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/10-filters/internal/filters"
)

/*
Debug Harness for Filters Module
*/
func main() {
	fmt.Println("=== Filters Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Println()

	// BREAKPOINT: Connect
	fmt.Println("Step 1: Connecting...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("         Connected!")
	fmt.Println()

	// BREAKPOINT: Step into Run()
	fmt.Println("Step 2: Setting up filters...")

	result, err := filters.Run(ctx, client, filters.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("Step 3: Results")
	fmt.Printf("Result: %+v\n", result)
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. FilterQuery specifies which logs to retrieve")
	fmt.Println("2. FromBlock/ToBlock define the block range")
	fmt.Println("3. Topics filter by event signature and indexed params")
	fmt.Println("4. Addresses filter by emitting contract")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/11-storage")
}
