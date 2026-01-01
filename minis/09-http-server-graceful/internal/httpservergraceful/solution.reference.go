//go:build reference

/*
Problem: Build an HTTP server with routes, middleware, and graceful shutdown

Requirements:
1. REST endpoints for key-value storage
2. Request counting middleware
3. Graceful shutdown on SIGINT/SIGTERM
4. JSON request/response handling

Algorithm: HTTP Request Handling
- Route requests to appropriate handlers
- Apply middleware for cross-cutting concerns
- Handle JSON encoding/decoding
- Coordinate graceful shutdown

Graceful Shutdown Algorithm:
- Receive shutdown signal (SIGINT/SIGTERM)
- Stop accepting new connections
- Wait for in-flight requests to complete
- Close server cleanly

Middleware Pattern:
- Wrap handler functions
- Execute before/after main handler
- Common uses: logging, metrics, authentication
*/

package httpservergraceful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

// Store defines key-value operations.
// BREAKPOINT: Set breakpoint in methods to trace storage operations
type Store interface {
	Put(key, val string) error
	Get(key string) (string, bool)
}

var reqCount atomic.Int64

// RegisterRoutes sets up HTTP handlers.
//
// Algorithm:
// 1. Create handler with store dependency
// 2. Wrap handler with middleware
// 3. Register wrapped handler with mux
//
// BREAKPOINT: Set breakpoint here to trace route registration
// DEBUG: Watch 's' to see store implementation
func RegisterRoutes(mux *http.ServeMux, s Store) {
	// DEBUG: Creating handler with store dependency
	// BREAKPOINT: Set breakpoint here to see handler creation
	mux.HandleFunc("/kv", withMiddleware(kvHandler(s)))
}

// withMiddleware wraps a handler with request counting.
//
// Middleware Pattern:
// 1. Execute before handler (increment counter, set header)
// 2. Call wrapped handler
// 3. Execute after handler (if needed)
//
// BREAKPOINT: Set breakpoint here to trace middleware execution
// DEBUG: Watch 'count' to see request counter
func withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// BREAKPOINT: Set breakpoint here to trace request counting
		// DEBUG: Watch reqCount before and after increment
		count := reqCount.Add(1)
		// DEBUG: Current request number

		// BREAKPOINT: Set breakpoint here to see header being set
		// DEBUG: Watch 'count' to see what value is set
		w.Header().Set("X-Req-Count", fmt.Sprintf("%d", count))

		// DEBUG: Calling next handler
		// BREAKPOINT: Set breakpoint here before calling wrapped handler
		next(w, r)
		// DEBUG: Handler completed
	}
}

// kvHandler creates a handler for key-value operations.
//
// Algorithm:
// - POST: Decode JSON body, store key-value pair
// - GET: Extract key from query params, return value
//
// BREAKPOINT: Set breakpoint in returned function to trace requests
// DEBUG: Watch 'r.Method' to see HTTP method
func kvHandler(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// BREAKPOINT: Set breakpoint here to trace request routing
		// DEBUG: Watch 'r.Method' to see which case executes

		switch r.Method {
		case http.MethodPost:
			// BREAKPOINT: Set breakpoint here to trace POST requests
			// DEBUG: Watch 'req' to see decoded JSON
			var req struct {
				Key string `json:"key"`
				Val string `json:"val"`
			}

			// BREAKPOINT: Set breakpoint here to trace JSON decoding
			// DEBUG: Watch 'r.Body' before decode
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				// DEBUG: JSON decode failed
				// BREAKPOINT: Set breakpoint here to inspect decode error
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// DEBUG: Watch 'req.Key' and 'req.Val' to see decoded values

			// BREAKPOINT: Set breakpoint here to trace storage
			// DEBUG: Watch 's.Put' arguments
			if err := s.Put(req.Key, req.Val); err != nil {
				// DEBUG: Storage failed
				// BREAKPOINT: Set breakpoint here to inspect storage error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// DEBUG: Request successful, returning 201
			// BREAKPOINT: Set breakpoint here to see success response
			w.WriteHeader(http.StatusCreated)

		case http.MethodGet:
			// BREAKPOINT: Set breakpoint here to trace GET requests
			// DEBUG: Watch 'key' to see requested key
			key := r.URL.Query().Get("k")
			// DEBUG: Extracted key from query parameter

			// BREAKPOINT: Set breakpoint here to trace retrieval
			// DEBUG: Watch 's.Get' arguments and return values
			val, ok := s.Get(key)
			if !ok {
				// DEBUG: Key not found
				// BREAKPOINT: Set breakpoint here to trace 404 response
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			// DEBUG: Watch 'val' to see retrieved value
			// BREAKPOINT: Set breakpoint here to trace JSON encoding
			json.NewEncoder(w).Encode(map[string]string{"val": val})
			// DEBUG: Response encoded and sent

		default:
			// DEBUG: Unsupported method
			// BREAKPOINT: Set breakpoint here to trace unsupported methods
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// NewServer creates an HTTP server.
//
// BREAKPOINT: Set breakpoint here to trace server creation
// DEBUG: Watch 'addr' to see bind address
// DEBUG: Watch 'mux' to see registered routes
func NewServer(addr string, mux *http.ServeMux) *http.Server {
	// DEBUG: Creating HTTP server with addr and mux
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

// GracefulShutdown shuts down the server gracefully.
//
// Graceful Shutdown Algorithm:
// 1. Call srv.Shutdown(ctx) to initiate shutdown
// 2. Server stops accepting new connections
// 3. Server waits for in-flight requests to complete
// 4. Server closes when all requests are done or context times out
//
// BREAKPOINT: Set breakpoint here to trace shutdown
// DEBUG: Watch 'ctx' to see shutdown timeout
// DEBUG: Watch return error to see if shutdown succeeded
func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	// BREAKPOINT: Set breakpoint here before shutdown
	// DEBUG: Step over to see shutdown complete
	return srv.Shutdown(ctx)
	// DEBUG: Shutdown complete or timed out
}

// MemStore is an in-memory Store.
// BREAKPOINT: Set breakpoint in methods to trace storage operations
type MemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemStore creates a new in-memory store.
// BREAKPOINT: Set breakpoint here to trace store creation
func NewMemStore() Store {
	// DEBUG: Creating new MemStore with empty map
	return &MemStore{data: make(map[string]string)}
}

// Put stores a key-value pair.
//
// BREAKPOINT: Set breakpoint here to trace writes
// DEBUG: Watch 'key' and 'val' to see what's being stored
func (m *MemStore) Put(key, val string) error {
	// BREAKPOINT: Set breakpoint here before acquiring lock
	m.mu.Lock()
	defer m.mu.Unlock()

	// DEBUG: Watch 'm.data' before update
	m.data[key] = val
	// DEBUG: Watch 'm.data' after update to see new value

	return nil
}

// Get retrieves a value by key.
//
// BREAKPOINT: Set breakpoint here to trace reads
// DEBUG: Watch 'key' to see lookup key
// DEBUG: Watch 'val' and 'ok' to see result
func (m *MemStore) Get(key string) (string, bool) {
	// BREAKPOINT: Set breakpoint here before acquiring read lock
	m.mu.RLock()
	defer m.mu.RUnlock()

	// DEBUG: Watch 'm.data[key]' to see lookup
	val, ok := m.data[key]
	// DEBUG: Watch 'ok' to see if key exists

	return val, ok
}
