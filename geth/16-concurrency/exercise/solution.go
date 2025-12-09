//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

/*
Problem: Probe multiple endpoints concurrently using a bounded worker pool.

When building Ethereum tooling, you often need to query multiple RPC endpoints,
check health of multiple nodes, or fetch data from multiple sources. Doing this
sequentially is slow. Doing it with unbounded goroutines risks exhausting resources
or hitting rate limits. A worker pool is the idiomatic Go solution.

Computer science principles highlighted:
  - Concurrency patterns: Worker pool with channels prevents unbounded goroutine creation
  - Resource management: Bounded workers respect system limits and RPC rate limits
  - Context propagation: Timeouts and cancellation cascade through concurrent operations
  - Safe concurrent access: Mutex-protected maps prevent data races when aggregating results
*/
func Run(ctx context.Context, p Prober, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Why validate inputs? This function is a library API that will be called by
	// other code. We can't trust callers to always pass valid inputs. This is the
	// same defensive programming pattern from modules 01-stack and 06-eip1559.
	//
	// Context handling: If ctx is nil, we provide context.Background() as a safe
	// default. This ensures that context propagation works even if callers forget
	// to pass a context. Never allow nil contexts to reach RPC calls or worker
	// goroutines - they would panic or misbehave.
	//
	// This pattern repeats: Every function accepting context should validate it.
	// This is a fundamental Go idiom that prevents runtime panics.
	if ctx == nil {
		ctx = context.Background()
	}

	// Prober validation: The Prober interface is our dependency injection point.
	// Without it, we can't perform health checks. This is a critical validation
	// that follows the same pattern as client validation in previous modules.
	//
	// Dependency injection: By accepting an interface, we make this function
	// testable (mock Prober) and flexible (any health check implementation).
	// This is a key software engineering principle: "Program to interfaces, not
	// implementations."
	if p == nil {
		return nil, errors.New("prober is nil")
	}

	// Worker count validation with sensible default: 4 workers is a reasonable
	// default that balances parallelism with resource consumption. This prevents
	// creating thousands of goroutines when endpoints list is huge.
	//
	// Why 4? It's a balance:
	//   - Too few (1-2): Underutilizes CPU cores, doesn't achieve much concurrency
	//   - Too many (100+): Contention on shared resources, potential rate limits
	//   - 4-8: Sweet spot for I/O-bound operations like RPC calls
	//
	// This pattern repeats: Provide sensible defaults that work for common cases
	// while allowing advanced users to override via config.
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}

	// Timeout validation with sensible default: 5 seconds is reasonable for
	// network health checks without being too aggressive. This bounds the entire
	// operation - all probes must complete within this time.
	//
	// Why 5 seconds? Network operations timing:
	//   - Too short (< 1s): Flaky results, false negatives from slow networks
	//   - Too long (> 30s): Poor UX, hangs visible to users
	//   - 5s: Long enough for most networks, short enough for good UX
	//
	// Timeout strategy: We'll use half this timeout for per-request timeouts
	// (explained in Step 3), giving buffer for retries and overhead.
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}

	// ============================================================================
	// STEP 2: Create Child Context with Timeout - Resource Management
	// ============================================================================
	// Context with timeout: This creates a child context that will automatically
	// cancel after cfg.Timeout. This is crucial for preventing goroutine leaks
	// and ensuring the entire operation completes in bounded time.
	//
	// How it works: context.WithTimeout returns (childCtx, cancelFunc).
	//   - childCtx: Inherits from parent, adds timeout behavior
	//   - cancelFunc: Must be called to free resources (timers, goroutines)
	//
	// defer cancel(): ALWAYS defer cancel after creating context. This prevents
	// context leaks. Even though timeout will eventually fire, explicit cancellation
	// is immediate and frees resources sooner.
	//
	// Why child context? Allows this function to enforce its own timeout
	// independently of the parent context. If caller passes a 1-minute timeout
	// but cfg.Timeout is 5 seconds, we'll respect the 5-second timeout.
	//
	// Context cancellation propagation: When this context times out:
	//   1. ctx.Done() channel closes
	//   2. All goroutines select on ctx.Done() wake up
	//   3. Workers stop processing new jobs
	//   4. In-flight probes receive canceled context and abort
	//
	// This pattern repeats: Always create child contexts for operations with
	// their own timeout requirements. This is fundamental to Go's cancellation
	// model.
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	// ============================================================================
	// STEP 3: Create Jobs Channel - Worker Pool Pattern
	// ============================================================================
	// Jobs channel: This is the core of the worker pool pattern. Workers pull
	// endpoints from this channel, process them, then pull the next one.
	//
	// Buffered vs unbuffered:
	//   - Unbuffered (size 0): Sender blocks until worker receives. Wastes CPU cycles.
	//   - Buffered (size > 0): Sender can enqueue jobs without blocking. More efficient.
	//
	// Why buffer size = cfg.Workers? This allows the job producer to enqueue one
	// job per worker without blocking. Since workers pull continuously, this keeps
	// the pipeline full and maximizes throughput.
	//
	// Channel as queue: This implements a work-stealing queue pattern. As soon
	// as a worker finishes, it steals the next job from the queue. This naturally
	// load-balances - faster workers process more jobs.
	//
	// Memory consideration: Each buffered channel slot holds a string (endpoint).
	// For 4 workers, this is negligible. For 1000+ workers, consider unbuffered
	// channels to limit memory.
	jobs := make(chan string, cfg.Workers)

	// Mutex for result aggregation: Maps in Go are not safe for concurrent access.
	// Multiple goroutines writing to the same map simultaneously causes data races
	// and panics. A mutex provides exclusive access - only one goroutine can hold
	// the lock at a time.
	//
	// Why mutex instead of channels? For aggregation, mutex is simpler and more
	// efficient than channels. Channels are better for streaming data; mutex is
	// better for shared state.
	//
	// Alternative design: Each worker could send results to a channel, and a
	// single aggregator goroutine receives and aggregates. This avoids mutex but
	// adds complexity (extra goroutine, result channel, aggregation loop).
	//
	// Critical section: Lock → write map → unlock. Keep this short! Long critical
	// sections create contention and hurt parallelism.
	var mu sync.Mutex

	// Result struct: Initialized with empty maps. We'll populate these as workers
	// complete probes.
	//
	// Design choice: Return partial results on timeout. If 100 endpoints are
	// probed but only 50 complete before timeout, we return those 50 results
	// plus an error indicating timeout. This is more useful than returning nothing.
	//
	// Map initialization: make() creates empty maps. Without this, the maps would
	// be nil and writes would panic. This is a common Go gotcha - always initialize
	// maps before writing.
	res := &Result{
		Successes: make(map[string]time.Duration),
		Failures:  make(map[string]error),
	}

	// ============================================================================
	// STEP 4: Start Worker Pool - Concurrency Pattern
	// ============================================================================
	// WaitGroup: Tracks the number of active worker goroutines. Main goroutine
	// calls wg.Wait() to block until all workers signal completion with wg.Done().
	//
	// Why WaitGroup? Without it, main goroutine would return immediately while
	// workers are still running. Those workers would be killed mid-processing.
	// WaitGroup ensures graceful completion.
	//
	// How it works:
	//   1. wg.Add(n): Increments counter by n (number of workers)
	//   2. wg.Done(): Decrements counter by 1 (called by each worker on exit)
	//   3. wg.Wait(): Blocks until counter reaches 0 (all workers exited)
	//
	// Pattern: Add before launching goroutine, Done in defer, Wait in main.
	var wg sync.WaitGroup

	// Launch worker pool: Create cfg.Workers goroutines. Each worker runs
	// independently, pulling jobs from the channel until it closes.
	//
	// Why loop instead of fixed number? This allows dynamic worker count based
	// on config. For testing, might use 1 worker. For production, might use 50.
	//
	// Goroutine cost: Each goroutine has ~2KB stack. For 4 workers, that's 8KB.
	// Even 1000 workers is only 2MB. Go's goroutines are lightweight compared
	// to OS threads (8MB+ stack per thread).
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			// defer wg.Done(): ALWAYS defer Done after Add. This ensures
			// Done is called even if worker panics. Without this, wg.Wait()
			// would hang forever waiting for a worker that will never signal.
			//
			// This pattern repeats: defer cleanup functions (Done, unlock,
			// close) to ensure they execute even on panic or early return.
			defer wg.Done()

			// Range over jobs channel: Worker pulls endpoints one at a time.
			// When channel closes, range loop exits. This is the standard
			// worker pool pattern in Go.
			//
			// Why range instead of select? Range is simpler for this use case.
			// Select is needed when monitoring multiple channels or ctx.Done().
			// Here, ctx.Done() is checked in the job producer (Step 5).
			for endpoint := range jobs {
				// Record start time: We'll calculate latency after probe completes.
				// time.Now() captures current timestamp. time.Since(start) computes
				// elapsed duration.
				//
				// Why measure latency? Health checks should not only report success/
				// failure but also performance. Slow endpoints might need investigation.
				start := time.Now()

				// Per-request timeout: Create a child context with half the overall
				// timeout. This prevents one slow endpoint from consuming the entire
				// budget.
				//
				// Why half timeout? Math:
				//   - Overall timeout: 5s (for all probes)
				//   - Per-request timeout: 2.5s (for each probe)
				//   - If 4 workers each hit slow endpoints, worst case is 2.5s
				//   - This is less than 5s overall timeout, so we have buffer
				//
				// Why not overall timeout? If one probe uses the full 5s, other
				// workers sit idle. By limiting each probe to 2.5s, we ensure
				// multiple probes can run within the 5s budget.
				//
				// Context hierarchy: ctx (5s) → reqCtx (2.5s). If parent cancels,
				// child cancels immediately. If child cancels, parent is unaffected.
				reqCtx, cancelReq := context.WithTimeout(ctx, cfg.Timeout/2)

				// Call Probe: This is the actual health check. Implementation is
				// injected via the Prober interface. Could be HTTP GET, TCP connect,
				// Ethereum RPC call, etc.
				//
				// Context propagation: We pass reqCtx so Probe respects the timeout.
				// If Probe makes network calls, they should honor ctx.Done().
				err := p.Probe(reqCtx, endpoint)

				// Cancel request context: Free resources immediately. Even though
				// reqCtx will timeout eventually, explicit cancellation is faster
				// and prevents goroutine leaks in the context implementation.
				//
				// Why not defer? defer would delay cancellation until worker exits.
				// We want immediate cancellation after each probe to free resources
				// for the next probe.
				cancelReq()

				// Calculate latency: Time elapsed since start. This measures total
				// time including network round-trip, processing, and any retries
				// within Probe.
				latency := time.Since(start)

				// Aggregate result with mutex protection: Lock before writing to
				// shared map, unlock after. This prevents data races.
				//
				// Critical section: mu.Lock() → write map → mu.Unlock()
				// Keep this short! Long critical sections hurt parallelism.
				//
				// Lock/Unlock ordering: ALWAYS pair Lock with Unlock. Forgetting
				// Unlock causes deadlock (other goroutines block forever waiting
				// for lock). Using defer mu.Unlock() is safer but we avoid it here
				// to minimize lock hold time.
				mu.Lock()
				if err != nil {
					// Failure: Record error with context. fmt.Errorf with %w wraps
					// the error, preserving the error chain for debugging.
					//
					// Error wrapping pattern: Same as previous modules. Adds context
					// ("probe failed") while preserving original error (network timeout,
					// connection refused, etc.).
					res.Failures[endpoint] = fmt.Errorf("probe failed: %w", err)
				} else {
					// Success: Record latency. This allows callers to identify slow
					// endpoints even if they're "healthy."
					//
					// Map insert: If endpoint already exists (duplicate), this
					// overwrites. In practice, endpoints should be unique.
					res.Successes[endpoint] = latency
				}
				mu.Unlock()
			}
			// Range loop exits when jobs channel closes (see Step 5).
			// wg.Done() called by defer, signaling this worker completed.
		}()
	}

	// ============================================================================
	// STEP 5: Send Jobs to Workers - Producer Pattern
	// ============================================================================
	// Job producer goroutine: Feeds endpoints into jobs channel. Runs
	// independently of workers, allowing them to start processing immediately.
	//
	// Why separate goroutine? If we sent jobs in the main goroutine before
	// calling wg.Wait(), we'd block if the channel buffer fills. By using a
	// separate goroutine, we can call wg.Wait() immediately while jobs are sent
	// in the background.
	//
	// Deadlock prevention: Main goroutine waits on workers (wg.Wait). Workers
	// wait on jobs (range jobs). Producer sends jobs and closes channel. This
	// chain ensures no circular waiting.
	go func() {
		// defer close(jobs): CRITICAL! Closing the channel signals workers
		// to exit their range loops. Without this, workers would block forever
		// waiting for more jobs, and wg.Wait() would never return.
		//
		// When to close: After all jobs are sent OR after timeout. The select
		// below handles both cases.
		//
		// Close semantics: Closing a channel doesn't block. It's instant.
		// After close, sends panic but receives return zero value. Range loops
		// exit cleanly when channel closes.
		defer close(jobs)

		// Send all endpoints: Iterate through cfg.Endpoints and send each to
		// the jobs channel.
		for _, ep := range cfg.Endpoints {
			// Select with timeout check: This is crucial for cancellation.
			// If ctx times out mid-iteration, we detect it and exit early.
			//
			// How select works:
			//   - case <-ctx.Done(): Chosen if context cancels/times out. Exit loop.
			//   - case jobs <- ep: Chosen if channel is ready to receive. Send job.
			//
			// Why check ctx.Done()? If overall timeout fires while we're sending
			// jobs, we should stop immediately rather than enqueuing more work.
			// Workers are already exiting, no point sending more jobs.
			//
			// Race condition: If timeout fires between iteration and select, the
			// select will catch it on the next iteration. Slight delay but harmless.
			select {
			case <-ctx.Done():
				// Context canceled or timed out. Stop sending jobs and exit.
				// defer close(jobs) will run, signaling workers to exit.
				return
			case jobs <- ep:
				// Job sent successfully. Continue to next endpoint.
				// If channel buffer is full, this blocks until a worker pulls a job.
			}
		}
		// All jobs sent. defer close(jobs) signals workers to finish and exit.
	}()

	// ============================================================================
	// STEP 6: Wait for Workers to Complete - Synchronization
	// ============================================================================
	// wg.Wait(): Blocks until all workers call wg.Done(). This ensures all
	// probes have completed (or timed out) before we return results.
	//
	// What happens during Wait:
	//   1. Workers are processing jobs from the channel
	//   2. Producer is sending jobs (or already finished and closed channel)
	//   3. Main goroutine blocks here
	//   4. Workers finish jobs, exit range loop, call wg.Done()
	//   5. When all workers Done, Wait returns
	//
	// Timeout interaction: If ctx times out:
	//   1. Producer's select detects ctx.Done(), exits, closes channel
	//   2. Workers' per-request contexts cancel, probes fail fast
	//   3. Workers finish their current job (with error), exit range loop
	//   4. Workers call wg.Done(), Wait returns
	//
	// Pattern: This is the standard synchronization point in worker pool pattern.
	// Main goroutine launches workers, sends jobs, then waits for completion.
	wg.Wait()

	// ============================================================================
	// STEP 7: Check for Timeout and Return Results
	// ============================================================================
	// Check context error: After workers complete, check if the context was
	// canceled or timed out. ctx.Err() returns nil if context is still valid,
	// or context.Canceled/context.DeadlineExceeded if canceled/timed out.
	//
	// Why check after Wait? We want to return partial results even on timeout.
	// If 100 endpoints were probed but only 50 completed before timeout, we
	// return those 50 results plus an error indicating timeout. This is more
	// useful than returning nothing.
	//
	// Error wrapping: fmt.Errorf with %w wraps ctx.Err(), which is either
	// context.Canceled (explicit cancellation) or context.DeadlineExceeded
	// (timeout). Caller can use errors.Is() to check which case occurred.
	//
	// Partial results on error: Note we return res even when error is non-nil.
	// This is intentional - caller gets whatever we completed before timeout.
	if ctx.Err() != nil {
		return res, fmt.Errorf("run aborted: %w", ctx.Err())
	}

	// Success: All probes completed within timeout. Return results and nil error.
	//
	// Result contents:
	//   - Successes: Map of endpoint → latency for successful probes
	//   - Failures: Map of endpoint → error for failed probes
	//
	// Building on previous concepts:
	//   - We validated all inputs (Step 1) → defensive programming from module 01
	//   - We used contexts for timeouts (Step 2) → context patterns from module 06
	//   - We handled errors consistently → error wrapping from all previous modules
	//   - We protected concurrent access (Step 4) → new concurrency pattern
	//   - We used worker pools (Steps 4-5) → new concurrency pattern
	//
	// How concepts build:
	//   - Module 01: Single RPC call with context
	//   - Module 06: Single transaction with context
	//   - Module 16: Multiple concurrent operations with contexts and worker pools
	//
	// This pattern repeats: Worker pools with channels and WaitGroup are the
	// idiomatic Go solution for bounded concurrency. You'll see this in indexers
	// (module 17), monitors, and any tool that needs to process multiple items
	// concurrently with resource limits.
	return res, nil
}
