# 03: Keys and Addresses

## What Is This Project About?

This module teaches you the cryptographic foundations of Ethereum accounts. You'll learn how private keys are generated, how public keys are derived using elliptic curve cryptography, and how Ethereum addresses are computed from public keys. Understanding this flow is essential for secure key management and wallet development.

## Why Is This Important?

Key management is the most security-critical aspect of Ethereum development. Understanding:
- How keys are generated and why randomness matters
- The mathematical relationship between private and public keys
- How addresses are derived and why they're shorter than public keys
- Best practices for key storage and handling

...protects users' funds and prevents catastrophic security failures.

## Real-World Problems This Solves

- **Wallet development**: Generate and manage keys securely
- **Address verification**: Validate address checksums and formats
- **HD wallets**: Understand the foundation for hierarchical deterministic wallets
- **Security auditing**: Identify key management vulnerabilities

## Key Concepts You'll Learn

- **Private keys**: 256-bit random numbers that control accounts
- **secp256k1**: The elliptic curve used by Ethereum (and Bitcoin)
- **Public key derivation**: One-way function from private to public key
- **Keccak-256 hashing**: The hash function used for address derivation
- **Checksum encoding**: EIP-55 mixed-case checksum for addresses

## Prerequisites

- Completion of `geth/01-stack` and `geth/02-rpc-basics`
- Basic understanding of cryptographic concepts

## Project Structure

```
geth/03-keys-addresses/
├── cmd/
│   ├── app/
│   │   └── main.go
│   └── dev/
│       └── main.go
├── internal/
│   └── keysaddresses/
│       ├── exercise.go
│       ├── exercise_test.go
│       ├── solution.reference.go
│       └── solution_no_err.reference.go
└── .vscode/
    └── launch.json
```

## How to Run

```bash
# Generate new key pair
go run ./cmd/app/main.go

# Derive address from existing private key
go run ./cmd/app/main.go 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
```

## Testing

```bash
go test -v ./...
```

## Additional Resources

- [Mastering Ethereum - Keys and Addresses](https://github.com/ethereumbook/ethereumbook/blob/develop/04keys-addresses.asciidoc)
- [EIP-55: Mixed-case checksum address encoding](https://eips.ethereum.org/EIPS/eip-55)
