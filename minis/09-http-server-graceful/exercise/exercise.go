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
	// TODO: Implement the function.

	// This function orchestrates the setup of your HTTP routes and middleware.

	// Step 1: Create the main handler for the "/kv" route.
	// - The `kvHandler` function (which you'll also implement) will be responsible for the core logic of your key-value store.
	// - It needs access to the `Store` interface, so you'll create it by calling `kvHandler(s)`.

	// Step 2: Create the middleware.
	// - The `withMiddleware` function (which you'll also implement) will wrap your main handler.
	// - It will add functionality that runs on every request, like incrementing a counter.
	// - You'll wrap your `kvHandler` by calling `withMiddleware(kvHandler(s))`.

	// Step 3: Register the handler with the `http.ServeMux`.
	// - `mux.HandleFunc("/kv", ...)`
	// - The second argument will be the result of your wrapped handler from Step 2.
}

// NewServer creates an HTTP server with graceful shutdown support.
func NewServer(addr string, mux *http.ServeMux) *http.Server {
	// TODO: Implement this function.
	// - Create and return a new `http.Server`.
	// - The `Addr` field should be set to the `addr` parameter.
	// - The `Handler` field should be set to the `mux` parameter.
	return nil
}

// GracefulShutdown listens for signals and shuts down the server cleanly.
func GracefulShutdown(ctx context.Context, srv *http.Server) error {
	// TODO: Implement this function.
	// - Call `srv.Shutdown(ctx)`.
	// - The `context` is used to set a deadline for the shutdown process.
	// - `Shutdown` gracefully shuts down the server without interrupting any active connections.
	//   It works by first closing all open listeners, then closing all idle connections, and then waiting indefinitely for connections to return to idle and then shut down.
	// - Return the error from `srv.Shutdown`.
	return nil
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
	// TODO: Implement this function.
	// - It should have the signature `func withMiddleware(next http.HandlerFunc) http.HandlerFunc`.
	// - It takes one handler (`next`) and returns a new handler.
	// - The returned handler should:
	//   1. Increment a global atomic counter (`reqCount.Add(1)`). An atomic integer is used to prevent race conditions when multiple requests come in at the same time.
	//   2. Set a response header: `w.Header().Set("X-Req-Count", ...)`
	//   3. Call the `next` handler to process the rest of the request: `next(w, r)`.
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func kvHandler(s Store) http.HandlerFunc {
	// TODO: Implement this function.
	// - It should have the signature `func kvHandler(s Store) http.HandlerFunc`.
	// - It takes the `Store` and returns a handler function. This is a common pattern in Go for giving handlers access to dependencies.
	// - The returned handler should use a `switch r.Method` statement to handle different HTTP methods.
	//   - `case http.MethodPost`:
	//     - Decode the JSON body of the request into a struct. Use `json.NewDecoder(r.Body).Decode(&yourStruct)`.
	//     - Call `s.Put(...)` to store the key and value.
	//     - Handle potential errors from decoding and storing, returning appropriate HTTP status codes (`http.StatusBadRequest`, `http.StatusInternalServerError`).
	//     - If successful, respond with `http.StatusCreated`.
	//   - `case http.MethodGet`:
	//     - Get the key from the URL query parameters: `r.URL.Query().Get("k")`.
	//     - Call `s.Get(...)` to retrieve the value.
	//     - If the key is not found, return `http.StatusNotFound`.
	//     - If found, encode the value into a JSON response: `json.NewEncoder(w).Encode(yourMapOrStruct)`.
	//   - `default`:
	//     - For any other method, return `http.StatusMethodNotAllowed`.
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
