package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/01-stack/internal/stack"
)

func main() {
	// Deterministic debug harness:
	// - fixed RPC URL
	// - fixed optional block number (nil => latest)
	//
	// BREAKPOINT: change inputs here
	rpcURL := "https://eth.llamarpc.com"
	var blockNumber *big.Int = nil

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// BREAKPOINT: dial RPC
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// BREAKPOINT: run
	out, err := stack.Run(ctx, client, stack.Config{BlockNumber: blockNumber})
	if err != nil {
		panic(err)
	}

	fmt.Println("ChainID:", out.ChainID)
	fmt.Println("Header.Number:", out.Header.Number)
}
