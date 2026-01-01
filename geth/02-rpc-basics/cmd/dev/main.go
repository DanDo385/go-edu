package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/02-rpc-basics/internal/rpcbasics"
)

func main() {
	// BREAKPOINT: deterministic inputs
	rpcURL := "https://eth.llamarpc.com"
	retries := 3

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := rpcbasics.Run(ctx, client, rpcbasics.Config{Retries: retries})
	if err != nil {
		panic(err)
	}

	fmt.Println("NetworkID:", out.NetworkID)
	fmt.Println("BlockNumber:", out.BlockNumber)
}
