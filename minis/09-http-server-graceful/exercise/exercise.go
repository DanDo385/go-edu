//go:build !solution
// +build !solution

package exercise

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

// RegisterRoutes sets up HTTP handlers on mux:
//   POST /kv - accepts {"key":"k","val":"v"}
//   GET /kv?k=... - returns {"val":"..."}  or 404
//
// Middleware: Add X-Req-Count header with request counter
func RegisterRoutes(mux *http.ServeMux, s Store) {
	mux.HandleFunc("/kv", withMiddleware(kvHandler(s)))
}

// NewServer creates an HTTP server with graceful shutdown support.
func NewServer(addr string, mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

// GracefulShutdown listens for signals and shuts down the server cleanly.
func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	return srv.Shutdown(ctx)
}

// MemStore is an in-memory Store implementation for testing.
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
			var body struct {
				Key string `json:"key"`
				Val string `json:"val"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := s.Put(body.Key, body.Val); err != nil {
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
