package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
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
	// geth/19-devnets
	//
	// This module is typically used with a local devnet (e.g. `geth --dev`).
	// This cmd/app confirms connectivity + prints chain metadata.
	//
	// Usage:
	//   go run ./geth/19-devnets/cmd/app <RPC_URL>
	//
	// Example (local dev chain):
	//   geth --dev --http --http.api eth,net,web3,personal
	//   go run ./geth/19-devnets/cmd/app http://127.0.0.1:8545
	//
	// BREAKPOINT
	rpcURL := rpcURLFromArgs("http://127.0.0.1:8545")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer c.Close()

	chainID, _ := c.ChainID(ctx)
	netID, _ := c.NetworkID(ctx)
	head, _ := c.HeaderByNumber(ctx, nil)

	fmt.Println("ChainID:", chainID)
	fmt.Println("NetworkID:", netID)
	if head != nil {
		fmt.Println("Head:", head.Number.Uint64())
	}
}
