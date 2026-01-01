package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"geth/04-accounts-balances/internal/accountsbalances"
)

/*
Debug Harness for Accounts and Balances Module
*/
func main() {
	fmt.Println("=== Accounts and Balances Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	// Vitalik's address - a well-known, funded address
	testAddress := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Address: %s (Vitalik)\n", testAddress.Hex())
	fmt.Println()

	// BREAKPOINT: Connect to Ethereum
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
	fmt.Println("Step 2: Querying account state...")

	result, err := accountsbalances.Run(ctx, client, accountsbalances.Config{
		Address: testAddress,
	})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("Step 3: Results")
	fmt.Println()
	fmt.Println("=== Account State ===")
	fmt.Printf("Balance: %s wei\n", result.Balance.String())
	fmt.Printf("Nonce:   %d (number of transactions sent)\n", result.Nonce)
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Balances are stored in wei (1 ETH = 10^18 wei)")
	fmt.Println("2. Nonce tracks number of transactions from an account")
	fmt.Println("3. eth_getBalance and eth_getTransactionCount are the RPC methods used")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/05-tx-nonces")
}
