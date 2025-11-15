//go:build solution
// +build solution

/*
Project 20: Select, Fan-In, and Fan-Out - Solutions

This file contains complete solutions to all exercises with detailed explanations.

Key Go Concepts Demonstrated:
1. Select statement for channel multiplexing
2. Non-blocking operations with default case
3. Timeout patterns with time.After
4. Fan-in pattern (merging channels)
5. Fan-out pattern (worker pools)
6. Pipeline architecture
7. Priority selection

Why Go is well-suited for this:
- Select statement is a first-class language feature
- Channels are lightweight and built into the runtime
- Goroutines make concurrent patterns trivial
- No callback hell or promise chains
- Clean, readable concurrency patterns

Compared to other languages:
- Python: No select equivalent (need asyncio.wait with complex setup)
- JavaScript: Promise.race() similar but less flexible
- Rust: select! macro in tokio, similar power but more verbose
- Java: No direct equivalent (need CompletableFuture.anyOf with boilerplate)
*/

package exercise

import (
	"sync"
	"time"
)

// ============================================================================
// Exercise 1: SelectFirst
// ============================================================================

/*
Problem: Receive from whichever channel has data first

Architecture:
- Use select to wait on multiple channels simultaneously
- Add timeout case to prevent indefinite blocking
- Return value and success indicator

Complexity:
- Time: O(1) - select operation is constant time
- Space: O(1) - no additional data structures

Three-Input Iteration Table:

Input 1: ch1 has value immediately
  ch1: "fast"
  ch2: (empty)
  select: ch1 case executes
  return: "fast", true

Input 2: Timeout occurs
  ch1: (empty)
  ch2: (empty)
  time passes...
  select: timeout case executes
  return: "", false

Input 3: ch2 sends after delay
  ch1: (empty)
  ch2: sends "slow" after 500ms
  timeout: 1 second
  select: ch2 case executes at 500ms
  return: "slow", true
*/

func SelectFirst(ch1, ch2 <-chan string, timeout time.Duration) (string, bool) {
	// Use select to wait on multiple channel operations
	// Select will execute whichever case becomes ready first
	select {
	case v := <-ch1:
		// ch1 sent a value first
		return v, true

	case v := <-ch2:
		// ch2 sent a value first
		return v, true

	case <-time.After(timeout):
		// Neither channel sent within timeout
		// time.After returns a channel that receives after duration
		return "", false
	}
}

// ============================================================================
// Exercise 2: NonBlockingSend
// ============================================================================

/*
Problem: Try to send without blocking

Architecture:
- Use select with default case
- Default case executes immediately if send would block
- Return success indicator

Complexity:
- Time: O(1) - select with default never blocks
- Space: O(1)

Three-Input Iteration Table:

Input 1: Unbuffered channel, no receiver
  ch: unbuffered, no goroutine receiving
  select: send would block
  default: executes
  return: false

Input 2: Buffered channel with space
  ch: buffered, space available
  select: send case executes
  return: true

Input 3: Buffered channel full
  ch: buffered, full
  select: send would block
  default: executes
  return: false
*/

func NonBlockingSend(ch chan<- int, value int) bool {
	// Select with default makes this non-blocking
	select {
	case ch <- value:
		// Send succeeded
		return true
	default:
		// Channel not ready to receive (would block)
		return false
	}
}

// ============================================================================
// Exercise 3: FanIn
// ============================================================================

/*
Problem: Merge multiple channels into one

Architecture:
- Start goroutine for each input channel
- Each goroutine forwards values to output
- Use WaitGroup to track completion
- Close output when all inputs exhausted

Complexity:
- Time: O(n) where n = total values across all channels
- Space: O(m) where m = number of input channels (goroutines)

Three-Input Iteration Table:

Input 1: Two channels, 3 values each
  ch1: 1, 2, 3
  ch2: 4, 5, 6
  goroutine1: forwards 1, 2, 3
  goroutine2: forwards 4, 5, 6
  output: receives all 6 values (order may vary)

Input 2: One channel closes early
  ch1: 1, 2, close
  ch2: 3, 4, 5, 6, close
  goroutine1: forwards 1, 2, exits
  goroutine2: forwards 3, 4, 5, 6, exits
  wg.Wait(): waits for both
  output: closes after all values sent

Input 3: Empty channels
  ch1: close immediately
  ch2: close immediately
  goroutine1: exits immediately
  goroutine2: exits immediately
  output: closes (no values sent)
*/

func FanIn(channels ...<-chan int) <-chan int {
	// Create output channel
	out := make(chan int)

	// WaitGroup to track completion of all input channels
	var wg sync.WaitGroup

	// Start goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()

			// Forward all values from input to output
			for v := range c {
				out <- v
			}
		}(ch)
	}

	// Close output when all inputs are exhausted
	go func() {
		wg.Wait()    // Wait for all forwarders to finish
		close(out)   // Signal to receiver that no more values coming
	}()

	return out
}

// ============================================================================
// Exercise 4: FanOut
// ============================================================================

/*
Problem: Distribute work to multiple workers

Architecture:
- Start numWorkers goroutines
- All workers read from shared input channel
- Each worker processes values and sends to shared output
- Use WaitGroup to track worker completion
- Close output when all workers done

Complexity:
- Time: O(n/w) where n = values, w = workers (speedup from parallelism)
- Space: O(w) for worker goroutines

Three-Input Iteration Table:

Input 1: 10 values, 3 workers, square function
  input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
  worker1: processes some values, squares them
  worker2: processes some values, squares them
  worker3: processes some values, squares them
  output: 1, 4, 9, 16, 25, 36, 49, 64, 81, 100 (order may vary)

Input 2: Fewer values than workers
  input: 1, 2
  workers: 5 (3 will be idle)
  worker1: might process 1
  worker2: might process 2
  workers 3-5: idle
  output: results from processed values

Input 3: Input closes immediately
  input: close immediately
  workers: all read closed channel, exit
  output: closes (no values)
*/

func FanOut(input <-chan int, numWorkers int, process func(int) int) <-chan int {
	// Create output channel
	out := make(chan int)

	// WaitGroup to track worker completion
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Read from shared input channel
			// Go's runtime ensures fair distribution among workers
			for v := range input {
				// Process value
				result := process(v)

				// Send result to output
				out <- result
			}
		}()
	}

	// Close output when all workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ============================================================================
// Exercise 5: OrChannel
// ============================================================================

/*
Problem: Create channel that closes when any input closes

Architecture:
- Start goroutine for each input channel
- Each goroutine blocks on its channel
- First one to unblock closes output
- Use sync.Once to ensure output closes only once

Complexity:
- Time: O(1) - returns immediately, goroutines wait
- Space: O(n) where n = number of input channels

Three-Input Iteration Table:

Input 1: ch2 closes first
  ch1: open
  ch2: closes at t=1s
  ch3: open
  goroutine2: unblocks, closes output
  return: output channel (will close at t=1s)

Input 2: All already closed
  ch1: closed
  ch2: closed
  ch3: closed
  goroutine1: unblocks immediately
  (or 2 or 3, race condition)
  output: closes immediately

Input 3: None close
  ch1: never closes
  ch2: never closes
  ch3: never closes
  goroutines: all block forever
  output: never closes
*/

func OrChannel(channels ...<-chan struct{}) <-chan struct{} {
	// Handle edge cases
	switch len(channels) {
	case 0:
		// No channels, return closed channel
		ch := make(chan struct{})
		close(ch)
		return ch
	case 1:
		// Single channel, just return it
		return channels[0]
	}

	// Create output channel
	out := make(chan struct{})

	// Use sync.Once to ensure we only close output once
	var once sync.Once

	// Start goroutine for each input channel
	for _, ch := range channels {
		go func(c <-chan struct{}) {
			// Block until this channel closes
			<-c

			// Close output (only first goroutine succeeds)
			once.Do(func() {
				close(out)
			})
		}(ch)
	}

	return out
}

// ============================================================================
// Exercise 6: TryReceiveAll
// ============================================================================

/*
Problem: Non-blocking receive from multiple channels

Architecture:
- Iterate through all channels
- Use select with default for each (non-blocking)
- Collect successful receives in map

Complexity:
- Time: O(n) where n = number of channels
- Space: O(k) where k = number of channels with values

Three-Input Iteration Table:

Input 1: Some channels have values
  ch[0]: has value 10
  ch[1]: empty
  ch[2]: has value 30
  select on ch[0]: receives 10
  select on ch[1]: default case
  select on ch[2]: receives 30
  return: map[0:10, 2:30]

Input 2: All channels empty
  ch[0]: empty
  ch[1]: empty
  ch[2]: empty
  all selects: default case
  return: map[] (empty)

Input 3: All channels have values
  ch[0]: has value 1
  ch[1]: has value 2
  ch[2]: has value 3
  all selects: receive case
  return: map[0:1, 1:2, 2:3]
*/

func TryReceiveAll(channels []<-chan int) map[int]int {
	// Create result map
	result := make(map[int]int)

	// Try to receive from each channel
	for i, ch := range channels {
		select {
		case v := <-ch:
			// Value available, add to result
			result[i] = v
		default:
			// No value available, skip
		}
	}

	return result
}

// ============================================================================
// Exercise 7: RateLimiter
// ============================================================================

/*
Problem: Rate limiting with token bucket algorithm

Architecture:
- Buffered channel as token bucket (capacity = rate)
- Fill with initial tokens
- Goroutine refills tokens at specified rate
- Wait() blocks on token receive
- TryWait() uses select with default

Complexity:
- Wait: O(1) - channel receive
- TryWait: O(1) - select with default
- Space: O(rate) - token buffer

Three-Input Iteration Table:

Input 1: Wait() when tokens available
  tokens: 10 available
  Wait(): receives token immediately
  tokens: 9 remaining

Input 2: Wait() when no tokens
  tokens: 0 available
  Wait(): blocks
  refill goroutine: adds token after interval
  Wait(): unblocks, receives token

Input 3: TryWait() when no tokens
  tokens: 0 available
  TryWait(): select with default
  return: false (no token available)
*/

type RateLimiter struct {
	tokens chan struct{}
	rate   int
}

func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate),
		rate:   rate,
	}

	// Fill bucket with initial tokens
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}

	// Start refill goroutine
	go rl.refill()

	return rl
}

func (rl *RateLimiter) refill() {
	// Calculate interval between token additions
	// If rate = 10/sec, interval = 100ms
	interval := time.Second / time.Duration(rl.rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// Try to add token (non-blocking)
		select {
		case rl.tokens <- struct{}{}:
			// Token added
		default:
			// Bucket full, skip
		}
	}
}

func (rl *RateLimiter) Wait() {
	// Block until token available
	<-rl.tokens
}

func (rl *RateLimiter) TryWait() bool {
	// Try to get token without blocking
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// ============================================================================
// Exercise 8: Pipeline
// ============================================================================

/*
Problem: Build multi-stage processing pipeline

Architecture:
- Stage 1: Generator (1 to n)
- Stage 2: Processor (square, fan-out)
- Stage 3: Filter (even numbers only)

Complexity:
- Time: O(n/w) where w = workers
- Space: O(w) for worker goroutines

Three-Input Iteration Table:

Input 1: n=5, workers=2
  generate: 1, 2, 3, 4, 5
  square workers: 1, 4, 9, 16, 25
  filter: 4, 16 (only even)
  output: 4, 16

Input 2: n=10, workers=3
  generate: 1-10
  square workers: 1, 4, 9, 16, 25, 36, 49, 64, 81, 100
  filter: 4, 16, 36, 64, 100
  output: 4, 16, 36, 64, 100

Input 3: n=1, workers=5
  generate: 1
  square workers: 1 (most idle)
  filter: (odd, filtered out)
  output: (empty)
*/

func Pipeline(n int, numWorkers int) <-chan int {
	// Stage 1: Generate numbers 1 to n
	generated := generate(n)

	// Stage 2: Square numbers (fan-out)
	squared := FanOut(generated, numWorkers, func(x int) int {
		return x * x
	})

	// Stage 3: Filter even numbers
	filtered := filter(squared, func(x int) bool {
		return x%2 == 0
	})

	return filtered
}

// generate produces numbers 1 to n
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

// filter only passes values where predicate returns true
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

// ============================================================================
// Exercise 9: SelectWithPriority
// ============================================================================

/*
Problem: Prefer one channel over another

Architecture:
- First select: Check high with default (non-blocking)
- If high has value, return it
- Otherwise, blocking select on both

Complexity:
- Time: O(1) - two select operations
- Space: O(1)

Three-Input Iteration Table:

Input 1: Both have values
  high: has value 1
  low: has value 2
  first select: receives from high
  return: 1, true

Input 2: Only low has value
  high: empty
  low: has value 2
  first select: default case
  second select: receives from low
  return: 2, false

Input 3: Neither has value
  high: empty
  low: empty
  first select: default case
  second select: blocks
  (waits for value from either)
*/

func SelectWithPriority(high, low <-chan int) (int, bool) {
	// First: Try high-priority channel (non-blocking)
	select {
	case v := <-high:
		return v, true
	default:
		// High not ready, proceed to normal select
	}

	// Second: Wait on both channels
	select {
	case v := <-high:
		return v, true
	case v := <-low:
		return v, false
	}
}

// ============================================================================
// Exercise 10: Timeout
// ============================================================================

/*
Problem: Receive with timeout

Architecture:
- Use select with time.After
- Return value and success flag

Complexity:
- Time: O(1) or O(timeout) depending on which happens first
- Space: O(1)

Three-Input Iteration Table:

Input 1: Value arrives before timeout
  ch: sends value after 500ms
  timeout: 1 second
  select: ch case at 500ms
  return: value, true

Input 2: Timeout occurs
  ch: sends value after 2 seconds
  timeout: 1 second
  select: timeout case at 1 second
  return: 0, false

Input 3: Value arrives immediately
  ch: has buffered value
  timeout: 1 second
  select: ch case immediately
  return: value, true
*/

func Timeout(ch <-chan int, timeout time.Duration) (int, bool) {
	select {
	case v := <-ch:
		// Received value before timeout
		return v, true
	case <-time.After(timeout):
		// Timeout occurred
		return 0, false
	}
}

/*
Alternative Implementations & Trade-offs:

1. FanIn with reflect.Select:
   - Can handle dynamic number of channels in single select
   - More complex, uses reflection
   - Slightly slower than goroutine-per-channel

2. FanOut with worker pool pattern:
   - Pre-create workers, reuse for multiple tasks
   - Better for long-running services
   - More complex lifecycle management

3. OrChannel with recursion:
   func OrChannel(channels ...<-chan struct{}) <-chan struct{} {
       switch len(channels) {
       case 0:
           return nil
       case 1:
           return channels[0]
       case 2:
           return or(channels[0], channels[1])
       }
       m := len(channels) / 2
       return or(OrChannel(channels[:m]...), OrChannel(channels[m:]...))
   }
   - Logarithmic goroutine count
   - More complex, harder to understand

Go vs X:

Go vs Python (asyncio.wait):
  done, pending = await asyncio.wait(
      [task1, task2],
      return_when=asyncio.FIRST_COMPLETED
  )
  Pros: Similar capability
  Cons: More verbose, futures instead of channels
  Go: Select is cleaner, built into language

Go vs JavaScript (Promise.race):
  const result = await Promise.race([promise1, promise2]);
  Pros: Concise for simple cases
  Cons: No equivalent to default case, can't disable cases
  Go: More flexible, can dynamically enable/disable cases

Go vs Rust (tokio::select!):
  tokio::select! {
      val = ch1.recv() => { },
      val = ch2.recv() => { },
  }
  Pros: Similar power, compile-time safety
  Cons: Macro-based, more complex syntax
  Go: Simpler, part of language

Common Mistakes to Avoid:

1. time.After in loop:
   for {
       select {
       case <-ch:
       case <-time.After(1*time.Second):  // LEAK!
       }
   }
   Creates new timer each iteration. Use time.NewTimer instead.

2. Not closing channels:
   Fan-in will never close output if you forget close(input)

3. Forgetting WaitGroup.Done:
   go func() {
       wg.Add(1)
       // Forgot defer wg.Done()
   }()
   Program hangs on wg.Wait()

4. Select on nil channel:
   var ch chan int
   select {
   case <-ch:  // Never selected
   }
   Useful for disabling cases, but easy to forget

5. Unbuffered channel in fan-out:
   Can deadlock if worker can't send result
   Use buffered channels or separate fan-in
*/
