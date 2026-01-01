# 08: abigen - Typed Contract Bindings

## What Is This Project About?

This module teaches you how to use abigen-generated typed bindings for smart contract interaction. Instead of manually encoding/decoding ABI as in module 07, you'll use type-safe Go methods generated from contract ABIs. This is the production approach for contract interaction.

## Why Is This Important?

Typed bindings provide:
- Compile-time type checking
- IDE autocompletion
- Reduced boilerplate code
- Fewer encoding/decoding errors

## Real-World Problems This Solves

- **Production contract integration**: Type-safe contract calls
- **Rapid development**: Less boilerplate than manual encoding
- **Error prevention**: Catch type mismatches at compile time
- **Code maintainability**: Clear, readable contract interactions

## Key Concepts You'll Learn

- **abigen tool**: Generating Go bindings from ABI JSON
- **Contract instances**: Creating typed contract handles
- **Method calls**: Type-safe read and write operations
- **Event parsing**: Working with generated event types

## Prerequisites

- Completion of `geth/06-smart-contracts` and `geth/07-eth-call`
- Understanding of manual ABI encoding from module 07

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

## Testing

```bash
go test -v ./...
```

## Additional Resources

- [go-ethereum abigen Documentation](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings)
