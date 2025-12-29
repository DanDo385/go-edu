// Package main declares this as an executable program (not a library)
// When you run "go run" or "go build", Go looks for a main package with a main() function
package main

import (
	"flag" // Command-line argument parsing
	"fmt"  // Formatted I/O
	"time" // Time operations

	// Import our exercise package containing the Cache type and New function
	"github.com/example/go-10x-minis/minis/07-generic-lru-cache/exercise"
)

/*
===================================================================
Generic LRU Cache with TTL Demonstration
===================================================================

This program demonstrates a generic Least Recently Used (LRU) cache
implementation with Time-To-Live (TTL) expiration.

USAGE:
    go run ./cmd/generic-lru-cache/main.go
    go run ./cmd/generic-lru-cache/main.go -capacity=5 -ttl=5s
    go run ./cmd/generic-lru-cache/main.go -capacity=10 -ttl=10s -demo=eviction
    go run ./cmd/generic-lru-cache/main.go -capacity=2 -ttl=1s -demo=expiration

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see cache state changes
- Inspect cache internals to see LRU ordering and TTL timestamps
- See /RUN_DEBUG.md for comprehensive debugging guide

CACHE BEHAVIOR:
- LRU eviction: When cache is full, least recently used item is removed
- TTL expiration: Items expire after the specified duration
- Generic implementation: Works with any comparable key and value types
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect default values before parsing
	// DEBUG: Before flag.Parse(), all flags have their default values
	// DEBUG: In Variables panel, you'll see: capacity=3, ttl=2s, demo="all"

	capacity := flag.Int("capacity", 3, "Maximum number of items in cache")
	ttl := flag.Duration("ttl", 2*time.Second, "Time-to-live for cache entries")
	demo := flag.String("demo", "all", "Which demo to run: all, eviction, expiration")

	flag.Parse()

	// BREAKPOINT: Set breakpoint here to inspect parsed values
	// DEBUG: After flag.Parse(), variables now hold actual command-line values
	// DEBUG: Watch 'capacity', 'ttl', 'demo' - they may have changed from defaults
	// DEBUG: Hover over *capacity to see the dereferenced int value

	// ============================================================================
	// STEP 2: Create Cache Instance
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before cache creation
	// DEBUG: Step Into (F11) on exercise.New() to see cache initialization
	// DEBUG: Watch how the cache struct is allocated and initialized

	cache := exercise.New[string, int](*capacity, *ttl)

	// BREAKPOINT: Set breakpoint here after cache creation
	// DEBUG: Inspect 'cache' variable in Variables panel
	// DEBUG: Expand cache to see internal fields: capacity, ttl, items map, list
	// DEBUG: Initially, cache should be empty with size 0

	fmt.Println("=== LRU Cache Demo ===\n")
	fmt.Printf("Configuration: capacity=%d, ttl=%v\n\n", *capacity, *ttl)

	// ============================================================================
	// STEP 3: Basic Operations Demo
	// ============================================================================
	if *demo == "all" || *demo == "basic" {
		fmt.Println("--- Basic Operations ---")

		// BREAKPOINT: Set breakpoint here before first Set operation
		// DEBUG: Step Into (F11) on cache.Set() to see how items are stored
		// DEBUG: Watch how the internal map and doubly-linked list are updated

		cache.Set("a", 1)

		// BREAKPOINT: Set breakpoint here after first Set
		// DEBUG: Inspect cache.Len() - should be 1
		// DEBUG: Expand cache to see "a" in the items map

		cache.Set("b", 2)
		cache.Set("c", 3)

		// BREAKPOINT: Set breakpoint here after filling cache to capacity
		// DEBUG: Check cache.Len() - should equal *capacity
		// DEBUG: Inspect the doubly-linked list order (most recent to least recent)

		fmt.Println("Added: a=1, b=2, c=3")
		fmt.Printf("Size: %d/%d\n\n", cache.Len(), *capacity)
	}

	// ============================================================================
	// STEP 4: LRU Access Pattern
	// ============================================================================
	if *demo == "all" || *demo == "access" {
		fmt.Println("--- LRU Access Pattern ---")

		// BREAKPOINT: Set breakpoint here before Get operation
		// DEBUG: Step Into (F11) on cache.Get() to see how access updates LRU order
		// DEBUG: Watch how accessing an item moves it to the front of the list

		if val, ok := cache.Get("a"); ok {
			// BREAKPOINT: Set breakpoint here after successful Get
			// DEBUG: Inspect 'val' - should be 1
			// DEBUG: Inspect 'ok' - should be true
			// DEBUG: Check the list order - "a" should now be at the front

			fmt.Printf("Get 'a': %d (moved to front)\n", val)
		}

		// BREAKPOINT: Set breakpoint here after access
		// DEBUG: In cache internals, verify "a" is now the most recently used
		// DEBUG: The LRU order should now be: a -> c -> b (or similar)
	}

	// ============================================================================
	// STEP 5: Eviction Demo
	// ============================================================================
	if *demo == "all" || *demo == "eviction" {
		fmt.Println("\n--- Eviction Demo ---")

		// BREAKPOINT: Set breakpoint here before eviction trigger
		// DEBUG: Note the current size and which item is least recently used
		// DEBUG: Step Into (F11) on cache.Set() to see eviction logic

		cache.Set("d", 4)

		// BREAKPOINT: Set breakpoint here after adding item that triggers eviction
		// DEBUG: Cache size should still be at capacity
		// DEBUG: The least recently used item should have been evicted
		// DEBUG: Check which key was removed from the items map

		fmt.Println("Added: d=4")
		fmt.Printf("Size: %d/%d\n", cache.Len(), *capacity)

		// BREAKPOINT: Set breakpoint here before checking evicted item
		// DEBUG: Step Into (F11) on cache.Get() to see cache miss logic

		if _, ok := cache.Get("b"); !ok {
			// BREAKPOINT: Set breakpoint here when item not found
			// DEBUG: 'ok' should be false - item was evicted
			// DEBUG: This confirms LRU eviction is working correctly

			fmt.Println("'b' was evicted (LRU)\n")
		}
	}

	// ============================================================================
	// STEP 6: TTL Expiration Demo
	// ============================================================================
	if *demo == "all" || *demo == "expiration" {
		fmt.Println("\n--- TTL Expiration Demo ---")

		// BREAKPOINT: Set breakpoint here before sleep
		// DEBUG: Note the current time and cache contents
		// DEBUG: All items should have timestamps from when they were Set

		waitTime := *ttl + 1*time.Second
		fmt.Printf("Waiting %v for TTL expiration...\n", waitTime)

		// BREAKPOINT: Set breakpoint here during sleep (won't pause, but shows timing)
		// DEBUG: Use time.Now() in Debug Console to check elapsed time

		time.Sleep(waitTime)

		// BREAKPOINT: Set breakpoint here after sleep
		// DEBUG: Items should now be expired (created more than ttl ago)
		// DEBUG: Step Into (F11) on cache.Get() to see expiration check

		fmt.Println("Checking for expired items:")

		// BREAKPOINT: Set breakpoint here before first expiration check
		// DEBUG: Step Into (F11) to see how Get checks TTL

		if _, ok := cache.Get("a"); !ok {
			// BREAKPOINT: Set breakpoint here when expired item not found
			// DEBUG: 'ok' should be false - item expired
			// DEBUG: The item should have been removed from cache

			fmt.Println("'a' expired (TTL)")
		}

		// BREAKPOINT: Set breakpoint before checking remaining items
		if _, ok := cache.Get("c"); !ok {
			fmt.Println("'c' expired (TTL)")
		}
		if _, ok := cache.Get("d"); !ok {
			fmt.Println("'d' expired (TTL)")
		}

		// BREAKPOINT: Set breakpoint here after all expiration checks
		// DEBUG: cache.Len() should be 0 or very small (only non-expired items)
		// DEBUG: All items created before sleep should be gone

		fmt.Printf("\nFinal size: %d\n", cache.Len())
	}

	// ============================================================================
	// STEP 7: Advanced Demo - Multiple Operations
	// ============================================================================
	if *demo == "all" || *demo == "advanced" {
		fmt.Println("\n--- Advanced Demo: Mixed Operations ---")

		// Reset cache for clean demo
		cache = exercise.New[string, int](*capacity, *ttl)

		operations := []struct {
			op    string
			key   string
			value int
		}{
			{"set", "x", 10},
			{"set", "y", 20},
			{"get", "x", 0},
			{"set", "z", 30},
			{"set", "w", 40},
		}

		// BREAKPOINT: Set breakpoint here before operation loop
		// DEBUG: Expand 'operations' to see the sequence of operations
		// DEBUG: Watch how each operation affects cache state

		for i, op := range operations {
			// BREAKPOINT: Set breakpoint here inside loop
			// DEBUG: Watch 'i' and 'op' change with each iteration
			// DEBUG: Step Into (F11) to see operation implementation

			switch op.op {
			case "set":
				cache.Set(op.key, op.value)
				fmt.Printf("%d. Set %s=%d (size: %d)\n", i+1, op.key, op.value, cache.Len())

				// BREAKPOINT: Set breakpoint after Set
				// DEBUG: Check if eviction occurred
				// DEBUG: Inspect LRU ordering

			case "get":
				if val, ok := cache.Get(op.key); ok {
					fmt.Printf("%d. Get %s=%d (moved to front)\n", i+1, op.key, val)
				} else {
					fmt.Printf("%d. Get %s=<not found>\n", i+1, op.key)
				}

				// BREAKPOINT: Set breakpoint after Get
				// DEBUG: Verify the accessed item moved to front
			}
		}

		// BREAKPOINT: Set breakpoint here after all operations
		// DEBUG: Final cache state - which items remain?
		// DEBUG: What is the LRU order?

		fmt.Printf("Final cache size: %d/%d\n", cache.Len(), *capacity)
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise package functions
	// 5. Watch Variables panel to see cache state changes
	// 6. Add watch expressions: cache.Len(), *capacity, *ttl
	// 7. Use Call Stack panel to see function hierarchy
	//
	// CACHE-SPECIFIC DEBUGGING:
	// - Expand cache variable to see internal structure
	// - Watch the items map to see key-value pairs
	// - Inspect the doubly-linked list for LRU ordering
	// - Check timestamps to verify TTL logic
	// - Use Debug Console: cache.Len(), time.Now()
	//
	// UNDERSTANDING LRU:
	// - Most recently used items are at the front of the list
	// - Least recently used items are at the back
	// - When cache is full, back item is evicted
	// - Get and Set operations move items to front
	//
	// UNDERSTANDING TTL:
	// - Each item has a timestamp when it was set
	// - On Get, check if (now - timestamp) > ttl
	// - Expired items are removed during Get
	// - Use time.Now() in Debug Console to check current time
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}
