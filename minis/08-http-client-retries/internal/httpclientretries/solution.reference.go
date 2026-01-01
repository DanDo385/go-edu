//go:build reference

package httpclientretries

/*
Problem: Build a resilient HTTP client with retries and exponential backoff

Requirements:
1. Retry failed requests automatically
2. Exponential backoff: delay increases exponentially
3. Jitter: add randomness to prevent thundering herd
4. Context-aware: respect timeouts and cancellation
5. Generic JSON decoding

Algorithm: Exponential Backoff with Jitter
- Attempt request
- If fails and retryable: wait backoff duration, retry
- If fails and non-retryable: return error immediately
- Backoff formula: BaseDelay * 2^attempt ± 20% jitter
- Repeat up to MaxRetries times

Why Exponential Backoff:
- Linear backoff (1s, 2s, 3s) is too aggressive
- Exponential (100ms, 200ms, 400ms, 800ms) backs off quickly
- Gives servers time to recover from overload
- Industry standard for distributed systems

Why Jitter (Randomness):
- Prevents thundering herd (synchronized retries)
- Spreads retries over time
- Reduces load spikes on recovering servers

Time Complexity: O(retries * request_time)
Space Complexity: O(1)
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Client wraps an HTTP client with retry configuration.
// BREAKPOINT: Set breakpoint in methods to trace retry logic
type Client struct {
	HTTP       *http.Client
	MaxRetries int
	BaseDelay  time.Duration
}

// GetJSON fetches JSON with retries.
//
// Algorithm:
// 1. Attempt request
// 2. If success, return immediately
// 3. If retryable error, calculate backoff delay
// 4. Wait with context awareness
// 5. Retry up to MaxRetries times
//
// BREAKPOINT: Set breakpoint at function entry to trace retry attempts
// DEBUG: Watch 'attempt' to see current retry count
// DEBUG: Watch 'lastErr' to see error progression
// DEBUG: Watch 'delay' to see backoff calculation
func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error) {
	var zero T
	var lastErr error

	// BREAKPOINT: Set breakpoint here to trace retry loop
	// DEBUG: Watch 'attempt' as it increments
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		// DEBUG: Current attempt number
		// BREAKPOINT: Set breakpoint here to trace each request attempt

		// Try request
		result, err := doRequest[T](ctx, c.HTTP, url)
		if err == nil {
			// DEBUG: Request succeeded on this attempt
			return result, nil
		}

		// DEBUG: Request failed, saving error
		lastErr = err

		// Check if retryable
		// BREAKPOINT: Set breakpoint here to see retry decision
		// DEBUG: Watch 'isRetryable(err)' to see if we should retry
		if !isRetryable(err) {
			// DEBUG: Error is non-retryable, returning immediately
			return zero, err
		}

		// Last attempt failed
		// BREAKPOINT: Set breakpoint here to check if this is the last attempt
		// DEBUG: Watch 'attempt' vs 'c.MaxRetries'
		if attempt == c.MaxRetries {
			// DEBUG: All retries exhausted
			break
		}

		// Calculate backoff with jitter
		// BREAKPOINT: Set breakpoint here to trace backoff calculation
		// DEBUG: Watch 'c.BaseDelay' to see base delay
		// DEBUG: Watch '1<<uint(attempt)' to see exponential multiplier
		delay := c.BaseDelay * time.Duration(1<<uint(attempt))
		// DEBUG: Delay before jitter

		// Add jitter (±20%)
		// BREAKPOINT: Set breakpoint here to see jitter calculation
		// DEBUG: Watch 'rand.Float64()' to see random value [0, 1)
		// DEBUG: Watch 'rand.Float64()*0.4-0.2' to see jitter factor [-0.2, 0.2)
		jitter := time.Duration(rand.Float64()*0.4-0.2) * delay // ±20%
		delay += jitter
		// DEBUG: Final delay with jitter applied

		// Wait with context awareness
		// BREAKPOINT: Set breakpoint here to trace waiting period
		// DEBUG: Watch 'delay' to see how long we'll wait
		select {
		case <-time.After(delay):
			// DEBUG: Delay completed, continuing to next attempt
			// Continue to next attempt
		case <-ctx.Done():
			// DEBUG: Context cancelled during backoff
			// BREAKPOINT: Set breakpoint here to trace cancellation
			return zero, ctx.Err()
		}
	}

	// DEBUG: All retries failed, returning last error
	// BREAKPOINT: Set breakpoint here to inspect final error
	return zero, fmt.Errorf("all retries failed: %w", lastErr)
}

// doRequest executes a single HTTP GET request and decodes JSON.
//
// Algorithm:
// 1. Create HTTP request with context
// 2. Execute request
// 3. Read response body
// 4. Decode JSON into generic type T
//
// BREAKPOINT: Set breakpoint at function entry to trace HTTP requests
// DEBUG: Watch 'url' to see request destination
// DEBUG: Watch 'resp.StatusCode' to see HTTP status
func doRequest[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var zero T

	// BREAKPOINT: Set breakpoint here to trace request creation
	// DEBUG: Watch 'url' to see target
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		// DEBUG: Request creation failed
		return zero, err
	}

	// BREAKPOINT: Set breakpoint here to trace request execution
	// DEBUG: Step over to see request complete
	resp, err := client.Do(req)
	if err != nil {
		// DEBUG: Request failed with network error
		// BREAKPOINT: Set breakpoint here to inspect network error
		return zero, err
	}
	defer resp.Body.Close()

	// BREAKPOINT: Set breakpoint here to check HTTP status
	// DEBUG: Watch 'resp.StatusCode' to see HTTP response code
	if resp.StatusCode != http.StatusOK {
		// DEBUG: Non-200 status code received
		// BREAKPOINT: Set breakpoint here to inspect error status
		return zero, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// BREAKPOINT: Set breakpoint here to trace body reading
	// DEBUG: Watch 'body' to see response content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// DEBUG: Failed to read response body
		return zero, err
	}

	// BREAKPOINT: Set breakpoint here to trace JSON decoding
	// DEBUG: Watch 'body' to see raw JSON
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		// DEBUG: JSON decode failed
		// BREAKPOINT: Set breakpoint here to inspect decode error
		return zero, err
	}

	// DEBUG: Watch 'result' to see decoded value
	return result, nil
}

// isRetryable determines if an error should be retried.
//
// Algorithm:
// - For this simple implementation, retry all errors
// - Production: check error types (network, timeout, 5xx) vs non-retryable (4xx, parse errors)
//
// BREAKPOINT: Set breakpoint here to trace retry decisions
// DEBUG: Watch 'err' to see error type
func isRetryable(err error) bool {
	// DEBUG: Simple heuristic - retry all errors
	// BREAKPOINT: Set breakpoint here to see decision
	// Production would check:
	// - Network errors (DNS, connection): RETRY
	// - Timeout errors: RETRY
	// - 5xx status codes: RETRY
	// - 429 Too Many Requests: RETRY
	// - 4xx status codes: DON'T RETRY
	// - JSON decode errors: DON'T RETRY
	return err != nil
}

/*
Alternatives & Trade-offs:

1. Fixed delay (no exponential backoff):
   Algorithm: Same delay every retry
   Pros: Simpler
   Cons: Too aggressive on failing servers

2. No jitter:
   Algorithm: delay = BaseDelay * 2^attempt
   Pros: Deterministic, simpler
   Cons: Thundering herd problem

3. Fibonacci backoff:
   Algorithm: delays follow Fibonacci sequence
   Pros: Grows slower than exponential
   Cons: Less standard, more complex

4. Circuit breaker pattern:
   Algorithm: Track failure rate, stop retries if rate is high
   Pros: Protects backend from overload
   Cons: More complex, needs state management

5. Request hedging:
   Algorithm: Send duplicate request after delay
   Pros: Lower tail latency
   Cons: Increased backend load
*/
