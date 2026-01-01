package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/08-abigen/internal/abigen"
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
	// geth/08-abigen
	//
	// Prerequisites:
	// - geth/06-smart-contracts (console calls/txs)
	// - geth/07-eth-call (manual ABI intuition)
	//
	// Usage:
	//   go run ./geth/08-abigen/cmd/app <RPC_URL> [contract_address] [holder_address]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	contract := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	if len(os.Args) >= 3 {
		contract = common.HexToAddress(os.Args[2])
	}

	var holder *common.Address
	if len(os.Args) >= 4 {
		h := common.HexToAddress(os.Args[3])
		holder = &h
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// BREAKPOINT: run
	out, err := abigen.Run(ctx, client, abigen.Config{Contract: contract, Holder: holder})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Name:", out.Name)
	fmt.Println("Symbol:", out.Symbol)
	fmt.Println("Decimals:", out.Decimals)
	fmt.Println("TotalSupply:", out.TotalSupply)
	if out.Balance != nil {
		fmt.Println("Balance:", out.Balance)
	}
}
