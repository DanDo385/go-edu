package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/16-concurrency/internal/concurrency"
)

/*
===================================================================
Concurrent RPC Endpoint Health Checking
===================================================================

This program demonstrates how to probe multiple Ethereum RPC endpoints
concurrently using a worker pool pattern. This is essential for building
production-grade Ethereum infrastructure that needs to monitor node
health, implement failover, or distribute load across multiple endpoints.

USAGE:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go <demo>

Arguments:
    [demo]          - Demo to run: all, basic, timeout, performance (default: "all")

Examples:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go basic
    go run ./cmd/app/main.go timeout
    go run ./cmd/app/main.go performance

Common Use Cases:
    - Health Monitoring: Check if RPC endpoints are responsive
    - Load Balancing: Measure latency to select fastest endpoint
    - Failover: Detect unavailable endpoints and switch to backups
    - Network Discovery: Test connectivity to multiple networks

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter concurrency package functions
- Watch Variables panel to see worker pool execution
- Inspect Successes and Failures maps in results

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the concurrency package from internal/concurrency
- exercise.go, solution.reference.go are in package concurrency
- The concurrency.Run() function implements the worker pool pattern

KEY CONCEPTS:
- Worker Pool: Bounded concurrency prevents resource exhaustion
- Goroutines: Lightweight threads managed by Go runtime
- Channels: Communication mechanism between goroutines
- Mutex: Protects shared data from concurrent access (data races)
- WaitGroup: Synchronizes goroutine completion
- Context: Propagates cancellation and timeouts across goroutines
- Prober Interface: Dependency injection for testability

COMPUTER SCIENCE CONCEPTS:
- Concurrency vs Parallelism: Concurrent = multiple tasks in progress; Parallel = simultaneous execution
- Producer-Consumer Pattern: Jobs channel decouples job generation from execution
- Work-Stealing Queue: Workers pull jobs dynamically, balancing load automatically
- Critical Section: Code that accesses shared resources (must be protected by mutex)
- Race Condition: Bug where outcome depends on timing of concurrent operations
- Deadlock: Circular wait where goroutines block each other forever
*/

// SimpleProber implements the Prober interface for basic RPC health checks.
// This demonstrates dependency injection—the worker pool doesn't care HOW
// endpoints are probed, it just calls Probe(). This makes testing easy
// (mock Prober) and the code reusable (different Prober implementations).
type SimpleProber struct{}

func (p *SimpleProber) Probe(ctx context.Context, endpoint string) error {
	// BREAKPOINT: Set breakpoint here to see each probe execution
	// DEBUG: Each worker calls this method for every endpoint it processes
	// DEBUG: ctx contains timeout—if exceeded, this should return quickly

	// Create a child context with timeout to prevent hanging indefinitely
	// on unresponsive endpoints. This is defensive programming—never trust
	// external services to respond promptly.
	dialCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Attempt to connect to the RPC endpoint
	// BREAKPOINT: Set breakpoint here
	// DEBUG: Step Over (F10) to see if connection succeeds or fails
	// DEBUG: Watch 'err' variable to see specific error (timeout, refused, etc.)
	client, err := ethclient.Dial(dialCtx, endpoint)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer client.Close()

	// Verify the endpoint works by fetching ChainID
	// This ensures the endpoint is not just accepting connections but
	// actually responding to RPC requests.
	chainCtx, chainCancel := context.WithTimeout(ctx, 2*time.Second)
	defer chainCancel()

	_, err = client.ChainID(chainCtx)
	if err != nil {
		return fmt.Errorf("chain ID call failed: %w", err)
	}

	// Success: endpoint is responsive
	return nil
}

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the demo type (optional)

	demo := "all"
	if len(os.Args) > 1 {
		demo = os.Args[1]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'demo' variable in Variables panel
	fmt.Printf("=== Concurrent RPC Endpoint Health Checking ===\n")
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Example 1 - Basic Worker Pool with Multiple Endpoints
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through basic example
	// DEBUG: This demonstrates the fundamental worker pool pattern
	// DEBUG: Watch how multiple endpoints are probed concurrently

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Worker Pool ===")
		fmt.Println("Probing multiple RPC endpoints concurrently with 4 workers...")
		fmt.Println()

		// Create context for the entire operation
		// BREAKPOINT: Set breakpoint here
		ctx := context.Background()

		// Define endpoints to probe
		// These are well-known public RPC endpoints for different networks
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice we're testing multiple networks: Mainnet, Sepolia, Polygon
		endpoints := []string{
			"https://eth.llamarpc.com",                // Ethereum Mainnet
			"https://rpc.ankr.com/eth",                // Ethereum Mainnet (alt)
			"https://rpc.sepolia.org",                 // Sepolia Testnet
			"https://polygon-rpc.com",                 // Polygon Mainnet
			"https://rpc.ankr.com/polygon",            // Polygon Mainnet (alt)
			"https://invalid-endpoint-test.local",     // Invalid (will fail)
		}

		// Configure worker pool
		// BREAKPOINT: Set breakpoint here
		// DEBUG: 4 workers will process 6 endpoints concurrently
		// DEBUG: 10-second timeout for entire operation
		cfg := concurrency.Config{
			Endpoints: endpoints,
			Workers:   4,              // 4 concurrent workers
			Timeout:   10 * time.Second, // 10s total timeout
		}

		// Create prober
		prober := &SimpleProber{}

		// BREAKPOINT: Set breakpoint here before calling concurrency.Run
		// DEBUG: Step Into (F11) to see the worker pool implementation
		// DEBUG: You'll see: input validation, worker pool setup, job distribution
		fmt.Println("Starting worker pool...")
		start := time.Now()

		result, err := concurrency.Run(ctx, prober, cfg)

		elapsed := time.Since(start)

		// BREAKPOINT: Set breakpoint here after worker pool completes
		// DEBUG: Expand 'result' in Variables panel
		// DEBUG: Look at Successes map (endpoint → latency)
		// DEBUG: Look at Failures map (endpoint → error)

		if err != nil {
			fmt.Printf("⚠ Worker pool completed with error: %v\n", err)
		} else {
			fmt.Println("✓ Worker pool completed successfully")
		}
		fmt.Printf("Total time: %v\n\n", elapsed)

		// Display results
		fmt.Printf("Successes (%d):\n", len(result.Successes))
		for endpoint, latency := range result.Successes {
			fmt.Printf("  ✓ %s - %v\n", endpoint, latency)
		}
		fmt.Println()

		fmt.Printf("Failures (%d):\n", len(result.Failures))
		for endpoint, err := range result.Failures {
			fmt.Printf("  ✗ %s - %v\n", endpoint, err)
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice how all endpoints were probed concurrently
		// DEBUG: Total time is much less than sum of individual probe times
		// DEBUG: This demonstrates the power of concurrency for I/O-bound operations

		// Educational explanation
		fmt.Println("Worker Pool Benefits:")
		fmt.Println("  1. Bounded concurrency: Only 4 goroutines, not 6 (prevents resource exhaustion)")
		fmt.Println("  2. Load balancing: Faster workers process more endpoints automatically")
		fmt.Println("  3. Timeout control: Entire operation bounded by 10-second timeout")
		fmt.Println("  4. Partial results: On timeout, returns endpoints completed so far")
		fmt.Println()
	}

	// ============================================================================
	// STEP 3: Example 2 - Timeout Behavior
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through timeout example
	// DEBUG: This shows how contexts propagate timeouts through worker pools
	// DEBUG: Useful for understanding timeout vs completion behavior

	if demo == "all" || demo == "timeout" {
		fmt.Println("=== Example 2: Timeout Behavior ===")
		fmt.Println("Testing worker pool timeout with aggressive deadline...")
		fmt.Println()

		ctx := context.Background()

		// Use many endpoints with very short timeout
		// BREAKPOINT: Set breakpoint here
		// DEBUG: 1-second timeout is too short for 10 endpoints
		// DEBUG: Watch how worker pool handles timeout gracefully
		endpoints := []string{
			"https://eth.llamarpc.com",
			"https://rpc.ankr.com/eth",
			"https://eth-mainnet.public.blastapi.io",
			"https://rpc.sepolia.org",
			"https://sepolia.gateway.tenderly.co",
			"https://polygon-rpc.com",
			"https://rpc.ankr.com/polygon",
			"https://optimism-mainnet.public.blastapi.io",
			"https://arbitrum-one.public.blastapi.io",
			"https://base-mainnet.public.blastapi.io",
		}

		cfg := concurrency.Config{
			Endpoints: endpoints,
			Workers:   4,
			Timeout:   1 * time.Second, // Aggressive timeout
		}

		prober := &SimpleProber{}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Watch how the worker pool races against timeout
		fmt.Println("Starting worker pool with 1-second timeout...")
		start := time.Now()

		result, err := concurrency.Run(ctx, prober, cfg)

		elapsed := time.Since(start)

		// BREAKPOINT: Set breakpoint here
		// DEBUG: err will likely be non-nil (timeout exceeded)
		// DEBUG: result still contains partial results (endpoints completed before timeout)

		if err != nil {
			fmt.Printf("⚠ Expected timeout occurred: %v\n", err)
		} else {
			fmt.Println("✓ Completed within timeout (unexpected)")
		}
		fmt.Printf("Actual time: %v (target: 1s)\n\n", elapsed)

		fmt.Printf("Completed before timeout: %d/%d endpoints\n",
			len(result.Successes)+len(result.Failures), len(endpoints))
		fmt.Printf("  Successes: %d\n", len(result.Successes))
		fmt.Printf("  Failures: %d\n", len(result.Failures))
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice we got partial results despite timeout
		// DEBUG: This is valuable—better than no results at all

		fmt.Println("Timeout Handling Insights:")
		fmt.Println("  1. Context timeout propagates to all workers immediately")
		fmt.Println("  2. Workers stop processing new jobs when context cancels")
		fmt.Println("  3. In-flight probes are canceled via their own contexts")
		fmt.Println("  4. Partial results are returned (completed before timeout)")
		fmt.Println("  5. Error indicates timeout occurred but doesn't lose data")
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 3 - Performance Analysis
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through performance example
	// DEBUG: This compares different worker pool sizes to show trade-offs

	if demo == "all" || demo == "performance" {
		fmt.Println("=== Example 3: Performance Analysis ===")
		fmt.Println("Comparing different worker pool sizes...")
		fmt.Println()

		ctx := context.Background()

		endpoints := []string{
			"https://eth.llamarpc.com",
			"https://rpc.ankr.com/eth",
			"https://rpc.sepolia.org",
			"https://polygon-rpc.com",
			"https://rpc.ankr.com/polygon",
			"https://optimism-mainnet.public.blastapi.io",
		}

		// Test different worker pool sizes
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Watch how different worker counts affect performance
		workerCounts := []int{1, 2, 4, 8}

		for _, workers := range workerCounts {
			cfg := concurrency.Config{
				Endpoints: endpoints,
				Workers:   workers,
				Timeout:   15 * time.Second,
			}

			prober := &SimpleProber{}

			// BREAKPOINT: Set breakpoint here in loop
			// DEBUG: Step through each iteration to compare timings
			fmt.Printf("Testing with %d worker(s)...\n", workers)
			start := time.Now()

			result, err := concurrency.Run(ctx, prober, cfg)

			elapsed := time.Since(start)

			if err != nil {
				fmt.Printf("  ⚠ Error: %v\n", err)
			} else {
				fmt.Printf("  ✓ Success\n")
			}

			fmt.Printf("  Time: %v\n", elapsed)
			fmt.Printf("  Completed: %d/%d endpoints\n",
				len(result.Successes)+len(result.Failures), len(endpoints))
			fmt.Printf("  Throughput: %.2f endpoints/sec\n\n",
				float64(len(result.Successes)+len(result.Failures))/elapsed.Seconds())
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Compare the timing results across different worker counts

		fmt.Println("Performance Insights:")
		fmt.Println("  1. More workers = faster completion (up to a point)")
		fmt.Println("  2. Diminishing returns after workers >= CPU cores")
		fmt.Println("  3. For I/O-bound tasks, workers > cores is often beneficial")
		fmt.Println("  4. Too many workers = overhead from goroutine scheduling")
		fmt.Println("  5. Optimal worker count depends on task type and system")
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Understanding Concurrency Patterns
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand concurrency patterns
	// DEBUG: This section explains the patterns used in the worker pool

	if demo == "all" {
		fmt.Println("=== Understanding Concurrency Patterns ===")
		fmt.Println()

		fmt.Println("Worker Pool Pattern Components:")
		fmt.Println("  1. Jobs Channel: Queue of work items (endpoints)")
		fmt.Println("  2. Workers: Goroutines that pull from jobs channel")
		fmt.Println("  3. Results Aggregation: Mutex-protected map for thread-safe writes")
		fmt.Println("  4. WaitGroup: Synchronization primitive to wait for workers")
		fmt.Println("  5. Context: Timeout and cancellation propagation")
		fmt.Println()

		fmt.Println("Why Use Worker Pools?")
		fmt.Println("  ✓ Bounded concurrency: Limit goroutines to prevent resource exhaustion")
		fmt.Println("  ✓ Resource management: Respect system limits (memory, file descriptors)")
		fmt.Println("  ✓ Rate limiting: Control load on external services")
		fmt.Println("  ✓ Predictable performance: Consistent resource usage")
		fmt.Println("  ✓ Graceful degradation: Handle timeouts and errors cleanly")
		fmt.Println()

		fmt.Println("Alternative Approaches (and why worker pools are better):")
		fmt.Println("  Sequential processing: Simple but slow (no concurrency)")
		fmt.Println("  Unbounded goroutines: Fast but dangerous (can exhaust memory)")
		fmt.Println("  Semaphore pattern: Similar to worker pool but more complex")
		fmt.Println()

		fmt.Println("Real-World Applications:")
		fmt.Println("  - RPC endpoint health monitoring (this example)")
		fmt.Println("  - Block indexing with multiple RPC calls per block")
		fmt.Println("  - Transaction broadcast to multiple nodes")
		fmt.Println("  - Parallel data fetching from multiple APIs")
		fmt.Println("  - Batch processing with controlled parallelism")
		fmt.Println()
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter concurrency package functions
	// 5. Watch Variables panel to see worker pool state
	// 6. Inspect result maps to see successes and failures
	// 7. Use Call Stack panel to see goroutine execution
	//
	// CONCURRENCY DEBUGGING TIPS:
	// - Use -race flag when running tests to detect data races
	// - Add logging to see goroutine IDs: fmt.Printf("goroutine %d\n", goid())
	// - Use breakpoints in worker loop to see individual job processing
	// - Watch channel state: len(jobs), cap(jobs)
	// - Monitor WaitGroup counter (not directly visible, but affects Wait())
	//
	// COMMON CONCURRENCY BUGS TO WATCH FOR:
	// - Data races: Multiple goroutines accessing same variable without mutex
	// - Deadlock: Circular wait (all goroutines blocked)
	// - Goroutine leaks: Goroutines never exit (check with runtime.NumGoroutine())
	// - Channel deadlock: Send/receive on closed or blocking channel
	// - Context leaks: Forgetting to call cancel() on contexts
	//
	// STEPPING INTO CONCURRENCY PACKAGE:
	// - Set breakpoint at concurrency.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to solution.reference.go
	// - Step through: input validation → worker pool setup → job distribution
	// - Watch goroutines spawn in Call Stack panel
	// - See mutex lock/unlock protecting shared map
	//
	// MONITORING WORKER POOL EXECUTION:
	// - Set breakpoint inside worker loop (solution.reference.go:220)
	// - Hit F5 to continue—breakpoint will hit multiple times (one per worker)
	// - Each hit is a different goroutine processing a job
	// - Use "Goroutines" panel (if available) to see all active workers
	//
	// ETHEREUM RPC CONCEPTS TO EXPLORE:
	// - Why probe multiple endpoints? Redundancy, failover, load balancing
	// - What makes a good health check? Fast, reliable, representative
	// - How often to probe? Trade-off: freshness vs load
	// - What to do with results? Route traffic, alert on failures, metrics
}
