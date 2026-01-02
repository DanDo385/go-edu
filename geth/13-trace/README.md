# 13: Transaction Tracing

## What Is This Project About?

This module teaches you how to trace transaction execution to understand internal calls, state changes, and debug contract interactions. Tracing is an advanced debugging technique that shows exactly what happened during transaction execution.

## Why Is This Important?

Transaction tracing enables:
- Debugging failed transactions
- Understanding internal calls
- Analyzing gas usage
- MEV research and analysis

## Key Concepts You'll Learn

- **Call traces**: Internal message calls during execution
- **State diffs**: Storage changes made by transactions
- **Gas tracing**: Per-opcode gas consumption
- **debug_traceTransaction**: RPC method for tracing

## Prerequisites

- Completion of `geth/01-stack` through `geth/12-proofs`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
