package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"geth/09-events/internal/events"
)

/*
Debug Harness for Events Module
*/
func main() {
	fmt.Println("=== Events Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	contractAddr := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

	fmt.Printf("RPC URL:  %s\n", rpcURL)
	fmt.Printf("Contract: %s (USDC)\n", contractAddr.Hex())
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
	fmt.Println("Step 2: Querying events...")

	result, err := events.Run(ctx, client, events.Config{
		Contract: contractAddr,
	})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("Step 3: Results")
	fmt.Printf("Result: %+v\n", result)
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Events are emitted by contracts to log state changes")
	fmt.Println("2. Topics[0] is the event signature hash")
	fmt.Println("3. Topics[1-3] are indexed parameters")
	fmt.Println("4. Data contains non-indexed parameters")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/10-filters")
}
