# 24: Chain Monitor

## What Is This Project About?

This module teaches you how to build a real-time chain monitor that watches for new blocks and transactions. Monitors are fundamental for building responsive blockchain applications.

## Why Is This Important?

Chain monitoring enables:
- Real-time notifications
- Transaction tracking
- Block streaming
- Event-driven architectures

## Key Concepts You'll Learn

- **Block subscriptions**: WebSocket new block events
- **Transaction watching**: Tracking specific txs
- **Head tracking**: Following the chain tip
- **Subscription management**: Handling reconnects

## Prerequisites

- Completion of `geth/23-mempool`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
