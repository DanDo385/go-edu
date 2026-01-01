package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/15-receipts/internal/receipts"
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
	// geth/15-receipts
	//
	// Usage:
	//   go run ./geth/15-receipts/cmd/app <RPC_URL> [tx_hash]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	var txHash common.Hash
	if len(os.Args) >= 3 {
		txHash = common.HexToHash(os.Args[2])
	} else {
		b, err := client.BlockByNumber(ctx, nil)
		if err != nil || b == nil || len(b.Transactions()) == 0 {
			fmt.Fprintln(os.Stderr, "provide tx_hash: could not auto-select from latest block")
			os.Exit(2)
		}
		txHash = b.Transactions()[0].Hash()
	}

	// BREAKPOINT: run
	out, err := receipts.Run(ctx, client, receipts.Config{TxHash: txHash})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("TxHash:", out.TxHash.Hex())
	fmt.Println("StatusOK:", out.StatusOK)
	fmt.Println("GasUsed:", out.GasUsed)
	fmt.Println("Logs:", len(out.Logs))
}
