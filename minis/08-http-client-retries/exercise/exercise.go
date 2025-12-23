//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "context" for cancellation and timeouts
// - "encoding/json" for JSON decoding
// - "fmt" for error formatting
// - "io" for reading response bodies
// - "math/rand" for jitter calculation
// - "net/http" for HTTP client
// - "time" for delays and backoff
//
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
// HTTP CLIENT WITH RETRIES: Resilient Network Requests
// ============================================================================
//
// Retry logic is essential for resilient distributed systems:
// - Network failures are common and transient
// - Servers may be temporarily overloaded (503 Service Unavailable)
// - Automatic retries improve reliability without user intervention
//
// Why exponential backoff:
// - Linear backoff (1s, 2s, 3s): Too aggressive, hammers failing servers
// - Exponential backoff (100ms, 200ms, 400ms, 800ms): Backs off quickly
// - Gives servers time to recover from overload
// - Prevents thundering herd (many clients retrying simultaneously)
//
// Why jitter (randomness):
// - Without jitter: All clients retry at same time (synchronized)
// - With jitter: Spreads retries over time (desynchronized)
// - Reduces load spikes on recovering servers
// - Formula: delay ± 20% random variation
//
// Retry strategies:
// - Immediate retry: No delay (useful for quick transients)
// - Fixed delay: Same delay between retries (simple but suboptimal)
// - Exponential backoff: Delay grows exponentially (industry standard)
// - Exponential backoff with jitter: Best practice for most systems
//
// When NOT to retry:
// - Client errors (4xx): Bad request, unauthorized, not found
// - Permanent failures: These won't succeed on retry
// - Non-idempotent operations: POST may duplicate (use PUT/DELETE instead)
//
// Memory considerations:
// - Each retry attempt allocates a new HTTP request/response
// - Old response bodies must be closed to avoid leaks
// - Context timeout should exceed total retry time
//
// ============================================================================

// ============================================================================
// Exercise 1: Define Client Type
// ============================================================================

// Client wraps an HTTP client with retry configuration.
//
// TODO: Define the Client struct with:
// 1. HTTP: *http.Client (underlying HTTP client)
// 2. MaxRetries: int (maximum number of retry attempts)
// 3. BaseDelay: time.Duration (base delay for exponential backoff)
//
// Example configuration:
// - MaxRetries: 3 (total 4 attempts: initial + 3 retries)
// - BaseDelay: 100ms (delays: 100ms, 200ms, 400ms)
//
// Why separate HTTP client:
// - Allows custom transport configuration (proxies, TLS)
// - Connection pooling (reuse TCP connections)
// - Timeout configuration (separate from retry logic)
//
// type Client struct {
//     HTTP       *http.Client  // Underlying HTTP client
//     MaxRetries int           // Maximum number of retry attempts
//     BaseDelay  time.Duration // Base delay for exponential backoff (e.g., 100ms)
// }

// ============================================================================
// Exercise 2: Implement GetJSON
// ============================================================================

// GetJSON fetches JSON from url and decodes it into type T.
//
// TODO: Implement retry logic with exponential backoff and jitter:
// 1. Loop from attempt = 0 to MaxRetries (inclusive)
// 2. Call doRequest[T] to fetch and decode JSON
// 3. If successful (no error), return result immediately
// 4. Save error for potential return
// 5. Check if error is retryable (use isRetryable helper)
// 6. If non-retryable or last attempt, break loop
// 7. Calculate backoff delay:
//    - Base formula: BaseDelay * (2^attempt)
//    - Add jitter: ± 20% random variation
// 8. Wait for delay with context awareness:
//    - Use time.After for delay
//    - Use select to handle context cancellation
// 9. If all retries exhausted, return last error
//
// Parameters:
//   - ctx: Context for timeout/cancellation
//   - c: Client with retry configuration
//   - url: URL to fetch
//
// Returns:
//   - T: Decoded JSON response
//   - error: Non-nil if all retries fail
//
// Exponential backoff calculation:
// - Attempt 0: BaseDelay * 2^0 = BaseDelay * 1 = 100ms
// - Attempt 1: BaseDelay * 2^1 = BaseDelay * 2 = 200ms
// - Attempt 2: BaseDelay * 2^2 = BaseDelay * 4 = 400ms
// - Attempt 3: BaseDelay * 2^3 = BaseDelay * 8 = 800ms
//
// Jitter calculation:
// - rand.Float64() returns [0.0, 1.0)
// - rand.Float64() * 0.4 - 0.2 returns [-0.2, 0.2) (±20%)
// - jitter = delay * [-0.2, 0.2)
// - finalDelay = delay + jitter
//
// Example with BaseDelay=100ms, attempt=2:
// - delay = 100ms * 4 = 400ms
// - jitter = 400ms * (rand between -0.2 and 0.2) = -80ms to +80ms
// - finalDelay = 400ms + jitter = 320ms to 480ms
//
// Context awareness:
// - select on two channels: time.After and ctx.Done()
// - If context cancelled, return immediately with ctx.Err()
// - Prevents waiting full backoff duration after timeout
//
// func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error) {
//     var zero T
//     var lastErr error
//
//     for attempt := 0; attempt <= c.MaxRetries; attempt++ {
//         // Try request
//         val, err := doRequest[T](ctx, c.HTTP, url)
//         if err == nil {
//             return val, nil  // Success!
//         }
//         lastErr = err
//
//         // Check if retryable
//         if !isRetryable(err) || attempt == c.MaxRetries {
//             break  // Don't retry
//         }
//
//         // Calculate backoff with jitter
//         delay := c.BaseDelay * time.Duration(1<<uint(attempt))  // 2^attempt
//         jitter := time.Duration(rand.Float64()*0.4-0.2) * delay  // ±20%
//         delay += jitter
//
//         // Wait with context awareness
//         select {
//         case <-time.After(delay):
//             // Continue to next attempt
//         case <-ctx.Done():
//             return zero, ctx.Err()
//         }
//     }
//
//     return zero, fmt.Errorf("all retries failed: %w", lastErr)
// }

// ============================================================================
// Exercise 3: Implement doRequest
// ============================================================================

// doRequest executes a single HTTP GET request and decodes JSON.
//
// TODO: Implement HTTP request with generic JSON decoding:
// 1. Create HTTP request with context (http.NewRequestWithContext)
// 2. Execute request using the provided client
// 3. Ensure response body is closed (defer resp.Body.Close())
// 4. Check status code (only 200 OK is success)
// 5. Read entire response body (io.ReadAll)
// 6. Unmarshal JSON into type T
// 7. Return decoded value
//
// Generic JSON decoding:
// - T is the target type (e.g., map[string]interface{}, struct)
// - json.Unmarshal requires a pointer, so we use &result
// - Zero value pattern: var zero T returns zero value for type T
//
// Error handling:
// - Network errors: DNS lookup, connection refused, timeout
// - HTTP errors: Non-200 status codes
// - JSON errors: Invalid JSON, type mismatch
// - All errors propagate up to retry logic
//
// Memory management:
// - resp.Body.Close() is CRITICAL to avoid leaks
// - Use defer to ensure closure even on error
// - io.ReadAll allocates buffer for entire response
//
// func doRequest[T any](ctx context.Context, client *http.Client, url string) (T, error) {
//     var zero T
//
//     // Create request with context
//     req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//     if err != nil {
//         return zero, err
//     }
//
//     // Execute request
//     resp, err := client.Do(req)
//     if err != nil {
//         return zero, err
//     }
//     defer resp.Body.Close()  // CRITICAL: prevent resource leak
//
//     // Check status code
//     if resp.StatusCode != http.StatusOK {
//         return zero, fmt.Errorf("HTTP %d", resp.StatusCode)
//     }
//
//     // Read response body
//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return zero, err
//     }
//
//     // Decode JSON
//     var result T
//     if err := json.Unmarshal(body, &result); err != nil {
//         return zero, err
//     }
//
//     return result, nil
// }

// ============================================================================
// Exercise 4: Implement isRetryable
// ============================================================================

// isRetryable determines if an error should be retried.
//
// TODO: Implement retry decision logic:
// - For this simple implementation, retry all errors
// - In production, you would check error types:
//   * Network errors (DNS, connection): RETRY
//   * Timeout errors: RETRY
//   * 5xx status codes (server errors): RETRY
//   * 429 Too Many Requests: RETRY (with longer backoff)
//   * 4xx status codes (client errors): DON'T RETRY
//   * JSON decode errors: DON'T RETRY
//
// Simple implementation (retry everything):
// func isRetryable(err error) bool {
//     return err != nil
// }
//
// Production implementation would look like:
// func isRetryable(err error) bool {
//     if err == nil {
//         return false
//     }
//
//     // Check for specific error types
//     var netErr net.Error
//     if errors.As(err, &netErr) && netErr.Timeout() {
//         return true  // Timeout errors are retryable
//     }
//
//     // Check for HTTP status codes
//     if strings.Contains(err.Error(), "HTTP 5") {
//         return true  // 5xx server errors are retryable
//     }
//     if strings.Contains(err.Error(), "HTTP 429") {
//         return true  // Rate limit errors are retryable
//     }
//     if strings.Contains(err.Error(), "HTTP 4") {
//         return false  // 4xx client errors are NOT retryable
//     }
//
//     // Network errors (connection refused, DNS, etc.)
//     return true
// }

// ============================================================================
// Comparison with Other Approaches
// ============================================================================
//
// Alternative 1: Fixed delay (no exponential backoff)
//   for attempt := 0; attempt <= maxRetries; attempt++ {
//       result, err := doRequest(ctx, url)
//       if err == nil { return result, nil }
//       time.Sleep(1 * time.Second)  // Same delay every time
//   }
//   Pros: Simpler implementation
//   Cons: Too aggressive (hammers failing server)
//         No adaptation to severity of failure
//
// Alternative 2: No jitter
//   delay := baseDelay * (1 << attempt)
//   time.Sleep(delay)
//   Pros: Simpler; deterministic
//   Cons: Thundering herd problem (all clients retry at same time)
//
// Alternative 3: Fibonacci backoff
//   delays: 100ms, 100ms, 200ms, 300ms, 500ms, 800ms
//   Pros: Grows slower than exponential
//   Cons: More complex; less standard than exponential
//
// Alternative 4: Full jitter (not additive)
//   delay := rand.Duration(0, baseDelay * (1 << attempt))
//   Pros: Maximum desynchronization
//   Cons: Less predictable; may retry too quickly
//
// Alternative 5: Circuit breaker pattern
//   if circuitBreaker.IsOpen() {
//       return nil, ErrCircuitOpen
//   }
//   result, err := doRequest(ctx, url)
//   circuitBreaker.RecordResult(err)
//   Pros: Stops retries if failure rate is high (protects backend)
//   Cons: More complex; needs state management
//
// Alternative 6: Request hedging
//   After delay, send duplicate request
//   Return whichever completes first
//   Pros: Lower latency at high percentiles (p99)
//   Cons: Increased load on backend; wasteful
//
// ============================================================================
// Go vs Other Languages
// ============================================================================
//
// Go vs Python (requests + retry):
//   from requests.adapters import HTTPAdapter
//   from urllib3.util.retry import Retry
//
//   session = requests.Session()
//   retry = Retry(total=3, backoff_factor=0.1)
//   adapter = HTTPAdapter(max_retries=retry)
//   session.mount('http://', adapter)
//   response = session.get(url)
//
//   Pros: Built-in retry support in libraries
//   Cons: Less control over backoff strategy
//         No type-safe JSON decoding
//         Synchronous (blocks thread)
//   Go: More explicit control
//       Type-safe generics
//       Lightweight goroutines (can retry many requests concurrently)
//
// Go vs JavaScript (fetch + retry):
//   async function fetchWithRetry(url, maxRetries = 3) {
//       for (let attempt = 0; attempt <= maxRetries; attempt++) {
//           try {
//               const response = await fetch(url);
//               if (response.ok) return await response.json();
//           } catch (err) {
//               if (attempt === maxRetries) throw err;
//               await new Promise(r => setTimeout(r, 100 * (2 ** attempt)));
//           }
//       }
//   }
//
//   Pros: Similar structure with async/await
//   Cons: Single-threaded (can't parallelize retries)
//         No type safety without TypeScript
//         Manual retry implementation needed
//   Go: Built-in concurrency
//       Type-safe generics
//       Standard library HTTP client is production-grade
//
// Go vs Rust (reqwest + retry):
//   use reqwest_retry::{RetryTransientMiddleware, policies::ExponentialBackoff};
//
//   let retry_policy = ExponentialBackoff::builder()
//       .build_with_max_retries(3);
//   let client = ClientBuilder::new(reqwest::Client::new())
//       .with(RetryTransientMiddleware::new_with_policy(retry_policy))
//       .build();
//   let response = client.get(url).send().await?;
//
//   Pros: Zero-cost abstractions; excellent type safety
//   Cons: More complex (trait bounds, async/await, Result types)
//         Requires external crate (not in stdlib)
//         Steeper learning curve
//   Go: Simpler error handling (no Result type)
//       Standard library is comprehensive
//       Faster development iteration
//
// Go vs Java (OkHttp with retry):
//   OkHttpClient client = new OkHttpClient.Builder()
//       .addInterceptor(new RetryInterceptor(3, 100))
//       .build();
//   Request request = new Request.Builder().url(url).build();
//   Response response = client.newCall(request).execute();
//
//   Pros: Similar interceptor pattern
//   Cons: Much more verbose (builder pattern, checked exceptions)
//         Heavyweight threads (not goroutines)
//         No generics for JSON (requires Gson/Jackson)
//   Go: Cleaner syntax
//       Built-in JSON support
//       Lightweight goroutines
//
// ============================================================================
// Advanced Topics
// ============================================================================
//
// 1. Idempotency tokens:
//    - Attach unique ID to requests
//    - Server deduplicates based on token
//    - Safe to retry POST/PUT operations
//
// 2. Circuit breaker:
//    - Track failure rate over time
//    - "Open" circuit after threshold (stop retries)
//    - "Half-open" state to test recovery
//    - "Closed" state for normal operation
//
// 3. Rate limiting:
//    - Respect server's rate limit headers
//    - Adjust retry delay based on Retry-After header
//    - Client-side rate limiting to prevent overload
//
// 4. Observability:
//    - Metrics: retry count, success rate, latency
//    - Logging: log each retry attempt
//    - Tracing: propagate trace IDs across retries
//
// 5. Adaptive backoff:
//    - Adjust backoff based on server feedback
//    - Longer backoff if server is overloaded
//    - Shorter backoff if server is healthy
//
// ============================================================================
// After Implementation
// ============================================================================
//
// Testing your implementation:
// - Run: go test -v ./...
// - Test with failing server (simulate network errors)
// - Test with slow server (context timeout)
// - Test with transient failures (success after retry)
// - Verify exponential backoff timing
// - Verify jitter spreads retries
//
// Performance tips:
// - Set reasonable MaxRetries (3-5 is typical)
// - Choose BaseDelay based on your system (100ms-1s)
// - Total retry time should be < context timeout
// - Monitor retry rate in production (high rate = problem)
//
// Common mistakes:
// - Not closing response body (resource leak)
// - Retrying non-idempotent operations (POST)
// - Retrying 4xx errors (waste of resources)
// - No context timeout (retries can run forever)
// - Forgetting to add jitter (thundering herd)
// - BaseDelay too large (slow to recover from transients)
