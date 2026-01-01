package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"geth/24-monitor/internal/monitor"
)

/*
Debug Harness for Monitor Module
*/
func main() {
	fmt.Println("=== Monitor Debug Harness ===")
	fmt.Println()

	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer client.Close()

	// BREAKPOINT: Step into Run()
	result, err := monitor.Run(ctx, client, monitor.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("Next: Proceed to geth/25-toolbox")
}
