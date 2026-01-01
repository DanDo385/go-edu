# 09: Events - Contract Event Logs

## What Is This Project About?

This module teaches how to query and decode ERC20 Transfer events from blockchain logs. Events are append-only records emitted during contract execution. They're searchable via bloom filters and provide an efficient way to track state changes without querying contract storage.

This builds on `geth/06-smart-contracts` (where you learned about events in the console), `geth/07-eth-call` (manual encoding), and `geth/08-abigen` (typed bindings), showing you how to work with event logs programmatically.

## Why Is This Important?

Understanding events is crucial because:

- **Efficient indexing**: Events are designed for searching and filtering
- **Historical queries**: Can query past events without replaying transactions
- **Decentralized applications**: Events are the primary way contracts communicate
- **Analytics**: Most analytics tools rely on event logs

## Real-World Problems This Solves

- **Token transfer tracking**: Monitor all transfers of a token
- **Contract state changes**: Track when contract state changes
- **Analytics dashboards**: Build dashboards from event data
- **Indexers**: Build custom indexes from event logs

## Key Concepts You'll Learn

- **Event Structure**: Topics (indexed) vs Data (non-indexed)
- **Topic Hashes**: keccak256 of event signature
- **Bloom Filters**: Probabilistic data structures for efficient searching
- **Log Decoding**: Extracting event parameters from logs
- **Block Range Queries**: Querying events across multiple blocks

## Prerequisites

- Completion of `geth/06-smart-contracts` (understanding events in console)
- Completion of `geth/07-eth-call` (understanding ABI encoding)
- Completion of `geth/08-abigen` (understanding typed bindings)
- Understanding of hash functions and bloom filters

## Project Structure

```
09-events/
├── cmd/
│   ├── app/          # Application entry point (CLI arguments)
│   └── dev/          # Debug harness (fixed inputs)
├── internal/
│   └── events/       # Exercise implementation
│       ├── exercise.go
│       ├── exercise_test.go
│       ├── solution.reference.go
│       ├── solution_no_err.reference.go
│       └── types.go
└── .vscode/
    └── launch.json   # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)

```bash
# Query Transfer events
go run ./cmd/app/main.go <RPC_URL> <contract_address> <event_signature>

# Example:
go run ./cmd/app/main.go https://eth.llamarpc.com 0x... Transfer(address,address,uint256)
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments
2. Use VS Code debugger (F5) and select:
   - **"Debug: cmd/app"** - Debug with CLI arguments
   - **"Debug: cmd/dev"** - Debug with fixed inputs (recommended)
   - **"Test: Run All Tests"** - Debug tests
3. Step through event querying and decoding
4. Inspect log structures in Variables panel

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestFunctionName ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

## Exercises

Implement the `Run` function in `internal/events/exercise.go`:

1. **Calculate Event Signature Hash**: keccak256 of event signature
2. **Build Query**: Create FilterQuery with contract address and topics
3. **Query Logs**: Use FilterLogs to get matching logs
4. **Decode Events**: Extract Transfer parameters (from, to, value) from logs
5. **Handle Multiple Events**: Process all Transfer events in the range

## Connection to Previous Modules

- **geth/06-smart-contracts**: Console tutorial on events and logs (decoding Transfer events)
- **geth/07-eth-call**: Manual ABI encoding (similar decoding logic for events)
- **geth/08-abigen**: Typed bindings (can also be used for event decoding)

## Where This Goes Next

- **geth/10-filters**: Advanced event filtering with multiple topics
- **geth/15-receipts**: Understanding transaction receipts and their logs

## Additional Resources

- [Ethereum Events Documentation](https://ethereum.org/en/developers/docs/smart-contracts/anatomy/#events-and-logs)
- [Bloom Filters Explained](https://en.wikipedia.org/wiki/Bloom_filter)
- [go-ethereum FilterLogs Documentation](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient#Client.FilterLogs)
