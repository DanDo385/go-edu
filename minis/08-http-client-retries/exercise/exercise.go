package exercise

import (
	"context"
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
	// TODO: implement
	var zero T
	return zero, nil
}
