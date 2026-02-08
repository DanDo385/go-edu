package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/16-concurrency/internal/concurrency"
)

// ethProber wraps ethclient to satisfy concurrency.Prober.
// Students will implement the real probe logic in the exercise.
type ethProber struct{ c *ethclient.Client }

func (p *ethProber) Probe(ctx context.Context, endpoint string) error {
	return fmt.Errorf("not yet implemented - complete the exercise")
}

/*
Debug Harness for Concurrency Module
*/
func main() {
	fmt.Println("=== Concurrency Debug Harness ===")
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
	result, err := concurrency.Run(ctx, &ethProber{client}, concurrency.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("Next: Proceed to geth/17-indexer")
}
