//go:build !solution && !reference

package httpclientretries



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
	// TODO: Implement this function
	panic("unimplemented")
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
	// TODO: Implement this function
	panic("unimplemented")
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
	// TODO: Implement this function
	panic("unimplemented")
}


