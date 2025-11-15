package exercise

import (
	"sync"
)

// SafeCounter is a thread-safe counter that students need to implement.
type SafeCounter struct {
	// TODO: Add necessary fields for thread-safe counter
	// Hint: You can use either sync.Mutex or atomic.Int64
}

// SafeMap is a thread-safe map wrapper.
type SafeMap struct {
	// TODO: Add necessary fields for thread-safe map
	// Hint: Use sync.RWMutex and a regular map
}

// LazyInit demonstrates lazy initialization pattern.
type LazyInit struct {
	// TODO: Add necessary fields for thread-safe lazy initialization
	// Hint: Use sync.Once
}

// SafeSlice is a thread-safe slice wrapper.
type SafeSlice struct {
	// TODO: Add necessary fields for thread-safe slice
	// Hint: Use sync.Mutex and a regular slice
}

// URLCache is a concurrent URL fetcher with caching.
type URLCache struct {
	cache   map[string]string
	mu      sync.RWMutex
	fetcher func(url string) (string, error)
}

// Metrics tracks application metrics concurrently.
type Metrics struct {
	// TODO: Add necessary fields for concurrent metrics
	// Hint: Use atomic types for counters
}

// BankAccount simulates a bank account with concurrent deposits/withdrawals.
type BankAccount struct {
	// TODO: Add necessary fields for thread-safe bank account
	// Hint: Use sync.Mutex to protect balance
}
