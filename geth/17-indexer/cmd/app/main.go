package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
===================================================================
Ethereum Block Indexer
===================================================================

This program demonstrates how to build a blockchain indexer that fetches
and processes blocks from an Ethereum node. Indexers are fundamental
infrastructure for blockchain analytics, block explorers, and applications
that need to track on-chain events in real-time or historically.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <demo>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, single, range, live (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com single
    go run ./cmd/app/main.go https://eth.llamarpc.com range
    go run ./cmd/app/main.go https://rpc.sepolia.org live

Common RPC Endpoints:
    Ethereum Mainnet:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth

    Sepolia Testnet:
      - https://rpc.sepolia.org
      - https://sepolia.gateway.tenderly.co

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter indexing functions
- Watch Variables panel to see block data structures
- Inspect block headers, transactions, and receipts

KEY CONCEPTS:
- Block: Container for transactions with metadata and cryptographic proofs
- Block Header: Lightweight block metadata (number, hash, parent, timestamp, etc.)
- Block Body: Transactions and uncles (ommers) included in the block
- Block Indexing: Process of fetching and storing blockchain data for querying
- Full Block vs Header: Full blocks include all transactions; headers are lightweight
- Sequential Indexing: Process blocks one-by-one in order (simple but slow)
- Concurrent Indexing: Process multiple blocks in parallel (fast but complex)

COMPUTER SCIENCE CONCEPTS:
- Sequential vs Parallel Processing: Trade-off between simplicity and performance
- Data Pipeline: Fetch → Parse → Process → Store pattern
- Backfilling: Indexing historical data from genesis to current block
- Live Indexing: Continuously indexing new blocks as they're produced
- Checkpointing: Saving progress to resume after interruption
- Idempotency: Processing the same block multiple times yields same result

INDEXER ARCHITECTURE:
1. Block Fetcher: Retrieves blocks from RPC endpoint
2. Block Parser: Extracts relevant data (transactions, logs, traces)
3. Block Processor: Transforms data for storage
4. Storage Layer: Persists processed data (database, cache)
5. Progress Tracker: Remembers last indexed block

REAL-WORLD APPLICATIONS:
- Block Explorers: Etherscan, Blockscout
- Analytics Platforms: Dune, Nansen
- DeFi Applications: Track token transfers, price feeds, events
- NFT Marketplaces: Monitor mints, sales, transfers
- Infrastructure Monitoring: Node health, network statistics
*/

// BlockInfo represents processed block information for display
type BlockInfo struct {
	Number       *big.Int
	Hash         string
	ParentHash   string
	Timestamp    uint64
	Miner        string
	TxCount      int
	GasUsed      uint64
	GasLimit     uint64
	BaseFee      *big.Int
	Difficulty   *big.Int
	Size         uint64
}

// extractBlockInfo converts a full block into a structured BlockInfo
// BREAKPOINT: Set breakpoint here to see block parsing
// DEBUG: Watch how raw block data is transformed into structured info
func extractBlockInfo(block *types.Block) BlockInfo {
	header := block.Header()

	return BlockInfo{
		Number:       new(big.Int).Set(block.Number()),
		Hash:         block.Hash().Hex(),
		ParentHash:   header.ParentHash.Hex(),
		Timestamp:    header.Time,
		Miner:        header.Coinbase.Hex(),
		TxCount:      len(block.Transactions()),
		GasUsed:      header.GasUsed,
		GasLimit:     header.GasLimit,
		BaseFee:      new(big.Int).Set(header.BaseFee),
		Difficulty:   new(big.Int).Set(header.Difficulty),
		Size:         block.Size(),
	}
}

// displayBlockInfo prints formatted block information
func displayBlockInfo(info BlockInfo) {
	fmt.Printf("Block #%s\n", info.Number.String())
	fmt.Printf("  Hash:       %s\n", info.Hash)
	fmt.Printf("  Parent:     %s\n", info.ParentHash)
	fmt.Printf("  Timestamp:  %d (%s)\n", info.Timestamp, time.Unix(int64(info.Timestamp), 0))
	fmt.Printf("  Miner:      %s\n", info.Miner)
	fmt.Printf("  Txs:        %d\n", info.TxCount)
	fmt.Printf("  Gas:        %d / %d (%.2f%% full)\n",
		info.GasUsed, info.GasLimit,
		float64(info.GasUsed)/float64(info.GasLimit)*100)
	fmt.Printf("  Base Fee:   %s wei\n", info.BaseFee.String())
	fmt.Printf("  Difficulty: %s\n", info.Difficulty.String())
	fmt.Printf("  Size:       %d bytes\n", info.Size)
	fmt.Println()
}

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
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com single\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com range\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, single, range, live\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Block Indexer ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection
	// DEBUG: This is where we establish the RPC connection
	// DEBUG: Watch for any connection errors in the error variable

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
	// STEP 3: Example 1 - Index a Single Block
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through single block indexing
	// DEBUG: This demonstrates the basic building block of all indexers
	// DEBUG: Watch how we fetch and parse a single block

	if demo == "all" || demo == "single" {
		fmt.Println("=== Example 1: Index a Single Block ===")
		fmt.Println("Fetching a specific block by number...")
		fmt.Println()

		// Fetch a well-known historical block
		// Block 1,000,000 is a milestone block from 2016
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice we're using BlockByNumber with a specific number
		blockNumber := big.NewInt(1_000_000)

		blockCtx, blockCancel := context.WithTimeout(ctx, 5*time.Second)
		defer blockCancel()

		// BREAKPOINT: Set breakpoint here before fetching block
		// DEBUG: Step Into (F11) to see the RPC call
		block, err := client.BlockByNumber(blockCtx, blockNumber)
		if err != nil {
			log.Fatalf("Failed to fetch block: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after fetching block
		// DEBUG: Expand 'block' in Variables panel to inspect all fields
		// DEBUG: Notice the block contains header, transactions, uncles

		// Extract and display block information
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Step Into (F11) extractBlockInfo to see parsing logic
		info := extractBlockInfo(block)
		displayBlockInfo(info)

		// Display transaction details
		fmt.Printf("Transaction Details:\n")
		txs := block.Transactions()
		if len(txs) > 0 {
			// Show first 3 transactions as examples
			displayCount := 3
			if len(txs) < displayCount {
				displayCount = len(txs)
			}

			for i := 0; i < displayCount; i++ {
				tx := txs[i]
				fmt.Printf("  Tx #%d: %s\n", i+1, tx.Hash().Hex())
				fmt.Printf("    From: (derive from signature)\n")
				fmt.Printf("    To:   %s\n", tx.To().Hex())
				fmt.Printf("    Value: %s wei\n", tx.Value().String())
				fmt.Printf("    Gas:   %d\n", tx.Gas())
				fmt.Println()
			}

			if len(txs) > displayCount {
				fmt.Printf("  ... and %d more transactions\n\n", len(txs)-displayCount)
			}
		} else {
			fmt.Println("  No transactions in this block\n")
		}

		// Educational explanation
		fmt.Println("Single Block Indexing Insights:")
		fmt.Println("  1. BlockByNumber fetches full block with all transactions")
		fmt.Println("  2. Each block contains header (metadata) and body (transactions)")
		fmt.Println("  3. Block hash uniquely identifies the block")
		fmt.Println("  4. Parent hash links to previous block (forms the chain)")
		fmt.Println("  5. Transaction count varies (0 to thousands per block)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 2 - Index a Range of Blocks
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through range indexing
	// DEBUG: This demonstrates sequential indexing of multiple blocks
	// DEBUG: Watch how we process blocks one-by-one in order

	if demo == "all" || demo == "range" {
		fmt.Println("=== Example 2: Index a Range of Blocks ===")
		fmt.Println("Sequentially indexing blocks 1,000,000 - 1,000,010...")
		fmt.Println()

		startBlock := big.NewInt(1_000_000)
		endBlock := big.NewInt(1_000_010)

		// BREAKPOINT: Set breakpoint here before loop
		// DEBUG: Watch how we iterate through block numbers
		current := new(big.Int).Set(startBlock)
		one := big.NewInt(1)

		fetchStart := time.Now()
		blocksProcessed := 0

		for current.Cmp(endBlock) <= 0 {
			// BREAKPOINT: Set breakpoint here in loop
			// DEBUG: Watch 'current' variable increment each iteration
			// DEBUG: See how long each block fetch takes

			blockCtx, blockCancel := context.WithTimeout(ctx, 5*time.Second)

			block, err := client.BlockByNumber(blockCtx, current)
			blockCancel()

			if err != nil {
				fmt.Printf("⚠ Failed to fetch block %s: %v\n", current.String(), err)
				current.Add(current, one)
				continue
			}

			// Extract basic info
			info := extractBlockInfo(block)
			fmt.Printf("Indexed block #%s - %d txs - %s\n",
				info.Number.String(), info.TxCount,
				time.Unix(int64(info.Timestamp), 0).Format("2006-01-02 15:04:05"))

			blocksProcessed++
			current.Add(current, one)
		}

		fetchElapsed := time.Since(fetchStart)

		fmt.Println()
		fmt.Printf("✓ Indexed %d blocks in %v\n", blocksProcessed, fetchElapsed)
		fmt.Printf("Average: %v per block\n\n",
			time.Duration(int64(fetchElapsed)/int64(blocksProcessed)))

		// Educational explanation
		fmt.Println("Range Indexing Insights:")
		fmt.Println("  1. Sequential indexing processes blocks in order")
		fmt.Println("  2. Each RPC call takes time (network latency + processing)")
		fmt.Println("  3. 10 blocks takes ~10x time of 1 block (no parallelism)")
		fmt.Println("  4. For large ranges, sequential indexing is slow")
		fmt.Println("  5. Solution: Concurrent indexing (worker pool pattern)")
		fmt.Println()

		fmt.Println("Optimization Strategies:")
		fmt.Println("  - Concurrent fetching: Use worker pools (module 16)")
		fmt.Println("  - Batch requests: Some RPC methods support batching")
		fmt.Println("  - Header-only: Fetch headers first, full blocks later")
		fmt.Println("  - Checkpointing: Save progress to resume after interruption")
		fmt.Println("  - Caching: Store frequently accessed blocks locally")
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 3 - Live Block Indexing
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through live indexing
	// DEBUG: This demonstrates indexing new blocks as they're produced
	// DEBUG: Watch how we poll for the latest block number

	if demo == "all" || demo == "live" {
		fmt.Println("=== Example 3: Live Block Indexing ===")
		fmt.Println("Monitoring for new blocks (will index next 3 blocks)...")
		fmt.Println()

		// Get current block number
		headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
		currentHeader, err := client.HeaderByNumber(headerCtx, nil)
		headerCancel()
		if err != nil {
			log.Fatalf("Failed to get current block: %v\n", err)
		}

		lastSeenBlock := new(big.Int).Set(currentHeader.Number)
		blocksToIndex := 3
		blocksIndexed := 0

		fmt.Printf("Starting from block #%s\n", lastSeenBlock.String())
		fmt.Println("Waiting for new blocks...\n")

		// BREAKPOINT: Set breakpoint here before polling loop
		// DEBUG: Watch how we continuously poll for new blocks
		ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds
		defer ticker.Stop()

		for blocksIndexed < blocksToIndex {
			// BREAKPOINT: Set breakpoint here in polling loop
			// DEBUG: See how often we check for new blocks

			select {
			case <-ticker.C:
				// Check for new block
				headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
				latestHeader, err := client.HeaderByNumber(headerCtx, nil)
				headerCancel()

				if err != nil {
					fmt.Printf("⚠ Failed to fetch latest block: %v\n", err)
					continue
				}

				latestNumber := latestHeader.Number

				// BREAKPOINT: Set breakpoint here
				// DEBUG: Compare lastSeenBlock with latestNumber
				// DEBUG: Watch how we detect new blocks

				if latestNumber.Cmp(lastSeenBlock) > 0 {
					// New block(s) detected!
					// Process all new blocks
					current := new(big.Int).Add(lastSeenBlock, big.NewInt(1))

					for current.Cmp(latestNumber) <= 0 && blocksIndexed < blocksToIndex {
						blockCtx, blockCancel := context.WithTimeout(ctx, 5*time.Second)
						block, err := client.BlockByNumber(blockCtx, current)
						blockCancel()

						if err != nil {
							fmt.Printf("⚠ Failed to fetch block %s: %v\n", current.String(), err)
							current.Add(current, big.NewInt(1))
							continue
						}

						// BREAKPOINT: Set breakpoint here when new block found
						// DEBUG: Inspect the new block's data
						info := extractBlockInfo(block)

						fmt.Printf("📦 New block #%s mined!\n", info.Number.String())
						fmt.Printf("   Hash: %s\n", info.Hash)
						fmt.Printf("   Miner: %s\n", info.Miner)
						fmt.Printf("   Txs: %d\n", info.TxCount)
						fmt.Printf("   Timestamp: %s\n",
							time.Unix(int64(info.Timestamp), 0).Format("2006-01-02 15:04:05"))
						fmt.Println()

						blocksIndexed++
						current.Add(current, big.NewInt(1))
					}

					lastSeenBlock.Set(latestNumber)
				}
			}
		}

		fmt.Printf("✓ Indexed %d new blocks\n\n", blocksIndexed)

		// Educational explanation
		fmt.Println("Live Indexing Insights:")
		fmt.Println("  1. New blocks are produced ~every 12 seconds on Ethereum")
		fmt.Println("  2. Polling: Periodically check for latest block number")
		fmt.Println("  3. Detect new blocks: Compare with last seen block")
		fmt.Println("  4. Handle gaps: Process all missed blocks (network delays)")
		fmt.Println("  5. Alternative: Use subscriptions (eth_subscribe via WebSocket)")
		fmt.Println()

		fmt.Println("Production Considerations:")
		fmt.Println("  - Reorg handling: Blocks can be reorganized (module 18)")
		fmt.Println("  - Missed blocks: Always check for gaps in block numbers")
		fmt.Println("  - Confirmations: Wait N blocks before considering finalized")
		fmt.Println("  - WebSocket subscriptions: More efficient than polling")
		fmt.Println("  - Error recovery: Retry failed fetches, alert on persistent errors")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Indexer Architecture Patterns
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand indexer architecture
	// DEBUG: This section explains production indexer design patterns

	if demo == "all" {
		fmt.Println("=== Indexer Architecture Patterns ===")
		fmt.Println()

		fmt.Println("Indexer Pipeline Stages:")
		fmt.Println("  1. Block Fetcher: Retrieves blocks from RPC endpoint")
		fmt.Println("  2. Block Parser: Extracts transactions, logs, traces")
		fmt.Println("  3. Data Transformer: Converts to application-specific format")
		fmt.Println("  4. Storage Layer: Persists to database (PostgreSQL, MongoDB, etc.)")
		fmt.Println("  5. Cache Layer: Speeds up frequent queries (Redis, Memcached)")
		fmt.Println("  6. API Layer: Exposes indexed data to applications")
		fmt.Println()

		fmt.Println("Indexing Strategies:")
		fmt.Println("  Sequential: Simple, slow, guaranteed order")
		fmt.Println("  Concurrent: Fast, complex, requires ordering logic")
		fmt.Println("  Hybrid: Sequential for recent blocks, concurrent for backfill")
		fmt.Println("  Incremental: Index only new blocks, starting from checkpoint")
		fmt.Println("  Full resync: Reindex from genesis (for schema changes)")
		fmt.Println()

		fmt.Println("Performance Optimizations:")
		fmt.Println("  - Worker pools: Parallel block fetching (module 16)")
		fmt.Println("  - Batch inserts: Group database writes for efficiency")
		fmt.Println("  - Header-first: Fetch headers, then full blocks as needed")
		fmt.Println("  - Selective indexing: Only index relevant contracts/events")
		fmt.Println("  - Connection pooling: Reuse RPC connections")
		fmt.Println()

		fmt.Println("Data to Index:")
		fmt.Println("  - Blocks: Number, hash, timestamp, miner, difficulty")
		fmt.Println("  - Transactions: Hash, from, to, value, input data, gas")
		fmt.Println("  - Logs/Events: Contract events emitted during execution")
		fmt.Println("  - Traces: Internal transactions and call traces")
		fmt.Println("  - State: Account balances, contract storage (advanced)")
		fmt.Println()

		fmt.Println("Common Use Cases:")
		fmt.Println("  - Token transfers: Track ERC20/ERC721 transfers")
		fmt.Println("  - DeFi analytics: Monitor swaps, liquidations, yields")
		fmt.Println("  - NFT tracking: Index mints, sales, ownership")
		fmt.Println("  - Address history: All transactions for an address")
		fmt.Println("  - Contract interactions: Who called which contracts")
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
	// 4. Use F11 (Step Into) to enter helper functions
	// 5. Watch Variables panel to see block structures
	// 6. Inspect block.Header() to see block metadata
	// 7. Inspect block.Transactions() to see transaction list
	// 8. Use Call Stack panel to trace function calls
	//
	// BLOCK STRUCTURE EXPLORATION:
	// - Set breakpoint after BlockByNumber call
	// - Expand 'block' variable in Variables panel
	// - Examine Header: Number, Hash, ParentHash, Time, etc.
	// - Examine Transactions: List of all txs in block
	// - Examine Uncles: Ommer blocks (usually empty on modern chains)
	//
	// PERFORMANCE PROFILING:
	// - Measure time for each BlockByNumber call
	// - Compare sequential vs concurrent indexing performance
	// - Identify bottlenecks: Network latency, parsing, storage
	// - Use Go's pprof for CPU and memory profiling
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Why do blocks have numbers AND hashes? Numbers are sequential,
	//   hashes are cryptographic proofs. Both needed for different purposes.
	// - What's a block's parent hash? Links to previous block, forming chain.
	//   If parent hash changes, entire chain from that point is invalid.
	// - Why index blocks? RPC nodes don't provide advanced queries. Indexers
	//   enable: "Find all transfers to address X" or "Show NFT ownership history"
	// - What's the difference between full sync and light sync? Full sync
	//   validates all blocks, light sync trusts headers. Indexers need full data.
}
