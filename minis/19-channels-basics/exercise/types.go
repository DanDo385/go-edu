package exercise

import (
	"sync"
)

// Message represents a message with an ID and content.
type Message struct {
	ID      int
	Content string
}

// Result represents a computation result.
type Result struct {
	Value int
	Err   error
}

// Task represents a work item to be processed.
type Task struct {
	ID   int
	Data int
}

// Pipeline represents a multi-stage processing pipeline.
type Pipeline struct {
	stages []func(in <-chan int) <-chan int
}

// BoundedQueue represents a thread-safe queue with maximum capacity.
type BoundedQueue struct {
	ch chan int
}

// Broadcaster allows multiple receivers to receive the same messages.
type Broadcaster struct {
	mu        sync.RWMutex
	listeners []chan Message
	input     chan Message
	done      chan struct{}
}

// Barrier synchronizes multiple goroutines at a point.
type Barrier struct {
	n       int
	count   int
	ch      chan struct{}
	mu      sync.Mutex
	waiting chan struct{}
}
