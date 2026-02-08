package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/06-eip1559/internal/eip1559"
)

/*
Debug Harness for EIP-1559 Module
*/
func main() {
	fmt.Println("=== EIP-1559 Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"
	fmt.Printf("RPC URL: %s\n", rpcURL)
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
	fmt.Println("Step 2: Fetching EIP-1559 fee data...")

	result, err := eip1559.Run(ctx, client, eip1559.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("Step 3: Results")
	fmt.Println()
	fmt.Printf("Base Fee:      %s\n", result.BaseFee.String())
	fmt.Printf("Max Fee:       %s\n", result.MaxFee.String())
	fmt.Printf("Priority Fee:  %s\n", result.PriorityFee.String())
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Base Fee: Set by protocol, burned on inclusion")
	fmt.Println("2. Priority Fee (Tip): Goes to validator, incentivizes inclusion")
	fmt.Println("3. Max Fee: Maximum you're willing to pay per gas")
	fmt.Println("4. Effective gas price = min(base_fee + priority_fee, max_fee)")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/06-smart-contracts for console tutorial")
	fmt.Println("      Then geth/07-eth-call for Go contract calls")
}
