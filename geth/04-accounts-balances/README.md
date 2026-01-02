# Ethereum: 04 Accounts Balances

## Project Structure

```
04-accounts-balances/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with auto-demo
├── internal/
│   └── accountsbalances/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       └── solution.reference.go  # Reference implementation
└── README.md
```

## Getting Started

1. Navigate to this directory:
   ```bash
   cd geth/04-accounts-balances
   ```

2. Open the exercise file:
   ```bash
   code internal/accountsbalances/exercise.go
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
go run ./cmd/app/main.go <RPC_URL> <ADDRESS>
```

**Examples:**

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
```

```bash
*/
```

**Copy & Paste Examples:**

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
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

If you get stuck, check `internal/accountsbalances/solution.reference.go` for a complete
implementation with detailed explanations.

