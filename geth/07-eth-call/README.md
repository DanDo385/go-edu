# Ethereum: 07 Eth Call

## Problem

Query ERC20 token metadata using manual ABI encoding/decoding.

## Project Structure

```
07-eth-call/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with auto-demo
├── internal/
│   └── ethcall/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       └── solution.reference.go  # Reference implementation
└── README.md
```

## Getting Started

1. Navigate to this directory:
   ```bash
   cd geth/07-eth-call
   ```

2. Open the exercise file:
   ```bash
   code internal/ethcall/exercise.go
   ```

3. Implement the TODO functions

4. Run tests:
   ```bash
   go test -v ./...
   ```

## CLI Usage

### Using cmd/app/main.go

The `cmd/app/main.go` file provides a CLI interface with custom arguments.

**Usage:**
```bash
go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>
```

**Examples:**

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F
```

```bash
*/
```

**Copy & Paste Examples:**

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F
```

### Using cmd/dev/main.go

The `cmd/dev/main.go` file automatically demonstrates the project's capabilities
by running through different scenarios with pre-configured inputs.

**Run the demo:**

```bash
go run ./cmd/dev/main.go
```

**Debug with VS Code:**

1. Open `cmd/dev/main.go`
2. Set breakpoints at `// BREAKPOINT:` comments
3. Press F5 and select "Debug cmd/dev (Debug Harness)"

## Testing

Run all tests:

```bash
go test -v ./...
```

Run specific test:

```bash
go test -v -run TestFunctionName ./...
```

## Reference Solution

If you get stuck, check `internal/ethcall/solution.reference.go` for a complete
implementation with detailed explanations.

