# 09: Contract Events and Logs

## What Is This Project About?

Subscribe to and decode contract events. This module teaches you to listen for on-chain events, filter logs, and decode event data—building on console concepts from geth/06-smart-contracts.

## Why Is This Important?

Events are how contracts communicate state changes. Indexers, analytics tools, and real-time UIs all rely on event subscriptions to stay up-to-date with on-chain activity.

## Real-World Problems This Solves

- **Building real-time UIs that update when on-chain data changes**
- **Creating indexers that track contract state over time**
- **Monitoring specific addresses or events for analytics**

## Key Concepts You'll Learn

- **Event logs and topics (indexed vs non-indexed parameters)**: Event logs and topics (indexed vs non-indexed parameters)
- **Log filtering by address, topics, and block range**: Log filtering by address, topics, and block range
- **Event subscriptions (eth_subscribe for WebSocket)**: Event subscriptions (eth_subscribe for WebSocket)
- **Decoding event data from logs**: Decoding event data from logs

## Prerequisites

Completion of geth/06-smart-contracts, geth/07-eth-call, and geth/08-abigen

## Project Structure

```
geth/09-events/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── events/
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
go run ./cmd/app/main.go <RPC_URL> <contract_address> <event_signature>
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

See `internal/events/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Ethereum Development Documentation](https://ethereum.org/en/developers/)

## Next Steps

- **geth/10-filters**
- **geth/17-indexer**
