# 10: gRPC Telemetry Service

This project introduces you to **gRPC**, a high-performance, open-source universal RPC framework developed by Google. You will move beyond the request/response model of REST/HTTP and into the world of contract-first, streaming communication. You'll build a telemetry service where a client can stream a series of data points to a server, which processes them and returns a summary. This is a common pattern in monitoring and data-ingestion pipelines.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: Beyond REST/JSON](#the-big-picture-beyond-restjson)
- [First Principles: IDL and RPC](#first-principles-idl-and-rpc)
  - [Protocol Buffers (Protobuf)](#protocol-buffers-protobuf)
- [Project Structure](#project-structure)
- [Key Concepts in This Project](#key-concepts-in-this-project)
  - [The gRPC Workflow](#the-grpc-workflow)
  - [Defining the Service with `.proto`](#defining-the-service-with-proto)
  - [Code Generation with `protoc`](#code-generation-with-protoc)
  - [Implementing the Server](#implementing-the-server)
  - [Client-Side Streaming](#client-side-streaming)
- [Progression: The Microservice Backbone](#progression-the-microservice-backbone)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Explain the advantages of gRPC** over traditional REST+JSON for inter-service communication.
-   **Define a service contract** using Protocol Buffers (`.proto` files).
-   **Generate Go code from `.proto` files** using the `protoc` compiler.
-   **Implement a gRPC server** by satisfying the generated service interface.
-   **Implement a client-side streaming RPC**, where the client sends a sequence of messages to the server.
-   **Understand the role of gRPC** as a foundational technology for building robust microservice architectures.

## The Big Picture: Beyond REST/JSON

While REST over HTTP/1.1 is the lingua franca of the public web, it has drawbacks for internal, inter-service communication:
-   **Performance**: HTTP/1.1 is text-based and carries a lot of overhead. JSON parsing can be a CPU bottleneck.
-   **Type Safety**: JSON is flexible, but that flexibility can lead to bugs. There's no enforced contract between the client and server.
-   **Streaming**: True bidirectional streaming is not natively supported over HTTP/1.1.

**gRPC** was designed to solve these problems. It uses **HTTP/2** for its transport, allowing for multiplexed, bidirectional streaming. It uses **Protocol Buffers** for its data format, which is a binary, highly-efficient serialization format. The result is a system that is significantly faster and more robust for building the "nervous system" of a microservices architecture.

## First Principles: IDL and RPC

1.  **RPC (Remote Procedure Call)**: The core idea is to make a function call to another machine seem like a local function call. The framework (gRPC) abstracts away the networking, serialization, and deserialization.
2.  **IDL (Interface Definition Language)**: The practice of defining the API for a service in a language-neutral format. This definition is a "contract." Both the client and the server use this contract to generate code in their respective languages, guaranteeing they are compatible.

### Protocol Buffers (Protobuf)
Protobuf is Google's implementation of an IDL.
-   **Schema-First**: You *must* define your data structures (`messages`) and service endpoints (`rpc`s) in a `.proto` file before you write any code.
-   **Binary Format**: Protobuf encodes data into a compact binary format. This is much smaller and faster to parse than text-based formats like JSON or XML.
-   **Strictly Typed**: The schema ensures that a client and server can't get out of sync. If you change a field's type, you must update the contract, and regenerating the code will cause a compile-time error if there's a mismatch—catching bugs before they hit production.

## Project Structure

```
.
├── cmd/
│   └── dev/
│       └── main.go       # Runs a client and server to demonstrate the service.
├── internal/
│   ├── client/           # The gRPC client implementation.
│   └── server/           # The gRPC server implementation.
└── proto/
    └── telemetry.proto # The Protocol Buffers definition for our service.
```
-   **`proto/`**: The `.proto` file is the single source of truth for the API contract.
-   The generated Go code (`.pb.go`) will also live in the `proto` directory.

## Key Concepts in This Project

### The gRPC Workflow

```
+-------------------+   +-----------------+   +---------------------+   +-------------------+
|  telemetry.proto  |-->|     protoc      |-->|  telemetry.pb.go    |-->| Your Server Code  |
| (Define Service)  |   | (Code Generator)|   | (Generated Go Code) |   | (Implements       |
+-------------------+   +-----------------+   +---------------------+   |  Server Interface)|
                                                    |               +-------------------+
                                                    |
                                                    v
                                              +-------------------+
                                              | Your Client Code  |
                                              | (Uses Client Stub)|
                                              +-------------------+
```

### Defining the Service with `.proto`
The `.proto` file is where it all begins.

```protobuf
syntax = "proto3";
package proto;
option go_package = "./proto";

// The service definition.
service TelemetryService {
  // A client-streaming RPC.
  rpc RecordTelemetry(stream TelemetryData) returns (TelemetrySummary);
}

// A message for a single data point.
message TelemetryData {
  double value = 1;
  // ... other fields
}

// A message for the final summary.
message TelemetrySummary {
  int32 count = 1;
  double sum = 2;
}
```

### Code Generation with `protoc`
You use the `protoc` compiler with Go plugins to generate the necessary code.

```bash
# This command reads telemetry.proto and generates two files.
protoc --go_out=. --go-grpc_out=. proto/telemetry.proto
```
-   `telemetry.pb.go`: Contains the Go structs for your `message` types (`TelemetryData`, `TelemetrySummary`).
-   `telemetry_grpc.pb.go`: Contains the interfaces for the client and server.

### Implementing the Server
The generated code includes a server interface. Your job is to create a struct and implement the methods of that interface.

```go
// Our server struct will embed the generated unimplemented server
// for forward compatibility.
type TelemetryServer struct {
    pb.UnimplementedTelemetryServiceServer
}

// We implement the RecordTelemetry method.
func (s *TelemetryServer) RecordTelemetry(stream pb.TelemetryService_RecordTelemetryServer) error {
    // ... our logic here ...
}
```

### Client-Side Streaming
This is the core of this project. The client sends a *stream* of messages, not just one. The server's handler receives this stream as an object it can read from in a loop.

```go
func (s *TelemetryServer) RecordTelemetry(stream pb.TelemetryService_RecordTelemetryServer) error {
    var count int32
    var sum float64

    // Loop until the client closes the stream.
    for {
        // Receive the next message from the stream.
        data, err := stream.Recv()

        if err == io.EOF {
            // The client has finished sending data.
            // Send back the summary and close the connection.
            return stream.SendAndClose(&pb.TelemetrySummary{
                Count: count,
                Sum:   sum,
            })
        }
        if err != nil {
            return err // Handle other errors.
        }

        // Process the received data point.
        count++
        sum += data.GetValue()
    }
}
```

## Progression: The Microservice Backbone

This project is a leap into the world of modern, high-performance backend engineering. You are moving from the familiar world of REST/JSON to the technology that powers communication at companies like Google, Netflix, and Square. Understanding gRPC is essential for building fast, scalable, and robust microservice systems. The skills learned here are directly applicable to creating internal APIs that are both high-performance and easy to maintain.

## How to Run and Test

1.  **Generate the Go code from the `.proto` file:**
    You must have `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc` installed.
    ```bash
    go generate ./...
    ```

2.  **Run the dev harness:**
    This command will start a gRPC server and a client that connects to it, streams a few data points, and prints the summary.
    ```bash
    go run ./cmd/dev/main.go
    ```

3.  **Run the tests:**
    ```bash
    go test -v ./...
    ```

## Key Takeaways

-   **gRPC is a high-performance RPC framework** ideal for inter-service communication.
-   **Protocol Buffers (`.proto`) provide a contract-first, strongly-typed** way to define APIs.
-   The `protoc` tool **generates client and server code**, saving you from writing boilerplate.
-   gRPC has first-class support for **streaming**, enabling more complex communication patterns than simple request-response.
-   In a client stream, the server uses a `for` loop and `stream.Recv()` to process messages until `io.EOF`.

## Further Reading

-   [**gRPC Official Website**](https://grpc.io/)
-   [**Protocol Buffers Documentation**](https://protobuf.dev/overview/)
-   [**gRPC Go Quickstart**](https://grpc.io/docs/languages/go/quickstart/)
-   [**gRPC Concepts**](https://grpc.io/docs/what-is-grpc/core-concepts/)