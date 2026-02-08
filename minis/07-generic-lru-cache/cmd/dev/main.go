package main

import (
	"fmt"
	"time"

	"github.com/example/go-10x-minis/minis/07-generic-lru-cache/internal/genericlrucache"
)

/*
Debug Harness for Generic LRU Cache Module

This file demonstrates all cache operations.
Set breakpoints in exercise.go and press F5 in VS Code.
*/
func main() {
	fmt.Println("=== Generic LRU Cache Debug Harness ===")
	fmt.Println()

	// Demo 1: Basic operations
	fmt.Println("--- Demo 1: Basic Operations ---")
	cache := genericlrucache.New[string, string](3)

	cache.Set("user:1", "Alice")
	cache.Set("user:2", "Bob")
	cache.Set("user:3", "Charlie")

	for _, key := range []string{"user:1", "user:2", "user:3"} {
		if val, ok := cache.Get(key); ok {
			fmt.Printf("  %s = %s\n", key, val)
		}
	}
	fmt.Println()

	// Demo 2: LRU eviction
	fmt.Println("--- Demo 2: LRU Eviction ---")
	fmt.Println("Cache full (3 items), adding 4th...")
	cache.Set("user:4", "Diana")
	fmt.Println("Added user:4")

	for _, key := range []string{"user:1", "user:2", "user:3", "user:4"} {
		if val, ok := cache.Get(key); ok {
			fmt.Printf("  %s = %s (found)\n", key, val)
		} else {
			fmt.Printf("  %s = NOT FOUND (evicted)\n", key)
		}
	}
	fmt.Println()

	// Demo 3: TTL expiration
	fmt.Println("--- Demo 3: TTL Expiration ---")
	shortCache := genericlrucache.New[string, int](10)
	shortCache.Set("temp", 42)
	if val, ok := shortCache.Get("temp"); ok {
		fmt.Printf("  Before TTL: temp = %d\n", val)
	}

	fmt.Println("  Waiting 150ms for TTL to expire...")
	time.Sleep(150 * time.Millisecond)

	if _, ok := shortCache.Get("temp"); !ok {
		fmt.Println("  After TTL: temp = NOT FOUND (expired)")
	}
	fmt.Println()

	// Demo 4: Generic types
	fmt.Println("--- Demo 4: Different Types ---")
	intCache := genericlrucache.New[int, string](5)
	intCache.Set(1, "first")
	intCache.Set(2, "second")
	if val, ok := intCache.Get(1); ok {
		fmt.Printf("  intCache[1] = %s\n", val)
	}
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Go generics: [K comparable, V any]")
	fmt.Println("2. LRU = doubly-linked list + hash map")
	fmt.Println("3. sync.Mutex for thread safety")
	fmt.Println("4. TTL expiration with time.Time")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/08-http-client-retries")
}
