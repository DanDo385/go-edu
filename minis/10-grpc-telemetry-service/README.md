# Project 10: gRPC Telemetry Service

## 1. What Is This About?

### Real-World Scenario

You're building a monitoring system that collects metrics from thousands of servers:
- CPU usage, memory, disk I/O
- Need to collect 10,000 metrics per second
- Must aggregate in real-time (count, sum, avg, min, max)
- Query latest statistics on demand

**❌ REST API approach** (inefficient):
```
For each metric:
  HTTP POST /metrics
  {json payload}
  HTTP overhead: headers, connection setup, JSON parsing

10,000 requests/sec = massive overhead
```

**✅ gRPC streaming approach** (efficient):
```
Open one connection
Stream metrics continuously:
  Point{metric: "cpu", value: 75.2}
  Point{metric: "cpu", value: 80.1}
  Point{metric: "mem", value: 4096}
  ...

One connection, binary protocol, minimal overhead
```

This project teaches you **modern RPC** with:
- **gRPC**: High-performance RPC framework
- **Protocol Buffers**: Efficient binary serialization
- **Streaming**: Send many messages over one connection
- **Concurrency**: Thread-safe aggregation
- **Production patterns**: Time-windowed data, statistics

### What You'll Learn

1. **Protocol Buffers**: Define services in .proto files
2. **gRPC streaming**: Client-side streaming RPC
3. **Thread-safe aggregation**: Concurrent updates with mutexes
4. **Time windows**: Expire old data automatically
5. **Statistics**: Count, sum, avg, min, max calculation
6. **Benchmarking**: Measure gRPC throughput

### The Challenge

Build a telemetry service that:
- Accepts streaming metrics via gRPC
- Aggregates per-metric statistics (count, sum, avg, min, max)
- Expires metrics after time window (e.g., 1 hour)
- Provides query endpoint for current statistics
- Handles concurrent updates safely
- Processes 100,000+ points per second

---

## 2. First Principles: gRPC and Protocol Buffers

### What is RPC?

**RPC (Remote Procedure Call)** lets you call functions on remote servers as if they were local.

**Without RPC** (manual HTTP):
```go
// Client
url := "http://server/add?a=5&b=3"
resp, _ := http.Get(url)
body, _ := io.ReadAll(resp.Body)
result, _ := strconv.Atoi(string(body))

// Server
func handleAdd(w http.ResponseWriter, r *http.Request) {
    a, _ := strconv.Atoi(r.URL.Query().Get("a"))
    b, _ := strconv.Atoi(r.URL.Query().Get("b"))
    result := a + b
    fmt.Fprintf(w, "%d", result)
}
```

**With RPC** (looks like local call):
```go
// Client
result, _ := client.Add(context.Background(), &AddRequest{A: 5, B: 3})

// Server
func (s *Server) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
    return &AddResponse{Result: req.A + req.B}, nil
}
```

**Key insight**: RPC abstracts away networking details.

### What is Protocol Buffers?

**Protocol Buffers (protobuf)** is a language for defining data structures and services.

**Example** (.proto file):
```protobuf
message Point {
    string metric = 1;
    double value = 2;
    int64 timestamp = 3;
}
```

**Compiled to Go**:
```go
type Point struct {
    Metric    string
    Value     float64
    Timestamp int64
}
```

**Why use protobuf?**
- **Language-agnostic**: Generate code for Go, Python, Java, etc.
- **Binary format**: Smaller than JSON (60-80% size reduction)
- **Fast**: 5-10x faster than JSON
- **Schema**: Type-checked at compile time
- **Versioning**: Add fields without breaking old clients

### What is gRPC Streaming?

gRPC supports **four types of RPC**:

#### 1. Unary (like REST)
```
Client → Request → Server
Client ← Response ← Server
```

```protobuf
rpc GetUser(UserRequest) returns (UserResponse);
```

#### 2. Server Streaming
```
Client → Request → Server
Client ← Stream ← Server (many messages)
```

```protobuf
rpc ListUsers(Empty) returns (stream User);
```

**Use case**: Server sends updates continuously (stock prices, notifications)

#### 3. Client Streaming (This Project)
```
Client → Stream → Server (many messages)
Client ← Response ← Server
```

```protobuf
rpc PushPoints(stream Point) returns (Ack);
```

**Use case**: Client uploads many items in batch (telemetry, logs)

#### 4. Bidirectional Streaming
```
Client ↔ Stream ↔ Server (both send many messages)
```

```protobuf
rpc Chat(stream Message) returns (stream Message);
```

**Use case**: Real-time chat, gaming

### Why Client Streaming for Telemetry?

**Without streaming** (unary RPC):
```
For each metric:
  1. Open connection
  2. Send request
  3. Wait for response
  4. Close connection
```

10,000 metrics = 10,000 connections = slow!

**With streaming**:
```
1. Open one connection
2. Send 10,000 metrics
3. Receive one acknowledgment
4. Close connection
```

10,000 metrics = 1 connection = fast!

**Performance comparison**:
- Unary: ~1,000 requests/sec
- Streaming: ~100,000 points/sec

**100x faster!**

---

## 3. Breaking Down the Solution

### Step 1: Define Protocol Buffers

```protobuf
syntax = "proto3";

package telemetry;

message Point {
    string metric = 1;      // Metric name (e.g., "cpu", "memory")
    double value = 2;       // Metric value (e.g., 75.5)
    int64 timestamp = 3;    // Unix timestamp
}

message Ack {
    int32 points_received = 1;
}

message Empty {}

message MetricSummary {
    int32 count = 1;
    double sum = 2;
    double avg = 3;
    double min = 4;
    double max = 5;
}

message Report {
    map<string, MetricSummary> metrics = 1;
}

service TelemetryService {
    rpc PushPoints(stream Point) returns (Ack);
    rpc Summary(Empty) returns (Report);
}
```

**Key points**:
- `stream Point`: Client sends multiple Points
- `PushPoints` returns single `Ack`
- `Summary` returns current statistics

### Step 2: Data Structures

```go
type Aggregator struct {
    mu     sync.Mutex
    window time.Duration
    data   map[string]*metricData
}

type metricData struct {
    points []dataPoint
}

type dataPoint struct {
    value     float64
    timestamp time.Time
}
```

**Why slice of points?**
Need to remember all points within time window to calculate statistics.

### Step 3: PushPoint Logic

```
1. Lock mutex (thread safety)
2. Get or create metricData for this metric
3. Append point to slice
4. Evict points older than window
5. Unlock mutex
```

### Step 4: Summary Logic

```
1. Lock mutex
2. For each metric:
   a. Evict old points
   b. Calculate count, sum, avg, min, max
3. Build Report
4. Unlock mutex
5. Return Report
```

### Step 5: gRPC Server

```go
func (s *TelemetryServer) PushPoints(stream pb.TelemetryService_PushPointsServer) error {
    count := 0
    for {
        point, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.Ack{PointsReceived: int32(count)})
        }
        if err != nil {
            return err
        }

        s.agg.PushPoint(context.Background(), point)
        count++
    }
}
```

**How streaming works**:
1. `stream.Recv()` reads next Point from client
2. Returns `io.EOF` when client closes stream
3. `stream.SendAndClose()` sends response and closes

---

## 4. Complete Solution Walkthrough

### Aggregator Structure

```go
type Aggregator struct {
    mu     sync.Mutex
    window time.Duration
    data   map[string]*metricData
}
```

**Fields**:
- `mu`: Protects `data` from concurrent access
- `window`: How long to keep metrics (e.g., 1 hour)
- `data`: Map of metric name → data points

### PushPoint Implementation

```go
func (a *Aggregator) PushPoint(ctx context.Context, p *pb.Point) {
    a.mu.Lock()
    defer a.mu.Unlock()

    // Get or create metric data
    md, ok := a.data[p.Metric]
    if !ok {
        md = &metricData{points: make([]dataPoint, 0)}
        a.data[p.Metric] = md
    }

    // Add point
    md.points = append(md.points, dataPoint{
        value:     p.Value,
        timestamp: time.Unix(p.Timestamp, 0),
    })

    // Evict old points
    md.evictOld(a.window)
}
```

**Line-by-line**:
1. **Lock**: Ensures only one goroutine modifies `data` at a time
2. **Get or create**: If metric doesn't exist, create empty slice
3. **Append**: Add new point to slice
4. **Evict**: Remove points older than window

### Evict Old Points

```go
func (md *metricData) evictOld(window time.Duration) {
    cutoff := time.Now().Add(-window)

    // Find first point within window
    i := 0
    for i < len(md.points) && md.points[i].timestamp.Before(cutoff) {
        i++
    }

    // Keep only points after cutoff
    md.points = md.points[i:]
}
```

**Algorithm**:
```
Window = 1 hour
Now = 14:00

Points:
  [13:00, 75.0]  ← 1 hour old, remove
  [13:30, 80.0]  ← 30min old, remove
  [13:45, 85.0]  ← 15min old, keep
  [14:00, 90.0]  ← now, keep

After eviction: [[13:45, 85.0], [14:00, 90.0]]
```

**Why not delete one by one?**
Slice re-slicing (`md.points = md.points[i:]`) is O(1), very fast.

### Summary Implementation

```go
func (a *Aggregator) Summary(ctx context.Context) *pb.Report {
    a.mu.Lock()
    defer a.mu.Unlock()

    report := &pb.Report{
        Metrics: make(map[string]*pb.MetricSummary),
    }

    for name, md := range a.data {
        md.evictOld(a.window)

        if len(md.points) == 0 {
            continue
        }

        sum := 0.0
        min := md.points[0].value
        max := md.points[0].value

        for _, p := range md.points {
            sum += p.value
            if p.value < min {
                min = p.value
            }
            if p.value > max {
                max = p.value
            }
        }

        avg := sum / float64(len(md.points))

        report.Metrics[name] = &pb.MetricSummary{
            Count: int32(len(md.points)),
            Sum:   sum,
            Avg:   avg,
            Min:   min,
            Max:   max,
        }
    }

    return report
}
```

**Statistics calculation**:
- **Count**: `len(md.points)`
- **Sum**: Add all values
- **Avg**: `sum / count`
- **Min/Max**: Track while iterating

### gRPC Server Implementation

```go
type TelemetryServer struct {
    agg *Aggregator
}

func (s *TelemetryServer) PushPoints(stream pb.TelemetryService_PushPointsServer) error {
    count := 0

    for {
        point, err := stream.Recv()
        if err == io.EOF {
            // Client finished sending
            ack := &pb.Ack{PointsReceived: int32(count)}
            return stream.SendAndClose(ack)
        }
        if err != nil {
            return err
        }

        s.agg.PushPoint(context.Background(), point)
        count++
    }
}

func (s *TelemetryServer) Summary(ctx context.Context, _ *pb.Empty) (*pb.Report, error) {
    return s.agg.Summary(ctx), nil
}
```

**Streaming flow**:
```
Client                    Server
  |                         |
  |--- Point{cpu, 75} ---→|
  |                         | PushPoint("cpu", 75)
  |--- Point{cpu, 80} ---→|
  |                         | PushPoint("cpu", 80)
  |--- Point{mem, 4096} →  |
  |                         | PushPoint("mem", 4096)
  |--- EOF -------------→  |
  |                         | SendAndClose(Ack{3})
  |←-- Ack{3} -----------  |
```

---

## 5. Key Concepts Explained

### Concept 1: Protocol Buffers Field Numbers

```protobuf
message Point {
    string metric = 1;
    double value = 2;
    int64 timestamp = 3;
}
```

**What do the numbers mean?**
- Not default values!
- **Field tags** for binary encoding
- Can't change once deployed (breaks compatibility)

**Binary encoding**:
```
Point{metric: "cpu", value: 75.0, timestamp: 1234567890}

Binary (simplified):
  Field 1 (string): "cpu"
  Field 2 (double): 75.0
  Field 3 (int64): 1234567890
```

**Why important?**
If you change `double value = 2` to `double value = 5`, old clients will break!

### Concept 2: gRPC vs REST

| Feature | REST | gRPC |
|---------|------|------|
| Protocol | HTTP/1.1 | HTTP/2 |
| Format | JSON (text) | Protobuf (binary) |
| Size | Larger | 60-80% smaller |
| Speed | Slower | 5-10x faster |
| Streaming | Awkward (SSE, WebSocket) | Native |
| Schema | Optional (OpenAPI) | Required (.proto) |
| Browser support | Yes | Limited (gRPC-Web) |

**When to use gRPC**:
- Microservices (high throughput)
- Real-time systems
- Mobile apps (battery efficient)

**When to use REST**:
- Public APIs (wider compatibility)
- Browser-based apps
- Simple CRUD

### Concept 3: Time-Windowed Data

**Sliding window**:
```
Window = 1 hour

T=13:00: Points = []
T=13:30: Add point → Points = [13:30]
T=14:00: Add point → Points = [13:30, 14:00]
T=14:30: Add point, evict old → Points = [14:00, 14:30]  (13:30 removed)
```

**Why evict?**
- **Memory**: Don't store metrics forever
- **Relevance**: Old data less useful for monitoring
- **Performance**: Smaller slices = faster iteration

### Concept 4: Thread Safety with Mutexes

**Problem**:
```go
// Goroutine 1
data[metric] = append(data[metric], point1)

// Goroutine 2 (same time)
data[metric] = append(data[metric], point2)

// RACE CONDITION! One append might be lost
```

**Solution**:
```go
mu.Lock()
data[metric] = append(data[metric], point)
mu.Unlock()
```

**Why defer?**
```go
mu.Lock()
defer mu.Unlock()  // Ensures unlock even if panic
data[metric] = append(data[metric], point)
```

### Concept 5: gRPC Error Handling

**gRPC status codes**:
- `OK`: Success
- `NOT_FOUND`: Resource doesn't exist
- `INVALID_ARGUMENT`: Bad request
- `UNAVAILABLE`: Service down
- `DEADLINE_EXCEEDED`: Timeout

**Usage**:
```go
import "google.golang.org/grpc/codes"
import "google.golang.org/grpc/status"

func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*User, error) {
    user, err := s.db.GetUser(req.Id)
    if err == sql.ErrNoRows {
        return nil, status.Error(codes.NotFound, "user not found")
    }
    return user, nil
}
```

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Generic Streaming Aggregator

```go
type StreamAggregator[T any] struct {
    mu   sync.Mutex
    data []T
    fn   func([]T) any  // Aggregation function
}

func (sa *StreamAggregator[T]) Add(item T) {
    sa.mu.Lock()
    defer sa.mu.Unlock()
    sa.data = append(sa.data, item)
}

func (sa *StreamAggregator[T]) Result() any {
    sa.mu.Lock()
    defer sa.mu.Unlock()
    return sa.fn(sa.data)
}
```

### Pattern 2: Time-Series Database

```go
type TimeSeries struct {
    mu      sync.RWMutex
    data    []TimePoint
    maxAge  time.Duration
}

func (ts *TimeSeries) Add(value float64) {
    ts.mu.Lock()
    defer ts.mu.Unlock()

    ts.data = append(ts.data, TimePoint{
        value: value,
        time:  time.Now(),
    })

    ts.evictOld()
}

func (ts *TimeSeries) Stats() Stats {
    ts.mu.RLock()
    defer ts.mu.RUnlock()

    return calculateStats(ts.data)
}
```

### Pattern 3: gRPC Interceptor (Middleware)

```go
func LoggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    duration := time.Since(start)

    log.Printf("Method: %s, Duration: %v, Error: %v",
        info.FullMethod, duration, err)

    return resp, err
}

// Usage:
server := grpc.NewServer(
    grpc.UnaryInterceptor(LoggingInterceptor),
)
```

### Pattern 4: Backpressure Handling

```go
func (s *Server) PushPoints(stream pb.TelemetryService_PushPointsServer) error {
    const maxBuffer = 1000
    buffer := make([]*pb.Point, 0, maxBuffer)

    for {
        point, err := stream.Recv()
        if err == io.EOF {
            s.processBatch(buffer)
            return stream.SendAndClose(&pb.Ack{})
        }

        buffer = append(buffer, point)

        if len(buffer) >= maxBuffer {
            s.processBatch(buffer)
            buffer = buffer[:0]  // Reset
        }
    }
}
```

### Pattern 5: Metric Expiry with Timer

```go
type ExpiringCache struct {
    mu    sync.Mutex
    data  map[string]cacheEntry
}

type cacheEntry struct {
    value     interface{}
    expiresAt time.Time
    timer     *time.Timer
}

func (c *ExpiringCache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    timer := time.AfterFunc(ttl, func() {
        c.mu.Lock()
        delete(c.data, key)
        c.mu.Unlock()
    })

    c.data[key] = cacheEntry{
        value:     value,
        expiresAt: time.Now().Add(ttl),
        timer:     timer,
    }
}
```

---

## 7. Real-World Applications

### Application Performance Monitoring (APM)

Companies: Datadog, New Relic, Dynatrace

```go
type APMClient struct {
    client pb.TelemetryServiceClient
}

func (c *APMClient) ReportMetrics() {
    stream, _ := c.client.PushPoints(context.Background())

    for {
        metrics := collectSystemMetrics()
        for _, m := range metrics {
            stream.Send(&pb.Point{
                Metric:    m.Name,
                Value:     m.Value,
                Timestamp: time.Now().Unix(),
            })
        }
        time.Sleep(1 * time.Second)
    }
}
```

### Distributed Tracing

Companies: Jaeger, Zipkin, Honeycomb

```go
func (s *TracingService) ReportSpans(stream pb.TracingService_ReportSpansServer) error {
    for {
        span, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.Ack{})
        }

        s.storage.StoreSpan(span)
    }
}
```

### Log Aggregation

Companies: Splunk, Elasticsearch, Loki

```go
type LogAggregator struct {
    client pb.LogServiceClient
}

func (la *LogAggregator) StreamLogs(logs <-chan LogEntry) {
    stream, _ := la.client.SendLogs(context.Background())

    for log := range logs {
        stream.Send(&pb.LogEntry{
            Level:     log.Level,
            Message:   log.Message,
            Timestamp: time.Now().Unix(),
        })
    }

    stream.CloseAndRecv()
}
```

### IoT Sensor Data

Companies: AWS IoT, Google Cloud IoT

```go
type SensorClient struct {
    client pb.TelemetryServiceClient
}

func (sc *SensorClient) StreamSensorData() {
    stream, _ := sc.client.PushPoints(context.Background())

    for {
        temperature := readTemperatureSensor()
        humidity := readHumiditySensor()

        stream.Send(&pb.Point{Metric: "temperature", Value: temperature})
        stream.Send(&pb.Point{Metric: "humidity", Value: humidity})

        time.Sleep(5 * time.Second)
    }
}
```

### Financial Trading Systems

High-frequency trading needs ultra-low latency.

```go
type MarketDataFeed struct {
    client pb.MarketDataServiceClient
}

func (mdf *MarketDataFeed) StreamPrices() {
    stream, _ := mdf.client.PushPrices(context.Background())

    for price := range mdf.priceChannel {
        stream.Send(&pb.Price{
            Symbol:    price.Symbol,
            Price:     price.Value,
            Timestamp: time.Now().UnixNano(),
        })
    }
}
```

**Why gRPC?** 10x lower latency than REST.

---

## 8. Common Mistakes to Avoid

### Mistake 1: Not Handling EOF Correctly

**❌ Wrong**:
```go
func (s *Server) PushPoints(stream pb.Service_PushPointsServer) error {
    for {
        point, err := stream.Recv()
        if err != nil {
            return err  // Returns error on EOF!
        }
        process(point)
    }
}
```

**✅ Correct**:
```go
for {
    point, err := stream.Recv()
    if err == io.EOF {
        return stream.SendAndClose(&pb.Ack{})
    }
    if err != nil {
        return err
    }
    process(point)
}
```

### Mistake 2: Forgetting Mutex

**❌ Wrong**:
```go
func (a *Aggregator) PushPoint(p *pb.Point) {
    a.data[p.Metric] = append(a.data[p.Metric], p.Value)
    // DATA RACE!
}
```

**✅ Correct**:
```go
func (a *Aggregator) PushPoint(p *pb.Point) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.data[p.Metric] = append(a.data[p.Metric], p.Value)
}
```

### Mistake 3: Not Evicting Old Data

**❌ Wrong**:
```go
func (a *Aggregator) PushPoint(p *pb.Point) {
    a.data[p.Metric] = append(a.data[p.Metric], p)
    // Memory grows unbounded!
}
```

**✅ Correct**:
```go
func (a *Aggregator) PushPoint(p *pb.Point) {
    a.data[p.Metric] = append(a.data[p.Metric], p)
    a.evictOld()  // Remove old points
}
```

### Mistake 4: Changing Proto Field Numbers

**❌ Wrong**:
```protobuf
// Version 1
message Point {
    string metric = 1;
    double value = 2;
}

// Version 2 - BREAKING CHANGE!
message Point {
    string metric = 2;  // Changed from 1!
    double value = 1;   // Changed from 2!
}
```

**✅ Correct**: Never change field numbers. Add new fields only:
```protobuf
// Version 2
message Point {
    string metric = 1;
    double value = 2;
    string unit = 3;  // NEW field
}
```

### Mistake 5: Not Closing Stream

**❌ Wrong**:
```go
stream, _ := client.PushPoints(context.Background())
for _, point := range points {
    stream.Send(point)
}
// Forgot to close!
```

**✅ Correct**:
```go
stream, _ := client.PushPoints(context.Background())
for _, point := range points {
    stream.Send(point)
}
ack, _ := stream.CloseAndRecv()  // Close and receive response
```

---

## 9. Stretch Goals

### Goal 1: Add Prometheus Exporter ⭐⭐

Export metrics in Prometheus format.

**Hint**:
```go
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    report := s.agg.Summary(context.Background())

    for name, summary := range report.Metrics {
        fmt.Fprintf(w, "# TYPE %s_count counter\n", name)
        fmt.Fprintf(w, "%s_count %d\n", name, summary.Count)
        fmt.Fprintf(w, "# TYPE %s_sum counter\n", name)
        fmt.Fprintf(w, "%s_sum %f\n", name, summary.Sum)
    }
}
```

### Goal 2: Server-Side Streaming ⭐⭐⭐

Stream live updates to clients.

**Hint**:
```protobuf
rpc StreamMetrics(Empty) returns (stream MetricUpdate);
```

```go
func (s *Server) StreamMetrics(_ *pb.Empty, stream pb.Service_StreamMetricsServer) error {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        report := s.agg.Summary(context.Background())
        stream.Send(&pb.MetricUpdate{Report: report})
    }
}
```

### Goal 3: Add Histogram/Percentiles ⭐⭐⭐

Track p50, p95, p99 latencies.

**Hint**:
```go
import "github.com/montanaflynn/stats"

func (a *Aggregator) Percentiles(metric string) (p50, p95, p99 float64) {
    values := extractValues(a.data[metric])
    p50, _ = stats.Percentile(values, 50)
    p95, _ = stats.Percentile(values, 95)
    p99, _ = stats.Percentile(values, 99)
    return
}
```

### Goal 4: Persistent Storage ⭐⭐⭐

Store metrics to disk for historical queries.

**Hint**:
```go
type PersistentAggregator struct {
    agg *Aggregator
    db  *sql.DB
}

func (pa *PersistentAggregator) PushPoint(p *pb.Point) {
    pa.agg.PushPoint(p)

    // Async write to DB
    go func() {
        pa.db.Exec("INSERT INTO metrics VALUES (?, ?, ?)",
            p.Metric, p.Value, p.Timestamp)
    }()
}
```

### Goal 5: Distributed Aggregation ⭐⭐⭐⭐

Shard metrics across multiple servers.

**Hint**:
```go
type ShardedAggregator struct {
    shards []*Aggregator
}

func (sa *ShardedAggregator) PushPoint(p *pb.Point) {
    shard := sa.getShard(p.Metric)
    shard.PushPoint(p)
}

func (sa *ShardedAggregator) getShard(metric string) *Aggregator {
    hash := fnv.New32a()
    hash.Write([]byte(metric))
    return sa.shards[hash.Sum32()%uint32(len(sa.shards))]
}
```

---

## How to Run

```bash
# Run tests
go test ./minis/10-grpc-telemetry-service/...

# Run benchmarks
go test -bench=. -benchmem ./minis/10-grpc-telemetry-service/...

# Example benchmark output:
# BenchmarkAggregator_Push-8       5000000    250 ns/op    48 B/op   1 allocs/op
# BenchmarkAggregator_Summary-8    1000000    1200 ns/op   256 B/op  3 allocs/op
```

---

## Summary

**What you learned**:
- ✅ Protocol Buffers for service definition
- ✅ gRPC client-side streaming
- ✅ Thread-safe concurrent aggregation
- ✅ Time-windowed data management
- ✅ Real-time statistics calculation
- ✅ High-performance RPC (100,000+ points/sec)

**Why this matters**:
gRPC is the modern standard for microservices communication. It's 10x faster than REST, supports streaming natively, and provides type-safe contracts. Used by Google, Netflix, Square, and thousands of companies.

**Key insights**:
- Streaming > Unary for high-volume data
- Protobuf > JSON for performance
- Time windows prevent memory growth
- Mutexes ensure thread safety

**gRPC Performance**:
- **Throughput**: 100,000+ messages/sec
- **Latency**: Sub-millisecond
- **Size**: 60-80% smaller than JSON

**Congratulations!**
You've completed all 10 projects. You now have a solid foundation in Go, from basics (strings, maps) to advanced topics (concurrency, gRPC, generics). You're ready for production Go development!

**Next steps**:
- Build your own projects
- Contribute to open source
- Read production Go codebases (Kubernetes, Docker, Terraform)

Keep building! 🚀
