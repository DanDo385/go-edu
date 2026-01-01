package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
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
	// geth/17-indexer (lightweight demo)
	//
	// This module's internal exercise is a placeholder in this repo snapshot.
	// The cmd/app demonstrates a tiny "index" pass: fetch N recent blocks and count txs.
	//
	// Usage:
	//   go run ./geth/17-indexer/cmd/app <RPC_URL> [n]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")
	n := 5
	if len(os.Args) >= 3 {
		if v, err := strconv.Atoi(os.Args[2]); err == nil {
			n = v
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer c.Close()

	head, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "header:", err)
		os.Exit(1)
	}

	start := new(big.Int).Sub(head.Number, big.NewInt(int64(n-1)))
	var total int
	for i := new(big.Int).Set(start); i.Cmp(head.Number) <= 0; i.Add(i, big.NewInt(1)) {
		b, err := c.BlockByNumber(ctx, i)
		if err != nil {
			fmt.Fprintln(os.Stderr, "block:", i, err)
			continue
		}
		total += len(b.Transactions())
		fmt.Println("block", b.NumberU64(), "txs", len(b.Transactions()))
	}
	fmt.Println("total txs:", total)
}
