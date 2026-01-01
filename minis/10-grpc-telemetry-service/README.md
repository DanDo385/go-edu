# 10: gRPC Telemetry Service

## What Is This Project About?

This module teaches you how to build a gRPC service for collecting telemetry data. You'll learn protocol buffers, gRPC streaming, and service design.

## Why Is This Important?

gRPC is used for:
- Microservice communication
- High-performance APIs
- Streaming data
- Cross-language services

## Key Concepts You'll Learn

- **Protocol Buffers**: Defining service contracts
- **gRPC**: High-performance RPC framework
- **Streaming**: Server and client streaming
- **Error handling**: gRPC status codes

## Prerequisites

- Completion of `minis/09-http-server-graceful`
- protoc compiler installed

## How to Run

```bash
go run ./cmd/dev/main.go
```

## Testing

```bash
go test -v ./...
```
