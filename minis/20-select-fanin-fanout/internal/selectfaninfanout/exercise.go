//go:build !solution
// +build !solution

package selectfaninfanout

import (
	"time"
)

// Exercise 1: SelectFirst
func SelectFirst(ch1, ch2 <-chan string, timeout time.Duration) (string, bool) {
	// TODO: Implement this function.

	// The `select` statement lets a goroutine wait on multiple communication operations.
	// It blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.

	// Step 1: Use a `select` statement with three cases.
	// - `case v := <-ch1:`
	//   - This case will be chosen if a value is available on `ch1`.
	//   - `v` will be the value received. Return `v, true`.
	// - `case v := <-ch2:`
	//   - This case will be chosen if a value is available on `ch2`.
	//   - Return `v, true`.
	// - `case <-time.After(timeout):`
	//   - `time.After` creates a channel and sends the current time on it after the specified duration.
	//   - This case will be chosen if the timeout is reached before a value is received on `ch1` or `ch2`.
	//   - Return `"", false`.

	// Memory & Language Comparison:
	// - Go: `select` is a first-class language feature, making this kind of operation clean and idiomatic.
	// - JavaScript: You would use `Promise.race([promise1, promise2, timeoutPromise])`. It's similar but can be more complex to set up correctly.
	// - Rust: The `tokio::select!` macro provides very similar functionality and is also a core part of its async ecosystem.
	// - Python: You would use `asyncio.wait()` with `return_when=asyncio.FIRST_COMPLETED`, which is more verbose than Go's `select`.
	return "", false
}

// Exercise 2: NonBlockingSend
func NonBlockingSend(ch chan<- int, value int) bool {
	// TODO: Implement this function.

	// A `select` statement with a `default` case is non-blocking.
	// If one of the channel operations can proceed, it will be selected.
	// If none can proceed, the `default` case is executed immediately.

	// Step 1: Use a `select` statement with two cases.
	// - `case ch <- value:`
	//   - This case will be chosen if the send operation can complete without blocking. This can happen if the channel is buffered and has space, or if a receiver is waiting.
	//   - If the send is successful, return `true`.
	// - `default:`
	//   - This case will be chosen if the send operation would block.
	//   - Return `false`.

	// Memory & Language Comparison:
	// - Go: The `select` with `default` is the idiomatic way to perform non-blocking channel operations.
	// - Java: A `BlockingQueue` has a method `offer(value)`, which is the direct equivalent. It returns `true` if the value was added and `false` if the queue was full.
	// - C#: `Channel<T>.Writer.TryWrite(value)` provides the same functionality.
	// - Rust: `channel.try_send(value)` is the equivalent non-blocking send operation.
	return false
}

// Exercise 3: FanIn
func FanIn(channels ...<-chan int) <-chan int {
	// TODO: Implement this function.

	// This is the "fan-in" pattern: merging multiple input streams into a single output stream.

	// Step 1: Create the output channel.
	// - `out := make(chan int)`

	// Step 2: Create a `sync.WaitGroup`.
	// - This will be used to wait for all the input-reading goroutines to finish before we close the output channel.

	// Step 3: Create a helper function for forwarding.
	// - This function will take an input channel, read from it until it's closed, and send all values to the `out` channel.
	// - `func forward(c <-chan int) { ... }`
	// - Inside this function, use a `for v := range c` loop.
	// - Don't forget to call `wg.Done()` when the loop finishes.

	// Step 4: Start a goroutine for each input channel.
	// - Loop through the `channels` slice.
	// - For each channel, call `wg.Add(1)` and start a goroutine that calls your `forward` function.

	// Step 5: Start a final goroutine to close the output channel.
	// - This goroutine's only job is to wait for all the forwarding goroutines to finish and then close the `out` channel.
	// - `go func() { wg.Wait(); close(out) }()`

	// Step 6: Return the output channel.
	return nil
}

// Exercise 4: FanOut
func FanOut(input <-chan int, numWorkers int, process func(int) int) <-chan int {
	// TODO: Implement this function.

	// This is the "fan-out" or "worker pool" pattern. A single stream of work is distributed among multiple workers to be processed in parallel.

	// Step 1: Create the output channel.
	// - `out := make(chan int)`

	// Step 2: Create a `sync.WaitGroup`.
	// - This will be used to wait for all the worker goroutines to finish.

	// Step 3: Start the worker goroutines.
	// - Loop `numWorkers` times.
	// - In each iteration, call `wg.Add(1)` and launch a new goroutine.

	// Step 4: Implement the logic for each worker goroutine.
	// - The goroutine should use a `for v := range input` loop. The Go runtime automatically handles distributing the values from the `input` channel among the multiple worker goroutines that are all reading from it.
	// - Inside the loop, call the `process` function with the received value.
	// - Send the result to the `out` channel.
	// - When the `input` channel is closed, the loop will terminate. The goroutine should then call `wg.Done()`.

	// Step 5: Start a final goroutine to close the output channel.
	// - This is the same pattern as in `FanIn`.
	// - `go func() { wg.Wait(); close(out) }()`

	// Step 6: Return the output channel.
	return nil
}

// Exercise 5: OrChannel
func OrChannel(channels ...<-chan struct{}) <-chan struct{} {
	// TODO: Implement this function.

	// This pattern is useful for waiting for one of several events to occur.

	// Step 1: Handle the edge cases.
	// - If `len(channels)` is 0, what should you do? A channel that is already closed is a reasonable choice.
	// - If `len(channels)` is 1, you can just return the single input channel.

	// Step 2: Create the output channel.
	// - `out := make(chan struct{})`

	// Step 3: Create a `sync.Once`.
	// - `var once sync.Once`
	// - This is crucial because multiple input channels could close at nearly the same time. Each would trigger its goroutine to try to close `out`. Closing an already-closed channel causes a panic. `sync.Once` ensures that the `close(out)` operation is only ever performed one time.

	// Step 4: Start a goroutine for each input channel.
	// - For each `ch` in `channels`, launch a goroutine.
	// - The goroutine should simply wait to receive from its channel: `<-ch`.
	// - As soon as it unblocks (because the channel was closed), it should call `once.Do(func() { close(out) })`.

	// Step 5: Return the output channel.
	return nil
}

// Exercise 6: TryReceiveAll
func TryReceiveAll(channels []<-chan int) map[int]int {
	// TODO: Implement this function.

	// This function "polls" a set of channels to see if any have data ready right now.

	// Step 1: Create the result map.
	// - `result := make(map[int]int)`

	// Step 2: Loop through the slice of channels.
	// - Use a `for i, ch := range channels` loop.

	// Step 3: For each channel, use a non-blocking receive.
	// - Use a `select` statement with a `default` case.
	// - `case v := <-ch:`
	//   - If a value was immediately available, add it to the map with the channel's index as the key: `result[i] = v`.
	// - `default:`
	//   - If the channel would have blocked, do nothing.

	// Step 4: Return the result map.
	return nil
}

// Exercise 7: RateLimiter
type RateLimiter struct {
	tokens chan struct{}
	rate   int
}

func NewRateLimiter(rate int) *RateLimiter {
	// TODO: Implement this function.

	// This implements the "token bucket" algorithm. The buffered channel is the "bucket".

	// Step 1: Initialize the RateLimiter struct.
	// - The `tokens` channel should be a buffered channel with a capacity equal to `rate`.
	// - `tokens := make(chan struct{}, rate)`
	// - Using `struct{}` as the channel type is a common Go idiom for when you only care about the signal, not the value. An empty struct uses zero memory.

	// Step 2: Fill the bucket with initial tokens.
	// - Loop `rate` times and send an empty struct to the `tokens` channel.

	// Step 3: Start the refill goroutine.
	// - `go rl.refill()`
	// - This goroutine will add a new token to the bucket at regular intervals.

	// Step 4: Implement the `refill` method (or a function called by the goroutine).
	// - Create a `time.Ticker` with a duration of `time.Second / time.Duration(rate)`.
	// - Use a `for range ticker.C` loop to run at the desired rate.
	// - In each iteration, try to send a token to the `tokens` channel. Use a non-blocking send (`select` with `default`) so that you don't add more tokens than the capacity allows.

	return &RateLimiter{}
}

func (rl *RateLimiter) Wait() {
	// TODO: Implement this function.
	// - To wait for a token, simply receive from the `tokens` channel.
	// - `<-rl.tokens`
	// - This will block until a token is available in the bucket.
}

func (rl *RateLimiter) TryWait() bool {
	// TODO: Implement this function.
	// - To try to get a token without blocking, use a `select` with a `default` case.
	// - `case <-rl.tokens:` -> return `true`.
	// - `default:` -> return `false`.
	return false
}

// Exercise 8: Pipeline
func Pipeline(n int, numWorkers int) <-chan int {
	// TODO: Implement this function by composing the patterns you've already built.

	// A pipeline is a series of "stages", where each stage is a function that receives data from an input channel and sends results to an output channel.

	// Step 1: The Generator stage.
	// - Create a function `generate(n int) <-chan int` that produces numbers from 1 to `n` and sends them to a channel, then closes it.
	// - Call this function to get the first channel in your pipeline.

	// Step 2: The Processor stage.
	// - Use the `FanOut` function you wrote earlier.
	// - The `input` will be the channel from the `generate` stage.
	// - The `process` function should square the numbers (`func(x int) int { return x * x }`).

	// Step 3: The Filter stage.
	// - Create a function `filter(in <-chan int, predicate func(int) bool) <-chan int` (similar to the `Filter` exercise in the previous lesson).
	// - The `input` will be the channel from the `FanOut` stage.
	// - The `predicate` function should check for even numbers (`func(x int) bool { return x%2 == 0 }`).

	// Step 4: Return the final channel from the `filter` stage.
	return nil
}

// Helper for Pipeline: generates numbers 1 to n
func generate(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

// Helper for Pipeline: filters values based on a predicate
func filter(in <-chan int, predicate func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			if predicate(v) {
				out <- v
			}
		}
	}()
	return out
}

// Exercise 9: SelectWithPriority
func SelectWithPriority(high, low <-chan int) (int, bool) {
	// TODO: Implement this function.

	// This is a common pattern when you need to handle some events (like cancellation or high-priority messages) before others.

	// Step 1: First, try a non-blocking read on the high-priority channel.
	// - Use a `select` with a `default` case.
	// - `case v := <-high:` -> If a value is immediately available, you're done. Return `v, true`.
	// - `default:` -> If not, do nothing and fall through to the next step.

	// Step 2: If the high-priority channel wasn't ready, perform a blocking select on both channels.
	// - Use a `select` without a `default` case.
	// - `case v := <-high:` -> It's possible a value arrived on `high` while you were executing Step 1. Handle it and return `v, true`.
	// - `case v := <-low:` -> If `high` is not ready but `low` is, handle it and return `v, false`.
	// - This `select` will block until one of the two channels has a value.
	return 0, false
}

// Exercise 10: Timeout
func Timeout(ch <-chan int, timeout time.Duration) (int, bool) {
	// TODO: Implement this function.

	// This is the same pattern as `SelectFirst`, but with only one data channel.

	// Step 1: Use a `select` statement.
	// - `case v := <-ch:`
	//   - If a value is received before the timeout, return `v, true`.
	// - `case <-time.After(timeout):`
	//   - If the timeout is reached first, return `0, false`.

	// Memory & Language Comparison:
	// - The `time.After` function is convenient, but it's important to know how it works. It allocates a `Timer` object that is not garbage collected until after the timer fires.
	// - If you use `time.After` inside a loop where the timeout is rarely hit, you can create a memory leak by allocating many timers that haven't fired yet.
	// - In a loop, it's better to use `time.NewTimer` and reuse it with `timer.Reset()`.
	return 0, false
}
