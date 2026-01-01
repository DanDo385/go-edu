package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/22-peers/internal/peers"
)

func main() {
	// BREAKPOINT
	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := peers.Run(ctx, client, peers.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.PeerCount)
}
