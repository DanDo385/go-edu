//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "net/http" for HTTP handlers and middleware
// - "net" for IP address parsing
// - "strings" for string manipulation
// - "sync" for RWMutex (protecting the buckets map)
// - "sync/atomic" for atomic operations (lock-free counters)
// - "time" for time-based token refill
//
// import (
//     "net"
//     "net/http"
//     "strings"
//     "sync"
//     "sync/atomic"
//     "time"
// )

// ============================================================================
// TOKEN BUCKET RATE LIMITER: Lock-Free Concurrent Rate Limiting
// ============================================================================
//
// What is rate limiting?
// - Controls how many requests a client can make in a time window
// - Prevents abuse, DoS attacks, and resource exhaustion
// - Ensures fair resource allocation among clients
//
// Token Bucket Algorithm:
// 1. Each client gets a "bucket" that holds tokens
// 2. Bucket has a maximum capacity (allows bursts)
// 3. Tokens refill at a constant rate (sustained rate limit)
// 4. Each request consumes 1 token
// 5. If no tokens available → request denied (429)
//
// Example: capacity=100, rate=10/sec
// - Allows burst of 100 requests immediately
// - Then sustained rate of 10 requests/second
// - After idle period, bucket refills to 100 (burst capacity restored)
//
// Why this algorithm?
// - Allows burst traffic (better UX than fixed window)
// - Simple to implement and reason about
// - Used by AWS API Gateway, Stripe, GitHub
//
// Memory management:
// - TokenBucket uses atomic.Int64 for lock-free operations
// - RateLimiter uses sync.RWMutex for map protection
// - Each client gets their own TokenBucket (isolated rate limits)
//
// Concurrency:
// - Multiple goroutines can call Allow() concurrently
// - CompareAndSwap ensures only one goroutine consumes each token
// - RWMutex allows many readers, one writer (efficient for lookups)
//
// ============================================================================

// ============================================================================
// Exercise 1: Implement Token Bucket
// ============================================================================

// TokenBucket implements the token bucket algorithm for rate limiting.
// TODO: Define the TokenBucket struct with necessary fields.
//
// Required fields:
// - capacity: Maximum number of tokens (burst size)
// - tokens: Current number of tokens (use atomic.Int64 for thread safety)
// - rate: Tokens refilled per second (float64)
// - lastRefill: Unix nanosecond timestamp of last refill (use atomic.Int64)
//
// Memory layout:
// - capacity: 8 bytes (int64)
// - tokens: 8 bytes (atomic.Int64 wraps int64)
// - rate: 8 bytes (float64)
// - lastRefill: 8 bytes (atomic.Int64)
// - Total: ~32 bytes per bucket
//
// Why atomic.Int64?
// - Lock-free: No mutex needed for simple counter operations
// - Fast: Uses CPU atomic instructions (LOCK XADD, CMPXCHG)
// - Safe: Guarantees atomicity even with concurrent access
// - Must be 8-byte aligned (Go handles this automatically)
//
// type TokenBucket struct {
//     capacity   int64        // Maximum tokens (burst capacity)
//     tokens     atomic.Int64 // Current tokens (atomic for thread safety)
//     rate       float64      // Refill rate (tokens per second)
//     lastRefill atomic.Int64 // Last refill time (Unix nanoseconds)
// }

// NewTokenBucket creates a new TokenBucket with the specified capacity and refill rate.
// TODO: Implement this constructor.
//
// Parameters:
// - capacity: Maximum number of tokens (burst size)
// - rate: Tokens added per second (sustained rate limit)
//
// Memory allocation:
// - Returns a pointer (*TokenBucket) because:
//   1. Buckets are shared across goroutines (need same instance)
//   2. Atomic types must not be copied (copying breaks atomicity)
//   3. Allows RateLimiter to store pointers in map
//
// Initialization:
// - Start with a full bucket (capacity tokens available)
// - Set lastRefill to current time (time.Now().UnixNano())
// - Use .Store() method for atomic types
//
// Example usage:
//   bucket := NewTokenBucket(100, 10.0)  // 100 burst, 10/sec sustained
//
// func NewTokenBucket(capacity int64, rate float64) *TokenBucket {
//     tb := &TokenBucket{
//         capacity: capacity,
//         rate:     rate,
//     }
//     // Initialize with full bucket
//     tb.tokens.Store(capacity)
//     // Set initial refill time
//     tb.lastRefill.Store(time.Now().UnixNano())
//     return tb
// }

// refill adds tokens to the bucket based on elapsed time since last refill.
// TODO: Implement token refill logic.
//
// This is the core of the token bucket algorithm.
// Called before every Allow() check to ensure tokens are up-to-date.
//
// Algorithm:
// 1. Get current time as Unix nanoseconds
// 2. Load lastRefill timestamp (atomic)
// 3. Calculate elapsed time in seconds: (now - last) / 1e9
// 4. Calculate tokens to add: elapsed * rate
// 5. If tokensToAdd > 0:
//    a. Update lastRefill using CompareAndSwap (prevents race)
//    b. Add tokens using CompareAndSwap loop (caps at capacity)
//
// Why CompareAndSwap (CAS)?
// - Multiple goroutines might call refill() simultaneously
// - CAS ensures only one goroutine updates the timestamp
// - If CAS fails, another goroutine already refilled → skip
//
// CAS loop pattern:
//   for {
//       current := atomic.Load(&value)
//       new := calculateNew(current)
//       if atomic.CompareAndSwap(&value, current, new) {
//           break  // Success!
//       }
//       // CAS failed, another goroutine modified value
//       // Loop and retry with updated value
//   }
//
// Example:
//   Last refill: 2 seconds ago
//   Rate: 10 tokens/second
//   Tokens to add: 2 * 10 = 20 tokens
//
// func (b *TokenBucket) refill() {
//     now := time.Now().UnixNano()
//     last := b.lastRefill.Load()
//
//     // Calculate elapsed time in seconds
//     elapsed := float64(now-last) / float64(time.Second)
//
//     // Calculate tokens to add
//     tokensToAdd := int64(elapsed * b.rate)
//
//     if tokensToAdd > 0 {
//         // Try to update lastRefill timestamp
//         // If another goroutine already did it, that's fine
//         if b.lastRefill.CompareAndSwap(last, now) {
//             // Successfully updated timestamp, now add tokens
//             // Use CAS loop to handle concurrent modifications
//             for {
//                 current := b.tokens.Load()
//                 new := current + tokensToAdd
//
//                 // Cap at capacity (bucket can't overflow)
//                 if new > b.capacity {
//                     new = b.capacity
//                 }
//
//                 // Try to update tokens
//                 if b.tokens.CompareAndSwap(current, new) {
//                     break  // Success!
//                 }
//                 // If CAS failed, another goroutine modified tokens
//                 // Loop and retry with updated value
//             }
//         }
//     }
// }

// Allow attempts to consume one token from the bucket.
// TODO: Implement thread-safe token consumption.
//
// Returns true if a token was available (request allowed).
// Returns false if no tokens available (request should be rate limited).
//
// This method is thread-safe and can be called concurrently.
//
// Algorithm:
// 1. Call refill() to update token count based on elapsed time
// 2. Use CompareAndSwap loop to consume token:
//    a. Load current token count
//    b. If current < 1, return false (no tokens available)
//    c. Try to CompareAndSwap(current, current-1)
//    d. If CAS succeeds, return true (token consumed)
//    e. If CAS fails, retry (another goroutine took the token)
//
// Why CAS loop?
// - Multiple goroutines might try to consume tokens simultaneously
// - CAS ensures only one goroutine successfully consumes each token
// - Failed goroutines retry with updated value
//
// Example scenarios:
//   Scenario 1: 5 tokens available, 1 goroutine
//   - Load: current = 5
//   - CAS(5, 4) succeeds
//   - Return true
//
//   Scenario 2: 1 token available, 2 goroutines
//   - G1: Load current=1, CAS(1,0) succeeds → true
//   - G2: Load current=1, CAS(1,0) fails (G1 changed it)
//   - G2: Retry, Load current=0 → false (no tokens)
//
// func (b *TokenBucket) Allow() bool {
//     // First, refill tokens based on elapsed time
//     b.refill()
//
//     // Try to consume a token using CAS loop
//     for {
//         current := b.tokens.Load()
//
//         // No tokens available
//         if current < 1 {
//             return false
//         }
//
//         // Try to consume one token
//         if b.tokens.CompareAndSwap(current, current-1) {
//             return true  // Successfully consumed token
//         }
//
//         // CAS failed, another goroutine consumed the token
//         // Retry with updated value
//     }
// }

// ============================================================================
// Exercise 2: Implement Per-Client Rate Limiter
// ============================================================================

// RateLimiter manages rate limiting for multiple clients.
// TODO: Define the RateLimiter struct with necessary fields.
//
// Each client (identified by IP or user ID) gets their own TokenBucket.
//
// Required fields:
// - mu: sync.RWMutex to protect the buckets map
// - buckets: map[string]*TokenBucket (client ID → bucket)
// - capacity: Capacity for new buckets
// - rate: Rate for new buckets
//
// Memory considerations:
// - Map stores pointers to TokenBucket (not values)
// - Each client: ~32 bytes for TokenBucket + map overhead (~24 bytes)
// - Total per client: ~56 bytes
// - With 10,000 clients: ~560 KB memory
//
// Why RWMutex instead of Mutex?
// - Read-heavy workload (most requests are cache hits)
// - RLock allows multiple concurrent readers
// - Lock only needed when creating new bucket (cache miss)
// - Much better performance for high-traffic scenarios
//
// type RateLimiter struct {
//     mu       sync.RWMutex            // Protects the buckets map
//     buckets  map[string]*TokenBucket // Map of client ID to their token bucket
//     capacity int64                   // Capacity for new buckets
//     rate     float64                 // Rate for new buckets
// }

// NewRateLimiter creates a new RateLimiter with specified capacity and rate.
// TODO: Implement this constructor.
//
// All clients will share the same rate limit settings (capacity and rate).
//
// Initialization:
// - Create empty map with make(map[string]*TokenBucket)
// - Store capacity and rate for creating new buckets
//
// Memory:
// - Empty map: ~48 bytes overhead
// - Grows as clients are added
//
// func NewRateLimiter(capacity int64, rate float64) *RateLimiter {
//     return &RateLimiter{
//         buckets:  make(map[string]*TokenBucket),
//         capacity: capacity,
//         rate:     rate,
//     }
// }

// Allow checks if a request from the given client should be allowed.
// TODO: Implement client-specific rate limiting.
//
// clientID is typically an IP address or user identifier.
//
// Algorithm:
// 1. Get or create a TokenBucket for this client (call getBucket)
// 2. Call Allow() on the bucket
// 3. Return the result
//
// This is the main entry point for rate limiting.
//
// func (rl *RateLimiter) Allow(clientID string) bool {
//     bucket := rl.getBucket(clientID)
//     return bucket.Allow()
// }

// getBucket retrieves the TokenBucket for a client, creating one if it doesn't exist.
// TODO: Implement bucket retrieval/creation with double-checked locking.
//
// This uses the double-checked locking pattern for efficiency:
// 1. Fast path: RLock, check cache, RUnlock (most requests)
// 2. Slow path: Lock, double-check, create, Unlock (first request per client)
//
// Why double-check?
// - Between releasing RLock and acquiring Lock, another goroutine
//   might have created the bucket
// - Without double-check: multiple goroutines create duplicate buckets
// - With double-check: only first goroutine creates, others reuse
//
// Pattern:
//   Fast path (read lock):
//   1. Acquire RLock (allows concurrent reads)
//   2. Check if bucket exists in map
//   3. Release RLock
//   4. If exists, return bucket
//
//   Slow path (write lock):
//   1. Acquire Lock (exclusive access)
//   2. Defer Unlock
//   3. Double-check if bucket exists (another goroutine might have created it)
//   4. If exists, return it
//   5. Create new bucket with NewTokenBucket
//   6. Store in map
//   7. Return new bucket
//
// Example timeline (3 goroutines, same new client):
//   T0: G1, G2, G3 all acquire RLock
//   T1: G1, G2, G3 all see bucket doesn't exist
//   T2: G1, G2, G3 all release RLock
//   T3: G1 acquires Lock
//   T4: G1 double-checks (still doesn't exist), creates bucket
//   T5: G1 releases Lock
//   T6: G2 acquires Lock
//   T7: G2 double-checks (now exists!), returns existing bucket
//   T8: G2 releases Lock
//   T9: G3 acquires Lock, double-checks, returns existing bucket
//
// Without double-check: All 3 would create buckets, 2 would be lost
// With double-check: Only G1 creates, G2 and G3 reuse
//
// func (rl *RateLimiter) getBucket(clientID string) *TokenBucket {
//     // Fast path: try to get bucket with read lock
//     rl.mu.RLock()
//     bucket, exists := rl.buckets[clientID]
//     rl.mu.RUnlock()
//
//     if exists {
//         return bucket
//     }
//
//     // Slow path: create new bucket with write lock
//     rl.mu.Lock()
//     defer rl.mu.Unlock()
//
//     // Double-check: another goroutine might have created it
//     // between releasing RLock and acquiring Lock
//     if bucket, exists := rl.buckets[clientID]; exists {
//         return bucket
//     }
//
//     // Create new bucket
//     bucket = NewTokenBucket(rl.capacity, rl.rate)
//     rl.buckets[clientID] = bucket
//
//     return bucket
// }

// Stats returns statistics about the rate limiter.
// TODO: Implement stats collection.
//
// Useful for monitoring and debugging.
//
// Returns a map with:
// - "total_clients": Number of tracked clients
// - "capacity": Bucket capacity
// - "rate": Refill rate
//
// Thread safety:
// - Use RLock when reading the map size
// - Don't hold lock while building return value (copy size first)
//
// func (rl *RateLimiter) Stats() map[string]interface{} {
//     rl.mu.RLock()
//     totalClients := len(rl.buckets)
//     rl.mu.RUnlock()
//
//     return map[string]interface{}{
//         "total_clients": totalClients,
//         "capacity":      rl.capacity,
//         "rate":          rl.rate,
//     }
// }

// ============================================================================
// Exercise 3: Implement HTTP Middleware
// ============================================================================

// Middleware returns an HTTP middleware that applies rate limiting.
// TODO: Implement rate limiting middleware.
//
// Requests exceeding the rate limit receive a 429 Too Many Requests response.
//
// Algorithm:
// 1. Get client IP address (call getClientIP)
// 2. Check if request is allowed (call rl.Allow)
// 3. If not allowed:
//    a. Set rate limit headers (X-RateLimit-Limit, Retry-After)
//    b. Return 429 status code with error message
//    c. Return early (don't call next handler)
// 4. If allowed:
//    a. Call next.ServeHTTP(w, r) to continue request processing
//
// HTTP status codes:
// - 429 Too Many Requests: Rate limit exceeded
//
// Headers:
// - X-RateLimit-Limit: Maximum requests per time window
// - Retry-After: Seconds to wait before retrying
//
// Why middleware pattern?
// - Separates rate limiting from business logic
// - Easy to apply to multiple routes
// - Can be composed with other middleware (auth, logging)
//
// func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Extract client IP
//         clientIP := getClientIP(r)
//
//         // Check if request is allowed
//         if !rl.Allow(clientIP) {
//             // Set rate limit headers
//             w.Header().Set("X-RateLimit-Limit", "100")  // Informational
//             w.Header().Set("Retry-After", "1")          // Wait 1 second
//
//             // Return 429 Too Many Requests
//             http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
//             return
//         }
//
//         // Request allowed, continue to next handler
//         next.ServeHTTP(w, r)
//     })
// }

// getClientIP extracts the client's IP address from the request.
// TODO: Implement IP extraction with proxy support.
//
// Handles cases where the server is behind a proxy or load balancer.
//
// Priority order:
// 1. Check X-Forwarded-For header (set by proxies/load balancers)
//    - Format: "client, proxy1, proxy2"
//    - Take first IP (actual client)
//    - Split by comma, trim whitespace
// 2. Check X-Real-IP header (alternative proxy header)
// 3. Fall back to r.RemoteAddr
//    - Format: "ip:port"
//    - Use net.SplitHostPort to extract IP
//
// Why this matters:
// - Behind a load balancer, r.RemoteAddr is the load balancer's IP
// - Without checking headers, all clients appear to be the same IP
// - This would rate limit ALL clients together (not per-client)
//
// Security considerations:
// - X-Forwarded-For can be spoofed by malicious clients
// - In production, only trust if request comes from known proxy
// - Can configure load balancer to set X-Real-IP (harder to spoof)
//
// func getClientIP(r *http.Request) string {
//     // Check X-Forwarded-For header (set by proxies/load balancers)
//     if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
//         // Format: "client, proxy1, proxy2"
//         // Take the first IP (actual client)
//         ips := strings.Split(xff, ",")
//         return strings.TrimSpace(ips[0])
//     }
//
//     // Check X-Real-IP header (alternative proxy header)
//     if xri := r.Header.Get("X-Real-IP"); xri != "" {
//         return xri
//     }
//
//     // Fall back to RemoteAddr
//     // Format: "ip:port"
//     ip, _, err := net.SplitHostPort(r.RemoteAddr)
//     if err != nil {
//         return r.RemoteAddr  // Return as-is if parsing fails
//     }
//     return ip
// }

// ============================================================================
// Exercise 4: Memory Management - Cleanup
// ============================================================================

// Cleanup removes inactive clients from the rate limiter to free memory.
// TODO: Implement cleanup logic (STRETCH GOAL).
//
// A client is considered inactive if they haven't made a request recently.
// This should be called periodically (e.g., every 10 minutes) in production.
//
// Without cleanup:
// - Memory grows unbounded as new clients make requests
// - If 1M unique IPs make requests, you'll store 1M buckets
// - Each bucket: ~56 bytes → ~56 MB memory
// - Most clients are one-time visitors, don't need indefinite tracking
//
// With cleanup:
// - Remove buckets for clients inactive > threshold
// - Bounded memory usage (only active clients tracked)
// - Example: 10-minute threshold → only last 10 min of clients in memory
//
// Algorithm:
// 1. Acquire write lock (modifying map)
// 2. Iterate through all buckets
// 3. Check each bucket's lastRefill time
// 4. If time.Since(lastRefill) > inactiveThreshold, delete bucket
// 5. Release lock
//
// Alternative approaches:
// - Use an LRU cache (github.com/hashicorp/golang-lru)
//   * Automatically evicts least recently used clients
//   * Bounded memory usage without manual cleanup
//   * More complex but handles edge cases better
//
// - Use time-based expiration (github.com/patrickmn/go-cache)
//   * Each entry has TTL, automatically cleaned up
//   * Efficient background cleanup
//
// func (rl *RateLimiter) Cleanup(inactiveThreshold time.Duration) {
//     rl.mu.Lock()
//     defer rl.mu.Unlock()
//
//     now := time.Now().UnixNano()
//
//     for id, bucket := range rl.buckets {
//         lastRefill := bucket.lastRefill.Load()
//         elapsed := time.Duration(now - lastRefill)
//
//         if elapsed > inactiveThreshold {
//             delete(rl.buckets, id)
//         }
//     }
// }

// ============================================================================
// Testing and Running
// ============================================================================
//
// After implementing all functions:
// - Run: go test -v ./...
// - Run with race detector: go test -race ./...
// - Check: All tests should pass
//
// To run the server:
// - See cmd/server/main.go for example usage
// - Start server: go run cmd/server/main.go
// - Test with curl:
//   curl http://localhost:8080/
//   curl http://localhost:8080/  # Repeat to hit rate limit
//
// To test rate limiting:
// - Use Apache Bench: ab -n 200 -c 10 http://localhost:8080/
// - Use hey: hey -n 200 -c 10 http://localhost:8080/
// - Check stats: curl http://localhost:8080/stats
//
// Key concepts learned:
// 1. Lock-free programming with atomic.CompareAndSwap
// 2. Double-checked locking pattern
// 3. RWMutex for read-heavy workloads
// 4. Token bucket algorithm
// 5. HTTP middleware pattern
// 6. Client IP extraction with proxy support
// 7. Memory management and cleanup strategies
//
// Compare with solution.go to see detailed implementations and explanations.
