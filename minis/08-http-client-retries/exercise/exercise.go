//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// import (
//     "context"
//     "encoding/json"
//     "fmt"
//     "io"
//     "math/rand"
//     "net/http"
//     "time"
// )

// ============================================================================
// HTTP CLIENT WITH RETRIES: Exponential Backoff Algorithm
// ============================================================================
//
// Exponential Backoff Algorithm:
// - Delay grows exponentially: BaseDelay * 2^attempt
// - Example: 100ms, 200ms, 400ms, 800ms
// - Prevents overwhelming failing servers
// - Industry standard for distributed systems
//
// Jitter (Randomness):
// - Add ±20% variation to delay
// - Prevents thundering herd (synchronized retries)
// - Spreads load over time
//
// Retry Strategy:
// - Network errors: RETRY
// - Timeout errors: RETRY
// - 5xx server errors: RETRY
// - 4xx client errors: DON'T RETRY
// - JSON parse errors: DON'T RETRY
//
// Time Complexity: O(retries * request_time)
// Space Complexity: O(1)
//
// ============================================================================

// TODO: Define Client struct
// Fields: HTTP (*http.Client), MaxRetries (int), BaseDelay (time.Duration)

// TODO: Implement GetJSON[T any](ctx context.Context, c *Client, url string) (T, error)
// Algorithm:
// 1. Loop from attempt = 0 to MaxRetries
// 2. Call doRequest[T] to fetch and decode JSON
// 3. If successful, return result
// 4. Check if error is retryable
// 5. If non-retryable or last attempt, return error
// 6. Calculate backoff: BaseDelay * (2^attempt)
// 7. Add jitter: ±20% random variation
// 8. Wait with select on time.After and ctx.Done()
// 9. If all retries fail, return last error

// TODO: Implement doRequest[T any](ctx context.Context, client *http.Client, url string) (T, error)
// Algorithm:
// 1. Create request with http.NewRequestWithContext
// 2. Execute request with client.Do
// 3. Check status code (200 OK)
// 4. Read body with io.ReadAll
// 5. Unmarshal JSON with json.Unmarshal
// 6. Return decoded value

// TODO: Implement isRetryable(err error) bool
// Simple version: return err != nil
// Production: check error types for network/timeout/5xx vs 4xx/parse errors
