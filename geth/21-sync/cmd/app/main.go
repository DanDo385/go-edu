package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/21-sync/internal/sync"
)

/*
===================================================================
Ethereum Node Sync Progress Monitoring
===================================================================

This program demonstrates how to check if an Ethereum node is fully
synced by inspecting its sync progress. This is critical for production
systems because querying a syncing node returns stale data.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <DEMO>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, basic, detailed, monitoring (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com basic
    go run ./cmd/app/main.go https://eth.llamarpc.com detailed
    go run ./cmd/app/main.go https://rpc.sepolia.org monitoring

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
- Step Into (F11) to enter sync package functions
- Watch Variables panel to see sync progress structures
- Inspect IsSyncing boolean and Progress details

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the sync package from internal/sync
- exercise.go, solution.reference.go are in package sync
- The sync.Run() function checks node sync status

KEY CONCEPTS:
- Sync Progress: Tracks how far a node is in blockchain sync
- Nil as Sentinel: nil Progress = fully synced, non-nil = syncing
- Sync Modes: Full sync, Snap sync, Light sync
- State Healing: Snap sync downloads state then fills gaps
- Production Readiness: Only query fully synced nodes for accurate data
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
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com detailed\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, basic, detailed, monitoring\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Node Sync Progress Monitoring ===\n")
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
	// the SyncClient interface. It handles JSON-RPC communication with Ethereum nodes.

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
	// STEP 3: Example 1 - Basic Sync Status Check
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through basic sync check
	// DEBUG: This demonstrates the simplest use case—is the node synced?
	// DEBUG: Watch the 'result' variable to see IsSyncing and Progress

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Sync Status Check ===")
		fmt.Println("Checking if the node is fully synced...")
		fmt.Println()

		// Create basic configuration (no special settings needed)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg is empty—sync check requires no configuration
		cfg := sync.Config{}

		// BREAKPOINT: Set breakpoint here before calling sync.Run
		// DEBUG: Step Into (F11) to see the implementation in exercise.go or solution.reference.go
		// DEBUG: You'll see how SyncProgress() is called and nil is checked
		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := sync.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to check sync status: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful retrieval
		// DEBUG: Expand 'result' in Variables panel to inspect fields
		// DEBUG: Look at IsSyncing (true = syncing, false = fully synced)
		// DEBUG: If IsSyncing is true, Progress will be non-nil

		fmt.Printf("Node Sync Status:\n")
		if result.IsSyncing {
			fmt.Printf("  Status: SYNCING (not ready for production queries)\n")
			fmt.Printf("  The node is still catching up to the network\n")
		} else {
			fmt.Printf("  Status: FULLY SYNCED (ready for production queries)\n")
			fmt.Printf("  The node is up-to-date with the network\n")
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Now you can see the basic sync status
	}

	// ============================================================================
	// STEP 4: Example 2 - Detailed Sync Progress Information
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through detailed progress
	// DEBUG: This shows how to inspect sync progress when a node is syncing
	// DEBUG: Useful for monitoring sync operations and estimating completion time

	if demo == "all" || demo == "detailed" {
		fmt.Println("=== Example 2: Detailed Sync Progress Information ===")
		fmt.Println("Retrieving detailed sync progress metrics...")
		fmt.Println()

		cfg := sync.Config{}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling sync.Run
		// DEBUG: Step Into (F11) to see how Progress is populated
		result, err := sync.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to check sync progress: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Expand 'result.Progress' to see all sync metrics
		// DEBUG: StartingBlock, CurrentBlock, HighestBlock show progress
		// DEBUG: PulledStates/KnownStates show state trie sync (snap sync)

		fmt.Printf("Sync Status: %s\n\n", map[bool]string{true: "SYNCING", false: "FULLY SYNCED"}[result.IsSyncing])

		if result.Progress != nil {
			fmt.Println("Detailed Sync Progress:")
			fmt.Printf("  Starting Block: %d\n", result.Progress.StartingBlock)
			fmt.Printf("  Current Block:  %d\n", result.Progress.CurrentBlock)
			fmt.Printf("  Highest Block:  %d\n", result.Progress.HighestBlock)
			fmt.Println()

			// Calculate sync percentage
			// BREAKPOINT: Set breakpoint here to see percentage calculation
			// DEBUG: Watch how we compute progress from block numbers
			if result.Progress.HighestBlock > result.Progress.StartingBlock {
				total := result.Progress.HighestBlock - result.Progress.StartingBlock
				current := result.Progress.CurrentBlock - result.Progress.StartingBlock
				percentage := float64(current) / float64(total) * 100
				fmt.Printf("  Progress: %.2f%% complete\n", percentage)
				fmt.Printf("  Remaining: %d blocks\n", result.Progress.HighestBlock-result.Progress.CurrentBlock)
			}
			fmt.Println()

			// Display state sync progress (for snap sync)
			// BREAKPOINT: Set breakpoint here
			// DEBUG: State sync is specific to snap/fast sync modes
			if result.Progress.KnownStates > 0 {
				fmt.Println("State Sync Progress (Snap Sync):")
				fmt.Printf("  Pulled States: %d\n", result.Progress.PulledStates)
				fmt.Printf("  Known States:  %d\n", result.Progress.KnownStates)
				statePercentage := float64(result.Progress.PulledStates) / float64(result.Progress.KnownStates) * 100
				fmt.Printf("  Progress:      %.2f%%\n", statePercentage)
				fmt.Println()
			}
		} else {
			fmt.Println("The node is fully synced—no progress information available.")
			fmt.Println("This is the normal state for a production node.")
			fmt.Println()
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: You've now seen all available sync progress details
	}

	// ============================================================================
	// STEP 5: Example 3 - Production Monitoring Pattern
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to see production monitoring pattern
	// DEBUG: This demonstrates how to use sync status in production systems

	if demo == "all" || demo == "monitoring" {
		fmt.Println("=== Example 3: Production Monitoring Pattern ===")
		fmt.Println("Demonstrating how to monitor node sync status for production...")
		fmt.Println()

		// Simulate a monitoring loop (in production, this would run continuously)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This pattern is what you'd use in a health check endpoint
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			cfg := sync.Config{}

			runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
			result, err := sync.Run(runCtx, client, cfg)
			runCancel()

			if err != nil {
				fmt.Printf("Check %d/%d: ERROR - %v\n", i+1, maxRetries, err)
				time.Sleep(2 * time.Second)
				continue
			}

			// BREAKPOINT: Set breakpoint here in the monitoring loop
			// DEBUG: Watch how sync status is evaluated for readiness
			fmt.Printf("Check %d/%d:\n", i+1, maxRetries)
			if result.IsSyncing {
				fmt.Printf("  ❌ Node is SYNCING - NOT ready for production\n")
				if result.Progress != nil {
					blocksRemaining := result.Progress.HighestBlock - result.Progress.CurrentBlock
					fmt.Printf("  ⏳ Blocks remaining: %d\n", blocksRemaining)
				}
			} else {
				fmt.Printf("  ✓ Node is FULLY SYNCED - ready for production\n")
				fmt.Printf("  ✓ Queries will return current data\n")
			}
			fmt.Println()

			time.Sleep(2 * time.Second)
		}

		// Display production best practices
		fmt.Println("Production Monitoring Best Practices:")
		fmt.Println("  1. Check sync status before accepting traffic")
		fmt.Println("  2. Poll sync status every 30-60 seconds")
		fmt.Println("  3. Alert if node stays syncing for too long")
		fmt.Println("  4. Implement graceful degradation (fallback to other nodes)")
		fmt.Println("  5. Log sync progress for debugging sync issues")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Understanding Sync Modes and Their Trade-offs
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to learn about sync modes
	// DEBUG: This section explains different sync strategies

	if demo == "all" {
		fmt.Println("=== Understanding Ethereum Sync Modes ===")
		fmt.Println()

		fmt.Println("Full Sync:")
		fmt.Println("  - Downloads all blocks and executes every transaction")
		fmt.Println("  - Rebuilds complete state from genesis")
		fmt.Println("  - Slowest method (~weeks for mainnet)")
		fmt.Println("  - Highest trust and verification")
		fmt.Println("  - Required for archive nodes")
		fmt.Println()

		fmt.Println("Snap Sync (Snapshot Sync):")
		fmt.Println("  - Downloads state snapshots from peers")
		fmt.Println("  - Syncs recent blocks in reverse order")
		fmt.Println("  - Heals missing state trie nodes")
		fmt.Println("  - Fastest method (~hours for mainnet)")
		fmt.Println("  - Default in modern Geth versions")
		fmt.Println("  - Progress shown via PulledStates/KnownStates")
		fmt.Println()

		fmt.Println("Light Sync:")
		fmt.Println("  - Downloads only block headers")
		fmt.Println("  - Requests proofs for specific data on demand")
		fmt.Println("  - Minimal storage (~few GB)")
		fmt.Println("  - Relies on full nodes for data")
		fmt.Println("  - Good for mobile/resource-constrained devices")
		fmt.Println()

		fmt.Println("SyncProgress Fields Explained:")
		fmt.Println("  - StartingBlock: Block number when sync began")
		fmt.Println("  - CurrentBlock: Latest block we've downloaded")
		fmt.Println("  - HighestBlock: Latest block known to exist on network")
		fmt.Println("  - PulledStates: State trie entries downloaded (snap sync)")
		fmt.Println("  - KnownStates: Total state trie entries needed (snap sync)")
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
	// 4. Use F11 (Step Into) to enter sync package functions
	// 5. Watch Variables panel to see sync progress structures
	// 6. Inspect result.IsSyncing and result.Progress
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Sync Progress: How does a node catch up to the network?
	// - Nil as Sentinel: Why does nil mean "fully synced"?
	// - Sync Modes: What are the trade-offs between speed and trust?
	// - State Healing: How does snap sync fill in missing data?
	// - Production Readiness: When is a node safe to query?
	//
	// STEPPING INTO SYNC PACKAGE:
	// - Set breakpoint at sync.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see SyncProgress() RPC call
	// - Watch how nil check determines IsSyncing status
	// - See how Progress is safely returned without copying
	//
	// COMMON SYNC ISSUES TO DEBUG:
	// - Node stuck syncing: Check peer count, network connectivity
	// - Sync progress not advancing: Check disk space, CPU/RAM
	// - SyncProgress returns nil on syncing node: RPC client issue
	// - High CurrentBlock but still syncing: State healing in progress
}
