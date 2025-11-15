package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

/*
Project 16: Context Cancellation and Timeouts - Comprehensive Demo

This program demonstrates:
1. context.WithCancel - Manual cancellation
2. context.WithTimeout - Automatic cancellation after duration
3. context.WithDeadline - Cancellation at specific time
4. Preventing goroutine leaks
5. Real-world HTTP timeout patterns

Each demo is standalone and heavily commented.
*/

func main() {
	fmt.Println("=== Context Cancellation and Timeouts Demo ===\n")

	// Demo 1: Manual Cancellation
	fmt.Println("--- Demo 1: WithCancel (Manual Cancellation) ---")
	demoWithCancel()
	time.Sleep(500 * time.Millisecond)

	// Demo 2: Automatic Timeout
	fmt.Println("\n--- Demo 2: WithTimeout (Automatic Cancellation) ---")
	demoWithTimeout()
	time.Sleep(500 * time.Millisecond)

	// Demo 3: Deadline-Based Cancellation
	fmt.Println("\n--- Demo 3: WithDeadline (Cancel at Specific Time) ---")
	demoWithDeadline()
	time.Sleep(500 * time.Millisecond)

	// Demo 4: Preventing Goroutine Leaks
	fmt.Println("\n--- Demo 4: Preventing Goroutine Leaks ---")
	demoGoroutineLeak()
	time.Sleep(500 * time.Millisecond)

	// Demo 5: HTTP Request Timeouts
	fmt.Println("\n--- Demo 5: HTTP Request Timeouts ---")
	demoHTTPTimeout()
	time.Sleep(500 * time.Millisecond)

	// Demo 6: Coordinating Multiple Workers
	fmt.Println("\n--- Demo 6: Coordinating Multiple Workers ---")
	demoWorkerCoordination()
	time.Sleep(500 * time.Millisecond)

	// Demo 7: Context Values (Use Sparingly!)
	fmt.Println("\n--- Demo 7: Context Values (Request-Scoped Data) ---")
	demoContextValues()

	fmt.Println("\n=== All Demos Complete ===")
}

// ============================================================================
// Demo 1: context.WithCancel - Manual Cancellation
// ============================================================================

func demoWithCancel() {
	// context.WithCancel creates a cancellable context
	// Use when: You need manual control over when to stop operations
	//
	// Returns:
	//   - ctx: New context that will be cancelled when cancel() is called
	//   - cancel: Function to call to trigger cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// CRITICAL: Always defer cancel() to release resources
	// Even if you don't explicitly cancel, this cleans up internal goroutines/timers
	defer cancel()

	// Start a worker that will run until context is cancelled
	done := make(chan bool)
	go worker(ctx, "Worker-1", done)

	// Let worker run for 2 seconds
	fmt.Println("Letting worker run for 2 seconds...")
	time.Sleep(2 * time.Second)

	// Manually cancel the context
	fmt.Println("Calling cancel() to stop worker...")
	cancel()

	// Wait for worker to finish
	<-done
	fmt.Println("Worker stopped successfully")
}

// worker simulates a long-running task that checks context cancellation
func worker(ctx context.Context, name string, done chan bool) {
	defer func() { done <- true }()

	// Infinite loop that periodically checks if context is cancelled
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			// ctx.Done() returns a channel that's closed when context is cancelled
			// When the channel is closed, this case becomes ready
			fmt.Printf("%s: Context cancelled, stopping (processed %d items)\n", name, i)
			return

		default:
			// Context not cancelled, continue working
			fmt.Printf("%s: Working... (iteration %d)\n", name, i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// ============================================================================
// Demo 2: context.WithTimeout - Automatic Timeout
// ============================================================================

func demoWithTimeout() {
	// context.WithTimeout creates a context that automatically cancels
	// after the specified duration
	//
	// Use when: You want operations to fail if they take too long
	// Example: HTTP requests, database queries, RPC calls
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Always defer cancel() to release timer resources

	fmt.Println("Starting operation with 2-second timeout...")

	// Simulate an operation that takes 3 seconds
	// This will timeout because it exceeds the 2-second deadline
	err := slowOperation(ctx, "Slow-Op-1", 3*time.Second)

	if err != nil {
		// Check what kind of error occurred
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("❌ Operation timed out (as expected)")
		} else {
			fmt.Printf("❌ Operation failed: %v\n", err)
		}
	} else {
		fmt.Println("✅ Operation completed")
	}

	// Now try an operation that finishes in time
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	fmt.Println("\nStarting operation that will finish in time...")
	err = slowOperation(ctx2, "Fast-Op-1", 1*time.Second)

	if err != nil {
		fmt.Printf("❌ Operation failed: %v\n", err)
	} else {
		fmt.Println("✅ Operation completed successfully")
	}
}

// slowOperation simulates a slow operation that respects context cancellation
func slowOperation(ctx context.Context, name string, duration time.Duration) error {
	fmt.Printf("%s: Starting (will take %v)...\n", name, duration)

	// Use select to make the sleep interruptible by context
	select {
	case <-time.After(duration):
		// Duration elapsed, operation completed
		fmt.Printf("%s: Completed\n", name)
		return nil

	case <-ctx.Done():
		// Context was cancelled before operation completed
		fmt.Printf("%s: Cancelled\n", name)
		return ctx.Err() // Returns context.DeadlineExceeded or context.Canceled
	}
}

// ============================================================================
// Demo 3: context.WithDeadline - Cancel at Specific Time
// ============================================================================

func demoWithDeadline() {
	// context.WithDeadline creates a context that cancels at a specific time
	//
	// Use when: You need work to finish by a specific point in time
	// Example: Rate limit windows, SLA enforcement, scheduled tasks
	//
	// Note: context.WithTimeout(parent, duration) is equivalent to
	//       context.WithDeadline(parent, time.Now().Add(duration))

	// Set deadline to 2 seconds from now
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("Current time: %s\n", time.Now().Format("15:04:05"))
	fmt.Printf("Deadline set to: %s (2 seconds from now)\n", deadline.Format("15:04:05"))

	// Check if context has a deadline
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("Context deadline: %s\n", d.Format("15:04:05"))
		fmt.Printf("Time remaining: %v\n", time.Until(d))
	}

	// Start operation that will be cancelled by deadline
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Working... (time remaining: %v)\n", time.Until(deadline))

		case <-ctx.Done():
			fmt.Printf("❌ Deadline exceeded at %s\n", time.Now().Format("15:04:05"))
			fmt.Printf("Context error: %v\n", ctx.Err())
			return
		}
	}
}

// ============================================================================
// Demo 4: Preventing Goroutine Leaks
// ============================================================================

func demoGoroutineLeak() {
	fmt.Println("Demonstrating goroutine leak prevention...\n")

	// BAD EXAMPLE: Goroutine leak
	fmt.Println("❌ Bad example (would leak goroutine):")
	badSearch("leak example")

	time.Sleep(500 * time.Millisecond)

	// GOOD EXAMPLE: No leak
	fmt.Println("\n✅ Good example (no leak):")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result := goodSearch(ctx, "no leak example")
	fmt.Printf("Result: %s\n", result)

	time.Sleep(500 * time.Millisecond)
}

// badSearch demonstrates a GOROUTINE LEAK
// The goroutine continues running even after timeout
func badSearch(query string) string {
	resultCh := make(chan string) // Unbuffered channel

	go func() {
		// Simulate slow search (10 seconds)
		time.Sleep(10 * time.Second)
		result := fmt.Sprintf("Results for: %s", query)

		// This send will BLOCK FOREVER if nobody is receiving
		// The goroutine will leak!
		resultCh <- result
	}()

	// Timeout after 1 second
	select {
	case result := <-resultCh:
		return result
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out (but goroutine still running! LEAK!)")
		return "timeout"
	}
}

// goodSearch demonstrates PROPER cleanup with context
// The goroutine stops when context is cancelled
func goodSearch(ctx context.Context, query string) string {
	// Buffered channel (size 1) prevents sender from blocking
	// if receiver has already returned
	resultCh := make(chan string, 1)

	go func() {
		// Simulate slow search that respects context
		select {
		case <-time.After(10 * time.Second):
			// Search completed
			select {
			case resultCh <- fmt.Sprintf("Results for: %s", query):
				// Sent successfully
			case <-ctx.Done():
				// Context cancelled while trying to send, don't block
				fmt.Println("Goroutine: Context cancelled while sending, stopping")
			}

		case <-ctx.Done():
			// Context cancelled during search
			fmt.Println("Goroutine: Context cancelled during search, stopping")
			return
		}
	}()

	// Wait for result or context cancellation
	select {
	case result := <-resultCh:
		return result
	case <-ctx.Done():
		fmt.Printf("Main: Context cancelled: %v\n", ctx.Err())
		return "cancelled"
	}
}

// ============================================================================
// Demo 5: HTTP Request Timeouts
// ============================================================================

func demoHTTPTimeout() {
	// Create a test HTTP server that has configurable delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow server (3 seconds)
		delay := 3 * time.Second
		fmt.Printf("Server: Received request, will respond in %v\n", delay)

		select {
		case <-time.After(delay):
			fmt.Fprintln(w, "Hello from server!")
		case <-r.Context().Done():
			// Client disconnected or timed out
			fmt.Println("Server: Client disconnected, aborting")
			return
		}
	}))
	defer server.Close()

	// Example 1: Request with timeout (will timeout)
	fmt.Println("Example 1: Request with 1-second timeout (will timeout)...")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()

	body, err := fetchURL(ctx1, server.URL)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("❌ Request timed out (as expected)")
		} else {
			fmt.Printf("❌ Request failed: %v\n", err)
		}
	} else {
		fmt.Printf("✅ Response: %s\n", body)
	}

	time.Sleep(500 * time.Millisecond)

	// Example 2: Request with longer timeout (will succeed)
	fmt.Println("\nExample 2: Request with 5-second timeout (will succeed)...")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	body, err = fetchURL(ctx2, server.URL)
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
	} else {
		fmt.Printf("✅ Response: %s\n", body)
	}
}

// fetchURL fetches a URL with context support
// The HTTP request will be cancelled if context is cancelled
func fetchURL(ctx context.Context, url string) (string, error) {
	// Create HTTP request with context
	// This is critical: it ties the request to the context's lifecycle
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	fmt.Println("Client: Sending request...")

	// Execute request
	// If context is cancelled (timeout, deadline, or manual cancel),
	// the HTTP client will abort the request immediately
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Client: Received response")
	return string(body), nil
}

// ============================================================================
// Demo 6: Coordinating Multiple Workers
// ============================================================================

func demoWorkerCoordination() {
	// Create cancellable context for all workers
	ctx, cancel := context.WithCancel(context.Background())

	numWorkers := 3
	var wg sync.WaitGroup

	fmt.Printf("Starting %d workers...\n", numWorkers)

	// Start multiple workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			coordinatedWorker(ctx, id)
		}(i)
	}

	// Let workers run for 2 seconds
	time.Sleep(2 * time.Second)

	// Cancel all workers at once
	fmt.Println("\nCancelling all workers...")
	cancel()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers stopped")
}

// coordinatedWorker demonstrates a worker that stops when context is cancelled
func coordinatedWorker(ctx context.Context, id int) {
	fmt.Printf("Worker %d: Started\n", id)

	// Worker loop with context cancellation check
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Do work
			fmt.Printf("Worker %d: Processing...\n", id)

		case <-ctx.Done():
			// Context cancelled, clean up and exit
			fmt.Printf("Worker %d: Context cancelled, shutting down\n", id)
			return
		}
	}
}

// ============================================================================
// Demo 7: Context Values (Use Sparingly!)
// ============================================================================

// contextKey is a custom type for context keys
// This prevents collisions with other packages using context values
type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

func demoContextValues() {
	// context.WithValue stores request-scoped data
	//
	// Use ONLY for:
	// - Request IDs (for tracing/logging)
	// - User identity (authentication)
	// - Request-scoped loggers
	//
	// DO NOT use for:
	// - Function parameters (pass explicitly)
	// - Configuration (use structs)
	// - Large objects (pass by reference)

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, requestIDKey, "req-12345")
	ctx = context.WithValue(ctx, userIDKey, "user-67890")

	fmt.Println("Processing request with context values...")

	// Pass context through call chain
	handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
	// Extract values from context
	requestID := getRequestID(ctx)
	userID := getUserID(ctx)

	fmt.Printf("Request ID: %s\n", requestID)
	fmt.Printf("User ID: %s\n", userID)

	// Call downstream function with same context
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	// Values are available throughout the call chain
	requestID := getRequestID(ctx)
	userID := getUserID(ctx)

	fmt.Printf("Processing for user %s (request %s)\n", userID, requestID)
}

// Helper functions for type-safe context value access

func getRequestID(ctx context.Context) string {
	// Type assertion with fallback
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return "unknown"
	}
	return id
}

func getUserID(ctx context.Context) string {
	id, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "unknown"
	}
	return id
}

/*
Key Takeaways from This Demo:

1. **Always defer cancel()**
   - Even if you don't explicitly cancel, this cleans up resources
   - Prevents goroutine leaks from timers

2. **Check ctx.Done() in loops**
   - Long-running operations should periodically check for cancellation
   - Use select statement to check context and do work

3. **Use buffered channels for goroutine results**
   - Prevents sender from blocking if receiver has already returned
   - Size 1 is sufficient for single result

4. **http.NewRequestWithContext is critical**
   - Ties HTTP request to context lifecycle
   - Request is cancelled when context is cancelled

5. **Context values are for request-scoped metadata only**
   - Request IDs, user identity, loggers
   - NOT for function parameters or configuration

6. **Context propagates down the call stack**
   - Pass ctx as first parameter to all functions
   - Child contexts inherit parent's cancellation

7. **Never store context in a struct**
   - Context is request-scoped, not object-scoped
   - Always pass as function parameter

Common Patterns:

✅ Good:
  - ctx, cancel := context.WithTimeout(parent, 5*time.Second)
  - defer cancel()
  - select { case <-ctx.Done(): return ctx.Err() }

❌ Bad:
  - Not calling cancel() (resource leak)
  - Ignoring ctx.Done() in long operations (hangs)
  - Passing nil context (panic)
  - Storing context in struct (wrong lifetime)
*/
