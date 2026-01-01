package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/09-events/internal/events"
)

/*
===================================================================
ERC20 Transfer Event Query and Decoding
===================================================================

This program demonstrates how to query and decode ERC20 Transfer events
from blockchain logs. Events provide an append-only audit trail of state
changes and are efficiently searchable via bloom filters.

USAGE:
    go run ./cmd/app/main.go <RPC_URL> <TOKEN_ADDRESS> <FROM_BLOCK> <TO_BLOCK>
    go run ./cmd/app/main.go <RPC_URL> <TOKEN_ADDRESS> <FROM_BLOCK> <TO_BLOCK> <demo>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    <TOKEN_ADDRESS> - ERC20 token contract address (required)
    <FROM_BLOCK>    - Starting block number (required)
    <TO_BLOCK>      - Ending block number (required)
    [demo]          - Demo to run: all, basic, filtered, analysis (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 19000000 19000100
    go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F 18000000 18000050 basic
    go run ./cmd/app/main.go https://sepolia.gateway.tenderly.co 0x... 5000000 5000100 filtered

Common ERC20 Tokens (Mainnet):
    USDC:  0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
    DAI:   0x6B175474E89094C44Da98b954EedeAC495271d0F
    WETH:  0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
    USDT:  0xdAC17F958D2ee523a2206206994597C13D831ec7

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
- Step Into (F11) to enter events package functions
- Watch Variables panel to see event filtering and decoding
- Inspect Topics and Data fields in raw logs

KEY CONCEPTS:
- Events/Logs: Append-only records emitted during contract execution
- Bloom Filters: Probabilistic data structure for efficient log searching
- Topics: Indexed parameters (searchable, max 3 per event)
- Data: Non-indexed parameters (cheaper, not searchable)
- FilterQuery: Specifies which logs to retrieve

ERC20 TRANSFER EVENT:
- Signature: Transfer(address indexed from, address indexed to, uint256 value)
- Topics[0]: Event signature hash (keccak256("Transfer(...)"))
- Topics[1]: from address (indexed)
- Topics[2]: to address (indexed)
- Data: value (uint256, non-indexed)
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments

	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> <TOKEN_ADDRESS> <FROM_BLOCK> <TO_BLOCK> [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 19000000 19000100\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F 18000000 18000050 basic\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, basic, filtered, analysis\n")
		fmt.Fprintf(os.Stderr, "\nNote: Keep block range small (< 10,000 blocks) to avoid rate limits\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	tokenAddressHex := os.Args[2]
	fromBlockStr := os.Args[3]
	toBlockStr := os.Args[4]

	demo := "all"
	if len(os.Args) > 5 {
		demo = os.Args[5]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	fmt.Printf("=== ERC20 Transfer Event Query and Decoding ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Parse Token Contract Address
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before parsing address

	if !common.IsHexAddress(tokenAddressHex) {
		log.Fatalf("Invalid token address: %s\n", tokenAddressHex)
	}
	tokenAddress := common.HexToAddress(tokenAddressHex)

	fmt.Printf("✓ Parsed token address: %s\n", tokenAddress.Hex())

	// ============================================================================
	// STEP 3: Parse Block Range
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before parsing block numbers
	// DEBUG: Block range determines how many logs we'll query
	// DEBUG: Large ranges (>10k blocks) may hit RPC provider limits

	fromBlock, ok := new(big.Int).SetString(fromBlockStr, 10)
	if !ok {
		log.Fatalf("Invalid from block: %s\n", fromBlockStr)
	}

	toBlock, ok := new(big.Int).SetString(toBlockStr, 10)
	if !ok {
		log.Fatalf("Invalid to block: %s\n", toBlockStr)
	}

	// Validate block range
	// BREAKPOINT: Set breakpoint here
	blockRange := new(big.Int).Sub(toBlock, fromBlock).Int64()
	if blockRange < 0 {
		log.Fatalf("Invalid block range: from %s > to %s\n", fromBlockStr, toBlockStr)
	}
	if blockRange > 10000 {
		fmt.Printf("⚠ Warning: Large block range (%d blocks) may hit rate limits\n", blockRange)
	}

	fmt.Printf("✓ Block range: %s to %s (%d blocks)\n", fromBlock.String(), toBlock.String(), blockRange)
	fmt.Println()

	// ============================================================================
	// STEP 4: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection

	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set breakpoint here
	client, err := ethclient.Dial(dialCtx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC endpoint: %v\n", err)
	}
	defer client.Close()

	// BREAKPOINT: Set breakpoint here after successful connection
	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 5: Example 1 - Basic Transfer Event Query
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates querying all Transfer events in block range
	// DEBUG: No filtering by sender/receiver

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Transfer Event Query ===")
		fmt.Println("Querying all Transfer events in block range...")
		fmt.Println()

		// Create configuration for basic query
		// BREAKPOINT: Set breakpoint here
		// DEBUG: No FromHolder/ToHolder means query all transfers
		cfg := events.Config{
			Token:     tokenAddress,
			FromBlock: fromBlock,
			ToBlock:   toBlock,
		}

		// BREAKPOINT: Set breakpoint here before calling events.Run
		// DEBUG: Step Into (F11) to see FilterQuery construction and log decoding
		// DEBUG: You'll see bloom filter usage, topic filtering, and event parsing
		runCtx, runCancel := context.WithTimeout(ctx, 30*time.Second)
		defer runCancel()

		result, err := events.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query Transfer events: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful query
		// DEBUG: Expand 'result.Events' in Variables panel
		// DEBUG: Each event has: BlockNumber, TxHash, LogIndex, From, To, Value

		fmt.Printf("Found %d Transfer events\n", len(result.Events))
		fmt.Println()

		// Display first few events
		// BREAKPOINT: Set breakpoint here
		displayCount := 5
		if len(result.Events) < displayCount {
			displayCount = len(result.Events)
		}

		for i := 0; i < displayCount; i++ {
			event := result.Events[i]
			fmt.Printf("Transfer #%d:\n", i+1)
			fmt.Printf("  Block:    %d\n", event.BlockNumber)
			fmt.Printf("  Tx Hash:  %s\n", event.TxHash.Hex())
			fmt.Printf("  Log Idx:  %d\n", event.LogIndex)
			fmt.Printf("  From:     %s\n", event.From.Hex())
			fmt.Printf("  To:       %s\n", event.To.Hex())
			fmt.Printf("  Value:    %s\n", event.Value.String())
			fmt.Println()
		}

		if len(result.Events) > displayCount {
			fmt.Printf("... and %d more events\n", len(result.Events)-displayCount)
			fmt.Println()
		}

		fmt.Println("ℹ All transfers queried (no filtering)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Example 2 - Filtered Transfer Query (By Sender)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through filtered query
	// DEBUG: This demonstrates filtering transfers by sender address
	// DEBUG: Uses Topic[1] filtering (from address)

	if demo == "all" || demo == "filtered" {
		fmt.Println("=== Example 2: Filtered Transfer Query (By Sender) ===")
		fmt.Println("Querying transfers FROM a specific address...")
		fmt.Println()

		// Use a well-known address for filtering (Uniswap V2 Router as example)
		uniswapRouter := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")
		fmt.Printf("Filtering by sender: %s\n", uniswapRouter.Hex())
		fmt.Println()

		// Create configuration with FromHolder filter
		// BREAKPOINT: Set breakpoint here
		// DEBUG: FromHolder filters Topic[1] (from address)
		cfg := events.Config{
			Token:      tokenAddress,
			FromBlock:  fromBlock,
			ToBlock:    toBlock,
			FromHolder: &uniswapRouter,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 30*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling events.Run
		// DEBUG: Step Into (F11) to see how FromHolder becomes a topic filter
		result, err := events.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query filtered Transfer events: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful query
		// DEBUG: Notice fewer events than unfiltered query
		// DEBUG: All events should have From == uniswapRouter

		fmt.Printf("Found %d Transfer events FROM %s\n", len(result.Events), uniswapRouter.Hex())
		fmt.Println()

		// Display events and verify filtering
		displayCount := 3
		if len(result.Events) < displayCount {
			displayCount = len(result.Events)
		}

		for i := 0; i < displayCount; i++ {
			event := result.Events[i]
			// BREAKPOINT: Set breakpoint inside loop
			// DEBUG: Verify event.From == uniswapRouter
			fmt.Printf("Transfer #%d:\n", i+1)
			fmt.Printf("  From:     %s ✓\n", event.From.Hex())
			fmt.Printf("  To:       %s\n", event.To.Hex())
			fmt.Printf("  Value:    %s\n", event.Value.String())
			fmt.Printf("  Tx Hash:  %s\n", event.TxHash.Hex())
			fmt.Println()
		}

		if len(result.Events) > displayCount {
			fmt.Printf("... and %d more events\n", len(result.Events)-displayCount)
			fmt.Println()
		}

		fmt.Println("ℹ Filtered by FROM address (Topic[1])")
		fmt.Println("  - Bloom filter narrows search to matching blocks")
		fmt.Println("  - Only transfers FROM specified address")
		fmt.Println()
	}

	// ============================================================================
	// STEP 7: Example 3 - Transfer Volume Analysis
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through analysis example
	// DEBUG: This demonstrates analyzing event data
	// DEBUG: Shows how events provide historical state change tracking

	if demo == "all" || demo == "analysis" {
		fmt.Println("=== Example 3: Transfer Volume Analysis ===")
		fmt.Println("Analyzing transfer patterns in block range...")
		fmt.Println()

		// Query all events for analysis
		cfg := events.Config{
			Token:     tokenAddress,
			FromBlock: fromBlock,
			ToBlock:   toBlock,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 30*time.Second)
		defer runCancel()

		result, err := events.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query events for analysis: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: We'll analyze the events to compute statistics

		fmt.Printf("Total Events: %d\n", len(result.Events))
		fmt.Println()

		if len(result.Events) > 0 {
			// Calculate total volume
			// BREAKPOINT: Set breakpoint here
			totalVolume := new(big.Int)
			uniqueSenders := make(map[common.Address]bool)
			uniqueReceivers := make(map[common.Address]bool)
			largestTransfer := result.Events[0]

			for _, event := range result.Events {
				// BREAKPOINT: Set breakpoint inside loop (first iteration)
				// DEBUG: Watch how we accumulate statistics
				totalVolume.Add(totalVolume, event.Value)
				uniqueSenders[event.From] = true
				uniqueReceivers[event.To] = true

				if event.Value.Cmp(largestTransfer.Value) > 0 {
					largestTransfer = event
				}
			}

			// Display analysis results
			// BREAKPOINT: Set breakpoint here
			fmt.Println("Volume Analysis:")
			fmt.Printf("  Total Volume:      %s (raw)\n", totalVolume.String())
			fmt.Printf("  Unique Senders:    %d addresses\n", len(uniqueSenders))
			fmt.Printf("  Unique Receivers:  %d addresses\n", len(uniqueReceivers))
			fmt.Printf("  Avg per Transfer:  %s (raw)\n", new(big.Int).Div(totalVolume, big.NewInt(int64(len(result.Events)))).String())
			fmt.Println()

			fmt.Println("Largest Transfer:")
			fmt.Printf("  From:     %s\n", largestTransfer.From.Hex())
			fmt.Printf("  To:       %s\n", largestTransfer.To.Hex())
			fmt.Printf("  Value:    %s (raw)\n", largestTransfer.Value.String())
			fmt.Printf("  Tx Hash:  %s\n", largestTransfer.TxHash.Hex())
			fmt.Printf("  Block:    %d\n", largestTransfer.BlockNumber)
			fmt.Println()

			fmt.Println("ℹ Event analysis demonstrates:")
			fmt.Println("  - Historical state change tracking")
			fmt.Println("  - Transfer pattern analysis")
			fmt.Println("  - Audit trail for token movements")
			fmt.Println()
		} else {
			fmt.Println("No transfers found in this block range")
			fmt.Println()
		}
	}

	// ============================================================================
	// Understanding Events and Logs
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand events/logs

	if demo == "all" {
		fmt.Println("=== Understanding Events and Logs ===")
		fmt.Println()

		fmt.Println("What are Events/Logs?")
		fmt.Println("  - Append-only records emitted during contract execution")
		fmt.Println("  - Stored in blockchain, but not part of state trie")
		fmt.Println("  - Cheaper than storage (events ~375 gas, storage ~20,000 gas)")
		fmt.Println("  - Cannot be read by smart contracts")
		fmt.Println()

		fmt.Println("Event Structure:")
		fmt.Println("  - Address: Which contract emitted the event")
		fmt.Println("  - Topics: Indexed parameters (max 4, first is event signature)")
		fmt.Println("  - Data: Non-indexed parameters (cheaper, not searchable)")
		fmt.Println("  - BlockNumber, TxHash, LogIndex: Metadata")
		fmt.Println()

		fmt.Println("Indexed vs Non-Indexed Parameters:")
		fmt.Println("  Indexed (Topics):")
		fmt.Println("    ✓ Searchable via bloom filters")
		fmt.Println("    ✓ Max 3 indexed params (+ event signature)")
		fmt.Println("    ✗ More expensive (gas cost)")
		fmt.Println()
		fmt.Println("  Non-Indexed (Data):")
		fmt.Println("    ✓ Cheaper to emit")
		fmt.Println("    ✓ Can store large data")
		fmt.Println("    ✗ Not searchable (must retrieve and decode)")
		fmt.Println()

		fmt.Println("Bloom Filters:")
		fmt.Println("  - Probabilistic data structure in block headers")
		fmt.Println("  - Quick check: \"might contain\" vs \"definitely doesn't contain\"")
		fmt.Println("  - Enables efficient log searching across millions of blocks")
		fmt.Println("  - False positives possible, but no false negatives")
		fmt.Println()

		fmt.Println("Why Events Matter:")
		fmt.Println("  1. Audit Trail: Immutable record of state changes")
		fmt.Println("  2. Off-Chain Indexing: Build searchable databases from events")
		fmt.Println("  3. Event-Driven Apps: React to blockchain state changes")
		fmt.Println("  4. Cost-Effective: Much cheaper than storing data on-chain")
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
	// 4. Use F11 (Step Into) to enter events package functions
	// 5. Watch Variables panel to see FilterQuery and log decoding
	// 6. Inspect Topics and Data in raw logs
	// 7. Use Call Stack panel to see function hierarchy
	//
	// EVENT CONCEPTS TO EXPLORE:
	// - Topics vs Data: How do they differ?
	// - Bloom Filters: How do they enable efficient searching?
	// - Event Signature: How is it computed?
	// - Topic Filtering: How does it narrow the search?
	// - Event Ordering: How are events sorted?
	//
	// STEPPING INTO EVENTS PACKAGE:
	// - Set breakpoint at events.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to solution.reference.go
	// - Step through to see: FilterQuery construction, FilterLogs call, decoding
	// - Watch how topics are built from addresses
	// - See how logs are decoded to TransferEvent structs
	//
	// COMMON ERRORS TO DEBUG:
	// - "filter logs": Invalid block range, RPC limits exceeded
	// - "log missing topics": Malformed event (wrong event type?)
	// - "unexpected topic": Wrong event signature (not a Transfer?)
	// - "data too short": Incomplete event data
}
