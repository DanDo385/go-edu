//go:build solution
// +build solution

package exercise

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"sync/atomic"
	"time"
)

// =============================================================================
// Context Keys and Helpers
// =============================================================================

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userKey      contextKey = "user"
)

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUser retrieves the user from the context
func GetUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey).(*User)
	return user, ok
}

// =============================================================================
// Models
// =============================================================================

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// =============================================================================
// Response Writer Wrapper
// =============================================================================

// ResponseWriter wraps http.ResponseWriter to capture status code and bytes written
type ResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	bytesWritten  int
	headerWritten bool
}

// NewResponseWriter creates a new ResponseWriter wrapper
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}

// WriteHeader captures the status code before writing it
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	if !rw.headerWritten {
		rw.statusCode = statusCode
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.headerWritten = true
	}
}

// Write counts bytes written and ensures header is written
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		rw.WriteHeader(http.StatusOK)
	}

	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// StatusCode returns the captured status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// BytesWritten returns the total bytes written
func (rw *ResponseWriter) BytesWritten() int {
	return rw.bytesWritten
}

// =============================================================================
// Middleware Type
// =============================================================================

type Middleware func(http.Handler) http.Handler

// =============================================================================
// Middleware Implementations
// =============================================================================

// LoggingMiddleware logs request and response details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)

		// Get request ID from context if available
		requestID, _ := GetRequestID(r.Context())
		if requestID != "" {
			log.Printf("[%s] → %s %s", requestID, r.Method, r.URL.Path)
		} else {
			log.Printf("→ %s %s", r.Method, r.URL.Path)
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		if requestID != "" {
			log.Printf(
				"[%s] ← %s %s - Status: %d, Bytes: %d, Duration: %v",
				requestID,
				r.Method,
				r.URL.Path,
				rw.StatusCode(),
				rw.BytesWritten(),
				duration,
			)
		} else {
			log.Printf(
				"← %s %s - Status: %d, Bytes: %d, Duration: %v",
				r.Method,
				r.URL.Path,
				rw.StatusCode(),
				rw.BytesWritten(),
				duration,
			)
		}
	})
}

// RecoveryMiddleware catches panics and returns 500 instead of crashing
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := GetRequestID(r.Context())
				if requestID != "" {
					log.Printf("[%s] PANIC: %v\n%s", requestID, err, debug.Stack())
				} else {
					log.Printf("PANIC: %v\n%s", err, debug.Stack())
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RequestIDMiddleware assigns a unique ID to each request
var requestCounter uint64

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists in header
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate simple incremental ID for testing
			// In production, use UUID or similar
			counter := atomic.AddUint64(&requestCounter, 1)
			requestID = "req-" + formatUint64(counter)
		}

		// Add to context
		ctx := WithRequestID(r.Context(), requestID)

		// Add to response header
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to format uint64 as string
func formatUint64(n uint64) string {
	if n == 0 {
		return "0"
	}

	var buf [20]byte // enough for 64-bit int
	i := len(buf) - 1

	for n > 0 {
		buf[i] = byte('0' + n%10)
		n /= 10
		i--
	}

	return string(buf[i+1:])
}

// AuthMiddleware validates authorization and adds user to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// Simple token validation for exercise
		// In production: verify JWT, check database, etc.
		if token != "Bearer valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Create user (in production: extract from token)
		user := &User{ID: 1, Name: "Alice"}

		// Add user to context
		ctx := WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(allowOrigin string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// MethodMiddleware only allows specific HTTP methods
func MethodMiddleware(allowedMethods ...string) Middleware {
	// Build map for O(1) lookup
	allowed := make(map[string]bool)
	for _, method := range allowedMethods {
		allowed[method] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allowed[r.Method] {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// =============================================================================
// Chain Helper
// =============================================================================

// Chain applies middleware in order: first middleware wraps all others
//
// Example: Chain(handler, A, B, C) produces A(B(C(handler)))
// Execution order: A → B → C → handler
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middleware in reverse order so first middleware is outermost
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
