//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"net/http"
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
	// TODO: implement
}

// NewServer creates an HTTP server with graceful shutdown support.
func NewServer(addr string, mux *http.ServeMux) *http.Server {
	// TODO: implement
	return nil
}

// GracefulShutdown listens for signals and shuts down the server cleanly.
func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	// TODO: implement
	return nil
}

// MemStore is an in-memory Store implementation for testing.
type MemStore struct {
	// TODO: add fields
}

func NewMemStore() Store {
	// TODO: implement
	return nil
}

func (m *MemStore) Put(key, val string) error {
	// TODO: implement
	return nil
}

func (m *MemStore) Get(key string) (string, bool) {
	// TODO: implement
	return "", false
}
