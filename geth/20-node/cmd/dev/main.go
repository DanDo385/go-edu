package main

import (
	"context"
	"fmt"
	"time"

	"geth/20-node/internal/node"
)

/*
Debug Harness for Node Module
*/
func main() {
	fmt.Println("=== Node Debug Harness ===")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// BREAKPOINT: Step into Run()
	result, err := node.Run(ctx, node.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
	fmt.Println("Next: Proceed to geth/21-sync")
}
