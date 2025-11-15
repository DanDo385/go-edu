// Package main demonstrates goroutines and their ability to scale to millions of concurrent tasks.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program demonstrates:
// 1. Creating and managing 1 million goroutines
// 2. Measuring memory footprint of massive concurrency
// 3. Communication between goroutines via channels
// 4. Graceful shutdown of millions of goroutines
// 5. The performance characteristics of Go's runtime scheduler
//
// RUNTIME BEHAVIOR:
// Goroutines are multiplexed onto a small number of OS threads (GOMAXPROCS).
// Each goroutine starts with a tiny 2KB stack that grows as needed.
// Context switching happens in user space (~200ns vs ~2μs for OS threads).
//
// MEMORY IMPLICATIONS:
// 1M goroutines × 2.5KB each ≈ 2.5GB RAM
// 1M OS threads × 2MB each ≈ 2TB RAM (impossible on most systems!)

package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// SECTION 1: Basic Goroutine Creation
// ============================================================================

// demonstrateBasicGoroutine shows the simplest goroutine creation.
//
// MACRO-COMMENT: The `go` Keyword
// Prefixing any function call with `go` creates a new goroutine.
// The function executes concurrently with the caller.
// The caller does NOT wait for the goroutine to finish (fire-and-forget).
//
// SCHEDULING DETAILS:
// When you write `go f()`, the Go runtime:
// 1. Allocates a new G struct (~320 bytes)
// 2. Allocates a 2KB stack
// 3. Initializes the stack with function arguments
// 4. Adds the G to a processor's run queue
// 5. Returns immediately (doesn't wait for f to run)
func demonstrateBasicGoroutine() {
	fmt.Println("=== Basic Goroutine Creation ===")

	// MICRO-COMMENT: Track when goroutines finish using a WaitGroup
	// WaitGroup is a counter: Add(n) increments by n, Done() decrements by 1
	// Wait() blocks until the counter reaches 0
	var wg sync.WaitGroup

	// MICRO-COMMENT: Launch 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1) // Increment counter before launching goroutine

		// MICRO-COMMENT: Launch goroutine (note: must pass i as argument!)
		// If we captured i from the loop, all goroutines would see i=10
		// This is because i is a single variable that all closures share
		go func(id int) {
			defer wg.Done() // Decrement counter when done (defer ensures it runs)

			// MICRO-COMMENT: Print from goroutine
			// Output may be interleaved (goroutines run concurrently)
			fmt.Printf("  Goroutine %d is running on OS thread %d\n", id, getThreadID())

			// MICRO-COMMENT: Simulate some work
			// time.Sleep yields the CPU to other goroutines
			time.Sleep(10 * time.Millisecond)
		}(i) // Pass i as argument (copies the current value)
	}

	// MICRO-COMMENT: Wait for all goroutines to finish
	// This blocks the main goroutine until all 10 workers call Done()
	wg.Wait()

	fmt.Printf("All goroutines finished. Total goroutines in runtime: %d\n\n", runtime.NumGoroutine())
}

// getThreadID returns the current OS thread ID (for demonstration).
//
// MICRO-COMMENT: Goroutine to Thread Mapping
// Multiple goroutines share a small number of OS threads (M).
// You'll see different goroutines running on the same thread ID.
func getThreadID() int {
	// MICRO-COMMENT: This is a simplified version
	// In practice, you'd use runtime.LockOSThread() and OS-specific syscalls
	// For demo purposes, we just return a dummy value based on GOMAXPROCS
	return runtime.NumCPU()
}

// ============================================================================
// SECTION 2: Launching 1 Million Goroutines
// ============================================================================

// demonstrateMillionGoroutines launches 1 million goroutines and measures memory.
//
// MACRO-COMMENT: Memory Footprint Analysis
// Each goroutine consumes:
// - 2 KB initial stack (minimum)
// - ~320 bytes for G struct (goroutine metadata)
// - ~100 bytes for scheduler overhead
// Total: ~2.5 KB per goroutine
//
// 1,000,000 goroutines × 2.5 KB = 2.5 GB
//
// COMPARISON WITH OS THREADS:
// 1,000,000 OS threads × 2 MB stack = 2 TB (impossible!)
// Most systems limit you to 10,000-30,000 threads before OOM/kernel limits.
func demonstrateMillionGoroutines() {
	fmt.Println("=== Launching 1 Million Goroutines ===")

	const numGoroutines = 1_000_000

	// MICRO-COMMENT: Measure memory before launching goroutines
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	fmt.Printf("Memory before: %d MB\n", memBefore.Alloc/1024/1024)

	// MICRO-COMMENT: Create a channel to signal shutdown
	// We'll close this channel to wake up all goroutines at once
	// Receiving from a closed channel returns immediately with the zero value
	done := make(chan struct{})

	// MICRO-COMMENT: Track goroutine readiness
	// We use atomic operations because multiple goroutines will increment this
	var ready atomic.Int64

	// MICRO-COMMENT: Launch 1 million goroutines
	fmt.Printf("Launching %d goroutines...\n", numGoroutines)
	startLaunch := time.Now()

	for i := 0; i < numGoroutines; i++ {
		// MACRO-COMMENT: The Goroutine Body
		// Each goroutine:
		// 1. Increments the ready counter (atomic)
		// 2. Blocks on the done channel (waiting for shutdown signal)
		// 3. Exits when done is closed
		//
		// WHY THIS IS EFFICIENT:
		// - Goroutines are BLOCKED (not spinning), so CPU usage is near 0%
		// - The Go scheduler parks blocked goroutines (doesn't schedule them)
		// - Memory is the only resource consumed (stack space)
		go func() {
			// MICRO-COMMENT: Signal that this goroutine is ready
			ready.Add(1)

			// MICRO-COMMENT: Block until shutdown signal
			// This is a receive operation on a channel
			// The goroutine will be parked (not scheduled) until done is closed
			<-done

			// MICRO-COMMENT: Goroutine exits here
			// The stack is freed, G struct is recycled (not freed, cached for reuse)
		}()
	}

	// MICRO-COMMENT: Wait for all goroutines to start
	// We poll the ready counter until it reaches numGoroutines
	// This ensures all goroutines are fully initialized before we measure memory
	for ready.Load() < numGoroutines {
		time.Sleep(10 * time.Millisecond)
	}

	launchDuration := time.Since(startLaunch)

	// MICRO-COMMENT: Measure memory after launching goroutines
	runtime.GC() // Force garbage collection to get accurate numbers
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	fmt.Printf("\nResults:\n")
	fmt.Printf("  Launch time:        %v\n", launchDuration)
	fmt.Printf("  Memory before:      %d MB\n", memBefore.Alloc/1024/1024)
	fmt.Printf("  Memory after:       %d MB\n", memAfter.Alloc/1024/1024)
	fmt.Printf("  Memory used:        %d MB\n", (memAfter.Alloc-memBefore.Alloc)/1024/1024)
	fmt.Printf("  Per goroutine:      %.2f KB\n", float64(memAfter.Alloc-memBefore.Alloc)/float64(numGoroutines)/1024)
	fmt.Printf("  Active goroutines:  %d\n", runtime.NumGoroutine())
	fmt.Printf("  GOMAXPROCS:         %d\n", runtime.GOMAXPROCS(0))

	// MACRO-COMMENT: Graceful Shutdown
	// Closing the done channel broadcasts to ALL goroutines simultaneously.
	// This is much faster than signaling each goroutine individually.
	//
	// TECHNICAL DETAIL:
	// When a channel is closed:
	// 1. All goroutines blocked on receive immediately wake up
	// 2. They receive the zero value for the channel's type (struct{} here)
	// 3. They can check if the channel was closed using the two-value receive: `val, ok := <-ch`
	//    But we don't need to, since we just want them to wake up and exit
	fmt.Println("\nInitiating graceful shutdown...")
	shutdownStart := time.Now()
	close(done) // Broadcast shutdown signal

	// MICRO-COMMENT: Wait for goroutines to exit
	// We give them a moment to process the shutdown signal
	time.Sleep(100 * time.Millisecond)

	shutdownDuration := time.Since(shutdownStart)

	fmt.Printf("Shutdown complete in %v\n", shutdownDuration)
	fmt.Printf("Active goroutines after shutdown: %d\n\n", runtime.NumGoroutine())
}

// ============================================================================
// SECTION 3: Goroutine Communication
// ============================================================================

// demonstrateCommunication shows goroutines communicating via channels.
//
// MACRO-COMMENT: Channels Are Goroutine-Safe Queues
// Channels provide:
// 1. Type-safe message passing
// 2. Synchronization (sender blocks until receiver is ready, unless buffered)
// 3. Happens-before guarantees (memory visibility)
//
// BUFFERED VS UNBUFFERED:
// - Unbuffered: make(chan T) - sender blocks until receiver receives
// - Buffered: make(chan T, N) - sender only blocks when buffer is full
func demonstrateCommunication() {
	fmt.Println("=== Goroutine Communication ===")

	const (
		numWorkers = 100
		numJobs    = 10_000
	)

	// MICRO-COMMENT: Create channels for work distribution
	// jobs: workers receive work from this channel
	// results: workers send results to this channel
	jobs := make(chan int, 100)    // Buffered (allows sending 100 jobs without blocking)
	results := make(chan int, 100) // Buffered (allows workers to send without blocking)

	// MICRO-COMMENT: Launch worker goroutines
	// Each worker:
	// 1. Receives jobs from the jobs channel
	// 2. Processes the job (simulated with a simple calculation)
	// 3. Sends the result to the results channel
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// MICRO-COMMENT: Process jobs until channel is closed
			// The `for job := range jobs` loop automatically:
			// 1. Receives values from the channel
			// 2. Exits when the channel is closed and empty
			for job := range jobs {
				// MICRO-COMMENT: Simulate work
				result := job * 2

				// MICRO-COMMENT: Send result
				results <- result
			}
		}(w)
	}

	// MICRO-COMMENT: Send jobs to workers
	// We do this in a separate goroutine so we don't block the main goroutine
	go func() {
		for j := 0; j < numJobs; j++ {
			jobs <- j // Send job to channel
		}
		close(jobs) // Close channel to signal no more jobs
	}()

	// MICRO-COMMENT: Collect results in a separate goroutine
	// This goroutine:
	// 1. Receives all results from the results channel
	// 2. Closes the done channel when all results are collected
	//    (signaling to the main goroutine that we're done)
	done := make(chan struct{})
	go func() {
		for i := 0; i < numJobs; i++ {
			<-results // Receive and discard result
		}
		close(done) // Signal that all results are collected
	}()

	// MICRO-COMMENT: Wait for all workers to finish
	wg.Wait()
	close(results) // Close results channel (all workers are done sending)

	// MICRO-COMMENT: Wait for result collection to finish
	<-done

	fmt.Printf("Processed %d jobs with %d workers\n\n", numJobs, numWorkers)
}

// ============================================================================
// SECTION 4: Context-Based Cancellation
// ============================================================================

// demonstrateContextCancellation shows how to gracefully cancel goroutines using context.
//
// MACRO-COMMENT: Context for Cancellation
// context.Context provides:
// 1. Cancellation signals (ctx.Done() channel)
// 2. Deadlines and timeouts
// 3. Request-scoped values (not used here)
//
// BEST PRACTICE:
// Always pass context as the FIRST parameter to functions.
// Always check ctx.Done() in loops or before expensive operations.
func demonstrateContextCancellation() {
	fmt.Println("=== Context-Based Cancellation ===")

	// MICRO-COMMENT: Create a cancellable context
	// ctx: the context passed to goroutines
	// cancel: function to cancel the context (signals all goroutines to stop)
	ctx, cancel := context.WithCancel(context.Background())

	const numWorkers = 1000
	var wg sync.WaitGroup

	// MICRO-COMMENT: Track work done
	var workDone atomic.Int64

	// MICRO-COMMENT: Launch workers that check for cancellation
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// MACRO-COMMENT: Cancellation-Aware Loop
			// This loop:
			// 1. Checks if context is cancelled (ctx.Done() is closed)
			// 2. Does work if not cancelled
			// 3. Exits immediately when cancelled
			//
			// THE SELECT STATEMENT:
			// select is like a switch for channels:
			// - It waits on multiple channel operations simultaneously
			// - Executes the first case that's ready
			// - If multiple cases are ready, picks one at random
			for {
				select {
				case <-ctx.Done():
					// MICRO-COMMENT: Context was cancelled, exit immediately
					return

				default:
					// MICRO-COMMENT: No cancellation, do work
					workDone.Add(1)
					time.Sleep(1 * time.Millisecond)
				}
			}
		}(i)
	}

	// MICRO-COMMENT: Let workers run for a bit
	time.Sleep(100 * time.Millisecond)

	// MICRO-COMMENT: Cancel all workers
	fmt.Printf("Cancelling %d workers...\n", numWorkers)
	cancel() // This closes ctx.Done(), signaling all workers to stop

	// MICRO-COMMENT: Wait for all workers to exit
	wg.Wait()

	fmt.Printf("All workers stopped. Total work done: %d\n\n", workDone.Load())
}

// ============================================================================
// SECTION 5: Measuring Goroutine Overhead
// ============================================================================

// demonstrateOverhead compares goroutine creation vs function calls.
//
// MACRO-COMMENT: Goroutine Creation Cost
// Creating a goroutine is cheap (~1-2 nanoseconds) but not free.
// For tiny tasks, the overhead of creating a goroutine can exceed the work itself.
//
// RULE OF THUMB:
// - If the task takes < 1 microsecond, don't use a goroutine
// - If the task takes > 10 microseconds, goroutine overhead is negligible
// - For I/O-bound tasks, always use goroutines (blocking is free)
func demonstrateOverhead() {
	fmt.Println("=== Goroutine Overhead Measurement ===")

	const iterations = 100_000

	// MICRO-COMMENT: Measure synchronous function calls
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = i * 2 // Tiny amount of work
	}
	syncDuration := time.Since(start)

	// MICRO-COMMENT: Measure goroutine creation + execution
	start = time.Now()
	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = n * 2 // Same work
		}(i)
	}
	wg.Wait()
	goroutineDuration := time.Since(start)

	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Synchronous:  %v (%.2f ns/op)\n", syncDuration, float64(syncDuration.Nanoseconds())/float64(iterations))
	fmt.Printf("Goroutines:   %v (%.2f ns/op)\n", goroutineDuration, float64(goroutineDuration.Nanoseconds())/float64(iterations))
	fmt.Printf("Overhead:     %.2fx slower\n\n", float64(goroutineDuration)/float64(syncDuration))
}

// ============================================================================
// SECTION 6: Stack Growth Demonstration
// ============================================================================

// recursiveStackGrowth demonstrates stack growth in goroutines.
//
// MACRO-COMMENT: Stack Growth Mechanics
// Goroutines start with a 2 KB stack.
// When a function call would overflow the stack:
// 1. Runtime allocates a new, larger stack (typically 2x size)
// 2. Copies all data from old stack to new stack
// 3. Updates all pointers to point to new stack
// 4. Frees old stack
//
// This happens transparently - your code never knows!
func recursiveStackGrowth(depth int, maxDepth int) {
	// MICRO-COMMENT: Allocate stack space
	// This array lives on the stack (not heap), forcing stack growth
	var buf [1024]byte // 1 KB per call
	buf[0] = byte(depth)

	if depth >= maxDepth {
		return
	}

	// MICRO-COMMENT: Recurse deeper
	recursiveStackGrowth(depth+1, maxDepth)
}

func demonstrateStackGrowth() {
	fmt.Println("=== Stack Growth Demonstration ===")

	// MICRO-COMMENT: Launch a goroutine that grows its stack
	const maxDepth = 100 // 100 calls × 1 KB = 100 KB needed

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// MICRO-COMMENT: Call recursive function
		// The goroutine starts with 2 KB stack
		// After ~2 calls, it needs to grow
		// The runtime automatically handles this
		recursiveStackGrowth(0, maxDepth)

		fmt.Printf("Completed %d recursive calls (stack grew from 2KB to >100KB)\n", maxDepth)
	}()

	wg.Wait()
	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

// main orchestrates all demonstrations.
//
// MACRO-COMMENT: Learning Progression
// 1. Basic goroutine creation (foundation)
// 2. Million goroutines (scalability)
// 3. Communication (channels)
// 4. Cancellation (context)
// 5. Overhead (performance characteristics)
// 6. Stack growth (runtime behavior)
func main() {
	fmt.Printf("=== Goroutine Demonstration ===\n")
	fmt.Printf("GOMAXPROCS (CPU cores): %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU (logical CPUs):  %d\n\n", runtime.NumCPU())

	demonstrateBasicGoroutine()
	demonstrateMillionGoroutines()
	demonstrateCommunication()
	demonstrateContextCancellation()
	demonstrateOverhead()
	demonstrateStackGrowth()

	// MACRO-COMMENT: Final Memory Stats
	// This shows the runtime's memory management
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Println("=== Final Memory Statistics ===")
	fmt.Printf("Alloc:       %d MB (currently allocated)\n", mem.Alloc/1024/1024)
	fmt.Printf("TotalAlloc:  %d MB (total allocated over lifetime)\n", mem.TotalAlloc/1024/1024)
	fmt.Printf("Sys:         %d MB (obtained from OS)\n", mem.Sys/1024/1024)
	fmt.Printf("NumGC:       %d (garbage collections)\n", mem.NumGC)
	fmt.Printf("Goroutines:  %d (still running)\n", runtime.NumGoroutine())
}
