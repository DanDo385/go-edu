package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/14-explorer/internal/explorer"
)

/*
===================================================================
Ethereum Block Explorer
===================================================================

This program demonstrates how to build a simple block explorer by
fetching and displaying block data and transaction summaries.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <DEMO>
    go run ./cmd/app/main.go <RPC_URL> <BLOCK_NUMBER>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, latest, specific, txs, range (default: "all")
    [block_number]  - Specific block number to fetch

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com latest
    go run ./cmd/app/main.go https://eth.llamarpc.com specific
    go run ./cmd/app/main.go https://eth.llamarpc.com txs
    go run ./cmd/app/main.go https://eth.llamarpc.com 19000000

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
- Step Into (F11) to enter explorer package functions
- Watch Variables panel to see block and transaction data
- Inspect block headers and transaction summaries

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the explorer package from internal/explorer
- exercise.go, solution.reference.go are in package explorer
- The explorer.Run() function encapsulates all block fetching logic

KEY CONCEPTS:
- Block Structure: Header + Transactions + Uncles
- Block Header: Lightweight metadata about the block
- Transaction Summaries: Condensed transaction information
- Gas Usage: How much gas was consumed vs available
- Block Explorer: User-friendly view of blockchain data
- Optional Expansion: Fetch only what you need (header vs full block)
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the RPC URL (required)
	// DEBUG: os.Args[2] is the demo type or block number (optional)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [demo|block_number]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com latest\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com txs\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 19000000\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, latest, specific, txs, range\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	var blockNumber *big.Int
	if len(os.Args) > 2 {
		demo = os.Args[2]
		// Try to parse as block number
		if num, ok := new(big.Int).SetString(os.Args[2], 10); ok {
			blockNumber = num
			demo = "specific"
		}
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL', 'demo', and 'blockNumber' variables in Variables panel
	fmt.Printf("=== Ethereum Block Explorer ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	if blockNumber != nil {
		fmt.Printf("Block Number: %s\n", blockNumber.String())
	} else {
		fmt.Printf("Demo: %s\n", demo)
	}
	fmt.Println()

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
	// the RPCClient interface. It handles JSON-RPC communication with nodes.

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
	// STEP 3: Example 1 - Fetch Latest Block Without Transactions
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates fetching block metadata without transactions
	// DEBUG: Watch the 'result' variable to see block information
	//
	// CONCEPT: Block structure has three main components:
	//   - Header: Block metadata (number, hash, parent, gas, timestamps, etc.)
	//   - Transactions: List of transactions in this block
	//   - Uncles: List of uncle blocks (orphaned blocks)
	//
	// Fetching without transactions is much lighter and faster!

	if demo == "all" || demo == "latest" {
		fmt.Println("=== Example 1: Latest Block (Header Only) ===")
		fmt.Println("Fetching latest block without transaction details...")
		fmt.Println()

		// Create configuration for latest block without transactions
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice Number is nil (means latest)
		// DEBUG: Notice IncludeTxs is false (no transaction details)
		cfg := explorer.Config{
			Number:     nil,   // nil = latest block
			IncludeTxs: false, // Don't fetch transaction details
		}

		// BREAKPOINT: Set breakpoint here before calling explorer.Run
		// DEBUG: Step Into (F11) to see the implementation in solution.reference.go
		// DEBUG: You'll see how the function fetches block and extracts metadata
		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := explorer.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to fetch latest block: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful fetch
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at Number, Hash, Parent, TxCount, GasUsed, GasLimit
		// DEBUG: Notice Txs is empty (we didn't request transaction details)

		fmt.Printf("Block Number: %d\n", result.Number)
		fmt.Printf("Block Hash: %s\n", result.Hash.Hex())
		fmt.Printf("Parent Hash: %s\n", result.Parent.Hex())
		fmt.Printf("Transaction Count: %d\n", result.TxCount)
		fmt.Printf("Gas Used: %d (%.2f%%)\n", result.GasUsed,
			float64(result.GasUsed)/float64(result.GasLimit)*100)
		fmt.Printf("Gas Limit: %d\n", result.GasLimit)
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice how much information we get from just the header
		// DEBUG: This is what block explorers show on the main page
	}

	// ============================================================================
	// STEP 4: Example 2 - Fetch Block With Transaction Summaries
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through transaction example
	// DEBUG: This demonstrates fetching block with transaction details
	// DEBUG: Transactions are summarized for efficiency
	//
	// CONCEPT: Transaction summaries contain only essential information:
	//   - Hash: Transaction identifier
	//   - To: Recipient address (nil for contract creation)
	//   - Gas: Gas limit for this transaction
	//
	// Full transactions would include: from, value, data, signature, etc.

	if demo == "all" || demo == "txs" || blockNumber != nil {
		fmt.Println("=== Example 2: Block With Transaction Summaries ===")
		if blockNumber != nil {
			fmt.Printf("Fetching block %s with transaction details...\n", blockNumber.String())
		} else {
			fmt.Println("Fetching latest block with transaction details...")
		}
		fmt.Println()

		// Create configuration with transactions
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice IncludeTxs is true (we want transaction summaries)
		cfg := explorer.Config{
			Number:     blockNumber, // nil = latest, or specific block
			IncludeTxs: true,        // Include transaction summaries
		}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling explorer.Run
		// DEBUG: Step Into (F11) to see how transactions are processed
		result, err := explorer.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to fetch block: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful fetch
		// DEBUG: Expand result.Txs to see transaction summaries
		// DEBUG: Each TxSummary has Hash, To, and Gas fields

		fmt.Printf("Block Number: %d\n", result.Number)
		fmt.Printf("Block Hash: %s\n", result.Hash.Hex())
		fmt.Printf("Transaction Count: %d\n", result.TxCount)
		fmt.Println()

		// Display transaction summaries
		if len(result.Txs) > 0 {
			fmt.Println("Transaction Summaries:")
			// Show first 5 transactions (blocks can have 100s of transactions)
			maxToShow := 5
			if len(result.Txs) < maxToShow {
				maxToShow = len(result.Txs)
			}

			for i := 0; i < maxToShow; i++ {
				tx := result.Txs[i]
				fmt.Printf("  [%d] Hash: %s\n", i, tx.Hash.Hex())
				if tx.To != nil {
					fmt.Printf("      To: %s\n", tx.To.Hex())
				} else {
					fmt.Printf("      To: CONTRACT CREATION\n")
				}
				fmt.Printf("      Gas: %d\n", tx.Gas)
			}

			if len(result.Txs) > maxToShow {
				fmt.Printf("  ... and %d more transactions\n", len(result.Txs)-maxToShow)
			}
			fmt.Println()
		} else {
			fmt.Println("Block contains no transactions (empty block)\n")
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice how transaction summaries are lightweight
		// DEBUG: This is perfect for block explorer list views
	}

	// ============================================================================
	// STEP 5: Example 3 - Fetch Specific Historical Block
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through historical block example
	// DEBUG: This demonstrates fetching a specific block from history

	if demo == "all" || demo == "specific" {
		fmt.Println("=== Example 3: Specific Historical Block ===")
		fmt.Println("Fetching block 1,000,000 (early Ethereum history)...")
		fmt.Println()

		// Famous block 1,000,000 (first major milestone)
		// BREAKPOINT: Set breakpoint here
		historicalBlock := big.NewInt(1_000_000)

		cfg := explorer.Config{
			Number:     historicalBlock,
			IncludeTxs: true,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := explorer.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to fetch historical block: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This block is from 2016—early Ethereum days
		// DEBUG: Notice lower gas usage compared to modern blocks

		fmt.Printf("Block Number: %d\n", result.Number)
		fmt.Printf("Block Hash: %s\n", result.Hash.Hex())
		fmt.Printf("Parent Hash: %s\n", result.Parent.Hex())
		fmt.Printf("Transaction Count: %d\n", result.TxCount)
		fmt.Printf("Gas Used: %d\n", result.GasUsed)
		fmt.Printf("Gas Limit: %d\n", result.GasLimit)
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Historical blocks help understand blockchain evolution
	}

	// ============================================================================
	// STEP 6: Example 4 - Understanding Block Structure
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand block structure
	// DEBUG: This demonstrates the relationship between blocks

	if demo == "all" {
		fmt.Println("=== Example 4: Understanding Block Structure ===")
		fmt.Println()

		fmt.Println("Block Components:")
		fmt.Println("  1. Header: Metadata about the block")
		fmt.Println("    - Number: Block height in the chain")
		fmt.Println("    - Hash: Unique identifier (keccak256 of header)")
		fmt.Println("    - Parent: Hash of previous block (forms the chain!)")
		fmt.Println("    - StateRoot: Merkle root of account state")
		fmt.Println("    - TxHash: Merkle root of transactions")
		fmt.Println("    - ReceiptHash: Merkle root of receipts")
		fmt.Println("    - Timestamp: When the block was mined")
		fmt.Println("    - GasLimit: Maximum gas allowed in block")
		fmt.Println("    - GasUsed: Actual gas consumed by transactions")
		fmt.Println()
		fmt.Println("  2. Transactions: List of transactions in this block")
		fmt.Println("    - Each transaction has: hash, from, to, value, data, gas, etc.")
		fmt.Println()
		fmt.Println("  3. Uncles: Orphaned blocks (Ethereum PoW only)")
		fmt.Println("    - Blocks mined at same time but not included in main chain")
		fmt.Println()

		fmt.Println("Block Explorer Views:")
		fmt.Println("  - List View: Block number, hash, txCount, gas (header only)")
		fmt.Println("  - Detail View: All header fields + transaction summaries")
		fmt.Println("  - Full View: Complete block data + full transactions")
		fmt.Println()

		fmt.Println("Parent Hash Chain:")
		fmt.Println("  Block N → Parent = Hash of Block N-1")
		fmt.Println("  Block N-1 → Parent = Hash of Block N-2")
		fmt.Println("  ... all the way back to Genesis Block (Block 0)")
		fmt.Println()
		fmt.Println("  This chain of hashes makes the blockchain immutable!")
		fmt.Println("  Changing any historical block would change all subsequent hashes.")
		fmt.Println()

		fmt.Println("Gas Metrics:")
		fmt.Println("  - GasLimit: Maximum gas per block (set by protocol/miners)")
		fmt.Println("  - GasUsed: Actual gas consumed (sum of all tx gas)")
		fmt.Println("  - Utilization: GasUsed / GasLimit (how full is the block?)")
		fmt.Println("  - High utilization → Network congestion → Higher fees")
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Understanding block structure is essential for blockchain development
		// DEBUG: Block explorers just present this data in user-friendly formats
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
	// 4. Use F11 (Step Into) to enter explorer package functions
	// 5. Watch Variables panel to see block and transaction data
	// 6. Inspect result fields to understand block structure
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Block Header: What metadata is in every block?
	// - Parent Hash: How does it form the blockchain?
	// - Merkle Roots: How do they commit to transactions and state?
	// - Gas Limit vs Gas Used: What's the difference?
	// - Optional Expansion: When to fetch full data vs summaries?
	//
	// STEPPING INTO EXPLORER PACKAGE:
	// - Set breakpoint at explorer.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see: validation, BlockByNumber call, data extraction
	// - Watch how transaction summaries are built from full transactions
	// - See how we extract only needed fields for efficiency
	//
	// COMMON EXPLORER ERRORS TO DEBUG:
	// - Block not found: Verify block exists on this network
	// - Invalid block number: Must be non-negative
	// - Historical block unavailable: Some nodes don't keep old blocks
	// - Connection timeout: Network issues or slow endpoint
}
