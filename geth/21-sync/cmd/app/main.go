package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/21-sync/internal/sync"
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
	// geth/21-sync
	//
	// Usage:
	//   go run ./geth/21-sync/cmd/app <RPC_URL>
	//
	// BREAKPOINT
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	out, err := sync.Run(ctx, client, sync.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("IsSyncing:", out.IsSyncing)
	if out.Progress != nil {
		fmt.Printf("Progress: %+v\n", *out.Progress)
	}
}
