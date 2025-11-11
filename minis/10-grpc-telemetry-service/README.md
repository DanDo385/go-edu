# Project 10: grpc-telemetry-service

## What You're Building

A gRPC telemetry service with streaming metrics collection and real-time aggregation. This demonstrates modern RPC patterns and high-performance data processing.

## Concepts Covered

- Protocol Buffers (protobuf) for service definition
- gRPC streaming (client-side streaming)
- Thread-safe aggregation with mutexes
- Time-windowed data (rolling window)
- Benchmarking gRPC throughput

## How to Run

```bash
# Generate protobuf code (requires protoc and Go plugins)
# Already generated for you in this scaffold

# Run tests
go test ./minis/10-grpc-telemetry-service/...

# Run benchmarks
go test -bench=. -benchmem ./minis/10-grpc-telemetry-service/...
```

## Solution Explanation

### gRPC Streaming

**Client Streaming**: Client sends multiple messages, server responds once
- Push(stream Point) → Ack (batch upload)

**Unary**: Single request, single response
- Summary(Empty) → Report (query aggregates)

### Aggregation

Per-metric statistics:
- Count, Sum, Avg, Min, Max
- Rolling time window (old data expires)

## Stretch Goals

1. Add Prometheus exporter
2. Implement server-side streaming (live updates)
3. Add histogram/percentile tracking
