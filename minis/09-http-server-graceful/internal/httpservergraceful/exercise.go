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

func RegisterRoutes(mux *http.ServeMux, s Store) {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// TODO: Implement this function
	panic("not implemented")
}

func kvHandler(s Store) http.HandlerFunc {
	// TODO: Implement this function
	panic("not implemented")
}

func NewServer(addr string, mux *http.ServeMux) *http.Server {
	// TODO: Implement this function
	panic("not implemented")
}

func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	// TODO: Implement this function
	panic("not implemented")
}

func NewMemStore() Store {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *MemStore) Put(key, val string) error {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *MemStore) Get(key string) (string, bool) {
	// TODO: Implement this function
	panic("not implemented")
}
