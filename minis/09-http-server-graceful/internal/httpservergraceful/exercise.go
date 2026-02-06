//go:build !solution && !reference

package httpservergraceful

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

func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string]string)}
}

func (s *MemStore) Set(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
}

func (s *MemStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

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

// RunGracefulServer starts and runs an HTTP server that can be shut down gracefully.
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
