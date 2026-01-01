package racedetectiondemo

import (
	"sync"
	"sync/atomic"
)

// SafeCounter is a thread-safe counter that students need to implement.
type SafeCounter struct {
	value atomic.Int64
}

// SafeMap is a thread-safe map wrapper.
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

// LazyInit demonstrates lazy initialization pattern.
type LazyInit struct {
	once  sync.Once
	value any
}

// SafeSlice is a thread-safe slice wrapper.
type SafeSlice struct {
	mu   sync.RWMutex
	data []int
}

// URLCache is a concurrent URL fetcher with caching.
type URLCache struct {
	cache   map[string]string
	mu      sync.RWMutex
	fetcher func(url string) (string, error)
}

// Metrics tracks application metrics concurrently.
type Metrics struct {
	requests atomic.Int64
	errors   atomic.Int64
}

// BankAccount simulates a bank account with concurrent deposits/withdrawals.
type BankAccount struct {
	mu      sync.Mutex
	balance int64
}
