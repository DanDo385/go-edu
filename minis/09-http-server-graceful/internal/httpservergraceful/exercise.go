//go:build !solution && !reference

package httpservergraceful

/*
Problem: Build an HTTP server with routes, middleware, and graceful shutdown
Requirements:
1. REST endpoints for key-value storage
2. Request counting middleware
3. Graceful shutdown on SIGINT/SIGTERM
Algorithm: HTTP Request Handling
- Route requests to appropriate handlers
- Apply middleware for cross-cutting concerns
- Handle JSON encoding/decoding
- Coordinate graceful shutdown
*/

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

// RegisterRoutes - TODO: implement this function
func RegisterRoutes(mux *http.ServeMux, s Store) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
}

// withMiddleware - TODO: implement this function
func withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// kvHandler - TODO: implement this function
func kvHandler(s Store) http.HandlerFunc {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// NewServer - TODO: implement this function
func NewServer(addr string, mux *http.ServeMux) *http.Server {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// GracefulShutdown - TODO: implement this function
func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

type MemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemStore - TODO: implement this function
func NewMemStore() Store {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// Put - TODO: implement this function
func (m *MemStore) Put(key, val string) error {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// Get - TODO: implement this function
func (m *MemStore) Get(key string) (string, bool) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

