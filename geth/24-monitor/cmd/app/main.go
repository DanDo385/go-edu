package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/24-monitor/internal/monitor"
)

/*
===================================================================
Ethereum Node Health Monitoring and Staleness Detection
===================================================================

This program demonstrates how to monitor Ethereum node health by
checking block freshness. A stale node (lagging behind) returns
outdated data and shouldn't be used for production queries.

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
- Step Into (F11) to enter monitor package functions
- Watch Variables panel to see health status and lag metrics
- Inspect Status, LagSeconds, and BlockTimestamp fields

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the monitor package from internal/monitor
- exercise.go, solution.reference.go are in package monitor
- The monitor.Run() function performs health checks

KEY CONCEPTS:
- Block Freshness: How recent is the latest block?
- Staleness Detection: Is the node lagging behind the network?
- Time-based Health Checks: Using timestamps for monitoring
- Production Readiness: When is a node safe to query?
- Threshold-based Alerting: OK vs STALE classification
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
	fmt.Printf("=== Ethereum Node Health Monitoring ===\n")
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
	// the MonitorClient interface. It handles JSON-RPC communication with Ethereum nodes.

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
	// STEP 3: Example 1 - Basic Node Health Check
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through basic health check
	// DEBUG: This demonstrates the simplest use case—is the node fresh?
	// DEBUG: Watch the 'result' variable to see Status and LagSeconds

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Node Health Check ===")
		fmt.Println("Checking if node is providing fresh data...")
		fmt.Println()

		// Create configuration with default lag threshold (60 seconds)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: MaxLagSeconds defaults to 60 (about 5 blocks on mainnet)
		cfg := monitor.Config{
			MaxLagSeconds: 60,   // 60 seconds = ~5 blocks on mainnet
			BlockNumber:   nil,  // nil = check latest block
		}

		// BREAKPOINT: Set breakpoint here before calling monitor.Run
		// DEBUG: Step Into (F11) to see the implementation in exercise.go or solution.reference.go
		// DEBUG: You'll see how block timestamp is compared to current time
		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := monitor.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to perform health check: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful retrieval
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at Status (OK vs STALE), LagSeconds, BlockTimestamp

		fmt.Printf("Health Check Results:\n")
		fmt.Printf("  Status:          %s\n", result.Status)
		fmt.Printf("  Block Number:    %d\n", result.BlockNumber)
		fmt.Printf("  Block Timestamp: %s\n", result.BlockTimestamp.Format(time.RFC3339))
		fmt.Printf("  Block Age:       %d seconds\n", result.LagSeconds)
		fmt.Println()

		// Interpret results
		// BREAKPOINT: Set breakpoint here to see health interpretation
		// DEBUG: Watch how Status determines node usability
		if result.Status == "OK" {
			fmt.Println("✓ Node is HEALTHY")
			fmt.Println("  The node is up-to-date and providing fresh data")
			fmt.Println("  Safe to use for production queries")
		} else {
			fmt.Println("❌ Node is STALE")
			fmt.Println("  The node is lagging behind the network")
			fmt.Println("  Data may be outdated—use a different endpoint")
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: You've now seen basic health status
	}

	// ============================================================================
	// STEP 4: Example 2 - Detailed Health Analysis
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through detailed analysis
	// DEBUG: This shows comprehensive health metrics and interpretation
	// DEBUG: Useful for understanding what different lag times mean

	if demo == "all" || demo == "detailed" {
		fmt.Println("=== Example 2: Detailed Health Analysis ===")
		fmt.Println("Performing comprehensive node health assessment...")
		fmt.Println()

		cfg := monitor.Config{
			MaxLagSeconds: 60,
			BlockNumber:   nil,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling monitor.Run
		// DEBUG: Step Into (F11) to see lag calculation
		result, err := monitor.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to perform health check: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect result fields for detailed analysis

		fmt.Println("Detailed Health Metrics:")
		fmt.Printf("  Block Number:    #%d\n", result.BlockNumber)
		fmt.Printf("  Block Time:      %s\n", result.BlockTimestamp.Format("2006-01-02 15:04:05 MST"))
		fmt.Printf("  Current Time:    %s\n", time.Now().Format("2006-01-02 15:04:05 MST"))
		fmt.Printf("  Lag:             %d seconds\n", result.LagSeconds)
		fmt.Printf("  Status:          %s\n", result.Status)
		fmt.Println()

		// Calculate block lag in terms of expected blocks
		// BREAKPOINT: Set breakpoint here to see block lag calculation
		// DEBUG: Watch how we estimate how many blocks behind the node is
		expectedBlockTime := 12 // Ethereum mainnet ~12 second block time
		estimatedBlocksBehind := result.LagSeconds / int64(expectedBlockTime)

		fmt.Println("Lag Interpretation:")
		if result.LagSeconds < 0 {
			fmt.Println("  Lag Category: FUTURE BLOCK")
			fmt.Println("  This is usually due to clock skew—generally harmless")
			fmt.Printf("  Block timestamp is %d seconds in the future\n", -result.LagSeconds)
		} else if result.LagSeconds <= 15 {
			fmt.Println("  Lag Category: EXCELLENT (<15s)")
			fmt.Println("  Node is fully synced and current")
			fmt.Printf("  Approximately %d blocks behind (within normal variance)\n", estimatedBlocksBehind)
		} else if result.LagSeconds <= 30 {
			fmt.Println("  Lag Category: GOOD (15-30s)")
			fmt.Println("  Node is current, minor lag acceptable")
			fmt.Printf("  Approximately %d blocks behind\n", estimatedBlocksBehind)
		} else if result.LagSeconds <= 60 {
			fmt.Println("  Lag Category: ACCEPTABLE (30-60s)")
			fmt.Println("  Node is slightly behind but still usable")
			fmt.Printf("  Approximately %d blocks behind\n", estimatedBlocksBehind)
		} else if result.LagSeconds <= 180 {
			fmt.Println("  Lag Category: CONCERNING (1-3 minutes)")
			fmt.Println("  Node is lagging—data may be outdated")
			fmt.Printf("  Approximately %d blocks behind\n", estimatedBlocksBehind)
		} else {
			fmt.Println("  Lag Category: CRITICAL (>3 minutes)")
			fmt.Println("  Node is significantly behind the network")
			fmt.Printf("  Approximately %d blocks behind\n", estimatedBlocksBehind)
			fmt.Println("  Do NOT use for production—find another endpoint")
		}
		fmt.Println()

		// Impact analysis
		// BREAKPOINT: Set breakpoint here to see impact analysis
		// DEBUG: This explains what lag means for different use cases
		fmt.Println("Impact on Different Use Cases:")
		fmt.Println()

		if result.LagSeconds <= 30 {
			fmt.Println("  ✓ Block explorers: Excellent")
			fmt.Println("  ✓ Wallets: Excellent")
			fmt.Println("  ✓ DeFi: Excellent")
			fmt.Println("  ✓ Trading bots: Excellent")
		} else if result.LagSeconds <= 60 {
			fmt.Println("  ✓ Block explorers: Good")
			fmt.Println("  ✓ Wallets: Acceptable")
			fmt.Println("  ⚠️  DeFi: Risky (prices may be stale)")
			fmt.Println("  ❌ Trading bots: Not recommended")
		} else {
			fmt.Println("  ⚠️  Block explorers: Delayed data")
			fmt.Println("  ❌ Wallets: May show incorrect balances")
			fmt.Println("  ❌ DeFi: Dangerous (stale prices, failed txs)")
			fmt.Println("  ❌ Trading bots: Do not use")
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: You've now seen comprehensive health analysis
	}

	// ============================================================================
	// STEP 5: Example 3 - Continuous Health Monitoring
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to see monitoring pattern
	// DEBUG: This demonstrates how to monitor node health over time

	if demo == "all" || demo == "monitoring" {
		fmt.Println("=== Example 3: Continuous Health Monitoring ===")
		fmt.Println("Monitoring node health over time to detect degradation...")
		fmt.Println()

		// Simulate monitoring loop (production would run continuously)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This pattern is essential for production monitoring
		maxChecks := 5
		lagHistory := make([]int64, 0, maxChecks)
		statusHistory := make([]string, 0, maxChecks)

		for i := 0; i < maxChecks; i++ {
			cfg := monitor.Config{
				MaxLagSeconds: 60,
				BlockNumber:   nil,
			}

			runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
			result, err := monitor.Run(runCtx, client, cfg)
			runCancel()

			if err != nil {
				fmt.Printf("Check %d/%d: ERROR - %v\n", i+1, maxChecks, err)
				time.Sleep(3 * time.Second)
				continue
			}

			// BREAKPOINT: Set breakpoint here in monitoring loop
			// DEBUG: Watch how health metrics change over time
			lagHistory = append(lagHistory, result.LagSeconds)
			statusHistory = append(statusHistory, result.Status)

			fmt.Printf("Check %d/%d: ", i+1, maxChecks)
			fmt.Printf("Block #%d, Lag: %ds, Status: %s",
				result.BlockNumber, result.LagSeconds, result.Status)

			// Detect status change from previous check
			if i > 0 {
				prevStatus := statusHistory[i-1]
				if result.Status != prevStatus {
					if result.Status == "STALE" {
						fmt.Printf(" ⚠️  DEGRADED (was %s)", prevStatus)
					} else {
						fmt.Printf(" ✓ RECOVERED (was %s)", prevStatus)
					}
				}
			}
			fmt.Println()

			time.Sleep(3 * time.Second)
		}

		fmt.Println()

		// Analyze health trend
		// BREAKPOINT: Set breakpoint here to see trend analysis
		// DEBUG: Watch how we detect improving/degrading health
		if len(lagHistory) > 1 {
			var totalLag int64
			minLag := lagHistory[0]
			maxLag := lagHistory[0]
			staleCount := 0

			for i, lag := range lagHistory {
				totalLag += lag
				if lag < minLag {
					minLag = lag
				}
				if lag > maxLag {
					maxLag = lag
				}
				if statusHistory[i] == "STALE" {
					staleCount++
				}
			}

			avgLag := totalLag / int64(len(lagHistory))

			fmt.Println("Health Statistics:")
			fmt.Printf("  Average Lag:    %d seconds\n", avgLag)
			fmt.Printf("  Min Lag:        %d seconds\n", minLag)
			fmt.Printf("  Max Lag:        %d seconds\n", maxLag)
			fmt.Printf("  Stale Checks:   %d/%d (%.1f%%)\n",
				staleCount, len(statusHistory),
				float64(staleCount)/float64(len(statusHistory))*100)
			fmt.Println()

			// Health assessment
			if staleCount == 0 && avgLag <= 30 {
				fmt.Println("  Overall Health: EXCELLENT")
				fmt.Println("  Node is consistently healthy")
			} else if staleCount == 0 && avgLag <= 60 {
				fmt.Println("  Overall Health: GOOD")
				fmt.Println("  Node is healthy with acceptable lag")
			} else if staleCount <= len(statusHistory)/2 {
				fmt.Println("  Overall Health: DEGRADED")
				fmt.Println("  Node is intermittently stale—monitor closely")
			} else {
				fmt.Println("  Overall Health: POOR")
				fmt.Println("  Node is frequently stale—consider switching endpoints")
			}
		}

		fmt.Println()

		// Display production monitoring best practices
		fmt.Println("Production Monitoring Best Practices:")
		fmt.Println("  1. Check health every 30-60 seconds")
		fmt.Println("  2. Alert if status is STALE for 3+ consecutive checks")
		fmt.Println("  3. Track lag trends over time (improving vs degrading)")
		fmt.Println("  4. Implement automatic failover to backup endpoints")
		fmt.Println("  5. Correlate with sync status and peer count")
		fmt.Println("  6. Set different lag thresholds for different use cases")
		fmt.Println("  7. Log all health checks for post-incident analysis")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Understanding Block Timestamps and Health Monitoring
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to learn monitoring concepts
	// DEBUG: This section explains how timestamp-based monitoring works

	if demo == "all" {
		fmt.Println("=== Understanding Block Timestamps and Health Monitoring ===")
		fmt.Println()

		fmt.Println("How Block Timestamps Work:")
		fmt.Println("  - Set by miner/validator when creating block")
		fmt.Println("  - Unix timestamp (seconds since Jan 1, 1970)")
		fmt.Println("  - Must be greater than parent block timestamp")
		fmt.Println("  - Network rejects blocks with future timestamps (> 15s)")
		fmt.Println("  - Can be slightly in the past due to clock skew")
		fmt.Println()

		fmt.Println("Why Timestamp-Based Health Checks:")
		fmt.Println("  - Simple and effective staleness detection")
		fmt.Println("  - Doesn't require external time source (uses block time)")
		fmt.Println("  - Works across all networks (mainnet, testnets, L2s)")
		fmt.Println("  - Catches both sync issues and slow RPC endpoints")
		fmt.Println()

		fmt.Println("Typical Block Times by Network:")
		fmt.Println("  - Ethereum Mainnet: ~12 seconds")
		fmt.Println("  - Polygon: ~2 seconds")
		fmt.Println("  - Optimism: ~2 seconds (L2)")
		fmt.Println("  - Arbitrum: ~0.25 seconds (L2)")
		fmt.Println("  - Sepolia Testnet: ~12 seconds")
		fmt.Println()

		fmt.Println("Recommended Lag Thresholds by Network:")
		fmt.Println("  - Ethereum Mainnet: 60s (5 blocks)")
		fmt.Println("  - Polygon: 20s (10 blocks)")
		fmt.Println("  - Optimism/Arbitrum: 10s (5-40 blocks)")
		fmt.Println("  - Custom: 5 × average block time")
		fmt.Println()

		fmt.Println("Common Causes of Node Staleness:")
		fmt.Println("  1. Initial sync still in progress (check sync status)")
		fmt.Println("  2. Network connectivity issues (check peer count)")
		fmt.Println("  3. Disk I/O bottleneck (check system resources)")
		fmt.Println("  4. CPU/RAM exhaustion (check system resources)")
		fmt.Println("  5. RPC endpoint load balancer routing to stale backend")
		fmt.Println("  6. Database corruption or issues")
		fmt.Println()

		fmt.Println("Production Health Check Architecture:")
		fmt.Println("  1. Primary Check: Block freshness (this module)")
		fmt.Println("  2. Secondary Checks:")
		fmt.Println("     - Sync status (module 21)")
		fmt.Println("     - Peer count (module 22)")
		fmt.Println("     - Mempool size (module 23)")
		fmt.Println("  3. Automated Actions:")
		fmt.Println("     - Failover to backup RPC endpoints")
		fmt.Println("     - Alert ops team via PagerDuty/Slack")
		fmt.Println("     - Log incident for analysis")
		fmt.Println("  4. Recovery Verification:")
		fmt.Println("     - Wait for 3+ consecutive OK checks")
		fmt.Println("     - Gradually restore traffic")
		fmt.Println()

		fmt.Println("Advanced Health Check Patterns:")
		fmt.Println("  - Multi-endpoint comparison: Query multiple RPCs, compare results")
		fmt.Println("  - Canary queries: Test with known block numbers")
		fmt.Println("  - SLA tracking: Measure uptime over 7/30 days")
		fmt.Println("  - Latency monitoring: Track RPC response times")
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
	// 4. Use F11 (Step Into) to enter monitor package functions
	// 5. Watch Variables panel to see health metrics
	// 6. Inspect result.Status, LagSeconds, and BlockTimestamp
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Block Timestamps: How are they set and validated?
	// - Staleness Detection: What makes a node "stale"?
	// - Time-based Monitoring: Why use timestamps for health checks?
	// - Threshold-based Alerting: How to classify OK vs STALE?
	// - Production Patterns: How to build reliable monitoring?
	//
	// STEPPING INTO MONITOR PACKAGE:
	// - Set breakpoint at monitor.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see HeaderByNumber() RPC call
	// - Watch how lag is calculated from timestamps
	// - See how Status is determined from lag threshold
	//
	// COMMON HEALTH MONITORING ISSUES TO DEBUG:
	// - Node always STALE: Check if syncing or check peer count
	// - Negative lag: Clock skew, usually harmless
	// - Intermittent staleness: Network issues or RPC load balancer
	// - High lag on public RPC: Consider running your own node
}
