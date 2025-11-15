package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
Project 20: Select, Fan-In, and Fan-Out - Comprehensive Demo

This program demonstrates:
1. Select statement basics (multiplexing channels)
2. Non-blocking operations with default case
3. Timeout patterns with time.After
4. Fan-in pattern (merging multiple channels)
5. Fan-out pattern (distributing work to workers)
6. Combined pipeline (fan-out + fan-in)
7. Advanced patterns (quit channels, priority select)

Each demo is standalone and heavily commented.
*/

func main() {
	fmt.Println("=== Select, Fan-In, and Fan-Out Demo ===\n")

	// Demo 1: Select Basics
	fmt.Println("--- Demo 1: Select Statement Basics ---")
	demoSelectBasics()
	time.Sleep(500 * time.Millisecond)

	// Demo 2: Non-Blocking Operations
	fmt.Println("\n--- Demo 2: Non-Blocking Operations (Default Case) ---")
	demoNonBlocking()
	time.Sleep(500 * time.Millisecond)

	// Demo 3: Timeout Patterns
	fmt.Println("\n--- Demo 3: Timeout Patterns ---")
	demoTimeouts()
	time.Sleep(500 * time.Millisecond)

	// Demo 4: Fan-In Pattern
	fmt.Println("\n--- Demo 4: Fan-In Pattern (Merging Channels) ---")
	demoFanIn()
	time.Sleep(500 * time.Millisecond)

	// Demo 5: Fan-Out Pattern
	fmt.Println("\n--- Demo 5: Fan-Out Pattern (Worker Pool) ---")
	demoFanOut()
	time.Sleep(500 * time.Millisecond)

	// Demo 6: Pipeline (Fan-Out + Fan-In)
	fmt.Println("\n--- Demo 6: Pipeline Architecture ---")
	demoPipeline()
	time.Sleep(500 * time.Millisecond)

	// Demo 7: Quit Channel Pattern
	fmt.Println("\n--- Demo 7: Quit Channel (Graceful Shutdown) ---")
	demoQuitChannel()
	time.Sleep(500 * time.Millisecond)

	// Demo 8: Priority Select
	fmt.Println("\n--- Demo 8: Priority Select (Preferring Certain Cases) ---")
	demoPrioritySelect()

	fmt.Println("\n=== All Demos Complete ===")
}

// ============================================================================
// Demo 1: Select Statement Basics
// ============================================================================

func demoSelectBasics() {
	// The select statement lets you wait on multiple channel operations
	// It's like a switch statement for channels
	//
	// Key behaviors:
	// 1. Blocks until one case can proceed
	// 2. If multiple cases are ready, chooses one at random (fairness)
	// 3. Executes only one case per select

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Start goroutine that sends to ch1 after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "message from ch1"
	}()

	// Start goroutine that sends to ch2 after 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch2 <- "message from ch2"
	}()

	// Select will pick whichever channel sends first
	// In this case, ch2 will be ready first (500ms < 1s)
	select {
	case msg1 := <-ch1:
		fmt.Println("Received from ch1:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received from ch2:", msg2)
	}
	// Output: "Received from ch2: message from ch2"

	// Demonstrate random selection when both are ready
	fmt.Println("\nDemonstrating random selection (run multiple times):")

	ch3 := make(chan int, 1) // Buffered channels
	ch4 := make(chan int, 1)

	// Both channels have values immediately available
	ch3 <- 3
	ch4 <- 4

	// Select will randomly choose between ch3 and ch4
	// Run this program multiple times - you'll see different outputs
	select {
	case v := <-ch3:
		fmt.Printf("Randomly selected ch3: %d\n", v)
	case v := <-ch4:
		fmt.Printf("Randomly selected ch4: %d\n", v)
	}
}

// ============================================================================
// Demo 2: Non-Blocking Operations with Default Case
// ============================================================================

func demoNonBlocking() {
	// The default case makes select non-blocking
	// If no other case is ready, default executes immediately
	//
	// Use case: Poll channels without blocking

	ch := make(chan int, 1)

	// Example 1: Try to receive without blocking
	fmt.Println("Example 1: Non-blocking receive (empty channel)")
	select {
	case v := <-ch:
		fmt.Printf("Received: %d\n", v)
	default:
		fmt.Println("No value available (default case)")
	}

	// Example 2: Try to receive (channel has value)
	ch <- 42
	fmt.Println("\nExample 2: Non-blocking receive (channel has value)")
	select {
	case v := <-ch:
		fmt.Printf("Received: %d\n", v)
	default:
		fmt.Println("No value available")
	}

	// Example 3: Try to send without blocking
	ch <- 1 // Fill buffer
	fmt.Println("\nExample 3: Non-blocking send (full channel)")
	select {
	case ch <- 2:
		fmt.Println("Sent value")
	default:
		fmt.Println("Channel full, cannot send (default case)")
	}

	// Example 4: Polling loop
	fmt.Println("\nExample 4: Polling pattern")
	done := make(chan bool)

	go func() {
		time.Sleep(2 * time.Second)
		done <- true
	}()

	// Poll until done signal received
	for i := 0; i < 5; i++ {
		select {
		case <-done:
			fmt.Println("Done signal received!")
			return
		default:
			fmt.Printf("Polling... (iteration %d)\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// ============================================================================
// Demo 3: Timeout Patterns
// ============================================================================

func demoTimeouts() {
	// Timeout patterns prevent operations from hanging indefinitely
	//
	// time.After(duration) returns a channel that receives a value
	// after the specified duration

	// Example 1: Operation that times out
	fmt.Println("Example 1: Operation timeout")
	ch1 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "slow response"
	}()

	select {
	case msg := <-ch1:
		fmt.Println("Received:", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout! Operation took too long")
	}

	// Example 2: Operation that completes in time
	fmt.Println("\nExample 2: Operation completes in time")
	ch2 := make(chan string)

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch2 <- "fast response"
	}()

	select {
	case msg := <-ch2:
		fmt.Println("Received:", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout!")
	}

	// Example 3: Multiple operations with timeout
	fmt.Println("\nExample 3: Multiple operations with shared timeout")
	ch3 := make(chan int)
	ch4 := make(chan int)

	go func() {
		time.Sleep(300 * time.Millisecond)
		ch3 <- 3
	}()

	go func() {
		time.Sleep(600 * time.Millisecond)
		ch4 <- 4
	}()

	timeout := time.After(500 * time.Millisecond)

	// First select: ch3 arrives at 300ms (before timeout)
	select {
	case v := <-ch3:
		fmt.Printf("Received from ch3: %d\n", v)
	case v := <-ch4:
		fmt.Printf("Received from ch4: %d\n", v)
	case <-timeout:
		fmt.Println("Timeout!")
	}

	// Second select: ch4 would arrive at 600ms, but timeout is at 500ms
	select {
	case v := <-ch3:
		fmt.Printf("Received from ch3: %d\n", v)
	case v := <-ch4:
		fmt.Printf("Received from ch4: %d\n", v)
	case <-timeout:
		fmt.Println("Timeout! (ch4 took too long)")
	}
}

// ============================================================================
// Demo 4: Fan-In Pattern (Merging Multiple Channels)
// ============================================================================

func demoFanIn() {
	// Fan-In: Merge multiple input channels into a single output channel
	//
	// Real-world analogy: Multiple checkout lines merging into single exit
	//
	// Use cases:
	// - Collecting results from multiple workers
	// - Aggregating logs from multiple services
	// - Merging search results from multiple sources

	// Create two channels that produce data at different rates
	ch1 := boring("Alice", 500*time.Millisecond)
	ch2 := boring("Bob", 700*time.Millisecond)

	// Fan-in: merge ch1 and ch2 into a single channel
	merged := fanIn(ch1, ch2)

	// Receive 10 messages from merged channel
	fmt.Println("Receiving messages from merged channel:")
	for i := 0; i < 10; i++ {
		msg := <-merged
		fmt.Println(msg)
	}

	fmt.Println("\nFan-in complete!")
}

// boring creates a channel that sends messages at specified intervals
func boring(name string, interval time.Duration) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- fmt.Sprintf("%s says: message %d", name, i)
			time.Sleep(interval)
		}
	}()

	return ch
}

// fanIn merges two channels into one
// This demonstrates the select-based fan-in pattern
func fanIn(ch1, ch2 <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		// Keep reading until both channels are closed
		// We use nil channels to disable closed channel cases
		for ch1 != nil || ch2 != nil {
			select {
			case msg, ok := <-ch1:
				if !ok {
					fmt.Println("  [ch1 closed]")
					ch1 = nil // Disable this case by setting to nil
					continue
				}
				out <- msg

			case msg, ok := <-ch2:
				if !ok {
					fmt.Println("  [ch2 closed]")
					ch2 = nil // Disable this case
					continue
				}
				out <- msg
			}
		}
	}()

	return out
}

// ============================================================================
// Demo 5: Fan-Out Pattern (Distributing Work to Workers)
// ============================================================================

func demoFanOut() {
	// Fan-Out: Distribute work from single channel to multiple workers
	//
	// Real-world analogy: Single task queue distributed to multiple workers
	//
	// Use cases:
	// - Parallel processing of jobs
	// - Load balancing
	// - Concurrent HTTP request handling

	// Create a channel of tasks
	tasks := make(chan Task, 20)

	// Generate 20 tasks
	go func() {
		defer close(tasks)
		for i := 1; i <= 20; i++ {
			tasks <- Task{
				ID:   i,
				Data: fmt.Sprintf("task-%d", i),
			}
		}
	}()

	// Fan-out: Start 3 workers to process tasks concurrently
	const numWorkers = 3
	results := make(chan Result, 20)

	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Printf("Processing tasks with %d workers:\n", numWorkers)
	processed := 0
	for result := range results {
		fmt.Printf("  Task %d processed by Worker %d: %s\n",
			result.TaskID, result.WorkerID, result.Output)
		processed++
	}

	fmt.Printf("\nProcessed %d tasks total\n", processed)
}

// Task represents a unit of work
type Task struct {
	ID   int
	Data string
}

// Result represents the output of processing a task
type Result struct {
	TaskID   int
	WorkerID int
	Output   string
}

// worker processes tasks from the tasks channel
func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		// Simulate work (random duration between 50-150ms)
		duration := time.Duration(50+rand.Intn(100)) * time.Millisecond
		time.Sleep(duration)

		// Send result
		results <- Result{
			TaskID:   task.ID,
			WorkerID: id,
			Output:   fmt.Sprintf("processed '%s' in %v", task.Data, duration),
		}
	}
}

// ============================================================================
// Demo 6: Pipeline Architecture (Fan-Out + Fan-In)
// ============================================================================

func demoPipeline() {
	// Pipeline: Chain multiple processing stages together
	//
	// Stage 1: Generate numbers
	// Stage 2: Square numbers (fan-out to 3 workers)
	// Stage 3: Filter odd numbers (fan-out to 2 workers)
	// Stage 4: Collect results (fan-in)

	fmt.Println("Building a 3-stage pipeline:")
	fmt.Println("  Stage 1: Generate numbers (1-20)")
	fmt.Println("  Stage 2: Square numbers (3 workers)")
	fmt.Println("  Stage 3: Filter odd numbers (2 workers)")
	fmt.Println("")

	// Stage 1: Generate numbers
	numbers := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)

	// Stage 2: Square numbers (fan-out to 3 workers)
	squared1 := square(numbers)
	squared2 := square(numbers)
	squared3 := square(numbers)

	// Merge squared results
	allSquared := fanInInts(squared1, squared2, squared3)

	// Stage 3: Filter odd numbers (fan-out to 2 workers)
	filtered1 := filterOdd(allSquared)
	filtered2 := filterOdd(allSquared)

	// Final fan-in
	results := fanInInts(filtered1, filtered2)

	// Collect and display results
	fmt.Println("Pipeline results (even squares only):")
	count := 0
	for v := range results {
		fmt.Printf("%d ", v)
		count++
	}
	fmt.Printf("\nProcessed %d values through pipeline\n", count)
}

// generate creates a channel and sends numbers to it
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square reads from input channel, squares values, sends to output
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// filterOdd only sends even numbers to output
func filterOdd(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

// fanInInts merges multiple int channels into one
func fanInInts(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}

	// Close output channel when all inputs are exhausted
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ============================================================================
// Demo 7: Quit Channel Pattern (Graceful Shutdown)
// ============================================================================

func demoQuitChannel() {
	// Quit channel pattern: Signal goroutines to stop gracefully
	//
	// Use case: Shutting down worker goroutines without killing them
	// abruptly (allows cleanup, flushing buffers, etc.)

	quit := make(chan struct{})
	done := make(chan struct{})

	// Start a worker that runs until quit signal
	go func() {
		defer func() { done <- struct{}{} }()

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		iteration := 0
		for {
			select {
			case <-ticker.C:
				iteration++
				fmt.Printf("  Worker: iteration %d\n", iteration)

			case <-quit:
				fmt.Println("  Worker: Quit signal received, cleaning up...")
				// Simulate cleanup
				time.Sleep(100 * time.Millisecond)
				fmt.Println("  Worker: Cleanup complete, exiting")
				return
			}
		}
	}()

	// Let worker run for a while
	fmt.Println("Worker running...")
	time.Sleep(1 * time.Second)

	// Send quit signal
	fmt.Println("\nSending quit signal...")
	close(quit) // Closing channel broadcasts to all receivers

	// Wait for worker to finish cleanup
	<-done
	fmt.Println("Worker stopped gracefully")
}

// ============================================================================
// Demo 8: Priority Select (Preferring Certain Cases)
// ============================================================================

func demoPrioritySelect() {
	// Sometimes you want to prioritize certain channels over others
	// Standard select chooses randomly, so we use a two-tier approach:
	//
	// 1. First select: Check high-priority channel with default (non-blocking)
	// 2. Second select: Check all channels normally

	high := make(chan string, 10)
	low := make(chan string, 10)
	quit := make(chan struct{})

	// Producer: Send to both channels
	go func() {
		for i := 0; i < 5; i++ {
			high <- fmt.Sprintf("HIGH-%d", i)
			low <- fmt.Sprintf("low-%d", i)
			time.Sleep(100 * time.Millisecond)
		}
		close(high)
		close(low)
		close(quit)
	}()

	fmt.Println("Processing with priority (HIGH priority over low):")

	// Consumer with priority
	for {
		// First: Always check quit channel (highest priority)
		select {
		case <-quit:
			fmt.Println("\nQuit signal received")
			return
		default:
		}

		// Second: Check high-priority channel (non-blocking)
		select {
		case msg, ok := <-high:
			if ok {
				fmt.Printf("  Processed: %s (high priority)\n", msg)
				continue // Loop again to check high again
			}
		default:
		}

		// Third: Check both channels normally
		select {
		case msg, ok := <-high:
			if ok {
				fmt.Printf("  Processed: %s\n", msg)
			}
		case msg, ok := <-low:
			if ok {
				fmt.Printf("  Processed: %s\n", msg)
			}
		case <-quit:
			fmt.Println("\nQuit signal received")
			return
		}
	}
}

/*
Key Takeaways from These Demos:

1. **Select Statement**:
   - Multiplexes multiple channel operations
   - Random selection ensures fairness
   - Default case makes it non-blocking

2. **Timeout Patterns**:
   - time.After(duration) for timeouts
   - Prevents operations from hanging
   - Can combine multiple timeouts in one select

3. **Fan-In Pattern**:
   - Merge multiple channels into one
   - Use select to read from whichever is ready
   - Set channels to nil when closed to disable them

4. **Fan-Out Pattern**:
   - Distribute work to multiple workers
   - Workers read from shared channel
   - Collect results using channels or WaitGroup

5. **Pipeline Pattern**:
   - Chain stages: generate → process → filter → collect
   - Each stage is concurrent
   - Fan-out and fan-in at each stage

6. **Quit Channel**:
   - Broadcast shutdown signal by closing channel
   - Workers check <-quit in select
   - Allows graceful cleanup

7. **Priority Select**:
   - Two-tier select for prioritization
   - First check high-priority with default
   - Then check all channels normally

Common Patterns:

✅ Good:
  - Use select for multiple channels
  - Close channels to signal completion
  - Fan-in with WaitGroup for cleanup
  - Buffered channels to prevent blocking

❌ Avoid:
  - time.After in tight loops (creates many timers)
  - Not closing channels (receivers hang)
  - Unbuffered channels in fan-out (can deadlock)
  - Forgetting to handle channel close (!ok check)
*/
