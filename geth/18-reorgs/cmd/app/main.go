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
	// geth/18-reorgs (continuity check demo)
	//
	// This demo walks back N headers and verifies parent-hash links.
	// Real reorg handling requires persistence + backfill.
	//
	// Usage:
	//   go run ./geth/18-reorgs/cmd/app <RPC_URL> [n]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")
	n := 20
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

	prev := head
	ok := true
	for i := 0; i < n; i++ {
		num := new(big.Int).Sub(prev.Number, big.NewInt(1))
		h, err := c.HeaderByNumber(ctx, num)
		if err != nil {
			fmt.Fprintln(os.Stderr, "header:", num, err)
			break
		}
		if prev.ParentHash != h.Hash() {
			ok = false
			fmt.Println("MISMATCH at", prev.Number.Uint64(), "parent", prev.ParentHash.Hex(), "expected", h.Hash().Hex())
			break
		}
		prev = h
	}

	fmt.Println("continuity_ok:", ok)
}
