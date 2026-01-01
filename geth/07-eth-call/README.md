# 07: eth_call - Manual Contract Calls

## What Is This Project About?

This module teaches how to interact with smart contracts using manual ABI encoding and decoding. You'll learn to query ERC20 token metadata (name, symbol, decimals, totalSupply) by manually encoding function selectors and decoding return values, giving you a deep understanding of how contract calls work at the ABI level.

This builds directly on `geth/06-smart-contracts`, where you learned contract interaction concepts via the Geth console. Here, you'll implement the same concepts in Go using `go-ethereum/ethclient`.

## Why Is This Important?

Understanding manual ABI encoding/decoding is crucial because:

- **Deep understanding**: See exactly how contract calls are encoded as bytes
- **Foundation for libraries**: Understand what abigen and other tools do under the hood
- **Custom implementations**: Sometimes you need manual control for performance or special cases
- **Debugging**: When typed bindings fail, manual encoding helps debug issues

## Real-World Problems This Solves

- **Querying contracts without bindings**: When you don't have generated contract bindings
- **Custom encoding logic**: For MEV bots and other performance-critical applications
- **Understanding RPC calls**: See what `eth_call` actually does at the byte level
- **Debugging contract interactions**: When automatic decoding fails, manual decoding helps

## Key Concepts You'll Learn

- **Function Selectors**: First 4 bytes of keccak256(functionSignature)
- **ABI Encoding**: How function calls are encoded as bytes
- **eth_call**: Simulating contract execution without sending transactions
- **Manual Decoding**: Decoding dynamic types (strings) and static types (uint256) from raw bytes
- **Call vs Transaction**: Understanding read-only contract calls

## Prerequisites

- Completion of `geth/06-smart-contracts` (understanding contract interaction concepts)
- Completion of `geth/01-stack` through `geth/05-tx-nonces`
- Understanding of basic Ethereum concepts (addresses, transactions, contracts)
- Basic knowledge of ABI encoding (covered in geth/06-smart-contracts)

## Project Structure

```
07-eth-call/
├── cmd/
│   ├── app/          # Application entry point (CLI arguments)
│   └── dev/          # Debug harness (fixed inputs)
├── internal/
│   └── ethcall/      # Exercise implementation
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
# Query USDC contract on mainnet
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

# Query at specific block
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 12345
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs (USDC contract)
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments
2. Use VS Code debugger (F5) and select:
   - **"Debug: cmd/app"** - Debug with CLI arguments
   - **"Debug: cmd/dev"** - Debug with fixed inputs (recommended)
   - **"Test: Run All Tests"** - Debug tests
3. Step through ABI encoding/decoding logic
4. Inspect raw bytes in Variables panel

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

Implement the `Run` function in `internal/ethcall/exercise.go`:

1. **Create Helper Function**: Build a closure for making contract calls
2. **Call name()**: Encode selector and decode string return value
3. **Call symbol()**: Encode selector and decode string return value
4. **Call decimals()**: Encode selector and decode uint8 return value
5. **Call totalSupply()**: Encode selector and decode uint256 return value

## Connection to Previous Modules

- **geth/06-smart-contracts**: Console-based tutorial on contract interaction concepts
- **geth/01-stack**: RPC connection patterns
- **geth/05-tx-nonces**: Transaction building (this module focuses on calls, not transactions)

## Where This Goes Next

- **geth/08-abigen**: Using code generation for typed contract bindings
- **geth/09-events**: Listening to contract events and logs

## Additional Resources

- [ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [ERC20 Token Standard](https://eips.ethereum.org/EIPS/eip-20)
- [go-ethereum CallContract Documentation](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient#Client.CallContract)
