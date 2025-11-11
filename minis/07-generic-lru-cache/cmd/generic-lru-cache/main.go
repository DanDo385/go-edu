package main

import (
	"fmt"
	"time"

	"github.com/example/go-10x-minis/minis/07-generic-lru-cache/exercise"
)

func main() {
	// Create a cache with capacity 3 and 2-second TTL
	cache := exercise.New[string, int](3, 2*time.Second)

	fmt.Println("=== LRU Cache Demo ===\n")

	// Add items
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)
	fmt.Println("Added: a=1, b=2, c=3")
	fmt.Printf("Size: %d/3\n\n", cache.Len())

	// Access "a" (makes it most recent)
	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Get 'a': %d (moved to front)\n", val)
	}

	// Add "d" (should evict "b", the least recently used)
	cache.Set("d", 4)
	fmt.Println("Added: d=4")
	fmt.Printf("Size: %d/3\n", cache.Len())

	// Check if "b" was evicted
	if _, ok := cache.Get("b"); !ok {
		fmt.Println("'b' was evicted (LRU)\n")
	}

	// Wait for TTL expiration
	fmt.Println("Waiting 3 seconds for TTL expiration...")
	time.Sleep(3 * time.Second)

	// Try to get expired items
	if _, ok := cache.Get("a"); !ok {
		fmt.Println("'a' expired (TTL)")
	}
	if _, ok := cache.Get("c"); !ok {
		fmt.Println("'c' expired (TTL)")
	}
	if _, ok := cache.Get("d"); !ok {
		fmt.Println("'d' expired (TTL)")
	}

	fmt.Printf("\nFinal size: %d\n", cache.Len())
}
