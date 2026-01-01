package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/02-rpc-basics/internal/rpcbasics"
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
	// geth/02-rpc-basics
	//
	// Usage:
	//   go run ./geth/02-rpc-basics/cmd/app <RPC_URL> [retries]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	retries := 3
	if len(os.Args) >= 3 {
		if v, err := strconv.Atoi(os.Args[2]); err == nil {
			retries = v
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// BREAKPOINT: dial
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// BREAKPOINT: run
	out, err := rpcbasics.Run(ctx, client, rpcbasics.Config{Retries: retries})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("NetworkID:", out.NetworkID)
	fmt.Println("Latest BlockNumber:", out.BlockNumber)
	if out.Block != nil {
		fmt.Println("Block.Hash:", out.Block.Hash())
		fmt.Println("Block.TxCount:", len(out.Block.Transactions()))
		fmt.Println("Block.Number:", out.Block.Number())
	}
}
