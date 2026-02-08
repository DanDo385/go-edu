//go:build reference

package httpservergraceful

/*
Reference Solution - HTTP Key-Value Server with Graceful Shutdown
================================================================

This file implements an in-memory key-value HTTP API (GET/POST) and demonstrates
graceful shutdown: on SIGINT/SIGTERM or context cancel, the server stops accepting
new connections but allows in-flight requests to complete before exiting.

This connects to:
- net/http: ServeMux, Handler, JSON encoding, request/response lifecycle
- os/signal: capture SIGINT (Ctrl+C) and SIGTERM for graceful exit
- sync.RWMutex: concurrent read access, exclusive write for map
- atomic: lock-free request counter for X-Req-Count header

The exercise teaches:
- Graceful shutdown: Shutdown() stops listener, waits for active handlers
- Signal handling: signal.Notify, buffered channel (size 1) for sigCh
- Context: RunGracefulServer accepts ctx for programmatic cancellation
- Concurrency: MemStore protected by RWMutex; multiple readers, single writer
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type MemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemStore - In-memory key-value store, thread-safe.
func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string]string)}
}

// Set - Write key-value. Uses Lock (exclusive) for writes.
func (s *MemStore) Set(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
}

// Get - Read value by key. Uses RLock (shared) so multiple readers can run concurrently.
func (s *MemStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

/*
RegisterRoutes - Mount /kv Handler with Request Counter

POST /kv: JSON body {key, val} - stores key-value
GET /kv?k=key - returns {"val": "..."} or 404
X-Req-Count header: atomic counter of requests (teaching atomic operations)
*/
func RegisterRoutes(mux *http.ServeMux, store *MemStore) {
	var reqCount uint64
	kv := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var body struct {
				Key string `json:"key"`
				Val string `json:"val"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			if body.Key == "" {
				http.Error(w, "missing key", http.StatusBadRequest)
				return
			}
			store.Set(body.Key, body.Val)
			w.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			key := r.URL.Query().Get("k")
			if key == "" {
				http.Error(w, "missing key query param", http.StatusBadRequest)
				return
			}
			val, ok := store.Get(key)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"val": val})
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/kv", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddUint64(&reqCount, 1)
		w.Header().Set("X-Req-Count", fmt.Sprintf("%d", count))
		kv.ServeHTTP(w, r)
	}))
}

/*
RunGracefulServer - Run Server Until Signal or Context Cancel

1. Start server in goroutine (ListenAndServe blocks)
2. Wait for: startup error, SIGINT/SIGTERM, or ctx.Done()
3. Call Shutdown with 5s timeout - stops listener, waits for active requests
4. Return Shutdown error (or nil if clean)

Signal channel must be buffered (size 1) so signal.Notify doesn't block
when we're not reading. signal.Stop on defer prevents further delivery.
*/
func RunGracefulServer(ctx context.Context, srv *http.Server) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		fmt.Println("shutdown signal received")
	case <-ctx.Done():
		fmt.Println("context canceled; shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
