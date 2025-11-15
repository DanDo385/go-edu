//go:build solution
// +build solution

package exercise

import (
	"context"
	"sync"
	"time"
)

// Ping creates a channel and sends a single value, then closes it.
func Ping(value int) <-chan int {
	ch := make(chan int, 1)
	ch <- value
	close(ch)
	return ch
}

// PingPong creates two channels that play ping-pong n times.
func PingPong(n int) (chan<- int, <-chan int) {
	ping := make(chan int)
	pong := make(chan int)

	go func() {
		defer close(pong)
		for i := 0; i < n; i++ {
			value := <-ping
			pong <- value
		}
	}()

	return ping, pong
}

// Merge combines multiple input channels into a single output channel.
func Merge(channels ...<-chan int) <-chan int {
	output := make(chan int)
	var wg sync.WaitGroup

	// Launch a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				output <- v
			}
		}(ch)
	}

	// Close output when all inputs are drained
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// Filter creates a channel that only forwards values matching the predicate.
func Filter(input <-chan int, predicate func(int) bool) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for v := range input {
			if predicate(v) {
				output <- v
			}
		}
	}()

	return output
}

// Map creates a channel that transforms values using a function.
func Map(input <-chan int, transform func(int) int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for v := range input {
			output <- transform(v)
		}
	}()

	return output
}

// Take creates a channel that forwards at most n values from input.
func Take(input <-chan int, n int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		count := 0
		for v := range input {
			if count >= n {
				break
			}
			output <- v
			count++
		}
	}()

	return output
}

// OrDone wraps a channel and adds cancellation via context.
func OrDone(ctx context.Context, input <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-input:
				if !ok {
					return
				}
				select {
				case output <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return output
}

// Tee splits an input channel into two output channels.
func Tee(input <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)

		for v := range input {
			// Create local copies for goroutine
			val1, val2 := v, v

			// Send to both outputs (in parallel to avoid blocking)
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				out1 <- val1
			}()

			go func() {
				defer wg.Done()
				out2 <- val2
			}()

			wg.Wait()
		}
	}()

	return out1, out2
}

// Bridge flattens a channel of channels into a single channel.
func Bridge(input <-chan (<-chan int)) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)

		for ch := range input {
			for v := range ch {
				output <- v
			}
		}
	}()

	return output
}

// Debounce creates a channel that only forwards values if no new value
// arrives within the specified duration.
func Debounce(input <-chan int, duration time.Duration) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)

		var timer *time.Timer
		var lastValue int
		var hasValue bool

		for {
			select {
			case v, ok := <-input:
				if !ok {
					// Input closed, send pending value if any
					if hasValue && timer != nil {
						timer.Stop()
						output <- lastValue
					}
					return
				}

				// Reset timer
				if timer != nil {
					timer.Stop()
				}

				lastValue = v
				hasValue = true
				timer = time.AfterFunc(duration, func() {
					output <- lastValue
					hasValue = false
				})

			}
		}
	}()

	return output
}

// NewBoundedQueue creates a queue with a maximum capacity.
func NewBoundedQueue(capacity int) *BoundedQueue {
	return &BoundedQueue{
		ch: make(chan int, capacity),
	}
}

// Enqueue adds a value to the queue (blocks if full).
func (q *BoundedQueue) Enqueue(value int) {
	q.ch <- value
}

// Dequeue removes and returns a value from the queue (blocks if empty).
func (q *BoundedQueue) Dequeue() int {
	return <-q.ch
}

// TryEnqueue attempts to add a value without blocking.
func (q *BoundedQueue) TryEnqueue(value int) bool {
	select {
	case q.ch <- value:
		return true
	default:
		return false
	}
}

// TryDequeue attempts to remove a value without blocking.
func (q *BoundedQueue) TryDequeue() (int, bool) {
	select {
	case v := <-q.ch:
		return v, true
	default:
		return 0, false
	}
}

// NewBroadcaster creates a new broadcaster.
func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		listeners: make([]chan Message, 0),
		input:     make(chan Message, 100),
		done:      make(chan struct{}),
	}

	// Start broadcast goroutine
	go func() {
		for {
			select {
			case <-b.done:
				// Close all listener channels
				b.mu.Lock()
				for _, ch := range b.listeners {
					close(ch)
				}
				b.mu.Unlock()
				return

			case msg := <-b.input:
				// Broadcast to all listeners
				b.mu.RLock()
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
func (b *Broadcaster) Subscribe() <-chan Message {
	ch := make(chan Message, 10)

	b.mu.Lock()
	b.listeners = append(b.listeners, ch)
	b.mu.Unlock()

	return ch
}

// Unsubscribe removes a listener.
func (b *Broadcaster) Unsubscribe(ch <-chan Message) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, listener := range b.listeners {
		if listener == ch {
			// Remove from slice
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			close(listener)
			break
		}
	}
}

// Send broadcasts a message to all subscribers.
func (b *Broadcaster) Send(msg Message) {
	b.input <- msg
}

// Close stops the broadcaster.
func (b *Broadcaster) Close() {
	close(b.done)
}

// NewBarrier creates a barrier for n goroutines.
func NewBarrier(n int) *Barrier {
	return &Barrier{
		n:       n,
		count:   0,
		ch:      make(chan struct{}),
		waiting: make(chan struct{}),
	}
}

// Wait blocks until all n goroutines have called Wait.
func (b *Barrier) Wait() {
	b.mu.Lock()
	b.count++

	if b.count < b.n {
		// Not all goroutines have arrived yet
		b.mu.Unlock()

		// Wait for signal
		<-b.waiting
		return
	}

	// Last goroutine to arrive
	// Signal all waiting goroutines
	close(b.waiting)

	// Reset for next use
	b.count = 0
	b.waiting = make(chan struct{})

	b.mu.Unlock()
}
