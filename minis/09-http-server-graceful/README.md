# Project 09: HTTP Server with Graceful Shutdown

## 1. What Is This About?

### Real-World Scenario

You deploy a new version of your HTTP server. Without graceful shutdown:

**❌ What happens:**
1. You send kill signal to old server
2. Server immediately stops
3. 50 requests currently being processed get dropped
4. Users see "Connection reset" errors
5. Some database writes are half-complete
6. Angry users, corrupted data

**✅ With graceful shutdown:**
1. You send termination signal
2. Server stops accepting NEW requests
3. Server waits for existing 50 requests to finish (up to 5 seconds)
4. All requests complete successfully
5. Server shuts down cleanly
6. Zero errors, happy users

This project teaches you how to build **production-ready HTTP servers** with:
- **Routing**: Handle different endpoints (/users, /posts, etc.)
- **Middleware**: Logging, authentication, request counting
- **Graceful shutdown**: No dropped requests during deployments
- **Clean architecture**: Interface-based design for testability

### What You'll Learn

1. **HTTP server basics**: net/http.Server, ServeMux, handlers
2. **Middleware pattern**: Request/response interceptors
3. **Graceful shutdown**: os/signal, context propagation
4. **Interface-based design**: Testable, mockable storage layer
5. **Production patterns**: Timeouts, logging, error handling

### The Challenge

Build an HTTP server with:
- GET and POST endpoints for key-value storage
- Logging middleware (logs every request)
- Request counter middleware (tracks total requests)
- Graceful shutdown on SIGINT/SIGTERM
- Configurable shutdown timeout (default 5 seconds)
- In-memory storage with interface abstraction

---

## 2. First Principles: HTTP Servers in Go

### What is an HTTP Server?

An **HTTP server** listens for incoming HTTP requests and sends back HTTP responses.

**Basic anatomy**:
```
Client              Server
  |                   |
  |--- GET /users ---|
  |                   |
  |--- 200 OK ------→|
  |    [{user data}]  |
```

**In Go**:
```go
http.ListenAndServe(":8080", handler)
```

This starts a server on port 8080 that calls `handler` for every request.

### What is a Handler?

A **handler** is any type that implements:
```go
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```

**Example**:
```go
type HelloHandler struct{}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}

http.ListenAndServe(":8080", HelloHandler{})
```

**Handler function** (simpler):
```go
http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
})
```

### What is Middleware?

**Middleware** wraps handlers to add functionality:

```
Request → Middleware 1 → Middleware 2 → Handler → Response
              ↓             ↓              ↓
           Logging      Auth Check   Business Logic
```

**Pattern**:
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)  // Call next handler
        log.Printf("Response sent")
    })
}
```

**Usage**:
```go
handler := LoggingMiddleware(HelloHandler{})
http.ListenAndServe(":8080", handler)
```

**Visual**:
```
Request
  ↓
LoggingMiddleware (before)
  ↓
HelloHandler
  ↓
LoggingMiddleware (after)
  ↓
Response
```

### What is Graceful Shutdown?

**Graceful shutdown** means stopping a server without interrupting active requests.

**Without graceful shutdown**:
```go
// Server receives SIGINT
os.Exit(0)  // Immediately kills process
// In-flight requests are dropped
```

**With graceful shutdown**:
```go
// Server receives SIGINT
server.Shutdown(ctx)
// 1. Stop accepting new requests
// 2. Wait for active requests to finish (up to timeout)
// 3. Then exit
```

**Real-world impact**:
- **Zero downtime deployments**: Rolling updates without errors
- **Data integrity**: Transactions complete before shutdown
- **User experience**: No "connection reset" errors

### Why is Graceful Shutdown Hard?

**Problem 1**: How do we detect shutdown signal?
- **Solution**: Listen for OS signals (SIGINT, SIGTERM)

**Problem 2**: How do we stop accepting new requests?
- **Solution**: `server.Shutdown()` stops the listener

**Problem 3**: How long should we wait for existing requests?
- **Solution**: Use context with timeout (e.g., 5 seconds)

**Problem 4**: What if some requests take too long?
- **Solution**: Force shutdown after timeout (trade-off: drop those requests)

---

## 3. Breaking Down the Solution

### Step 1: Server Setup

```go
server := &http.Server{
    Addr:    ":8080",
    Handler: handler,
}

go func() {
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Server error: %v", err)
    }
}()
```

**Why in a goroutine?**
`ListenAndServe()` blocks until server stops. We need to continue execution to handle shutdown signals.

### Step 2: Signal Handling

```go
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

<-stop  // Block until signal received
log.Println("Shutting down gracefully...")
```

**OS signals**:
- `os.Interrupt`: Ctrl+C in terminal (SIGINT)
- `syscall.SIGTERM`: Termination signal from process manager (Docker, Kubernetes)

### Step 3: Graceful Shutdown

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
    log.Printf("Force shutdown: %v", err)
}
```

**What `Shutdown()` does**:
1. Closes listener (no new connections accepted)
2. Waits for active connections to become idle
3. Returns when all connections closed OR context times out

### Step 4: Middleware Chain

```go
handler = LoggingMiddleware(
    CounterMiddleware(
        routes,
    ),
)
```

**Execution order**:
```
Request
  ↓
Logging (before)
  ↓
Counter (before)
  ↓
Routes
  ↓
Counter (after)
  ↓
Logging (after)
  ↓
Response
```

---

## 4. Complete Solution Walkthrough

### Server Structure

```go
type Server struct {
    store   Store
    counter *Counter
    server  *http.Server
}
```

**Store interface**:
```go
type Store interface {
    Set(key, val string)
    Get(key string) (string, bool)
}
```

**Why interface?**
- **Testability**: Mock storage in tests
- **Flexibility**: Swap in-memory → Redis → PostgreSQL without changing handler code
- **Separation of concerns**: Server doesn't care about storage implementation

### Handler: POST /kv

```go
func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
    // Parse JSON body
    var req SetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Validate
    if req.Key == "" {
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }

    // Store
    s.store.Set(req.Key, req.Val)

    // Respond
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

**Key patterns**:
1. **Parse**: `json.NewDecoder(r.Body).Decode(&req)`
2. **Validate**: Check required fields
3. **Business logic**: `s.store.Set()`
4. **Respond**: Write status code, encode JSON

### Handler: GET /kv?k=name

```go
func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
    // Parse query param
    key := r.URL.Query().Get("k")
    if key == "" {
        http.Error(w, "Missing key parameter", http.StatusBadRequest)
        return
    }

    // Retrieve
    val, ok := s.store.Get(key)
    if !ok {
        http.Error(w, "Key not found", http.StatusNotFound)
        return
    }

    // Respond
    json.NewEncoder(w).Encode(map[string]string{"key": key, "val": val})
}
```

**URL query parsing**:
- `/kv?k=name`: `r.URL.Query().Get("k")` → `"name"`
- `/kv?k=name&v=Go`: `r.URL.Query().Get("v")` → `"Go"`

### Logging Middleware

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("→ %s %s", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)

        duration := time.Since(start)
        log.Printf("← %s %s (%v)", r.Method, r.URL.Path, duration)
    })
}
```

**Output**:
```
→ POST /kv
← POST /kv (2.3ms)
→ GET /kv?k=name
← GET /kv?k=name (0.5ms)
```

### Counter Middleware

```go
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}

func (c *Counter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c.Increment()
        next.ServeHTTP(w, r)
    })
}
```

**Why mutex?**
Multiple goroutines (one per request) access `count` concurrently. Without mutex = data race.

### Graceful Shutdown Implementation

```go
func (s *Server) Run() error {
    // Start server in background
    go func() {
        if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    log.Printf("Server started on %s", s.server.Addr)

    // Wait for interrupt signal
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    log.Println("Shutting down gracefully...")

    // Shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := s.server.Shutdown(ctx); err != nil {
        return fmt.Errorf("shutdown failed: %w", err)
    }

    log.Println("Server stopped")
    return nil
}
```

**Timeline**:
```
T=0s   : Server starts, listening on :8080
T=10s  : User presses Ctrl+C
T=10s  : signal.Notify detects SIGINT
T=10s  : server.Shutdown() called
         - Listener closes (no new connections)
         - Waits for active requests
T=12s  : Last request finishes
T=12s  : server.Shutdown() returns
T=12s  : log "Server stopped"
T=12s  : Program exits
```

**With timeout**:
```
T=0s   : Server starts
T=10s  : Ctrl+C received
T=10s  : server.Shutdown() called with 5s timeout
T=15s  : Timeout! Some requests still active
T=15s  : server.Shutdown() returns error
T=15s  : log "shutdown failed: context deadline exceeded"
T=15s  : Program exits (remaining requests dropped)
```

---

## 5. Key Concepts Explained

### Concept 1: http.Server Configuration

```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      handler,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

**Timeouts prevent resource exhaustion**:
- **ReadTimeout**: Max time to read request (headers + body)
- **WriteTimeout**: Max time to write response
- **IdleTimeout**: Max time keep-alive connection stays open

**Why timeouts matter**:
```go
// Without ReadTimeout:
// Malicious client sends headers very slowly (1 byte per minute)
// Server waits forever, goroutine never exits
// Eventually: out of memory

// With ReadTimeout:
// After 10 seconds, server closes connection
// Goroutine exits, memory freed
```

### Concept 2: ServeMux (Multiplexer)

```go
mux := http.NewServeMux()
mux.HandleFunc("/users", handleUsers)
mux.HandleFunc("/posts", handlePosts)
```

**Routing logic**:
```
Request: GET /users
         ↓
    ServeMux checks registered patterns
         ↓
    Matches "/users"
         ↓
    Calls handleUsers
```

**Pattern matching**:
```go
mux.HandleFunc("/", handleRoot)        // Matches /anything
mux.HandleFunc("/api/", handleAPI)     // Matches /api/anything
mux.HandleFunc("/users", handleUsers)  // Exact match /users
```

**Precedence**: Most specific pattern wins
- Request `/api/users` matches `/api/` (not `/`)

### Concept 3: Middleware Composition

**Pattern**:
```go
type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Usage:
handler := Chain(
    myHandler,
    LoggingMiddleware,
    AuthMiddleware,
    RateLimitMiddleware,
)
```

**Why reverse order?**
```
Chain(handler, A, B, C) should produce:
  A(B(C(handler)))

Execution: Request → A → B → C → handler → C → B → A → Response
```

### Concept 4: Signal Handling

```go
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
```

**Why buffered channel (size 1)?**
If signal arrives before we're ready to receive, it's queued (not lost).

**Common signals**:
- `SIGINT` (2): Interrupt (Ctrl+C)
- `SIGTERM` (15): Termination (polite request to stop)
- `SIGKILL` (9): Kill (forceful, **can't be caught**)
- `SIGHUP` (1): Hangup (terminal closed, often used for "reload config")

### Concept 5: Context in Shutdown

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

server.Shutdown(ctx)
```

**What happens inside `Shutdown()`**:
```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()  // Timeout!
    default:
        if allConnectionsIdle() {
            return nil  // Success
        }
        time.Sleep(10 * time.Millisecond)
    }
}
```

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Response Helper

```go
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
    RespondJSON(w, status, map[string]string{"error": message})
}
```

### Pattern 2: Recovery Middleware

```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", 500)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

### Pattern 3: Request ID Middleware

```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Pattern 4: CORS Middleware

```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Pattern 5: Timeout Middleware

```go
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()

            done := make(chan struct{})

            go func() {
                next.ServeHTTP(w, r.WithContext(ctx))
                close(done)
            }()

            select {
            case <-done:
                // Request completed
            case <-ctx.Done():
                http.Error(w, "Request timeout", http.StatusGatewayTimeout)
            }
        })
    }
}
```

---

## 7. Real-World Applications

### Microservices

Every HTTP-based microservice needs graceful shutdown.

```go
type UserService struct {
    db     *sql.DB
    cache  *redis.Client
    server *Server
}

func (s *UserService) Shutdown() {
    s.server.Shutdown(context.Background())
    s.db.Close()
    s.cache.Close()
}
```

Companies: Netflix, Uber, Stripe (all use graceful shutdown)

### API Gateways

Handle thousands of concurrent requests, must shut down gracefully.

```go
gateway := &http.Server{
    Addr:         ":8080",
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

Companies: Kong, AWS API Gateway

### Admin Dashboards

Internal tools that need middleware for auth, logging, metrics.

```go
handler := Chain(
    dashboard,
    AuthMiddleware,
    LoggingMiddleware,
    MetricsMiddleware,
)
```

### Webhooks

Receive webhooks from third parties, must be reliable.

```go
func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
    // Must respond 200 OK quickly
    // Process webhook asynchronously
    go s.processWebhook(r)

    w.WriteHeader(http.StatusOK)
}
```

Companies: Stripe, Shopify, GitHub (all send webhooks)

### Health Check Endpoints

Kubernetes, Docker, load balancers need health checks.

```go
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    if err := s.db.Ping(); err != nil {
        http.Error(w, "Database down", http.StatusServiceUnavailable)
        return
    }

    RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
```

---

## 8. Common Mistakes to Avoid

### Mistake 1: Not Using Timeouts

**❌ Wrong**:
```go
server := &http.Server{Addr: ":8080", Handler: handler}
```

**Problem**: Slow clients can exhaust server resources.

**✅ Correct**:
```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      handler,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
```

### Mistake 2: Calling `os.Exit()` Directly

**❌ Wrong**:
```go
if err != nil {
    os.Exit(1)  // Immediate exit, drops requests
}
```

**✅ Correct**:
```go
if err != nil {
    server.Shutdown(ctx)  // Graceful shutdown first
    return err
}
```

### Mistake 3: Not Buffering Signal Channel

**❌ Wrong**:
```go
stop := make(chan os.Signal)  // Unbuffered
signal.Notify(stop, os.Interrupt)
```

**Problem**: If signal arrives before `<-stop`, it might be lost.

**✅ Correct**:
```go
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt)
```

### Mistake 4: Middleware Order Matters

**❌ Wrong**:
```go
handler := LoggingMiddleware(RecoveryMiddleware(routes))
```

**Problem**: If panic occurs, recovery catches it BEFORE logging.

**✅ Correct**:
```go
handler := RecoveryMiddleware(LoggingMiddleware(routes))
// Order: Recovery → Logging → Routes
// Panic in routes → caught by Recovery
```

### Mistake 5: Forgetting `defer cancel()`

**❌ Wrong**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
server.Shutdown(ctx)
```

**Problem**: Context resources leak.

**✅ Correct**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Shutdown(ctx)
```

### Mistake 6: Not Checking Shutdown Error

**❌ Wrong**:
```go
server.Shutdown(ctx)
// Ignores error
```

**Problem**: Don't know if shutdown timed out (some requests dropped).

**✅ Correct**:
```go
if err := server.Shutdown(ctx); err != nil {
    log.Printf("Shutdown failed (requests may be dropped): %v", err)
}
```

---

## 9. Stretch Goals

### Goal 1: Add Metrics Endpoint ⭐

Expose `/metrics` endpoint with Prometheus metrics.

**Hint**:
```go
var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "http_requests_total"},
        []string{"method", "path", "status"},
    )
)

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rec := &statusRecorder{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(rec, r)
        requestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rec.statusCode)).Inc()
    })
}
```

### Goal 2: Add Rate Limiting ⭐⭐

Limit requests per client (by IP).

**Hint**:
```go
type RateLimiter struct {
    mu      sync.Mutex
    clients map[string]*rate.Limiter
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        limiter := rl.getLimiter(ip)

        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Goal 3: Add TLS Support ⭐⭐

Enable HTTPS with certificates.

**Hint**:
```go
server := &http.Server{Addr: ":443", Handler: handler}
server.ListenAndServeTLS("cert.pem", "key.pem")
```

### Goal 4: Structured Logging ⭐⭐⭐

Use structured logging (JSON) instead of plain text.

**Hint**:
```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
logger.Info("Request received",
    zap.String("method", r.Method),
    zap.String("path", r.URL.Path),
    zap.Duration("duration", duration),
)
```

### Goal 5: Distributed Tracing ⭐⭐⭐

Add OpenTelemetry tracing.

**Hint**:
```go
import "go.opentelemetry.io/otel"

func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, span := otel.Tracer("http-server").Start(r.Context(), r.URL.Path)
        defer span.End()

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## How to Run

```bash
# Run the server
go run ./minis/09-http-server-graceful/cmd/http-server-graceful

# Test with curl
curl -X POST http://localhost:8080/kv -d '{"key":"name","val":"Go"}'
curl http://localhost:8080/kv?k=name

# Stop with Ctrl+C (graceful shutdown)
```

---

## Summary

**What you learned**:
- ✅ HTTP server setup with net/http
- ✅ Middleware pattern for request interception
- ✅ Graceful shutdown with os/signal
- ✅ Context propagation for timeouts
- ✅ Interface-based storage abstraction
- ✅ Production patterns (logging, counting, error handling)

**Why this matters**:
Every production HTTP service needs graceful shutdown. Without it, deployments cause dropped requests and angry users. Middleware enables clean separation of concerns (logging, auth, metrics).

**Key takeaway**:
Graceful shutdown = Zero downtime deployments

**Next steps**:
- Project 10: Learn gRPC for high-performance RPC with streaming

Build reliably! 🚀
