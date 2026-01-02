# 12-proofs

**Merkle-Patricia Proofs**

Fetch Merkle-Patricia trie proofs for accounts and storage slots.

## What You'll Learn

- Merkle-Patricia trie structure
- Account proofs (balance, nonce, code hash, storage root)
- Storage proofs (specific slot values)
- Verifying proofs against state roots

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Fetch and inspect account/storage proofs |

## Project Structure

```
12-proofs/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/proofs/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/12-proofs

# Get account proof
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7
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
| `ADDRESS` | Yes | Account to get proof for |
| `STORAGE_KEYS` | No | Comma-separated storage keys |

## Quick Copy & Paste

```bash
# Get account proof
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Account Proof**: Path through state trie to account
2. **Storage Proof**: Path through storage trie to slot
3. **eth_getProof**: Combined account and storage proofs

## Next Steps

After completing this exercise, proceed to `geth/13-trace`.
