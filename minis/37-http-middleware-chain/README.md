# Project 37: HTTP Middleware Chain

## 1. What Is This About?

### Real-World Scenario

You're building a REST API for a production application. Every endpoint needs:
- **Logging**: Track all requests for debugging and analytics
- **Authentication**: Verify user tokens
- **Rate limiting**: Prevent abuse
- **Request IDs**: Track requests across microservices
- **Metrics**: Measure response times and error rates
- **CORS**: Allow browser requests from different origins
- **Recovery**: Catch panics and return 500 instead of crashing

**❌ Naive approach:** Copy-paste this logic into every handler
```go
func handleUsers(w http.ResponseWriter, r *http.Request) {
    // Logging
    log.Printf("Request: %s %s", r.Method, r.URL.Path)

    // Auth
    token := r.Header.Get("Authorization")
    if !isValid(token) {
        http.Error(w, "Unauthorized", 401)
        return
    }

    // Rate limiting
    if !checkRateLimit(r.RemoteAddr) {
        http.Error(w, "Too many requests", 429)
        return
    }

    // Business logic (finally!)
    users := getUsers()
    json.NewEncoder(w).Encode(users)
}

// Copy-paste the same 20 lines for every handler... 😱
```

**Problems:**
- **Code duplication**: Same logic in 50+ handlers
- **Hard to maintain**: Bug fix requires updating 50+ files
- **Error-prone**: Easy to forget auth check in one handler
- **Not composable**: Can't easily reorder or disable features

**✅ Middleware approach:**
```go
func handleUsers(w http.ResponseWriter, r *http.Request) {
    users := getUsers()
    json.NewEncoder(w).Encode(users)
}

// Build middleware chain once, reuse everywhere
handler := Chain(
    handleUsers,
    LoggingMiddleware,
    AuthMiddleware,
    RateLimitMiddleware,
    RecoveryMiddleware,
)
```

**Benefits:**
- **Reusable**: Write once, apply to all endpoints
- **Composable**: Mix and match middleware as needed
- **Maintainable**: Fix bugs in one place
- **Testable**: Test middleware independently
- **Readable**: Business logic is clear and focused

This project teaches you how to build **production-grade middleware chains** that are:
- **Composable**: Stack middleware like LEGO blocks
- **Type-safe**: Compile-time guarantees
- **Context-aware**: Pass request-scoped data between middleware
- **Performance-conscious**: Zero allocation patterns
- **Flexible**: Conditional and parameterized middleware

### What You'll Learn

1. **Middleware pattern**: Function wrappers for cross-cutting concerns
2. **Function composition**: Building chains with higher-order functions
3. **Request/response wrapping**: Intercepting and modifying HTTP traffic
4. **Context passing**: Sharing data between middleware layers
5. **Response writers**: Capturing status codes and bytes written
6. **Advanced patterns**: Conditional, parameterized, and async middleware

### The Challenge

Build a middleware system with:
- Basic middleware (logging, recovery, request ID)
- Middleware chaining (compose multiple middleware)
- Context-based data passing (request-scoped values)
- Response writer wrapping (capture status codes)
- Parameterized middleware (configurable behavior)
- Conditional middleware (apply based on conditions)

---

## 2. First Principles: The Middleware Pattern

### What is Middleware?

**Middleware** is a function that wraps an HTTP handler to add functionality before or after the handler executes.

**Core concept:**
```
Request → Middleware → Handler → Response
              ↓
          Before
          Handler runs
          After
```

**In Go:**
```go
type Middleware func(http.Handler) http.Handler
```

A middleware takes a handler and returns a new handler that wraps it.

### How Does Middleware Work?

**Simple example:**
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("→ %s %s", r.Method, r.URL.Path)  // BEFORE

        next.ServeHTTP(w, r)  // Call the wrapped handler

        log.Printf("← %s %s", r.Method, r.URL.Path)  // AFTER
    })
}
```

**What happens:**
1. LoggingMiddleware receives `next` (the handler to wrap)
2. Returns a new handler function
3. When request arrives:
   - Log "→ GET /users" (before)
   - Call `next.ServeHTTP(w, r)` (execute wrapped handler)
   - Log "← GET /users" (after)

**Visual flow:**
```
Client
  ↓
[Logging Middleware]
  ↓ log "→ GET /users"
  ↓
  ↓ call next.ServeHTTP()
  ↓
[Actual Handler]
  ↓ execute business logic
  ↓ return
  ↓
[Logging Middleware]
  ↓ log "← GET /users"
  ↓
Client
```

### Why is http.Handler an Interface?

**The http.Handler interface:**
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**Why interface?**
- **Flexibility**: Any type can be a handler
- **Composability**: Middleware wraps handlers, which are also handlers
- **Testability**: Easy to create mock handlers

**Three ways to create a handler:**

1. **Struct implementing ServeHTTP:**
```go
type MyHandler struct{}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
}
```

2. **HandlerFunc adapter:**
```go
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
})
```

3. **Helper methods:**
```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
})
```

### What is Middleware Chaining?

**Chaining** means wrapping a handler with multiple middleware layers.

**Nested wrapping:**
```go
handler := LoggingMiddleware(
    AuthMiddleware(
        RateLimitMiddleware(
            actualHandler,
        ),
    ),
)
```

**Execution order:**
```
Request
  ↓
[Logging] before
  ↓
[Auth] before
  ↓
[RateLimit] before
  ↓
[Handler]
  ↓
[RateLimit] after
  ↓
[Auth] after
  ↓
[Logging] after
  ↓
Response
```

**Chain helper function:**
```go
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    // Apply middleware in reverse order
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Usage:
handler := Chain(
    actualHandler,
    LoggingMiddleware,
    AuthMiddleware,
    RateLimitMiddleware,
)
```

**Why reverse order?**
```
Chain(handler, A, B, C) should execute as: A → B → C → handler

We want: A(B(C(handler)))

So we apply: C first, then B, then A (reverse order)
```

### What is Context?

**context.Context** carries request-scoped values, cancellation signals, and deadlines.

**Use cases in middleware:**

1. **Pass request ID:**
```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Later in handler:
func handleUsers(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value("request_id").(string)
    log.Printf("[%s] Getting users", requestID)
}
```

2. **Pass user info:**
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := authenticateUser(r)
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

3. **Timeout propagation:**
```go
func TimeoutMiddleware(timeout time.Duration) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Context rules:**
- **Immutable**: Don't modify context, create new one with `WithValue`
- **Request-scoped**: Context dies when request completes
- **Type safety**: Use typed keys (not strings) in production
- **Never nil**: Always pass a context, use `context.Background()` if needed

---

## 3. Breaking Down the Solution

### Step 1: Basic Middleware Pattern

**Template:**
```go
func MiddlewareName(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // BEFORE: Execute before handler

        next.ServeHTTP(w, r)  // Call wrapped handler

        // AFTER: Execute after handler
    })
}
```

**Example 1: Logging**
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

**Example 2: Recovery (Panic Handler)**
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

### Step 2: Parameterized Middleware

**Pattern:** Middleware that accepts configuration parameters.

**Template:**
```go
func MiddlewareWithConfig(config Config) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Use config here
            next.ServeHTTP(w, r)
        })
    }
}
```

**Example 1: Timeout**
```go
func TimeoutMiddleware(timeout time.Duration) Middleware {
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

// Usage:
handler := TimeoutMiddleware(5 * time.Second)(actualHandler)
```

**Example 2: CORS**
```go
func CORSMiddleware(allowOrigin string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// Usage:
handler := CORSMiddleware("*")(actualHandler)
```

### Step 3: Response Writer Wrapping

**Problem:** Standard `http.ResponseWriter` doesn't expose status code or bytes written.

**Solution:** Wrap `ResponseWriter` to capture this information.

**Implementation:**
```go
type ResponseWriter struct {
    http.ResponseWriter
    statusCode int
    bytesWritten int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
    return &ResponseWriter{
        ResponseWriter: w,
        statusCode:     http.StatusOK,  // Default 200
    }
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
    rw.statusCode = statusCode
    rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
    n, err := rw.ResponseWriter.Write(b)
    rw.bytesWritten += n
    return n, err
}

func (rw *ResponseWriter) StatusCode() int {
    return rw.statusCode
}

func (rw *ResponseWriter) BytesWritten() int {
    return rw.bytesWritten
}
```

**Usage in middleware:**
```go
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := NewResponseWriter(w)
        start := time.Now()

        next.ServeHTTP(rw, r)

        duration := time.Since(start)
        log.Printf(
            "%s %s - Status: %d, Bytes: %d, Duration: %v",
            r.Method, r.URL.Path,
            rw.StatusCode(),
            rw.BytesWritten(),
            duration,
        )
    })
}
```

### Step 4: Context-Based Data Passing

**Pattern:** Use context to pass data between middleware layers.

**Type-safe context keys:**
```go
type contextKey string

const (
    requestIDKey contextKey = "request_id"
    userKey      contextKey = "user"
)

// Helper functions
func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey).(string)
    return id, ok
}

func WithUser(ctx context.Context, user *User) context.Context {
    return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(userKey).(*User)
    return user, ok
}
```

**Middleware using context:**
```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := uuid.New().String()
        ctx := WithRequestID(r.Context(), id)
        w.Header().Set("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        user, err := authenticateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        ctx := WithUser(r.Context(), user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// In handler:
func handleProfile(w http.ResponseWriter, r *http.Request) {
    user, ok := GetUser(r.Context())
    if !ok {
        http.Error(w, "User not found", 500)
        return
    }

    json.NewEncoder(w).Encode(user)
}
```

### Step 5: Middleware Chaining

**Implementation:**
```go
type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
    // Apply middleware in reverse order so first middleware wraps all others
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}
```

**Why reverse?**
```
Input:  Chain(H, A, B, C)
Want:   A(B(C(H)))  // A wraps everything

Loop iteration:
  i=2: H = C(H)
  i=1: H = B(C(H))
  i=0: H = A(B(C(H)))  ✓
```

**Usage:**
```go
handler := Chain(
    handleUsers,
    RecoveryMiddleware,
    LoggingMiddleware,
    RequestIDMiddleware,
    AuthMiddleware,
)

// Execution order:
// Recovery → Logging → RequestID → Auth → handleUsers
```

**Alternative: Method chaining**
```go
type MiddlewareChain struct {
    handler http.Handler
}

func NewChain(handler http.Handler) *MiddlewareChain {
    return &MiddlewareChain{handler: handler}
}

func (c *MiddlewareChain) Use(middleware Middleware) *MiddlewareChain {
    c.handler = middleware(c.handler)
    return c
}

func (c *MiddlewareChain) Handler() http.Handler {
    return c.handler
}

// Usage:
handler := NewChain(handleUsers).
    Use(AuthMiddleware).
    Use(RequestIDMiddleware).
    Use(LoggingMiddleware).
    Use(RecoveryMiddleware).
    Handler()
```

---

## 4. Complete Solution Walkthrough

### Full Example: Production-Ready API

**main.go:**
```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

// User model
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// Handlers
func handleUsers(w http.ResponseWriter, r *http.Request) {
    users := []User{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }
    json.NewEncoder(w).Encode(users)
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
    user, ok := GetUser(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", 401)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func main() {
    // Public routes (no auth)
    publicMux := http.NewServeMux()
    publicMux.HandleFunc("/users", handleUsers)

    publicHandler := Chain(
        publicMux,
        RecoveryMiddleware,
        LoggingMiddleware,
        RequestIDMiddleware,
    )

    // Protected routes (auth required)
    protectedMux := http.NewServeMux()
    protectedMux.HandleFunc("/profile", handleProfile)

    protectedHandler := Chain(
        protectedMux,
        RecoveryMiddleware,
        LoggingMiddleware,
        RequestIDMiddleware,
        AuthMiddleware,
    )

    // Combine routes
    mux := http.NewServeMux()
    mux.Handle("/", publicHandler)
    mux.Handle("/profile", protectedHandler)

    // Start server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    log.Println("Server starting on :8080")
    log.Fatal(server.ListenAndServe())
}
```

### Middleware Implementations

**1. Recovery Middleware:**
```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

**2. Logging Middleware:**
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := NewResponseWriter(w)

        requestID, _ := GetRequestID(r.Context())
        log.Printf("[%s] → %s %s", requestID, r.Method, r.URL.Path)

        next.ServeHTTP(rw, r)

        duration := time.Since(start)
        log.Printf(
            "[%s] ← %s %s - %d (%v)",
            requestID,
            r.Method,
            r.URL.Path,
            rw.StatusCode(),
            duration,
        )
    })
}
```

**3. Request ID Middleware:**
```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := WithRequestID(r.Context(), requestID)
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**4. Auth Middleware:**
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")

        user, err := validateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        ctx := WithUser(r.Context(), user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func validateToken(token string) (*User, error) {
    // In production: verify JWT, check database, etc.
    if token == "" {
        return nil, errors.New("missing token")
    }

    return &User{ID: 1, Name: "Alice"}, nil
}
```

**5. Rate Limit Middleware:**
```go
type RateLimiter struct {
    mu      sync.Mutex
    clients map[string]*rate.Limiter
    rate    rate.Limit
    burst   int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
    return &RateLimiter{
        clients: make(map[string]*rate.Limiter),
        rate:    r,
        burst:   b,
    }
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    limiter, exists := rl.clients[ip]
    if !exists {
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.clients[ip] = limiter
    }

    return limiter
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        limiter := rl.getLimiter(r.RemoteAddr)

        if !limiter.Allow() {
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## 5. Key Concepts Explained

### Concept 1: Middleware Execution Order

**Middleware order matters!**

```go
handler := Chain(
    actualHandler,
    A, B, C,
)
```

**Creates:** `A(B(C(actualHandler)))`

**Execution flow:**
```
Request
  ↓
A before
  ↓
B before
  ↓
C before
  ↓
Handler
  ↓
C after
  ↓
B after
  ↓
A after
  ↓
Response
```

**Common ordering patterns:**

1. **Recovery should be outermost** (catch panics from all middleware)
```go
Chain(handler, Recovery, Logging, Auth)
```

2. **Logging should be early** (log even if auth fails)
```go
Chain(handler, Recovery, Logging, Auth)
```

3. **Auth before business logic** (reject unauthorized early)
```go
Chain(handler, Recovery, Logging, Auth, RateLimit)
```

### Concept 2: Short-Circuiting

**Middleware can stop the chain** by not calling `next.ServeHTTP()`.

**Example: Auth failure**
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", 401)
            return  // ← Chain stops here!
        }

        next.ServeHTTP(w, r)  // Only called if auth succeeds
    })
}
```

**Flow with auth failure:**
```
Request
  ↓
Recovery (before)
  ↓
Logging (before)
  ↓
Auth (before)
  ↓
Check auth → FAIL
  ↓
Return 401
  ↓
Logging (after)  ← Still executes!
  ↓
Recovery (after)
  ↓
Response
```

**Note:** "After" code in outer middleware still runs even if inner middleware short-circuits.

### Concept 3: Context Propagation

**Context flows through the entire chain.**

**Example:**
```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := WithRequestID(r.Context(), uuid.New().String())
        next.ServeHTTP(w, r.WithContext(ctx))  // ← Pass new context
    })
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID, _ := GetRequestID(r.Context())  // ← Read from context
        log.Printf("[%s] Request started", requestID)
        next.ServeHTTP(w, r)
    })
}
```

**Important:** Always pass context forward with `r.WithContext(ctx)`.

### Concept 4: ResponseWriter Interface

**http.ResponseWriter interface:**
```go
type ResponseWriter interface {
    Header() http.Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```

**Why wrap it?**
- Standard interface doesn't expose status code
- Can't tell what status was written
- Can't count bytes written
- Can't intercept writes

**Wrapper pattern:**
```go
type ResponseWriter struct {
    http.ResponseWriter  // Embed original
    statusCode     int
    bytesWritten   int
    headerWritten  bool
}
```

**Gotcha: WriteHeader is optional**
```go
func (rw *ResponseWriter) Write(b []byte) (int, error) {
    if !rw.headerWritten {
        rw.WriteHeader(http.StatusOK)  // Implicit 200
    }

    n, err := rw.ResponseWriter.Write(b)
    rw.bytesWritten += n
    return n, err
}
```

### Concept 5: Middleware Factories

**Pattern:** Functions that return middleware.

**Simple middleware:**
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    // ...
}
```

**Middleware factory:**
```go
func LoggingMiddleware(logger *log.Logger) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.Printf("Request: %s", r.URL.Path)
            next.ServeHTTP(w, r)
        })
    }
}

// Usage:
logger := log.New(os.Stdout, "[API] ", log.LstdFlags)
handler := LoggingMiddleware(logger)(actualHandler)
```

**Benefits:**
- **Configuration**: Pass options to middleware
- **Dependency injection**: Pass logger, DB, cache, etc.
- **Reusability**: Same middleware, different configs

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Conditional Middleware

Apply middleware only if condition is met.

```go
func ConditionalMiddleware(condition func(*http.Request) bool, middleware Middleware) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if condition(r) {
                middleware(next).ServeHTTP(w, r)
            } else {
                next.ServeHTTP(w, r)
            }
        })
    }
}

// Usage: Only apply auth to /api/* routes
authOnlyAPI := ConditionalMiddleware(
    func(r *http.Request) bool {
        return strings.HasPrefix(r.URL.Path, "/api/")
    },
    AuthMiddleware,
)
```

### Pattern 2: Per-Route Middleware

Different middleware for different routes.

```go
func SetupRoutes() http.Handler {
    mux := http.NewServeMux()

    // Public routes
    mux.Handle("/", Chain(handleHome, LoggingMiddleware))

    // API routes (auth + rate limit)
    mux.Handle("/api/", Chain(
        handleAPI,
        LoggingMiddleware,
        AuthMiddleware,
        RateLimitMiddleware,
    ))

    // Admin routes (auth + admin check)
    mux.Handle("/admin/", Chain(
        handleAdmin,
        LoggingMiddleware,
        AuthMiddleware,
        AdminOnlyMiddleware,
    ))

    return mux
}
```

### Pattern 3: Metrics Middleware

Track request metrics.

```go
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := NewResponseWriter(w)

        next.ServeHTTP(rw, r)

        duration := time.Since(start)

        // Record metrics
        httpRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
            strconv.Itoa(rw.StatusCode()),
        ).Observe(duration.Seconds())

        httpRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            strconv.Itoa(rw.StatusCode()),
        ).Inc()
    })
}
```

### Pattern 4: Request Timeout

Enforce timeout on request handling.

```go
func TimeoutMiddleware(timeout time.Duration) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()

            done := make(chan struct{})
            panicChan := make(chan interface{}, 1)

            go func() {
                defer func() {
                    if p := recover(); p != nil {
                        panicChan <- p
                    }
                }()

                next.ServeHTTP(w, r.WithContext(ctx))
                close(done)
            }()

            select {
            case p := <-panicChan:
                panic(p)
            case <-done:
                return
            case <-ctx.Done():
                http.Error(w, "Request Timeout", http.StatusGatewayTimeout)
                return
            }
        })
    }
}
```

### Pattern 5: Method Enforcement

Only allow specific HTTP methods.

```go
func MethodMiddleware(allowedMethods ...string) Middleware {
    allowed := make(map[string]bool)
    for _, method := range allowedMethods {
        allowed[method] = true
    }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !allowed[r.Method] {
                http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// Usage:
handler := MethodMiddleware("GET", "POST")(actualHandler)
```

### Pattern 6: Content-Type Enforcement

Require specific content type for POST/PUT requests.

```go
func RequireJSONMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" || r.Method == "PUT" {
            contentType := r.Header.Get("Content-Type")
            if contentType != "application/json" {
                http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
                return
            }
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## 7. Real-World Applications

### Production APIs

**Use case:** Every production REST API

Companies: Stripe, Twilio, GitHub, AWS

```go
// Stripe-like API setup
func SetupAPI() http.Handler {
    handler := Chain(
        router,
        RecoveryMiddleware,
        RequestIDMiddleware,
        LoggingMiddleware,
        CORSMiddleware("*"),
        AuthMiddleware,
        RateLimitMiddleware(100, 10), // 100 req/sec, burst 10
        MetricsMiddleware,
    )

    return handler
}
```

### Microservices

**Use case:** Service-to-service communication

Companies: Uber, Netflix, Airbnb

```go
// Internal service middleware
func InternalServiceMiddleware() http.Handler {
    return Chain(
        handler,
        RecoveryMiddleware,
        RequestIDPropagationMiddleware,  // Forward request ID
        ServiceAuthMiddleware,            // Verify service token
        TracingMiddleware,                // Distributed tracing
        MetricsMiddleware,
    )
}
```

### API Gateways

**Use case:** Single entry point for multiple services

Companies: Kong, AWS API Gateway

```go
func GatewayMiddleware() http.Handler {
    return Chain(
        proxyHandler,
        RecoveryMiddleware,
        LoggingMiddleware,
        RateLimitMiddleware(1000, 100),
        AuthMiddleware,
        CircuitBreakerMiddleware,
        RetryMiddleware,
        LoadBalancerMiddleware,
    )
}
```

### GraphQL Servers

**Use case:** GraphQL API servers

Companies: GitHub, Shopify

```go
func GraphQLMiddleware() http.Handler {
    return Chain(
        graphqlHandler,
        RecoveryMiddleware,
        CORSMiddleware("*"),
        ComplexityLimitMiddleware,  // Prevent expensive queries
        DepthLimitMiddleware,       // Prevent deeply nested queries
        AuthMiddleware,
        CachingMiddleware,
    )
}
```

### Webhooks

**Use case:** Receive webhooks from third parties

Companies: Stripe, Shopify, GitHub

```go
func WebhookMiddleware() http.Handler {
    return Chain(
        webhookHandler,
        RecoveryMiddleware,
        LoggingMiddleware,
        SignatureVerificationMiddleware,  // Verify webhook signature
        ReplayProtectionMiddleware,       // Prevent replay attacks
        RateLimitMiddleware(10, 5),
    )
}
```

---

## 8. Common Mistakes to Avoid

### Mistake 1: Ignoring Middleware Order

**❌ Wrong:**
```go
handler := Chain(
    actualHandler,
    LoggingMiddleware,
    RecoveryMiddleware,  // Recovery should be first!
)
```

**Problem:** If logging panics, recovery won't catch it.

**✅ Correct:**
```go
handler := Chain(
    actualHandler,
    RecoveryMiddleware,  // Outermost
    LoggingMiddleware,
)
```

### Mistake 2: Not Calling next.ServeHTTP

**❌ Wrong:**
```go
func BrokenMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Before")
        // Forgot next.ServeHTTP(w, r)!
        log.Println("After")
    })
}
```

**Problem:** Handler never runs. Request hangs.

**✅ Correct:**
```go
func WorkingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Before")
        next.ServeHTTP(w, r)
        log.Println("After")
    })
}
```

### Mistake 3: Modifying Request After Calling next

**❌ Wrong:**
```go
func BadMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)

        // Too late! Response already sent
        w.Header().Set("X-Custom", "value")
    })
}
```

**Problem:** Headers can't be set after response started.

**✅ Correct:**
```go
func GoodMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Custom", "value")  // Before next.ServeHTTP
        next.ServeHTTP(w, r)
    })
}
```

### Mistake 4: Context Value Type Assertions Without Checking

**❌ Wrong:**
```go
func handleUsers(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*User)  // Panic if not found!
    // ...
}
```

**Problem:** Panics if context value is missing or wrong type.

**✅ Correct:**
```go
func handleUsers(w http.ResponseWriter, r *http.Request) {
    user, ok := GetUser(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", 401)
        return
    }
    // ...
}
```

### Mistake 5: Sharing State Between Requests

**❌ Wrong:**
```go
func BrokenMiddleware(next http.Handler) http.Handler {
    var requestCount int  // Shared state!

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestCount++  // Race condition!
        log.Printf("Request #%d", requestCount)
        next.ServeHTTP(w, r)
    })
}
```

**Problem:** Data race. Multiple goroutines access `requestCount` concurrently.

**✅ Correct:**
```go
func WorkingMiddleware(next http.Handler) http.Handler {
    var (
        mu           sync.Mutex
        requestCount int
    )

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        mu.Lock()
        requestCount++
        count := requestCount
        mu.Unlock()

        log.Printf("Request #%d", count)
        next.ServeHTTP(w, r)
    })
}
```

### Mistake 6: Not Wrapping ResponseWriter Properly

**❌ Wrong:**
```go
type ResponseWriter struct {
    w          http.ResponseWriter
    statusCode int
}

func (rw *ResponseWriter) Header() http.Header {
    return rw.w.Header()
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
    return rw.w.Write(b)
}

func (rw *ResponseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.w.WriteHeader(code)
}
```

**Problem:** Doesn't handle interfaces like `http.Flusher`, `http.Hijacker`.

**✅ Correct:**
```go
type ResponseWriter struct {
    http.ResponseWriter  // Embed to get all methods
    statusCode int
}

// Only override methods you need to intercept
func (rw *ResponseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// Check for optional interfaces
func (rw *ResponseWriter) Flush() {
    if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
        flusher.Flush()
    }
}
```

---

## 9. Stretch Goals

### Goal 1: Implement Structured Logging Middleware ⭐⭐

Add structured JSON logging with request details.

**Hint:**
```go
import "go.uber.org/zap"

func StructuredLoggingMiddleware(logger *zap.Logger) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            rw := NewResponseWriter(w)

            next.ServeHTTP(rw, r)

            logger.Info("request completed",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.Int("status", rw.StatusCode()),
                zap.Duration("duration", time.Since(start)),
            )
        })
    }
}
```

### Goal 2: Add Distributed Tracing ⭐⭐⭐

Integrate OpenTelemetry for distributed tracing.

**Hint:**
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

### Goal 3: Circuit Breaker Middleware ⭐⭐⭐

Implement circuit breaker pattern in middleware.

**Hint:**
```go
func CircuitBreakerMiddleware(cb *CircuitBreaker) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !cb.Allow() {
                http.Error(w, "Service Unavailable", 503)
                return
            }

            rw := NewResponseWriter(w)
            next.ServeHTTP(rw, r)

            if rw.StatusCode() >= 500 {
                cb.RecordFailure()
            } else {
                cb.RecordSuccess()
            }
        })
    }
}
```

### Goal 4: Request Compression ⭐⭐

Add gzip compression middleware.

**Hint:**
```go
import "compress/gzip"

func CompressionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        w.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(w)
        defer gz.Close()

        gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
        next.ServeHTTP(gzw, r)
    })
}
```

### Goal 5: Request Validation Middleware ⭐⭐

Validate request schema before reaching handler.

**Hint:**
```go
func ValidationMiddleware(schema interface{}) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Method == "POST" || r.Method == "PUT" {
                var body interface{}
                if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                    http.Error(w, "Invalid JSON", 400)
                    return
                }

                if err := validate(body, schema); err != nil {
                    http.Error(w, err.Error(), 400)
                    return
                }

                // Re-encode body for handler
                encoded, _ := json.Marshal(body)
                r.Body = io.NopCloser(bytes.NewReader(encoded))
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## How to Run

```bash
# Run the demo
go run ./minis/37-http-middleware-chain/cmd/middleware-demo

# Test with curl
curl http://localhost:8080/
curl -H "Authorization: Bearer token" http://localhost:8080/protected

# Run tests
go test ./minis/37-http-middleware-chain/exercise

# Run with verbose output
go test -v ./minis/37-http-middleware-chain/exercise

# Run with race detector
go test -race ./minis/37-http-middleware-chain/exercise
```

---

## Summary

**What you learned:**
- ✅ Middleware pattern for cross-cutting concerns
- ✅ Function composition and higher-order functions
- ✅ Request/response wrapping and interception
- ✅ Context-based data passing between layers
- ✅ ResponseWriter wrapping to capture metadata
- ✅ Middleware chaining and execution order
- ✅ Common patterns (auth, logging, recovery, metrics)

**Why this matters:**
Middleware is the foundation of modern web frameworks. It enables clean separation of concerns, code reuse, and composable architecture. Every production API uses middleware for logging, auth, metrics, and more.

**Key patterns:**
- **Wrapper pattern**: `func(http.Handler) http.Handler`
- **Chain pattern**: Apply multiple middleware in order
- **Context pattern**: Pass request-scoped data
- **ResponseWriter pattern**: Intercept response metadata

**Next steps:**
- Build a full REST API with middleware
- Add authentication and authorization
- Implement rate limiting and circuit breakers
- Add distributed tracing

Build composable systems! 🚀
