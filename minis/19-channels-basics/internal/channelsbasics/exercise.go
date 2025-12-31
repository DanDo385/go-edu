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
	// TODO: Implement this.

	// Step 1: Create a buffered channel of integers with a capacity of 1.
	// - `ch := make(chan int, 1)`
	// - Why buffered? A buffered channel allows the sender to send a value without a receiver immediately being ready. In this case, it lets us send the value and close the channel in the same function without blocking forever.

	// Step 2: Send the value to the channel.
	// - `ch <- value`

	// Step 3: Close the channel.
	// - `close(ch)`
	// - Closing a channel is important. It signals to receivers that no more values will be sent. A receiver can test for this with the "comma ok" idiom: `val, ok := <-ch`. If `ok` is `false`, the channel is closed and empty.

	// Step 4: Return the channel.
	// - The function signature returns a `<-chan int`, which is a "receive-only" channel. This is a compile-time safety feature that prevents the caller from accidentally sending to a channel they are only supposed to be receiving from.

	// Memory & Language Comparison:
	// - Go: Channels are a core concurrency primitive, providing a safe way for goroutines to communicate. They are fundamental to Go's "communicating sequential processes" (CSP) philosophy.
	// - C#: `System.Threading.Channels` provides a very similar concept.
	// - Java: You would typically use a `BlockingQueue` from the `java.util.concurrent` package to achieve similar results.
	// - Rust: The `std::sync::mpsc` (multiple producer, single consumer) and libraries like `crossbeam-channel` provide channel functionality.
	// - Python: `asyncio.Queue` provides queues for asynchronous programming.
	return nil
}

// PingPong creates two channels that play ping-pong n times.
func PingPong(n int) (chan<- int, <-chan int) {
	// TODO: Implement this.

	// This exercise demonstrates two-way communication between goroutines.

	// Step 1: Create two unbuffered channels.
	// - `ping := make(chan int)`
	// - `pong := make(chan int)`
	// - We use unbuffered channels here, which means that a sender will block until a receiver is ready. This is how the "game" is synchronized.

	// Step 2: Launch a new goroutine for the "pong" player.
	// - `go func() { ... }()`
	// - This goroutine will run concurrently with the calling function.

	// Step 3: Implement the pong player's logic inside the goroutine.
	// - Use a `for` loop that runs `n` times.
	// - In each iteration, wait to receive a value from the `ping` channel: `value := <-ping`.
	// - Then, send the same value back on the `pong` channel: `pong <- value`.
	// - After the loop, close the `pong` channel to signal that the game is over: `defer close(pong)`.

	// Step 4: Return the two channels from the main function.
	// - The `ping` channel is returned as a `chan<- int` (send-only).
	// - The `pong` channel is returned as a `<-chan int` (receive-only).
	// - This provides type safety, preventing the caller from using the channels in the wrong direction.

	// Memory & Language Comparison:
	// - This is a classic CSP pattern. The two goroutines are completely decoupled and communicate only through channels.
	// - In a thread-based model (like C++ or Java without high-level concurrency libraries), you might use a shared state variable with mutexes and condition variables to achieve this synchronization, which is much more complex and error-prone.
	return nil, nil
}

// Merge combines multiple input channels into a single output channel.
func Merge(channels ...<-chan int) <-chan int {
	// TODO: Implement this.

	// This is the "fan-in" pattern. You have multiple sources of data (the input channels) and you want to combine them into a single stream.

	// Step 1: Create the output channel.
	// - `output := make(chan int)`

	// Step 2: Create a `sync.WaitGroup`.
	// - A `WaitGroup` is used to wait for a collection of goroutines to finish.

	// Step 3: Loop through the input channels.
	// - For each input channel, launch a new goroutine.
	// - `wg.Add(1)` before starting each goroutine.

	// Step 4: Implement the logic for each input goroutine.
	// - The goroutine should read all values from its assigned input channel using a `for v := range c { ... }` loop.
	// - It should send each value to the `output` channel.
	// - When the input channel is closed, the `range` loop will terminate. The goroutine should then call `wg.Done()`.
	// - **Important:** You must pass the channel as an argument to the goroutine to avoid a common loop variable closure bug: `go func(c <-chan int) { ... }(ch)`.

	// Step 5: Create a final goroutine to close the output channel.
	// - This goroutine's only job is to wait for all the other goroutines to finish and then close the `output` channel.
	// - `go func() { wg.Wait(); close(output) }()`
	// - `wg.Wait()` will block until `wg.Done()` has been called for every `wg.Add(1)`.
	// - This ensures that `output` is only closed after all data has been forwarded.

	// Step 6: Return the output channel.
	return nil
}

// Filter creates a channel that only forwards values matching the predicate.
func Filter(input <-chan int, predicate func(int) bool) <-chan int {
	// TODO: Implement this.

	// This function is a "pipeline stage". It receives data, processes it (by filtering), and passes it on.

	// Step 1: Create the output channel.
	// - `output := make(chan int)`

	// Step 2: Launch a goroutine to do the work.
	// - `go func() { ... }()`

	// Step 3: Implement the goroutine's logic.
	// - Use a `for v := range input` loop to read from the input channel until it's closed.
	// - Inside the loop, check if the value satisfies the predicate: `if predicate(v)`.
	// - If it does, send the value to the `output` channel.
	// - After the loop finishes (which happens when `input` is closed), close the `output` channel to signal to downstream consumers that you're done. `defer close(output)` is a good way to ensure this happens.

	// Step 4: Return the output channel immediately.
	// - The caller gets the output channel right away and can start listening to it. The goroutine works in the background.
	return nil
}

// Map creates a channel that transforms values using a function.
func Map(input <-chan int, transform func(int) int) <-chan int {
	// TODO: Implement this.

	// This is another "pipeline stage", similar to `Filter`.

	// Step 1: Create the output channel.
	// - `output := make(chan int)`

	// Step 2: Launch a goroutine.
	// - `go func() { ... }()`

	// Step 3: Implement the goroutine's logic.
	// - Use a `for v := range input` loop to read from the input channel.
	// - For each value, apply the transformation: `transformedValue := transform(v)`.
	// - Send the `transformedValue` to the `output` channel.
	// - When the `input` channel is closed, the loop will end. Close the `output` channel. `defer close(output)` is a good pattern here.

	// Step 4: Return the output channel.

	// Memory & Language Comparison:
	// - This pattern is common in many languages.
	// - In functional languages (and languages with functional features), this is known as a "map" operation on a sequence.
	// - Go's implementation is unique because the sequence is a channel, and the operation happens concurrently in a separate goroutine.
	// - In Rust, you would use the `.map()` adapter on an `Iterator`.
	// - In Python, you could use a generator expression: `(transform(v) for v in input)`.
	return nil
}

// Take creates a channel that forwards at most n values from input.
func Take(input <-chan int, n int) <-chan int {
	// TODO: Implement this.

	// This is another pipeline stage that modifies the stream of data, in this case by truncating it.

	// Step 1: Create the output channel.
	// - `output := make(chan int)`

	// Step 2: Launch a goroutine.
	// - `go func() { ... }()`

	// Step 3: Implement the goroutine's logic.
	// - Use a counter variable, initialized to 0.
	// - Use a `for v := range input` loop.
	// - Inside the loop, check if the counter has reached `n`. If it has, `break` the loop.
	// - If the counter is less than `n`, send the value to the `output` channel and increment the counter.
	// - When the loop finishes (either because the `input` channel was closed or because you broke out of it), close the `output` channel. `defer close(output)`.

	// Step 4: Return the output channel.

	// Memory & Language Comparison:
	// - Most languages with support for lazy iteration have a "take" or "limit" operation.
	// - In Rust, you would use `.take(n)` on an `Iterator`.
	// - In Python, you can use `itertools.islice(input, n)`.
	// - In C# (with LINQ), you would use `.Take(n)`.
	// - Go's channel-based approach is different because it's inherently concurrent.
	return nil
}

// OrDone wraps a channel and adds cancellation via context.
func OrDone(ctx context.Context, input <-chan int) <-chan int {
	// TODO: Implement this.

	// This is a very important pattern in Go for making concurrent operations cancellable.

	// Step 1: Create the output channel.
	// - `output := make(chan int)`

	// Step 2: Launch a goroutine.
	// - `go func() { ... }()`

	// Step 3: Implement the goroutine's logic using a `select` statement.
	// - Use an infinite `for { ... }` loop.
	// - Inside the loop, `select` on two cases:
	//   - `case <-ctx.Done():`
	//     - The context has been cancelled. `return` from the goroutine to stop processing.
	//   - `case v, ok := <-input:`
	//     - A value has been received from the input channel.
	//     - Check the `ok` flag. If it's `false`, the input channel has been closed, so `return` from the goroutine.
	//     - If `ok` is `true`, you need to send the value to the `output` channel.
	//     - **Critical:** Sending to `output` might block. If the context is cancelled while you are blocked, you need to be able to stop. So, you must use another `select` statement:
	//       ```
	//       select {
	//       case output <- v:
	//       case <-ctx.Done():
	//           return
	//       }
	//       ```
	// - Don't forget to `defer close(output)` in the goroutine.

	// Step 4: Return the output channel.
	return nil
}

// Tee splits an input channel into two output channels.
func Tee(input <-chan int) (<-chan int, <-chan int) {
	// TODO: Implement this.

	// This is the "fan-out" pattern. You are splitting one stream of data into two identical streams.

	// Step 1: Create the two output channels.
	// - `out1 := make(chan int)`
	// - `out2 := make(chan int)`

	// Step 2: Launch a goroutine.
	// - `go func() { ... }()`

	// Step 3: Implement the goroutine's logic.
	// - Use a `for v := range input` loop to read from the input.
	// - When the loop ends, close both output channels: `defer close(out1)` and `defer close(out2)`.
	// - Inside the loop, you have a value `v` that you need to send to *both* `out1` and `out2`.
	// - **The Challenge:** What if a receiver is ready on `out1`, but not on `out2`? If you write `out1 <- v; out2 <- v`, your goroutine could block forever on the second send.
	// - **The Solution:** Send to both channels concurrently.
	//   - Use a `sync.WaitGroup` for each value.
	//   - Launch two new goroutines, one for sending to `out1` and one for sending to `out2`.
	//   - Wait for both sends to complete before reading the next value from `input`.

	// Step 4: Return the two output channels.
	return nil, nil
}

// Bridge flattens a channel of channels into a single channel.
func Bridge(input <-chan (<-chan int)) <-chan int {
	// TODO: Implement this.

	// This pattern is useful when you have a series of concurrent operations that each produce their own output channel, and you want to combine all their results into one stream.

	// Step 1: Create the output channel.
	// - `output := make(chan int)`

	// Step 2: Launch a goroutine.
	// - `go func() { ... }()`

	// Step 3: Implement the goroutine's logic.
	// - Use a `for...range` loop to read from the `input` channel. Each item `ch` will be a channel (`<-chan int`).
	// - For each `ch` received, use another `for...range` loop to read all the values from it: `for v := range ch`.
	// - Send each value `v` to the `output` channel.
	// - When the outer loop finishes (because `input` was closed), close the `output` channel. `defer close(output)`.

	// Step 4: Return the output channel.
	return nil
}

// Debounce creates a channel that only forwards values if no new value
// arrives within the specified duration.
func Debounce(input <-chan int, duration time.Duration) <-chan int {
	// TODO: Implement this.

	// Debouncing is a common pattern for handling rapid events (like keystrokes or sensor data) where you only care about the final event in a burst.

	// Step 1: Create the output channel.
	// - `output := make(chan int)`

	// Step 2: Launch a goroutine.
	// - `go func() { ... }()`

	// Step 3: Implement the goroutine's logic.
	// - You'll need a `*time.Timer` and a variable to hold the last received value.
	// - Use an infinite `for` loop and a `select` statement.
	// - `case v, ok := <-input:`
	//   - If the input channel is closed (`!ok`), send any pending value and return.
	//   - If you receive a new value, stop the previous timer if it exists (`timer.Stop()`).
	//   - Store the new value.
	//   - Create a new timer with `time.AfterFunc(duration, func() { ... })`. The function passed to `AfterFunc` will be executed after `duration` has passed. This function should send the stored value to the `output` channel.
	// - Don't forget to `defer close(output)`.

	// Step 4: Return the output channel.
	return nil
}

// NewBoundedQueue creates a queue with a maximum capacity.
func NewBoundedQueue(capacity int) *BoundedQueue {
	// TODO: Implement this.

	// A buffered channel is a natural fit for a bounded queue.
	// - The channel's buffer acts as the storage for the queue.
	// - The channel's capacity enforces the size limit.
	// - The channel's blocking behavior on send/receive naturally handles the "full" and "empty" conditions.

	// Step 1: Create a `BoundedQueue` struct.
	// Step 2: Initialize its `ch` field with a buffered channel of the given `capacity`.
	// Step 3: Return a pointer to the struct.
	return nil
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
	// TODO: Implement this.

	// The broadcaster pattern involves a central goroutine that manages subscribers and message distribution.

	// Step 1: Initialize the Broadcaster struct.
	// - `listeners`: a new, empty slice of channels.
	// - `input`: a channel for receiving messages to be broadcast.
	// - `done`: a channel to signal shutdown.

	// Step 2: Launch the central management goroutine.
	// - `go func() { ... }()`
	// - This goroutine will loop forever, using a `select` statement to handle incoming messages and shutdown signals.

	// Step 3: Implement the goroutine's logic.
	// - `case msg := <-b.input:`
	//   - When a message arrives, iterate through the `listeners` slice and send the message to each listener channel.
	//   - **Important:** You must use a read lock (`b.mu.RLock()`) around the iteration to prevent the `listeners` slice from being modified while you're using it.
	// - `case <-b.done:`
	//   - When the shutdown signal is received, lock the `listeners` slice (`b.mu.Lock()`), close all listener channels, and return from the goroutine.

	// Step 4: Return the initialized broadcaster.
	return nil
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
	// TODO: Implement this.

	// Step 1: Initialize the Barrier struct.
	// - `n`: the number of goroutines to wait for.
	// - `count`: the number of goroutines currently waiting, initialized to 0.
	// - `waiting`: a channel that waiting goroutines will block on.
	// - A `sync.Mutex` is also needed to protect the `count` field.
	return nil
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
