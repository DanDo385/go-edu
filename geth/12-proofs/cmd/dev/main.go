package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/example/go-10x-minis/geth/12-proofs/internal/proofs"
)

func main() {
	// BREAKPOINT: deterministic inputs
	rpcURL := "https://eth.llamarpc.com"
	account := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer rpcClient.Close()

	gc := gethclient.New(rpcClient)

	out, err := proofs.Run(ctx, gc, proofs.Config{Account: account})
	if err != nil {
		panic(err)
	}

	fmt.Println("Nonce:", out.Account.Nonce)
}
