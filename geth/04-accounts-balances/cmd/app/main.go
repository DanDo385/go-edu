package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/04-accounts-balances/internal/accountsbalances"
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
	// geth/04-accounts-balances
	//
	// Usage:
	//   go run ./geth/04-accounts-balances/cmd/app <RPC_URL> [addr1] [addr2] ...
	//
	// If no addresses are provided, we query a known EOA + a known contract.
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	addrs := []common.Address{common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")}
	if len(os.Args) >= 3 {
		addrs = nil
		for _, s := range os.Args[2:] {
			addrs = append(addrs, common.HexToAddress(s))
		}
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
	out, err := accountsbalances.Run(ctx, client, accountsbalances.Config{Addresses: addrs})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	for _, a := range out.Accounts {
		fmt.Println("Address:", a.Address.Hex())
		fmt.Println("Type:", a.Type)
		fmt.Println("BalanceWei:", a.Balance)
		fmt.Println("CodeBytes:", len(a.Code))
		fmt.Println()
	}
}
