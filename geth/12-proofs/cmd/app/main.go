package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/example/go-10x-minis/geth/12-proofs/internal/proofs"
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
	// geth/12-proofs
	//
	// Usage:
	//   go run ./geth/12-proofs/cmd/app <RPC_URL> [account_address]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	account := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	if len(os.Args) >= 3 {
		account = common.HexToAddress(os.Args[2])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// BREAKPOINT: dial raw RPC (needed for gethclient)
	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial rpc:", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	gc := gethclient.New(rpcClient)

	// BREAKPOINT: run
	out, err := proofs.Run(ctx, gc, proofs.Config{Account: account})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Balance:", out.Account.Balance)
	fmt.Println("Nonce:", out.Account.Nonce)
	fmt.Println("CodeHash:", out.Account.CodeHash.Hex())
	fmt.Println("StorageHash:", out.Account.StorageHash.Hex())
}
