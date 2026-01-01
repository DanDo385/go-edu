package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/14-explorer/internal/explorer"
)

func main() {
	// BREAKPOINT
	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := explorer.Run(ctx, client, explorer.Config{IncludeTxs: true})
	if err != nil {
		panic(err)
	}

	fmt.Println("Block:", out.Number, "txs:", out.TxCount)
}
