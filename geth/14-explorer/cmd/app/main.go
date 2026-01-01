package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/14-explorer/internal/explorer"
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
	// geth/14-explorer
	//
	// Usage:
	//   go run ./geth/14-explorer/cmd/app <RPC_URL> [block_number] [--txs]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")
	includeTxs := false
	var block *big.Int
	for _, a := range os.Args[2:] {
		if a == "--txs" {
			includeTxs = true
			continue
		}
		if v, ok := new(big.Int).SetString(a, 10); ok {
			block = v
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	out, err := explorer.Run(ctx, client, explorer.Config{Number: block, IncludeTxs: includeTxs})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Block:", out.Number, out.Hash.Hex())
	fmt.Println("TxCount:", out.TxCount)
	if includeTxs {
		for i, tx := range out.Txs {
			if i >= 5 {
				break
			}
			fmt.Println(tx.Hash.Hex(), tx.Gas)
		}
	}
}
