package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/22-peers/internal/peers"
)

/*
===================================================================
Ethereum Node Peer Connectivity Monitoring
===================================================================

This program demonstrates how to check the number of peers connected
to an Ethereum node. Peer count is a critical health indicator for
P2P network connectivity and data propagation.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <DEMO>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, basic, analysis, monitoring (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com basic
    go run ./cmd/app/main.go https://eth.llamarpc.com analysis
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
- Step Into (F11) to enter peers package functions
- Watch Variables panel to see peer count values
- Inspect PeerCount and health status classifications

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the peers package from internal/peers
- exercise.go, solution.reference.go are in package peers
- The peers.Run() function queries peer connectivity

KEY CONCEPTS:
- P2P Networks: Decentralized mesh of interconnected nodes
- Gossip Protocols: How transactions and blocks propagate
- Peer Discovery: DHT (Kademlia) and DNS-based discovery
- Network Health: Peer count indicates connectivity quality
- Connection Types: Inbound vs outbound peer connections
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
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com analysis\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, basic, analysis, monitoring\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Node Peer Connectivity Monitoring ===\n")
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
	// the PeerClient interface. It handles JSON-RPC communication with Ethereum nodes.

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
	// STEP 3: Example 1 - Basic Peer Count Query
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through basic peer count
	// DEBUG: This demonstrates the simplest use case—how many peers?
	// DEBUG: Watch the 'result' variable to see PeerCount

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Peer Count Query ===")
		fmt.Println("Retrieving number of connected peers...")
		fmt.Println()

		// Create basic configuration (no special settings needed)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg is empty—peer count requires no configuration
		cfg := peers.Config{}

		// BREAKPOINT: Set breakpoint here before calling peers.Run
		// DEBUG: Step Into (F11) to see the implementation in exercise.go or solution.reference.go
		// DEBUG: You'll see how PeerCount() RPC is called
		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := peers.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query peer count: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful retrieval
		// DEBUG: Expand 'result' in Variables panel to inspect PeerCount
		// DEBUG: The count is a uint64 representing active peer connections

		fmt.Printf("Connected Peers: %d\n", result.PeerCount)
		fmt.Println()

		// Provide basic health assessment
		// BREAKPOINT: Set breakpoint here to see peer count interpretation
		// DEBUG: Watch how different peer counts indicate different health levels
		if result.PeerCount == 0 {
			fmt.Println("⚠️  WARNING: No peers connected!")
			fmt.Println("   Node is completely isolated from the network")
			fmt.Println("   Check firewall settings and network connectivity")
		} else if result.PeerCount < 5 {
			fmt.Println("⚠️  Low peer count (< 5)")
			fmt.Println("   Network connectivity is limited")
			fmt.Println("   May experience slow block/transaction propagation")
		} else if result.PeerCount < 25 {
			fmt.Println("✓ Good peer count (5-25)")
			fmt.Println("  Healthy connectivity for normal operations")
		} else {
			fmt.Println("✓ Excellent peer count (25+)")
			fmt.Println("  Strong network connectivity and redundancy")
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: You've now seen the basic peer count and health assessment
	}

	// ============================================================================
	// STEP 4: Example 2 - Peer Count Analysis and Interpretation
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through peer analysis
	// DEBUG: This shows how to interpret peer count in context
	// DEBUG: Useful for understanding P2P network topology

	if demo == "all" || demo == "analysis" {
		fmt.Println("=== Example 2: Peer Count Analysis ===")
		fmt.Println("Analyzing peer connectivity and network topology...")
		fmt.Println()

		cfg := peers.Config{}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling peers.Run
		// DEBUG: Step Into (F11) to see PeerCount implementation
		result, err := peers.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query peer count: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect result.PeerCount to see current connectivity

		fmt.Printf("Current Peer Count: %d\n\n", result.PeerCount)

		// Detailed peer count interpretation
		fmt.Println("Peer Count Interpretation:")
		fmt.Println()

		fmt.Println("Connection Quality Assessment:")
		if result.PeerCount == 0 {
			fmt.Println("  Status: CRITICAL - No Network Connectivity")
			fmt.Println("  Impact: Cannot send/receive transactions or blocks")
			fmt.Println("  Action: Check firewall, NAT, and discovery bootstrap nodes")
		} else if result.PeerCount < 3 {
			fmt.Println("  Status: POOR - Minimal Connectivity")
			fmt.Println("  Impact: Vulnerable to network partitions")
			fmt.Println("  Action: Wait for more peer connections or check node configuration")
		} else if result.PeerCount < 10 {
			fmt.Println("  Status: FAIR - Limited Connectivity")
			fmt.Println("  Impact: Acceptable for testing, not ideal for production")
			fmt.Println("  Action: Monitor peer count growth")
		} else if result.PeerCount < 50 {
			fmt.Println("  Status: GOOD - Healthy Connectivity")
			fmt.Println("  Impact: Suitable for production use")
			fmt.Println("  Action: No action needed")
		} else {
			fmt.Println("  Status: EXCELLENT - Strong Connectivity")
			fmt.Println("  Impact: Optimal for production and fast propagation")
			fmt.Println("  Action: No action needed")
		}
		fmt.Println()

		// Display P2P network concepts
		// BREAKPOINT: Set breakpoint here to learn P2P concepts
		// DEBUG: This section explains how peers affect node operation
		fmt.Println("Why Peer Count Matters:")
		fmt.Println("  - Data Propagation: More peers = faster block/tx spread")
		fmt.Println("  - Redundancy: Multiple peers prevent single point of failure")
		fmt.Println("  - Partition Resistance: Higher peer count resists network splits")
		fmt.Println("  - Discovery: Peers help discover more peers via DHT")
		fmt.Println()

		fmt.Println("Typical Peer Count Ranges by Node Type:")
		fmt.Println("  - Full node (default): 25-50 peers")
		fmt.Println("  - Archive node: 50-100 peers (higher capacity)")
		fmt.Println("  - Light client: 5-10 peers (lower resource usage)")
		fmt.Println("  - Mining/validator: 50-100 peers (needs fast propagation)")
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: You've now seen comprehensive peer count analysis
	}

	// ============================================================================
	// STEP 5: Example 3 - Continuous Peer Monitoring
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to see monitoring pattern
	// DEBUG: This demonstrates how to monitor peer count over time

	if demo == "all" || demo == "monitoring" {
		fmt.Println("=== Example 3: Continuous Peer Monitoring ===")
		fmt.Println("Monitoring peer count over time to detect connectivity issues...")
		fmt.Println()

		// Simulate monitoring loop (production would run continuously)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This pattern is what you'd use for health monitoring
		maxChecks := 5
		peerHistory := make([]uint64, 0, maxChecks)

		for i := 0; i < maxChecks; i++ {
			cfg := peers.Config{}

			runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
			result, err := peers.Run(runCtx, client, cfg)
			runCancel()

			if err != nil {
				fmt.Printf("Check %d/%d: ERROR - %v\n", i+1, maxChecks, err)
				time.Sleep(2 * time.Second)
				continue
			}

			// BREAKPOINT: Set breakpoint here in monitoring loop
			// DEBUG: Watch how peer count changes over time
			peerHistory = append(peerHistory, result.PeerCount)

			fmt.Printf("Check %d/%d: ", i+1, maxChecks)
			fmt.Printf("%d peers", result.PeerCount)

			// Detect changes from previous check
			if i > 0 {
				prevCount := peerHistory[i-1]
				diff := int64(result.PeerCount) - int64(prevCount)
				if diff > 0 {
					fmt.Printf(" (+%d from last check)", diff)
				} else if diff < 0 {
					fmt.Printf(" (%d from last check)", diff)
				} else {
					fmt.Printf(" (no change)")
				}
			}
			fmt.Println()

			time.Sleep(2 * time.Second)
		}

		fmt.Println()

		// Analyze peer count trend
		// BREAKPOINT: Set breakpoint here to see trend analysis
		// DEBUG: Watch how we compute average and detect patterns
		if len(peerHistory) > 0 {
			var sum uint64
			minPeers := peerHistory[0]
			maxPeers := peerHistory[0]

			for _, count := range peerHistory {
				sum += count
				if count < minPeers {
					minPeers = count
				}
				if count > maxPeers {
					maxPeers = count
				}
			}

			avgPeers := float64(sum) / float64(len(peerHistory))

			fmt.Println("Peer Count Statistics:")
			fmt.Printf("  Average: %.1f peers\n", avgPeers)
			fmt.Printf("  Minimum: %d peers\n", minPeers)
			fmt.Printf("  Maximum: %d peers\n", maxPeers)
			fmt.Printf("  Range:   %d peers\n", maxPeers-minPeers)
			fmt.Println()

			// Stability assessment
			variability := maxPeers - minPeers
			if variability <= 2 {
				fmt.Println("  ✓ Peer count is STABLE (low variability)")
			} else if variability <= 10 {
				fmt.Println("  ✓ Peer count is RELATIVELY STABLE (moderate variability)")
			} else {
				fmt.Println("  ⚠️  Peer count is UNSTABLE (high variability)")
				fmt.Println("     May indicate network issues or peer churn")
			}
		}

		fmt.Println()

		// Display production monitoring best practices
		fmt.Println("Production Monitoring Best Practices:")
		fmt.Println("  1. Poll peer count every 60 seconds")
		fmt.Println("  2. Alert if peer count drops to 0 for > 1 minute")
		fmt.Println("  3. Alert if peer count stays < 5 for > 5 minutes")
		fmt.Println("  4. Track peer count trends over time")
		fmt.Println("  5. Correlate with sync status and block lag")
		fmt.Println("  6. Set up automatic reconnection if peers drop")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Understanding P2P Networking in Ethereum
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to learn P2P concepts
	// DEBUG: This section explains Ethereum's P2P architecture

	if demo == "all" {
		fmt.Println("=== Understanding Ethereum P2P Networking ===")
		fmt.Println()

		fmt.Println("Peer Discovery Process:")
		fmt.Println("  1. Bootstrap: Connect to hard-coded bootstrap nodes")
		fmt.Println("  2. DHT Lookup: Use Kademlia DHT to find peers")
		fmt.Println("  3. DNS Discovery: Query DNS records for peer lists")
		fmt.Println("  4. Peer Exchange: Learn about peers from connected peers")
		fmt.Println("  5. Connection: Establish TCP connections with discovered peers")
		fmt.Println()

		fmt.Println("Connection Types:")
		fmt.Println("  Outbound Connections:")
		fmt.Println("    - Node initiates connection to discovered peers")
		fmt.Println("    - Typically 25-50 outbound connections")
		fmt.Println("    - More reliable and controllable")
		fmt.Println()
		fmt.Println("  Inbound Connections:")
		fmt.Println("    - Other nodes connect to this node")
		fmt.Println("    - Requires publicly reachable IP/port")
		fmt.Println("    - May need port forwarding/UPnP")
		fmt.Println("    - Typically 0-50 inbound connections")
		fmt.Println()

		fmt.Println("Protocol Negotiation:")
		fmt.Println("  - eth/67: Main Ethereum protocol (blocks, txs, headers)")
		fmt.Println("  - snap/1: State snapshot protocol (fast sync)")
		fmt.Println("  - Peers negotiate which protocols to use")
		fmt.Println("  - Incompatible versions are disconnected")
		fmt.Println()

		fmt.Println("Gossip Protocol:")
		fmt.Println("  - New transactions: Gossiped to subset of peers")
		fmt.Println("  - New blocks: Gossiped to all peers")
		fmt.Println("  - Deduplication: Peers track seen messages")
		fmt.Println("  - Propagation time: ~1-2 seconds for full network")
		fmt.Println()

		fmt.Println("Peer Limitations:")
		fmt.Println("  - Public RPC endpoints:")
		fmt.Println("    * Often show 0 peers (no direct P2P access)")
		fmt.Println("    * Load balancer fronts multiple backend nodes")
		fmt.Println("    * You're not connecting to a real P2P node")
		fmt.Println()
		fmt.Println("  - Your own node:")
		fmt.Println("    * Shows actual peer count")
		fmt.Println("    * Peer details via admin_peers RPC (if enabled)")
		fmt.Println("    * Full control over peer connections")
		fmt.Println()

		fmt.Println("Common Peer Issues:")
		fmt.Println("  - Zero peers: Firewall, NAT, or discovery issues")
		fmt.Println("  - Low peers: Network issues or restrictive config")
		fmt.Println("  - Peer churn: High connect/disconnect rate")
		fmt.Println("  - Max peers reached: Need to increase --maxpeers")
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
	// 4. Use F11 (Step Into) to enter peers package functions
	// 5. Watch Variables panel to see peer count values
	// 6. Inspect result.PeerCount over time
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - P2P Networks: How does decentralized networking work?
	// - Peer Discovery: How do nodes find each other?
	// - Gossip Protocols: How does information spread?
	// - Network Health: What does peer count tell us?
	// - Connection Types: Inbound vs outbound—what's the difference?
	//
	// STEPPING INTO PEERS PACKAGE:
	// - Set breakpoint at peers.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see PeerCount() RPC call
	// - Watch how uint64 is returned (no copying needed)
	// - See error handling for network issues
	//
	// COMMON PEER CONNECTIVITY ISSUES TO DEBUG:
	// - Zero peers: Check if public RPC (they often show 0)
	// - Low peers: Network connectivity or firewall issues
	// - Peer count not increasing: Discovery problems
	// - Fluctuating peer count: Network instability or peer churn
}
