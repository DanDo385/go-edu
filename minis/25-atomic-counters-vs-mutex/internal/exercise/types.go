package exercise

import (
	"sync/atomic"
	"unsafe"
)

// AtomicCounter represents a thread-safe counter using atomic operations.
type AtomicCounter struct {
	value int64
}

// AtomicFlag represents a thread-safe boolean flag.
type AtomicFlag struct {
	value int64 // 0 = false, 1 = true
}

// RateLimiter implements a token bucket rate limiter using atomic operations.
type RateLimiter struct {
	tokens     int64
	maxTokens  int64
	refillRate int64 // tokens per second
	lastRefill int64 // Unix timestamp in seconds
}

// AtomicMax tracks the maximum value seen using atomic operations.
type AtomicMax struct {
	max int64
}

// LockFreeQueue represents a simple lock-free queue.
type LockFreeQueue struct {
	head unsafe.Pointer // *queueNode
	tail unsafe.Pointer // *queueNode
}

// queueNode is an internal node for the lock-free queue.
type queueNode struct {
	value int
	next  unsafe.Pointer // *queueNode
}

// SpinLock implements a simple spinlock using atomic operations.
type SpinLock struct {
	state int64 // 0 = unlocked, 1 = locked
}

// AtomicState represents a state machine with atomic transitions.
type AtomicState struct {
	current int64
}

// State constants for AtomicState
const (
	StateIdle = iota
	StateRunning
	StatePaused
	StateStopped
)

// ReferenceCounter tracks references to a resource.
type ReferenceCounter struct {
	count int64
}

// ConfigManager manages configuration updates atomically.
type ConfigManager struct {
	config atomic.Value // stores *Config
}

// Config represents application configuration.
type Config struct {
	Timeout    int
	MaxRetries int
	BatchSize  int
}

// AtomicBitmap represents a bitmap with atomic bit operations.
type AtomicBitmap struct {
	bits [4]int64 // 256 bits (4 * 64)
}

// LoadBalancer distributes work across workers using atomic operations.
type LoadBalancer struct {
	counter int64
	workers int64
}

// CircularBuffer implements a lock-free circular buffer.
type CircularBuffer struct {
	buffer   []int64
	head     int64
	tail     int64
	size     int64
	capacity int64
}
