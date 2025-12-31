//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for channels.

package channelsbasics

import (
	"context"
	"time"
)

// Ping creates a channel and sends a single value, then closes it.
func Ping(value int) <-chan int {
	// TODO: Implement Ping
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// PingPong creates two channels that play ping-pong n times.
func PingPong(n int) (chan<- int, <-chan int) {
	// TODO: Implement PingPong
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Merge combines multiple input channels into a single output channel.
func Merge(channels ...<-chan int) <-chan int {
	// TODO: Implement Merge
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Filter creates a channel that only forwards values matching the predicate.
func Filter(input <-chan int, predicate func(int) bool) <-chan int {
	// TODO: Implement Filter
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Map creates a channel that transforms values using a function.
func Map(input <-chan int, transform func(int) int) <-chan int {
	// TODO: Implement Map
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Take creates a channel that forwards at most n values from input.
func Take(input <-chan int, n int) <-chan int {
	// TODO: Implement Take
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// OrDone wraps a channel and adds cancellation via context.
func OrDone(ctx context.Context, input <-chan int) <-chan int {
	// TODO: Implement OrDone
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Tee splits an input channel into two output channels.
func Tee(input <-chan int) (<-chan int, <-chan int) {
	// TODO: Implement Tee
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Bridge flattens a channel of channels into a single channel.
func Bridge(input <-chan (<-chan int)) <-chan int {
	// TODO: Implement Bridge
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Debounce creates a channel that only forwards values if no new value
// arrives within the specified duration.
func Debounce(input <-chan int, duration time.Duration) <-chan int {
	// TODO: Implement Debounce
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// NewBoundedQueue creates a queue with a maximum capacity.
func NewBoundedQueue(capacity int) *BoundedQueue {
	// TODO: Implement NewBoundedQueue
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Enqueue adds a value to the queue (blocks if full).
func (q *BoundedQueue) Enqueue(value int) {
	// TODO: Implement this.
	// - Simply send the value to the channel: `q.ch <- value`.
	// - If the channel's buffer is full, this operation will block until another goroutine dequeues a value, making space.
}

// Dequeue removes and returns a value from the queue (blocks if empty).
func (q *BoundedQueue) Dequeue() int {
	// TODO: Implement this.
	// - Simply receive a value from the channel: `return <-q.ch`.
	// - If the channel's buffer is empty, this operation will block until another goroutine enqueues a value.
	return 0
}

// TryEnqueue attempts to add a value without blocking.
// Returns true if successful, false if queue is full.
func (q *BoundedQueue) TryEnqueue(value int) bool {
	// TODO: Implement this.

	// This demonstrates non-blocking channel operations.

	// Step 1: Use a `select` statement.
	// - `case q.ch <- value:`
	//   - If the send succeeds immediately (because the buffer is not full), this case is chosen. Return `true`.
	// - `default:`
	//   - If the send would block (because the buffer is full), the `default` case is chosen immediately. Return `false`.
	return false
}

// TryDequeue attempts to remove a value without blocking.
// Returns (value, true) if successful, (0, false) if queue is empty.
func (q *BoundedQueue) TryDequeue() (int, bool) {
	// TODO: Implement this.

	// This is the non-blocking receive equivalent of `TryEnqueue`.

	// Step 1: Use a `select` statement.
	// - `case v := <-q.ch:`
	//   - If the receive succeeds immediately (because the buffer is not empty), this case is chosen. Return `v, true`.
	// - `default:`
	//   - If the receive would block (because the buffer is empty), the `default` case is chosen. Return `0, false`.
	return 0, false
}

// NewBroadcaster creates a new broadcaster.
func NewBroadcaster() *Broadcaster {
	// TODO: Implement NewBroadcaster
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Subscribe adds a new listener and returns its channel.
func (b *Broadcaster) Subscribe() <-chan Message {
	// TODO: Implement this.

	// Step 1: Create a new channel for the subscriber.
	// Step 2: Lock the mutex (`b.mu.Lock()`) because you are modifying the `listeners` slice.
	// Step 3: Append the new channel to the `b.listeners` slice.
	// Step 4: Unlock the mutex (`b.mu.Unlock()`).
	// Step 5: Return the new channel.
	return nil
}

// Unsubscribe removes a listener.
func (b *Broadcaster) Unsubscribe(ch <-chan Message) {
	// TODO: Implement this.

	// Step 1: Lock the mutex.
	// Step 2: Find the channel in the `listeners` slice.
	// Step 3: Remove it from the slice. A common way to do this is `b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)`.
	// Step 4: Close the channel that is being removed.
	// Step 5: Unlock the mutex.
}

// Send broadcasts a message to all subscribers.
func (b *Broadcaster) Send(msg Message) {
	// TODO: Implement this.
	// - Simply send the message to the `input` channel. The central goroutine will handle the rest.
}

// Close stops the broadcaster.
func (b *Broadcaster) Close() {
	// TODO: Implement this.
	// - Simply close the `done` channel. This will signal the central goroutine to shut down.
}

// NewBarrier creates a barrier for n goroutines.
func NewBarrier(n int) *Barrier {
	// TODO: Implement NewBarrier
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Wait blocks until all n goroutines have called Wait.
func (b *Barrier) Wait() {
	// TODO: Implement this.

	// This is the core logic of the barrier.

	// Step 1: Lock the mutex to safely modify the count.
	// - `b.mu.Lock()`

	// Step 2: Increment the count of waiting goroutines.
	// - `b.count++`

	// Step 3: Check if this is the last goroutine.
	// - `if b.count < b.n { ... }` (The "early" goroutines)
	//   - If it's not the last one, this goroutine must wait.
	//   - **Crucially**, you must unlock the mutex *before* waiting on the channel, otherwise the last goroutine will never be able to acquire the lock. `b.mu.Unlock()`.
	//   - Wait on the `waiting` channel: `<-b.waiting`. This will block.
	//   - After being unblocked, `return`.
	//
	// - `else { ... }` (The "last" goroutine)
	//   - If this is the last goroutine, it's time to open the gate for everyone else.
	//   - `close(b.waiting)`. Closing a channel unblocks all goroutines that are currently waiting to receive from it.
	//   - Reset the barrier for the next use: set `b.count = 0` and create a new `waiting` channel.
	//   - Unlock the mutex: `b.mu.Unlock()`.

	// Memory & Language Comparison:
	// - Go: The standard library doesn't have a barrier primitive, so this is a common pattern to implement one using mutexes and channels.
	// - Java: `java.util.concurrent.CyclicBarrier` is a standard library class that provides this exact functionality.
	// - C++: `std::barrier` was introduced in C++20.
	// - Python: `threading.Barrier` provides this functionality.
}
