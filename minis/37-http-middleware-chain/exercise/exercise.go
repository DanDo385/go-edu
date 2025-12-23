//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"net/http"
)

// =============================================================================
// Context Keys and Helpers
// =============================================================================

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userKey      contextKey = "user"
)

// TODO: Implement WithRequestID to add request ID to context
func WithRequestID(ctx context.Context, id string) context.Context {
	// TODO: Use `context.WithValue` to return a new context that holds the ID.
	// The key should be the `requestIDKey` constant to avoid collisions.
	return context.WithValue(ctx, requestIDKey, id)
}

// TODO: Implement GetRequestID to retrieve request ID from context
func GetRequestID(ctx context.Context) (string, bool) {
	// TODO: Use `ctx.Value()` to get the value.
	// Perform a type assertion to check if the value is a string and if it exists.
	// The "comma ok" idiom (`value, ok := ...`) is the standard way to do this safely.
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

// TODO: Implement WithUser to add user to context
func WithUser(ctx context.Context, user *User) context.Context {
	// TODO: Use `context.WithValue` with the `userKey`.
	return context.WithValue(ctx, userKey, user)
}

// TODO: Implement GetUser to retrieve user from context
func GetUser(ctx context.Context) (*User, bool) {
	// TODO: Use `ctx.Value()` and a type assertion to `*User`.
	user, ok := ctx.Value(userKey).(*User)
	return user, ok
}

// =============================================================================
// Models
// =============================================================================

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// =============================================================================
// Response Writer Wrapper
// =============================================================================

// ResponseWriter wraps http.ResponseWriter to capture status code and bytes written
type ResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	bytesWritten  int
	headerWritten bool
}

// TODO: Implement NewResponseWriter to create a ResponseWriter wrapper
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	// TODO: Initialize and return a `ResponseWriter`.
	// - The embedded `ResponseWriter` should be `w`.
	// - The default status code for an HTTP response is 200 OK. Initialize `statusCode` to this value.
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// TODO: Implement WriteHeader to capture status code
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	// TODO: This method intercepts the call to `WriteHeader`.
	// - First, check if the header has already been written. If so, do nothing.
	// - If not, store the `statusCode` in the struct.
	// - Then, call the *original* `WriteHeader` method from the embedded `ResponseWriter`.
	// - Finally, set the `headerWritten` flag to true.
}

// TODO: Implement Write to count bytes written
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	// TODO: This method intercepts the call to `Write`.
	// - If `WriteHeader` has not been called yet, call it with a default of `http.StatusOK`.
	// - Call the original `Write` method.
	// - Add the number of bytes written (`n`) to `rw.bytesWritten`.
	// - Return the results of the original `Write` call.
	return 0, nil
}

// StatusCode returns the captured status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// BytesWritten returns the total bytes written
func (rw *ResponseWriter) BytesWritten() int {
	return rw.bytesWritten
}

// =============================================================================
// Middleware Type
// =============================================================================

type Middleware func(http.Handler) http.Handler

// =============================================================================
// Middleware Implementations
// =============================================================================

// TODO: Implement LoggingMiddleware to log request/response details
func LoggingMiddleware(next http.Handler) http.Handler {
	// TODO: Implement this middleware.
	// - The function should return an `http.HandlerFunc`.
	// - The handler func should:
	//   1. Record the start time: `start := time.Now()`.
	//   2. Create your `ResponseWriter` wrapper to capture the status code and bytes written.
	//   3. Log the incoming request (e.g., `log.Printf("→ %s %s", r.Method, r.URL.Path)`).
	//   4. Call the next handler in the chain: `next.ServeHTTP(rw, r)`.
	//   5. After the handler returns, calculate the duration (`time.Since(start)`).
	//   6. Log the outgoing response, including the status code and bytes written from your wrapper and the duration.
	return nil
}

// TODO: Implement RecoveryMiddleware to catch panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	// TODO: Implement this middleware.
	// - This middleware prevents a panic in a downstream handler from crashing the whole server.

	// Step 1: Return an `http.HandlerFunc`.
	// Step 2: Inside the handler, use `defer` to set up a function that will run *after* the `next` handler has finished (or panicked).
	//   - `defer func() { ... }()`
	// Step 3: Inside the deferred function, call `recover()`.
	//   - `if err := recover(); err != nil { ... }`
	//   - `recover()` will be non-nil only if a panic is currently being handled.
	// Step 4: If a panic was recovered:
	//   - Log the error and a stack trace for debugging (`debug.Stack()` is useful here).
	//   - Respond to the client with a generic error, like `http.Error(w, "Internal Server Error", http.StatusInternalServerError)`. Do not send the raw panic message to the client.
	// Step 5: After the `defer`, call the `next` handler as usual: `next.ServeHTTP(w, r)`.
	return nil
}

// TODO: Implement RequestIDMiddleware to assign unique ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	// TODO: Implement this middleware.

	// This middleware is responsible for ensuring every request has a unique ID, which is crucial for tracing and debugging in a complex system.

	// Step 1: Return an `http.HandlerFunc`.
	// Step 2: Inside the handler:
	//   - Check if the incoming request already has an `X-Request-ID` header. If so, use it.
	//   - If not, generate a new unique ID. For this exercise, a simple atomic counter is fine.
	//     - `id := atomic.AddUint64(&requestCounter, 1)`
	//   - Set the `X-Request-ID` header on the *response* so the client knows the ID.
	//   - Create a new context that contains the request ID using your `WithRequestID` helper.
	//   - Call the `next` handler, but with the *new* context: `next.ServeHTTP(w, r.WithContext(ctx))`.
	return nil
}

// TODO: Implement AuthMiddleware to validate authorization
func AuthMiddleware(next http.Handler) http.Handler {
	// TODO: Implement this middleware.

	// This middleware protects endpoints by checking for a valid authentication token.

	// Step 1: Return an `http.HandlerFunc`.
	// Step 2: Inside the handler:
	//   - Get the token from the `Authorization` header: `r.Header.Get("Authorization")`.
	//   - For this exercise, we'll use a simple check: if the token is not exactly `"Bearer valid-token"`, then it's invalid.
	//   - If the token is invalid, use `http.Error` to send a `401 Unauthorized` response and `return` immediately to stop further processing.
	//   - If the token is valid, create a sample `User` object.
	//   - Add the user to the request's context using your `WithUser` helper.
	//   - Call the `next` handler, passing the new context.
	return nil
}

// TODO: Implement CORSMiddleware to add CORS headers
func CORSMiddleware(allowOrigin string) Middleware {
	// TODO: This is a "middleware factory". It takes configuration and returns a `Middleware`.
	// - Return a `Middleware` function `func(next http.Handler) http.Handler { ... }`.
	// - That function, in turn, should return an `http.HandlerFunc`.
	// - Inside the `http.HandlerFunc`:
	//   1. Set the necessary CORS headers on the `ResponseWriter`:
	//      - `Access-Control-Allow-Origin`
	//      - `Access-Control-Allow-Methods`
	//      - `Access-Control-Allow-Headers`
	//   2. Handle the `OPTIONS` preflight request. If `r.Method == "OPTIONS"`, write `http.StatusOK` and return without calling `next`.
	//   3. For all other methods, call `next.ServeHTTP(w, r)`.
	return nil
}

// TODO: Implement MethodMiddleware to restrict HTTP methods
func MethodMiddleware(allowedMethods ...string) Middleware {
	// TODO: This is another "middleware factory".
	// - Return a `Middleware` function.
	// - Inside that function, for efficiency, you might want to convert the `allowedMethods` slice into a map for quick lookups.
	// - The `Middleware` should return an `http.HandlerFunc` that:
	//   1. Checks if `r.Method` is in the map of allowed methods.
	//   2. If not, responds with `http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)` and returns.
	//   3. If it is allowed, calls `next.ServeHTTP(w, r)`.
	return nil
}

// =============================================================================
// Chain Helper
// =============================================================================

// TODO: Implement Chain to compose multiple middleware
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// TODO: Implement this helper function.

	// The goal is to wrap the handler with all the middlewares.
	// If you have `Chain(h, M1, M2, M3)`, you want to produce `M1(M2(M3(h)))`.
	// This means the first middleware in the list is the outermost one, and it runs first.

	// Step 1: To achieve this, you need to apply the middleware in reverse order.
	// - Loop through the `middlewares` slice from the end to the beginning.
	// - `for i := len(middlewares) - 1; i >= 0; i-- { ... }`

	// Step 2: In each iteration, wrap the current `handler` with the middleware.
	// - `handler = middlewares[i](handler)`
	// - On the first iteration (with M3), `handler` becomes `M3(h)`.
	// - On the second iteration (with M2), `handler` becomes `M2(M3(h))`.
	// - On the final iteration (with M1), `handler` becomes `M1(M2(M3(h)))`.

	// Step 3: Return the fully wrapped handler.
	return nil
}
