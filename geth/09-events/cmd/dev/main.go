package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/09-events/internal/events"
)

func main() {
	// BREAKPOINT: deterministic inputs
	rpcURL := "https://eth.llamarpc.com"
	token := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	fromBlock := big.NewInt(0) // set me if you want

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		panic(err)
	}
	toBlock := header.Number
	if fromBlock.Sign() == 0 {
		fromBlock = new(big.Int).Sub(toBlock, big.NewInt(500))
	}

	out, err := events.Run(ctx, client, events.Config{Token: token, FromBlock: fromBlock, ToBlock: toBlock})
	if err != nil {
		panic(err)
	}

	fmt.Println("Events:", len(out.Events))
}
