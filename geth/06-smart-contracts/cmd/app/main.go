package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/06-smart-contracts/internal/smartcontracts"
)

/*
Smart Contract Interaction Demo

This application demonstrates the concepts learned in the Geth console tutorial.
It queries ERC20 token metadata using eth_call, showing the Go equivalent of
what you did interactively in the console.

Usage:

	go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>

Examples:

	# Query USDC on Ethereum mainnet
	go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

	# Query DAI on Ethereum mainnet
	go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EesedCDAE5F


	# Query WETH on Ethereum mainnet
	go run ./cmd/app/main.go https://eth.llamarpc.com 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2

Arguments:

	RPC_URL          - Ethereum RPC endpoint (e.g., https://eth.llamarpc.com)
	CONTRACT_ADDRESS - ERC20 token contract address to query

This demonstrates:
  - Connecting to an Ethereum node (same as geth attach)
  - Making eth_call requests (same as myContract.name() in console)
  - Decoding ABI-encoded responses (done automatically in console)
*/
func main() {
	// ============================================================================
	// STEP 1: Parse Command Line Arguments
	// ============================================================================
	// BREAKPOINT: Set a breakpoint here to inspect command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  # Query USDC on Ethereum mainnet")
		fmt.Println("  go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
		fmt.Println()
		fmt.Println("  # Query DAI on Ethereum mainnet")
		fmt.Println("  go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedinCDAE5F")
		fmt.Println()
		fmt.Println("This is the Go equivalent of the Geth console tutorial.")
		fmt.Println("See README.md for the console-based learning experience.")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	contractAddr := os.Args[2]

	// ============================================================================
	// STEP 2: Connect to Ethereum Node
	// ============================================================================
	// This is equivalent to running: geth attach <RPC_URL>
	//
	// BREAKPOINT: Set a breakpoint here to step into the connection
	fmt.Printf("Connecting to %s...\n", rpcURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("Connected!")

	// ============================================================================
	// STEP 3: Parse Contract Address
	// ============================================================================
	// In the console, you typed: var contractAddress = "0x..."
	// In Go, we use common.HexToAddress() to parse the string.
	//
	// BREAKPOINT: Inspect the parsed address
	contract := common.HexToAddress(contractAddr)
	fmt.Printf("Querying contract: %s\n", contract.Hex())
	fmt.Println()

	// ============================================================================
	// STEP 4: Query Contract (This is the Core Concept)
	// ============================================================================
	// In the console, you called:
	//   myContract.name()
	//   myContract.symbol()
	//   myContract.decimals()
	//   myContract.totalSupply()
	//
	// The smartcontracts.Run() function does the same thing, but shows you
	// the underlying mechanics:
	//   1. Compute function selectors from signatures
	//   2. Build CallMsg with contract address and selector
	//   3. Execute eth_call via client.CallContract
	//   4. Decode ABI-encoded responses
	//
	// BREAKPOINT: Step into Run() to see the contract call mechanics
	result, err := smartcontracts.Run(ctx, client, smartcontracts.Config{
		Contract: contract,
	})
	if err != nil {
		fmt.Printf("Failed to query contract: %v\n", err)
		os.Exit(1)
	}

	// ============================================================================
	// STEP 5: Display Results
	// ============================================================================
	// This is what you saw in the console after each call.
	fmt.Println("=== ERC20 Token Metadata ===")
	fmt.Printf("Name:         %s\n", result.Name)
	fmt.Printf("Symbol:       %s\n", result.Symbol)
	fmt.Printf("Decimals:     %d\n", result.Decimals)
	fmt.Printf("Total Supply: %s\n", result.TotalSupply.String())
	fmt.Println()

	// ============================================================================
	// Key Takeaways
	// ============================================================================
	fmt.Println("=== What Happened Under the Hood ===")
	fmt.Println("1. Connected to Ethereum node (like 'geth attach')")
	fmt.Println("2. Computed function selectors: keccak256('name()')[:4], etc.")
	fmt.Println("3. Sent eth_call requests (read-only, no gas consumed)")
	fmt.Println("4. Decoded ABI-encoded responses (strings, uint8, uint256)")
	fmt.Println()
	fmt.Println("Next: See geth/07-eth-call for deeper Go implementation details")
}
