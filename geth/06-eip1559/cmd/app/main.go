package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("Project: 06-eip1559")
	
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <rpc_url> [args...]\n", os.Args[0])
		fmt.Println("Example: go run ./cmd/app/main.go https://eth.llamarpc.com")
		// Don't exit, just print info
	} else {
		rpcURL := os.Args[1]
		fmt.Printf("Connecting to %s...\n", rpcURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		client, err := ethclient.DialContext(ctx, rpcURL)
		if err != nil {
			log.Printf("Failed to connect: %v", err)
		} else {
			defer client.Close()
			id, err := client.ChainID(ctx)
			if err != nil {
				log.Printf("Connected, but failed to get ChainID: %v", err)
			} else {
				fmt.Printf("Successfully connected! ChainID: %s\n", id.String())
			}
		}
	}

	fmt.Println("\nSee internal/eip1559/exercise.go for the implementation.")
}
