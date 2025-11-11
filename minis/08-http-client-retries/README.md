# Project 08: HTTP Client with Retries

## 1. What Is This About?

### Real-World Scenario

You're building a service that calls a third-party API to fetch user data. The API is generally reliable, but occasionally:
- Returns `503 Service Unavailable` (servers overloaded)
- Times out (slow response)
- Returns network errors (temporary connectivity issues)

**❌ Naive approach:** Fail immediately on first error
```go
resp, err := http.Get(url)
if err != nil {
    return err  // Fails on transient errors!
}
```

**Result**: Your service fails even though the problem was temporary. If you had waited 100ms and tried again, it would have worked.

**✅ Resilient approach:** Retry with exponential backoff
```go
for attempt := 0; attempt < 3; attempt++ {
    resp, err := http.Get(url)
    if err == nil {
        return resp  // Success!
    }
    time.Sleep(100ms * (2^attempt))  // 100ms, 200ms, 400ms
}
return err  // All retries failed
```

**Result**: Transient errors are automatically handled. Your service is much more reliable.

This project teaches you how to build **production-grade HTTP clients** that are:
- **Resilient**: Automatically retry transient failures
- **Smart**: Use exponential backoff to avoid overwhelming servers
- **Fair**: Add jitter to prevent thundering herd
- **Context-aware**: Respect timeouts and cancellation
- **Type-safe**: Generic JSON decoding without type assertions

### What You'll Learn

1. **Retry strategies**: When and how to retry failed requests
2. **Exponential backoff**: Why delays should increase exponentially
3. **Jitter**: Preventing thundering herd problem
4. **Error classification**: Retryable vs non-retryable errors
5. **Context propagation**: Timeouts and cancellation in HTTP clients
6. **Generic functions**: Type-safe JSON decoding with generics

### The Challenge

Build an HTTP client that:
- Retries failed requests up to N times
- Uses exponential backoff (100ms, 200ms, 400ms, ...)
- Adds jitter (±20% randomness) to backoff delays
- Respects context timeouts and cancellation
- Only retries retriable errors (not client bugs)
- Provides type-safe JSON decoding

---

## 2. First Principles: Network Resilience

### Why Do Network Requests Fail?

Networks are unreliable. Requests can fail for many reasons:

**Transient failures** (temporary, worth retrying):
- 📡 **Network blip**: Packet loss, DNS hiccup, routing issue
- ⏱️ **Timeout**: Server took too long to respond
- 🔄 **Server overload**: `503 Service Unavailable`, `429 Too Many Requests`
- 🔧 **Server restart**: Brief downtime during deployment

**Permanent failures** (won't succeed even if retried):
- 🚫 **Client error**: `400 Bad Request` (malformed request)
- 🔐 **Auth error**: `401 Unauthorized`, `403 Forbidden`
- 🔍 **Not found**: `404 Not Found`
- ❌ **Invalid URL**: Malformed endpoint

**Key insight**: Retry transient failures, fail fast on permanent failures.

### What is Exponential Backoff?

**Exponential backoff** means doubling the delay after each retry.

**Linear backoff** (bad):
```
Attempt 1: wait 1 second
Attempt 2: wait 2 seconds  (+1)
Attempt 3: wait 3 seconds  (+1)
```

**Exponential backoff** (good):
```
Attempt 1: wait 100ms
Attempt 2: wait 200ms  (×2)
Attempt 3: wait 400ms  (×2)
Attempt 4: wait 800ms  (×2)
```

**Formula**: `delay = baseDelay × 2^attempt`

**Why exponential?**

1. **Gives servers time to recover**: If a server is overloaded, hitting it harder (linear) makes it worse. Exponential gives it breathing room.

2. **Balances retry speed and server load**:
   - First retry (100ms): Quick, catches brief hiccups
   - Later retries (800ms, 1.6s): Slower, for serious outages

3. **Industry standard**: Used by AWS, Google Cloud, Kubernetes, etc.

### What is Jitter?

**Jitter** adds randomness to retry delays.

**Without jitter (Thundering Herd Problem)**:

Imagine 10,000 clients all make a request to a server:
1. Server crashes or returns 503
2. **All 10,000 clients** wait exactly 100ms (exponential backoff)
3. **All 10,000 clients** retry at the exact same millisecond
4. Server gets slammed with 10,000 simultaneous requests
5. Server crashes again
6. **All 10,000 clients** wait exactly 200ms
7. Repeat... server can never recover!

**With jitter**:
```
Attempt 1:
- Client A: wait 95ms   (100ms - 5ms jitter)
- Client B: wait 102ms  (100ms + 2ms jitter)
- Client C: wait 88ms   (100ms - 12ms jitter)
```

Retries are **spread out over time** instead of all hitting at once → server can recover.

**Formula**: `delay = (baseDelay × 2^attempt) × (1 ± random(20%))`

**Common jitter algorithms**:
- **Full jitter**: `delay = random(0, base × 2^attempt)`
- **Decorrelated jitter**: More complex, used by AWS SDK
- **±20% jitter**: `delay = base × 2^attempt × (0.8 to 1.2)` ← We use this

### What is Context?

`context.Context` provides deadlines, cancellation, and request-scoped values.

**Use case 1: Overall timeout**
```go
// Total request must complete in 5 seconds (including retries)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := GetJSON[User](ctx, client, url)
```

**Use case 2: User cancellation**
```go
// User presses "Cancel" button
ctx, cancel := context.WithCancel(context.Background())

go func() {
    <-cancelButton
    cancel()  // Abort request immediately
}()

result, err := GetJSON[User](ctx, client, url)
```

**Key insight**: Context propagates through the entire call stack. When you cancel at the top, everything stops.

---

## 3. Breaking Down the Solution

### Step 1: Client Configuration

```go
type Client struct {
    HTTP       *http.Client   // Underlying HTTP client
    MaxRetries int            // Maximum retry attempts (e.g., 3)
    BaseDelay  time.Duration  // Initial delay (e.g., 100ms)
}
```

**Why wrap `http.Client`?**
- Add retry logic on top of Go's standard HTTP client
- Keep the standard client's connection pooling, TLS, etc.

### Step 2: Retry Loop

```
for attempt from 0 to MaxRetries:
    1. Try request
    2. If success → return result
    3. If non-retryable error → return error immediately
    4. If last attempt → return error
    5. Calculate backoff with jitter
    6. Sleep (with context awareness)
    7. Repeat
```

**Pseudocode**:
```
attempt = 0
while attempt <= MaxRetries:
    result, err = doRequest()
    if err == nil:
        return result  // Success!

    if not retryable(err):
        return err  // Permanent failure

    if attempt == MaxRetries:
        return err  // Out of retries

    delay = baseDelay * (2 ^ attempt)
    delay += jitter(delay)

    wait(delay) or context cancelled

    attempt++
```

### Step 3: Exponential Backoff Calculation

```go
delay := c.BaseDelay * time.Duration(1 << uint(attempt))
```

**What is `1 << uint(attempt)`?**
- Left bit shift = multiply by 2^attempt
- `1 << 0` = 1
- `1 << 1` = 2
- `1 << 2` = 4
- `1 << 3` = 8

**Example** (BaseDelay = 100ms):
- Attempt 0: `100ms * (1 << 0)` = 100ms × 1 = 100ms
- Attempt 1: `100ms * (1 << 1)` = 100ms × 2 = 200ms
- Attempt 2: `100ms * (1 << 2)` = 100ms × 4 = 400ms
- Attempt 3: `100ms * (1 << 3)` = 100ms × 8 = 800ms

### Step 4: Adding Jitter

```go
jitter := time.Duration(rand.Float64()*0.4-0.2) * delay  // ±20%
delay += jitter
```

**How this works**:
- `rand.Float64()` returns value in [0.0, 1.0)
- `rand.Float64() * 0.4` returns value in [0.0, 0.4)
- `rand.Float64() * 0.4 - 0.2` returns value in [-0.2, 0.2) → ±20%

**Example** (delay = 100ms):
- `jitter = rand() * 0.4 - 0.2` → random value in [-0.2, 0.2]
- If jitter = -0.15: `delay = 100ms + (100ms × -0.15)` = 85ms
- If jitter = +0.18: `delay = 100ms + (100ms × 0.18)` = 118ms

Result: Delays range from 80ms to 120ms (±20% of 100ms)

### Step 5: Context-Aware Sleep

```go
select {
case <-time.After(delay):
    // Delay finished, continue to next attempt
case <-ctx.Done():
    return zero, ctx.Err()  // Context cancelled or timed out
}
```

**Why use `select` instead of `time.Sleep`?**

`time.Sleep(delay)` would block even if context is cancelled.

`select` allows us to abort the sleep if context is cancelled.

**Example**:
- Total timeout: 1 second
- We're on attempt 3, delay = 800ms
- But only 200ms left in timeout
- After 200ms, `ctx.Done()` is closed → we return immediately instead of waiting the full 800ms

---

## 4. Complete Solution Walkthrough

### GetJSON Function Signature

```go
func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error)
```

**Type parameter `T`**: The expected JSON response type

**Usage**:
```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

client := &Client{
    HTTP:       &http.Client{Timeout: 10 * time.Second},
    MaxRetries: 3,
    BaseDelay:  100 * time.Millisecond,
}

user, err := GetJSON[User](ctx, client, "https://api.example.com/user/123")
```

**Why generic?**

Without generics:
```go
func GetJSON(ctx context.Context, c *Client, url string) (interface{}, error)

result, _ := GetJSON(ctx, client, url)
user := result.(User)  // Type assertion required, runtime error if wrong type
```

With generics:
```go
user, err := GetJSON[User](ctx, client, url)
// `user` is type `User`, no type assertion needed, compile-time type safety
```

### Retry Loop Implementation

```go
var zero T
var lastErr error

for attempt := 0; attempt <= c.MaxRetries; attempt++ {
    result, err := doRequest[T](ctx, c.HTTP, url)
    if err == nil {
        return result, nil  // Success!
    }

    lastErr = err

    if !isRetryable(err) {
        return zero, err  // Non-retryable error
    }

    if attempt == c.MaxRetries {
        break  // Out of retries
    }

    // Backoff calculation...
    // Wait...
}

return zero, fmt.Errorf("all retries failed: %w", lastErr)
```

**Line-by-line**:

1. **`var zero T`**: Zero value of type T (for error returns)
2. **`for attempt := 0; attempt <= c.MaxRetries`**:
   - MaxRetries = 3 means 4 total attempts (0, 1, 2, 3)
   - Attempt 0 = initial try
   - Attempts 1-3 = retries
3. **`result, err := doRequest[T](...)`**: Try the request
4. **`if err == nil { return result, nil }`**: Success! Return immediately
5. **`lastErr = err`**: Save error in case all retries fail
6. **`if !isRetryable(err) { return zero, err }`**: Permanent error, don't retry
7. **`if attempt == c.MaxRetries { break }`**: Out of retries, exit loop
8. **Backoff and sleep** (shown next)

### doRequest Function

```go
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
```

**Key points**:

1. **`http.NewRequestWithContext(ctx, ...)`**: Attaches context to request
   - If context is cancelled, HTTP client aborts the request
   - If context has deadline, HTTP client enforces it

2. **`defer resp.Body.Close()`**: Always close response body
   - Prevents connection leaks
   - Returns connection to pool

3. **`if resp.StatusCode != http.StatusOK`**: Check HTTP status
   - In production, check specific codes (200, 201, etc.)

4. **`json.Unmarshal(body, &result)`**: Decode JSON into type T
   - Type-safe thanks to generics
   - Compile error if T is not JSON-decodable

### Error Classification

```go
func isRetryable(err error) bool {
    // Simple heuristic: retry on network errors, not on parse errors
    // In production, check specific error types
    return err != nil
}
```

**Production implementation**:

```go
func isRetryable(err error) bool {
    if err == nil {
        return false
    }

    // Check for specific error types
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Timeout() {
        return true  // Network timeout
    }

    // Check HTTP status codes
    if strings.Contains(err.Error(), "HTTP 503") ||
       strings.Contains(err.Error(), "HTTP 429") ||
       strings.Contains(err.Error(), "HTTP 500") {
        return true  // Server error
    }

    // Check for specific errors
    if errors.Is(err, context.DeadlineExceeded) {
        return false  // Overall timeout exceeded, don't retry
    }

    return false  // Unknown error, don't retry
}
```

**Retryable errors**:
- Network timeouts
- Connection refused (server down)
- `500 Internal Server Error`
- `503 Service Unavailable`
- `429 Too Many Requests`

**Non-retryable errors**:
- `400 Bad Request` (client bug)
- `401 Unauthorized` (auth issue)
- `404 Not Found` (wrong URL)
- JSON parse errors (response malformed)
- Context deadline exceeded (out of time)

---

## 5. Key Concepts Explained

### Concept 1: Retry Budgets

In production, you don't want infinite retries. You need **retry budgets**.

**Problem**: If your service retries aggressively, it can amplify load on downstream services.

**Example**:
- Your service gets 1000 req/sec
- Downstream API is down (100% failure rate)
- You retry 3 times per request
- Downstream API receives 1000 × (1 + 3) = **4000 req/sec**
- You've 4x the load on an already-struggling service!

**Solution**: Limit retries as a percentage of total requests.

```go
type RetryBudget struct {
    mu              sync.Mutex
    totalRequests   int
    retriedRequests int
    maxRetryRatio   float64  // e.g., 0.1 = max 10% retries
}

func (rb *RetryBudget) CanRetry() bool {
    rb.mu.Lock()
    defer rb.mu.Unlock()

    if rb.totalRequests == 0 {
        return true
    }

    ratio := float64(rb.retriedRequests) / float64(rb.totalRequests)
    return ratio < rb.maxRetryRatio
}
```

### Concept 2: Circuit Breaker Pattern

**Circuit breaker** stops retrying if failure rate is consistently high.

**States**:
1. **Closed** (normal): Requests go through
2. **Open** (failing): Requests immediately fail (no attempt)
3. **Half-open** (testing): Allow a few requests to test if service recovered

**State transitions**:
```
Closed → Open: If error rate > threshold
Open → Half-open: After cooldown period
Half-open → Closed: If test requests succeed
Half-open → Open: If test requests fail
```

**Why use it?**

Prevents cascading failures. If a downstream service is dead, stop hammering it and fail fast.

### Concept 3: Request Hedging

**Hedging**: Send duplicate requests to improve latency.

**Algorithm**:
1. Send request to server A
2. Wait 50ms
3. If no response yet, send duplicate request to server B
4. Return whichever responds first

**Use case**: When you need low latency and have spare capacity.

**Trade-off**: Increases server load, but reduces tail latency.

### Concept 4: Adaptive Backoff

**Adaptive backoff** adjusts delay based on server feedback.

**Example**: `Retry-After` header
```
HTTP/1.1 503 Service Unavailable
Retry-After: 120
```

Server says "retry after 120 seconds," so you wait exactly that (instead of exponential backoff).

**Benefits**:
- Respects server's capacity
- Faster recovery (no need to guess)

### Concept 5: Generic Functions in Go

**Before Go 1.18**:
```go
func GetJSON(ctx context.Context, client *http.Client, url string) (interface{}, error) {
    // ...
    var result interface{}
    json.Unmarshal(body, &result)
    return result, nil
}

// Usage:
rawUser, _ := GetJSON(ctx, client, url)
user := rawUser.(User)  // Runtime type assertion
```

**With Go 1.18+ generics**:
```go
func GetJSON[T any](ctx context.Context, client *http.Client, url string) (T, error) {
    // ...
    var result T
    json.Unmarshal(body, &result)
    return result, nil
}

// Usage:
user, _ := GetJSON[User](ctx, client, url)  // Compile-time type safety
```

**Benefits**:
- Type safety: Compiler catches type mismatches
- No type assertions: Cleaner code
- Code reuse: One function for all types

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Generic HTTP Client

```go
type HTTPClient interface {
    Do(*http.Request) (*http.Response, error)
}

func Fetch[T any](ctx context.Context, client HTTPClient, method, url string, body io.Reader) (T, error) {
    var zero T

    req, err := http.NewRequestWithContext(ctx, method, url, body)
    if err != nil {
        return zero, err
    }

    resp, err := client.Do(req)
    if err != nil {
        return zero, err
    }
    defer resp.Body.Close()

    var result T
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return zero, err
    }

    return result, nil
}
```

### Pattern 2: Retry With Custom Predicate

```go
func RetryWithPredicate[T any](
    ctx context.Context,
    fn func() (T, error),
    shouldRetry func(error) bool,
    maxRetries int,
    baseDelay time.Duration,
) (T, error) {
    var zero T

    for attempt := 0; attempt <= maxRetries; attempt++ {
        result, err := fn()
        if err == nil {
            return result, nil
        }

        if !shouldRetry(err) {
            return zero, err
        }

        if attempt == maxRetries {
            return zero, err
        }

        delay := baseDelay * time.Duration(1<<uint(attempt))
        select {
        case <-time.After(delay):
        case <-ctx.Done():
            return zero, ctx.Err()
        }
    }

    return zero, fmt.Errorf("unreachable")
}
```

### Pattern 3: Circuit Breaker

```go
type CircuitBreaker struct {
    mu            sync.Mutex
    state         State  // Closed, Open, HalfOpen
    failures      int
    threshold     int
    cooldown      time.Duration
    lastFailTime  time.Time
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.state == Open {
        if time.Since(cb.lastFailTime) > cb.cooldown {
            cb.state = HalfOpen
        } else {
            return errors.New("circuit breaker open")
        }
    }

    err := fn()

    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()

        if cb.failures >= cb.threshold {
            cb.state = Open
        }
    } else {
        cb.failures = 0
        cb.state = Closed
    }

    return err
}
```

### Pattern 4: Request Timeout with Context

```go
func FetchWithTimeout[T any](client *http.Client, url string, timeout time.Duration) (T, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    return GetJSON[T](ctx, &Client{HTTP: client}, url)
}
```

### Pattern 5: Parallel Requests with Fastest Wins

```go
func FetchFastest[T any](ctx context.Context, client *http.Client, urls []string) (T, error) {
    type result struct {
        data T
        err  error
    }

    results := make(chan result, len(urls))

    for _, url := range urls {
        go func(u string) {
            data, err := GetJSON[T](ctx, &Client{HTTP: client}, u)
            results <- result{data, err}
        }(url)
    }

    select {
    case r := <-results:
        return r.data, r.err
    case <-ctx.Done():
        return *new(T), ctx.Err()
    }
}
```

---

## 7. Real-World Applications

### Microservices Communication

**Use case**: Service A calls Service B's API

```go
type UserService struct {
    client *Client
}

func (us *UserService) GetUser(ctx context.Context, userID int) (*User, error) {
    url := fmt.Sprintf("http://user-service/api/users/%d", userID)
    return GetJSON[User](ctx, us.client, url)
}
```

Companies using this: Netflix, Uber, Airbnb (all use retries between microservices)

### External API Integration

**Use case**: Calling third-party APIs (Stripe, Twilio, SendGrid)

```go
type StripeClient struct {
    client *Client
    apiKey string
}

func (sc *StripeClient) CreateCharge(ctx context.Context, amount int) (*Charge, error) {
    // Stripe recommends exponential backoff for retries
    return GetJSON[Charge](ctx, sc.client, "https://api.stripe.com/v1/charges")
}
```

Companies: Every SaaS that integrates with external APIs

### Database Connection Pools

**Use case**: Retry database queries on transient errors

```go
func QueryWithRetry[T any](ctx context.Context, db *sql.DB, query string, args ...interface{}) ([]T, error) {
    return RetryWithPredicate(
        ctx,
        func() ([]T, error) {
            return executeQuery[T](ctx, db, query, args...)
        },
        isTransientDBError,
        3,
        100*time.Millisecond,
    )
}
```

### Cloud Storage Downloads

**Use case**: Download files from S3, GCS with retries on network failures

```go
func DownloadFile(ctx context.Context, bucket, key string) ([]byte, error) {
    client := &Client{
        HTTP:       &http.Client{Timeout: 30 * time.Second},
        MaxRetries: 3,
        BaseDelay:  200 * time.Millisecond,
    }

    url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
    return GetJSON[[]byte](ctx, client, url)
}
```

### Webhook Delivery

**Use case**: Deliver webhooks to customer endpoints (might be unreliable)

```go
type WebhookDelivery struct {
    client *Client
}

func (wd *WebhookDelivery) Send(ctx context.Context, url string, payload interface{}) error {
    // Retry up to 5 times with exponential backoff
    _, err := PostJSON[Response](ctx, wd.client, url, payload)
    return err
}
```

Companies: GitHub, Stripe, Shopify (all retry webhook delivery)

---

## 8. Common Mistakes to Avoid

### Mistake 1: Retrying Non-Retryable Errors

**❌ Wrong**:
```go
for attempt := 0; attempt < 3; attempt++ {
    resp, err := http.Get(url)
    if err == nil {
        return resp
    }
    time.Sleep(100 * time.Millisecond)
}
```

**Problem**: Retries `400 Bad Request` (client bug), `404 Not Found`, etc. Wastes time and resources.

**✅ Correct**:
```go
for attempt := 0; attempt < 3; attempt++ {
    resp, err := http.Get(url)
    if err == nil {
        return resp
    }

    if !isRetryable(err) {
        return nil, err  // Fail fast on permanent errors
    }

    time.Sleep(100 * time.Millisecond)
}
```

### Mistake 2: No Jitter

**❌ Wrong**:
```go
delay := 100 * time.Millisecond * time.Duration(1<<uint(attempt))
time.Sleep(delay)
```

**Problem**: Thundering herd. All clients retry at exactly the same time.

**✅ Correct**:
```go
delay := 100 * time.Millisecond * time.Duration(1<<uint(attempt))
jitter := time.Duration(rand.Float64()*0.4-0.2) * delay
delay += jitter
time.Sleep(delay)
```

### Mistake 3: Ignoring Context

**❌ Wrong**:
```go
time.Sleep(5 * time.Second)  // Sleeps even if context is cancelled
```

**Problem**: If user cancels request, we still wait 5 seconds.

**✅ Correct**:
```go
select {
case <-time.After(5 * time.Second):
case <-ctx.Done():
    return ctx.Err()
}
```

### Mistake 4: Unbounded Retries

**❌ Wrong**:
```go
for {
    resp, err := http.Get(url)
    if err == nil {
        return resp
    }
    time.Sleep(100 * time.Millisecond)
}
```

**Problem**: Infinite loop if service is permanently down.

**✅ Correct**:
```go
for attempt := 0; attempt <= maxRetries; attempt++ {
    // ...
}
```

### Mistake 5: Not Using Exponential Backoff

**❌ Wrong**:
```go
for attempt := 0; attempt < 3; attempt++ {
    resp, err := http.Get(url)
    if err == nil {
        return resp
    }
    time.Sleep(100 * time.Millisecond)  // Fixed delay
}
```

**Problem**: Hammers server with constant rate → can prevent recovery.

**✅ Correct**:
```go
delay := 100 * time.Millisecond * time.Duration(1<<uint(attempt))
time.Sleep(delay)
```

### Mistake 6: Forgetting to Close Response Body

**❌ Wrong**:
```go
resp, err := http.Get(url)
if err != nil {
    return err
}
// Forgot defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
```

**Problem**: Connection leak. Eventually runs out of file descriptors.

**✅ Correct**:
```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
```

### Mistake 7: Using time.Sleep Instead of time.After

**❌ Wrong**:
```go
time.Sleep(delay)  // Can't be interrupted by context
```

**✅ Correct**:
```go
select {
case <-time.After(delay):
case <-ctx.Done():
    return ctx.Err()
}
```

---

## 9. Stretch Goals

### Goal 1: Implement Circuit Breaker ⭐⭐

Add circuit breaker pattern to stop retrying if failure rate is consistently high.

**Hint**:
```go
type CircuitBreaker struct {
    mu            sync.Mutex
    state         string  // "closed", "open", "half-open"
    failures      int
    threshold     int
    cooldown      time.Duration
    lastFailTime  time.Time
}

func (cb *CircuitBreaker) Allow() bool {
    // Check if request should be allowed
}
```

### Goal 2: Add Request Hedging ⭐⭐⭐

Send duplicate requests after a delay to reduce tail latency.

**Hint**:
```go
func GetJSONHedged[T any](ctx context.Context, c *Client, url string, hedgeDelay time.Duration) (T, error) {
    resultCh := make(chan result[T], 2)

    go func() {
        data, err := GetJSON[T](ctx, c, url)
        resultCh <- result[T]{data, err}
    }()

    select {
    case r := <-resultCh:
        return r.data, r.err
    case <-time.After(hedgeDelay):
        // Send hedged request
        go func() {
            data, err := GetJSON[T](ctx, c, url)
            resultCh <- result[T]{data, err}
        }()
    }

    r := <-resultCh
    return r.data, r.err
}
```

### Goal 3: Add Prometheus Metrics ⭐⭐

Export retry metrics to Prometheus.

**Hint**:
```go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"status"},
    )
    httpRetriesTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "http_retries_total",
            Help: "Total retry attempts",
        },
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal, httpRetriesTotal)
}
```

### Goal 4: Implement Retry Budget ⭐⭐⭐

Limit retries to a percentage of total requests to prevent retry storms.

**Hint**: See "Concept 1: Retry Budgets" section above.

### Goal 5: Add Adaptive Backoff ⭐⭐⭐

Respect `Retry-After` header from server responses.

**Hint**:
```go
func parseRetryAfter(resp *http.Response) (time.Duration, bool) {
    header := resp.Header.Get("Retry-After")
    if header == "" {
        return 0, false
    }

    if seconds, err := strconv.Atoi(header); err == nil {
        return time.Duration(seconds) * time.Second, true
    }

    if t, err := time.Parse(time.RFC1123, header); err == nil {
        return time.Until(t), true
    }

    return 0, false
}
```

---

## How to Run

```bash
# Run tests
go test ./minis/08-http-client-retries/...

# Run with verbose output
go test -v ./minis/08-http-client-retries/...

# Run with race detector
go test -race ./minis/08-http-client-retries/...

# Run benchmarks
go test -bench=. ./minis/08-http-client-retries/...
```

---

## Summary

**What you learned**:
- ✅ Retry strategies for resilient network programming
- ✅ Exponential backoff prevents overwhelming servers
- ✅ Jitter prevents thundering herd problem
- ✅ Error classification (retryable vs permanent)
- ✅ Context propagation for timeouts and cancellation
- ✅ Generic functions for type-safe JSON decoding

**Why this matters**:
Network requests fail. Production systems must handle failures gracefully. Retries with exponential backoff and jitter are industry-standard resilience patterns used by every major tech company.

**Key formulas**:
- **Exponential backoff**: `delay = baseDelay × 2^attempt`
- **Jitter**: `delay = delay × (1 ± random(20%))`
- **Retry budget**: `retries / total < threshold`

**Next steps**:
- Project 09: Build HTTP servers with middleware and graceful shutdown
- Project 10: Learn gRPC for high-performance RPC

Stay resilient! 🚀
