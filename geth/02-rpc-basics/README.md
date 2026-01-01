# 02: RPC Basics

## What Is This Project About?

This module deepens your understanding of how Ethereum clients communicate via JSON-RPC. You'll learn the structure of RPC requests and responses, common RPC methods, and how the go-ethereum library abstracts these calls. Understanding the RPC layer is crucial for debugging, optimizing, and building custom tooling.

## Why Is This Important?

While libraries like ethclient abstract RPC calls, understanding the underlying protocol helps you:
- Debug connection and response issues
- Build custom RPC methods for specialized nodes
- Optimize batch requests for performance
- Understand rate limiting and error handling

## Real-World Problems This Solves

- **Custom node integration**: Work with nodes that have non-standard RPC methods
- **Performance optimization**: Batch multiple RPC calls to reduce latency
- **Error debugging**: Understand RPC error codes and responses
- **Protocol compliance**: Ensure your tools work with any Ethereum-compatible node

## Key Concepts You'll Learn

- **JSON-RPC 2.0 protocol**: Request/response format
- **Common RPC methods**: eth_blockNumber, eth_getBlockByNumber, etc.
- **Error handling**: Understanding RPC error responses
- **Batch requests**: Sending multiple requests in one call

## Prerequisites

- Completion of `geth/01-stack`
- Understanding of JSON format

## Project Structure

```
geth/02-rpc-basics/
├── cmd/
│   ├── app/
│   │   └── main.go
│   └── dev/
│       └── main.go
├── internal/
│   └── rpcbasics/
│       ├── exercise.go
│       ├── exercise_test.go
│       ├── solution.reference.go
│       └── solution_no_err.reference.go
└── .vscode/
    └── launch.json
```

## How to Run

```bash
# Using CLI
go run ./cmd/app/main.go https://eth.llamarpc.com

# Using debug harness
go run ./cmd/dev/main.go
```

## Testing

```bash
go test -v ./...
```

## Additional Resources

- [Ethereum JSON-RPC Specification](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
