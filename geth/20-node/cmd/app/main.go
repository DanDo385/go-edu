package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
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
	// geth/20-node
	//
	// This demo calls `web3_clientVersion` via raw JSON-RPC.
	//
	// Usage:
	//   go run ./geth/20-node/cmd/app <RPC_URL>
	//
	// BREAKPOINT
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer c.Close()

	var v string
	if err := c.CallContext(ctx, &v, "web3_clientVersion"); err != nil {
		fmt.Fprintln(os.Stderr, "web3_clientVersion:", err)
		os.Exit(1)
	}
	fmt.Println("clientVersion:", v)
}
