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
	// geth/01-stack
	//
	// Usage:
	//   go run ./geth/01-stack/cmd/app <RPC_URL> [block_number]
	//
	// Examples:
	//   go run ./geth/01-stack/cmd/app https://eth.llamarpc.com
	//   go run ./geth/01-stack/cmd/app https://eth.llamarpc.com 19000000
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	var blockNumber *big.Int
	if len(os.Args) >= 3 {
		// BREAKPOINT: parse optional block number
		var n big.Int
		if _, ok := n.SetString(os.Args[2], 10); ok {
			blockNumber = &n
		} else {
			fmt.Fprintln(os.Stderr, "invalid block_number (expected base-10 uint)")
			os.Exit(2)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// BREAKPOINT: dial RPC
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// BREAKPOINT: call internal module logic
	out, err := stack.Run(ctx, client, stack.Config{BlockNumber: blockNumber})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("ChainID:", out.ChainID)
	fmt.Println("NetworkID:", out.NetworkID)
	fmt.Println("Header.Number:", out.Header.Number)
	fmt.Println("Header.Hash:", out.Header.Hash())
}
