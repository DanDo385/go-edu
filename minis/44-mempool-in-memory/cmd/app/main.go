package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/44-mempool-in-memory/internal/mempoolinmemory"
)

func main() {
	fmt.Println("=== In-Memory Mempool Demo ===")
	fmt.Println()
	
	// Create a new mempool
	mp := mempoolinmemory.NewMempool(1000)
	
	fmt.Printf("Created mempool with capacity: %d\n", 1000)
	fmt.Printf("Current size: %d\n", mp.Size())
	fmt.Println()
	fmt.Println("See internal/mempoolinmemory/exercise.go for implementation details.")
}
