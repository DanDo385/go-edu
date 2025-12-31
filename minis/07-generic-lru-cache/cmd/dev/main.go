// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/genericlrucache.go
// 2. Open this file (cmd/debug/main.go)
// 3. Use "Debug Main Program (Current Package)" configuration
// 4. Press F5 - that's it! The debugger will stop at your breakpoints
//
// Usage:
//   go run ./cmd/debug
//   # Or just press F5 in VS Code

package main

import (
	"fmt"
	"time"

	"github.com/example/go-10x-minis/minis/07-generic-lru-cache/internal/genericlrucache"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	capacity := 3
	ttl := time.Duration(0) // No TTL for this test

	fmt.Println("=== Debugging LRU Cache ===")
	fmt.Printf("Capacity: %d, TTL: %v\n\n", capacity, ttl)

	// Set breakpoint in genericlrucache.go at New function
	cache := genericlrucache.New[string, int](capacity, ttl)

	// Test Set operations
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Test Get operation
	val, ok := cache.Get("a")
	fmt.Printf("Get('a'): val=%d, ok=%v\n", val, ok)

	// Test eviction (add 4th item, should evict 'a')
	cache.Set("d", 4)
	val, ok = cache.Get("a")
	fmt.Printf("Get('a') after eviction: val=%d, ok=%v\n", val, ok)

	fmt.Printf("Cache length: %d\n", cache.Len())
}

