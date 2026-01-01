package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

/*
geth/09-events: cmd/dev (Debug Harness)

Fixed test inputs for debugging. No CLI arguments needed.

Usage:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5, select "Debug: cmd/dev (Debug Harness)"
  3. Step through with F10/F11

BREAKPOINT: Start here
*/

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  geth/09-events: Debug Harness")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// BREAKPOINT: Fixed test inputs
	const rpcURL = "https://eth.llamarpc.com"
	
	fmt.Printf("Connecting to: %s\n", rpcURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("✓ Connected")
	fmt.Println()

	// BREAKPOINT: Exercise with fixed inputs
	fmt.Println("Running with test inputs...")
	// TODO: Add exercise-specific test logic
	// See internal/*/exercise.go for implementation
	
	fmt.Println()
	fmt.Println("✓ Complete")
	fmt.Println()
	fmt.Println("Tips: F10=Step Over, F11=Step Into, Watch Variables panel")
}
