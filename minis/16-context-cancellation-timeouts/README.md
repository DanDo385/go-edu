# Project 16: Context Cancellation and Timeouts

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a web server that makes several database queries and API calls to handle each request:

**❌ Without context management:**
```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    user := db.GetUser(userID)           // Takes 100ms
    posts := db.GetPosts(userID)         // Takes 200ms
    comments := api.GetComments(postIDs) // Takes 300ms
    // Total: 600ms even if client disconnected after 50ms!
}
```

**Problem**: If the client disconnects (closes browser, network timeout, user navigates away), your server keeps working on a response nobody wants. This wastes resources (database connections, API quotas, CPU time).

**✅ With context management:**
```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // Get request context

    user := db.GetUserContext(ctx, userID)
    if ctx.Err() != nil {
        return // Client disconnected, stop immediately
    }

    posts := db.GetPostsContext(ctx, userID)
    if ctx.Err() != nil {
        return // Client disconnected, stop immediately
    }

    comments := api.GetCommentsContext(ctx, postIDs)
    // All operations respect context cancellation
}
```

**Result**: When the client disconnects, all operations stop immediately. No wasted work!

This project teaches you **Go's context package**, which is the standard way to:
- **Propagate cancellation** across function boundaries
- **Set timeouts and deadlines** for operations
- **Prevent goroutine leaks** when work is no longer needed
- **Pass request-scoped values** (use sparingly!)

### What You'll Learn

1. **Context fundamentals**: What context is and why it exists
2. **Context cancellation**: How to stop work when it's no longer needed
3. **Timeouts and deadlines**: Automatic cancellation after time limits
4. **Context propagation**: Passing context through call stacks
5. **Goroutine lifecycle management**: Preventing leaks
6. **Real-world HTTP patterns**: Timeouts in web services

### The Challenge

Build functions that demonstrate:
- Manual cancellation with `context.WithCancel`
- Timeout-based cancellation with `context.WithTimeout`
- Deadline-based cancellation with `context.WithDeadline`
- Preventing goroutine leaks when operations are cancelled
- Combining timeouts with HTTP requests

---

## 2. First Principles: Understanding Context

### What is Context?

`context.Context` is an **interface** that carries:
1. **Cancellation signals**: "Stop what you're doing"
2. **Deadlines**: "Stop at this specific time"
3. **Values**: Request-scoped data (authentication tokens, request IDs, etc.)

**The interface** (you don't need to implement this yourself):
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

**Key methods explained**:

1. **`Done() <-chan struct{}`**: Returns a channel that's closed when context is cancelled
   - You don't read from this channel
   - You use it in `select` statements to detect cancellation

2. **`Err() error`**: Returns why context was cancelled
   - `nil` if not cancelled
   - `context.Canceled` if manually cancelled
   - `context.DeadlineExceeded` if timeout/deadline passed

3. **`Deadline() (time.Time, bool)`**: Returns when context will be cancelled
   - `ok` is `false` if no deadline

4. **`Value(key) any`**: Gets request-scoped values (avoid using except for specific cases)

### Why Does Context Exist?

**Problem 1: No standard way to cancel work**

Before Go 1.7 (when context was added to stdlib):
```go
func DoWork(stopChan <-chan bool) {
    // Every function invents its own cancellation mechanism
    select {
    case <-stopChan:
        return
    }
}

func DoOtherWork(cancel *atomic.Bool) {
    // Different cancellation mechanism!
    if cancel.Load() {
        return
    }
}
```

Each library used different patterns. No interoperability.

**After context**:
```go
func DoWork(ctx context.Context) {
    // Standard cancellation mechanism
    select {
    case <-ctx.Done():
        return
    }
}
```

Every library uses `context.Context`. Cancellation propagates seamlessly.

**Problem 2: Resource leaks from orphaned goroutines**

```go
// BAD: Goroutine leak
func SearchWithoutContext(query string) string {
    resultCh := make(chan string)

    go func() {
        result := slowSearch(query) // Takes 30 seconds
        resultCh <- result
    }()

    select {
    case result := <-resultCh:
        return result
    case <-time.After(1 * time.Second):
        return "timeout" // Goroutine still running!
    }
}
```

The goroutine continues running for 29 more seconds even though we've returned. Over time, thousands of leaked goroutines accumulate → memory leak.

**With context**:
```go
// GOOD: No leak
func SearchWithContext(ctx context.Context, query string) string {
    resultCh := make(chan string)

    go func() {
        result := slowSearchContext(ctx, query) // Stops when ctx is cancelled
        resultCh <- result
    }()

    select {
    case result := <-resultCh:
        return result
    case <-ctx.Done():
        return "cancelled" // Goroutine stops too!
    }
}
```

### Context Tree Structure

Contexts form a **tree**: child contexts are derived from parent contexts.

```
context.Background() (root)
    │
    ├─ WithTimeout(5s)
    │   │
    │   ├─ WithCancel()
    │   │   │
    │   │   └─ WithValue("userID", 123)
    │   │
    │   └─ WithTimeout(1s)  ← Inherits parent's 5s timeout
    │
    └─ WithCancel()
```

**Key rule**: When a parent context is cancelled, **all children are cancelled too**.

**Example**:
```go
parent, cancel := context.WithTimeout(context.Background(), 5*time.Second)
child1, _ := context.WithTimeout(parent, 10*time.Second)
child2, _ := context.WithTimeout(parent, 10*time.Second)

cancel() // Cancels parent, child1, and child2 immediately
// Even though child1 and child2 had 10s timeouts, parent cancellation overrides
```

### When to Use Context

**✅ DO use context for**:
- HTTP request handling (propagate client disconnection)
- Database queries (stop query if request is cancelled)
- RPC calls (timeout for slow remote services)
- Goroutine coordination (tell workers to stop)
- Long-running operations (respect user cancellation)

**❌ DON'T use context for**:
- Passing dependencies (use struct fields or function parameters)
- Optional parameters (use function options pattern)
- Storing large objects (context values are for small metadata)
- Application-wide settings (use configuration structs)

**Context Convention**: Context is **always the first parameter** of a function:
```go
func DoWork(ctx context.Context, arg1 string, arg2 int) error {
    // ...
}
```

---

## 3. Context Creation Functions

### context.Background()

**Purpose**: The root of all contexts. Never cancelled, no deadline, no values.

**When to use**:
- In `main()` function
- In test setup (unless testing timeouts)
- When you don't have a parent context

**Example**:
```go
func main() {
    ctx := context.Background()
    results, err := FetchData(ctx)
}
```

### context.TODO()

**Purpose**: Placeholder when you're not sure which context to use.

**When to use**:
- During development/refactoring
- When you plan to add proper context later
- To make code compile while you figure out the right context

**Example**:
```go
func LegacyFunction() {
    // TODO: Get context from caller
    ctx := context.TODO()
    NewFunction(ctx)
}
```

**In production**: `context.Background()` and `context.TODO()` are identical. `TODO()` is just documentation.

### context.WithCancel(parent)

**Purpose**: Create a cancellable context.

**Signature**:
```go
func WithCancel(parent context.Context) (ctx context.Context, cancel context.CancelFunc)
```

**Returns**:
- `ctx`: New context that is cancelled when `cancel()` is called
- `cancel`: Function to call when you want to cancel

**When to use**:
- You need manual control over cancellation
- Multiple goroutines that should all stop together
- Coordinating worker shutdown

**Example**:
```go
ctx, cancel := context.WithCancel(context.Background())

go worker(ctx) // Worker will stop when context is cancelled
go worker(ctx)
go worker(ctx)

time.Sleep(5 * time.Second)
cancel() // Stop all workers

// IMPORTANT: Always defer cancel() to prevent leaks
ctx, cancel := context.WithCancel(parent)
defer cancel() // Releases resources even if not called explicitly
```

**Critical rule**: **Always call `cancel()` eventually**, even if not explicitly cancelling.

Why? Context implementations allocate resources (goroutines for timers, channels, etc.). Calling `cancel()` cleans them up.

**Pattern**:
```go
ctx, cancel := context.WithCancel(parent)
defer cancel() // ← ALWAYS do this
// ... use ctx ...
```

### context.WithTimeout(parent, duration)

**Purpose**: Automatically cancel after a duration.

**Signature**:
```go
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
```

**Example**:
```go
// Automatically cancel after 3 seconds
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

result, err := FetchData(ctx)
if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("Operation timed out")
}
```

**When to use**:
- HTTP requests (don't wait forever for slow servers)
- Database queries (prevent slow queries from blocking)
- RPC calls (timeout on slow services)

**What happens**:
1. Timer starts when context is created
2. After `timeout` elapses, `ctx.Done()` channel is closed
3. `ctx.Err()` returns `context.DeadlineExceeded`
4. All child contexts are also cancelled

### context.WithDeadline(parent, time)

**Purpose**: Cancel at a specific point in time.

**Signature**:
```go
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc)
```

**Example**:
```go
// Cancel at exactly 3:00 PM
deadline := time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

**When to use**:
- Rate limiting (requests must complete before rate limit window ends)
- Scheduled tasks (work must finish before next schedule)
- SLA enforcement (response must be ready by specific time)

**Relationship to WithTimeout**:
```go
// These are equivalent:
ctx1, _ := context.WithTimeout(parent, 5*time.Second)
ctx2, _ := context.WithDeadline(parent, time.Now().Add(5*time.Second))
```

`WithTimeout` is just a convenience wrapper around `WithDeadline`.

### context.WithValue(parent, key, value)

**Purpose**: Store request-scoped values in context.

**Signature**:
```go
func WithValue(parent context.Context, key, val any) context.Context
```

**When to use (rarely!)**:
- Request IDs for tracing
- Authentication tokens
- User identity
- Request-scoped loggers

**When NOT to use**:
- Configuration (use structs)
- Optional parameters (use function options)
- Large objects (pass explicitly)

**Example**:
```go
type contextKey string

const requestIDKey contextKey = "requestID"

// Set value
ctx := context.WithValue(context.Background(), requestIDKey, "req-12345")

// Get value
requestID := ctx.Value(requestIDKey).(string)
```

**Best practices**:
1. **Use custom types for keys** (not strings) to avoid collisions:
   ```go
   type contextKey string // unexported type
   const myKey contextKey = "mykey"
   ```

2. **Provide helper functions**:
   ```go
   func WithRequestID(ctx context.Context, id string) context.Context {
       return context.WithValue(ctx, requestIDKey, id)
   }

   func GetRequestID(ctx context.Context) string {
       id, ok := ctx.Value(requestIDKey).(string)
       if !ok {
           return ""
       }
       return id
   }
   ```

3. **Keep values small** (pointers, IDs, tokens - not large structs)

---

## 4. Using Context Effectively

### Pattern 1: Checking for Cancellation

**In loops**:
```go
func ProcessItems(ctx context.Context, items []Item) error {
    for _, item := range items {
        // Check cancellation before expensive operation
        select {
        case <-ctx.Done():
            return ctx.Err() // Return context error
        default:
            // Context not cancelled, continue
        }

        process(item)
    }
    return nil
}
```

**In long-running operations**:
```go
func LongTask(ctx context.Context) error {
    for {
        // Do some work
        work()

        // Periodically check cancellation
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue
        }
    }
}
```

**In select with multiple channels**:
```go
func Worker(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case <-ctx.Done():
            return // Context cancelled
        case job := <-jobs:
            process(job)
        }
    }
}
```

### Pattern 2: Passing Context to Functions

**HTTP Handlers**:
```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // Get request context

    // Pass to database
    user, err := db.GetUserContext(ctx, userID)

    // Pass to API client
    posts, err := api.GetPostsContext(ctx, userID)
}
```

**Database Queries**:
```go
func GetUser(ctx context.Context, db *sql.DB, id int) (*User, error) {
    var user User

    query := "SELECT name, email FROM users WHERE id = ?"
    err := db.QueryRowContext(ctx, query, id).Scan(&user.Name, &user.Email)

    return &user, err
}
```

**HTTP Requests**:
```go
func FetchURL(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Pattern 3: Preventing Goroutine Leaks

**❌ LEAK: Goroutine never stops**:
```go
func Search(query string) string {
    ch := make(chan string)

    go func() {
        result := slowSearch(query) // Never checks cancellation
        ch <- result
    }()

    select {
    case result := <-ch:
        return result
    case <-time.After(1 * time.Second):
        return "timeout" // Goroutine still running! LEAK!
    }
}
```

**✅ NO LEAK: Goroutine stops with context**:
```go
func Search(ctx context.Context, query string) string {
    ch := make(chan string, 1) // Buffered to prevent sender blocking

    go func() {
        result := slowSearchContext(ctx, query) // Respects context
        select {
        case ch <- result:
        case <-ctx.Done():
            // Context cancelled, don't send (avoid blocking)
        }
    }()

    select {
    case result := <-ch:
        return result
    case <-ctx.Done():
        return "cancelled" // Goroutine will also stop
    }
}
```

**Key techniques**:
1. Pass context to goroutine
2. Check `ctx.Done()` inside goroutine
3. Use buffered channel (size 1) if sender might not have receiver
4. Use `select` when sending to channel in goroutine

### Pattern 4: Timeout for HTTP Requests

**Per-request timeout**:
```go
func FetchWithTimeout(url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

**Combined timeout and cancellation**:
```go
func FetchWithCancellation(parentCtx context.Context, url string) ([]byte, error) {
    // Create child context with timeout
    // Will cancel when:
    // 1. Parent is cancelled, OR
    // 2. 5 seconds elapse
    ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        // Check what kind of error
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, fmt.Errorf("request timed out: %w", err)
        }
        if errors.Is(err, context.Canceled) {
            return nil, fmt.Errorf("request cancelled: %w", err)
        }
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Pattern 5: Coordinating Multiple Goroutines

**Stop all workers when done**:
```go
func RunWorkers(ctx context.Context, numWorkers int) {
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(ctx, id)
        }(i)
    }

    // Wait for all workers to finish
    wg.Wait()
}

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d stopping\n", id)
            return
        default:
            // Do work
            time.Sleep(100 * time.Millisecond)
        }
    }
}
```

**Cancel all on first error**:
```go
func FetchAll(urls []string) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    errCh := make(chan error, len(urls))

    for _, url := range urls {
        go func(u string) {
            _, err := FetchWithCancellation(ctx, u)
            if err != nil {
                cancel() // Cancel all other requests
                errCh <- err
            }
        }(url)
    }

    select {
    case err := <-errCh:
        return err
    case <-time.After(10 * time.Second):
        return fmt.Errorf("timeout")
    }
}
```

---

## 5. Common Patterns You Can Reuse

### Pattern 1: Timeout Wrapper

```go
// WithTimeout wraps any function with a timeout
func WithTimeout[T any](fn func(context.Context) (T, error), timeout time.Duration) (T, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    return fn(ctx)
}

// Usage:
result, err := WithTimeout(func(ctx context.Context) (string, error) {
    return slowOperation(ctx)
}, 5*time.Second)
```

### Pattern 2: Graceful Shutdown

```go
func Server() {
    srv := &http.Server{Addr: ":8080"}

    // Start server
    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}
```

### Pattern 3: Context-Aware Sleep

```go
// Sleep that can be interrupted by context
func Sleep(ctx context.Context, duration time.Duration) error {
    select {
    case <-time.After(duration):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Usage:
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(1 * time.Second)
    cancel()
}()

if err := Sleep(ctx, 10*time.Second); err != nil {
    fmt.Println("Sleep interrupted:", err)
}
```

### Pattern 4: Request ID Propagation

```go
type contextKey string

const requestIDKey contextKey = "requestID"

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
    id, ok := ctx.Value(requestIDKey).(string)
    if !ok {
        return ""
    }
    return id
}

// Usage:
func handler(w http.ResponseWriter, r *http.Request) {
    requestID := uuid.New().String()
    ctx := WithRequestID(r.Context(), requestID)

    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    log.Printf("[%s] Processing request", GetRequestID(ctx))
}
```

### Pattern 5: Fan-Out with Context

```go
func FanOut[T any](ctx context.Context, inputs []T, fn func(context.Context, T) error) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(inputs))

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    for _, input := range inputs {
        wg.Add(1)
        go func(in T) {
            defer wg.Done()
            if err := fn(ctx, in); err != nil {
                errCh <- err
                cancel() // Cancel all others on first error
            }
        }(input)
    }

    wg.Wait()
    close(errCh)

    // Return first error (if any)
    for err := range errCh {
        return err
    }
    return nil
}
```

---

## 6. Real-World Applications

### Web Servers

**Use case**: Handle client disconnection gracefully

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // If client disconnects, context is cancelled
    user, err := db.GetUser(ctx, userID)
    if err != nil {
        if errors.Is(err, context.Canceled) {
            // Client disconnected, don't bother responding
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}
```

Companies: Google, Facebook, Netflix (all use context for request handling)

### Microservices

**Use case**: Timeout for inter-service calls

```go
type UserService struct {
    client *http.Client
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Each service call has 2s timeout
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    url := fmt.Sprintf("http://user-service/users/%d", id)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    resp, err := s.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    return &user, nil
}
```

Companies: Uber, Airbnb (timeout all inter-service calls)

### Database Operations

**Use case**: Prevent slow queries from blocking

```go
func GetOrders(ctx context.Context, db *sql.DB, userID int) ([]Order, error) {
    // Query must complete in 5 seconds
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        SELECT id, user_id, total, created_at
        FROM orders
        WHERE user_id = ?
        ORDER BY created_at DESC
        LIMIT 100
    `

    rows, err := db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []Order
    for rows.Next() {
        var o Order
        if err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.CreatedAt); err != nil {
            return nil, err
        }
        orders = append(orders, o)
    }

    return orders, rows.Err()
}
```

### Background Workers

**Use case**: Gracefully stop workers on shutdown

```go
func Worker(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case <-ctx.Done():
            log.Println("Worker shutting down...")
            return
        case job := <-jobs:
            process(job)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    jobs := make(chan Job)

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            Worker(ctx, jobs)
        }()
    }

    // Wait for interrupt
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt)
    <-sigCh

    // Stop all workers
    cancel()
    wg.Wait()

    log.Println("All workers stopped")
}
```

### ETL Pipelines

**Use case**: Cancel entire pipeline on failure

```go
func RunPipeline(ctx context.Context, input <-chan Record) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Stage 1: Extract
    extracted := extract(ctx, input)

    // Stage 2: Transform
    transformed := transform(ctx, extracted)

    // Stage 3: Load
    return load(ctx, transformed)
}

func extract(ctx context.Context, input <-chan Record) <-chan Record {
    output := make(chan Record)
    go func() {
        defer close(output)
        for record := range input {
            select {
            case <-ctx.Done():
                return
            case output <- record:
            }
        }
    }()
    return output
}
```

---

## 7. Common Mistakes to Avoid

### Mistake 1: Not Calling cancel()

**❌ Wrong**:
```go
func badExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    // Forgot to call cancel()!

    doWork(ctx)
}
```

**Problem**: Timer goroutine leaks. Even if operation finishes early, timer keeps running.

**✅ Correct**:
```go
func goodExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // ALWAYS defer cancel()

    doWork(ctx)
}
```

### Mistake 2: Ignoring Context in Long Operations

**❌ Wrong**:
```go
func processItems(ctx context.Context, items []Item) {
    for _, item := range items {
        // Never checks context!
        process(item) // Each takes 1 second
    }
}
```

**Problem**: If context is cancelled after first item, function runs for full duration anyway.

**✅ Correct**:
```go
func processItems(ctx context.Context, items []Item) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        process(item)
    }
    return nil
}
```

### Mistake 3: Passing nil Context

**❌ Wrong**:
```go
func doWork(ctx context.Context) {
    // ...
}

func caller() {
    doWork(nil) // PANIC when ctx.Done() is called!
}
```

**Problem**: `nil` context will panic when any method is called on it.

**✅ Correct**:
```go
func caller() {
    doWork(context.Background()) // Always use a real context
}
```

### Mistake 4: Storing Context in Struct

**❌ Wrong**:
```go
type Server struct {
    ctx context.Context // DON'T DO THIS
}
```

**Problem**: Context is request-scoped, not object-scoped. Storing in struct means all requests share same context.

**✅ Correct**:
```go
type Server struct {
    // No context field
}

func (s *Server) HandleRequest(ctx context.Context) {
    // Pass context as parameter
}
```

### Mistake 5: Using Context for Function Parameters

**❌ Wrong**:
```go
func getUser(ctx context.Context, db *sql.DB) (*User, error) {
    // Getting db from context
    db := ctx.Value("db").(*sql.DB)

    // Getting userID from context
    userID := ctx.Value("userID").(int)

    // ...
}
```

**Problem**: Makes function signature unclear. Hard to test. Fragile (type assertions can panic).

**✅ Correct**:
```go
func getUser(ctx context.Context, db *sql.DB, userID int) (*User, error) {
    // Explicit parameters are clear and type-safe
    // ...
}
```

### Mistake 6: Creating Context Inside Goroutine

**❌ Wrong**:
```go
func spawn(parentCtx context.Context) {
    go func() {
        ctx, cancel := context.WithCancel(parentCtx)
        defer cancel()

        doWork(ctx)
    }()
    // Goroutine might never run, cancel() might never be called
}
```

**Problem**: If goroutine doesn't run immediately, `cancel()` isn't called, resources leak.

**✅ Correct**:
```go
func spawn(parentCtx context.Context) {
    ctx, cancel := context.WithCancel(parentCtx)

    go func() {
        defer cancel()
        doWork(ctx)
    }()
}
```

### Mistake 7: Forgetting Buffered Channel for Goroutine Results

**❌ LEAK**:
```go
func search(ctx context.Context) (string, error) {
    ch := make(chan string) // Unbuffered

    go func() {
        result := slowOp()
        ch <- result // Blocks forever if nobody is receiving!
    }()

    select {
    case result := <-ch:
        return result, nil
    case <-ctx.Done():
        return "", ctx.Err() // Goroutine still blocked on send!
    }
}
```

**✅ NO LEAK**:
```go
func search(ctx context.Context) (string, error) {
    ch := make(chan string, 1) // Buffered (size 1)

    go func() {
        result := slowOp()
        select {
        case ch <- result: // Won't block even if nobody receives
        case <-ctx.Done():
        }
    }()

    select {
    case result := <-ch:
        return result, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

---

## 8. Stretch Goals

### Goal 1: Implement Retry with Exponential Backoff ⭐

Add context-aware retry logic with exponential backoff.

**Hint**:
```go
func RetryWithContext(ctx context.Context, fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        if err := fn(); err == nil {
            return nil
        }

        delay := time.Duration(1<<uint(attempt)) * 100 * time.Millisecond

        select {
        case <-time.After(delay):
            // Continue to next attempt
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return fmt.Errorf("max retries exceeded")
}
```

### Goal 2: Add Request Tracing ⭐⭐

Implement distributed tracing with context values.

**Hint**:
```go
type Trace struct {
    TraceID string
    SpanID  string
}

func WithTrace(ctx context.Context, trace Trace) context.Context {
    return context.WithValue(ctx, traceKey, trace)
}

func GetTrace(ctx context.Context) Trace {
    trace, ok := ctx.Value(traceKey).(Trace)
    if !ok {
        return Trace{}
    }
    return trace
}
```

### Goal 3: Context-Aware Rate Limiter ⭐⭐⭐

Build a rate limiter that respects context cancellation.

**Hint**:
```go
type RateLimiter struct {
    tokens chan struct{}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Goal 4: Deadline Propagation Middleware ⭐⭐

HTTP middleware that sets deadline based on remaining time.

**Hint**:
```go
func DeadlineMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        if deadline, ok := ctx.Deadline(); ok {
            remaining := time.Until(deadline)
            ctx, cancel := context.WithTimeout(ctx, remaining*9/10) // Leave 10% buffer
            defer cancel()
            r = r.WithContext(ctx)
        }

        next.ServeHTTP(w, r)
    })
}
```

### Goal 5: Parallel Task Executor ⭐⭐⭐

Execute tasks in parallel with context, return on first error.

**Hint**:
```go
func Parallel(ctx context.Context, tasks ...func(context.Context) error) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    errCh := make(chan error, len(tasks))

    for _, task := range tasks {
        go func(t func(context.Context) error) {
            if err := t(ctx); err != nil {
                errCh <- err
                cancel()
            }
        }(task)
    }

    select {
    case err := <-errCh:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

## How to Run

```bash
# Run the demo
go run ./minis/16-context-cancellation-timeouts/cmd/context-demo/main.go

# Run tests
go test ./minis/16-context-cancellation-timeouts/...

# Run with verbose output
go test -v ./minis/16-context-cancellation-timeouts/...

# Run with race detector
go test -race ./minis/16-context-cancellation-timeouts/...
```

---

## Summary

**What you learned**:
- ✅ Context is Go's standard way to handle cancellation and timeouts
- ✅ Always pass context as first parameter
- ✅ Always `defer cancel()` to prevent resource leaks
- ✅ Check `ctx.Done()` in long-running operations
- ✅ Use context values sparingly (only for request-scoped metadata)
- ✅ Context propagates through call chains (parent cancels all children)

**Why this matters**:
Context is fundamental to writing production Go code. Every major Go library supports context:
- `net/http`: Request handling
- `database/sql`: Query cancellation
- `google.golang.org/grpc`: RPC timeouts
- Cloud SDKs (AWS, GCP, Azure): API call timeouts

**Key rules**:
1. Context is always the first parameter: `func Do(ctx context.Context, ...)`
2. Always `defer cancel()` after creating context
3. Never store context in a struct
4. Never pass `nil` context (use `context.Background()` or `context.TODO()`)
5. Check `ctx.Done()` in loops and long operations

**Next steps**:
- Project 17: File streaming with bufio (context-aware I/O)
- Project 18: Million goroutines demo (coordinating with context)
- Project 19: Channels basics (combining channels and context)

Master context, master Go! 🚀
