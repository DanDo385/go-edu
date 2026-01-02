# minis/10-grpc-telemetry-service

## Problem

Problem: Build a gRPC telemetry aggregator with streaming and time windows

Requirements:
1. Accept streaming telemetry points (metric name, value, timestamp)
2. Aggregate statistics per metric (count, sum, avg, min, max)
3. Support rolling time window (exclude old data)
4. Thread-safe concurrent access
5. gRPC service implementation

Algorithm: Rolling Time Window Aggregation
- Store points with timestamps
- Filter points within time window on Summary
- Calculate statistics on filtered points
- Use RWMutex for concurrent access

Time Window Algorithm:
- cutoff = now - window duration
- For each point: if timestamp > cutoff, include in stats
- This is "lazy" cleanup (filter on read, not proactive deletion)

Statistics Calculation:
- Iterate filtered points once
- Accumulate sum, track min/max
- Calculate average = sum / count

Concurrency Strategy:
- Write lock for PushPoint (modifies map/slices)
- Read lock for Summary (reads data, doesn't modify)
- RWMutex allows concurrent readers

## Quickstart

```bash
cd minis/10-grpc-telemetry-service
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/grpctelemetryservice/exercise.go`: implement the TODOs here
- `internal/grpctelemetryservice/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
