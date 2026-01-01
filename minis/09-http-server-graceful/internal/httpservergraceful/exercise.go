//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package httpservergraceful

import (
	"sync/atomic"
	"context"

	"net/http"
	"sync"
)

type Store interface {
	Put(key, val string) error
	Get(key string) (string, bool)
}

var reqCount atomic.Int64
// TODO: implement RegisterRoutes.
func RegisterRoutes(mux *http.ServeMux, s Store) { panic("TODO: implement") }
// TODO: implement withMiddleware.
func withMiddleware(next http.HandlerFunc) http.HandlerFunc { panic("TODO: implement") }
// TODO: implement kvHandler.
func kvHandler(s Store) http.HandlerFunc { panic("TODO: implement") }
// TODO: implement NewServer.
func NewServer(addr string, mux *http.ServeMux) *http.Server { panic("TODO: implement") }
// TODO: implement GracefulShutdown.
func GracefulShutdown(ctx context.Context, srv *http.Server) error { panic("TODO: implement") }

type MemStore struct {
	mu   sync.RWMutex
	data map[string]string
}
// TODO: implement NewMemStore.
func NewMemStore() Store { panic("TODO: implement") }
// TODO: implement Put.
func (m *MemStore) Put(key, val string) error { panic("TODO: implement") }
// TODO: implement Get.
func (m *MemStore) Get(key string) (string, bool) { panic("TODO: implement") }
