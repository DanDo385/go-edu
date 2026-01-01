# 02: RPC Basics and Retry Logic

## What Is This Project About?

Learn to make various RPC calls and implement retry logic for fault tolerance. This module builds on the connectivity concepts from module 01 by introducing multiple RPC methods and resilience patterns.

## Why Is This Important?

Production Ethereum applications must handle transient network failures gracefully. This module teaches you to build robust RPC clients that retry failed requests and handle various error conditions.

## Real-World Problems This Solves

- **Handling rate limits from public RPC providers**
- **Building resilient applications that survive transient network failures**
- **Implementing exponential backoff for distributed systems**

## Key Concepts You'll Learn

- **Multiple RPC methods**:  NetworkID, BlockNumber, BlockByNumber
- **Retry logic with exponential backoff**: Retry logic with exponential backoff
- **Context-aware retry loops (respecting cancellation)**: Context-aware retry loops (respecting cancellation)
- **Fault tolerance patterns in distributed systems**: Fault tolerance patterns in distributed systems

## Prerequisites

Completion of geth/01-stack

## Project Structure

```
geth/02-rpc-basics/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── rpcbasics/
│       ├── exercise.go      # Your implementation
│       ├── exercise_test.go # Tests
│       ├── solution.reference.go      # Reference solution
│       ├── solution_no_err.reference.go # Simplified reference
│       └── types.go         # Type definitions
└── .vscode/
    └── launch.json          # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)

```bash
go run ./cmd/app/main.go <RPC_URL>
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5)
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments
2. Press F5 in VS Code, select "Debug: cmd/dev (Debug Harness)"
3. Step through with F10 (Step Over) and F11 (Step Into)
4. Watch variables in the Variables panel

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

## Exercises

See `internal/rpcbasics/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Ethereum Development Documentation](https://ethereum.org/en/developers/)

## Next Steps

- **geth/03-keys-addresses**
- **geth/04-accounts-balances**
