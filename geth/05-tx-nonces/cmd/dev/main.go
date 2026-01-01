package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/05-tx-nonces/internal/txnonces"
)

/*
Debug Harness for Transaction Nonces Module
*/
func main() {
	fmt.Println("=== Transaction Nonces Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	testAddress := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Address: %s\n", testAddress.Hex())
	fmt.Println()

	// BREAKPOINT: Connect
	fmt.Println("Step 1: Connecting...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("         Connected!")
	fmt.Println()

	// BREAKPOINT: Step into Run()
	fmt.Println("Step 2: Querying nonces...")

	result, err := txnonces.Run(ctx, client, txnonces.Config{
		Address: testAddress,
	})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("Step 3: Results")
	fmt.Println()
	fmt.Printf("Confirmed Nonce: %d\n", result.Nonce)
	fmt.Printf("Pending Nonce:   %d\n", result.PendingNonce)
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Nonce = number of transactions sent from an account")
	fmt.Println("2. Confirmed nonce = nonce of mined transactions")
	fmt.Println("3. Pending nonce = next nonce to use (includes mempool)")
	fmt.Println("4. Nonces must be sequential - gaps cause transactions to get stuck")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/06-smart-contracts for console tutorial")
}
