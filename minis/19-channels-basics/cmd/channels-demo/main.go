// Package main demonstrates channels and their various use cases.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program demonstrates:
// 1. Unbuffered vs buffered channels (synchronization vs queuing)
// 2. Send and receive operations (blocking semantics)
// 3. Channel closing and detection (signaling completion)
// 4. Range over channels (consuming until closed)
// 5. Select statement (multiplexing multiple channels)
// 6. Pipeline patterns (composing concurrent stages)
// 7. Common patterns (timeouts, cancellation, fan-out/fan-in)
//
// CHANNELS VS SHARED MEMORY:
// Channels provide type-safe, synchronized communication between goroutines.
// Unlike shared memory + mutexes, channels make data flow explicit and
// prevent many concurrency bugs by design.
//
// "Don't communicate by sharing memory; share memory by communicating."

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ============================================================================
// SECTION 1: Unbuffered Channels
// ============================================================================

// demonstrateUnbufferedChannel shows synchronous send/receive behavior.
//
// MACRO-COMMENT: Unbuffered Channel Semantics
// An unbuffered channel (created with make(chan T)) has ZERO capacity.
// This means:
// 1. Sends BLOCK until a receiver is ready
// 2. Receives BLOCK until a sender sends
// 3. Send and receive happen SIMULTANEOUSLY (synchronization point)
//
// VISUAL:
//   Sender              Receiver
//     |                    |
//     | ch <- 42           |
//     | (BLOCKS)           |
//     |                    | value := <-ch
//     | ← BOTH COMPLETE →  | (receives 42)
//     |                    |
//
// USE CASES:
// - Signaling events (completion, errors)
// - Synchronizing goroutines (handoff pattern)
// - Ensuring work is done before proceeding
func demonstrateUnbufferedChannel() {
	fmt.Println("=== Unbuffered Channels ===")

	// MICRO-COMMENT: Create unbuffered channel
	// The absence of a capacity argument means capacity = 0
	ch := make(chan int)

	// MACRO-COMMENT: Launch sender goroutine
	// This goroutine will BLOCK on send until the receiver is ready.
	// If we didn't launch this in a goroutine, the main goroutine would
	// block forever (deadlock) because there's no receiver yet.
	go func() {
		fmt.Println("  Sender: About to send (will block until receiver ready)...")
		ch <- 42 // BLOCKS until main goroutine receives
		fmt.Println("  Sender: Send complete (receiver received the value)")
	}()

	// MICRO-COMMENT: Give sender time to start (for demo purposes)
	// In real code, you don't need this - channels handle synchronization
	time.Sleep(100 * time.Millisecond)

	// MICRO-COMMENT: Receive from channel
	// This will unblock the sender and complete the handoff
	fmt.Println("  Receiver: About to receive (will unblock sender)...")
	value := <-ch
	fmt.Println("  Receiver: Received", value)

	// MICRO-COMMENT: Wait for sender to print completion message
	time.Sleep(50 * time.Millisecond)
	fmt.Println()
}

// ============================================================================
// SECTION 2: Buffered Channels
// ============================================================================

// demonstrateBufferedChannel shows asynchronous send/receive with buffering.
//
// MACRO-COMMENT: Buffered Channel Semantics
// A buffered channel (created with make(chan T, N)) has capacity N.
// This means:
// 1. Sends BLOCK only when buffer is FULL
// 2. Receives BLOCK only when buffer is EMPTY
// 3. Send and receive are DECOUPLED (sender doesn't wait for receiver)
//
// BUFFER STATE DIAGRAM:
//   make(chan int, 3) → [_, _, _] (3 empty slots)
//   ch <- 1           → [1, _, _] (send succeeds immediately)
//   ch <- 2           → [1, 2, _] (send succeeds immediately)
//   ch <- 3           → [1, 2, 3] (send succeeds, buffer now FULL)
//   ch <- 4           → BLOCKS (buffer full, must wait for receive)
//   v := <-ch         → [_, 2, 3] (receive gets 1, send 4 can proceed)
//
// USE CASES:
// - Work queues (bounded capacity prevents memory explosion)
// - Smoothing out producer/consumer speed differences
// - Batching (accumulate N items before processing)
func demonstrateBufferedChannel() {
	fmt.Println("=== Buffered Channels ===")

	// MICRO-COMMENT: Create buffered channel with capacity 3
	// This channel can hold 3 values before blocking sends
	ch := make(chan int, 3)

	// MICRO-COMMENT: Send 3 values (all succeed immediately)
	fmt.Println("  Sending 3 values to buffered channel (capacity 3)...")
	ch <- 1
	fmt.Println("    Sent 1 (buffer: [1, _, _])")
	ch <- 2
	fmt.Println("    Sent 2 (buffer: [1, 2, _])")
	ch <- 3
	fmt.Println("    Sent 3 (buffer: [1, 2, 3]) - buffer now FULL")

	// MACRO-COMMENT: Attempting to send to a full buffer
	// If we tried `ch <- 4` here in the main goroutine, it would DEADLOCK
	// because there's no receiver to make space in the buffer.
	// Instead, we'll launch a receiver goroutine first.
	fmt.Println("  Buffer is full. Launching receiver to make space...")

	// MICRO-COMMENT: Receive one value to make space
	go func() {
		time.Sleep(100 * time.Millisecond) // Delay to show the blocking
		value := <-ch
		fmt.Println("    Received", value, "(buffer: [2, 3, _]) - space available!")
	}()

	// MICRO-COMMENT: This send will block until the receiver makes space
	fmt.Println("  Attempting to send 4 (will block until space available)...")
	ch <- 4
	fmt.Println("  Sent 4 (buffer: [2, 3, 4])")

	// MICRO-COMMENT: Receive remaining values
	fmt.Println("  Receiving remaining values...")
	for i := 0; i < 3; i++ {
		fmt.Println("    Received", <-ch)
	}
	fmt.Println()
}

// ============================================================================
// SECTION 3: Closing Channels
// ============================================================================

// demonstrateClosingChannels shows how to close channels and detect closure.
//
// MACRO-COMMENT: Channel Closing Semantics
// Closing a channel signals "no more values will be sent".
//
// KEY RULES:
// 1. Only the SENDER should close (receiver doesn't know if sender is done)
// 2. Closing is OPTIONAL (GC will collect unclosed channels)
// 3. Close EXACTLY ONCE (closing a closed channel panics)
// 4. Never send after close (panics)
//
// BEHAVIOR AFTER CLOSING:
// - Receive returns zero value + false: value, ok := <-ch (ok == false)
// - Send PANICS: ch <- 42 → panic: send on closed channel
// - Close again PANICS: close(ch) → panic: close of closed channel
// - Range exits: for v := range ch { ... } (loop ends)
//
// USE CASES:
// - Signal completion to receivers
// - Exit range loops gracefully
// - Broadcast to multiple receivers (close wakes ALL waiters)
func demonstrateClosingChannels() {
	fmt.Println("=== Closing Channels ===")

	ch := make(chan int, 3)

	// MICRO-COMMENT: Send some values
	ch <- 1
	ch <- 2
	ch <- 3

	// MICRO-COMMENT: Close the channel
	// This signals "no more values will be sent"
	// Receivers can still drain buffered values
	close(ch)
	fmt.Println("  Channel closed (buffered values still available)")

	// MICRO-COMMENT: Receive buffered values (still works after close)
	v1 := <-ch
	v2 := <-ch
	v3 := <-ch
	fmt.Printf("  Received buffered values: %d, %d, %d\n", v1, v2, v3)

	// MACRO-COMMENT: Detecting Closure with Two-Value Receive
	// After draining buffered values, receives return zero value + false
	v4, ok := <-ch
	fmt.Printf("  After draining: value=%d, ok=%t (channel closed and empty)\n", v4, ok)

	// MICRO-COMMENT: Additional receives keep returning zero value
	v5 := <-ch
	v6 := <-ch
	fmt.Printf("  More receives: %d, %d (always zero value)\n", v5, v6)

	// MACRO-COMMENT: What Happens If We Send After Close?
	// Uncommenting the next line would PANIC:
	// ch <- 99  // panic: send on closed channel

	fmt.Println()
}

// demonstrateBroadcastViaClose shows closing as a broadcast mechanism.
//
// MICRO-COMMENT: Close as Broadcast
// Closing a channel wakes up ALL goroutines waiting on it.
// This is more efficient than sending N values to N goroutines.
func demonstrateBroadcastViaClose() {
	fmt.Println("=== Broadcast via Close ===")

	// MICRO-COMMENT: Signal channel (empty struct uses 0 bytes)
	// We don't care about the value, only the signal
	done := make(chan struct{})

	// MICRO-COMMENT: Launch 5 goroutines, all waiting on done
	for i := 1; i <= 5; i++ {
		go func(id int) {
			<-done // Block until done is closed
			fmt.Printf("  Worker %d: Received shutdown signal\n", id)
		}(i)
	}

	// MICRO-COMMENT: Let workers start
	time.Sleep(100 * time.Millisecond)

	// MICRO-COMMENT: Broadcast shutdown by closing
	// All 5 goroutines unblock simultaneously!
	fmt.Println("  Broadcasting shutdown to all workers...")
	close(done)

	// MICRO-COMMENT: Wait for workers to print
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// ============================================================================
// SECTION 4: Range Over Channels
// ============================================================================

// demonstrateRangeOverChannel shows the range loop pattern.
//
// MACRO-COMMENT: Range on Channels
// The `for range` loop on a channel:
// 1. Receives values from the channel
// 2. Loops until the channel is CLOSED
// 3. Exits when closed AND empty
//
// EQUIVALENT CODE:
//   for v := range ch { ... }
//
// Is the same as:
//   for {
//       v, ok := <-ch
//       if !ok {
//           break  // Channel closed
//       }
//       // Use v...
//   }
//
// CRITICAL: If the channel is NEVER closed, range blocks forever (goroutine leak)
func demonstrateRangeOverChannel() {
	fmt.Println("=== Range Over Channels ===")

	ch := make(chan int)

	// MACRO-COMMENT: Producer goroutine
	// Sends values and CLOSES the channel when done
	// Closing is ESSENTIAL for the range loop to exit
	go func() {
		fmt.Println("  Producer: Sending values 1..5")
		for i := 1; i <= 5; i++ {
			ch <- i
			time.Sleep(50 * time.Millisecond) // Slow down for visibility
		}
		fmt.Println("  Producer: Closing channel")
		close(ch) // Signal completion
	}()

	// MICRO-COMMENT: Consumer using range
	// This loop exits when ch is closed
	fmt.Println("  Consumer: Waiting for values...")
	for value := range ch {
		fmt.Printf("  Consumer: Received %d\n", value)
	}
	fmt.Println("  Consumer: Channel closed, exiting loop")
	fmt.Println()
}

// ============================================================================
// SECTION 5: Select Statement
// ============================================================================

// demonstrateSelect shows multiplexing multiple channel operations.
//
// MACRO-COMMENT: The Select Statement
// Select waits on multiple channel operations simultaneously.
//
// BEHAVIOR:
// 1. Evaluates all cases
// 2. If one or more are READY, picks one at RANDOM and executes it
// 3. If none are ready, BLOCKS until one becomes ready
// 4. If there's a DEFAULT case, executes it immediately (non-blocking)
//
// SYNTAX:
//   select {
//   case v := <-ch1:
//       // Received from ch1
//   case ch2 <- value:
//       // Sent to ch2
//   case <-time.After(1*time.Second):
//       // Timeout
//   default:
//       // No cases ready (non-blocking)
//   }
//
// USE CASES:
// - Timeouts (select with time.After)
// - Cancellation (select with context.Done())
// - Non-blocking operations (select with default)
// - Multiplexing inputs (select on multiple receive channels)
func demonstrateSelect() {
	fmt.Println("=== Select Statement ===")

	ch1 := make(chan string)
	ch2 := make(chan int)

	// MICRO-COMMENT: Send to ch1 after delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "hello from ch1"
	}()

	// MICRO-COMMENT: Send to ch2 after delay
	go func() {
		time.Sleep(150 * time.Millisecond)
		ch2 <- 42
	}()

	// MACRO-COMMENT: Select on multiple channels
	// Whichever channel sends first wins
	fmt.Println("  Waiting on ch1 or ch2 (whichever is ready first)...")
	select {
	case msg := <-ch1:
		fmt.Printf("  Received from ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("  Received from ch2: %d\n", num)
	}

	// MICRO-COMMENT: Receive from the other channel
	// (one is still pending)
	fmt.Println("  Waiting on remaining channel...")
	select {
	case msg := <-ch1:
		fmt.Printf("  Received from ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("  Received from ch2: %d\n", num)
	}

	fmt.Println()
}

// demonstrateSelectWithTimeout shows timeout pattern.
//
// MICRO-COMMENT: Timeout Pattern
// time.After(d) returns a channel that sends after duration d.
// Combine with select to implement timeouts.
func demonstrateSelectWithTimeout() {
	fmt.Println("=== Select with Timeout ===")

	ch := make(chan int)

	// MICRO-COMMENT: Send after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		ch <- 42
	}()

	// MICRO-COMMENT: Try to receive with 1-second timeout
	fmt.Println("  Waiting for value with 1-second timeout...")
	select {
	case v := <-ch:
		fmt.Printf("  Received: %d\n", v)
	case <-time.After(1 * time.Second):
		fmt.Println("  Timeout! (no value received within 1 second)")
	}

	fmt.Println()
}

// demonstrateSelectWithDefault shows non-blocking operations.
//
// MICRO-COMMENT: Non-Blocking Select
// If no cases are ready and there's a default case,
// the default executes immediately (no blocking).
func demonstrateSelectWithDefault() {
	fmt.Println("=== Select with Default (Non-Blocking) ===")

	ch := make(chan int)

	// MICRO-COMMENT: Try non-blocking receive
	select {
	case v := <-ch:
		fmt.Printf("  Received: %d\n", v)
	default:
		fmt.Println("  No value available (non-blocking)")
	}

	// MICRO-COMMENT: Send a value
	go func() {
		ch <- 99
	}()

	// MICRO-COMMENT: Try again after value is sent
	time.Sleep(50 * time.Millisecond)
	select {
	case v := <-ch:
		fmt.Printf("  Received: %d\n", v)
	default:
		fmt.Println("  No value available")
	}

	fmt.Println()
}

// ============================================================================
// SECTION 6: Pipeline Patterns
// ============================================================================

// demonstratePipeline shows a multi-stage processing pipeline.
//
// MACRO-COMMENT: Pipeline Pattern
// A pipeline is a series of stages connected by channels:
// 1. Each stage reads from an input channel
// 2. Performs a transformation
// 3. Writes to an output channel
// 4. Runs concurrently with other stages
//
// BENEFITS:
// - Parallelism: Stages run simultaneously
// - Modularity: Each stage is independent
// - Backpressure: Slow stages naturally slow down fast stages
//
// STRUCTURE:
//   generate → square → sum
//      ↓         ↓       ↓
//   [goroutine] [goroutine] [main]
//
// DATA FLOW:
//   generate: 1, 2, 3, 4, 5
//   square:   1, 4, 9, 16, 25
//   sum:      55
func demonstratePipeline() {
	fmt.Println("=== Pipeline Pattern ===")

	// MACRO-COMMENT: Stage 1 - Generate numbers
	// This function creates a channel and launches a goroutine that
	// sends numbers 1..n, then closes the channel to signal completion.
	//
	// PATTERN: Generator
	// Returns a RECEIVE-ONLY channel (<-chan int) to prevent callers
	// from sending to it (type safety).
	generate := func(n int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out) // ALWAYS close when done sending
			fmt.Print("  Stage 1 (generate): ")
			for i := 1; i <= n; i++ {
				out <- i
				fmt.Printf("%d ", i)
			}
			fmt.Println()
		}()
		return out
	}

	// MACRO-COMMENT: Stage 2 - Square numbers
	// This function reads from an input channel, squares each value,
	// and sends to an output channel.
	//
	// PATTERN: Transformer
	// Takes RECEIVE-ONLY input (<-chan int) and returns RECEIVE-ONLY output
	// Input type restriction: Can't accidentally send to input
	// Output type restriction: Caller can't send to our output
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out) // Close output when input is exhausted
			fmt.Print("  Stage 2 (square):   ")
			for v := range in { // Range until input is closed
				squared := v * v
				out <- squared
				fmt.Printf("%d ", squared)
			}
			fmt.Println()
		}()
		return out
	}

	// MACRO-COMMENT: Stage 3 - Sum numbers
	// This function reads all values from input and returns the sum.
	//
	// PATTERN: Reducer
	// This is a BLOCKING function (not a goroutine stage).
	// It consumes the entire input channel and returns a final result.
	sum := func(in <-chan int) int {
		total := 0
		for v := range in { // Range until input is closed
			total += v
		}
		return total
	}

	// MICRO-COMMENT: Build and execute pipeline
	fmt.Println("  Building 3-stage pipeline: generate → square → sum")

	nums := generate(5)    // Stage 1: 1, 2, 3, 4, 5
	squares := square(nums) // Stage 2: 1, 4, 9, 16, 25
	result := sum(squares)  // Stage 3: 55 (blocks until input closed)

	fmt.Printf("  Stage 3 (sum):      %d\n", result)
	fmt.Println()
}

// demonstrateFanOutFanIn shows work distribution and collection pattern.
//
// MACRO-COMMENT: Fan-Out / Fan-In Pattern
//
// FAN-OUT: Distribute work from one channel to multiple workers
//   jobs channel → [worker1, worker2, worker3, ...]
//
// FAN-IN: Collect results from multiple workers to one channel
//   [worker1, worker2, worker3, ...] → results channel
//
// USE CASE:
// Parallelize CPU-bound work by running multiple workers,
// then collect all results.
func demonstrateFanOutFanIn() {
	fmt.Println("=== Fan-Out / Fan-In Pattern ===")

	const (
		numJobs    = 10
		numWorkers = 3
	)

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// MACRO-COMMENT: Fan-Out - Launch workers
	// All workers read from the SAME jobs channel.
	// Go's runtime ensures only ONE worker receives each job.
	fmt.Printf("  Launching %d workers...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for job := range jobs {
				// MICRO-COMMENT: Simulate work
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				result := job * 2
				fmt.Printf("    Worker %d: processed job %d → result %d\n", id, job, result)
				results <- result
			}
		}(w)
	}

	// MICRO-COMMENT: Send jobs
	fmt.Printf("  Sending %d jobs...\n", numJobs)
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Signal no more jobs (workers will exit range loop)

	// MACRO-COMMENT: Fan-In - Collect results
	// We know exactly how many results to expect (numJobs).
	// Collect them all before proceeding.
	fmt.Println("  Collecting results...")
	for i := 0; i < numJobs; i++ {
		<-results // Receive and discard (just demonstrating)
	}
	fmt.Println("  All jobs complete!")
	fmt.Println()
}

// ============================================================================
// SECTION 7: Common Patterns
// ============================================================================

// demonstrateQuitChannel shows graceful shutdown pattern.
//
// MICRO-COMMENT: Quit Channel Pattern
// Use a channel to signal goroutines to stop.
// Close the channel to broadcast shutdown to all listeners.
func demonstrateQuitChannel() {
	fmt.Println("=== Quit Channel Pattern ===")

	quit := make(chan struct{})
	workDone := make(chan int)

	// MICRO-COMMENT: Worker that checks quit channel
	go func() {
		count := 0
		for {
			select {
			case <-quit:
				fmt.Println("  Worker: Received quit signal, stopping...")
				workDone <- count
				return
			default:
				// MICRO-COMMENT: Do work
				count++
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// MICRO-COMMENT: Let worker run for a bit
	fmt.Println("  Worker running...")
	time.Sleep(300 * time.Millisecond)

	// MICRO-COMMENT: Signal quit
	fmt.Println("  Sending quit signal...")
	close(quit)

	// MICRO-COMMENT: Wait for worker to finish
	count := <-workDone
	fmt.Printf("  Worker stopped after %d iterations\n", count)
	fmt.Println()
}

// demonstrateFuturePattern shows async computation with channels.
//
// MICRO-COMMENT: Future/Promise Pattern
// Return a channel immediately, compute result asynchronously.
// Caller can do other work and retrieve result later.
func demonstrateFuturePattern() {
	fmt.Println("=== Future/Promise Pattern ===")

	// MACRO-COMMENT: Future Type
	// A future is just a channel that will eventually receive a value.
	// We use a buffered channel (capacity 1) so the goroutine doesn't block.
	type Future chan int

	expensiveComputation := func() Future {
		future := make(chan int, 1) // Buffered so sender doesn't block
		go func() {
			fmt.Println("  Computing result asynchronously...")
			time.Sleep(500 * time.Millisecond) // Simulate work
			result := 42
			future <- result
			fmt.Println("  Computation complete!")
		}()
		return future
	}

	// MICRO-COMMENT: Start computation (non-blocking)
	fmt.Println("  Starting async computation...")
	future := expensiveComputation()

	// MICRO-COMMENT: Do other work while computation runs
	fmt.Println("  Doing other work...")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("  Other work complete!")

	// MICRO-COMMENT: Retrieve result (blocks if not ready)
	fmt.Println("  Waiting for async result...")
	result := <-future
	fmt.Printf("  Result: %d\n", result)
	fmt.Println()
}

// ============================================================================
// SECTION 8: Directional Channels
// ============================================================================

// demonstrateDirectionalChannels shows send-only and receive-only channels.
//
// MACRO-COMMENT: Channel Direction
// Channels can be restricted to send-only or receive-only in type signatures.
//
// TYPES:
//   chan T      Bidirectional (can send and receive)
//   chan<- T    Send-only (can only send)
//   <-chan T    Receive-only (can only receive)
//
// BENEFITS:
// 1. Type safety: Compiler prevents misuse
// 2. API clarity: Function signature shows intent
// 3. Best practice: Always use most restrictive type
//
// CONVERSION:
// Bidirectional channels implicitly convert to directional channels:
//   ch := make(chan int)         // Bidirectional
//   sendTo(ch)                    // Converts to chan<- int
//   receiveFrom(ch)               // Converts to <-chan int
func demonstrateDirectionalChannels() {
	fmt.Println("=== Directional Channels ===")

	// MICRO-COMMENT: Producer function (send-only parameter)
	// This function can ONLY send to the channel, not receive.
	// Trying to receive would be a compile error.
	producer := func(ch chan<- int) {
		fmt.Println("  Producer: Sending values...")
		for i := 1; i <= 3; i++ {
			ch <- i
			fmt.Printf("    Sent: %d\n", i)
		}
		close(ch)
		fmt.Println("  Producer: Closed channel")

		// MICRO-COMMENT: This would be a compile error:
		// value := <-ch  // Error: cannot receive from send-only channel
	}

	// MICRO-COMMENT: Consumer function (receive-only parameter)
	// This function can ONLY receive from the channel, not send.
	// Trying to send would be a compile error.
	consumer := func(ch <-chan int) {
		fmt.Println("  Consumer: Receiving values...")
		for value := range ch {
			fmt.Printf("    Received: %d\n", value)
		}
		fmt.Println("  Consumer: Channel closed")

		// MICRO-COMMENT: This would be a compile error:
		// ch <- 99  // Error: cannot send to receive-only channel

		// MICRO-COMMENT: This would also be a compile error:
		// close(ch)  // Error: cannot close receive-only channel
	}

	// MICRO-COMMENT: Create bidirectional channel
	ch := make(chan int)

	// MICRO-COMMENT: Pass to functions (implicit conversion)
	// The bidirectional channel converts to directional types
	go producer(ch) // chan int → chan<- int
	consumer(ch)    // chan int → <-chan int

	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	fmt.Println("=== Channel Basics Demonstration ===\n")

	// SECTION 1: Unbuffered vs Buffered
	demonstrateUnbufferedChannel()
	demonstrateBufferedChannel()

	// SECTION 2: Closing
	demonstrateClosingChannels()
	demonstrateBroadcastViaClose()

	// SECTION 3: Range
	demonstrateRangeOverChannel()

	// SECTION 4: Select
	demonstrateSelect()
	demonstrateSelectWithTimeout()
	demonstrateSelectWithDefault()

	// SECTION 5: Pipelines
	demonstratePipeline()
	demonstrateFanOutFanIn()

	// SECTION 6: Common Patterns
	demonstrateQuitChannel()
	demonstrateFuturePattern()

	// SECTION 7: Type Safety
	demonstrateDirectionalChannels()

	fmt.Println("=== All Demonstrations Complete ===")
}
