# 41-signed-transactions-ed25519

**Ed25519 Signed Transactions**

Build and sign transactions with Ed25519.

## What You'll Learn

- Ed25519 signatures
- Transaction structure
- Signing and verification
- Key generation

## Functions to Implement

| Function | Description |
|----------|-------------|
| Sign transactions | Ed25519 signing |
| Verify signatures | Validate transaction |

## Project Structure

```
41-signed-transactions-ed25519/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/signedtransactionsed25519/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/41-signed-transactions-ed25519

# Generate keypair
go run ./cmd/app/main.go keygen

# Sign transaction
go run ./cmd/app/main.go sign --key private.key --to "alice" --amount 100

# Verify signature
go run ./cmd/app/main.go verify --tx transaction.json

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Generate keys
go run ./cmd/app/main.go keygen

# Sign a transaction
go run ./cmd/app/main.go sign --to alice --amount 100

# Verify
go run ./cmd/app/main.go verify --tx tx.json

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Ed25519**: Fast, secure signatures
2. **crypto/ed25519**: Standard library
3. **Public/Private Keys**: 32/64 bytes
4. **Signature**: 64 bytes

## Next Steps

After completing this exercise, proceed to `minis/42-simple-block-struct-hashing`.
