# 10-grpc-telemetry-service

**gRPC Telemetry Service**

Build a gRPC service for collecting telemetry data.

## What You'll Learn

- gRPC service definition
- Protocol Buffers
- Streaming RPC
- gRPC server/client

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement gRPC service | Handle telemetry data collection |

## Project Structure

```
10-grpc-telemetry-service/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/grpctelemetryservice/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
├── proto/
│   ├── telemetry.proto  # Protocol definition
│   └── telemetry.pb.go  # Generated code
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/10-grpc-telemetry-service

# Start gRPC server
go run ./cmd/app/main.go --port 50051
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Argument | Description |
|----------|-------------|
| `--port` | gRPC port (default: 50051) |

## Quick Copy & Paste

```bash
# Start server
go run ./cmd/app/main.go --port 50051

# Run benchmarks
go test -bench=. -benchmem ./internal/grpctelemetryservice/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Protocol Buffers**: Binary serialization
2. **gRPC**: High-performance RPC framework
3. **Streaming**: Client/server/bidirectional streams
4. **Code Generation**: protoc + go plugins

## Next Steps

After completing this exercise, proceed to `minis/11-slices-internals-capacity-growth`.
