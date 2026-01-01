package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/23-mempool/internal/mempool"
)

/*
===================================================================
Ethereum Mempool (Transaction Pool) Inspection
===================================================================

This program demonstrates how to inspect the mempool (transaction pool)
to understand pending transactions. The mempool contains transactions
that have been broadcast but not yet included in blocks.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <DEMO>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, basic, congestion, monitoring (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com basic
    go run ./cmd/app/main.go https://eth.llamarpc.com congestion
    go run ./cmd/app/main.go https://rpc.sepolia.org monitoring

Common RPC Endpoints:
    Ethereum Mainnet:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth
      - https://eth-mainnet.public.blastapi.io

    Sepolia Testnet:
      - https://rpc.sepolia.org
      - https://sepolia.gateway.tenderly.co

NOTE: Most public RPC endpoints restrict mempool access for privacy and
security reasons. Your own node will provide full mempool visibility.

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter mempool package functions
- Watch Variables panel to see pending transaction counts
- Inspect PendingCount and congestion analysis

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the mempool package from internal/mempool
- exercise.go, solution.reference.go are in package mempool
- The mempool.Run() function queries pending transaction count

KEY CONCEPTS:
- Mempool/TxPool: Queue of pending transactions waiting for inclusion
- MEV (Maximal Extractable Value): Profit from reordering transactions
- Gas Price Markets: Supply/demand dynamics for block space
- Transaction Lifecycle: Broadcast → Mempool → Block → Confirmation
- Privacy vs Transparency: Trade-offs in mempool visibility
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
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com basic\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com congestion\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, basic, congestion, monitoring\n")
		fmt.Fprintf(os.Stderr, "\nNOTE: Most public RPCs restrict mempool access\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Mempool (Transaction Pool) Inspection ===\n")
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
	// the MempoolClient interface. It handles JSON-RPC communication with Ethereum nodes.

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
	// STEP 3: Example 1 - Basic Mempool Size Query
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through basic mempool query
	// DEBUG: This demonstrates the simplest use case—how many pending txs?
	// DEBUG: Watch the 'result' variable to see PendingCount

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Mempool Size Query ===")
		fmt.Println("Retrieving number of pending transactions...")
		fmt.Println()

		// Create basic configuration (no special settings needed)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg has Limit set to 0 (no limit)
		cfg := mempool.Config{
			Limit: 0, // 0 = no limit, return total count
		}

		// BREAKPOINT: Set breakpoint here before calling mempool.Run
		// DEBUG: Step Into (F11) to see the implementation in exercise.go or solution.reference.go
		// DEBUG: You'll see how PendingTransactionCount() RPC is called
		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := mempool.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query mempool: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful retrieval
		// DEBUG: Expand 'result' in Variables panel to inspect PendingCount
		// DEBUG: The count represents transactions waiting for block inclusion

		fmt.Printf("Pending Transactions: %d\n", result.PendingCount)
		fmt.Println()

		// Provide basic congestion assessment
		// BREAKPOINT: Set breakpoint here to see congestion interpretation
		// DEBUG: Watch how different pending counts indicate network activity
		if result.PendingCount == 0 {
			fmt.Println("✓ No pending transactions")
			fmt.Println("  Network is idle or mempool access is restricted")
		} else if result.PendingCount < 1000 {
			fmt.Println("✓ Low mempool size (< 1,000)")
			fmt.Println("  Minimal network congestion")
			fmt.Println("  Transactions confirm quickly with base fees")
		} else if result.PendingCount < 10000 {
			fmt.Println("⚠️  Moderate mempool size (1,000-10,000)")
			fmt.Println("  Normal network activity")
			fmt.Println("  Standard gas prices should confirm in reasonable time")
		} else if result.PendingCount < 50000 {
			fmt.Println("⚠️  High mempool size (10,000-50,000)")
			fmt.Println("  Network congestion present")
			fmt.Println("  May need higher gas prices for timely confirmation")
		} else {
			fmt.Println("❌ Extreme mempool size (> 50,000)")
			fmt.Println("  Severe network congestion")
			fmt.Println("  Only high-fee transactions confirm quickly")
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: You've now seen the basic mempool size and assessment
	}

	// ============================================================================
	// STEP 4: Example 2 - Network Congestion Analysis
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through congestion analysis
	// DEBUG: This shows how to interpret mempool size for gas pricing
	// DEBUG: Useful for understanding when to send transactions

	if demo == "all" || demo == "congestion" {
		fmt.Println("=== Example 2: Network Congestion Analysis ===")
		fmt.Println("Analyzing mempool size to assess network congestion...")
		fmt.Println()

		cfg := mempool.Config{Limit: 0}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling mempool.Run
		// DEBUG: Step Into (F11) to see PendingTransactionCount implementation
		result, err := mempool.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query mempool: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect result.PendingCount for congestion analysis

		fmt.Printf("Current Mempool Size: %d pending transactions\n\n", result.PendingCount)

		// Detailed congestion analysis
		fmt.Println("Congestion Analysis:")
		fmt.Println()

		// Categorize congestion level
		congestionLevel := "Unknown"
		gasStrategy := ""

		if result.PendingCount == 0 {
			congestionLevel = "IDLE"
			gasStrategy = "Use minimum gas price (may be restricted RPC endpoint)"
		} else if result.PendingCount < 1000 {
			congestionLevel = "LOW"
			gasStrategy = "Base fee + small priority fee (1-2 gwei)"
		} else if result.PendingCount < 5000 {
			congestionLevel = "MODERATE"
			gasStrategy = "Base fee + standard priority fee (2-5 gwei)"
		} else if result.PendingCount < 20000 {
			congestionLevel = "HIGH"
			gasStrategy = "Base fee + elevated priority fee (5-10 gwei)"
		} else if result.PendingCount < 50000 {
			congestionLevel = "VERY HIGH"
			gasStrategy = "Base fee + high priority fee (10-20 gwei)"
		} else {
			congestionLevel = "EXTREME"
			gasStrategy = "Base fee + very high priority fee (20+ gwei) or wait"
		}

		fmt.Printf("Congestion Level: %s\n", congestionLevel)
		fmt.Printf("Recommended Gas Strategy: %s\n\n", gasStrategy)

		// Explain implications
		fmt.Println("What This Means:")
		fmt.Println()

		if result.PendingCount < 1000 {
			fmt.Println("  Transaction Confirmation:")
			fmt.Println("    - Expected confirmation: 1-2 blocks (~15-30 seconds)")
			fmt.Println("    - Base fee likely decreasing or stable")
			fmt.Println("    - Good time to send non-urgent transactions")
		} else if result.PendingCount < 10000 {
			fmt.Println("  Transaction Confirmation:")
			fmt.Println("    - Expected confirmation: 2-5 blocks (~30-75 seconds)")
			fmt.Println("    - Base fee may be rising slowly")
			fmt.Println("    - Standard transactions confirm normally")
		} else if result.PendingCount < 50000 {
			fmt.Println("  Transaction Confirmation:")
			fmt.Println("    - Expected confirmation: 5-20 blocks (~1-4 minutes)")
			fmt.Println("    - Base fee likely rising")
			fmt.Println("    - Need higher priority fee for faster confirmation")
			fmt.Println("    - Consider waiting if transaction is not urgent")
		} else {
			fmt.Println("  Transaction Confirmation:")
			fmt.Println("    - Expected confirmation: Highly variable (minutes to hours)")
			fmt.Println("    - Base fee rising rapidly")
			fmt.Println("    - Only high-fee transactions confirm quickly")
			fmt.Println("    - Strongly recommend waiting unless critical")
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: You've now seen comprehensive congestion analysis
	}

	// ============================================================================
	// STEP 5: Example 3 - Mempool Monitoring Over Time
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to see monitoring pattern
	// DEBUG: This demonstrates tracking mempool size trends

	if demo == "all" || demo == "monitoring" {
		fmt.Println("=== Example 3: Mempool Monitoring Over Time ===")
		fmt.Println("Tracking mempool size to detect congestion trends...")
		fmt.Println()

		// Simulate monitoring loop (production would run continuously)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This pattern helps detect mempool growth/shrinkage
		maxChecks := 5
		mempoolHistory := make([]uint, 0, maxChecks)

		for i := 0; i < maxChecks; i++ {
			cfg := mempool.Config{Limit: 0}

			runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
			result, err := mempool.Run(runCtx, client, cfg)
			runCancel()

			if err != nil {
				fmt.Printf("Check %d/%d: ERROR - %v\n", i+1, maxChecks, err)
				time.Sleep(3 * time.Second)
				continue
			}

			// BREAKPOINT: Set breakpoint here in monitoring loop
			// DEBUG: Watch how mempool size changes over time
			mempoolHistory = append(mempoolHistory, result.PendingCount)

			fmt.Printf("Check %d/%d: ", i+1, maxChecks)
			fmt.Printf("%d pending txs", result.PendingCount)

			// Detect changes from previous check
			if i > 0 {
				prevCount := mempoolHistory[i-1]
				diff := int(result.PendingCount) - int(prevCount)
				if diff > 0 {
					fmt.Printf(" (+%d, GROWING)", diff)
				} else if diff < 0 {
					fmt.Printf(" (%d, SHRINKING)", diff)
				} else {
					fmt.Printf(" (no change)")
				}
			}
			fmt.Println()

			time.Sleep(3 * time.Second)
		}

		fmt.Println()

		// Analyze mempool trend
		// BREAKPOINT: Set breakpoint here to see trend analysis
		// DEBUG: Watch how we detect mempool growth/shrinkage patterns
		if len(mempoolHistory) > 1 {
			firstCount := mempoolHistory[0]
			lastCount := mempoolHistory[len(mempoolHistory)-1]
			totalChange := int(lastCount) - int(firstCount)

			fmt.Println("Mempool Trend Analysis:")
			fmt.Printf("  Initial size: %d txs\n", firstCount)
			fmt.Printf("  Final size:   %d txs\n", lastCount)
			fmt.Printf("  Net change:   %+d txs\n\n", totalChange)

			// Determine trend
			if totalChange > 100 {
				fmt.Println("  Trend: GROWING (congestion increasing)")
				fmt.Println("  Recommendation: Gas prices likely rising, consider waiting")
			} else if totalChange < -100 {
				fmt.Println("  Trend: SHRINKING (congestion decreasing)")
				fmt.Println("  Recommendation: Good time to send transactions")
			} else {
				fmt.Println("  Trend: STABLE (congestion steady)")
				fmt.Println("  Recommendation: Current gas strategy should work")
			}
		}

		fmt.Println()

		// Display production monitoring best practices
		fmt.Println("Production Monitoring Best Practices:")
		fmt.Println("  1. Poll mempool size every 30-60 seconds")
		fmt.Println("  2. Track mempool size trends (growing vs shrinking)")
		fmt.Println("  3. Correlate with gas price recommendations")
		fmt.Println("  4. Alert on extreme mempool sizes (> 50,000)")
		fmt.Println("  5. Use mempool data to optimize transaction timing")
		fmt.Println("  6. Remember: Public RPCs often hide mempool data")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Understanding Mempools and Transaction Lifecycle
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to learn mempool concepts
	// DEBUG: This section explains how mempools work

	if demo == "all" {
		fmt.Println("=== Understanding Mempools and Transaction Lifecycle ===")
		fmt.Println()

		fmt.Println("Transaction Lifecycle:")
		fmt.Println("  1. Creation: User creates and signs transaction")
		fmt.Println("  2. Broadcast: Transaction sent to network via RPC")
		fmt.Println("  3. Propagation: Gossipped to peers via P2P network")
		fmt.Println("  4. Validation: Nodes verify signature, nonce, balance")
		fmt.Println("  5. Mempool: Valid transactions stored in local pool")
		fmt.Println("  6. Selection: Miner/validator picks txs for block")
		fmt.Println("  7. Inclusion: Transaction included in mined block")
		fmt.Println("  8. Confirmation: Block finalized on chain")
		fmt.Println()

		fmt.Println("Mempool Management:")
		fmt.Println("  Priority Sorting:")
		fmt.Println("    - Transactions sorted by gas price (fee/gas)")
		fmt.Println("    - Higher price = higher priority for inclusion")
		fmt.Println("    - Miners maximize revenue by selecting high-fee txs")
		fmt.Println()
		fmt.Println("  Size Limits:")
		fmt.Println("    - Default Geth limit: ~4GB or ~1M transactions")
		fmt.Println("    - When full, lowest-fee txs are evicted")
		fmt.Println("    - Each node can set custom limits")
		fmt.Println()
		fmt.Println("  Eviction Rules:")
		fmt.Println("    - Lowest gas price transactions evicted first")
		fmt.Println("    - Old transactions (> 3 hours) may be dropped")
		fmt.Println("    - Invalid txs (bad nonce, insufficient balance) rejected")
		fmt.Println()

		fmt.Println("MEV and Mempool Privacy:")
		fmt.Println("  What is MEV?")
		fmt.Println("    - Maximal Extractable Value: profit from tx ordering")
		fmt.Println("    - Searchers scan mempool for profitable opportunities")
		fmt.Println("    - Examples: arbitrage, liquidations, sandwich attacks")
		fmt.Println()
		fmt.Println("  Privacy Trade-offs:")
		fmt.Println("    - Public mempool: Transparency but enables front-running")
		fmt.Println("    - Private mempools: Reduced MEV but less transparency")
		fmt.Println("    - Flashbots: Private tx submission to avoid front-running")
		fmt.Println()

		fmt.Println("Why Public RPCs Hide Mempool Data:")
		fmt.Println("  1. Privacy: Protect users from MEV exploitation")
		fmt.Println("  2. Resources: Returning thousands of txs is expensive")
		fmt.Println("  3. Security: Full mempool dumps could enable DoS")
		fmt.Println("  4. Business: Some operators monetize mempool access")
		fmt.Println()

		fmt.Println("Transaction Replacement (RBF):")
		fmt.Println("  - Replace by Fee: Submit new tx with same nonce")
		fmt.Println("  - Requires 10%+ higher gas price (configurable)")
		fmt.Println("  - Useful for stuck transactions (fee too low)")
		fmt.Println("  - Can cancel by sending 0 ETH to self with higher fee")
		fmt.Println()

		fmt.Println("Available RPC Methods:")
		fmt.Println("  Public (sometimes available):")
		fmt.Println("    - eth_getTransactionByHash: Check specific tx")
		fmt.Println("    - Basic mempool stats (if not restricted)")
		fmt.Println()
		fmt.Println("  Restricted (rarely on public RPCs):")
		fmt.Println("    - txpool_content: Full mempool contents")
		fmt.Println("    - txpool_status: Detailed mempool stats")
		fmt.Println("    - txpool_inspect: Transaction summaries")
		fmt.Println("    - eth_pendingTransactions: All pending txs")
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
	// 4. Use F11 (Step Into) to enter mempool package functions
	// 5. Watch Variables panel to see pending transaction counts
	// 6. Inspect result.PendingCount over time
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Mempool/TxPool: How do pending transactions queue?
	// - Gas Price Markets: How does supply/demand affect fees?
	// - MEV: What is Maximal Extractable Value?
	// - Transaction Lifecycle: From creation to confirmation
	// - Privacy Trade-offs: Why hide mempool data?
	//
	// STEPPING INTO MEMPOOL PACKAGE:
	// - Set breakpoint at mempool.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see PendingTransactionCount() RPC call
	// - Watch how uint is returned (no copying needed)
	// - See error handling for restricted endpoints
	//
	// COMMON MEMPOOL ISSUES TO DEBUG:
	// - Zero pending txs: Likely public RPC with restricted access
	// - Error querying mempool: Endpoint doesn't support method
	// - High pending count: Network congestion, plan gas accordingly
	// - Mempool growing rapidly: Congestion increasing, fees rising
}
