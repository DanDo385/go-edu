//go:build !solution && !reference

package httpservergraceful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

type Store interface {
	Put(key, val string) error
	Get(key string) (string, bool)
}

type MemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// RegisterRoutes implements the exercise.
//
// TODO: Implement this function
func RegisterRoutes(mux *http.ServeMux, s Store) {
	// TODO: Implement
}

// NewServer implements the exercise.
//
// TODO: Implement this function
func NewServer(addr string, mux *http.ServeMux) *http.Server {
	// TODO: Implement
	return nil
}

// GracefulShutdown implements the exercise.
//
// TODO: Implement this function
func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	// TODO: Implement
	return nil
}

// NewMemStore implements the exercise.
//
// TODO: Implement this function
func NewMemStore() Store {
	// TODO: Implement
	return Store{}
}

// Put implements the exercise.
//
// TODO: Implement this function
func (m *MemStore) Put(key string, val string) error {
	// TODO: Implement
	return nil
}

// Get implements the exercise.
//
// TODO: Implement this function
func (m *MemStore) Get(key string) (string, bool) {
	// TODO: Implement
	return "", false
}
