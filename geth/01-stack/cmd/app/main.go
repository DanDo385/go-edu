package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/01-stack/internal/stack"
)

/*
geth/01-stack: cmd/app

Demonstrates basic Ethereum RPC connectivity and stack information retrieval.

This application proves RPC connectivity by reading:
  - Chain ID (replay protection identifier)
  - Network ID (legacy P2P network identifier)
  - Block header (cryptographic commitments to state)

Usage:
  go run ./cmd/app/main.go <RPC_URL> [block_number]

Examples:
  # Query latest block from public RPC
  go run ./cmd/app/main.go https://eth.llamarpc.com

  # Query specific block number
  go run ./cmd/app/main.go https://eth.llamarpc.com 12345

  # Query Sepolia testnet
  go run ./cmd/app/main.go https://ethereum-sepolia-rpc.publicnode.com

  # Query local Geth node
  go run ./cmd/app/main.go http://localhost:8545

Arguments:
  RPC_URL       Ethereum RPC endpoint URL (required)
  block_number  Specific block number to query (optional, defaults to latest)

BREAKPOINT: Set breakpoints at key locations marked with "// BREAKPOINT:" comments.
*/

func main() {
	// ========================================================================
	// STEP 1: Parse CLI Arguments
	// ========================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [block_number]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 12345\n", os.Args[0])
		os.Exit(1)
	}

	rpcURL := os.Args[1]

	// Parse optional block number argument
	var blockNumber *big.Int
	if len(os.Args) >= 3 {
		n := new(big.Int)
		_, ok := n.SetString(os.Args[2], 10)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: Invalid block number: %s\n", os.Args[2])
			os.Exit(1)
		}
		blockNumber = n
	}

	// ========================================================================
	// STEP 2: Connect to Ethereum Node
	// ========================================================================
	// BREAKPOINT: Set breakpoint here to see RPC connection being established
	fmt.Printf("Connecting to Ethereum node: %s\n", rpcURL)

	// Create context with timeout to prevent hanging forever
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Dial RPC endpoint
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to RPC endpoint: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("✓ Connected successfully")
	fmt.Println()

	// ========================================================================
	// STEP 3: Call Our Exercise Function
	// ========================================================================
	// BREAKPOINT: Set breakpoint here to step into the Run function
	cfg := stack.Config{
		BlockNumber: blockNumber,
	}

	fmt.Println("Retrieving Ethereum stack information...")
	result, err := stack.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving stack info: %v\n", err)
		os.Exit(1)
	}

	// ========================================================================
	// STEP 4: Display Results
	// ========================================================================
	// BREAKPOINT: Set breakpoint here to inspect the result
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║           Ethereum Stack Information                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Chain ID: Used for transaction signing (EIP-155 replay protection)
	fmt.Printf("Chain ID:    %s\n", result.ChainID.String())
	fmt.Printf("             (1=mainnet, 11155111=Sepolia, 1337=dev)\n")
	fmt.Println()

	// Network ID: Legacy identifier for P2P networking
	fmt.Printf("Network ID:  %s\n", result.NetworkID.String())
	fmt.Println()

	// Block Header Information
	fmt.Println("Block Header:")
	fmt.Printf("  Number:      %s\n", result.Header.Number.String())
	fmt.Printf("  Hash:        %s\n", result.Header.Hash().Hex())
	fmt.Printf("  Parent Hash: %s\n", result.Header.ParentHash.Hex())
	fmt.Printf("  State Root:  %s\n", result.Header.Root.Hex())
	fmt.Printf("  Timestamp:   %s (%d)\n", 
		time.Unix(int64(result.Header.Time), 0).Format(time.RFC3339), 
		result.Header.Time)
	fmt.Printf("  Gas Limit:   %d\n", result.Header.GasLimit)
	fmt.Printf("  Gas Used:    %d (%.2f%% full)\n", 
		result.Header.GasUsed,
		float64(result.Header.GasUsed)/float64(result.Header.GasLimit)*100)
	
	// EIP-1559 fields (if present)
	if result.Header.BaseFee != nil {
		fmt.Printf("  Base Fee:    %s wei\n", result.Header.BaseFee.String())
	}
	
	fmt.Println()
	fmt.Println("✓ Successfully retrieved Ethereum stack information")

	// ========================================================================
	// Conceptual Notes
	// ========================================================================
	// This demonstrates:
	// - RPC connectivity (ethclient.DialContext)
	// - Context-based timeouts (prevents hanging on slow/dead endpoints)
	// - Reading chain metadata (ChainID, NetworkID)
	// - Reading block headers (lightweight cryptographic commitments)
	// - Proper resource cleanup (defer client.Close())
	//
	// Next steps:
	// - geth/02-rpc-basics: Understanding different RPC methods
	// - geth/03-keys-addresses: Working with keys and addresses
	// - geth/04-accounts-balances: Reading account states
}
