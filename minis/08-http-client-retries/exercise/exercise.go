//go:build !solution
// +build !solution

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
	HTTP       *http.Client  // Underlying HTTP client
	MaxRetries int           // Maximum number of retry attempts
	BaseDelay  time.Duration // Base delay for exponential backoff (e.g., 100ms)
}

// GetJSON fetches JSON from url and decodes it into type T.
// Retries on failure with exponential backoff and jitter.
//
// Backoff formula: delay = BaseDelay * (2^attempt) * (1 ± 20% jitter)
//
// Parameters:
//   - ctx: Context for timeout/cancellation
//   - c: Client with retry configuration
//   - url: URL to fetch
//
// Returns:
//   - T: Decoded JSON response
//   - error: Non-nil if all retries fail
func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error) {
	var zero T
	var lastErr error

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		val, err := doRequest[T](ctx, c.HTTP, url)
		if err == nil {
			return val, nil
		}
		lastErr = err

		if !isRetryable(err) || attempt == c.MaxRetries {
			break
		}

		delay := c.BaseDelay * time.Duration(1<<uint(attempt))
		jitter := time.Duration(rand.Float64()*0.4-0.2) * delay
		delay += jitter

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}

	return zero, fmt.Errorf("all retries failed: %w", lastErr)
}

func doRequest[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var zero T
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var out T
	if err := json.Unmarshal(body, &out); err != nil {
		return zero, err
	}
	return out, nil
}

func isRetryable(err error) bool {
	return err != nil
}
