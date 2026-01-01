//go:build !solution && !reference

package httpclientretries

/*
Problem: Build a resilient HTTP client with retries and exponential backoff
Requirements:
1. Retry failed requests automatically
2. Exponential backoff: delay increases exponentially
3. Jitter: add randomness to prevent thundering herd
Algorithm: Exponential Backoff with Jitter
- Attempt request
- If fails and retryable: wait backoff duration, retry
- If fails and non-retryable: return error immediately
- Backoff formula: BaseDelay * 2^attempt ± 20% jitter
- Repeat up to MaxRetries times
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

type Client struct {
	HTTP       *http.Client
	MaxRetries int
	BaseDelay  time.Duration
}

// GetJSON - TODO: implement this function
func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// doRequest - TODO: implement this function
func doRequest[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// isRetryable - TODO: implement this function
func isRetryable(err error) bool {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return false
}

