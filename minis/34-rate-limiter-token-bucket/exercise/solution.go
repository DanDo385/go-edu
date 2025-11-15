//go:build solution
// +build solution

/*
Problem: Implement a production-grade rate limiter using the token bucket algorithm

Requirements:
1. Thread-safe token bucket with atomic operations
2. Per-client rate limiting with independent buckets
3. HTTP middleware integration
4. Automatic token refill based on elapsed time
5. Proper client IP extraction (handle proxies/load balancers)

Why Go is well-suited:
- sync/atomic: Lock-free atomic operations for high performance
- Goroutines: Each HTTP request runs concurrently, needs thread safety
- sync.RWMutex: Efficient read-heavy access patterns
- http.Handler: Clean middleware composition

Compared to other languages:
- Python: GIL limits true concurrency, harder to implement lock-free algorithms
- Node.js: Single-threaded, easier but can't utilize multiple cores effectively
- Rust: More control, but more complex with ownership/borrowing
- Java: Similar capabilities, but more verbose

Token Bucket Algorithm:
- Bucket has maximum capacity (allows bursts)
- Tokens refill at constant rate (sustained rate limit)
- Each request costs 1 token
- If no tokens available, request is denied (429)

Real-world usage:
- AWS API Gateway: Token bucket rate limiting
- Stripe API: ~25 req/s sustained, 100 req/s burst
- GitHub API: 5000 req/hour authenticated
*/

package exercise

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// TokenBucket implements the token bucket algorithm for rate limiting.
type TokenBucket struct {
	capacity   int64        // Maximum number of tokens in the bucket
	tokens     atomic.Int64 // Current number of tokens (atomic for thread safety)
	rate       float64      // Tokens added per second
	lastRefill atomic.Int64 // Unix nanosecond timestamp of last refill (atomic)
}

// NewTokenBucket creates a new TokenBucket with the specified capacity and refill rate.
func NewTokenBucket(capacity int64, rate float64) *TokenBucket {
	tb := &TokenBucket{
		capacity: capacity,
		rate:     rate,
	}

	// Initialize with full bucket
	tb.tokens.Store(capacity)

	// Set initial refill time
	tb.lastRefill.Store(time.Now().UnixNano())

	return tb
}

// refill adds tokens to the bucket based on elapsed time since last refill.
func (b *TokenBucket) refill() {
	now := time.Now().UnixNano()
	last := b.lastRefill.Load()

	// Calculate elapsed time in seconds
	elapsed := float64(now-last) / float64(time.Second)

	// Calculate tokens to add
	tokensToAdd := int64(elapsed * b.rate)

	if tokensToAdd > 0 {
		// Try to update lastRefill timestamp
		// If another goroutine already did it, that's fine
		if b.lastRefill.CompareAndSwap(last, now) {
			// Successfully updated timestamp, now add tokens
			// Use CAS loop to handle concurrent modifications
			for {
				current := b.tokens.Load()
				new := current + tokensToAdd

				// Cap at capacity
				if new > b.capacity {
					new = b.capacity
				}

				// Try to update tokens
				if b.tokens.CompareAndSwap(current, new) {
					break
				}
				// If CAS failed, another goroutine modified tokens
				// Loop and retry with updated value
			}
		}
	}
}

// Allow attempts to consume one token from the bucket.
// Returns true if a token was available (request allowed).
// Returns false if no tokens available (request should be rate limited).
func (b *TokenBucket) Allow() bool {
	// First, refill tokens based on elapsed time
	b.refill()

	// Try to consume a token using CAS loop
	for {
		current := b.tokens.Load()

		// No tokens available
		if current < 1 {
			return false
		}

		// Try to consume one token
		if b.tokens.CompareAndSwap(current, current-1) {
			return true
		}

		// CAS failed, another goroutine consumed the token
		// Retry with updated value
	}
}

// RateLimiter manages rate limiting for multiple clients.
type RateLimiter struct {
	mu       sync.RWMutex            // Protects the buckets map
	buckets  map[string]*TokenBucket // Map of client ID to their token bucket
	capacity int64                   // Capacity for new buckets
	rate     float64                 // Rate for new buckets
}

// NewRateLimiter creates a new RateLimiter with specified capacity and rate.
func NewRateLimiter(capacity int64, rate float64) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}

// Allow checks if a request from the given client should be allowed.
func (rl *RateLimiter) Allow(clientID string) bool {
	bucket := rl.getBucket(clientID)
	return bucket.Allow()
}

// getBucket retrieves the TokenBucket for a client, creating one if it doesn't exist.
func (rl *RateLimiter) getBucket(clientID string) *TokenBucket {
	// Fast path: try to get bucket with read lock
	rl.mu.RLock()
	bucket, exists := rl.buckets[clientID]
	rl.mu.RUnlock()

	if exists {
		return bucket
	}

	// Slow path: create new bucket with write lock
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check: another goroutine might have created it
	// between releasing RLock and acquiring Lock
	if bucket, exists := rl.buckets[clientID]; exists {
		return bucket
	}

	// Create new bucket
	bucket = NewTokenBucket(rl.capacity, rl.rate)
	rl.buckets[clientID] = bucket

	return bucket
}

// Stats returns statistics about the rate limiter.
func (rl *RateLimiter) Stats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"total_clients": len(rl.buckets),
		"capacity":      rl.capacity,
		"rate":          rl.rate,
	}
}

// Middleware returns an HTTP middleware that applies rate limiting.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client IP
		clientIP := getClientIP(r)

		// Check if request is allowed
		if !rl.Allow(clientIP) {
			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", rl.rate*60))
			w.Header().Set("Retry-After", "1")

			// Return 429 Too Many Requests
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Request allowed, continue to next handler
		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the client's IP address from the request.
// Handles cases where the server is behind a proxy or load balancer.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Format: "client, proxy1, proxy2"
		// Take the first IP (actual client)
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header (alternative proxy header)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	// Format: "ip:port"
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// Cleanup removes inactive clients from the rate limiter to free memory.
// Should be called periodically (e.g., every 10 minutes) in production.
func (rl *RateLimiter) Cleanup(inactiveThreshold time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().UnixNano()

	for id, bucket := range rl.buckets {
		lastRefill := bucket.lastRefill.Load()
		elapsed := time.Duration(now - lastRefill)

		if elapsed > inactiveThreshold {
			delete(rl.buckets, id)
		}
	}
}
