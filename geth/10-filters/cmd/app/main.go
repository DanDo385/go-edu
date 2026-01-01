package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/10-filters/internal/filters"
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
	// geth/10-filters
	//
	// This module can use either:
	// - WebSocket subscriptions (real-time)
	// - HTTP polling (works everywhere; slower)
	//
	// Usage:
	//   go run ./geth/10-filters/cmd/app <RPC_URL>
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// Default to polling so this works with HTTP endpoints.
	out, err := filters.Run(ctx, client, filters.Config{MaxHeads: 5, PollInterval: 2 * time.Second, PollMode: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Mode:", out.Mode)
	for _, h := range out.Heads {
		fmt.Println(h.Number, h.Hash.Hex(), "reorg=", h.Reorg)
	}
}
