# Project 26: sync.Once and the Singleton Pattern

## 1. What Is This About?

### Real-World Scenario

Imagine you're building an application that needs to:
- Connect to a database (expensive operation, should only happen once)
- Load configuration from a file (slow I/O, no need to repeat)
- Initialize a logging system (single global logger needed)
- Set up a connection pool (resource-intensive, one-time setup)

**❌ Without proper initialization:**
```go
var db *Database

func GetDB() *Database {
    if db == nil {
        db = ConnectToDatabase() // RACE CONDITION!
        // Multiple goroutines might all see db==nil and initialize multiple times!
    }
    return db
}
```

**Problem**: When multiple goroutines call `GetDB()` simultaneously:
1. Goroutine A checks `db == nil` (true)
2. Goroutine B checks `db == nil` (true) ← **RACE!**
3. Both goroutines initialize the database connection
4. Resources leak, connections multiply, chaos ensues

**✅ With sync.Once:**
```go
var (
    db   *Database
    once sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        db = ConnectToDatabase() // Guaranteed to run EXACTLY once
    })
    return db
}
```

**Result**: No matter how many goroutines call `GetDB()`:
- The initialization runs **exactly once**
- All goroutines wait for initialization to complete
- No race conditions, no duplicate work
- Thread-safe, efficient, idempotent

This project teaches you **sync.Once**, Go's elegant solution for one-time initialization, and the **singleton pattern**, a design pattern for ensuring only one instance of a type exists.

### What You'll Learn

1. **sync.Once mechanics**: How it guarantees exactly-once execution
2. **Singleton pattern**: Creating single instances with lazy initialization
3. **Lazy initialization**: Deferring expensive work until needed
4. **Idempotent guarantees**: Operations that are safe to attempt multiple times
5. **Thread safety**: Preventing race conditions in initialization
6. **Real-world patterns**: Configuration loading, connection pools, caches

### The Challenge

Build thread-safe singletons with:
- Lazy initialization using `sync.Once`
- Configuration loaders that read files only once
- Database connection managers
- Metrics collectors with singleton instances
- Proper error handling in initialization

---

## 2. First Principles: Understanding Initialization

### What Is Initialization?

**Initialization** is the process of setting up a resource before first use:
- Allocating memory
- Opening file handles
- Establishing network connections
- Reading configuration
- Computing derived values

**The fundamental problem**: Initialization often has costs:
- **Time**: Database connections take milliseconds or seconds
- **Resources**: File handles, network sockets are limited
- **Correctness**: Some resources should only exist once (global config, singleton caches)

### Types of Initialization

**1. Eager Initialization (at program start)**
```go
var config = LoadConfig() // Runs when program starts

func main() {
    // config is already loaded
}
```

**Pros**: Simple, no race conditions
**Cons**: Slow startup, wastes resources if never used

**2. Lazy Initialization (on first use)**
```go
var config *Config

func GetConfig() *Config {
    if config == nil {
        config = LoadConfig() // Runs on first call
    }
    return config
}
```

**Pros**: Fast startup, only initialize what you use
**Cons**: Race conditions in multi-threaded code!

**3. Thread-Safe Lazy Initialization (sync.Once)**
```go
var (
    config *Config
    once   sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        config = LoadConfig()
    })
    return config
}
```

**Pros**: Fast startup, thread-safe, exactly-once guarantee
**Cons**: Slightly more complex

### Why sync.Once Exists

**Problem Statement**: In a concurrent program, you want to:
1. Initialize a resource **exactly once**
2. Block other goroutines until initialization completes
3. Guarantee that all goroutines see the initialized value
4. Do this with minimal overhead

**Before sync.Once**, you'd write:
```go
var (
    config *Config
    mu     sync.Mutex
    loaded bool
)

func GetConfig() *Config {
    mu.Lock()
    if !loaded {
        config = LoadConfig()
        loaded = true
    }
    mu.Unlock()
    return config
}
```

**Problems**:
- Lock is acquired on **every call** (expensive)
- Manual bookkeeping with `loaded` flag
- Easy to get wrong (forget the flag, double-check locking bugs)

**With sync.Once**:
```go
var (
    config *Config
    once   sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        config = LoadConfig()
    })
    return config
}
```

**Advantages**:
- Lock is only acquired during first call (fast path is lock-free)
- Built-in bookkeeping (no manual flag)
- Impossible to misuse (the API enforces correctness)

---

## 3. sync.Once Deep Dive

### The Interface

```go
type Once struct {
    // Has unexported fields
}

func (o *Once) Do(f func())
```

**That's it!** One method: `Do(f func())`

**Contract**:
- The function `f` will be called **at most once**
- If multiple goroutines call `Do` simultaneously:
  - One goroutine runs `f`
  - All other goroutines **block** until `f` returns
  - After `f` returns, all subsequent calls to `Do` return immediately
- `f` is executed exactly once **even if it panics**
  - If `f` panics, `Once` still considers it "done"
  - Subsequent calls to `Do` will not call `f`

### How It Works (Simplified)

**Internal structure** (conceptual):
```go
type Once struct {
    done uint32       // Atomic flag: 0 = not done, 1 = done
    m    sync.Mutex   // Mutex for synchronization
}
```

**Algorithm** (fast path + slow path):
```go
func (o *Once) Do(f func()) {
    // Fast path: check if already done (no lock!)
    if atomic.LoadUint32(&o.done) == 1 {
        return // Already initialized, return immediately
    }

    // Slow path: not done yet, need synchronization
    o.m.Lock()
    defer o.m.Unlock()

    // Double-check (another goroutine might have initialized while we waited for lock)
    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1) // Mark as done AFTER f returns
        f() // Execute the initialization
    }
}
```

**Why is this fast?**
- **After first initialization**: Only an atomic load (no lock!)
- **During first initialization**: Lock protects against races
- **Fast path dominates**: 99.99% of calls take the fast path

### Guarantees

**1. Exactly-Once Execution**
```go
var (
    counter int
    once    sync.Once
)

func increment() {
    once.Do(func() {
        counter++
    })
}

// Call increment() from 1000 goroutines
// counter will be exactly 1 (not 1000)
```

**2. Happens-Before Guarantee**

From Go memory model:
> A single call of f() from once.Do(f) happens before any call of once.Do(f) returns.

**Translation**:
- If goroutine A runs `f`, and goroutine B calls `once.Do(f)` later:
  - B sees all writes from A's execution of `f`
  - No need for additional synchronization

**Example**:
```go
var (
    config *Config
    once   sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        config = &Config{Setting: "value"} // Write
    })
    return config // Read - guaranteed to see the write
}
```

**3. Panic Handling**
```go
var once sync.Once

func init() {
    once.Do(func() {
        panic("initialization failed")
    })
}

// Even though f panicked, once considers it "done"
// Calling once.Do(f) again will NOT call f
```

**Why?** If `Once` retried after panic, you could get into infinite panic loops. The design philosophy: fail fast, don't retry automatically.

---

## 4. The Singleton Pattern

### What Is a Singleton?

**Definition**: A design pattern that ensures:
1. A class/type has **exactly one instance**
2. Provides a **global access point** to that instance

**Real-world analogies**:
- President of a country (only one at a time)
- Sun in the solar system (only one)
- Database connection pool (single shared pool)

### Why Use Singletons?

**Valid use cases**:
1. **Configuration**: One global config for the entire application
2. **Logging**: Single logger to avoid mixing output
3. **Caching**: One shared cache to maximize hit rate
4. **Connection pools**: Share limited resources (database connections)
5. **Metrics**: Centralized metrics collection

**Anti-patterns** (DON'T use singletons for):
- User objects (multiple users exist!)
- Business logic (harder to test, tight coupling)
- Anything that naturally has multiple instances

### Singleton in Go

**Pattern structure**:
```go
package mypackage

import "sync"

// Private instance (unexported)
var instance *MySingleton
var once sync.Once

// MySingleton is the type (exported so methods can be called)
type MySingleton struct {
    // fields
}

// GetInstance returns the singleton instance (lazy initialization)
func GetInstance() *MySingleton {
    once.Do(func() {
        instance = &MySingleton{
            // initialization
        }
    })
    return instance
}
```

**Key elements**:
1. **Private variable**: `instance` is unexported (lowercase)
2. **sync.Once**: Ensures one-time initialization
3. **Public accessor**: `GetInstance()` is exported
4. **Lazy loading**: Instance created on first `GetInstance()` call

### Singleton Variations

**1. Simple Singleton (no initialization work)**
```go
var instance = &MySingleton{} // Eager initialization

func GetInstance() *MySingleton {
    return instance
}
```

**When to use**: When initialization is cheap and you always need it

**2. Lazy Singleton (deferred initialization)**
```go
var (
    instance *MySingleton
    once     sync.Once
)

func GetInstance() *MySingleton {
    once.Do(func() {
        instance = &MySingleton{}
    })
    return instance
}
```

**When to use**: When initialization is expensive or rarely needed

**3. Singleton with Configuration**
```go
var (
    instance *MySingleton
    once     sync.Once
)

func GetInstance(config Config) *MySingleton {
    once.Do(func() {
        instance = NewMySingleton(config)
    })
    return instance
}
```

**WARNING**: Only the first call's `config` is used! Subsequent calls ignore config.

**Better approach** (fail fast):
```go
func GetInstance(config Config) *MySingleton {
    once.Do(func() {
        instance = NewMySingleton(config)
    })
    if instance.config != config {
        panic("singleton already initialized with different config")
    }
    return instance
}
```

---

## 5. Common Patterns and Best Practices

### Pattern 1: Configuration Singleton

**Use case**: Load config file once, share across application

```go
package config

import (
    "encoding/json"
    "os"
    "sync"
)

var (
    cfg  *Config
    once sync.Once
)

type Config struct {
    DatabaseURL string
    APIKey      string
    Port        int
}

func Get() *Config {
    once.Do(func() {
        cfg = &Config{}
        file, err := os.ReadFile("config.json")
        if err != nil {
            panic(err) // Or use default config
        }
        if err := json.Unmarshal(file, cfg); err != nil {
            panic(err)
        }
    })
    return cfg
}

// Usage:
// dbURL := config.Get().DatabaseURL
```

**Why this works**:
- Config file read only once (fast startup)
- Thread-safe access (no races)
- Global access (convenient)

### Pattern 2: Database Connection Singleton

**Use case**: Single connection pool shared by all goroutines

```go
package database

import (
    "database/sql"
    "sync"
)

var (
    db   *sql.DB
    once sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        var err error
        db, err = sql.Open("postgres", "connection-string")
        if err != nil {
            panic(err)
        }
        // Configure pool
        db.SetMaxOpenConns(25)
        db.SetMaxIdleConns(5)
    })
    return db
}

// Usage:
// row := database.GetDB().QueryRow("SELECT ...")
```

### Pattern 3: Logger Singleton

**Use case**: Centralized logging with consistent configuration

```go
package logger

import (
    "log"
    "os"
    "sync"
)

var (
    logger *log.Logger
    once   sync.Once
)

func Get() *log.Logger {
    once.Do(func() {
        file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
            panic(err)
        }
        logger = log.New(file, "APP: ", log.LstdFlags)
    })
    return logger
}

// Usage:
// logger.Get().Println("Application started")
```

### Pattern 4: Metrics Collector Singleton

**Use case**: Centralized metrics for observability

```go
package metrics

import (
    "sync"
    "sync/atomic"
)

var (
    collector *Collector
    once      sync.Once
)

type Collector struct {
    requestCount uint64
    errorCount   uint64
}

func Get() *Collector {
    once.Do(func() {
        collector = &Collector{}
    })
    return collector
}

func (c *Collector) IncrementRequests() {
    atomic.AddUint64(&c.requestCount, 1)
}

func (c *Collector) IncrementErrors() {
    atomic.AddUint64(&c.errorCount, 1)
}

// Usage:
// metrics.Get().IncrementRequests()
```

### Pattern 5: Error Handling in Initialization

**Problem**: What if initialization fails?

**Option 1: Panic (fail fast)**
```go
func GetDB() *sql.DB {
    once.Do(func() {
        var err error
        db, err = sql.Open("postgres", "connection-string")
        if err != nil {
            panic(err) // Application can't function without DB
        }
    })
    return db
}
```

**When to use**: Critical resources that must succeed

**Option 2: Return error (caller handles)**
```go
var (
    db      *sql.DB
    once    sync.Once
    initErr error
)

func GetDB() (*sql.DB, error) {
    once.Do(func() {
        db, initErr = sql.Open("postgres", "connection-string")
    })
    return db, initErr
}
```

**When to use**: Optional resources or graceful degradation

**Option 3: Retry mechanism (NOT with sync.Once)**
```go
var (
    db *sql.DB
    mu sync.Mutex
)

func GetDB() (*sql.DB, error) {
    mu.Lock()
    defer mu.Unlock()

    if db != nil {
        return db, nil
    }

    var err error
    for i := 0; i < 3; i++ {
        db, err = sql.Open("postgres", "connection-string")
        if err == nil {
            return db, nil
        }
        time.Sleep(time.Second)
    }
    return nil, err
}
```

**When to use**: Flaky resources that might succeed on retry (NOT suitable for sync.Once)

---

## 6. Advanced Concepts

### sync.Once vs Other Synchronization

**vs. Mutex**
```go
// With mutex (slow)
var (
    config *Config
    mu     sync.Mutex
)

func GetConfig() *Config {
    mu.Lock()         // Every call acquires lock
    defer mu.Unlock()
    if config == nil {
        config = loadConfig()
    }
    return config
}

// With sync.Once (fast)
var (
    config *Config
    once   sync.Once
)

func GetConfig() *Config {
    once.Do(func() {      // Only first call uses lock
        config = loadConfig()
    })
    return config
}
```

**Performance**: `sync.Once` is **10-100x faster** after initialization (lock-free fast path)

**vs. sync.RWMutex**
```go
var (
    config *Config
    mu     sync.RWMutex
)

func GetConfig() *Config {
    mu.RLock()
    if config != nil {
        defer mu.RUnlock()
        return config
    }
    mu.RUnlock()

    mu.Lock()
    defer mu.Unlock()
    if config == nil {
        config = loadConfig()
    }
    return config
}
```

**Complexity**: More complex, easy to get wrong
**Performance**: Better than Mutex, but still slower than sync.Once

**Verdict**: Use `sync.Once` for one-time initialization. It's designed for this!

### Memory Ordering and Visibility

**The subtle problem**:
```go
var (
    config *Config
    done   bool // Regular bool, NOT atomic
)

func GetConfig() *Config {
    if done {
        return config // MIGHT SEE NIL! (memory ordering issue)
    }

    mu.Lock()
    if !done {
        config = &Config{}
        done = true // Write might not be visible to other goroutines!
    }
    mu.Unlock()
    return config
}
```

**Problem**: Without proper memory barriers:
- Goroutine A sets `done = true`
- Goroutine B sees `done == true` but `config == nil` (reordering!)

**sync.Once solves this**:
```go
var (
    config *Config
    once   sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        config = &Config{}
    })
    return config // GUARANTEED to see initialized config
}
```

**How?** `sync.Once` uses atomic operations and memory barriers to ensure visibility.

### Testing Singletons

**Problem**: Singletons make testing hard (global state!)

**Solution 1: Reset function (for tests only)**
```go
var (
    instance *MySingleton
    once     sync.Once
)

func GetInstance() *MySingleton {
    once.Do(func() {
        instance = &MySingleton{}
    })
    return instance
}

// ONLY for tests - NOT thread-safe!
func resetForTest() {
    instance = nil
    once = sync.Once{}
}
```

**Solution 2: Dependency injection**
```go
type App struct {
    DB     *sql.DB
    Logger *log.Logger
}

// In production:
app := &App{
    DB:     database.GetDB(),
    Logger: logger.Get(),
}

// In tests:
app := &App{
    DB:     testDB,
    Logger: testLogger,
}
```

**Solution 3: Interface-based design**
```go
type Database interface {
    Query(string) *sql.Rows
}

type App struct {
    DB Database // Interface, not concrete singleton
}
```

**Best practice**: Use dependency injection for testability, singletons for convenience.

---

## 7. Real-World Applications

### Application Configuration

**Scenario**: Web application with config file

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "sync"
)

var (
    config *AppConfig
    once   sync.Once
)

type AppConfig struct {
    Port        int
    DatabaseURL string
    APIKeys     map[string]string
}

func GetConfig() *AppConfig {
    once.Do(func() {
        config = &AppConfig{}
        data, err := os.ReadFile("config.json")
        if err != nil {
            panic(fmt.Sprintf("failed to read config: %v", err))
        }
        if err := json.Unmarshal(data, config); err != nil {
            panic(fmt.Sprintf("failed to parse config: %v", err))
        }
    })
    return config
}

func handler(w http.ResponseWriter, r *http.Request) {
    cfg := GetConfig() // Fast after first call
    fmt.Fprintf(w, "API Key: %s", cfg.APIKeys["default"])
}

func main() {
    http.HandleFunc("/", handler)
    cfg := GetConfig()
    http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
```

**Used by**: Almost every production Go application

### Connection Pool Management

**Scenario**: Redis client singleton

```go
package cache

import (
    "github.com/go-redis/redis/v8"
    "sync"
)

var (
    client *redis.Client
    once   sync.Once
)

func GetRedis() *redis.Client {
    once.Do(func() {
        client = redis.NewClient(&redis.Options{
            Addr:         "localhost:6379",
            PoolSize:     10,
            MinIdleConns: 5,
        })
    })
    return client
}

// Usage across application:
// cache.GetRedis().Set(ctx, key, value, expiration)
```

**Used by**: Docker, Kubernetes, most microservices

### Rate Limiter Initialization

**Scenario**: Global rate limiter for API

```go
package ratelimit

import (
    "golang.org/x/time/rate"
    "sync"
)

var (
    limiter *rate.Limiter
    once    sync.Once
)

func GetLimiter() *rate.Limiter {
    once.Do(func() {
        // 100 requests per second, burst of 200
        limiter = rate.NewLimiter(100, 200)
    })
    return limiter
}

// Usage in HTTP middleware:
// if !ratelimit.GetLimiter().Allow() {
//     http.Error(w, "Too Many Requests", 429)
// }
```

### Feature Flag System

**Scenario**: Load feature flags once from remote service

```go
package features

import (
    "net/http"
    "sync"
)

var (
    flags map[string]bool
    once  sync.Once
)

func IsEnabled(feature string) bool {
    once.Do(func() {
        flags = fetchFlagsFromAPI()
    })
    return flags[feature]
}

func fetchFlagsFromAPI() map[string]bool {
    // HTTP call to feature flag service
    // Returns map of feature names to enabled status
    return map[string]bool{
        "new_ui":       true,
        "beta_feature": false,
    }
}

// Usage:
// if features.IsEnabled("new_ui") {
//     renderNewUI()
// }
```

---

## 8. Common Mistakes and Gotchas

### Mistake 1: Passing Different Arguments

**❌ Wrong**:
```go
var (
    config *Config
    once   sync.Once
)

func GetConfig(filename string) *Config {
    once.Do(func() {
        config = loadConfig(filename)
    })
    return config
}

// First call:
cfg1 := GetConfig("config.json") // Loads config.json

// Second call:
cfg2 := GetConfig("other.json") // STILL RETURNS config.json!
```

**Problem**: `once.Do()` only runs the first call's function. Subsequent arguments are ignored.

**✅ Correct**:
```go
func GetConfig(filename string) *Config {
    once.Do(func() {
        config = loadConfig(filename)
    })
    if config.filename != filename {
        panic(fmt.Sprintf("config already loaded from %s, cannot load %s",
            config.filename, filename))
    }
    return config
}
```

### Mistake 2: Not Handling Initialization Errors

**❌ Wrong**:
```go
var (
    db   *sql.DB
    once sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        db, _ = sql.Open("postgres", "connection-string")
        // Ignoring error! db might be nil!
    })
    return db // Could return nil!
}
```

**✅ Correct (Option 1 - Panic)**:
```go
func GetDB() *sql.DB {
    once.Do(func() {
        var err error
        db, err = sql.Open("postgres", "connection-string")
        if err != nil {
            panic(fmt.Sprintf("failed to connect to DB: %v", err))
        }
    })
    return db
}
```

**✅ Correct (Option 2 - Return error)**:
```go
var (
    db      *sql.DB
    once    sync.Once
    initErr error
)

func GetDB() (*sql.DB, error) {
    once.Do(func() {
        db, initErr = sql.Open("postgres", "connection-string")
    })
    return db, initErr
}
```

### Mistake 3: Trying to Retry with sync.Once

**❌ Wrong**:
```go
var (
    db   *sql.DB
    once sync.Once
)

func GetDB() (*sql.DB, error) {
    once.Do(func() {
        db, err := sql.Open("postgres", "connection-string")
        if err != nil {
            // Want to retry on next call - but can't!
            return
        }
    })
    if db == nil {
        return nil, errors.New("db not initialized")
    }
    return db, nil
}
```

**Problem**: `once.Do()` will never run again, even if first attempt failed.

**✅ Correct (use Mutex for retry logic)**:
```go
var (
    db *sql.DB
    mu sync.Mutex
)

func GetDB() (*sql.DB, error) {
    mu.Lock()
    defer mu.Unlock()

    if db != nil {
        return db, nil
    }

    var err error
    db, err = sql.Open("postgres", "connection-string")
    if err != nil {
        db = nil // Reset so we retry next time
        return nil, err
    }
    return db, nil
}
```

### Mistake 4: Race Condition in Lazy Initialization

**❌ Wrong**:
```go
var config *Config

func GetConfig() *Config {
    if config == nil { // RACE! Check and set are not atomic
        config = &Config{}
    }
    return config
}
```

**Problem**: Multiple goroutines might all see `config == nil` and initialize multiple times.

**✅ Correct**:
```go
var (
    config *Config
    once   sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        config = &Config{}
    })
    return config
}
```

### Mistake 5: Forgetting Package-Level Variables

**❌ Wrong**:
```go
func GetConfig() *Config {
    var (
        config *Config
        once   sync.Once // NEW Once EVERY CALL!
    )

    once.Do(func() {
        config = &Config{} // Runs every time!
    })
    return config
}
```

**Problem**: `once` is local to the function, so it's recreated on each call.

**✅ Correct**:
```go
var (
    config *Config
    once   sync.Once // Package-level, persists across calls
)

func GetConfig() *Config {
    once.Do(func() {
        config = &Config{}
    })
    return config
}
```

---

## 9. Performance Characteristics

### Benchmark: sync.Once vs Mutex vs RWMutex

**Setup**:
```go
func BenchmarkOnce(b *testing.B) {
    var once sync.Once
    var value int

    once.Do(func() { value = 1 })

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            once.Do(func() { value++ })
        }
    })
}

func BenchmarkMutex(b *testing.B) {
    var mu sync.Mutex
    var value int
    var initialized bool

    mu.Lock()
    value = 1
    initialized = true
    mu.Unlock()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mu.Lock()
            if !initialized {
                value++
            }
            mu.Unlock()
        }
    })
}
```

**Results** (approximate):
```
BenchmarkOnce-8      1000000000    0.50 ns/op   (fast path: atomic load)
BenchmarkMutex-8       100000000   15.0  ns/op   (lock on every call)
BenchmarkRWMutex-8     200000000    7.0  ns/op   (read lock on every call)
```

**Takeaway**: `sync.Once` is **30x faster** than Mutex, **14x faster** than RWMutex for the common case (after initialization).

### Memory Overhead

**sync.Once**:
```go
type Once struct {
    done uint32      // 4 bytes
    m    sync.Mutex  // ~8 bytes (1 int32 + 1 int32)
}
// Total: ~12 bytes
```

**Negligible** - you can have thousands of `sync.Once` instances with minimal memory impact.

---

## 10. Stretch Goals

### Goal 1: Configuration Hot Reload ⭐⭐

Implement a singleton that detects config file changes and reloads.

**Hint**: Use `fsnotify` to watch file changes. Can't use `sync.Once` for reload (use `sync.RWMutex` instead).

### Goal 2: Lazy Field Initialization ⭐

Implement a struct where each field is initialized lazily with `sync.Once`.

**Hint**:
```go
type LazyConfig struct {
    dbOnce sync.Once
    db     *sql.DB

    loggerOnce sync.Once
    logger     *log.Logger
}

func (c *LazyConfig) GetDB() *sql.DB {
    c.dbOnce.Do(func() {
        c.db, _ = sql.Open("postgres", "...")
    })
    return c.db
}
```

### Goal 3: Singleton with Context ⭐⭐

Implement a singleton that respects context cancellation during initialization.

**Hint**: Can't use `sync.Once` directly (no context support). Use atomic CAS instead.

### Goal 4: Registry Pattern ⭐⭐⭐

Implement a singleton registry that manages multiple named singletons.

**Hint**:
```go
type Registry struct {
    mu        sync.RWMutex
    instances map[string]interface{}
}

func (r *Registry) GetOrCreate(name string, factory func() interface{}) interface{} {
    // Check if exists (read lock)
    // If not, create with factory (write lock + sync.Once per name)
}
```

### Goal 5: Thread-Safe Lazy Map ⭐⭐⭐

Implement a map where each value is lazily initialized using `sync.Once`.

**Hint**: `sync.Map` + per-key `sync.Once`

---

## How to Run

```bash
# Run the demo
go run ./minis/26-sync-once-singleton/cmd/once-demo/main.go

# Run tests
go test ./minis/26-sync-once-singleton/exercise/...

# Run with verbose output
go test -v ./minis/26-sync-once-singleton/exercise/...

# Run with race detector (important for concurrency!)
go test -race ./minis/26-sync-once-singleton/exercise/...

# Run benchmarks
go test -bench=. ./minis/26-sync-once-singleton/exercise/...
```

---

## Summary

**What you learned**:
- ✅ `sync.Once` guarantees exactly-once execution in concurrent code
- ✅ Singleton pattern ensures single instances with lazy initialization
- ✅ Lazy initialization improves startup time and resource usage
- ✅ Idempotent initialization prevents duplicate work and race conditions
- ✅ Proper error handling in initialization is critical
- ✅ `sync.Once` is faster than Mutex/RWMutex for one-time initialization

**Why this matters**:
Almost every production Go application uses `sync.Once`:
- **Configuration loading**: Kubernetes, Docker, Prometheus
- **Database connections**: All web frameworks (Gin, Echo, Fiber)
- **Logging**: Uber's Zap, Logrus
- **Metrics**: Prometheus client library
- **Feature flags**: LaunchDarkly, Split.io

**Key rules**:
1. Use `sync.Once` for one-time initialization (don't roll your own)
2. Package-level variables (not function-local)
3. Handle initialization errors explicitly (panic or return error)
4. Don't try to retry with `sync.Once` (it never runs twice)
5. Test with `-race` flag to catch concurrency bugs

**Next steps**:
- Project 27: sync.Map (concurrent map with lazy initialization)
- Project 28: Worker pools with initialization
- Project 29: Resource lifecycle management

Master `sync.Once`, master Go's concurrency primitives! 🚀
