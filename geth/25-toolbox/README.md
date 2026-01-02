# 25-toolbox

**Ethereum Swiss Army Knife**

Build a multi-command CLI that combines multiple node operations.

## What You'll Learn

- Building multi-command CLIs
- Combining multiple operations
- Formatting output for humans
- Command routing patterns

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Route to appropriate handler |
| `handleStatus(ctx, client)` | Get node status |
| `handleBlock(ctx, client, args)` | Get block info |
| `handleTx(ctx, client, args)` | Get transaction info |

## Project Structure

```
25-toolbox/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/toolbox/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/25-toolbox

# Status command
go run ./cmd/app/main.go https://eth.llamarpc.com status

# Block command
go run ./cmd/app/main.go https://eth.llamarpc.com block 18000000

# Transaction command
go run ./cmd/app/main.go https://eth.llamarpc.com tx 0x<hash>
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `RPC_URL` | Yes | Ethereum RPC endpoint URL |
| `COMMAND` | Yes | `status`, `block`, or `tx` |
| `ARGS` | Varies | Command-specific arguments |

## Quick Copy & Paste

```bash
# Get node status
go run ./cmd/app/main.go https://eth.llamarpc.com status

# Get block info
go run ./cmd/app/main.go https://eth.llamarpc.com block 18000000

# Get latest block
go run ./cmd/app/main.go https://eth.llamarpc.com block latest

# Debug harness (runs all commands)
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Command Routing**: Switch on first argument
2. **Subcommands**: status, block, tx
3. **Unified Interface**: One tool, many operations

## Congratulations!

You've completed the geth track! You now have a solid foundation in Ethereum Go development.
