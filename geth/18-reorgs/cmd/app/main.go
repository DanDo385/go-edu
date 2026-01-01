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
Blockchain Reorganization (Reorg) Detection and Handling
===================================================================

This program demonstrates how to detect and handle blockchain reorganizations
(reorgs), a critical concept for any application that processes blockchain
data in real-time. Reorgs occur when the network reaches consensus on a
different chain than what was previously seen, causing blocks to be replaced.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <demo>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, detection, confirmation, rollback (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com detection
    go run ./cmd/app/main.go https://rpc.sepolia.org confirmation
    go run ./cmd/app/main.go https://rpc.sepolia.org rollback

Common RPC Endpoints:
    Ethereum Mainnet:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth

    Sepolia Testnet (more frequent reorgs):
      - https://rpc.sepolia.org
      - https://sepolia.gateway.tenderly.co

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter reorg detection functions
- Watch Variables panel to see block hash changes
- Inspect parent hash chains to trace reorganizations

KEY CONCEPTS:
- Blockchain: Chain of blocks linked by parent hashes
- Fork: Multiple competing chains at the same height
- Reorganization (Reorg): Network switches to different fork as canonical chain
- Canonical Chain: The "main" chain accepted by network consensus
- Uncle/Ommer: Valid block that's not in the canonical chain
- Finality: Point after which blocks are extremely unlikely to reorg
- Confirmations: Number of blocks built on top (more = less likely to reorg)

COMPUTER SCIENCE CONCEPTS:
- Consensus: Distributed agreement on canonical state
- Eventual Consistency: System eventually converges to consistent state
- CAP Theorem: Trade-offs between Consistency, Availability, Partition tolerance
- Probabilistic Finality: Probability of reorg decreases with each confirmation
- State Rollback: Reverting state changes when reorg occurs

WHY REORGS HAPPEN:
1. Network Latency: Two miners produce blocks simultaneously
2. Chain Splits: Network temporarily divided (e.g., network partition)
3. Mining Randomness: Longest chain rule creates natural forks
4. Attacks: 51% attack (attacker controls majority of hash power)

REORG DEPTH:
- 1-block reorg: Common, happens several times per day
- 2-3 block reorg: Uncommon but normal, happens occasionally
- 4+ block reorg: Rare, indicates network issues or attack
- 6+ block reorg: Very rare on mainnet, likely attack (Bitcoin uses 6 confirmations)

REAL-WORLD IMPACT:
- Exchanges: Wait N confirmations before crediting deposits
- DeFi: Oracle price feeds may temporarily report wrong prices
- NFT Marketplaces: Sale transactions may be reversed
- Block Explorers: Must update displayed data when reorg occurs
- Indexers: Must roll back indexed data and reprocess

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging
- Watch Variables panel to see block hash changes
- Trace parent hash chains during reorganizations
*/

// BlockSummary represents simplified block information for reorg tracking
type BlockSummary struct {
	Number     *big.Int
	Hash       string
	ParentHash string
	Timestamp  uint64
}

// extractBlockSummary extracts key fields for reorg detection
// BREAKPOINT: Set breakpoint here to see block data extraction
func extractBlockSummary(header *types.Header) BlockSummary {
	return BlockSummary{
		Number:     new(big.Int).Set(header.Number),
		Hash:       header.Hash().Hex(),
		ParentHash: header.ParentHash.Hex(),
		Timestamp:  header.Time,
	}
}

// BlockHistory tracks historical blocks for reorg detection
type BlockHistory struct {
	blocks map[int64]BlockSummary // blockNumber -> BlockSummary
}

// NewBlockHistory creates a new block history tracker
func NewBlockHistory() *BlockHistory {
	return &BlockHistory{
		blocks: make(map[int64]BlockSummary),
	}
}

// Store saves a block in history
func (bh *BlockHistory) Store(summary BlockSummary) {
	bh.blocks[summary.Number.Int64()] = summary
}

// Get retrieves a block from history
func (bh *BlockHistory) Get(number int64) (BlockSummary, bool) {
	summary, exists := bh.blocks[number]
	return summary, exists
}

// DetectReorg checks if a new block represents a reorganization
// Returns true if the parent hash doesn't match our stored block at Number-1
// BREAKPOINT: Set breakpoint here to see reorg detection logic
func (bh *BlockHistory) DetectReorg(newBlock BlockSummary) bool {
	// Get the expected parent block (Number - 1)
	expectedParentNumber := new(big.Int).Sub(newBlock.Number, big.NewInt(1)).Int64()

	storedParent, exists := bh.Get(expectedParentNumber)
	if !exists {
		// No history for parent block, can't detect reorg
		return false
	}

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Compare newBlock.ParentHash with storedParent.Hash
	// DEBUG: If they don't match, a reorg occurred

	// Reorg detected if parent hashes don't match
	return newBlock.ParentHash != storedParent.Hash
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
		fmt.Fprintf(os.Stderr, "  %s https://rpc.sepolia.org detection\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, detection, confirmation, rollback\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Blockchain Reorganization Detection ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Connect to Ethereum RPC Endpoint
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

	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 3: Example 1 - Understanding Reorg Detection
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through reorg detection
	// DEBUG: This demonstrates the basic concept of reorg detection

	if demo == "all" || demo == "detection" {
		fmt.Println("=== Example 1: Reorg Detection Mechanism ===")
		fmt.Println("Demonstrating how to detect reorganizations...")
		fmt.Println()

		// Fetch a range of blocks to understand parent-child relationships
		// BREAKPOINT: Set breakpoint here
		startBlock := big.NewInt(1_000_000)
		blockCount := 5

		fmt.Println("Fetching block chain to examine parent-child links:")
		fmt.Println()

		var previousBlock *BlockSummary

		for i := 0; i < blockCount; i++ {
			blockNumber := new(big.Int).Add(startBlock, big.NewInt(int64(i)))

			headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
			header, err := client.HeaderByNumber(headerCtx, blockNumber)
			headerCancel()

			if err != nil {
				log.Fatalf("Failed to fetch block %s: %v\n", blockNumber.String(), err)
			}

			// BREAKPOINT: Set breakpoint here in loop
			// DEBUG: Watch how each block's ParentHash matches previous block's Hash
			summary := extractBlockSummary(header)

			fmt.Printf("Block #%s:\n", summary.Number.String())
			fmt.Printf("  Hash:       %s\n", summary.Hash)
			fmt.Printf("  Parent:     %s\n", summary.ParentHash)

			if previousBlock != nil {
				// Check if parent hash matches previous block hash
				// BREAKPOINT: Set breakpoint here
				// DEBUG: In normal chain, summary.ParentHash == previousBlock.Hash
				if summary.ParentHash == previousBlock.Hash {
					fmt.Printf("  ✓ Links to previous block (normal chain)\n")
				} else {
					fmt.Printf("  ⚠ REORG: Parent doesn't match previous block!\n")
					fmt.Printf("    Expected parent: %s\n", previousBlock.Hash)
					fmt.Printf("    Actual parent:   %s\n", summary.ParentHash)
				}
			}
			fmt.Println()

			previousBlock = &summary
		}

		// Educational explanation
		fmt.Println("Reorg Detection Insights:")
		fmt.Println("  1. Each block contains its parent's hash")
		fmt.Println("  2. In normal chain: Block N's parent = Block N-1's hash")
		fmt.Println("  3. During reorg: Block N's parent ≠ our stored Block N-1's hash")
		fmt.Println("  4. This indicates Block N-1 was replaced with different block")
		fmt.Println("  5. Must reprocess all affected blocks after detecting reorg")
		fmt.Println()

		fmt.Println("Parent Hash Chain:")
		fmt.Println("  Block 100: Hash=A, ParentHash=<99's hash>")
		fmt.Println("  Block 101: Hash=B, ParentHash=A  ✓ Valid chain")
		fmt.Println("  Block 102: Hash=C, ParentHash=B  ✓ Valid chain")
		fmt.Println()
		fmt.Println("  If Block 102 reorgs:")
		fmt.Println("  Block 102: Hash=D, ParentHash=X  ⚠ Reorg! X ≠ B")
		fmt.Println("  This means Block 101 was replaced too")
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 2 - Confirmation Depth Strategy
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through confirmation strategy
	// DEBUG: This demonstrates why confirmations matter for safety

	if demo == "all" || demo == "confirmation" {
		fmt.Println("=== Example 2: Confirmation Depth Strategy ===")
		fmt.Println("Understanding how confirmations reduce reorg risk...")
		fmt.Println()

		// Get latest block
		headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
		latestHeader, err := client.HeaderByNumber(headerCtx, nil)
		headerCancel()
		if err != nil {
			log.Fatalf("Failed to get latest block: %v\n", err)
		}

		latestNumber := latestHeader.Number
		latestSummary := extractBlockSummary(latestHeader)

		fmt.Printf("Latest Block #%s:\n", latestSummary.Number.String())
		fmt.Printf("  Hash: %s\n\n", latestSummary.Hash)

		// Show blocks at different confirmation depths
		// BREAKPOINT: Set breakpoint here
		confirmationDepths := []int64{0, 1, 6, 12, 32, 64}

		fmt.Println("Confirmation Depth Analysis:")
		fmt.Println()

		for _, depth := range confirmationDepths {
			blockNumber := new(big.Int).Sub(latestNumber, big.NewInt(depth))

			// Skip if block number would be negative
			if blockNumber.Sign() < 0 {
				continue
			}

			headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
			header, err := client.HeaderByNumber(headerCtx, blockNumber)
			headerCancel()

			if err != nil {
				fmt.Printf("⚠ Failed to fetch block %s: %v\n", blockNumber.String(), err)
				continue
			}

			summary := extractBlockSummary(header)

			// BREAKPOINT: Set breakpoint here in loop
			// DEBUG: Watch how confirmation count affects reorg likelihood

			var reorgLikelihood string
			var recommendation string

			switch {
			case depth == 0:
				reorgLikelihood = "Very High (just mined)"
				recommendation = "⚠ Do NOT trust - likely to change"
			case depth == 1:
				reorgLikelihood = "High (1 confirmation)"
				recommendation = "⚠ Wait for more confirmations"
			case depth >= 6 && depth < 12:
				reorgLikelihood = "Low (6+ confirmations)"
				recommendation = "✓ Safe for most applications"
			case depth >= 12 && depth < 32:
				reorgLikelihood = "Very Low (12+ confirmations)"
				recommendation = "✓ Safe for high-value transactions"
			case depth >= 32:
				reorgLikelihood = "Extremely Low (32+ confirmations)"
				recommendation = "✓ Practically final (near finality)"
			}

			fmt.Printf("%d Confirmations (Block #%s):\n", depth, summary.Number.String())
			fmt.Printf("  Hash: %s\n", summary.Hash[:16]+"...")
			fmt.Printf("  Reorg Likelihood: %s\n", reorgLikelihood)
			fmt.Printf("  Recommendation: %s\n\n", recommendation)
		}

		// Educational explanation
		fmt.Println("Confirmation Depth Insights:")
		fmt.Println("  - 0 confirmations: Block just mined, high reorg risk")
		fmt.Println("  - 1-5 confirmations: Reorgs still possible")
		fmt.Println("  - 6 confirmations: Bitcoin standard, low risk")
		fmt.Println("  - 12 confirmations: Very low risk for Ethereum")
		fmt.Println("  - 32 confirmations: Ethereum finality checkpoint")
		fmt.Println("  - 64+ confirmations: Practically irreversible")
		fmt.Println()

		fmt.Println("Practical Guidelines:")
		fmt.Println("  - Small transactions: 6 confirmations")
		fmt.Println("  - Large transactions: 12+ confirmations")
		fmt.Println("  - Exchange deposits: 12-30 confirmations")
		fmt.Println("  - Critical operations: Wait for finality (32+ blocks)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 3 - Rollback Simulation
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through rollback simulation
	// DEBUG: This demonstrates how to handle reorgs in applications

	if demo == "all" || demo == "rollback" {
		fmt.Println("=== Example 3: Reorg Rollback Handling ===")
		fmt.Println("Simulating how applications should handle reorgs...")
		fmt.Println()

		// Initialize block history tracker
		history := NewBlockHistory()

		// Fetch initial chain
		headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
		latestHeader, err := client.HeaderByNumber(headerCtx, nil)
		headerCancel()
		if err != nil {
			log.Fatalf("Failed to get latest block: %v\n", err)
		}

		// Store last 10 blocks in history
		// BREAKPOINT: Set breakpoint here
		startNumber := new(big.Int).Sub(latestHeader.Number, big.NewInt(9))

		fmt.Println("Building block history (last 10 blocks)...")

		for i := int64(0); i < 10; i++ {
			blockNumber := new(big.Int).Add(startNumber, big.NewInt(i))

			headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
			header, err := client.HeaderByNumber(headerCtx, blockNumber)
			headerCancel()

			if err != nil {
				log.Printf("Failed to fetch block %s: %v\n", blockNumber.String(), err)
				continue
			}

			summary := extractBlockSummary(header)
			history.Store(summary)

			fmt.Printf("  Stored block #%s: %s\n",
				summary.Number.String(), summary.Hash[:16]+"...")
		}

		fmt.Println()

		// Simulate reorg detection
		// BREAKPOINT: Set breakpoint here
		fmt.Println("Simulating reorg detection on next block fetch...")
		fmt.Println()

		// In a real application, you would:
		// 1. Continuously fetch new blocks
		// 2. Check if each new block's parent matches stored history
		// 3. If mismatch detected, reorg occurred
		// 4. Roll back application state and reprocess from reorg point

		fmt.Println("Reorg Handling Steps:")
		fmt.Println()

		fmt.Println("1. Detection Phase:")
		fmt.Println("   - Fetch new block N")
		fmt.Println("   - Check: Does block N's parent hash match stored block N-1?")
		fmt.Println("   - If NO → Reorg detected!")
		fmt.Println()

		fmt.Println("2. Identification Phase:")
		fmt.Println("   - Walk backwards from block N")
		fmt.Println("   - Find the common ancestor (last matching block)")
		fmt.Println("   - Reorg depth = N - common_ancestor_number")
		fmt.Println()

		fmt.Println("3. Rollback Phase:")
		fmt.Println("   - Mark blocks after common ancestor as 'orphaned'")
		fmt.Println("   - Reverse state changes from orphaned blocks")
		fmt.Println("   - Delete or mark invalid any derived data")
		fmt.Println()

		fmt.Println("4. Reprocessing Phase:")
		fmt.Println("   - Fetch new canonical blocks from common ancestor + 1")
		fmt.Println("   - Process new blocks as if seeing them first time")
		fmt.Println("   - Update all derived data and indexes")
		fmt.Println()

		fmt.Println("5. Notification Phase:")
		fmt.Println("   - Log reorg event with depth and affected blocks")
		fmt.Println("   - Alert monitoring systems")
		fmt.Println("   - Notify users if their transactions were affected")
		fmt.Println()

		// Educational explanation
		fmt.Println("Rollback Implementation Patterns:")
		fmt.Println()

		fmt.Println("Pattern 1: Database Transactions")
		fmt.Println("  - Use DB transactions for atomic rollback")
		fmt.Println("  - Keep last N blocks in staging table")
		fmt.Println("  - Only commit to production after confirmations")
		fmt.Println()

		fmt.Println("Pattern 2: Event Sourcing")
		fmt.Println("  - Store all events (block processed, reorg detected)")
		fmt.Println("  - Rebuild state by replaying events")
		fmt.Println("  - Easy to rewind and reprocess")
		fmt.Println()

		fmt.Println("Pattern 3: Versioned State")
		fmt.Println("  - Keep state snapshots at each block height")
		fmt.Println("  - On reorg, revert to snapshot at common ancestor")
		fmt.Println("  - Trade memory for fast rollback")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Production Reorg Handling Best Practices
	// ============================================================================
	// BREAKPOINT: Set breakpoint here

	if demo == "all" {
		fmt.Println("=== Production Reorg Handling Best Practices ===")
		fmt.Println()

		fmt.Println("Detection Strategy:")
		fmt.Println("  ✓ Always check parent hash when processing new blocks")
		fmt.Println("  ✓ Maintain history of recent blocks (at least 100)")
		fmt.Println("  ✓ Monitor for uncle blocks (indicators of forks)")
		fmt.Println("  ✓ Subscribe to reorg events if RPC supports it")
		fmt.Println("  ✓ Use multiple RPC endpoints to detect inconsistencies")
		fmt.Println()

		fmt.Println("Prevention Strategy:")
		fmt.Println("  ✓ Wait for confirmations before trusting data")
		fmt.Println("  ✓ Use higher confirmation counts for high-value operations")
		fmt.Println("  ✓ Mark unconfirmed data as 'pending' in UI")
		fmt.Println("  ✓ Implement confirmation-based workflows")
		fmt.Println()

		fmt.Println("Recovery Strategy:")
		fmt.Println("  ✓ Implement idempotent processing (can reprocess safely)")
		fmt.Println("  ✓ Use database transactions for atomic rollback")
		fmt.Println("  ✓ Keep audit logs of all processing decisions")
		fmt.Println("  ✓ Test rollback procedures regularly")
		fmt.Println("  ✓ Have alerting for deep reorgs (4+ blocks)")
		fmt.Println()

		fmt.Println("Data Management:")
		fmt.Println("  ✓ Soft delete orphaned data (don't hard delete)")
		fmt.Println("  ✓ Track block hash with all derived data")
		fmt.Println("  ✓ Index by both block number AND hash")
		fmt.Println("  ✓ Keep raw block data for reprocessing")
		fmt.Println()

		fmt.Println("User Communication:")
		fmt.Println("  ✓ Show confirmation count in UI")
		fmt.Println("  ✓ Warn about unconfirmed transactions")
		fmt.Println("  ✓ Notify users if their transaction was reorged out")
		fmt.Println("  ✓ Provide transaction status: pending → confirmed → finalized")
		fmt.Println()

		fmt.Println("Monitoring and Alerts:")
		fmt.Println("  ✓ Track reorg frequency and depth")
		fmt.Println("  ✓ Alert on unusual reorg activity")
		fmt.Println("  ✓ Monitor RPC consistency across providers")
		fmt.Println("  ✓ Log all reorg events for analysis")
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
	// 4. Use F11 (Step Into) to enter reorg detection functions
	// 5. Watch Variables panel to see block hash changes
	// 6. Inspect parent hash chains
	// 7. Trace block history through reorganizations
	//
	// REORG DETECTION DEBUGGING:
	// - Set breakpoint in DetectReorg function
	// - Compare newBlock.ParentHash with storedParent.Hash
	// - If they differ, reorg occurred
	// - Walk backwards to find common ancestor
	//
	// TESTING REORGS:
	// - Use testnets (Sepolia) - more frequent reorgs
	// - Monitor during high network congestion
	// - Create local testnet and manually create forks
	// - Use archived node to fetch historical reorgs
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Why do reorgs happen? Network latency, mining randomness
	// - How often? Mainnet: several 1-block reorgs per day
	// - What's an uncle/ommer? Valid block not in canonical chain
	// - What's finality? Point after which reorg is extremely unlikely
	// - How does proof-of-stake change reorgs? Faster finality, explicit finalization
}
