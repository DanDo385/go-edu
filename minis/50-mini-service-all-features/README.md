# Project 50: Mini Service - All Features (Capstone Project)

## 1. What Is This About?

### The Ultimate Challenge

You've completed 49 mini-projects covering everything from basic strings to blockchain concepts. Now it's time to **put it all together** in a production-grade microservice that demonstrates mastery of Go.

This capstone project integrates:
- **Projects 1-5**: Core Go (strings, maps, CSV, JSON, file handling)
- **Projects 6-10**: Concurrency (worker pools, caching, HTTP, gRPC, graceful shutdown)
- **Projects 11-20**: Deep Go (slices, pointers, interfaces, contexts, channels, goroutines)
- **Projects 21-30**: Advanced concurrency (race detection, mutexes, atomics, sync patterns, profiling)
- **Projects 31-40**: Web & networking (HTTP servers, WebSockets, TCP, middleware, config, crypto)
- **Projects 41-49**: Specialized topics (blockchain, P2P, transactions, merkle trees)

### Real-World Scenario

You're building a **production microservice** for a fintech company that needs:

**Requirements:**
- Handle 10,000 requests per second
- Sub-100ms latency for 99th percentile
- Zero downtime deployments
- Comprehensive observability (logs, metrics, traces)
- Secure authentication and authorization
- Rate limiting to prevent abuse
- Database integration with connection pooling
- Graceful degradation when dependencies fail
- Health checks for load balancers
- Configuration via YAML and environment variables

**❌ Naive approach:**
```go
func main() {
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        // No logging
        // No metrics
        // No auth
        // No rate limiting
        // No graceful shutdown
        // Direct DB calls (no pooling)
        db.Query("SELECT * FROM users")
        // ...
    })
    http.ListenAndServe(":8080", nil) // Crashes on Ctrl+C
}
```

**✅ Production approach:**
```go
func main() {
    // Load config from YAML + env vars
    cfg := config.Load()

    // Initialize structured logging
    logger := setupLogger(cfg)

    // Initialize Prometheus metrics
    metrics := setupMetrics()

    // Initialize database with pooling
    db := setupDatabase(cfg)
    defer db.Close()

    // Build middleware chain
    handler := middleware.Chain(
        router,
        middleware.Recovery(logger),
        middleware.RequestID(),
        middleware.Logging(logger),
        middleware.Metrics(metrics),
        middleware.CORS(cfg.CORS),
        middleware.RateLimit(cfg.RateLimit),
    )

    // Create server with timeouts
    server := &http.Server{
        Addr:         cfg.Server.Addr,
        Handler:      handler,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    // Graceful shutdown
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    go func() {
        logger.Info().Msgf("Server starting on %s", cfg.Server.Addr)
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            logger.Fatal().Err(err).Msg("Server failed")
        }
    }()

    <-shutdown
    logger.Info().Msg("Shutting down gracefully...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        logger.Error().Err(err).Msg("Shutdown failed")
    }

    logger.Info().Msg("Server stopped")
}
```

### What You'll Build

A complete microservice with:

1. **Structured Logging** (zerolog)
   - JSON formatted logs
   - Request tracing with correlation IDs
   - Log levels (debug, info, warn, error)
   - Performance optimized (zero allocation)

2. **Prometheus Metrics**
   - HTTP request duration histogram
   - Request counter by status code
   - Active connections gauge
   - Custom business metrics

3. **Configuration Management**
   - YAML configuration files
   - Environment variable overrides
   - Validation and defaults
   - Hot reload support

4. **Graceful Shutdown**
   - Signal handling (SIGTERM, SIGINT)
   - Drain existing connections
   - Close database connections
   - Flush metrics and logs

5. **HTTP Middleware Chain**
   - Recovery (panic handler)
   - Request ID generation
   - Structured logging
   - Metrics collection
   - CORS handling
   - JWT authentication
   - Rate limiting

6. **Health & Readiness Endpoints**
   - `/health` - liveness probe
   - `/ready` - readiness probe (checks DB, dependencies)
   - `/metrics` - Prometheus metrics

7. **Database Integration**
   - Connection pooling
   - Prepared statements
   - Transaction support
   - Retry logic
   - Health checks

8. **JWT Authentication**
   - Token generation and validation
   - Claims extraction
   - Token refresh
   - User context propagation

9. **Request Tracing**
   - Distributed tracing with request IDs
   - Trace context propagation
   - Correlation across services

10. **Error Handling**
    - Structured error responses
    - Error wrapping and context
    - HTTP status code mapping
    - Client vs server error distinction

11. **Comprehensive Testing**
    - Unit tests
    - Integration tests
    - HTTP handler tests
    - Middleware tests
    - Mock database
    - Test fixtures

---

## 2. First Principles: Microservice Architecture

### What is a Microservice?

A **microservice** is a small, independently deployable service that does one thing well.

**Monolith vs Microservices:**

```
MONOLITH:
┌─────────────────────────────┐
│  Single Application         │
│  ┌─────────────────────┐   │
│  │ User Management     │   │
│  │ Order Processing    │   │
│  │ Payment System      │   │
│  │ Inventory           │   │
│  │ Notifications       │   │
│  └─────────────────────┘   │
│  All in one database        │
└─────────────────────────────┘

MICROSERVICES:
┌──────────┐  ┌──────────┐  ┌──────────┐
│  User    │  │  Order   │  │ Payment  │
│ Service  │  │ Service  │  │ Service  │
│  + DB    │  │  + DB    │  │  + DB    │
└──────────┘  └──────────┘  └──────────┘
     │             │              │
     └─────────────┴──────────────┘
              API Gateway
```

**Benefits:**
- **Independent deployment**: Update one service without touching others
- **Technology diversity**: Use different languages/databases per service
- **Scalability**: Scale only the services that need it
- **Resilience**: One service failure doesn't crash everything
- **Team autonomy**: Different teams own different services

**Challenges:**
- **Distributed system complexity**: Network calls, latency, partial failures
- **Data consistency**: No ACID transactions across services
- **Observability**: Logs and metrics scattered across services
- **Testing**: Need integration tests across services

### The 12-Factor App Methodology

Modern microservices follow **12-Factor App** principles:

1. **Codebase**: One codebase in version control
2. **Dependencies**: Explicitly declare dependencies (go.mod)
3. **Config**: Store config in environment variables
4. **Backing services**: Treat databases as attached resources
5. **Build/release/run**: Strict separation of stages
6. **Processes**: Execute as stateless processes
7. **Port binding**: Export services via port binding
8. **Concurrency**: Scale out via process model
9. **Disposability**: Fast startup and graceful shutdown
10. **Dev/prod parity**: Keep environments similar
11. **Logs**: Treat logs as event streams
12. **Admin processes**: Run admin tasks as one-off processes

**This project implements all 12 factors.**

### Observability: The Three Pillars

Production microservices need **observability**:

```
┌────────────────────────────────────┐
│         OBSERVABILITY              │
├────────────┬───────────┬───────────┤
│   LOGS     │  METRICS  │  TRACES   │
└────────────┴───────────┴───────────┘

LOGS:
  What happened?
  [INFO] User 123 logged in
  [ERROR] Database connection failed

  Use: Debugging, audit trails

METRICS:
  How many? How fast?
  http_requests_total{status="200"} 1523
  http_request_duration_seconds{p99} 0.095

  Use: Alerting, capacity planning

TRACES:
  Where is the time spent?
  Request → API Gateway → User Service → Database
           50ms          30ms           20ms

  Use: Performance optimization, dependency mapping
```

**This project implements all three pillars.**

### Middleware Architecture

**Middleware** wraps handlers to add cross-cutting concerns:

```
Request
  ↓
┌─────────────────┐
│ Recovery        │ ← Catch panics
└─────────────────┘
  ↓
┌─────────────────┐
│ Request ID      │ ← Generate correlation ID
└─────────────────┘
  ↓
┌─────────────────┐
│ Logging         │ ← Log request/response
└─────────────────┘
  ↓
┌─────────────────┐
│ Metrics         │ ← Record metrics
└─────────────────┘
  ↓
┌─────────────────┐
│ CORS            │ ← Set CORS headers
└─────────────────┘
  ↓
┌─────────────────┐
│ Auth (optional) │ ← Verify JWT
└─────────────────┘
  ↓
┌─────────────────┐
│ Rate Limit      │ ← Check rate limits
└─────────────────┘
  ↓
┌─────────────────┐
│ Business Logic  │ ← Your handler
└─────────────────┘
  ↓
Response
```

**Each middleware:**
- Runs before handler (top to bottom)
- Runs after handler (bottom to top)
- Can short-circuit the chain
- Can modify request/response

### Graceful Shutdown

**Why graceful shutdown matters:**

**Without graceful shutdown:**
```
1. SIGTERM received
2. Server stops immediately
3. In-flight requests dropped
4. Database connections left open
5. Metrics not flushed
```

**With graceful shutdown:**
```
1. SIGTERM received
2. Stop accepting new requests
3. Wait for in-flight requests to complete (with timeout)
4. Close database connections
5. Flush logs and metrics
6. Exit cleanly
```

**Kubernetes integration:**
```
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: myservice
    lifecycle:
      preStop:
        exec:
          command: ["/bin/sh", "-c", "sleep 5"]
    # Kubernetes sends SIGTERM, waits 30s (terminationGracePeriodSeconds)
    terminationGracePeriodSeconds: 30
```

### Database Integration Patterns

**Connection pooling:**
```go
db, _ := sql.Open("postgres", connStr)

// Configure pool
db.SetMaxOpenConns(25)        // Max connections
db.SetMaxIdleConns(5)         // Keep idle connections
db.SetConnMaxLifetime(5 * time.Minute)  // Recycle connections
```

**Why pooling?**
- Opening connections is expensive (TCP + TLS + auth)
- Reuse connections across requests
- Limit total connections (database has limits)
- Automatic retry and health checks

**Prepared statements:**
```go
// ❌ SQL injection risk
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)
db.Query(query)

// ✅ Safe with placeholders
stmt, _ := db.Prepare("SELECT * FROM users WHERE id = $1")
defer stmt.Close()
stmt.Query(userID)
```

**Transactions:**
```go
tx, _ := db.Begin()
defer tx.Rollback()  // Rollback if not committed

tx.Exec("INSERT INTO orders ...")
tx.Exec("UPDATE inventory ...")

tx.Commit()  // Atomic commit
```

---

## 3. Architecture Overview

### Project Structure

```
50-mini-service-all-features/
├── README.md
├── go.mod
├── go.sum
├── config.yaml                    # Configuration file
├── cmd/
│   └── service/
│       └── main.go                # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Config loading & validation
│   ├── middleware/
│   │   ├── middleware.go          # Middleware chain
│   │   ├── recovery.go            # Panic recovery
│   │   ├── requestid.go           # Request ID generation
│   │   ├── logging.go             # Request/response logging
│   │   ├── metrics.go             # Prometheus metrics
│   │   ├── cors.go                # CORS headers
│   │   ├── auth.go                # JWT authentication
│   │   └── ratelimit.go           # Rate limiting
│   ├── handlers/
│   │   ├── health.go              # Health check endpoints
│   │   ├── auth.go                # Authentication handlers
│   │   └── users.go               # User CRUD handlers
│   ├── models/
│   │   ├── user.go                # User model
│   │   └── errors.go              # Error types
│   ├── database/
│   │   ├── database.go            # DB connection & pooling
│   │   └── users.go               # User repository
│   └── metrics/
│       └── metrics.go             # Prometheus metrics registry
└── exercise/
    ├── exercise.go                # Skeleton for students
    ├── exercise_test.go           # Tests
    └── solution.go                # Complete solution
```

### Component Interactions

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ↓
┌─────────────────────────────────┐
│      HTTP Server                │
│  (graceful shutdown support)    │
└──────┬──────────────────────────┘
       │
       ↓
┌─────────────────────────────────┐
│    Middleware Chain             │
│  Recovery → RequestID →         │
│  Logging → Metrics → CORS →     │
│  Auth → RateLimit               │
└──────┬──────────────────────────┘
       │
       ↓
┌─────────────────────────────────┐
│       HTTP Handlers             │
│  /health  /ready  /metrics      │
│  /login   /users                │
└──────┬──────────────────────────┘
       │
       ↓
┌─────────────────────────────────┐
│    Business Logic               │
│  (models, validation)           │
└──────┬──────────────────────────┘
       │
       ↓
┌─────────────────────────────────┐
│      Database Layer             │
│  (connection pool, queries)     │
└─────────────────────────────────┘
```

### Request Flow Example

```
1. Client sends: GET /users/123
   Header: Authorization: Bearer eyJhbG...

2. Recovery middleware: Set up panic recovery

3. RequestID middleware: Generate request ID "req-abc-123"

4. Logging middleware: Log "→ GET /users/123 [req-abc-123]"

5. Metrics middleware: Start timer

6. CORS middleware: Set Access-Control-Allow-Origin

7. Auth middleware:
   - Extract JWT from Authorization header
   - Validate signature
   - Extract user claims
   - Store in context

8. RateLimit middleware:
   - Check user's rate limit
   - If exceeded, return 429
   - Else, allow through

9. Handler:
   - Get user ID from URL
   - Query database
   - Return JSON response

10. Metrics middleware: Record duration (42ms), status (200)

11. Logging middleware: Log "← GET /users/123 [req-abc-123] 200 42ms"

12. Response sent to client
```

---

## 4. Complete Implementation Walkthrough

### Step 1: Configuration Management

**config.yaml:**
```yaml
server:
  addr: ":8080"
  read_timeout: 10s
  write_timeout: 10s
  shutdown_timeout: 30s

database:
  host: localhost
  port: 5432
  name: myapp
  user: postgres
  password: postgres
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

logging:
  level: info  # debug, info, warn, error
  format: json # json, console

cors:
  allowed_origins:
    - "*"
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
  allowed_headers:
    - Content-Type
    - Authorization

rate_limit:
  requests_per_second: 100
  burst: 10

jwt:
  secret: your-secret-key
  expiration: 24h
```

**internal/config/config.go:**
```go
package config

import (
    "fmt"
    "os"
    "time"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Server    ServerConfig    `yaml:"server"`
    Database  DatabaseConfig  `yaml:"database"`
    Logging   LoggingConfig   `yaml:"logging"`
    CORS      CORSConfig      `yaml:"cors"`
    RateLimit RateLimitConfig `yaml:"rate_limit"`
    JWT       JWTConfig       `yaml:"jwt"`
}

type ServerConfig struct {
    Addr            string        `yaml:"addr"`
    ReadTimeout     time.Duration `yaml:"read_timeout"`
    WriteTimeout    time.Duration `yaml:"write_timeout"`
    ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
    Host            string        `yaml:"host"`
    Port            int           `yaml:"port"`
    Name            string        `yaml:"name"`
    User            string        `yaml:"user"`
    Password        string        `yaml:"password"`
    MaxOpenConns    int           `yaml:"max_open_conns"`
    MaxIdleConns    int           `yaml:"max_idle_conns"`
    ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type LoggingConfig struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
}

type CORSConfig struct {
    AllowedOrigins []string `yaml:"allowed_origins"`
    AllowedMethods []string `yaml:"allowed_methods"`
    AllowedHeaders []string `yaml:"allowed_headers"`
}

type RateLimitConfig struct {
    RequestsPerSecond float64 `yaml:"requests_per_second"`
    Burst             int     `yaml:"burst"`
}

type JWTConfig struct {
    Secret     string        `yaml:"secret"`
    Expiration time.Duration `yaml:"expiration"`
}

// Load reads config from file and env vars
func Load(configPath string) (*Config, error) {
    // Read YAML file
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("read config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }

    // Override with environment variables
    if addr := os.Getenv("SERVER_ADDR"); addr != "" {
        cfg.Server.Addr = addr
    }
    if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
        cfg.Database.Host = dbHost
    }
    if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
        cfg.Database.Password = dbPass
    }
    if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
        cfg.JWT.Secret = jwtSecret
    }

    // Validate
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("validate config: %w", err)
    }

    return &cfg, nil
}

func (c *Config) Validate() error {
    if c.Server.Addr == "" {
        return fmt.Errorf("server.addr is required")
    }
    if c.JWT.Secret == "" {
        return fmt.Errorf("jwt.secret is required")
    }
    if c.Database.Host == "" {
        return fmt.Errorf("database.host is required")
    }
    return nil
}

func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        c.Host, c.Port, c.User, c.Password, c.Name,
    )
}
```

**Key concepts:**
- YAML for human-readable config
- Environment variables override YAML (12-factor app)
- Validation ensures required fields are set
- DSN builder for database connection string

### Step 2: Structured Logging with Zerolog

**Setup logger:**
```go
package main

import (
    "os"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func setupLogger(cfg *config.LoggingConfig) zerolog.Logger {
    // Set log level
    level, err := zerolog.ParseLevel(cfg.Level)
    if err != nil {
        level = zerolog.InfoLevel
    }
    zerolog.SetGlobalLevel(level)

    // Configure output format
    var logger zerolog.Logger
    if cfg.Format == "console" {
        logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
    } else {
        logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
    }

    return logger
}
```

**Usage in middleware:**
```go
func LoggingMiddleware(logger zerolog.Logger) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            requestID := GetRequestID(r.Context())

            // Log request
            logger.Info().
                Str("request_id", requestID).
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Str("remote_addr", r.RemoteAddr).
                Msg("request started")

            // Wrap response writer
            rw := NewResponseWriter(w)

            next.ServeHTTP(rw, r)

            // Log response
            duration := time.Since(start)
            logger.Info().
                Str("request_id", requestID).
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Int("status", rw.StatusCode()).
                Int("bytes", rw.BytesWritten()).
                Dur("duration", duration).
                Msg("request completed")
        })
    }
}
```

**Log output (JSON):**
```json
{"level":"info","request_id":"req-abc-123","method":"GET","path":"/users/123","remote_addr":"192.168.1.1","time":"2025-01-15T10:30:00Z","message":"request started"}
{"level":"info","request_id":"req-abc-123","method":"GET","path":"/users/123","status":200,"bytes":156,"duration":42.3,"time":"2025-01-15T10:30:00Z","message":"request completed"}
```

### Step 3: Prometheus Metrics

**internal/metrics/metrics.go:**
```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
    HTTPRequestsTotal   *prometheus.CounterVec
    HTTPRequestDuration *prometheus.HistogramVec
    HTTPActiveRequests  prometheus.Gauge
}

func New() *Metrics {
    return &Metrics{
        HTTPRequestsTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),
        HTTPRequestDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "http_request_duration_seconds",
                Help:    "HTTP request duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "path", "status"},
        ),
        HTTPActiveRequests: promauto.NewGauge(
            prometheus.GaugeOpts{
                Name: "http_active_requests",
                Help: "Number of active HTTP requests",
            },
        ),
    }
}
```

**Metrics middleware:**
```go
func MetricsMiddleware(m *metrics.Metrics) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Increment active requests
            m.HTTPActiveRequests.Inc()
            defer m.HTTPActiveRequests.Dec()

            // Wrap response writer
            rw := NewResponseWriter(w)

            next.ServeHTTP(rw, r)

            // Record metrics
            duration := time.Since(start).Seconds()
            status := strconv.Itoa(rw.StatusCode())

            m.HTTPRequestsTotal.WithLabelValues(
                r.Method, r.URL.Path, status,
            ).Inc()

            m.HTTPRequestDuration.WithLabelValues(
                r.Method, r.URL.Path, status,
            ).Observe(duration)
        })
    }
}
```

**Prometheus metrics output:**
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/users",status="200"} 1523

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/users",status="200",le="0.005"} 100
http_request_duration_seconds_bucket{method="GET",path="/users",status="200",le="0.01"} 450
http_request_duration_seconds_sum{method="GET",path="/users",status="200"} 64.3
http_request_duration_seconds_count{method="GET",path="/users",status="200"} 1523

# HELP http_active_requests Number of active HTTP requests
# TYPE http_active_requests gauge
http_active_requests 12
```

### Step 4: Middleware Chain

**internal/middleware/middleware.go:**
```go
package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Chain applies middleware in order
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
    // Apply in reverse so first middleware wraps all others
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}
```

**Recovery middleware:**
```go
func Recovery(logger zerolog.Logger) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    logger.Error().
                        Interface("panic", err).
                        Bytes("stack", debug.Stack()).
                        Msg("panic recovered")

                    http.Error(w, "Internal Server Error",
                        http.StatusInternalServerError)
                }
            }()

            next.ServeHTTP(w, r)
        })
    }
}
```

**Request ID middleware:**
```go
type contextKey string

const requestIDKey contextKey = "request_id"

func RequestID() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Generate unique ID
            requestID := uuid.New().String()

            // Store in context
            ctx := context.WithValue(r.Context(), requestIDKey, requestID)

            // Add to response header
            w.Header().Set("X-Request-ID", requestID)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}
```

**CORS middleware:**
```go
func CORS(cfg config.CORSConfig) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Set CORS headers
            if len(cfg.AllowedOrigins) > 0 {
                w.Header().Set("Access-Control-Allow-Origin",
                    cfg.AllowedOrigins[0])
            }

            w.Header().Set("Access-Control-Allow-Methods",
                strings.Join(cfg.AllowedMethods, ", "))

            w.Header().Set("Access-Control-Allow-Headers",
                strings.Join(cfg.AllowedHeaders, ", "))

            // Handle preflight
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**Rate limiting middleware:**
```go
func RateLimit(cfg config.RateLimitConfig) Middleware {
    limiter := rate.NewLimiter(
        rate.Limit(cfg.RequestsPerSecond),
        cfg.Burst,
    )

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Too Many Requests",
                    http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### Step 5: JWT Authentication

**internal/middleware/auth.go:**
```go
package middleware

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Exp      int64  `json:"exp"`
}

// GenerateJWT creates a JWT token
func GenerateJWT(userID int, username, secret string, expiration time.Duration) (string, error) {
    // Create header
    header := map[string]string{
        "alg": "HS256",
        "typ": "JWT",
    }
    headerJSON, _ := json.Marshal(header)
    headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

    // Create claims
    claims := Claims{
        UserID:   userID,
        Username: username,
        Exp:      time.Now().Add(expiration).Unix(),
    }
    claimsJSON, _ := json.Marshal(claims)
    claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

    // Create signature
    message := headerB64 + "." + claimsB64
    signature := createSignature(message, secret)

    // Combine
    token := message + "." + signature
    return token, nil
}

// ValidateJWT verifies and parses a JWT token
func ValidateJWT(token, secret string) (*Claims, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid token format")
    }

    // Verify signature
    message := parts[0] + "." + parts[1]
    expectedSignature := createSignature(message, secret)
    if parts[2] != expectedSignature {
        return nil, fmt.Errorf("invalid signature")
    }

    // Parse claims
    claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return nil, fmt.Errorf("decode claims: %w", err)
    }

    var claims Claims
    if err := json.Unmarshal(claimsJSON, &claims); err != nil {
        return nil, fmt.Errorf("parse claims: %w", err)
    }

    // Check expiration
    if time.Now().Unix() > claims.Exp {
        return nil, fmt.Errorf("token expired")
    }

    return &claims, nil
}

func createSignature(message, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(message))
    return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

const userKey contextKey = "user"

// Auth middleware validates JWT
func Auth(secret string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")
            if token == authHeader {
                http.Error(w, "Invalid authorization format",
                    http.StatusUnauthorized)
                return
            }

            // Validate token
            claims, err := ValidateJWT(token, secret)
            if err != nil {
                http.Error(w, "Invalid token: "+err.Error(),
                    http.StatusUnauthorized)
                return
            }

            // Store user in context
            ctx := context.WithValue(r.Context(), userKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUser(ctx context.Context) (*Claims, bool) {
    claims, ok := ctx.Value(userKey).(*Claims)
    return claims, ok
}
```

### Step 6: Database Integration

**internal/database/database.go:**
```go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
}

func New(cfg config.DatabaseConfig) (*DB, error) {
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

    // Ping to verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping database: %w", err)
    }

    return &DB{db}, nil
}

// Health checks database connectivity
func (db *DB) Health(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
    defer cancel()

    return db.PingContext(ctx)
}
```

**User repository:**
```go
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func (db *DB) GetUser(ctx context.Context, id int) (*User, error) {
    query := `
        SELECT id, username, email, created_at
        FROM users
        WHERE id = $1
    `

    var user User
    err := db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, fmt.Errorf("query user: %w", err)
    }

    return &user, nil
}

func (db *DB) ListUsers(ctx context.Context) ([]User, error) {
    query := `
        SELECT id, username, email, created_at
        FROM users
        ORDER BY id
    `

    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("query users: %w", err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
            return nil, fmt.Errorf("scan user: %w", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return users, nil
}
```

### Step 7: HTTP Handlers

**Health check handlers:**
```go
package handlers

func Health(logger zerolog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "healthy",
        })
    }
}

func Ready(db *database.DB, logger zerolog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check database
        if err := db.Health(r.Context()); err != nil {
            logger.Error().Err(err).Msg("database health check failed")
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]string{
                "status": "not ready",
                "reason": "database unavailable",
            })
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "ready",
        })
    }
}
```

**User handlers:**
```go
func GetUser(db *database.DB, logger zerolog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get user ID from URL
        idStr := r.URL.Path[len("/users/"):]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        // Query database
        user, err := db.GetUser(r.Context(), id)
        if err != nil {
            logger.Error().Err(err).Int("user_id", id).Msg("get user failed")
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        // Return JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}

func ListUsers(db *database.DB, logger zerolog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        users, err := db.ListUsers(r.Context())
        if err != nil {
            logger.Error().Err(err).Msg("list users failed")
            http.Error(w, "Internal Server Error",
                http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}
```

### Step 8: Graceful Shutdown

**cmd/service/main.go:**
```go
func main() {
    // Load config
    cfg, err := config.Load("config.yaml")
    if err != nil {
        log.Fatal().Err(err).Msg("load config failed")
    }

    // Setup logger
    logger := setupLogger(&cfg.Logging)

    // Setup metrics
    m := metrics.New()

    // Setup database
    db, err := database.New(cfg.Database)
    if err != nil {
        logger.Fatal().Err(err).Msg("database connection failed")
    }
    defer db.Close()

    // Setup router
    router := setupRouter(cfg, db, logger, m)

    // Create server
    server := &http.Server{
        Addr:         cfg.Server.Addr,
        Handler:      router,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    // Start server in goroutine
    go func() {
        logger.Info().Msgf("Server starting on %s", cfg.Server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal().Err(err).Msg("server failed")
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    logger.Info().Msg("Shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(
        context.Background(),
        cfg.Server.ShutdownTimeout,
    )
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        logger.Error().Err(err).Msg("server shutdown failed")
    }

    logger.Info().Msg("Server stopped")
}

func setupRouter(
    cfg *config.Config,
    db *database.DB,
    logger zerolog.Logger,
    m *metrics.Metrics,
) http.Handler {
    mux := http.NewServeMux()

    // Health endpoints (no auth)
    mux.HandleFunc("/health", handlers.Health(logger))
    mux.HandleFunc("/ready", handlers.Ready(db, logger))
    mux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

    // Auth endpoints (no auth)
    mux.HandleFunc("/login", handlers.Login(db, cfg.JWT, logger))

    // Protected endpoints (auth required)
    protectedMux := http.NewServeMux()
    protectedMux.HandleFunc("/users", handlers.ListUsers(db, logger))
    protectedMux.HandleFunc("/users/", handlers.GetUser(db, logger))

    protected := middleware.Chain(
        protectedMux,
        middleware.Auth(cfg.JWT.Secret),
    )
    mux.Handle("/users", protected)
    mux.Handle("/users/", protected)

    // Apply global middleware
    handler := middleware.Chain(
        mux,
        middleware.Recovery(logger),
        middleware.RequestID(),
        middleware.Logging(logger),
        middleware.Metrics(m),
        middleware.CORS(cfg.CORS),
        middleware.RateLimit(cfg.RateLimit),
    )

    return handler
}
```

---

## 5. Key Patterns Demonstrated

### Pattern 1: Dependency Injection

**Don't:**
```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    db := getGlobalDB()  // Global state, hard to test
    logger := getGlobalLogger()
    // ...
}
```

**Do:**
```go
func GetUser(db *database.DB, logger zerolog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Use injected dependencies
        // Easy to test with mocks
    }
}
```

### Pattern 2: Context Propagation

**Don't:**
```go
var globalRequestID string

func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        globalRequestID = uuid.New().String()  // Race condition!
        next.ServeHTTP(w, r)
    })
}
```

**Do:**
```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := uuid.New().String()
        ctx := context.WithValue(r.Context(), requestIDKey, id)
        next.ServeHTTP(w, r.WithContext(ctx))  // Thread-safe
    })
}
```

### Pattern 3: Graceful Degradation

**Don't:**
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    cache.Get(key)  // Panic if cache is down
    // Service crashes
}
```

**Do:**
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    val, err := cache.Get(key)
    if err != nil {
        logger.Warn().Err(err).Msg("cache miss, using database")
        val = db.Get(key)  // Fallback
    }
    // Service continues
}
```

### Pattern 4: Structured Errors

**Don't:**
```go
if err != nil {
    http.Error(w, err.Error(), 500)  // Exposes internal errors
}
```

**Do:**
```go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

if err != nil {
    logger.Error().Err(err).Msg("internal error")
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(AppError{
        Code:    "INTERNAL_ERROR",
        Message: "An error occurred",  // Safe message
    })
}
```

### Pattern 5: Request Timeouts

**Don't:**
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    result := expensiveOperation()  // Might take forever
    // ...
}
```

**Do:**
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    result, err := expensiveOperationWithContext(ctx)
    if err == context.DeadlineExceeded {
        http.Error(w, "Request timeout", http.StatusGatewayTimeout)
        return
    }
    // ...
}
```

---

## 6. Testing Strategy

### Unit Tests

**Testing middleware:**
```go
func TestRecoveryMiddleware(t *testing.T) {
    logger := zerolog.Nop()  // No-op logger for tests

    handler := middleware.Recovery(logger)(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            panic("test panic")
        }),
    )

    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()

    handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusInternalServerError {
        t.Errorf("expected 500, got %d", rec.Code)
    }
}
```

### Integration Tests

**Testing with mock database:**
```go
type MockDB struct {
    users map[int]*User
}

func (m *MockDB) GetUser(ctx context.Context, id int) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}

func TestGetUserHandler(t *testing.T) {
    db := &MockDB{
        users: map[int]*User{
            1: {ID: 1, Username: "alice"},
        },
    }

    handler := handlers.GetUser(db, zerolog.Nop())

    req := httptest.NewRequest("GET", "/users/1", nil)
    rec := httptest.NewRecorder()

    handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rec.Code)
    }

    var user User
    json.NewDecoder(rec.Body).Decode(&user)

    if user.Username != "alice" {
        t.Errorf("expected alice, got %s", user.Username)
    }
}
```

---

## 7. Deployment Considerations

### Docker

**Dockerfile:**
```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o service ./cmd/service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/service .
COPY config.yaml .

EXPOSE 8080
CMD ["./service"]
```

### Kubernetes

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myservice
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myservice
  template:
    metadata:
      labels:
        app: myservice
    spec:
      containers:
      - name: myservice
        image: myservice:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres-service
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: myservice-secrets
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
```

---

## 8. Performance Benchmarks

Expected performance on modern hardware:

```
CPU: 4 cores @ 2.5 GHz
RAM: 8 GB

Throughput:  10,000 req/sec
Latency p50: 5ms
Latency p95: 15ms
Latency p99: 50ms

Memory usage: ~50 MB
CPU usage: ~20%
```

**Optimization tips:**
1. Use connection pooling (avoid per-request connections)
2. Enable HTTP/2 (multiplexing)
3. Use zerolog (zero allocation logging)
4. Minimize middleware (each adds ~1-2ms)
5. Use prepared statements (avoid SQL parsing)
6. Enable gzip compression for large responses
7. Add caching layer (Redis) for hot data

---

## 9. Stretch Goals

### Goal 1: Add Distributed Tracing ⭐⭐⭐

Integrate OpenTelemetry for request tracing.

**Hint:** Use `go.opentelemetry.io/otel` to trace requests across services.

### Goal 2: Add Redis Caching ⭐⭐

Cache database queries in Redis.

**Hint:** Check cache before database, set TTL on cached items.

### Goal 3: Add Background Jobs ⭐⭐⭐

Use worker pool for async tasks.

**Hint:** Create job queue with channels, process in goroutines.

### Goal 4: Add GraphQL API ⭐⭐⭐⭐

Expose GraphQL endpoint alongside REST.

**Hint:** Use `github.com/graphql-go/graphql` library.

### Goal 5: Add Multi-Tenancy ⭐⭐⭐⭐

Support multiple tenants with data isolation.

**Hint:** Extract tenant ID from JWT, filter all queries by tenant.

---

## 10. How to Run

```bash
# Install dependencies
go mod download

# Run the service
go run ./cmd/service

# In another terminal, test endpoints

# Health check
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready

# Login (get JWT token)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret"}'

# Get users (requires auth)
curl http://localhost:8080/users \
  -H "Authorization: Bearer YOUR_TOKEN"

# Metrics
curl http://localhost:8080/metrics

# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

---

## 11. Summary

**Congratulations!** You've built a production-grade microservice that demonstrates:

✅ **Structured logging** with zerolog (zero allocation)
✅ **Prometheus metrics** (counters, histograms, gauges)
✅ **Configuration** (YAML + env vars, validation)
✅ **Graceful shutdown** (drain connections, close resources)
✅ **Middleware chain** (recovery, logging, metrics, CORS, auth, rate limit)
✅ **Health checks** (/health, /ready)
✅ **Database integration** (connection pooling, prepared statements)
✅ **JWT authentication** (token generation, validation, claims)
✅ **Request tracing** (correlation IDs)
✅ **Error handling** (structured errors, proper status codes)
✅ **Comprehensive testing** (unit, integration, mocks)

**What makes this production-ready?**
- Observability (logs, metrics, traces)
- Resilience (graceful shutdown, error handling)
- Security (JWT auth, rate limiting, input validation)
- Performance (connection pooling, zero-allocation logging)
- Operability (health checks, metrics, configuration)

**You've mastered:**
- All 49 previous projects synthesized into one
- Production patterns used at companies like Google, Netflix, Uber
- 12-factor app methodology
- Cloud-native microservice architecture

**Next steps:**
1. Deploy to Kubernetes
2. Add observability stack (Prometheus, Grafana, Jaeger)
3. Implement CI/CD pipeline
4. Load test with tools like `k6` or `vegeta`
5. Build more microservices and connect them

You're now ready for production Go development! 🚀

**This is the culmination of your Go journey. Keep building amazing things!**
