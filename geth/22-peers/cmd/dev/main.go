package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("Debug Harness: 22-peers")
	
	// Hardcoded for debugging convenience
	rpcURL := "https://eth.llamarpc.com"
	
	fmt.Printf("Connecting to %s...\n", rpcURL)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer client.Close()
	
	id, err := client.ChainID(ctx)
	if err != nil {
		log.Printf("Connected, but failed to get ChainID: %v", err)
	} else {
		fmt.Printf("Successfully connected! ChainID: %s\n", id.String())
	}
	
	// BREAKPOINT: Set breakpoint here to inspect client state
	fmt.Println("Ready to debug internal logic...")
}
