# 03-keys-addresses

**Ethereum Key Management**

Generate, manage, and understand Ethereum keys and addresses. Learn the cryptographic foundations of Ethereum accounts.

## What You'll Learn

- Generating private keys with `crypto.GenerateKey()`
- Deriving public keys and addresses
- Keystore file management
- Cryptographic hashing with Keccak-256

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(cfg)` | Generate keys, derive addresses, manage keystore files |

## Project Structure

```
03-keys-addresses/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/keysaddresses/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/03-keys-addresses

# Generate new key
go run ./cmd/app/main.go generate

# From existing private key
go run ./cmd/app/main.go fromkey <PRIVATE_KEY_HEX>
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
| `command` | Yes | `generate` or `fromkey` |
| `PRIVATE_KEY` | For `fromkey` | Hex-encoded private key |

## Quick Copy & Paste

```bash
# Generate new random key
go run ./cmd/app/main.go generate

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **ECDSA secp256k1**: Ethereum's signature algorithm
2. **Keccak-256**: Hash function for address derivation
3. **Keystore Files**: Encrypted key storage (UTC JSON)
4. **Address Derivation**: public key → Keccak-256 → last 20 bytes

## Next Steps

After completing this exercise, proceed to `geth/04-accounts-balances`.
