// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/atomiccountersvsmutex.go
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
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("=== Debugging Atomic Counters vs Mutex ===")

	// Fixed default values - modify these directly if you want to test different inputs
	var atomicCounter int64
	var mutexCounter int64
	var mu sync.Mutex

	fmt.Println("Testing atomic vs mutex counter increment...")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Atomic increment
			atomic.AddInt64(&atomicCounter, 1)
			// Mutex increment
			mu.Lock()
			mutexCounter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("Atomic counter: %d\n", atomic.LoadInt64(&atomicCounter))
	fmt.Printf("Mutex counter: %d\n", mutexCounter)
	fmt.Println("\nNote: Use cmd/app/main.go for comprehensive atomic vs mutex demos.")
}

