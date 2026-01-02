package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/07-eth-call/internal/ethcall"
)

/*
Debug Harness for eth_call Module

This is the Go implementation of what you learned in geth/06-smart-contracts.
Step through to see manual ABI encoding/decoding.
*/
func main() {
	fmt.Println("=== eth_call Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	// USDC contract address
	contractAddr := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

	fmt.Printf("RPC URL:  %s\n", rpcURL)
	fmt.Printf("Contract: %s (USDC)\n", contractAddr.Hex())
	fmt.Println()

	// BREAKPOINT: Connect
	fmt.Println("Step 1: Connecting...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("         Connected!")
	fmt.Println()

	// BREAKPOINT: Step into Run() to see:
	// - Function selector computation (keccak256[:4])
	// - CallMsg construction
	// - eth_call execution
	// - ABI decoding (strings, uint8, uint256)
	fmt.Println("Step 2: Calling contract methods...")
	fmt.Println("        - name()")
	fmt.Println("        - symbol()")
	fmt.Println("        - decimals()")
	fmt.Println("        - totalSupply()")
	fmt.Println()

	result, err := ethcall.Run(ctx, client, ethcall.Config{
		Contract: contractAddr,
	})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("Step 3: Results")
	fmt.Println()
	fmt.Println("=== USDC Token Metadata ===")
	fmt.Printf("Name:         %s\n", result.Name)
	fmt.Printf("Symbol:       %s\n", result.Symbol)
	fmt.Printf("Decimals:     %d\n", result.Decimals)
	fmt.Printf("Total Supply: %s\n", result.TotalSupply.String())
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Function selectors: keccak256('name()')[:4] = 0x06fdde03")
	fmt.Println("2. eth_call: Read-only execution, no gas cost")
	fmt.Println("3. ABI decoding: offset+length+data for strings, 32-byte words for ints")
	fmt.Println("4. This is exactly what the console does, but in Go")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/08-abigen for typed bindings")
}
