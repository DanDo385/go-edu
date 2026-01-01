package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/15-receipts/internal/receipts"
)

func main() {
	// BREAKPOINT: deterministic inputs
	rpcURL := "https://eth.llamarpc.com"
	txHash := common.Hash{} // optionally set

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	if txHash == (common.Hash{}) {
		b, err := client.BlockByNumber(ctx, nil)
		if err != nil || b == nil || len(b.Transactions()) == 0 {
			panic("no txs in latest block")
		}
		txHash = b.Transactions()[0].Hash()
	}

	out, err := receipts.Run(ctx, client, receipts.Config{TxHash: txHash})
	if err != nil {
		panic(err)
	}

	fmt.Println("Logs:", len(out.Logs))
}
