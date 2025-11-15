//go:build !solution
// +build !solution

package exercise

import (
	"time"
)

// Exercise 1: SelectFirst
//
// Implement a function that receives from two channels and returns
// whichever value arrives first. If neither arrives within the timeout,
// return the empty string and false.
//
// Parameters:
//   - ch1, ch2: Channels to receive from
//   - timeout: Maximum time to wait
//
// Returns:
//   - string: The value that arrived first
//   - bool: true if a value was received, false if timeout
//
// Example:
//   ch1 := make(chan string, 1)
//   ch2 := make(chan string, 1)
//   ch1 <- "fast"
//   value, ok := SelectFirst(ch1, ch2, 1*time.Second)
//   // value == "fast", ok == true
func SelectFirst(ch1, ch2 <-chan string, timeout time.Duration) (string, bool) {
	// TODO: implement
	return "", false
}

// Exercise 2: NonBlockingSend
//
// Attempt to send a value to a channel without blocking.
// If the channel is not ready to receive, return false immediately.
//
// Parameters:
//   - ch: Channel to send to
//   - value: Value to send
//
// Returns:
//   - bool: true if sent successfully, false if channel not ready
//
// Example:
//   ch := make(chan int)
//   sent := NonBlockingSend(ch, 42)  // false (no receiver)
//
//   ch := make(chan int, 1)
//   sent := NonBlockingSend(ch, 42)  // true (buffered channel)
func NonBlockingSend(ch chan<- int, value int) bool {
	// TODO: implement
	return false
}

// Exercise 3: FanIn
//
// Merge multiple input channels into a single output channel.
// The output channel should receive all values from all input channels.
// Close the output channel when all input channels are closed.
//
// Parameters:
//   - channels: Variable number of input channels
//
// Returns:
//   - <-chan int: Output channel with merged values
//
// Example:
//   ch1 := make(chan int)
//   ch2 := make(chan int)
//   merged := FanIn(ch1, ch2)
//
//   go func() {
//       ch1 <- 1
//       ch1 <- 2
//       close(ch1)
//   }()
//
//   go func() {
//       ch2 <- 3
//       ch2 <- 4
//       close(ch2)
//   }()
//
//   for v := range merged {
//       fmt.Println(v)  // Prints 1, 2, 3, 4 (order may vary)
//   }
func FanIn(channels ...<-chan int) <-chan int {
	// TODO: implement
	return nil
}

// Exercise 4: FanOut
//
// Distribute values from input channel to multiple worker goroutines.
// Each worker processes values using the provided function.
// Return a channel of results.
//
// Parameters:
//   - input: Channel of input values
//   - numWorkers: Number of worker goroutines
//   - process: Function to process each value
//
// Returns:
//   - <-chan int: Channel of processed results
//
// Example:
//   input := make(chan int)
//   results := FanOut(input, 3, func(n int) int {
//       return n * n
//   })
//
//   go func() {
//       for i := 1; i <= 10; i++ {
//           input <- i
//       }
//       close(input)
//   }()
//
//   for result := range results {
//       fmt.Println(result)  // Prints squares of 1-10
//   }
func FanOut(input <-chan int, numWorkers int, process func(int) int) <-chan int {
	// TODO: implement
	return nil
}

// Exercise 5: OrChannel
//
// Create a channel that closes when any of the input channels close.
// This is useful for combining cancellation signals.
//
// Parameters:
//   - channels: Variable number of input channels
//
// Returns:
//   - <-chan struct{}: Channel that closes when any input closes
//
// Example:
//   ch1 := make(chan struct{})
//   ch2 := make(chan struct{})
//   ch3 := make(chan struct{})
//   done := OrChannel(ch1, ch2, ch3)
//
//   go func() {
//       time.Sleep(1 * time.Second)
//       close(ch2)  // This will close 'done'
//   }()
//
//   <-done  // Unblocks after 1 second
func OrChannel(channels ...<-chan struct{}) <-chan struct{} {
	// TODO: implement
	return nil
}

// Exercise 6: TryReceiveAll
//
// Try to receive values from multiple channels without blocking.
// Return a map of channel indices to values received.
// Only include channels that had values available.
//
// Parameters:
//   - channels: Slice of channels to receive from
//
// Returns:
//   - map[int]int: Map of channel index to value received
//
// Example:
//   ch1 := make(chan int, 1)
//   ch2 := make(chan int, 1)
//   ch3 := make(chan int, 1)
//
//   ch1 <- 10
//   ch3 <- 30
//
//   result := TryReceiveAll([]<-chan int{ch1, ch2, ch3})
//   // result == map[int]int{0: 10, 2: 30}
func TryReceiveAll(channels []<-chan int) map[int]int {
	// TODO: implement
	return nil
}

// Exercise 7: RateLimiter
//
// Implement a rate limiter that allows up to 'rate' operations per second.
// Use channels and select to implement the token bucket algorithm.
//
// Methods:
//   - Wait(): Block until a token is available
//   - TryWait(): Try to get a token without blocking (return false if unavailable)
//
// Example:
//   limiter := NewRateLimiter(10)  // 10 ops/sec
//
//   for i := 0; i < 100; i++ {
//       limiter.Wait()  // Blocks if rate exceeded
//       doWork()
//   }
type RateLimiter struct {
	// TODO: add fields
}

func NewRateLimiter(rate int) *RateLimiter {
	// TODO: implement
	return &RateLimiter{}
}

func (rl *RateLimiter) Wait() {
	// TODO: implement
}

func (rl *RateLimiter) TryWait() bool {
	// TODO: implement
	return false
}

// Exercise 8: Pipeline
//
// Build a 3-stage pipeline:
//   1. Generator: Produces numbers 1 to n
//   2. Processor: Squares each number (fan-out to numWorkers)
//   3. Filter: Only keeps even numbers
//
// Parameters:
//   - n: Number of values to generate
//   - numWorkers: Number of workers for processing stage
//
// Returns:
//   - <-chan int: Channel of filtered results
//
// Example:
//   results := Pipeline(10, 3)
//   for v := range results {
//       fmt.Println(v)  // Prints: 4, 16, 36, 64, 100
//   }
func Pipeline(n int, numWorkers int) <-chan int {
	// TODO: implement
	return nil
}

// Exercise 9: SelectWithPriority
//
// Implement select with priority: prefer values from 'high' channel
// over 'low' channel. If neither has a value, block.
//
// Parameters:
//   - high: High-priority channel
//   - low: Low-priority channel
//
// Returns:
//   - int: Value received
//   - bool: true if from high priority, false if from low
//
// Example:
//   high := make(chan int, 1)
//   low := make(chan int, 1)
//   high <- 1
//   low <- 2
//   value, isHigh := SelectWithPriority(high, low)
//   // value == 1, isHigh == true
func SelectWithPriority(high, low <-chan int) (int, bool) {
	// TODO: implement
	return 0, false
}

// Exercise 10: Timeout
//
// Wait for a value from the channel with a timeout.
// Return the value and true if received within timeout.
// Return zero value and false if timeout occurs.
//
// Parameters:
//   - ch: Channel to receive from
//   - timeout: Maximum time to wait
//
// Returns:
//   - int: Value received (or 0 if timeout)
//   - bool: true if received, false if timeout
//
// Example:
//   ch := make(chan int)
//   go func() {
//       time.Sleep(2 * time.Second)
//       ch <- 42
//   }()
//   value, ok := Timeout(ch, 1*time.Second)
//   // value == 0, ok == false (timeout)
func Timeout(ch <-chan int, timeout time.Duration) (int, bool) {
	// TODO: implement
	return 0, false
}
