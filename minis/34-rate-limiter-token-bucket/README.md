# Project 34: Rate Limiter with Token Bucket Algorithm

## 1. What Is This About?

### Real-World Scenario

You build an API for your startup. Within a week:

**❌ Without rate limiting:**
1. One user writes a buggy script hitting your API 10,000 times per second
2. Your servers get overwhelmed
3. Legitimate users can't access the API
4. Your AWS bill explodes to $10,000
5. Database connections exhausted
6. Site goes down, customers angry

**✅ With rate limiting:**
1. Buggy script hits your API repeatedly
2. Rate limiter blocks requests after 100/minute threshold
3. User gets "429 Too Many Requests" response
4. Legitimate users unaffected
5. Servers remain stable
6. You email the user about their script issue

This project teaches you how to build **production-grade rate limiters** using:
- **Token bucket algorithm**: Fair, bursty request handling
- **HTTP middleware**: Transparent rate limiting layer
- **Atomic operations**: Lock-free concurrent counters
- **Per-client tracking**: Independent limits for each user

### What You'll Learn

1. **Rate limiting fundamentals**: Why, when, and how to limit requests
2. **Token bucket algorithm**: Mathematical model for rate limiting
3. **Atomic operations**: `sync/atomic` for lock-free programming
4. **Time-based logic**: Refilling tokens, tracking time windows
5. **Middleware patterns**: Composable HTTP request interceptors
6. **Production patterns**: Memory cleanup, sliding windows

### The Challenge

Build a rate limiter with:
- Token bucket implementation (100 requests per minute per client)
- HTTP middleware that tracks clients by IP
- Atomic operations for thread-safe token management
- Automatic token refill over time
- Proper HTTP responses (429 Too Many Requests)
- Memory-efficient client cleanup

---

## 2. First Principles: What Is Rate Limiting?

### The Core Problem

**Problem**: Servers have finite resources (CPU, memory, network, database connections).

**Reality**: One bad actor can consume all resources, denying service to everyone else.

**Solution**: Limit the rate at which any single client can make requests.

### Why Rate Limiting Matters

**Real incidents**:
- **Twitter (2012)**: DDoS attack, no rate limits → site down for hours
- **GitHub (2018)**: Unintentional traffic spike → added stricter limits
- **Cloudflare**: Blocks 72 billion cyber threats/day using rate limiting

**Use cases**:
1. **DDoS protection**: Block malicious floods
2. **Fair resource allocation**: Prevent one user hogging resources
3. **Cost control**: Limit expensive API operations (OpenAI charges per token)
4. **Quality of service**: Guarantee baseline performance for all users
5. **Prevent abuse**: Stop scrapers, brute force attacks

### Common Rate Limiting Algorithms

| Algorithm | How It Works | Pros | Cons |
|-----------|-------------|------|------|
| **Fixed Window** | Allow N requests per time window (e.g., 100/min) | Simple | Burst at window edges |
| **Sliding Window** | Track requests in rolling time window | Smooth distribution | More memory |
| **Token Bucket** | Refill tokens over time, consume per request | Handles bursts gracefully | Slightly complex |
| **Leaky Bucket** | Requests drain at constant rate | Smooth output | Can lose requests |

**We'll implement Token Bucket** because it's:
- Industry standard (used by AWS, Stripe, GitHub)
- Handles legitimate bursts (user loads page with 10 API calls)
- Fair over time (doesn't penalize brief spikes)

---

## 3. The Token Bucket Algorithm

### Mental Model: Bucket of Tokens

Imagine a physical bucket:

```
         🪣 Bucket (capacity: 10 tokens)
         ┌─────────────────────┐
         │ 🪙 🪙 🪙 🪙 🪙      │ ← Currently 5 tokens
         │                     │
         │                     │
         └─────────────────────┘
              ↓ Refill: +2 tokens/second
```

**Rules**:
1. **Bucket has a maximum capacity** (e.g., 10 tokens)
2. **Tokens refill at a constant rate** (e.g., 2 tokens/second)
3. **Each request costs 1 token** (remove from bucket)
4. **If bucket empty → request denied** (429 error)
5. **Overflow tokens are discarded** (bucket can't exceed capacity)

### Example Timeline

**Settings**: Capacity = 10 tokens, Refill rate = 2 tokens/second

```
T=0s:  Bucket has 10 tokens
       User makes 5 requests → 5 tokens consumed
       Bucket now has 5 tokens

T=1s:  +2 tokens refilled
       Bucket now has 7 tokens
       User makes 1 request → 1 token consumed
       Bucket now has 6 tokens

T=2s:  +2 tokens refilled
       Bucket now has 8 tokens
       User makes 10 requests → first 8 succeed, next 2 fail
       Bucket now has 0 tokens

T=3s:  +2 tokens refilled
       Bucket now has 2 tokens
       User makes 1 request → succeeds
       Bucket now has 1 token
```

### Why Token Bucket Works Well

**1. Handles bursts gracefully**
```
User loads a dashboard → 10 simultaneous API calls
With token bucket: All 10 succeed if bucket is full
With fixed window: Might be spread across 2 windows, some fail
```

**2. Fair over time**
```
Attacker sends 1000 requests instantly
First N requests succeed (bucket capacity)
Remaining 1000-N fail immediately
Attacker can't consume more than refill rate over time
```

**3. Simple to reason about**
```
"You can make bursts of up to 100 requests,
but sustained rate is limited to 10 requests/second"
```

### Mathematical Model

**Key variables**:
- `capacity`: Maximum tokens in bucket
- `rate`: Tokens added per second
- `tokens`: Current tokens in bucket (0 ≤ tokens ≤ capacity)
- `lastRefill`: Timestamp of last refill

**Token refill calculation**:
```go
elapsed := time.Since(lastRefill).Seconds()
tokensToAdd := elapsed * rate
tokens = min(tokens + tokensToAdd, capacity)
```

**Request handling**:
```go
if tokens >= 1 {
    tokens -= 1
    return true  // Allow request
} else {
    return false // Deny request
}
```

### Token Bucket vs Other Algorithms

**Fixed Window Counter**:
```
Minute 1: [Request 1, Request 2, ..., Request 100] ✓
Minute 1: [Request 101] ✗ (limit exceeded)
          └─ Problem: At 00:59, user makes 100 requests
          └─ At 01:01, user makes 100 more requests
          └─ Result: 200 requests in 2 seconds!
```

**Token Bucket**:
```
User makes 100 requests at T=0
Bucket empties
At T=1s, only 10 tokens refilled
User can only make 10 more requests
Result: Sustained rate is 10/second (as intended)
```

---

## 4. Breaking Down the Solution

### Step 1: Token Bucket Data Structure

```go
type TokenBucket struct {
    capacity   int64         // Maximum tokens
    tokens     atomic.Int64  // Current tokens (atomic for concurrency)
    rate       float64       // Tokens per second
    lastRefill atomic.Int64  // Unix nano timestamp (atomic)
}
```

**Why atomic types?**

Without atomic:
```go
// Thread 1                 // Thread 2
tokens := b.tokens         tokens := b.tokens
tokens -= 1                tokens -= 1
b.tokens = tokens          b.tokens = tokens
// Data race! Both threads read the same value,
// decrement it, and write back. One decrement is lost!
```

With atomic:
```go
// Thread 1                 // Thread 2
b.tokens.Add(-1)           b.tokens.Add(-1)
// Atomic operations are serialized by CPU.
// Both decrements are applied correctly.
```

### Step 2: Token Refill Logic

```go
func (b *TokenBucket) refill() {
    now := time.Now().UnixNano()
    last := b.lastRefill.Load()

    elapsed := float64(now-last) / float64(time.Second)
    tokensToAdd := int64(elapsed * b.rate)

    if tokensToAdd > 0 {
        // Update timestamp atomically
        if b.lastRefill.CompareAndSwap(last, now) {
            // Add tokens, capping at capacity
            for {
                current := b.tokens.Load()
                new := current + tokensToAdd
                if new > b.capacity {
                    new = b.capacity
                }
                if b.tokens.CompareAndSwap(current, new) {
                    break
                }
            }
        }
    }
}
```

**Why CompareAndSwap (CAS)?**

CAS is an atomic operation: "If value is still X, update to Y; otherwise, do nothing."

```go
// Without CAS (race condition):
current := b.tokens.Load()  // Read: 50
new := current + 10         // Calculate: 60
// Another goroutine changes tokens to 45 here
b.tokens.Store(new)         // Write: 60 (WRONG! Lost the update to 45)

// With CAS (safe):
for {
    current := b.tokens.Load()  // Read: 50
    new := current + 10         // Calculate: 60
    if b.tokens.CompareAndSwap(current, new) {
        break  // Success! Value was 50, now 60
    }
    // If another goroutine changed value, CAS fails, loop retries
}
```

### Step 3: Allow Function

```go
func (b *TokenBucket) Allow() bool {
    b.refill()  // Always refill first

    // Try to consume 1 token
    for {
        current := b.tokens.Load()
        if current < 1 {
            return false  // No tokens available
        }
        if b.tokens.CompareAndSwap(current, current-1) {
            return true  // Successfully consumed token
        }
        // CAS failed (another goroutine changed value), retry
    }
}
```

**Why loop in Allow?**

Multiple goroutines might try to consume tokens simultaneously:
```
Goroutine A: current = 5, tries to CAS(5, 4)
Goroutine B: current = 5, tries to CAS(5, 4)

Result: Only one CAS succeeds. The other retries with updated value.
```

### Step 4: Rate Limiter (Per-Client Tracking)

```go
type RateLimiter struct {
    mu      sync.RWMutex
    buckets map[string]*TokenBucket  // key: client IP
    rate    float64
    capacity int64
}

func (rl *RateLimiter) Allow(clientID string) bool {
    bucket := rl.getBucket(clientID)
    return bucket.Allow()
}

func (rl *RateLimiter) getBucket(clientID string) *TokenBucket {
    // Try read lock first (fast path)
    rl.mu.RLock()
    bucket, exists := rl.buckets[clientID]
    rl.mu.RUnlock()

    if exists {
        return bucket
    }

    // Create new bucket (slow path, needs write lock)
    rl.mu.Lock()
    defer rl.mu.Unlock()

    // Double-check (another goroutine might have created it)
    if bucket, exists := rl.buckets[clientID]; exists {
        return bucket
    }

    bucket = NewTokenBucket(rl.capacity, rl.rate)
    rl.buckets[clientID] = bucket
    return bucket
}
```

**Why RWMutex?**

Most requests are from existing clients (read-heavy workload):
- **RLock**: Multiple goroutines can read simultaneously (fast)
- **Lock**: Only one goroutine can write, blocks all readers (slow)

### Step 5: HTTP Middleware

```go
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        clientIP := getClientIP(r)

        if !rl.Allow(clientIP) {
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", rl.rate))
            w.Header().Set("Retry-After", "1")
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func getClientIP(r *http.Request) string {
    // Check X-Forwarded-For header (proxy/load balancer)
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        ips := strings.Split(xff, ",")
        return strings.TrimSpace(ips[0])
    }

    // Check X-Real-IP header
    if xri := r.Header.Get("X-Real-IP"); xri != "" {
        return xri
    }

    // Fallback to RemoteAddr
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}
```

**Middleware execution flow**:
```
HTTP Request
    ↓
Rate Limiter Middleware
    ↓
  Allow()?
    ↓
  Yes: next.ServeHTTP() → Your handler
  No:  429 error → Request rejected
```

---

## 5. Key Concepts Explained

### Concept 1: Atomic Operations

**What are atomic operations?**

Operations that complete in a single, indivisible step. No other thread can see the operation "half-done."

**Non-atomic example**:
```go
// 64-bit increment on 32-bit system (2 steps)
var counter int64 = 100

// Thread 1: counter++
// Step 1: Read low 32 bits
// Step 2: Read high 32 bits  ← Thread 2 interrupts here
// Step 3: Increment
// Step 4: Write low 32 bits
// Step 5: Write high 32 bits

// Thread 2: counter++
// Reads partially updated value = DATA RACE
```

**Atomic example**:
```go
var counter atomic.Int64
counter.Store(100)

// Thread 1: counter.Add(1)  // Single atomic operation
// Thread 2: counter.Add(1)  // Single atomic operation
// Result: Always correct (101 or 102, never corrupted)
```

**Common atomic operations in Go**:
```go
var val atomic.Int64

val.Store(42)           // Set value
x := val.Load()         // Get value
val.Add(5)              // Atomic add (returns new value)
val.CompareAndSwap(old, new)  // CAS: if val==old, set to new
```

**When to use atomic vs mutex?**

| Use Atomic | Use Mutex |
|------------|-----------|
| Single variable | Multiple variables |
| Counters, flags | Complex state |
| Very hot path | Less frequent |
| Lock-free algorithms | Simpler code |

**Example: Atomic for counter**:
```go
var requestCount atomic.Int64

func handleRequest() {
    requestCount.Add(1)  // Fast, lock-free
    // ...
}
```

**Example: Mutex for complex state**:
```go
type Stats struct {
    mu      sync.Mutex
    count   int
    errors  int
    avgTime float64
}

func (s *Stats) Record(duration time.Duration, err error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.count++
    if err != nil {
        s.errors++
    }
    s.avgTime = (s.avgTime*float64(s.count-1) + duration.Seconds()) / float64(s.count)
}
```

### Concept 2: Lock-Free Programming

**Lock-free algorithm**: Progress is guaranteed for at least one thread, even if others are delayed.

**Token bucket is lock-free**:
```go
func (b *TokenBucket) Allow() bool {
    b.refill()  // Lock-free (uses CAS)

    for {
        current := b.tokens.Load()
        if current < 1 {
            return false
        }
        if b.tokens.CompareAndSwap(current, current-1) {
            return true  // Success!
        }
        // Retry if CAS failed (another thread won)
    }
}
```

**Why it's lock-free**:
- No mutex → No thread can block others indefinitely
- If 10 threads compete, at least 1 CAS succeeds per iteration
- Failed threads retry immediately (no waiting)

**Benefits**:
1. **High throughput**: No lock contention
2. **Scalability**: Performance doesn't degrade with more threads
3. **No deadlocks**: No locks to deadlock on

**Drawbacks**:
1. **Complexity**: Harder to reason about correctness
2. **ABA problem**: Value changes from A→B→A, CAS thinks nothing changed
3. **Retry loops**: Can waste CPU in extreme contention

### Concept 3: Memory Cleanup

**Problem**: Rate limiter tracks every client in memory. Over time, this grows unbounded.

```go
// After 1 million requests from different IPs:
rl.buckets contains 1 million entries
Each entry: ~100 bytes
Total memory: 100 MB (and growing!)
```

**Solution 1: Lazy cleanup (simple)**
```go
func (rl *RateLimiter) cleanup() {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    for id, bucket := range rl.buckets {
        if bucket.isIdle() {  // No requests in last 10 minutes
            delete(rl.buckets, id)
        }
    }
}

// Run periodically
go func() {
    for {
        time.Sleep(10 * time.Minute)
        rl.cleanup()
    }
}()
```

**Solution 2: LRU cache (sophisticated)**
```go
import "github.com/hashicorp/golang-lru"

type RateLimiter struct {
    buckets *lru.Cache  // Automatically evicts least recently used
}

func NewRateLimiter(maxClients int) *RateLimiter {
    cache, _ := lru.New(maxClients)
    return &RateLimiter{buckets: cache}
}
```

**Solution 3: Sliding window (memory-efficient)**
```go
type SlidingWindow struct {
    mu       sync.Mutex
    requests []time.Time  // Timestamps of recent requests
    limit    int
    window   time.Duration
}

func (sw *SlidingWindow) Allow() bool {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-sw.window)

    // Remove old requests
    i := 0
    for ; i < len(sw.requests); i++ {
        if sw.requests[i].After(cutoff) {
            break
        }
    }
    sw.requests = sw.requests[i:]

    if len(sw.requests) < sw.limit {
        sw.requests = append(sw.requests, now)
        return true
    }
    return false
}
```

### Concept 4: HTTP Status Codes for Rate Limiting

**429 Too Many Requests**: Standard status code for rate limiting.

**Important headers**:
```go
w.Header().Set("X-RateLimit-Limit", "100")        // Max requests
w.Header().Set("X-RateLimit-Remaining", "23")     // Remaining in window
w.Header().Set("X-RateLimit-Reset", "1614556800") // Unix timestamp
w.Header().Set("Retry-After", "60")               // Seconds until retry
```

**Example response**:
```
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1614556860
Retry-After: 60
Content-Type: application/json

{"error": "Rate limit exceeded. Try again in 60 seconds."}
```

**Client-side handling**:
```javascript
async function callAPI() {
    const response = await fetch('/api/data');

    if (response.status === 429) {
        const retryAfter = response.headers.get('Retry-After');
        console.log(`Rate limited. Retry after ${retryAfter} seconds`);

        // Exponential backoff
        await sleep(retryAfter * 1000);
        return callAPI();  // Retry
    }

    return response.json();
}
```

### Concept 5: Distributed Rate Limiting

**Problem**: Single-server rate limiter doesn't work with multiple servers.

```
Client → Load Balancer → Server 1 (allows 100 req/min)
                      → Server 2 (allows 100 req/min)
                      → Server 3 (allows 100 req/min)

Result: Client can make 300 req/min (bypasses limit!)
```

**Solution 1: Redis-based rate limiter**
```go
func (rl *RateLimiter) Allow(ctx context.Context, clientID string) bool {
    key := fmt.Sprintf("ratelimit:%s", clientID)

    // Lua script for atomic operations in Redis
    script := `
        local current = redis.call('GET', KEYS[1])
        if current and tonumber(current) >= tonumber(ARGV[1]) then
            return 0
        end
        redis.call('INCR', KEYS[1])
        redis.call('EXPIRE', KEYS[1], ARGV[2])
        return 1
    `

    result, err := rl.redis.Eval(ctx, script, []string{key}, rl.limit, rl.window).Int()
    return err == nil && result == 1
}
```

**Solution 2: Token bucket in Redis**
```go
// Store token bucket state in Redis
func (rl *RateLimiter) Allow(ctx context.Context, clientID string) bool {
    key := fmt.Sprintf("bucket:%s", clientID)

    // Get bucket state
    val, _ := rl.redis.Get(ctx, key).Result()
    bucket := parseBucket(val)

    bucket.refill()
    if bucket.tokens < 1 {
        return false
    }

    bucket.tokens--

    // Save bucket state
    rl.redis.Set(ctx, key, bucket.serialize(), time.Hour)
    return true
}
```

**Solution 3: Sticky sessions (simple but limited)**
```nginx
# Nginx configuration
upstream backend {
    ip_hash;  # Same client always goes to same server
    server backend1:8080;
    server backend2:8080;
}
```

---

## 6. Real-World Applications

### Public APIs

Every major API uses rate limiting.

**GitHub API**:
```
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4999
X-RateLimit-Reset: 1614556860

Authenticated: 5,000 requests/hour
Unauthenticated: 60 requests/hour
```

**OpenAI API**:
```
Rate Limit: 3,500 requests/minute
Token Limit: 90,000 tokens/minute

Uses token bucket for smooth handling
```

**Stripe API**:
```
Rate Limit: 100 requests/second (burst)
Sustained: ~25 requests/second

Returns 429 with Retry-After header
```

### DDoS Protection

**Cloudflare**: Rate limits at edge to block DDoS attacks.

```go
// Aggressive rate limiting for suspicious IPs
if isSuspicious(clientIP) {
    limiter = NewRateLimiter(10, 1)  // 10 requests/second
} else {
    limiter = NewRateLimiter(100, 10)  // 100 requests/second
}
```

### Preventing Brute Force

**Login endpoint**: Prevent password guessing.

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")

    // Rate limit per username (not IP, to prevent distributed attacks)
    if !loginRateLimiter.Allow(username) {
        http.Error(w, "Too many login attempts", 429)
        return
    }

    // Check password...
}
```

**Settings**:
- 5 attempts per minute per username
- After 5 failures, block for 15 minutes

### Cost Control

**Expensive operations**: Database queries, external API calls.

```go
// Limit expensive search queries
func searchHandler(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    // Allow 10 searches per minute
    if !searchRateLimiter.Allow(userID) {
        http.Error(w, "Search rate limit exceeded", 429)
        return
    }

    results := performExpensiveSearch(r.FormValue("query"))
    json.NewEncoder(w).Encode(results)
}
```

### Webhook Rate Limiting

**Receiving webhooks**: Prevent malicious webhook floods.

```go
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    signature := r.Header.Get("X-Webhook-Signature")

    // Verify signature first
    if !verifySignature(r.Body, signature) {
        http.Error(w, "Invalid signature", 401)
        return
    }

    // Rate limit per webhook source
    source := getWebhookSource(r)
    if !webhookRateLimiter.Allow(source) {
        http.Error(w, "Webhook rate limit exceeded", 429)
        return
    }

    processWebhook(r)
    w.WriteHeader(200)
}
```

---

## 7. Common Mistakes to Avoid

### Mistake 1: Forgetting Cleanup

**❌ Wrong**:
```go
type RateLimiter struct {
    buckets map[string]*TokenBucket
}

// Memory grows forever!
```

**✅ Correct**:
```go
type RateLimiter struct {
    buckets    map[string]*TokenBucket
    lastAccess map[string]time.Time
}

func (rl *RateLimiter) cleanup() {
    for id, lastSeen := range rl.lastAccess {
        if time.Since(lastSeen) > 1*time.Hour {
            delete(rl.buckets, id)
            delete(rl.lastAccess, id)
        }
    }
}
```

### Mistake 2: Using Plain int64 Instead of atomic.Int64

**❌ Wrong**:
```go
type TokenBucket struct {
    tokens int64  // Race condition!
}

func (b *TokenBucket) Allow() bool {
    if b.tokens < 1 {
        return false
    }
    b.tokens--  // NOT THREAD-SAFE
    return true
}
```

**✅ Correct**:
```go
type TokenBucket struct {
    tokens atomic.Int64  // Thread-safe
}

func (b *TokenBucket) Allow() bool {
    for {
        current := b.tokens.Load()
        if current < 1 {
            return false
        }
        if b.tokens.CompareAndSwap(current, current-1) {
            return true
        }
    }
}
```

### Mistake 3: Rate Limiting by IP Behind Load Balancer

**❌ Wrong**:
```go
func getClientIP(r *http.Request) string {
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip  // Returns load balancer IP, not client!
}
```

**✅ Correct**:
```go
func getClientIP(r *http.Request) string {
    // Check X-Forwarded-For header first
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        ips := strings.Split(xff, ",")
        return strings.TrimSpace(ips[0])  // First IP is client
    }

    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}
```

### Mistake 4: Inconsistent Refill Timing

**❌ Wrong**:
```go
// Refill every second in background goroutine
go func() {
    for {
        time.Sleep(1 * time.Second)
        b.tokens = b.capacity  // Resets completely!
    }
}()
```

**Problem**: Allows bursts every second, not smooth refill.

**✅ Correct**:
```go
// Refill on-demand, proportional to elapsed time
func (b *TokenBucket) refill() {
    elapsed := time.Since(b.lastRefill).Seconds()
    tokensToAdd := elapsed * b.rate
    b.tokens = min(b.tokens + tokensToAdd, b.capacity)
    b.lastRefill = time.Now()
}
```

### Mistake 5: Not Handling 429 on Client Side

**❌ Wrong**:
```javascript
fetch('/api/data')
    .then(res => res.json())
    .then(data => console.log(data))
    .catch(err => console.error(err));
// Doesn't check status code, treats 429 as error
```

**✅ Correct**:
```javascript
async function fetchWithRetry(url, maxRetries = 3) {
    for (let i = 0; i < maxRetries; i++) {
        const res = await fetch(url);

        if (res.status === 429) {
            const retryAfter = parseInt(res.headers.get('Retry-After') || '1');
            console.log(`Rate limited, retrying after ${retryAfter}s`);
            await sleep(retryAfter * 1000);
            continue;
        }

        if (res.ok) {
            return res.json();
        }

        throw new Error(`HTTP ${res.status}`);
    }
    throw new Error('Max retries exceeded');
}
```

### Mistake 6: Integer Overflow in Time Calculations

**❌ Wrong**:
```go
elapsed := int(time.Since(b.lastRefill).Seconds())  // Truncates!
tokensToAdd := elapsed * b.rate
```

**Problem**: For fractional rates (e.g., 0.5 tokens/second), int truncation means no refill.

**✅ Correct**:
```go
elapsed := time.Since(b.lastRefill).Seconds()  // float64
tokensToAdd := int64(elapsed * b.rate)
```

---

## 8. Stretch Goals

### Goal 1: Add Remaining Tokens Header ⭐

Return how many requests the client has left.

**Hint**:
```go
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        bucket := rl.getBucket(getClientIP(r))
        bucket.refill()

        remaining := bucket.tokens.Load()
        w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

        if !bucket.Allow() {
            http.Error(w, "Rate limit exceeded", 429)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Goal 2: Sliding Window Log ⭐⭐

Implement sliding window algorithm for comparison.

**Hint**:
```go
type SlidingWindow struct {
    mu       sync.Mutex
    requests []time.Time
    limit    int
    window   time.Duration
}

func (sw *SlidingWindow) Allow() bool {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-sw.window)

    // Remove expired requests
    i := 0
    for ; i < len(sw.requests); i++ {
        if sw.requests[i].After(cutoff) {
            break
        }
    }
    sw.requests = sw.requests[i:]

    if len(sw.requests) < sw.limit {
        sw.requests = append(sw.requests, now)
        return true
    }

    return false
}
```

### Goal 3: Redis-Based Distributed Rate Limiter ⭐⭐⭐

Share rate limit state across multiple servers.

**Hint**:
```go
import "github.com/go-redis/redis/v8"

type DistributedRateLimiter struct {
    redis  *redis.Client
    limit  int64
    window time.Duration
}

func (rl *DistributedRateLimiter) Allow(ctx context.Context, key string) bool {
    script := `
        local current = redis.call('incr', KEYS[1])
        if current == 1 then
            redis.call('expire', KEYS[1], ARGV[1])
        end
        return current <= tonumber(ARGV[2])
    `

    result, _ := rl.redis.Eval(ctx, script,
        []string{key},
        int(rl.window.Seconds()),
        rl.limit,
    ).Int()

    return result == 1
}
```

### Goal 4: Adaptive Rate Limiting ⭐⭐⭐

Adjust rate limits based on server load.

**Hint**:
```go
type AdaptiveRateLimiter struct {
    baseRate    float64
    currentRate atomic.Int64  // Adjusted dynamically
    load        func() float64  // Returns current CPU/memory usage
}

func (rl *AdaptiveRateLimiter) adjustRate() {
    load := rl.load()

    var rate float64
    switch {
    case load < 0.5:
        rate = rl.baseRate * 1.5  // Allow more requests
    case load < 0.7:
        rate = rl.baseRate
    case load < 0.9:
        rate = rl.baseRate * 0.5  // Reduce requests
    default:
        rate = rl.baseRate * 0.1  // Emergency throttle
    }

    rl.currentRate.Store(int64(rate))
}
```

### Goal 5: Per-Endpoint Rate Limits ⭐⭐

Different limits for different endpoints.

**Hint**:
```go
type EndpointRateLimiter struct {
    limiters map[string]*RateLimiter
}

func (erl *EndpointRateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        limiter := erl.getLimiterForPath(r.URL.Path)

        if !limiter.Allow(getClientIP(r)) {
            http.Error(w, "Rate limit exceeded", 429)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (erl *EndpointRateLimiter) getLimiterForPath(path string) *RateLimiter {
    switch {
    case strings.HasPrefix(path, "/api/search"):
        return erl.limiters["search"]  // 10 req/min
    case strings.HasPrefix(path, "/api/upload"):
        return erl.limiters["upload"]  // 5 req/min
    default:
        return erl.limiters["default"]  // 100 req/min
    }
}
```

---

## How to Run

```bash
# Run the demo server
cd /home/user/go-edu/minis/34-rate-limiter-token-bucket
go run ./cmd/rate-limiter

# In another terminal, test the rate limiter
# First 10 requests succeed
for i in {1..10}; do
    curl http://localhost:8080/api/data
    echo ""
done

# Next requests get rate limited
for i in {11..15}; do
    curl -i http://localhost:8080/api/data
    echo ""
done

# Run tests
cd exercise
go test -v

# Run with solution
go test -v -tags=solution
```

---

## Summary

**What you learned**:
- ✅ Rate limiting fundamentals and real-world importance
- ✅ Token bucket algorithm with mathematical model
- ✅ Atomic operations for lock-free programming
- ✅ CompareAndSwap for safe concurrent updates
- ✅ HTTP middleware pattern for rate limiting
- ✅ Per-client tracking with memory cleanup
- ✅ Proper HTTP status codes and headers (429, Retry-After)

**Why this matters**:
Every production API needs rate limiting. Without it, a single bad actor can take down your entire service. Token bucket is the industry standard algorithm used by AWS, Stripe, GitHub, and others.

**Key takeaway**:
Rate limiting = Protect your service + Fair resource allocation + Cost control

**Companies using this**:
- AWS API Gateway: Token bucket rate limiting
- Stripe: Token bucket for payment API
- GitHub: Tiered rate limits (5000/hour authenticated)
- Cloudflare: Rate limiting at edge for DDoS protection
- Twitter: Rate limits on all API endpoints

**Next steps**:
Apply this pattern to your own APIs. Start with generous limits, monitor usage, and adjust based on actual traffic patterns.

Build reliably!
