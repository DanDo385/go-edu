package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/02-rpc-basics/internal/rpcbasics"
)

/*
===================================================================
Ethereum RPC Basics - Full Block Retrieval with Retry Logic
===================================================================

This program demonstrates how to fetch full Ethereum blocks (with all
transaction data) using a resilient RPC client with retry logic for
handling transient network failures.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> [demo]

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, latest, historical, compare, retry (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com latest
    go run ./cmd/app/main.go https://eth.llamarpc.com historical
    go run ./cmd/app/main.go https://eth.llamarpc.com compare
    go run ./cmd/app/main.go https://eth.llamarpc.com retry

Common RPC Endpoints:
    Ethereum Mainnet:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth
      - https://eth-mainnet.public.blastapi.io

    Sepolia Testnet:
      - https://rpc.sepolia.org
      - https://sepolia.gateway.tenderly.co

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter rpcbasics package functions
- Watch Variables panel to see RPC responses
- Inspect Block structure with Header and Transactions

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the rpcbasics package from internal/rpcbasics
- exercise.go, solution.reference.go are in package rpcbasics
- The rpcbasics.Run() function encapsulates all RPC logic with retry

KEY CONCEPTS:
- Full Blocks vs Headers: Blocks contain complete transaction data
- Fault Tolerance: Retry logic handles transient network failures
- Context-Aware Operations: Respects cancellation and timeouts
- Defensive Copying: Prevents mutations of shared RPC data
- Transaction Composition: A Block contains a Header and Transactions array
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the RPC URL (required)
	// DEBUG: os.Args[2] is the demo type (optional)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com latest\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com historical\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, latest, historical, compare, retry\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum RPC Basics - Full Block Retrieval ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection
	// DEBUG: This is where we establish the RPC connection
	// DEBUG: Watch for any connection errors in the error variable
	//
	// Context with timeout: We use context.WithTimeout to prevent hanging
	// indefinitely if the RPC endpoint is slow or unreachable. This is a
	// critical pattern for production code—always set timeouts on network calls.
	//
	// ethclient.Dial: This is the standard go-ethereum client that implements
	// the RPCClient interface. It handles JSON-RPC communication with Ethereum nodes.

	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Step Over (F10) to execute the Dial call
	// DEBUG: If it fails, inspect 'err' to see the connection error
	client, err := ethclient.Dial(dialCtx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC endpoint: %v\n", err)
	}
	defer client.Close()

	// BREAKPOINT: Set breakpoint here after successful connection
	// DEBUG: Connection is established—client is ready to use
	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 3: Example 1 - Retrieve Latest Full Block
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates fetching the latest full block with all transactions
	// DEBUG: Watch the 'result' variable to see NetworkID, BlockNumber, and Block

	if demo == "all" || demo == "latest" {
		fmt.Println("=== Example 1: Latest Full Block ===")
		fmt.Println("Retrieving latest full block with all transactions...")
		fmt.Println()

		// Create configuration with retry logic
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg.Retries is set to 3—this enables fault tolerance
		// DEBUG: If the first RPC call fails, it will retry up to 3 times
		cfg := rpcbasics.Config{
			Retries: 3,
		}

		// BREAKPOINT: Set breakpoint here before calling rpcbasics.Run
		// DEBUG: Step Into (F11) to see the implementation in exercise.go or solution.reference.go
		// DEBUG: You'll see how the function makes RPC calls with retry logic
		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := rpcbasics.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve latest block: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful retrieval
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at result.Block—it contains both Header and Transactions
		// DEBUG: Notice result.Block.Transactions() returns all tx data

		fmt.Println("Network Information:")
		fmt.Printf("  Network ID:   %s\n", result.NetworkID.String())
		fmt.Printf("  Block Number: %d\n\n", result.BlockNumber)

		fmt.Println("Block Header:")
		header := result.Block.Header()
		fmt.Printf("  Block Number: %s\n", header.Number.String())
		fmt.Printf("  Block Hash:   %s\n", header.Hash().Hex())
		fmt.Printf("  Parent Hash:  %s\n", header.ParentHash.Hex())
		fmt.Printf("  Timestamp:    %d (%s)\n", header.Time, time.Unix(int64(header.Time), 0))
		fmt.Printf("  Gas Limit:    %d\n", header.GasLimit)
		fmt.Printf("  Gas Used:     %d (%.2f%% full)\n",
			header.GasUsed,
			float64(header.GasUsed)/float64(header.GasLimit)*100)
		if header.BaseFee != nil {
			fmt.Printf("  Base Fee:     %s wei\n", header.BaseFee.String())
		}
		fmt.Println()

		// Display transaction information
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This shows the difference between full blocks and headers
		// DEBUG: Full blocks include all transaction data
		txs := result.Block.Transactions()
		fmt.Printf("Transactions: %d total\n", len(txs))
		if len(txs) > 0 {
			fmt.Println("  First 3 transactions:")
			limit := 3
			if len(txs) < limit {
				limit = len(txs)
			}
			for i := 0; i < limit; i++ {
				tx := txs[i]
				fmt.Printf("    [%d] Hash: %s, Value: %s wei\n", i, tx.Hash().Hex(), tx.Value().String())
			}
		}
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 2 - Retrieve Historical Block by Number
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through historical block example
	// DEBUG: This shows how to fetch a specific historical block with all its data

	if demo == "all" || demo == "historical" {
		fmt.Println("=== Example 2: Historical Block (Block 15,000,000) ===")
		fmt.Println("Retrieving historical full block...")
		fmt.Println()

		// Block 15,000,000 was mined on September 15, 2022 - The Merge!
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This is a significant block in Ethereum history
		cfg := rpcbasics.Config{
			Retries: 3,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := rpcbasics.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve historical block: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect the block to see historical transaction data
		// DEBUG: Notice how the block structure is identical to latest blocks

		header := result.Block.Header()
		fmt.Printf("Block #%s Information:\n", header.Number.String())
		fmt.Printf("  Block Hash:        %s\n", header.Hash().Hex())
		fmt.Printf("  Timestamp:         %d (%s)\n", header.Time, time.Unix(int64(header.Time), 0))
		fmt.Printf("  Transaction Count: %d\n", len(result.Block.Transactions()))
		fmt.Printf("  Gas Used:          %d / %d (%.2f%%)\n",
			header.GasUsed, header.GasLimit,
			float64(header.GasUsed)/float64(header.GasLimit)*100)
		fmt.Printf("  Miner:             %s\n", header.Coinbase.Hex())
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 3 - Compare Block Header vs Full Block Data
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand block composition
	// DEBUG: This demonstrates the difference between headers and full blocks

	if demo == "all" || demo == "compare" {
		fmt.Println("=== Example 3: Block Composition - Header vs Full Block ===")
		fmt.Println()

		cfg := rpcbasics.Config{
			Retries: 3,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := rpcbasics.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve block: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice that result.Block contains both Header and Transactions
		// DEBUG: The Header is lightweight metadata, Transactions are the actual data

		fmt.Println("Understanding Block Composition:")
		fmt.Println()
		fmt.Println("1. Block Header (Lightweight - ~500 bytes):")
		fmt.Println("   - Block number, hash, parent hash")
		fmt.Println("   - Timestamps, difficulty, nonce")
		fmt.Println("   - Gas limit, gas used")
		fmt.Println("   - Merkle roots (state, transactions, receipts)")
		fmt.Println("   - Miner address, extra data")
		fmt.Println()

		fmt.Println("2. Full Block (Heavy - can be megabytes):")
		fmt.Println("   - Header (above)")
		fmt.Println("   - Complete transaction array")
		fmt.Println("   - Each transaction includes:")
		fmt.Println("     * Sender/recipient addresses")
		fmt.Println("     * Value, gas price, gas limit")
		fmt.Println("     * Nonce, signature (v, r, s)")
		fmt.Println("     * Input data (for contract calls)")
		fmt.Println()

		txs := result.Block.Transactions()

		fmt.Printf("Current Block Statistics:\n")
		fmt.Printf("  Header Size:       ~500 bytes\n")
		fmt.Printf("  Transaction Count: %d\n", len(txs))
		fmt.Printf("  Estimated Block Size: ~%d KB (with all tx data)\n", len(txs)*200/1024)
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This demonstrates why we sometimes fetch headers only
		// DEBUG: Full blocks are much larger but contain all transaction details
		fmt.Println("Why This Matters:")
		fmt.Println("  - Use Headers for: lightweight sync, block navigation, chain verification")
		fmt.Println("  - Use Full Blocks for: transaction indexing, contract event processing")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Example 4 - Demonstrate Retry Logic and Fault Tolerance
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand retry logic
	// DEBUG: This shows how the retry mechanism provides fault tolerance

	if demo == "all" || demo == "retry" {
		fmt.Println("=== Example 4: Retry Logic and Fault Tolerance ===")
		fmt.Println()

		fmt.Println("Understanding Retry Logic:")
		fmt.Println("  1. RPC calls can fail due to network issues, rate limits, or node problems")
		fmt.Println("  2. Retry logic makes your application more resilient")
		fmt.Println("  3. The retry pattern includes exponential backoff and context awareness")
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: We'll demonstrate retry with different retry counts
		fmt.Println("Demo: Fetching block with different retry configurations")
		fmt.Println()

		// First attempt with 0 retries
		fmt.Println("Attempt 1: No retries (Retries=0)")
		cfg1 := rpcbasics.Config{Retries: 0}
		runCtx1, runCancel1 := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel1()

		result1, err := rpcbasics.Run(runCtx1, client, cfg1)
		if err != nil {
			fmt.Printf("  ✗ Failed: %v\n\n", err)
		} else {
			fmt.Printf("  ✓ Success: Retrieved block %d\n\n", result1.BlockNumber)
		}

		// Second attempt with 3 retries
		fmt.Println("Attempt 2: With retries (Retries=3)")
		cfg2 := rpcbasics.Config{Retries: 3}
		runCtx2, runCancel2 := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel2()

		// BREAKPOINT: Set breakpoint here before retry attempt
		// DEBUG: Step Into (F11) to see retry logic in action
		// DEBUG: Watch how it handles failures and retries with delays
		result2, err := rpcbasics.Run(runCtx2, client, cfg2)
		if err != nil {
			fmt.Printf("  ✗ Failed: %v\n\n", err)
		} else {
			fmt.Printf("  ✓ Success: Retrieved block %d\n", result2.BlockNumber)
			fmt.Printf("  Transaction count: %d\n\n", len(result2.Block.Transactions()))
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Compare the results to understand retry behavior
		fmt.Println("Retry Best Practices:")
		fmt.Println("  1. Always use retries for production applications")
		fmt.Println("  2. Implement exponential backoff to avoid overwhelming servers")
		fmt.Println("  3. Respect context cancellation for responsive shutdown")
		fmt.Println("  4. Log retry attempts for debugging and monitoring")
		fmt.Println("  5. Set maximum retry limits to prevent infinite loops")
		fmt.Println()
	}

	// ============================================================================
	// STEP 7: Understanding Defensive Copying in Full Blocks
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand defensive copying
	// DEBUG: This demonstrates why defensive copying is critical for full blocks

	if demo == "all" {
		fmt.Println("=== Understanding Defensive Copying for Full Blocks ===")
		fmt.Println()

		cfg := rpcbasics.Config{Retries: 3}
		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := rpcbasics.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve data: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice that result.Block is a complete copy
		// DEBUG: Modifying it won't affect the RPC client's internal cache

		fmt.Println("Why Defensive Copying Matters for Full Blocks:")
		fmt.Println("  1. Full blocks contain mutable data structures (big.Int, slices)")
		fmt.Println("  2. Without copying, multiple callers could mutate shared data")
		fmt.Println("  3. types.CopyBlock() performs deep copy of header AND transactions")
		fmt.Println("  4. Each transaction is also deep-copied to prevent mutations")
		fmt.Println("  5. This prevents data races in concurrent applications")
		fmt.Println()

		fmt.Println("Example: Block data is safely copied")
		fmt.Printf("  Block Number: %s\n", result.Block.Number().String())
		fmt.Printf("  Transaction Count: %d\n", len(result.Block.Transactions()))
		fmt.Println("  ✓ This data is our copy—we can safely read/modify it")
		fmt.Println()
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter rpcbasics package functions
	// 5. Watch Variables panel to see RPC responses
	// 6. Inspect result.Block to see full block structure
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Full Blocks vs Headers: When to use each?
	// - Retry Logic: How does it handle transient failures?
	// - Context Awareness: How does it enable cancellation?
	// - Transaction Data: What information is in each transaction?
	// - Block Composition: How are headers and transactions related?
	//
	// STEPPING INTO RPCBASICS PACKAGE:
	// - Set breakpoint at rpcbasics.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see RPC calls: NetworkID(), BlockNumber(), BlockByNumber()
	// - Watch the retry loop handle failures with exponential backoff
	// - See how context cancellation is respected during retry delays
	// - Observe defensive copying with types.CopyBlock()
	//
	// COMMON RPC ERRORS TO DEBUG:
	// - Connection timeout: Check network connectivity and RPC URL
	// - Invalid RPC endpoint: Verify the URL is correct
	// - Rate limiting: Some public endpoints have rate limits (retry helps!)
	// - Block not found: Historical blocks might not be available on all nodes
	// - Insufficient resources: Full blocks are large—ensure adequate memory
}
