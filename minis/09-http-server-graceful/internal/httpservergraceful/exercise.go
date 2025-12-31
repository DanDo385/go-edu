//go:build !solution
// +build !solution

package httpservergraceful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

// Store defines key-value storage operations.
type Store interface {
	Put(key, val string) error
	Get(key string) (string, bool)
}

// ============================================================================
// HTTP SERVER: REST API with Middleware and Graceful Shutdown
// ============================================================================
//
// REST API Pattern:
// - POST /kv with JSON body → Store key-value
// - GET /kv?k=key → Retrieve value
//
// Middleware Pattern:
// - Wrap handlers to add cross-cutting functionality
// - Execute before/after main handler
// - Common uses: logging, metrics, auth
//
// Graceful Shutdown Algorithm:
// 1. Receive shutdown signal
// 2. Stop accepting new connections
// 3. Wait for in-flight requests to complete
// 4. Close server cleanly
//
// ============================================================================

// TODO: Implement RegisterRoutes(mux *http.ServeMux, s Store)
// Algorithm:
// 1. Create handler: kvHandler(s)
// 2. Wrap with middleware: withMiddleware(kvHandler(s))
// 3. Register with mux: mux.HandleFunc("/kv", ...)

// TODO: Implement NewServer(addr string, mux *http.ServeMux) *http.Server
// Create http.Server with Addr=addr, Handler=mux

// TODO: Implement GracefulShutdown(ctx context.Context, srv *http.Server) error
// Call srv.Shutdown(ctx) to gracefully close server

// MemStore is an in-memory Store implementation.
type MemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemStore() Store {
	return &MemStore{data: make(map[string]string)}
}

func (m *MemStore) Put(key, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
	return nil
}

func (m *MemStore) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

var reqCount atomic.Int64

// TODO: Implement withMiddleware(next http.HandlerFunc) http.HandlerFunc
// Algorithm:
// 1. Return new handler function
// 2. Increment reqCount atomically
// 3. Set X-Req-Count header
// 4. Call next(w, r)

func withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

// TODO: Implement kvHandler(s Store) http.HandlerFunc
// Algorithm:
// 1. Return handler that switches on r.Method
// 2. POST:
//    - Decode JSON body into struct{Key, Val string}
//    - Call s.Put(key, val)
//    - Return 201 Created
// 3. GET:
//    - Get key from r.URL.Query().Get("k")
//    - Call s.Get(key)
//    - If not found, return 404
//    - Encode JSON response with value
// 4. Default: return 405 Method Not Allowed

func kvHandler(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
