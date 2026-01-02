# 10: Event Filters

## What Is This Project About?

This module teaches you how to create and use event filters for querying historical events and subscribing to new events in real-time. Filters allow you to efficiently query large ranges of blocks for specific events.

## Why Is This Important?

Filters are essential for:
- Efficiently querying event history
- Real-time event subscriptions
- Building responsive DApps
- Indexer performance optimization

## Key Concepts You'll Learn

- **FilterQuery**: Specifying filter criteria
- **Block ranges**: FromBlock and ToBlock parameters
- **Topic filtering**: Matching events by indexed parameters
- **Address filtering**: Matching events by emitting contract

## Prerequisites

- Completion of `geth/09-events`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
