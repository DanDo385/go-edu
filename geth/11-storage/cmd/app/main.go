package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/11-storage/internal/storage"
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
	// geth/11-storage
	//
	// Usage:
	//   go run ./geth/11-storage/cmd/app <RPC_URL> [contract_address] [slot]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	contract := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	if len(os.Args) >= 3 {
		contract = common.HexToAddress(os.Args[2])
	}

	slot := big.NewInt(0)
	if len(os.Args) >= 4 {
		if v, ok := new(big.Int).SetString(os.Args[3], 10); ok {
			slot = v
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

	// BREAKPOINT: run
	out, err := storage.Run(ctx, client, storage.Config{Contract: contract, Slot: slot})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("ResolvedSlot:", out.ResolvedSlot.Hex())
	fmt.Printf("Value (hex): 0x%x\n", out.Value)
}
