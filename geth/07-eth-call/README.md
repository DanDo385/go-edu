# 07: eth_call - Manual Contract Interaction

## What Is This Project About?

This module teaches you how to interact with smart contracts using manual ABI encoding and decoding. You'll learn to construct function selectors, build CallMsg structs, execute eth_call requests, and decode the returned bytes. This gives you deep understanding of what libraries like abigen abstract away.

**Important**: Before starting this module, complete the console tutorial in `geth/06-smart-contracts`. The console experience provides the conceptual foundation that makes this Go implementation much clearer.

## Why Is This Important?

Understanding manual contract calls is essential for:
- Debugging contract interaction issues
- Building tools that work with any contract (not just those with bindings)
- Understanding gas estimation and call simulation
- Implementing custom encoding for non-standard ABIs

## Real-World Problems This Solves

- **Universal contract queries**: Query any contract without pre-generated bindings
- **ABI debugging**: Understand exactly what bytes are being sent/received
- **Custom encoding**: Handle non-standard or dynamic ABIs
- **MEV/arbitrage**: Build low-latency contract interaction tools

## Key Concepts You'll Learn

- **Function selectors**: keccak256(signature)[:4]
- **ABI encoding**: How arguments are serialized to bytes
- **ABI decoding**: How return values are parsed from bytes
- **eth_call**: Read-only contract execution without state changes
- **CallMsg**: Structure describing a contract call

## Prerequisites

- Completion of `geth/06-smart-contracts` (console tutorial)
- Completion of `geth/01-stack` through `geth/06-eip1559`

## Project Structure

```
geth/07-eth-call/
├── cmd/
│   ├── app/
│   │   └── main.go
│   └── dev/
│       └── main.go
├── internal/
│   └── ethcall/
│       ├── exercise.go      # Manual ABI encoding/decoding
│       ├── exercise_test.go
│       ├── solution.reference.go
│       └── types.go
└── .vscode/
    └── launch.json
```

## How to Run

```bash
# Query USDC token metadata
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

# Query DAI token metadata
go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F
```

## Testing

```bash
go test -v ./...
```

## Additional Resources

- [Ethereum ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [go-ethereum CallContract](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient#Client.CallContract)
