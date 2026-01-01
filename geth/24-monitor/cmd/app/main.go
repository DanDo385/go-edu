package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/24-monitor/internal/monitor"
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
	// geth/24-monitor
	//
	// Usage:
	//   go run ./geth/24-monitor/cmd/app <RPC_URL>
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

	out, err := monitor.Run(ctx, client, monitor.Config{MaxLagSeconds: 120})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Status:", out.Status)
	fmt.Println("BlockNumber:", out.BlockNumber)
	fmt.Println("LagSeconds:", out.LagSeconds)
}
