//go:build reference

package channelsbasics

/*
Problem: Master Go channels for goroutine communication and synchronization

Requirements:
1. Understand channel creation (buffered vs unbuffered)
2. Channel direction types (send-only, receive-only)
3. Channel closing and range iteration
4. Channel-based patterns (merge, filter, map, tee, bridge)
5. Building concurrent data structures (queues, barriers, broadcasters)
6. Context cancellation patterns with channels

Data Structure:
- Channels: Typed conduits for goroutine communication
- Unbuffered: Synchronous send/receive (rendezvous)
- Buffered: Asynchronous with capacity (non-blocking until full)
- Direction: chan T (bidirectional), chan<- T (send-only), <-chan T (receive-only)

Time/Space Complexity:
- Channel operations: O(1) (constant time)
- Merge/Fan-in: O(1) per value (parallel forwarding)
- Filter/Map: O(1) per value (pipelined processing)
- Space: O(buffer capacity) per channel

Algorithm: Channel Communication Patterns
- Producer/Consumer: One goroutine sends, another receives
- Fan-in: Multiple channels → single channel (merge)
- Fan-out: Single channel → multiple channels (tee, broadcast)
- Pipeline: Chain of processing stages
- Generator: Function returning channel (generator pattern)

Why channels are powerful:
- Built-in synchronization (no explicit locks needed)
- Type-safe communication (compile-time guarantees)
- Select statement for multiplexing
- Context cancellation support
- Idiomatic Go concurrency model

Critical Best Practices:
1. Always close channels to signal completion (prevents goroutine leaks)
2. Use buffered channels when producer/consumer speeds differ
3. Range over channels (automatically stops on close)
4. Use select for non-blocking operations and timeouts
5. Never send on closed channel (panic)
6. Receiving from closed channel returns zero value immediately
*/

import (
	"context"
	"sync"
	"time"
)

// ============================================================================
// Exercise 1: Basic Channel Operations
// ============================================================================

// Ping creates a channel, sends a value, and closes it.
//
// Pattern: Generator function (returns channel for lazy evaluation)
// Use Case: Create one-shot value channels, lazy initialization
//
// BREAKPOINT: Set breakpoint here to trace channel creation
// DEBUG: Watch 'value' parameter
func Ping(value int) <-chan int {
	// BREAKPOINT: Set breakpoint here before channel creation
	// DEBUG: Buffered channel (capacity 1) allows non-blocking send
	ch := make(chan int, 1)
	// BREAKPOINT: Set breakpoint here before send
	// DEBUG: Send to buffered channel (non-blocking if space available)
	ch <- value
	// BREAKPOINT: Set breakpoint here before close
	// DEBUG: Close channel to signal no more values (important for receivers)
	close(ch)
	// DEBUG: Return receive-only channel type (<-chan int)
	return ch
}

// PingPong creates two channels that pass values back and forth n times.
//
// Pattern: Bidirectional communication between goroutines
// Use Case: Request-response patterns, cooperative algorithms
//
// BREAKPOINT: Set breakpoint here to trace ping-pong setup
// DEBUG: Watch 'n' parameter (number of rounds)
func PingPong(n int) (chan<- int, <-chan int) {
	// BREAKPOINT: Set breakpoint here before channel creation
	// DEBUG: Unbuffered channels = synchronous communication
	ping := make(chan int)
	pong := make(chan int)

	// BREAKPOINT: Set breakpoint here before goroutine launch
	go func() {
		// BREAKPOINT: Set breakpoint in goroutine to trace ping-pong logic
		defer close(pong) // Ensure channel is closed even on panic
		// DEBUG: Loop for n iterations
		for i := 0; i < n; i++ {
			// BREAKPOINT: Set breakpoint here to see value received from ping
			// DEBUG: Block until value arrives from ping channel
			value := <-ping
			// BREAKPOINT: Set breakpoint here to see value sent to pong
			// DEBUG: Forward value to pong channel
			pong <- value
		}
	}()

	// BREAKPOINT: Set breakpoint here before return
	// DEBUG: Return send-only ping (chan<- int) and receive-only pong (<-chan int)
	return ping, pong
}

// ============================================================================
// Exercise 2: Channel Transformation Patterns
// ============================================================================

// Merge combines multiple input channels into a single output channel (Fan-in pattern).
//
// Pattern: Fan-in (multiple producers → single consumer)
// Use Case: Aggregating results from parallel workers, combining data streams
// Complexity: O(1) per value (parallel forwarding via goroutines)
//
// BREAKPOINT: Set breakpoint here to trace merge operation
// DEBUG: Watch 'channels' variadic parameter
func Merge(channels ...<-chan int) <-chan int {
	// BREAKPOINT: Set breakpoint here before output channel creation
	output := make(chan int)
	var wg sync.WaitGroup // Track when all forwarders complete

	// BREAKPOINT: Set breakpoint here before launching forwarder goroutines
	// Launch a goroutine for each input channel
	for _, ch := range channels {
		// BREAKPOINT: Set breakpoint here in loop to see each channel
		wg.Add(1)
		go func(c <-chan int) {
			// BREAKPOINT: Set breakpoint in forwarder goroutine
			defer wg.Done()
			// DEBUG: Range automatically stops when channel closes
			for v := range c {
				// BREAKPOINT: Set breakpoint here to see values forwarded
				// DEBUG: Forward each value from input to output
				output <- v
			}
		}(ch) // Important: pass ch as parameter (closure capture gotcha)
	}

	// BREAKPOINT: Set breakpoint here before closer goroutine
	// Close output when all inputs are drained
	go func() {
		// BREAKPOINT: Set breakpoint in closer goroutine
		// DEBUG: Wait for all forwarders to complete
		wg.Wait()
		// BREAKPOINT: Set breakpoint here before close
		// DEBUG: Close output channel (signals completion to receivers)
		close(output)
	}()

	return output
}

// Filter creates a channel that only forwards values matching the predicate.
//
// Pattern: Pipeline stage (functional filter)
// Use Case: Removing unwanted values, conditional forwarding
//
// BREAKPOINT: Set breakpoint here to trace filter operation
// DEBUG: Watch 'predicate' function
func Filter(input <-chan int, predicate func(int) bool) <-chan int {
	// BREAKPOINT: Set breakpoint here before output channel creation
	output := make(chan int)

	go func() {
		// BREAKPOINT: Set breakpoint in filter goroutine
		defer close(output)
		// DEBUG: Range over input channel
		for v := range input {
			// BREAKPOINT: Set breakpoint here to test predicate
			// DEBUG: Watch 'predicate(v)' return value
			if predicate(v) {
				// BREAKPOINT: Set breakpoint here when value passes filter
				// DEBUG: Forward value only if predicate returns true
				output <- v
			}
			// DEBUG: Values where predicate returns false are dropped
		}
	}()

	return output
}

// Map creates a channel that transforms values using a function.
//
// Pattern: Pipeline stage (functional map)
// Use Case: Transforming values, converting types, applying computations
//
// BREAKPOINT: Set breakpoint here to trace map operation
// DEBUG: Watch 'transform' function
func Map(input <-chan int, transform func(int) int) <-chan int {
	// BREAKPOINT: Set breakpoint here before output channel creation
	output := make(chan int)

	go func() {
		// BREAKPOINT: Set breakpoint in map goroutine
		defer close(output)
		// DEBUG: Range over input channel
		for v := range input {
			// BREAKPOINT: Set breakpoint here before transform
			// DEBUG: Watch 'v' input value
			transformed := transform(v)
			// BREAKPOINT: Set breakpoint here after transform
			// DEBUG: Watch 'transformed' output value
			output <- transformed
		}
	}()

	return output
}

// Take creates a channel that forwards at most n values from input.
//
// Pattern: Limiting stream (take n elements)
// Use Case: Sampling, rate limiting, testing with limited data
//
// BREAKPOINT: Set breakpoint here to trace take operation
// DEBUG: Watch 'n' parameter
func Take(input <-chan int, n int) <-chan int {
	// BREAKPOINT: Set breakpoint here before output channel creation
	output := make(chan int)

	go func() {
		// BREAKPOINT: Set breakpoint in take goroutine
		defer close(output)
		count := 0
		// DEBUG: Range over input but break after n values
		for v := range input {
			// BREAKPOINT: Set breakpoint here to check count
			// DEBUG: Watch 'count' increment
			if count >= n {
				// BREAKPOINT: Set breakpoint here when limit reached
				// DEBUG: Break loop when n values forwarded
				break
			}
			// BREAKPOINT: Set breakpoint here before forwarding
			output <- v
			count++
		}
	}()

	return output
}

// OrDone wraps a channel and adds cancellation via context.
//
// Pattern: Context cancellation wrapper (defensive channel reading)
// Use Case: Adding timeout/cancellation to any channel operation
// Critical: Prevents goroutine leaks when context cancels
//
// BREAKPOINT: Set breakpoint here to trace OrDone wrapper
// DEBUG: Watch 'ctx' context
func OrDone(ctx context.Context, input <-chan int) <-chan int {
	// BREAKPOINT: Set breakpoint here before output channel creation
	output := make(chan int)

	go func() {
		// BREAKPOINT: Set breakpoint in OrDone goroutine
		defer close(output)
		for {
			// BREAKPOINT: Set breakpoint here before select
			select {
			case <-ctx.Done():
				// BREAKPOINT: Set breakpoint here when context cancelled
				// DEBUG: Context cancelled - exit immediately
				return
			case v, ok := <-input:
				// BREAKPOINT: Set breakpoint here when value received
				// DEBUG: Watch 'ok' - false if channel closed
				if !ok {
					// BREAKPOINT: Set breakpoint here when input closed
					// DEBUG: Input channel closed - exit
					return
				}
				// BREAKPOINT: Set breakpoint here before forwarding
				// Nested select: also check context when sending
				select {
				case output <- v:
					// BREAKPOINT: Set breakpoint here when successfully forwarded
					// DEBUG: Value forwarded to output
				case <-ctx.Done():
					// BREAKPOINT: Set breakpoint here if context cancelled during send
					// DEBUG: Context cancelled during send - exit
					return
				}
			}
		}
	}()

	return output
}

// Tee splits an input channel into two output channels (Fan-out pattern).
//
// Pattern: Fan-out (single producer → multiple consumers)
// Use Case: Broadcasting values, duplicating stream, parallel processing
// Challenge: Both outputs must receive each value (parallel sends needed)
//
// BREAKPOINT: Set breakpoint here to trace tee operation
func Tee(input <-chan int) (<-chan int, <-chan int) {
	// BREAKPOINT: Set breakpoint here before output channel creation
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		// BREAKPOINT: Set breakpoint in tee goroutine
		defer close(out1) // Close both outputs when done
		defer close(out2)

		// DEBUG: Range over input channel
		for v := range input {
			// BREAKPOINT: Set breakpoint here before splitting
			// DEBUG: Watch 'v' - value to split to both channels
			// Create local copies for goroutines (closure capture)
			val1, val2 := v, v

			// BREAKPOINT: Set breakpoint here before parallel sends
			// Send to both outputs in parallel to avoid blocking
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				// BREAKPOINT: Set breakpoint in out1 sender
				defer wg.Done()
				out1 <- val1
			}()

			go func() {
				// BREAKPOINT: Set breakpoint in out2 sender
				defer wg.Done()
				out2 <- val2
			}()

			// BREAKPOINT: Set breakpoint here before wait
			// DEBUG: Wait for both sends to complete before next value
			wg.Wait()
		}
	}()

	return out1, out2
}

// Bridge flattens a channel of channels into a single channel.
//
// Pattern: Channel flattening (concatenate multiple streams)
// Use Case: Processing batches, combining streams from dynamic sources
//
// BREAKPOINT: Set breakpoint here to trace bridge operation
func Bridge(input <-chan (<-chan int)) <-chan int {
	// BREAKPOINT: Set breakpoint here before output channel creation
	output := make(chan int)

	go func() {
		// BREAKPOINT: Set breakpoint in bridge goroutine
		defer close(output)

		// DEBUG: Outer loop: receive channels from input
		for ch := range input {
			// BREAKPOINT: Set breakpoint here when new channel received
			// DEBUG: Inner loop: receive values from each channel
			for v := range ch {
				// BREAKPOINT: Set breakpoint here to see values forwarded
				// DEBUG: Forward all values from inner channel to output
				output <- v
			}
			// DEBUG: When inner channel closes, continue to next channel
		}
	}()

	return output
}

// Debounce creates a channel that only forwards values if no new value
// arrives within the specified duration.
//
// Pattern: Debouncing (wait for quiet period)
// Use Case: Search input, resize events, button clicks (prevent spam)
//
// BREAKPOINT: Set breakpoint here to trace debounce operation
// DEBUG: Watch 'duration' parameter
func Debounce(input <-chan int, duration time.Duration) <-chan int {
	// BREAKPOINT: Set breakpoint here before output channel creation
	output := make(chan int)

	go func() {
		// BREAKPOINT: Set breakpoint in debounce goroutine
		defer close(output)

		var timer *time.Timer
		var lastValue int
		var hasValue bool

		// DEBUG: Loop forever until input closes
		for {
			// BREAKPOINT: Set breakpoint here before select
			select {
			case v, ok := <-input:
				// BREAKPOINT: Set breakpoint here when value received
				// DEBUG: Watch 'ok' - false if channel closed
				if !ok {
					// BREAKPOINT: Set breakpoint here when input closed
					// Input closed, send pending value if any
					if hasValue && timer != nil {
						// BREAKPOINT: Set breakpoint here before sending final value
						timer.Stop()
						output <- lastValue
					}
					return
				}

				// BREAKPOINT: Set breakpoint here before timer reset
				// Reset timer (cancel previous if exists)
				if timer != nil {
					timer.Stop()
				}

				// BREAKPOINT: Set breakpoint here to store new value
				// DEBUG: Store value but don't send yet
				lastValue = v
				hasValue = true
				// BREAKPOINT: Set breakpoint here before timer setup
				// DEBUG: Start timer - will send value after duration if no new value arrives
				timer = time.AfterFunc(duration, func() {
					// BREAKPOINT: Set breakpoint here in timer callback
					// DEBUG: Timer fired - send value and reset flag
					output <- lastValue
					hasValue = false
				})
			}
		}
	}()

	return output
}

// ============================================================================
// Exercise 3: Concurrent Data Structures
// ============================================================================

// NewBoundedQueue creates a queue with a maximum capacity.
//
// Pattern: Channel-based queue (buffered channel as queue)
// Use Case: Work queues, rate limiting, backpressure
//
// BREAKPOINT: Set breakpoint here to trace queue creation
// DEBUG: Watch 'capacity' parameter
func NewBoundedQueue(capacity int) *BoundedQueue {
	// BREAKPOINT: Set breakpoint here before returning
	// DEBUG: Buffered channel = queue with capacity limit
	return &BoundedQueue{
		ch: make(chan int, capacity),
	}
}

// Enqueue adds a value to the queue (blocks if full).
//
// BREAKPOINT: Set breakpoint here to trace enqueue
// DEBUG: Watch 'value' parameter
func (q *BoundedQueue) Enqueue(value int) {
	// BREAKPOINT: Set breakpoint here before send
	// DEBUG: Send to buffered channel (blocks if full)
	q.ch <- value
	// DEBUG: Blocks until space available (backpressure)
}

// Dequeue removes and returns a value from the queue (blocks if empty).
//
// BREAKPOINT: Set breakpoint here to trace dequeue
func (q *BoundedQueue) Dequeue() int {
	// BREAKPOINT: Set breakpoint here before receive
	// DEBUG: Receive from channel (blocks if empty)
	return <-q.ch
}

// TryEnqueue attempts to add a value without blocking.
//
// Pattern: Non-blocking operation (select with default)
//
// BREAKPOINT: Set breakpoint here to trace try-enqueue
// DEBUG: Watch 'value' parameter and return value
func (q *BoundedQueue) TryEnqueue(value int) bool {
	// BREAKPOINT: Set breakpoint here before select
	select {
	case q.ch <- value:
		// BREAKPOINT: Set breakpoint here when enqueue succeeds
		// DEBUG: Space available - value enqueued
		return true
	default:
		// BREAKPOINT: Set breakpoint here when queue full
		// DEBUG: Queue full - non-blocking return
		return false
	}
}

// TryDequeue attempts to remove a value without blocking.
//
// Pattern: Non-blocking operation (select with default)
//
// BREAKPOINT: Set breakpoint here to trace try-dequeue
func (q *BoundedQueue) TryDequeue() (int, bool) {
	// BREAKPOINT: Set breakpoint here before select
	select {
	case v := <-q.ch:
		// BREAKPOINT: Set breakpoint here when dequeue succeeds
		// DEBUG: Value available - return it
		return v, true
	default:
		// BREAKPOINT: Set breakpoint here when queue empty
		// DEBUG: Queue empty - return zero value and false
		return 0, false
	}
}

// NewBroadcaster creates a new broadcaster for pub-sub pattern.
//
// Pattern: Publisher-Subscriber (one-to-many message distribution)
// Use Case: Event broadcasting, notification systems, observer pattern
//
// BREAKPOINT: Set breakpoint here to trace broadcaster creation
func NewBroadcaster() *Broadcaster {
	// BREAKPOINT: Set breakpoint here before struct creation
	b := &Broadcaster{
		listeners: make([]chan Message, 0),
		input:     make(chan Message, 100), // Buffered for non-blocking sends
		done:      make(chan struct{}),     // Signal for shutdown
	}

	// BREAKPOINT: Set breakpoint here before goroutine launch
	// Start broadcast goroutine
	go func() {
		// BREAKPOINT: Set breakpoint in broadcast goroutine
		for {
			// BREAKPOINT: Set breakpoint here before select
			select {
			case <-b.done:
				// BREAKPOINT: Set breakpoint here when shutdown requested
				// Close all listener channels
				b.mu.Lock()
				// BREAKPOINT: Set breakpoint here before closing listeners
				// DEBUG: Close all subscriber channels (signals completion)
				for _, ch := range b.listeners {
					close(ch)
				}
				b.mu.Unlock()
				return

			case msg := <-b.input:
				// BREAKPOINT: Set breakpoint here when message received
				// DEBUG: Watch 'msg' - message to broadcast
				// Broadcast to all listeners
				b.mu.RLock() // Read lock (multiple readers OK)
				// BREAKPOINT: Set breakpoint here before broadcasting
				// DEBUG: Send message to all subscribers
				for _, ch := range b.listeners {
					ch <- msg
				}
				b.mu.RUnlock()
			}
		}
	}()

	return b
}

// Subscribe adds a new listener and returns its channel.
//
// BREAKPOINT: Set breakpoint here to trace subscription
func (b *Broadcaster) Subscribe() <-chan Message {
	// BREAKPOINT: Set breakpoint here before channel creation
	// DEBUG: Create buffered channel for subscriber (non-blocking for broadcaster)
	ch := make(chan Message, 10)

	// BREAKPOINT: Set breakpoint here before adding listener
	b.mu.Lock()
	// DEBUG: Add channel to listeners slice
	b.listeners = append(b.listeners, ch)
	b.mu.Unlock()

	return ch
}

// Unsubscribe removes a listener.
//
// BREAKPOINT: Set breakpoint here to trace unsubscription
// DEBUG: Watch 'ch' parameter
func (b *Broadcaster) Unsubscribe(ch <-chan Message) {
	// BREAKPOINT: Set breakpoint here before lock
	b.mu.Lock()
	defer b.mu.Unlock()

	// BREAKPOINT: Set breakpoint here before search loop
	// DEBUG: Find and remove channel from listeners
	for i, listener := range b.listeners {
		if listener == ch {
			// BREAKPOINT: Set breakpoint here when listener found
			// Remove from slice
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			// BREAKPOINT: Set breakpoint here before closing channel
			// DEBUG: Close channel (signals subscriber that broadcaster stopped)
			close(listener)
			break
		}
	}
}

// Send broadcasts a message to all subscribers.
//
// BREAKPOINT: Set breakpoint here to trace send
// DEBUG: Watch 'msg' parameter
func (b *Broadcaster) Send(msg Message) {
	// BREAKPOINT: Set breakpoint here before send
	// DEBUG: Send to input channel (buffered, so non-blocking)
	b.input <- msg
}

// Close stops the broadcaster.
//
// BREAKPOINT: Set breakpoint here to trace close
func (b *Broadcaster) Close() {
	// BREAKPOINT: Set breakpoint here before closing done channel
	// DEBUG: Close done channel (signals broadcast goroutine to stop)
	close(b.done)
}

// NewBarrier creates a barrier for synchronizing n goroutines.
//
// Pattern: Barrier synchronization (all goroutines wait at barrier)
// Use Case: Phased algorithms, parallel sections, synchronization points
//
// BREAKPOINT: Set breakpoint here to trace barrier creation
// DEBUG: Watch 'n' parameter (number of goroutines to wait for)
func NewBarrier(n int) *Barrier {
	// BREAKPOINT: Set breakpoint here before returning
	return &Barrier{
		n:       n,
		count:   0, // Number of goroutines that have arrived
		ch:      make(chan struct{}),
		waiting: make(chan struct{}), // Channel for waiting goroutines
	}
}

// Wait blocks until all n goroutines have called Wait.
//
// Algorithm:
// 1. Increment count (mutex-protected)
// 2. If not last goroutine: wait on waiting channel
// 3. If last goroutine: close waiting channel (unblocks all waiters), reset
//
// BREAKPOINT: Set breakpoint here to trace barrier wait
func (b *Barrier) Wait() {
	// BREAKPOINT: Set breakpoint here before lock
	b.mu.Lock()
	b.count++ // Increment arrival count

	if b.count < b.n {
		// BREAKPOINT: Set breakpoint here when not all goroutines arrived
		// Not all goroutines have arrived yet
		waitingChan := b.waiting // Capture before unlock
		b.mu.Unlock()

		// BREAKPOINT: Set breakpoint here before blocking
		// DEBUG: Block until last goroutine closes waiting channel
		<-waitingChan
		return
	}

	// BREAKPOINT: Set breakpoint here when last goroutine arrives
	// Last goroutine to arrive
	// Signal all waiting goroutines
	// DEBUG: Closing channel unblocks all receivers
	close(b.waiting)

	// BREAKPOINT: Set breakpoint here before reset
	// Reset for next use
	b.count = 0
	b.waiting = make(chan struct{}) // Create new channel for next barrier

	b.mu.Unlock()
}
