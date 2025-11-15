//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for channels.

package exercise

import (
	"context"
	"time"
)

// Ping creates a channel and sends a single value, then closes it.
//
// REQUIREMENTS:
// - Create a buffered channel with capacity 1
// - Send the value to the channel
// - Close the channel
// - Return the channel
//
// EXAMPLE:
//   ch := Ping(42)
//   value := <-ch  // 42
//   _, ok := <-ch  // false (channel closed)
//
// HINT: Use make(chan int, 1) to create a buffered channel
func Ping(value int) <-chan int {
	// TODO: Implement this
	return nil
}

// PingPong creates two channels that play ping-pong n times.
//
// REQUIREMENTS:
// - Create two channels: ping and pong
// - Launch a goroutine that receives from ping, sends to pong (n times)
// - Return both channels
// - Close channels when done
//
// EXAMPLE:
//   ping, pong := PingPong(3)
//   ping <- 1       // Start the game
//   <-pong          // Receive from pong
//   ping <- 2
//   <-pong
//   // ... continues n times
//
// HINT: Use a goroutine to alternate between receiving and sending
func PingPong(n int) (chan<- int, <-chan int) {
	// TODO: Implement this
	return nil, nil
}

// Merge combines multiple input channels into a single output channel.
//
// REQUIREMENTS:
// - Create an output channel
// - Launch a goroutine for each input that forwards values to output
// - Close output when ALL inputs are closed
// - Return the output channel
//
// EXAMPLE:
//   ch1 := make(chan int)
//   ch2 := make(chan int)
//   merged := Merge(ch1, ch2)
//   // Values from ch1 and ch2 appear in merged
//
// HINT: Use sync.WaitGroup to wait for all input goroutines to finish
// HINT: Close output in a separate goroutine after WaitGroup.Wait()
func Merge(channels ...<-chan int) <-chan int {
	// TODO: Implement this
	return nil
}

// Filter creates a channel that only forwards values matching the predicate.
//
// REQUIREMENTS:
// - Create an output channel
// - Launch a goroutine that reads from input
// - Only send values where predicate(value) is true
// - Close output when input is closed
// - Return the output channel
//
// EXAMPLE:
//   input := make(chan int)
//   evens := Filter(input, func(x int) bool { return x%2 == 0 })
//   // Only even numbers from input appear in evens
//
// HINT: Use range to read from input until it's closed
func Filter(input <-chan int, predicate func(int) bool) <-chan int {
	// TODO: Implement this
	return nil
}

// Map creates a channel that transforms values using a function.
//
// REQUIREMENTS:
// - Create an output channel
// - Launch a goroutine that reads from input
// - Apply the transform function to each value
// - Send transformed values to output
// - Close output when input is closed
// - Return the output channel
//
// EXAMPLE:
//   input := make(chan int)
//   doubled := Map(input, func(x int) int { return x * 2 })
//   // Each value from input is doubled in the output
//
// HINT: Similar to Filter, but always send (after transforming)
func Map(input <-chan int, transform func(int) int) <-chan int {
	// TODO: Implement this
	return nil
}

// Take creates a channel that forwards at most n values from input.
//
// REQUIREMENTS:
// - Create an output channel
// - Launch a goroutine that reads from input
// - Forward at most n values to output
// - Close output after n values or when input closes (whichever comes first)
// - Return the output channel
//
// EXAMPLE:
//   input := make(chan int)
//   first5 := Take(input, 5)
//   // Only the first 5 values from input appear in first5
//
// HINT: Use a counter and break after n values
func Take(input <-chan int, n int) <-chan int {
	// TODO: Implement this
	return nil
}

// OrDone wraps a channel and adds cancellation via context.
//
// REQUIREMENTS:
// - Create an output channel
// - Launch a goroutine that reads from input
// - Forward values to output
// - Stop and close output when ctx.Done() is signaled OR input closes
// - Return the output channel
//
// EXAMPLE:
//   ctx, cancel := context.WithCancel(context.Background())
//   input := make(chan int)
//   output := OrDone(ctx, input)
//   cancel()  // Stops forwarding, closes output
//
// HINT: Use select to monitor both input and ctx.Done()
func OrDone(ctx context.Context, input <-chan int) <-chan int {
	// TODO: Implement this
	return nil
}

// Tee splits an input channel into two output channels.
//
// REQUIREMENTS:
// - Create two output channels
// - Launch a goroutine that reads from input
// - Send each value to BOTH output channels
// - Close both outputs when input is closed
// - Return both output channels
//
// EXAMPLE:
//   input := make(chan int)
//   out1, out2 := Tee(input)
//   // Each value from input is sent to both out1 and out2
//
// HINT: For each value, send to out1 then out2 (or use select)
func Tee(input <-chan int) (<-chan int, <-chan int) {
	// TODO: Implement this
	return nil, nil
}

// Bridge flattens a channel of channels into a single channel.
//
// REQUIREMENTS:
// - Create an output channel
// - Launch a goroutine that reads channels from input
// - For each channel received, read all its values and send to output
// - Close output when input is closed and all sub-channels are drained
// - Return the output channel
//
// EXAMPLE:
//   input := make(chan <-chan int)
//   flattened := Bridge(input)
//   // Values from all channels sent to input appear in flattened
//
// HINT: Use nested loops - outer for channels, inner for values
func Bridge(input <-chan (<-chan int)) <-chan int {
	// TODO: Implement this
	return nil
}

// Debounce creates a channel that only forwards values if no new value
// arrives within the specified duration.
//
// REQUIREMENTS:
// - Create an output channel
// - Launch a goroutine that reads from input
// - Use a timer to delay forwarding
// - If a new value arrives before the timer fires, reset the timer
// - Only send a value when the timer fires (no new value arrived)
// - Close output when input is closed
// - Return the output channel
//
// EXAMPLE:
//   input := make(chan int)
//   debounced := Debounce(input, 100*time.Millisecond)
//   // Only forwards values that have 100ms of "quiet time" after them
//
// HINT: Use time.NewTimer and select with timer.C
func Debounce(input <-chan int, duration time.Duration) <-chan int {
	// TODO: Implement this
	return nil
}

// NewBoundedQueue creates a queue with a maximum capacity.
//
// REQUIREMENTS:
// - Create a BoundedQueue with a buffered channel of size capacity
// - Return a pointer to the queue
//
// EXAMPLE:
//   queue := NewBoundedQueue(10)
//   queue.Enqueue(42)  // Non-blocking (until capacity reached)
//   val := queue.Dequeue()  // 42
func NewBoundedQueue(capacity int) *BoundedQueue {
	// TODO: Implement this
	return nil
}

// Enqueue adds a value to the queue (blocks if full).
func (q *BoundedQueue) Enqueue(value int) {
	// TODO: Implement this
}

// Dequeue removes and returns a value from the queue (blocks if empty).
func (q *BoundedQueue) Dequeue() int {
	// TODO: Implement this
	return 0
}

// TryEnqueue attempts to add a value without blocking.
// Returns true if successful, false if queue is full.
func (q *BoundedQueue) TryEnqueue(value int) bool {
	// TODO: Implement this
	return false
}

// TryDequeue attempts to remove a value without blocking.
// Returns (value, true) if successful, (0, false) if queue is empty.
func (q *BoundedQueue) TryDequeue() (int, bool) {
	// TODO: Implement this
	return 0, false
}

// NewBroadcaster creates a new broadcaster.
//
// REQUIREMENTS:
// - Initialize the broadcaster's fields
// - Launch a goroutine that reads from input and sends to all listeners
// - Return a pointer to the broadcaster
//
// EXAMPLE:
//   bc := NewBroadcaster()
//   ch1 := bc.Subscribe()
//   ch2 := bc.Subscribe()
//   bc.Send(Message{ID: 1, Content: "hello"})
//   // Both ch1 and ch2 receive the message
func NewBroadcaster() *Broadcaster {
	// TODO: Implement this
	return nil
}

// Subscribe adds a new listener and returns its channel.
func (b *Broadcaster) Subscribe() <-chan Message {
	// TODO: Implement this
	return nil
}

// Unsubscribe removes a listener.
func (b *Broadcaster) Unsubscribe(ch <-chan Message) {
	// TODO: Implement this
}

// Send broadcasts a message to all subscribers.
func (b *Broadcaster) Send(msg Message) {
	// TODO: Implement this
}

// Close stops the broadcaster.
func (b *Broadcaster) Close() {
	// TODO: Implement this
}

// NewBarrier creates a barrier for n goroutines.
//
// REQUIREMENTS:
// - Initialize a barrier that waits for n goroutines
// - Return a pointer to the barrier
//
// EXAMPLE:
//   barrier := NewBarrier(3)
//   // Launch 3 goroutines that call barrier.Wait()
//   // All 3 block until all have called Wait(), then all proceed
func NewBarrier(n int) *Barrier {
	// TODO: Implement this
	return nil
}

// Wait blocks until all n goroutines have called Wait.
func (b *Barrier) Wait() {
	// TODO: Implement this
}
