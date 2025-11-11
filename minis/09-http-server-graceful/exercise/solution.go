/*
Problem: Build an HTTP server with routes, middleware, and graceful shutdown

Requirements:
1. REST endpoints for key-value storage
2. Request counting middleware
3. Graceful shutdown on SIGINT/SIGTERM
4. JSON request/response handling

Why Go is well-suited:
- http.ServeMux: Built-in routing
- http.Server.Shutdown: First-class graceful shutdown
- Middleware: Simple function wrapping

Compared to other languages:
- Node.js: Express.js similar, but requires library
- Python: Flask/FastAPI similar, ASGI for async
- Rust: Axum/Actix are powerful but more complex
*/

package exercise

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

// Store defines key-value operations.
type Store interface {
	Put(key, val string) error
	Get(key string) (string, bool)
}

var reqCount atomic.Int64

// RegisterRoutes sets up HTTP handlers.
func RegisterRoutes(mux *http.ServeMux, s Store) {
	mux.HandleFunc("/kv", withMiddleware(kvHandler(s)))
}

func withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count := reqCount.Add(1)
		w.Header().Set("X-Req-Count", fmt.Sprintf("%d", count))
		next(w, r)
	}
}

func kvHandler(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req struct {
				Key string `json:"key"`
				Val string `json:"val"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := s.Put(req.Key, req.Val); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

		case http.MethodGet:
			key := r.URL.Query().Get("k")
			val, ok := s.Get(key)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"val": val})

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// NewServer creates an HTTP server.
func NewServer(addr string, mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

// GracefulShutdown shuts down the server gracefully.
func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	return srv.Shutdown(ctx)
}

// MemStore is an in-memory Store.
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
