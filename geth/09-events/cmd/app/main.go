package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/09-events/internal/events"
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
	// geth/09-events
	//
	// Prerequisites:
	// - geth/06-smart-contracts (receipts/logs in console)
	// - geth/07-eth-call
	// - geth/08-abigen
	//
	// Usage:
	//   go run ./geth/09-events/cmd/app <RPC_URL> [token_address] [from_block] [to_block]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	token := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	if len(os.Args) >= 3 {
		token = common.HexToAddress(os.Args[2])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// Pick a recent window if not provided.
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "header:", err)
		os.Exit(1)
	}
	toBlock := new(big.Int).Set(header.Number)
	fromBlock := new(big.Int).Sub(toBlock, big.NewInt(500))

	if len(os.Args) >= 5 {
		if v, ok := new(big.Int).SetString(os.Args[3], 10); ok {
			fromBlock = v
		}
		if v, ok := new(big.Int).SetString(os.Args[4], 10); ok {
			toBlock = v
		}
	}

	// BREAKPOINT: run
	out, err := events.Run(ctx, client, events.Config{Token: token, FromBlock: fromBlock, ToBlock: toBlock})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Events:", len(out.Events))
	for i, ev := range out.Events {
		if i >= 5 {
			break
		}
		fmt.Println(ev.BlockNumber, ev.TxHash.Hex(), ev.From.Hex(), "->", ev.To.Hex(), ev.Value)
	}
}
