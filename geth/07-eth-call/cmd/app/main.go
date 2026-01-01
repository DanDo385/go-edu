package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/07-eth-call/internal/ethcall"
)

func rpcURLFromArgs(defaultURL string) string {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		return os.Args[1]
	}
	if v := os.Getenv("RPC_URL"); v != "" {
		return v
	}
	return defaultURL
}

func main() {
	// geth/07-eth-call
	//
	// Prerequisite: geth/06-smart-contracts (console intuition for calls vs txs).
	//
	// Usage:
	//   go run ./geth/07-eth-call/cmd/app <RPC_URL> [contract_address]
	//
	// Example:
	//   go run ./geth/07-eth-call/cmd/app https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	contract := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	if len(os.Args) >= 3 {
		contract = common.HexToAddress(os.Args[2])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// BREAKPOINT: run
	out, err := ethcall.Run(ctx, client, ethcall.Config{Contract: contract})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Name:", out.Name)
	fmt.Println("Symbol:", out.Symbol)
	fmt.Println("Decimals:", out.Decimals)
	fmt.Println("TotalSupply:", out.TotalSupply)
}
