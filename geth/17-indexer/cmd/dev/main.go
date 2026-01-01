package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// BREAKPOINT
	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	h, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("latest:", h.Number.Uint64())
}
