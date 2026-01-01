package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"geth/07-eth-call/internal/ethcall"
)

/*
geth/07-eth-call: Query ERC20 token metadata using manual ABI encoding/decoding.

This module teaches:
- Manual ABI encoding/decoding
- Function selectors
- eth_call RPC method
- Reading contract state

Prerequisites: Complete geth/06-smart-contracts to understand contract interaction concepts.

Usage:
  go run ./cmd/app/main.go <RPC_URL> <contract_address> [block_number]

Examples:
  # Query USDC on mainnet
  go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

  # Query at specific block
  go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 12345
*/

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> <contract_address> [block_number]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 12345\n", os.Args[0])
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	contractAddress := common.HexToAddress(os.Args[2])
	var blockNumber *big.Int

	if len(os.Args) >= 4 {
		blockNum, ok := new(big.Int).SetString(os.Args[3], 10)
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid block number: %s\n", os.Args[3])
			os.Exit(1)
		}
		blockNumber = blockNum
	}

	fmt.Printf("RPC URL:         %s\n", rpcURL)
	fmt.Printf("Contract:        %s\n", contractAddress.Hex())
	if blockNumber != nil {
		fmt.Printf("Block Number:    %s\n", blockNumber.String())
	} else {
		fmt.Printf("Block Number:    latest\n")
	}
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)

	cfg := ethcall.Config{
		Contract:    contractAddress,
		BlockNumber: blockNumber,
	}

	result, err := ethcall.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== ERC20 Token Metadata ===")
	fmt.Printf("Name:         %s\n", result.Name)
	fmt.Printf("Symbol:      %s\n", result.Symbol)
	fmt.Printf("Decimals:    %d\n", result.Decimals)
	fmt.Printf("Total Supply: %s\n", result.TotalSupply.String())
}
