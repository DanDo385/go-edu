package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"geth/07-eth-call/internal/ethcall"
)

/*
Debug harness for geth/07-eth-call.

Uses USDC contract on mainnet as a test case.
*/

func main() {
	fmt.Println("geth/07-eth-call: Debug Harness")
	fmt.Println("================================")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	// USDC contract address on mainnet
	contractAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

	fmt.Printf("RPC URL:  %s\n", rpcURL)
	fmt.Printf("Contract: %s (USDC)\n", contractAddress.Hex())
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
		BlockNumber: nil, // latest
	}

	result, err := ethcall.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Result ===")
	fmt.Printf("Name:         %s\n", result.Name)
	fmt.Printf("Symbol:      %s\n", result.Symbol)
	fmt.Printf("Decimals:    %d\n", result.Decimals)
	fmt.Printf("Total Supply: %s\n", result.TotalSupply.String())
}
