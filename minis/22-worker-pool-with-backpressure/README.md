# Project 22: Worker Pool with Backpressure

## What Is This Project About?

This project teaches you **backpressure** - one of the most critical concepts for building production-ready concurrent systems. You'll learn:

1. **Bounded channels** (preventing unbounded queue growth)
2. **Backpressure strategies** (block, drop, reject, timeout)
3. **Overflow handling** (what to do when queues fill up)
4. **Rate limiting** (controlling throughput)
5. **Graceful degradation** (maintaining system stability under load)

By the end, you'll understand how to build concurrent systems that remain stable and performant even under extreme load.

---

## The Fundamental Problem: Unbounded Queues Kill Systems

### First Principles: What Is Backpressure?

Imagine a restaurant:
- **Kitchen** (workers) can prepare 10 meals per minute
- **Orders** (requests) arrive at 20 per minute
- **Order queue** grows infinitely

**What happens?**
1. Queue grows: 10, 20, 30, 40, 50... orders waiting
2. Wait times increase: 1 min, 2 min, 3 min, 4 min...
3. Customers leave angry
4. Kitchen eventually crashes (out of space)

**Backpressure** is the restaurant saying **"We're full, please wait"** or **"Sorry, we can't take more orders right now"**.

### The Critical Insight: Producer-Consumer Imbalance

**Core problem**: Producers create work faster than consumers can process it.

Without backpressure:
```
Producers → [∞ queue] → Consumers
  Fast         Grows      Slow
               Forever
```

With backpressure:
```
Producers → [Bounded Queue] → Consumers
  Fast       Fixed Size         Slow
             ↑
      Pushes back on producers
```

**In Go terms**:
```go
// ❌ UNBOUNDED: Will consume all memory eventually
jobs := make(chan Job)  // No limit!

for i := 0; i < 1000000; i++ {
    jobs <- Job{ID: i}  // Never blocks, queue grows infinitely
}

// ✅ BOUNDED: Queue has fixed capacity
jobs := make(chan Job, 100)  // Max 100 items

for i := 0; i < 1000000; i++ {
    jobs <- Job{ID: i}  // Blocks when queue is full (backpressure!)
}
```

---

## Problem 1: Understanding Channel Capacity and Blocking

### How Buffered Channels Provide Backpressure

**Unbuffered channel** (capacity 0):
```go
ch := make(chan int)  // No buffer
ch <- 42  // BLOCKS immediately until someone receives
```

**Analogy**: Handoff between two people. Sender can't let go until receiver grabs it.

**Buffered channel** (capacity > 0):
```go
ch := make(chan int, 3)  // Buffer size 3
ch <- 1  // Doesn't block (buffer: [1])
ch <- 2  // Doesn't block (buffer: [1, 2])
ch <- 3  // Doesn't block (buffer: [1, 2, 3])
ch <- 4  // BLOCKS! (buffer full)
```

**Analogy**: Conveyor belt with limited space. Can add items until belt is full.

### The Blocking Behavior

When a buffered channel is full, **send operations block**:

```go
jobs := make(chan int, 2)

// Producer
go func() {
    for i := 0; i < 10; i++ {
        fmt.Printf("Sending %d...\n", i)
        jobs <- i  // Will block when buffer is full
        fmt.Printf("Sent %d\n", i)
    }
    close(jobs)
}()

// Slow consumer
time.Sleep(1 * time.Second)  // Simulate slow startup
for job := range jobs {
    fmt.Printf("Processing %d\n", job)
    time.Sleep(500 * time.Millisecond)  // Slow processing
}
```

**Output**:
```
Sending 0...
Sent 0
Sending 1...
Sent 1
Sending 2...
Sent 2
Sending 3...
[Blocks here - buffer full]
Processing 0
Sent 3
Sending 4...
[Blocks again]
Processing 1
...
```

**Key insight**: The producer automatically slows down to match consumer speed. This is **automatic backpressure**.

---

## Problem 2: Backpressure Strategies

There are **four main strategies** for handling backpressure:

### Strategy 1: Block (Wait Until Space Available)

**What**: Producer blocks until consumer catches up.

**When to use**:
- Work must not be lost
- Producer can afford to wait
- System should self-regulate

**Example**: Database transaction queue

```go
func blockStrategy(jobs chan<- int, job int) error {
    jobs <- job  // Blocks if full (built-in behavior)
    return nil
}
```

**Pros**:
- Simple (no code needed)
- No data loss
- Automatic flow control

**Cons**:
- Producer can stall indefinitely
- Latency increases under load
- Can cause cascading slowdowns

### Strategy 2: Drop (Discard New Items)

**What**: If queue is full, drop the new item (don't add it).

**When to use**:
- Losing some data is acceptable
- Newer data is less important than processing current data
- Metrics, logs, telemetry

**Example**: Metrics aggregation

```go
func dropStrategy(jobs chan<- int, job int) error {
    select {
    case jobs <- job:
        return nil  // Sent successfully
    default:
        return fmt.Errorf("queue full, dropped job %d", job)
    }
}
```

**Pros**:
- Producer never blocks
- System remains responsive
- Prevents queue growth

**Cons**:
- Data loss
- No feedback to producer
- Can lose important work

### Strategy 3: Reject (Refuse New Items, Return Error)

**What**: If queue is full, reject new items and inform the producer.

**When to use**:
- Producer can retry later
- Need to signal overload to caller
- HTTP servers (return 503 Service Unavailable)

**Example**: HTTP request handler

```go
func rejectStrategy(jobs chan<- int, job int) error {
    select {
    case jobs <- job:
        return nil  // Accepted
    default:
        return fmt.Errorf("queue full, request rejected")
    }
}
```

**Pros**:
- Producer knows about failure
- Can retry or handle gracefully
- System stays stable

**Cons**:
- Caller must handle errors
- May need retry logic
- User-visible failures

### Strategy 4: Timeout (Wait for Limited Time)

**What**: Try to send for a limited time, then give up.

**When to use**:
- Want to wait a bit, but not forever
- SLAs require bounded latency
- Hybrid between blocking and rejecting

**Example**: SLA-constrained request processing

```go
func timeoutStrategy(jobs chan<- int, job int, timeout time.Duration) error {
    select {
    case jobs <- job:
        return nil  // Sent successfully
    case <-time.After(timeout):
        return fmt.Errorf("timeout sending job %d", job)
    }
}
```

**Pros**:
- Bounded waiting time
- Balance between blocking and dropping
- Predictable latency

**Cons**:
- Still can fail
- More complex logic
- Timer allocation overhead

---

## Problem 3: Worker Pool with Bounded Channels

### The Complete Pattern

```go
type WorkerPool struct {
    jobs    chan Job
    results chan Result
    workers int
}

func NewWorkerPool(queueSize, numWorkers int) *WorkerPool {
    return &WorkerPool{
        jobs:    make(chan Job, queueSize),  // Bounded queue
        results: make(chan Result, queueSize),
        workers: numWorkers,
    }
}

func (p *WorkerPool) Start(ctx context.Context) {
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < p.workers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case job, ok := <-p.jobs:
                    if !ok {
                        return  // Jobs channel closed
                    }
                    result := processJob(job)
                    p.results <- result
                }
            }
        }(i)
    }

    // Close results when all workers done
    go func() {
        wg.Wait()
        close(p.results)
    }()
}

func (p *WorkerPool) Submit(job Job) error {
    select {
    case p.jobs <- job:
        return nil  // Submitted successfully
    default:
        return ErrQueueFull  // Backpressure applied
    }
}

func (p *WorkerPool) Results() <-chan Result {
    return p.results
}

func (p *WorkerPool) Close() {
    close(p.jobs)
}
```

**Key components**:
1. **Bounded `jobs` channel**: Limits queue size
2. **Non-blocking submit**: Uses `select` with `default`
3. **Worker goroutines**: Process jobs concurrently
4. **Results channel**: Return processed results
5. **Graceful shutdown**: Close channels, wait for workers

---

## Problem 4: Overflow Handling Patterns

### Pattern 1: Buffering with Spillover Storage

When channel is full, store excess items in secondary storage:

```go
type BufferedPool struct {
    jobs     chan Job
    overflow []Job
    mu       sync.Mutex
}

func (p *BufferedPool) Submit(job Job) error {
    select {
    case p.jobs <- job:
        return nil  // Sent to channel
    default:
        p.mu.Lock()
        p.overflow = append(p.overflow, job)
        p.mu.Unlock()
        return nil  // Stored in overflow
    }
}

func (p *BufferedPool) drainOverflow() {
    for {
        p.mu.Lock()
        if len(p.overflow) == 0 {
            p.mu.Unlock()
            time.Sleep(100 * time.Millisecond)
            continue
        }
        job := p.overflow[0]
        p.overflow = p.overflow[1:]
        p.mu.Unlock()

        p.jobs <- job  // Blocks until space available
    }
}
```

**Use case**: Temporary bursts that exceed capacity.

### Pattern 2: Priority Queue

Process important jobs first when queue is full:

```go
type PriorityPool struct {
    highPriority chan Job
    lowPriority  chan Job
}

func (p *PriorityPool) worker() {
    for {
        select {
        case job := <-p.highPriority:
            process(job)
        default:
            select {
            case job := <-p.highPriority:
                process(job)  // High priority still checked
            case job := <-p.lowPriority:
                process(job)
            }
        }
    }
}
```

**Use case**: Critical vs. non-critical work.

### Pattern 3: Adaptive Worker Scaling

Increase workers when queue fills up:

```go
type AdaptivePool struct {
    jobs       chan Job
    workers    int
    maxWorkers int
    mu         sync.Mutex
}

func (p *AdaptivePool) monitor() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        queueLoad := float64(len(p.jobs)) / float64(cap(p.jobs))

        if queueLoad > 0.8 && p.workers < p.maxWorkers {
            p.addWorker()
        } else if queueLoad < 0.2 && p.workers > 1 {
            // Signal worker to stop (implementation omitted)
        }
    }
}
```

**Use case**: Variable load patterns.

---

## Problem 5: Rate Limiting (Controlling Throughput)

### Token Bucket Algorithm

```go
type RateLimiter struct {
    tokens   chan struct{}
    rate     time.Duration
    capacity int
}

func NewRateLimiter(rate time.Duration, capacity int) *RateLimiter {
    rl := &RateLimiter{
        tokens:   make(chan struct{}, capacity),
        rate:     rate,
        capacity: capacity,
    }

    // Fill bucket initially
    for i := 0; i < capacity; i++ {
        rl.tokens <- struct{}{}
    }

    // Refill tokens periodically
    go func() {
        ticker := time.NewTicker(rate)
        defer ticker.Stop()

        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
                // Added token
            default:
                // Bucket full
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil  // Got token
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (rl *RateLimiter) TryAcquire() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}
```

**Usage with worker pool**:

```go
limiter := NewRateLimiter(100*time.Millisecond, 10)  // 10 req/sec

for _, job := range jobs {
    if err := limiter.Wait(ctx); err != nil {
        return err
    }
    pool.Submit(job)
}
```

---

## Problem 6: Combining Strategies

### Hybrid: Timeout + Retry with Exponential Backoff

```go
func submitWithRetry(pool *WorkerPool, job Job, maxRetries int) error {
    backoff := 10 * time.Millisecond

    for attempt := 0; attempt < maxRetries; attempt++ {
        select {
        case pool.jobs <- job:
            return nil  // Success
        case <-time.After(backoff):
            backoff *= 2  // Exponential backoff
            if backoff > 1*time.Second {
                backoff = 1 * time.Second  // Cap at 1s
            }
        }
    }

    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

### Circuit Breaker Pattern

Stop sending requests if failures exceed threshold:

```go
type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    failures     int
    lastFailure  time.Time
    mu           sync.Mutex
    state        string  // "closed", "open", "half-open"
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.resetTimeout {
            cb.state = "half-open"
        } else {
            return fmt.Errorf("circuit breaker open")
        }
    }

    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }

    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

---

## Real-World Applications

### 1. HTTP Server with Backpressure

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    job := parseRequest(r)

    err := pool.Submit(job)
    if err == ErrQueueFull {
        http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
        return
    }

    result := <-pool.Results()
    json.NewEncoder(w).Encode(result)
}
```

**Real systems**: NGINX, HAProxy, Envoy all implement backpressure.

### 2. Message Queue Consumer

```go
func consumeMessages(queue MessageQueue, pool *WorkerPool) {
    for msg := range queue.Messages() {
        err := pool.Submit(msg)
        if err != nil {
            // Message stays in queue (not acknowledged)
            log.Printf("Queue full, will retry: %v", err)
            time.Sleep(1 * time.Second)
        } else {
            msg.Ack()  // Remove from queue
        }
    }
}
```

**Real systems**: RabbitMQ, Kafka, SQS support backpressure via ack delays.

### 3. Database Connection Pool

```go
type DBPool struct {
    conns chan *sql.DB
}

func (p *DBPool) Acquire(ctx context.Context) (*sql.DB, error) {
    select {
    case conn := <-p.conns:
        return conn, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func (p *DBPool) Release(conn *sql.DB) {
    select {
    case p.conns <- conn:
    default:
        conn.Close()  // Pool full, close connection
    }
}
```

**Real systems**: PgBouncer, connection pooling in database drivers.

### 4. Log Aggregation with Sampling

```go
type LogAggregator struct {
    logs       chan LogEntry
    sampleRate float64
}

func (la *LogAggregator) Log(entry LogEntry) {
    if la.shouldSample() {
        select {
        case la.logs <- entry:
            // Logged
        default:
            // Dropped (backpressure)
        }
    }
}

func (la *LogAggregator) shouldSample() bool {
    return rand.Float64() < la.sampleRate
}
```

**Real systems**: Datadog, Splunk use sampling under load.

---

## Common Mistakes to Avoid

### Mistake 1: Unbounded Queues in Production

```go
// ❌ BAD: Will consume all memory under load
jobs := make(chan Job)  // Unbounded

// ✅ GOOD: Fixed capacity
jobs := make(chan Job, 1000)
```

### Mistake 2: Ignoring Backpressure Signals

```go
// ❌ BAD: Blindly blocks, can stall system
for _, job := range jobs {
    pool.jobs <- job  // Blocks if full!
}

// ✅ GOOD: Handle backpressure
for _, job := range jobs {
    if err := pool.Submit(job); err != nil {
        handleOverload(job)  // Drop, reject, or retry
    }
}
```

### Mistake 3: No Monitoring

```go
// ✅ GOOD: Monitor queue depth
func (p *WorkerPool) QueueDepth() int {
    return len(p.jobs)
}

func (p *WorkerPool) QueueUtilization() float64 {
    return float64(len(p.jobs)) / float64(cap(p.jobs))
}

// Alert when utilization > 80%
if pool.QueueUtilization() > 0.8 {
    log.Warn("Queue filling up, consider scaling")
}
```

### Mistake 4: Wrong Queue Size

```go
// ❌ TOO SMALL: Excessive backpressure
jobs := make(chan Job, 10)  // Only 10 items

// ❌ TOO LARGE: Defeats purpose of backpressure
jobs := make(chan Job, 1000000)  // 1M items!

// ✅ GOOD: Based on latency tolerance and throughput
// Queue should hold: (requests/sec) * (acceptable delay in seconds)
// Example: 1000 req/s * 0.1s delay = 100 item buffer
jobs := make(chan Job, 100)
```

---

## Performance Characteristics

### Queue Size Selection

**Formula**: `Queue Size = Throughput × Acceptable Delay`

**Example**:
- Throughput: 500 requests/second
- Acceptable delay: 200ms (0.2 seconds)
- Queue size: 500 × 0.2 = **100 items**

**Trade-offs**:
- **Small queue**: More backpressure, lower latency, less memory
- **Large queue**: Less backpressure, higher latency, more memory

### Worker Count Selection

**Formula**: `Workers = (CPU Cores) × (1 + Wait/Compute Ratio)`

**I/O-bound work** (e.g., HTTP requests):
- Wait/Compute ratio: 10:1 (90% waiting, 10% computing)
- Workers: 8 cores × (1 + 10) = **88 workers**

**CPU-bound work** (e.g., image processing):
- Wait/Compute ratio: 0:1 (no waiting)
- Workers: 8 cores × (1 + 0) = **8 workers**

**Mixed workload**:
- Experiment and measure
- Monitor CPU and queue depth
- Adjust based on metrics

---

## How to Run

```bash
# Run the demonstration
go run ./minis/22-worker-pool-with-backpressure/cmd/worker-pool/main.go

# Run tests
go test ./minis/22-worker-pool-with-backpressure/...

# Run with race detector
go test -race ./minis/22-worker-pool-with-backpressure/...

# Benchmark different strategies
go test -bench=. ./minis/22-worker-pool-with-backpressure/...
```

---

## Summary

**What you learned**:
- ✅ Unbounded queues cause memory exhaustion and crashes
- ✅ Bounded channels provide automatic backpressure
- ✅ Four backpressure strategies: block, drop, reject, timeout
- ✅ Overflow handling: spillover storage, priority queues, adaptive scaling
- ✅ Rate limiting controls throughput (token bucket algorithm)
- ✅ Queue sizing: based on throughput and acceptable latency
- ✅ Worker count: based on CPU cores and I/O ratio

**Key principles**:
1. **Always use bounded queues** in production
2. **Handle backpressure explicitly** (don't ignore it)
3. **Monitor queue depth** and latency
4. **Choose strategy based on requirements**:
   - Block: No data loss allowed
   - Drop: Lossy metrics/logs acceptable
   - Reject: User-facing with retry capability
   - Timeout: SLA-driven latency bounds
5. **Size queues based on math**, not guesswork

**Next steps**:
- Project 23: Bounded channel semaphore patterns
- Project 24: Mutex vs RWMutex for shared state
- Project 25: Atomic counters for lock-free programming

Build resilient systems! 🚀
