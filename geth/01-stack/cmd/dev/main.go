package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/01-stack/internal/stack"
)

/*
geth/01-stack: cmd/dev (Debug Harness)

This is a debug harness with fixed, deterministic inputs for stepping through code.

Why use this instead of cmd/app?
  - No command-line arguments to remember
  - Same inputs every time (deterministic)
  - Easier to set breakpoints and step through
  - Perfect for understanding the code flow

Usage:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5 in VS Code
  3. Select "Debug: cmd/dev (Debug Harness)"
  4. Step through with F10 (Step Over) and F11 (Step Into)

BREAKPOINT: Set your first breakpoint at the start of main() below.
*/

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  geth/01-stack: Debug Harness")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// ========================================================================
	// Fixed Test Inputs (No CLI Arguments Needed)
	// ========================================================================
	// BREAKPOINT: Set breakpoint here to see the fixed test inputs
	const rpcURL = "https://eth.llamarpc.com"
	
	// Test cases: nil (latest) and specific block
	testCases := []struct {
		name        string
		blockNumber *big.Int
	}{
		{
			name:        "Latest Block",
			blockNumber: nil,
		},
		{
			name:        "Specific Block (12345)",
			blockNumber: big.NewInt(12345),
		},
	}

	// ========================================================================
	// Connect to Ethereum Node
	// ========================================================================
	// BREAKPOINT: Set breakpoint here to watch RPC connection
	fmt.Printf("Connecting to: %s\n", rpcURL)
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("❌ Error connecting: %v\n", err)
		fmt.Println("\nTip: Check your internet connection and try a different RPC URL")
		return
	}
	defer client.Close()

	fmt.Println("✓ Connected successfully")
	fmt.Println()

	// ========================================================================
	// Run Test Cases
	// ========================================================================
	for i, tc := range testCases {
		fmt.Printf("─────────────────────────────────────────────────────────────\n")
		fmt.Printf("Test Case %d: %s\n", i+1, tc.name)
		fmt.Printf("─────────────────────────────────────────────────────────────\n")

		// BREAKPOINT: Set breakpoint here to step into Run function
		cfg := stack.Config{
			BlockNumber: tc.blockNumber,
		}

		result, err := stack.Run(ctx, client, cfg)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			fmt.Println()
			continue
		}

		// BREAKPOINT: Set breakpoint here to inspect result
		displayResult(result)
		fmt.Println()
	}

	// ========================================================================
	// Debugging Tips
	// ========================================================================
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  Debugging Tips:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("1. Use F10 (Step Over) to execute line by line")
	fmt.Println("2. Use F11 (Step Into) to enter function calls")
	fmt.Println("3. Watch the Variables panel to see data change")
	fmt.Println("4. Hover over variables to see their values")
	fmt.Println("5. Use the Call Stack panel to see function call hierarchy")
	fmt.Println()
	fmt.Println("Next: Open internal/stack/exercise.go and step through Run()")
}

// displayResult formats and prints the stack information
func displayResult(result *stack.Result) {
	fmt.Println()
	fmt.Println("Results:")
	fmt.Printf("  Chain ID:    %s\n", result.ChainID.String())
	fmt.Printf("  Network ID:  %s\n", result.NetworkID.String())
	fmt.Println()
	fmt.Println("  Block Header:")
	fmt.Printf("    Number:      %s\n", result.Header.Number.String())
	fmt.Printf("    Hash:        %s\n", result.Header.Hash().Hex())
	fmt.Printf("    Parent Hash: %s\n", result.Header.ParentHash.Hex())
	fmt.Printf("    State Root:  %s\n", result.Header.Root.Hex())
	fmt.Printf("    Timestamp:   %d (%s)\n", 
		result.Header.Time,
		time.Unix(int64(result.Header.Time), 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("    Gas Used:    %d / %d (%.1f%%)\n",
		result.Header.GasUsed,
		result.Header.GasLimit,
		float64(result.Header.GasUsed)/float64(result.Header.GasLimit)*100)
	
	if result.Header.BaseFee != nil {
		fmt.Printf("    Base Fee:    %s wei\n", result.Header.BaseFee.String())
	}
}
