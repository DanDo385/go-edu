package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/13-trace/internal/trace"
)

// traceAdapter wraps ethclient to satisfy trace.TraceClient.
// Students will implement the real trace logic in the exercise.
type traceAdapter struct{ c *ethclient.Client }

func (a *traceAdapter) TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error) {
	return nil, fmt.Errorf("not yet implemented - complete the exercise")
}

/*
Debug Harness for Trace Module
*/
func main() {
	fmt.Println("=== Trace Debug Harness ===")
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
	result, err := trace.Run(ctx, &traceAdapter{client}, trace.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("Next: Proceed to geth/14-explorer")
}
