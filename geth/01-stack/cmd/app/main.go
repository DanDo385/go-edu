package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/01-stack/internal/stack"
)

/*
===================================================================
Ethereum RPC Stack Verification
===================================================================

This program demonstrates how to connect to an Ethereum RPC endpoint
and verify the execution stack by retrieving fundamental network
identifiers and block headers.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <BLOCK_NUMBER>
    go run ./cmd/app/main.go <RPC_URL> latest
    go run ./cmd/app/main.go <RPC_URL> specific

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, latest, specific, network (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com latest
    go run ./cmd/app/main.go https://eth.llamarpc.com specific
    go run ./cmd/app/main.go https://eth.llamarpc.com network

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
- Step Into (F11) to enter stack package functions
- Watch Variables panel to see RPC responses
- Inspect header, chainID, and networkID structures

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the stack package from internal/stack
- exercise.go, solution.reference.go are in package stack
- The stack.Run() function encapsulates all RPC logic

KEY CONCEPTS:
- Chain ID: Replay protection mechanism (EIP-155)
- Network ID: Legacy P2P network identifier
- Block Headers: Lightweight block metadata with cryptographic commitments
- RPC Client: Interface for communicating with Ethereum nodes
- Context: Go's standard way to handle timeouts and cancellation
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
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com specific\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, latest, specific, network\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum RPC Stack Verification ===\n")
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
	// STEP 3: Example 1 - Retrieve Latest Block Information
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates the most common use case—getting latest block info
	// DEBUG: Watch the 'result' variable to see Chain ID, Network ID, and Header

	if demo == "all" || demo == "latest" {
		fmt.Println("=== Example 1: Latest Block Information ===")
		fmt.Println("Retrieving latest block header and network identifiers...")
		fmt.Println()

		// Create configuration for latest block (BlockNumber: nil means "latest")
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg.BlockNumber is nil—this tells the RPC to fetch the latest block
		cfg := stack.Config{
			BlockNumber: nil, // nil = latest block
		}

		// BREAKPOINT: Set breakpoint here before calling stack.Run
		// DEBUG: Step Into (F11) to see the implementation in exercise.go or solution.reference.go
		// DEBUG: You'll see how the function validates inputs, makes RPC calls, and handles errors
		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := stack.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve stack information: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful retrieval
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at ChainID (1 = Ethereum Mainnet, 11155111 = Sepolia)
		// DEBUG: Examine Header fields: Number, Hash, ParentHash, StateRoot, etc.

		fmt.Println("Network Identifiers:")
		fmt.Printf("  Chain ID:   %s\n", result.ChainID.String())
		fmt.Printf("  Network ID: %s\n\n", result.NetworkID.String())

		fmt.Println("Latest Block Header:")
		fmt.Printf("  Block Number: %s\n", result.Header.Number.String())
		fmt.Printf("  Block Hash:   %s\n", result.Header.Hash().Hex())
		fmt.Printf("  Parent Hash:  %s\n", result.Header.ParentHash.Hex())
		fmt.Printf("  Timestamp:    %d (%s)\n", result.Header.Time, time.Unix(int64(result.Header.Time), 0))
		fmt.Printf("  Gas Limit:    %d\n", result.Header.GasLimit)
		fmt.Printf("  Gas Used:     %d (%.2f%% full)\n",
			result.Header.GasUsed,
			float64(result.Header.GasUsed)/float64(result.Header.GasLimit)*100)
		fmt.Printf("  Base Fee:     %s wei\n", result.Header.BaseFee.String())
		fmt.Println()

		// Display Merkle roots (cryptographic commitments)
		fmt.Println("Cryptographic Commitments (Merkle Roots):")
		fmt.Printf("  State Root:        %s\n", result.Header.Root.Hex())
		fmt.Printf("  Transactions Root: %s\n", result.Header.TxHash.Hex())
		fmt.Printf("  Receipts Root:     %s\n", result.Header.ReceiptHash.Hex())
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Now you can see all the header data displayed
		// DEBUG: Notice how Merkle roots commit to full blockchain state
	}

	// ============================================================================
	// STEP 4: Example 2 - Retrieve Specific Block by Number
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through specific block example
	// DEBUG: This shows how to fetch a historical block by number
	// DEBUG: Useful for deterministic testing and historical analysis

	if demo == "all" || demo == "specific" {
		fmt.Println("=== Example 2: Specific Block by Number ===")
		fmt.Println("Retrieving historical block (block 1,000,000)...")
		fmt.Println()

		// Create configuration for specific block number
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg.BlockNumber is set to 1000000 (a specific historical block)
		// DEBUG: This is useful for deterministic tests and historical analysis
		cfg := stack.Config{
			BlockNumber: big.NewInt(1_000_000), // Specific block number
		}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling stack.Run
		// DEBUG: Step Into (F11) to see how the same function handles different configs
		result, err := stack.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve block 1,000,000: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Compare this result with the latest block result
		// DEBUG: Notice the timestamp is from 2016 (Ethereum's early days)

		fmt.Println("Block #1,000,000 Information:")
		fmt.Printf("  Block Hash:   %s\n", result.Header.Hash().Hex())
		fmt.Printf("  Parent Hash:  %s\n", result.Header.ParentHash.Hex())
		fmt.Printf("  Timestamp:    %d (%s)\n", result.Header.Time, time.Unix(int64(result.Header.Time), 0))
		fmt.Printf("  Gas Used:     %d\n", result.Header.GasUsed)
		fmt.Printf("  Miner:        %s\n", result.Header.Coinbase.Hex())
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 3 - Network Identification and Chain Detection
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through network detection
	// DEBUG: This demonstrates how to identify which Ethereum network you're connected to

	if demo == "all" || demo == "network" {
		fmt.Println("=== Example 3: Network Identification ===")
		fmt.Println("Detecting which Ethereum network we're connected to...")
		fmt.Println()

		cfg := stack.Config{BlockNumber: nil}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := stack.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve network information: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect result.ChainID to see the network identifier
		// DEBUG: Different networks have different Chain IDs

		networkName := "Unknown Network"
		switch result.ChainID.Int64() {
		case 1:
			networkName = "Ethereum Mainnet"
		case 11155111:
			networkName = "Sepolia Testnet"
		case 5:
			networkName = "Goerli Testnet (deprecated)"
		case 137:
			networkName = "Polygon Mainnet"
		case 10:
			networkName = "Optimism Mainnet"
		case 42161:
			networkName = "Arbitrum One"
		case 8453:
			networkName = "Base Mainnet"
		}

		fmt.Printf("Connected Network: %s\n", networkName)
		fmt.Printf("  Chain ID:   %s\n", result.ChainID.String())
		fmt.Printf("  Network ID: %s\n", result.NetworkID.String())
		fmt.Println()

		// Display additional context
		fmt.Println("Chain ID Significance:")
		fmt.Println("  - Introduced by EIP-155 for replay attack prevention")
		fmt.Println("  - Included in transaction signatures")
		fmt.Println("  - Prevents cross-chain transaction replay")
		fmt.Println("  - Each Ethereum-compatible chain has unique Chain ID")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Understanding Defensive Copying
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand defensive copying pattern
	// DEBUG: This demonstrates why the stack package returns copies, not original data

	if demo == "all" {
		fmt.Println("=== Understanding Defensive Copying ===")
		fmt.Println()

		cfg := stack.Config{BlockNumber: nil}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := stack.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve data: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice that result contains copies of the data
		// DEBUG: Modifying result won't affect the RPC client's internal state

		fmt.Println("Why Defensive Copying Matters:")
		fmt.Println("  1. Prevents data races in concurrent code")
		fmt.Println("  2. Ensures callers can't mutate shared state")
		fmt.Println("  3. Makes code safer and easier to reason about")
		fmt.Println("  4. Follows Go's best practices for API design")
		fmt.Println()

		fmt.Println("Example: Modifying returned data doesn't affect original")
		originalChainID := new(big.Int).Set(result.ChainID)
		result.ChainID.Add(result.ChainID, big.NewInt(1)) // Modify the copy

		fmt.Printf("  Original Chain ID: %s\n", originalChainID.String())
		fmt.Printf("  Modified Chain ID: %s\n", result.ChainID.String())
		fmt.Println("  ✓ Modification only affects our copy, not RPC client's data")
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
	// 4. Use F11 (Step Into) to enter stack package functions
	// 5. Watch Variables panel to see RPC responses
	// 6. Inspect result.Header to see all block header fields
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Chain ID vs Network ID: What's the difference?
	// - Block Headers vs Full Blocks: Why use headers?
	// - Merkle Roots: How do they provide cryptographic commitments?
	// - Context: How does it enable cancellation and timeouts?
	// - Defensive Copying: Why is immutability important?
	//
	// STEPPING INTO STACK PACKAGE:
	// - Set breakpoint at stack.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see RPC calls: ChainID(), NetworkID(), HeaderByNumber()
	// - Watch how errors are handled and data is validated
	// - See defensive copying in action with big.Int and Header copies
	//
	// COMMON RPC ERRORS TO DEBUG:
	// - Connection timeout: Check network connectivity and RPC URL
	// - Invalid RPC endpoint: Verify the URL is correct
	// - Rate limiting: Some public endpoints have rate limits
	// - Block not found: Historical blocks might not be available on all nodes
}
