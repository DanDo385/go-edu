# 10: Grpc Telemetry Service

## Core Concepts

- The concrete problem in Grpc Telemetry Service and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Grpc Telemetry Service patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for grpc telemetry service.

At this point in the arc:
Lesson 10 introduces a sharper systems concern so later modules can assume this mental model is stable.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Define the smallest valid behavior and reject invalid input or impossible state early.

### Step 2: Why This Approach
Pick a direct design that keeps control flow and data flow visible for debugging and testing.

### Step 3: Memory / Pointer Impact
Call out where data is copied versus aliased, and where mutable shared state needs synchronization.

### Step 4: What Changed
Produce a stable result shape and explicit error behavior that downstream code can rely on.

## Pointer and Indirection

- Explain * and & in this module when they appear in code or docs.
- Show memory-before and memory-after when data ownership changes.
- Clarify common misconceptions: Go stays pass-by-value even when pointer values are copied.
- Primer link: docs/MEMORY_POINTERS_PRIMER.md

## Verify


a) learner path


go test -v ./...


b) reference path


go test -tags=reference -v ./...


This project introduces you to **gRPC**, a high-performance RPC framework that is a cornerstone of modern microservice architectures. You will move beyond REST/JSON to build a contract-first, streaming API.

**IMPORTANT: This lesson requires external tools to be installed on your system.**

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## Prerequisites

Before you begin, you **must** install the Protocol Buffers compiler, `protoc`.

-   **On macOS (using Homebrew):**
    ```bash
    brew install protobuf
    ```
-   **On Linux (using apt):**
    ```bash
    sudo apt-get install -y protobuf-compiler
    ```
-   **On Windows (using Chocolatey):**
    ```bash
    choco install protoc
    ```

You also need the Go plugins for `protoc`:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

## What You'll Learn

- The advantages of **gRPC** over REST+JSON.
- How to define a service contract using **Protocol Buffers (`.proto`)**.
- How to **generate Go code** from `.proto` files.
- How to implement a gRPC server for a **client-side streaming** RPC.

## The Challenge: Beyond REST

For internal communication between microservices, REST over HTTP/1.1 can be slow and error-prone. gRPC is designed to be much faster and safer:
- **Performance:** It uses HTTP/2 and a compact binary format (Protocol Buffers) instead of text-based JSON.
- **Type Safety:** You define a strict "contract" in a `.proto` file. This contract is used to code-generate both the client and server, guaranteeing they are compatible.
- **Streaming:** gRPC has first-class support for streaming data, which is difficult with traditional REST.

## Core Concepts

### The "Contract First" Workflow
With gRPC, you always start by defining your API in a language-neutral `.proto` file. This is your single source of truth.

```protobuf
// In proto/telemetry.proto
service TelemetryService {
  // The client will stream TelemetryData messages to the server.
  // The server will respond with a single TelemetrySummary.
  rpc RecordTelemetry(stream TelemetryData) returns (TelemetrySummary);
}

message TelemetryData { double value = 1; }
message TelemetrySummary { int32 count = 1; double sum = 2; }
```

### Code Generation
You use `protoc` to generate Go code from your `.proto` file. In this project, you can run `go generate ./...` to trigger this process.

### Implementing the Server
The generated code gives you a `TelemetryServiceServer` interface. Your job is to create a struct that implements this interface.

## Your Task

Your task is to implement the `RecordTelemetry` method in `internal/grpctelemetryservice/server/exercise.go`.

1.  **Generate Code:** After installing the prerequisites, run `go generate ./...` from the root of the `go-edu` project.
2.  **Inspect the generated code:** Look at the `proto/telemetry.pb.go` and `proto/telemetry_grpc.pb.go` files to see the interfaces and structs that `protoc` created for you.
3.  **Implement `RecordTelemetry`:**
    - The function receives a `stream` argument.
    - You need to loop, calling `stream.Recv()` to get messages from the client.
    - The loop should end when `stream.Recv()` returns an `io.EOF` error.
    - As you receive messages, aggregate the data (e.g., count the messages and sum their values).
    - When the stream ends, use `stream.SendAndClose()` to send back a `TelemetrySummary` message.

## How to Verify Your Work

1.  **First, generate the Go code:**
    ```bash
    go generate ./...
    ```

2.  **Run the tests:**
    From this lesson's directory (`minis/10-grpc-telemetry-service`), run:
    ```bash
    go test -v ./...
    ```
    The tests will start a server with your implementation and use a client to test it.

## Related Lessons
- Previous: `minis/09-http-server-graceful`
- Next: `minis/11-slices-internals-capacity-growth`
