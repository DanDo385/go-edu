package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/10-filters/internal/filters"
)

func main() {
	// BREAKPOINT: deterministic inputs
	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := filters.Run(ctx, client, filters.Config{MaxHeads: 5, PollInterval: 2 * time.Second, PollMode: true})
	if err != nil {
		panic(err)
	}

	fmt.Println("Heads:", len(out.Heads), "mode:", out.Mode)
}
