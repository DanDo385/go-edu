package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/04-accounts-balances/internal/accountsbalances"
)

func main() {
	// BREAKPOINT: deterministic inputs
	rpcURL := "https://eth.llamarpc.com"
	addrs := []common.Address{
		common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
		common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := accountsbalances.Run(ctx, client, accountsbalances.Config{Addresses: addrs})
	if err != nil {
		panic(err)
	}

	fmt.Println("Accounts:", len(out.Accounts))
	for _, a := range out.Accounts {
		fmt.Println(a.Address.Hex(), a.Type, a.Balance)
	}
}
