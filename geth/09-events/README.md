# 09: Events and Logs

## What Is This Project About?

This module teaches you how to work with Ethereum events (logs). Events are how smart contracts communicate state changes to external observers. You'll learn to query historical events, decode event data, and understand the log structure.

## Why Is This Important?

Events are the primary mechanism for:
- Off-chain indexing of on-chain activity
- Real-time notifications of contract state changes
- Building block explorers and analytics tools
- Debugging contract interactions

## Real-World Problems This Solves

- **Token transfers tracking**: Monitor ERC20/ERC721 Transfer events
- **DeFi monitoring**: Track swaps, liquidity changes, etc.
- **Indexer development**: Build searchable databases of on-chain activity
- **Notification systems**: Alert users of relevant on-chain events

## Key Concepts You'll Learn

- **Log structure**: Topics array and data field
- **Indexed parameters**: Stored in topics for filtering
- **Non-indexed parameters**: Stored in data field
- **Event signatures**: topics[0] = keccak256(EventSignature)
- **Log filtering**: Query events by address, topics, and block range

## Prerequisites

- Completion of `geth/06-smart-contracts` through `geth/08-abigen`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

## Testing

```bash
go test -v ./...
```
