package exercise

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestTokenBucket_SingleRequest tests basic token consumption
func TestTokenBucket_SingleRequest(t *testing.T) {
	bucket := NewTokenBucket(10, 5.0)
	if bucket == nil {
		t.Fatal("NewTokenBucket returned nil")
	}

	// First request should be allowed
	if !bucket.Allow() {
		t.Error("First request should be allowed")
	}

	// Verify token was consumed
	tokens := bucket.tokens.Load()
	if tokens != 9 {
		t.Errorf("Expected 9 tokens after one request, got %d", tokens)
	}
}

// TestTokenBucket_Capacity tests that bucket respects capacity limit
func TestTokenBucket_Capacity(t *testing.T) {
	capacity := int64(5)
	bucket := NewTokenBucket(capacity, 10.0)
	if bucket == nil {
		t.Fatal("NewTokenBucket returned nil")
	}

	// Consume all tokens
	for i := int64(0); i < capacity; i++ {
		if !bucket.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be denied
	if bucket.Allow() {
		t.Error("Request beyond capacity should be denied")
	}

	// Verify bucket is empty
	tokens := bucket.tokens.Load()
	if tokens != 0 {
		t.Errorf("Expected 0 tokens after consuming capacity, got %d", tokens)
	}
}

// TestTokenBucket_Refill tests token refill over time
func TestTokenBucket_Refill(t *testing.T) {
	// Create bucket: capacity=10, rate=5 tokens/second
	bucket := NewTokenBucket(10, 5.0)
	if bucket == nil {
		t.Fatal("NewTokenBucket returned nil")
	}

	// Consume all tokens
	for i := 0; i < 10; i++ {
		bucket.Allow()
	}

	// Verify empty
	if bucket.Allow() {
		t.Error("Bucket should be empty")
	}

	// Wait 1 second (should refill 5 tokens)
	time.Sleep(1 * time.Second)

	// Should be able to make 5 requests
	allowed := 0
	for i := 0; i < 10; i++ {
		if bucket.Allow() {
			allowed++
		}
	}

	// Allow some tolerance (4-6 tokens due to timing)
	if allowed < 4 || allowed > 6 {
		t.Errorf("Expected ~5 requests after 1 second, got %d", allowed)
	}
}

// TestTokenBucket_RefillCap tests that refill doesn't exceed capacity
func TestTokenBucket_RefillCap(t *testing.T) {
	bucket := NewTokenBucket(10, 5.0)
	if bucket == nil {
		t.Fatal("NewTokenBucket returned nil")
	}

	// Consume 3 tokens
	for i := 0; i < 3; i++ {
		bucket.Allow()
	}

	// Wait long enough to refill beyond capacity
	time.Sleep(3 * time.Second)

	// Should only have capacity worth of tokens
	allowed := 0
	for i := 0; i < 15; i++ {
		if bucket.Allow() {
			allowed++
		}
	}

	if allowed != 10 {
		t.Errorf("Expected 10 tokens (capacity), got %d", allowed)
	}
}

// TestTokenBucket_Concurrent tests thread safety
func TestTokenBucket_Concurrent(t *testing.T) {
	bucket := NewTokenBucket(100, 50.0)
	if bucket == nil {
		t.Fatal("NewTokenBucket returned nil")
	}

	var wg sync.WaitGroup
	allowed := 0
	var mu sync.Mutex

	// 10 goroutines each trying 20 requests
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				if bucket.Allow() {
					mu.Lock()
					allowed++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	// Should have consumed all 100 tokens
	if allowed != 100 {
		t.Errorf("Expected 100 allowed requests, got %d", allowed)
	}

	// Verify bucket is empty
	if bucket.Allow() {
		t.Error("Bucket should be empty after consuming all tokens")
	}
}

// TestRateLimiter_SingleClient tests rate limiting for one client
func TestRateLimiter_SingleClient(t *testing.T) {
	limiter := NewRateLimiter(10, 5.0)
	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	clientID := "192.168.1.1"

	// First request should be allowed
	if !limiter.Allow(clientID) {
		t.Error("First request should be allowed")
	}

	// Should be able to make 9 more requests (total capacity = 10)
	for i := 0; i < 9; i++ {
		if !limiter.Allow(clientID) {
			t.Errorf("Request %d should be allowed", i+2)
		}
	}

	// Next request should be denied
	if limiter.Allow(clientID) {
		t.Error("Request beyond capacity should be denied")
	}
}

// TestRateLimiter_MultipleClients tests independent rate limiting per client
func TestRateLimiter_MultipleClients(t *testing.T) {
	limiter := NewRateLimiter(5, 10.0)

	client1 := "192.168.1.1"
	client2 := "192.168.1.2"

	// Client 1 consumes all tokens
	for i := 0; i < 5; i++ {
		if !limiter.Allow(client1) {
			t.Errorf("Client 1 request %d should be allowed", i+1)
		}
	}

	// Client 1 should be rate limited
	if limiter.Allow(client1) {
		t.Error("Client 1 should be rate limited")
	}

	// Client 2 should still have full capacity
	for i := 0; i < 5; i++ {
		if !limiter.Allow(client2) {
			t.Errorf("Client 2 request %d should be allowed", i+1)
		}
	}

	// Client 2 should now be rate limited
	if limiter.Allow(client2) {
		t.Error("Client 2 should be rate limited")
	}
}

// TestRateLimiter_Stats tests statistics reporting
func TestRateLimiter_Stats(t *testing.T) {
	limiter := NewRateLimiter(10, 5.0)

	// Make requests from 3 different clients
	limiter.Allow("192.168.1.1")
	limiter.Allow("192.168.1.2")
	limiter.Allow("192.168.1.3")

	stats := limiter.Stats()
	if stats == nil {
		t.Fatal("Stats returned nil")
	}

	// Check total_clients
	if clients, ok := stats["total_clients"].(int); !ok || clients != 3 {
		t.Errorf("Expected 3 clients, got %v", stats["total_clients"])
	}

	// Check capacity
	if capacity, ok := stats["capacity"].(int64); !ok || capacity != 10 {
		t.Errorf("Expected capacity 10, got %v", stats["capacity"])
	}

	// Check rate
	if rate, ok := stats["rate"].(float64); !ok || rate != 5.0 {
		t.Errorf("Expected rate 5.0, got %v", stats["rate"])
	}
}

// TestMiddleware_AllowRequest tests middleware allows valid requests
func TestMiddleware_AllowRequest(t *testing.T) {
	limiter := NewRateLimiter(10, 5.0)

	// Create a simple handler that returns 200 OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with rate limiting middleware
	middleware := limiter.Middleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	// Make request
	middleware.ServeHTTP(rec, req)

	// Should be allowed (200 OK)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "success" {
		t.Errorf("Expected 'success', got '%s'", rec.Body.String())
	}
}

// TestMiddleware_RateLimitRequest tests middleware blocks rate limited requests
func TestMiddleware_RateLimitRequest(t *testing.T) {
	limiter := NewRateLimiter(3, 1.0) // Small capacity for easy testing

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := limiter.Middleware(handler)

	// Make requests until rate limited
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/data", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d should be allowed, got status %d", i+1, rec.Code)
		}
	}

	// Next request should be rate limited (429)
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", rec.Code)
	}

	// Check for rate limit headers
	if rec.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("Expected X-RateLimit-Limit header")
	}

	if rec.Header().Get("Retry-After") == "" {
		t.Error("Expected Retry-After header")
	}
}

// TestMiddleware_XForwardedFor tests IP extraction from X-Forwarded-For header
func TestMiddleware_XForwardedFor(t *testing.T) {
	limiter := NewRateLimiter(2, 1.0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := limiter.Middleware(handler)

	// Request from client 1 via proxy
	req1 := httptest.NewRequest("GET", "/api/data", nil)
	req1.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")
	rec1 := httptest.NewRecorder()
	middleware.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Error("First request should be allowed")
	}

	// Another request from same client (different proxy)
	req2 := httptest.NewRequest("GET", "/api/data", nil)
	req2.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.3")
	rec2 := httptest.NewRecorder()
	middleware.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Error("Second request from same client should be allowed")
	}

	// Third request should be rate limited (capacity = 2)
	req3 := httptest.NewRequest("GET", "/api/data", nil)
	req3.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.4")
	rec3 := httptest.NewRecorder()
	middleware.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusTooManyRequests {
		t.Errorf("Third request should be rate limited, got status %d", rec3.Code)
	}
}

// TestGetClientIP tests IP address extraction logic
func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		forwardedFor   string
		realIP         string
		expectedIP     string
	}{
		{
			name:       "RemoteAddr only",
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "192.168.1.1",
		},
		{
			name:         "X-Forwarded-For single IP",
			remoteAddr:   "10.0.0.1:12345",
			forwardedFor: "192.168.1.1",
			expectedIP:   "192.168.1.1",
		},
		{
			name:         "X-Forwarded-For multiple IPs",
			remoteAddr:   "10.0.0.1:12345",
			forwardedFor: "192.168.1.1, 10.0.0.2, 10.0.0.3",
			expectedIP:   "192.168.1.1",
		},
		{
			name:       "X-Real-IP",
			remoteAddr: "10.0.0.1:12345",
			realIP:     "192.168.1.1",
			expectedIP: "192.168.1.1",
		},
		{
			name:         "X-Forwarded-For takes precedence",
			remoteAddr:   "10.0.0.1:12345",
			forwardedFor: "192.168.1.1",
			realIP:       "192.168.1.2",
			expectedIP:   "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr

			if tt.forwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.forwardedFor)
			}

			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			ip := getClientIP(req)

			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}

// BenchmarkTokenBucket_Allow benchmarks token consumption
func BenchmarkTokenBucket_Allow(b *testing.B) {
	bucket := NewTokenBucket(int64(b.N), float64(b.N))
	if bucket == nil {
		b.Fatal("NewTokenBucket returned nil")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bucket.Allow()
	}
}

// BenchmarkTokenBucket_Concurrent benchmarks concurrent access
func BenchmarkTokenBucket_Concurrent(b *testing.B) {
	bucket := NewTokenBucket(int64(b.N), float64(b.N))
	if bucket == nil {
		b.Fatal("NewTokenBucket returned nil")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bucket.Allow()
		}
	})
}

// BenchmarkRateLimiter_Allow benchmarks rate limiter with single client
func BenchmarkRateLimiter_Allow(b *testing.B) {
	limiter := NewRateLimiter(int64(b.N), float64(b.N))
	if limiter == nil {
		b.Fatal("NewRateLimiter returned nil")
	}
	clientID := "192.168.1.1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(clientID)
	}
}

// BenchmarkRateLimiter_MultiClient benchmarks rate limiter with many clients
func BenchmarkRateLimiter_MultiClient(b *testing.B) {
	limiter := NewRateLimiter(1000, 100.0)
	if limiter == nil {
		b.Fatal("NewRateLimiter returned nil")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			clientID := fmt.Sprintf("192.168.1.%d", i%256)
			limiter.Allow(clientID)
			i++
		}
	})
}

// BenchmarkMiddleware benchmarks HTTP middleware
func BenchmarkMiddleware(b *testing.B) {
	limiter := NewRateLimiter(int64(b.N), float64(b.N))
	if limiter == nil {
		b.Fatal("NewRateLimiter returned nil")
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := limiter.Middleware(handler)

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		middleware.ServeHTTP(rec, req)
	}
}
