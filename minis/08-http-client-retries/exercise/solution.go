//go:build solution
// +build solution

/*
Problem: Build a resilient HTTP client with retries and exponential backoff

Requirements:
1. Retry failed requests automatically
2. Exponential backoff: delay increases exponentially
3. Jitter: add randomness to prevent thundering herd
4. Context-aware: respect timeouts and cancellation
5. Generic JSON decoding

Algorithm:
- Attempt request
- If fails and retryable: wait backoff duration, retry
- If fails and non-retryable: return error immediately
- Repeat up to MaxRetries times

Time Complexity: O(retries * request_time)
Space Complexity: O(1)

Why Go is well-suited:
- net/http: Production-grade HTTP client built-in
- context.Context: Standardized cancellation
- Generics: Type-safe JSON decoding

Compared to other languages:
- Python: requests library similar, but no built-in generics
- JavaScript: fetch() requires manual retry logic
- Rust: reqwest crate is similar, more complex types
*/

package exercise

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
type Client struct {
	HTTP       *http.Client
	MaxRetries int
	BaseDelay  time.Duration
}

// GetJSON fetches JSON with retries.
func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error) {
	var zero T
	var lastErr error

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		// Try request
		result, err := doRequest[T](ctx, c.HTTP, url)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Check if retryable
		if !isRetryable(err) {
			return zero, err
		}

		// Last attempt failed
		if attempt == c.MaxRetries {
			break
		}

		// Calculate backoff with jitter
		delay := c.BaseDelay * time.Duration(1<<uint(attempt))
		jitter := time.Duration(rand.Float64()*0.4-0.2) * delay // ±20%
		delay += jitter

		// Wait with context awareness
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}

	return zero, fmt.Errorf("all retries failed: %w", lastErr)
}

func doRequest[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var zero T

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return zero, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, err
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return zero, err
	}

	return result, nil
}

func isRetryable(err error) bool {
	// Simple heuristic: retry on network errors, not on parse errors
	// In production, check specific error types
	return err != nil
}

/*
Alternatives:

1. Circuit breaker: Stop retrying if failure rate is high
2. Request hedging: Send duplicate requests after delay
3. Adaptive backoff: Adjust based on server feedback

Go vs X:
- Python requests: Similar, but no generics
- JS fetch: Manual retry logic needed
- Rust reqwest: More complex types
*/
