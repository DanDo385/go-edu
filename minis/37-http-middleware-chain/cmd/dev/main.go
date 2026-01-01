package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// =============================================================================
// Context Keys
// =============================================================================

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userKey      contextKey = "user"
)

// Helper functions for type-safe context access
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

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

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	if !rw.headerWritten {
		rw.statusCode = statusCode
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.headerWritten = true
	}
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		rw.WriteHeader(http.StatusOK)
	}

	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) BytesWritten() int {
	return rw.bytesWritten
}

// Support http.Flusher interface
func (rw *ResponseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// =============================================================================
// Middleware Type
// =============================================================================

type Middleware func(http.Handler) http.Handler

// =============================================================================
// Middleware: Recovery
// =============================================================================

// RecoveryMiddleware catches panics and returns 500 instead of crashing
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := GetRequestID(r.Context())
				log.Printf("[%s] PANIC: %v\n%s", requestID, err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// =============================================================================
// Middleware: Request ID
// =============================================================================

// RequestIDMiddleware assigns a unique ID to each request
var requestIDCounter uint64

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists (from header)
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Add to context
		ctx := WithRequestID(r.Context(), requestID)

		// Add to response header
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Use atomic counter + random bytes for uniqueness
	counter := atomic.AddUint64(&requestIDCounter, 1)

	// Generate 8 random bytes
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to counter only if random generation fails
		return fmt.Sprintf("req-%d", counter)
	}

	return fmt.Sprintf("req-%d-%s", counter, hex.EncodeToString(randomBytes)[:8])
}

// =============================================================================
// Middleware: Logging
// =============================================================================

// LoggingMiddleware logs request and response details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)

		requestID, _ := GetRequestID(r.Context())
		log.Printf("[%s] → %s %s %s", requestID, r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"[%s] ← %s %s - Status: %d, Bytes: %d, Duration: %v",
			requestID,
			r.Method,
			r.URL.Path,
			rw.StatusCode(),
			rw.BytesWritten(),
			duration,
		)
	})
}

// =============================================================================
// Middleware: Authentication
// =============================================================================

// AuthMiddleware validates authorization token and adds user to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// Simple token validation (in production: verify JWT, check DB, etc.)
		user, err := validateToken(token)
		if err != nil {
			requestID, _ := GetRequestID(r.Context())
			log.Printf("[%s] Auth failed: %v", requestID, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(token string) (*User, error) {
	// In production: verify JWT signature, check expiration, query database
	if token == "" || token != "Bearer secret-token" {
		return nil, fmt.Errorf("invalid token")
	}

	return &User{ID: 1, Name: "Alice"}, nil
}

// =============================================================================
// Middleware: CORS
// =============================================================================

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

// =============================================================================
// Middleware: Rate Limiting
// =============================================================================

type RequestCounter struct {
	mu    sync.Mutex
	count int
}

func (rc *RequestCounter) Increment() int {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.count++
	return rc.count
}

func (rc *RequestCounter) Get() int {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return rc.count
}

// RequestCounterMiddleware adds request count header
func RequestCounterMiddleware(counter *RequestCounter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := counter.Increment()
			w.Header().Set("X-Request-Count", fmt.Sprintf("%d", count))
			next.ServeHTTP(w, r)
		})
	}
}

// =============================================================================
// Middleware: Timeout
// =============================================================================

// TimeoutMiddleware enforces request timeout
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			done := make(chan struct{})
			panicChan := make(chan interface{}, 1)

			go func() {
				defer func() {
					if p := recover(); p != nil {
						panicChan <- p
					}
				}()

				next.ServeHTTP(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case p := <-panicChan:
				panic(p)
			case <-done:
				return
			case <-ctx.Done():
				requestID, _ := GetRequestID(r.Context())
				log.Printf("[%s] Request timeout after %v", requestID, timeout)
				http.Error(w, "Request Timeout", http.StatusGatewayTimeout)
				return
			}
		})
	}
}

// =============================================================================
// Middleware: Method Filter
// =============================================================================

// MethodMiddleware only allows specific HTTP methods
func MethodMiddleware(allowedMethods ...string) Middleware {
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
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply in reverse so first middleware is outermost
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// =============================================================================
// Handlers
// =============================================================================

func handleHome(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Welcome to the HTTP Middleware Demo!",
		"path":    r.URL.Path,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUser(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handlePanic(w http.ResponseWriter, r *http.Request) {
	panic("This is a deliberate panic to test RecoveryMiddleware!")
}

func handleSlow(w http.ResponseWriter, r *http.Request) {
	// Simulate slow processing
	select {
	case <-time.After(10 * time.Second):
		w.Write([]byte("This should not be reached due to timeout"))
	case <-r.Context().Done():
		// Context cancelled (timeout)
		return
	}
}

func handleStats(counter *RequestCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := map[string]interface{}{
			"total_requests": counter.Get(),
			"timestamp":      time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

// =============================================================================
// Main
// =============================================================================

func main() {
	// Request counter for stats
	counter := &RequestCounter{}

	// Create router
	mux := http.NewServeMux()

	// Public routes (no auth)
	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/", handleHome)
	publicMux.HandleFunc("/users", handleUsers)
	publicMux.HandleFunc("/panic", handlePanic)
	publicMux.HandleFunc("/stats", handleStats(counter))

	publicHandler := Chain(
		publicMux,
		RecoveryMiddleware,
		RequestIDMiddleware,
		LoggingMiddleware,
		RequestCounterMiddleware(counter),
		CORSMiddleware("*"),
	)

	// Protected routes (auth required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/profile", handleProfile)

	protectedHandler := Chain(
		protectedMux,
		RecoveryMiddleware,
		RequestIDMiddleware,
		LoggingMiddleware,
		RequestCounterMiddleware(counter),
		AuthMiddleware,
	)

	// Slow route with timeout
	slowMux := http.NewServeMux()
	slowMux.HandleFunc("/slow", handleSlow)

	slowHandler := Chain(
		slowMux,
		RecoveryMiddleware,
		RequestIDMiddleware,
		LoggingMiddleware,
		TimeoutMiddleware(3*time.Second),
	)

	// Mount all handlers
	mux.Handle("/profile", protectedHandler)
	mux.Handle("/slow", slowHandler)
	mux.Handle("/", publicHandler)

	// Create server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		log.Println("Server starting on :8080")
		log.Println("")
		log.Println("Try these endpoints:")
		log.Println("  curl http://localhost:8080/")
		log.Println("  curl http://localhost:8080/users")
		log.Println("  curl http://localhost:8080/stats")
		log.Println("  curl http://localhost:8080/panic")
		log.Println("  curl http://localhost:8080/slow")
		log.Println("  curl -H 'Authorization: Bearer secret-token' http://localhost:8080/profile")
		log.Println("")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("\nShutting down gracefully...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
