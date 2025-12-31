// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/syncmutexvsrwmutex.go
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
	"time"
)

func main() {
	fmt.Println("=== Debugging Mutex vs RWMutex ===")

	// Fixed default values - modify these directly if you want to test different inputs
	var mu sync.Mutex
	var counter int

	fmt.Println("Testing Mutex...")

	// Multiple goroutines incrementing counter
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("Counter value: %d (expected: 5)\n", counter)
	fmt.Println("\nNote: Use cmd/app/main.go for comprehensive mutex vs rwmutex demos.")
}

