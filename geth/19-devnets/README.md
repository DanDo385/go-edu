# 19: Development Networks

## What Is This Project About?

This module teaches you how to set up and work with local development networks for testing Ethereum applications. Devnets provide instant feedback and free gas for development.

## Why Is This Important?

Development networks enable:
- Fast iteration during development
- Testing without real funds
- Reproducible test scenarios
- Isolated debugging environments

## Key Concepts You'll Learn

- **Geth dev mode**: Local instant-mining chain
- **Account management**: Pre-funded accounts
- **Time manipulation**: Advancing block time
- **State reset**: Starting fresh for tests

## Prerequisites

- Geth installed locally
- Completion of previous modules

## How to Run

```bash
# Start local devnet first
geth --dev --http --http.api eth,net,web3,personal

# Then run the demo
go run ./cmd/app/main.go http://localhost:8545
```

## Testing

```bash
go test -v ./...
```
