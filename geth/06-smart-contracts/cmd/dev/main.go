package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/06-smart-contracts/internal/smartcontracts"
)

/*
Debug Harness for Smart Contract Interaction

This file provides a deterministic debug environment with fixed inputs.
Use this for stepping through the code with VS Code debugger.

How to use:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5 in VS Code
  3. Select "Debug: cmd/dev (Debug Harness)"
  4. Step through with F10 (Step Over) and F11 (Step Into)

This harness:
  - Uses a public RPC endpoint (no setup required)
  - Queries USDC on Ethereum mainnet (well-known, stable contract)
  - Has fixed inputs for reproducible debugging sessions
*/
func main() {
	fmt.Println("=== Smart Contracts Debug Harness ===")
	fmt.Println()

	// ============================================================================
	// Fixed Test Configuration
	// ============================================================================
	// These values are fixed for reproducible debugging.
	// USDC is chosen because it's a stable, well-known ERC20 token.
	//
	// BREAKPOINT: Inspect these values to understand the test setup
	rpcURL := "https://eth.llamarpc.com"
	contractAddr := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" // USDC

	fmt.Printf("RPC URL:  %s\n", rpcURL)
	fmt.Printf("Contract: %s (USDC)\n", contractAddr)
	fmt.Println()

	// ============================================================================
	// Connect to Ethereum
	// ============================================================================
	// BREAKPOINT: Step into DialContext to see connection establishment
	fmt.Println("Step 1: Connecting to Ethereum node...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect: %v\n", err)
		fmt.Println()
		fmt.Println("Troubleshooting:")
		fmt.Println("  - Check your internet connection")
		fmt.Println("  - Try a different RPC URL")
		fmt.Println("  - The public endpoint may be rate-limited")
		return
	}
	defer client.Close()

	fmt.Println("         Connected successfully!")
	fmt.Println()

	// ============================================================================
	// Parse Contract Address
	// ============================================================================
	// BREAKPOINT: Inspect how hex string becomes common.Address
	fmt.Println("Step 2: Parsing contract address...")
	contract := common.HexToAddress(contractAddr)
	fmt.Printf("         Parsed: %s\n", contract.Hex())
	fmt.Println()

	// ============================================================================
	// Query Contract
	// ============================================================================
	// BREAKPOINT: This is the main entry point - step into Run()
	fmt.Println("Step 3: Querying contract (eth_call)...")
	fmt.Println("        This is the Go equivalent of:")
	fmt.Println("          myContract.name()")
	fmt.Println("          myContract.symbol()")
	fmt.Println("          myContract.decimals()")
	fmt.Println("          myContract.totalSupply()")
	fmt.Println()

	result, err := smartcontracts.Run(ctx, client, smartcontracts.Config{
		Contract: contract,
	})
	if err != nil {
		fmt.Printf("ERROR: Query failed: %v\n", err)
		fmt.Println()
		fmt.Println("Troubleshooting:")
		fmt.Println("  - Verify the contract address is correct")
		fmt.Println("  - Ensure the contract is an ERC20 token")
		fmt.Println("  - Check if the RPC endpoint supports eth_call")
		return
	}

	// ============================================================================
	// Display Results
	// ============================================================================
	// BREAKPOINT: Inspect the decoded result struct
	fmt.Println("Step 4: Results (decoded from ABI encoding)")
	fmt.Println()
	fmt.Println("=== USDC Token Metadata ===")
	fmt.Printf("Name:         %s\n", result.Name)
	fmt.Printf("Symbol:       %s\n", result.Symbol)
	fmt.Printf("Decimals:     %d\n", result.Decimals)
	fmt.Printf("Total Supply: %s\n", result.TotalSupply.String())
	fmt.Println()

	// ============================================================================
	// Educational Notes
	// ============================================================================
	fmt.Println("=== What You Just Learned ===")
	fmt.Println()
	fmt.Println("1. FUNCTION SELECTORS")
	fmt.Println("   - Each function has a 4-byte selector: keccak256(signature)[:4]")
	fmt.Println("   - name() -> 0x06fdde03")
	fmt.Println("   - symbol() -> 0x95d89b41")
	fmt.Println("   - decimals() -> 0x313ce567")
	fmt.Println("   - totalSupply() -> 0x18160ddd")
	fmt.Println()
	fmt.Println("2. eth_call vs eth_sendTransaction")
	fmt.Println("   - eth_call: Read-only, no gas, no signature needed")
	fmt.Println("   - eth_sendTransaction: State change, costs gas, needs signature")
	fmt.Println("   - All queries above used eth_call")
	fmt.Println()
	fmt.Println("3. ABI DECODING")
	fmt.Println("   - Strings: offset + length + data (dynamic type)")
	fmt.Println("   - uint8/uint256: 32-byte word, value right-aligned (static type)")
	fmt.Println()
	fmt.Println("Next: Complete the Geth console tutorial in README.md")
	fmt.Println("      Then proceed to geth/07-eth-call for deeper Go implementation")
}
