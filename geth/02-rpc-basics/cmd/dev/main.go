package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"geth/02-rpc-basics/internal/rpcbasics"
)

/*
Debug harness for geth/02-rpc-basics with fixed test inputs.
*/

func main() {
	fmt.Println("geth/02-rpc-basics: Debug Harness")
	fmt.Println("==================================")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	retries := 3

	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Retries: %d\n", retries)
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)

	cfg := rpcbasics.Config{
		Retries: retries,
	}

	result, err := rpcbasics.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Result ===")
	fmt.Printf("Network ID:   %s\n", result.NetworkID.String())
	fmt.Printf("Block Number: %d\n", result.BlockNumber)
}
