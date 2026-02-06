package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/07-generic-lru-cache/internal/genericlrucache"
)

/*
Generic LRU Cache Demo

Usage:
  go run ./cmd/app/main.go

This demo shows the LRU cache in action with different data types.
*/
func main() {
	fmt.Println("=== Generic LRU Cache Demo ===")
	fmt.Println()

	// Create a cache with capacity 3 and 5 second TTL
	cache := genericlrucache.New[string, int](3)

	fmt.Println("--- Adding Items (capacity: 3) ---")
	cache.Set("one", 1)
	fmt.Println("Set 'one' = 1")

	cache.Set("two", 2)
	fmt.Println("Set 'two' = 2")

	cache.Set("three", 3)
	fmt.Println("Set 'three' = 3")

	// Access 'one' to make it recently used
	if val, ok := cache.Get("one"); ok {
		fmt.Printf("Get 'one' = %d (moves to front)\n", val)
	}
	fmt.Println()

	// Add fourth item - should evict 'two' (least recently used)
	fmt.Println("--- Adding 4th Item (triggers eviction) ---")
	cache.Set("four", 4)
	fmt.Println("Set 'four' = 4")

	// Check what's still in cache
	fmt.Println()
	fmt.Println("--- Checking Cache Contents ---")
	for _, key := range []string{"one", "two", "three", "four"} {
		if val, ok := cache.Get(key); ok {
			fmt.Printf("  '%s' = %d (found)\n", key, val)
		} else {
			fmt.Printf("  '%s' = NOT FOUND (evicted)\n", key)
		}
	}

	fmt.Println()
	fmt.Println("=== Key Concepts ===")
	fmt.Println("1. LRU evicts least recently used items")
	fmt.Println("2. O(1) Get and Set operations")
	fmt.Println("3. Generics allow type-safe caching of any types")
	fmt.Println("4. Thread-safe with sync.Mutex")
}
