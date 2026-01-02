# 03-keys-addresses: Key Management and Address Derivation

## Overview

Learn how to generate private keys, derive public keys and Ethereum addresses, and securely store keys using the keystore format. This is fundamental to understanding Ethereum account management.

## Learning Objectives

- Generate cryptographically secure private keys using secp256k1
- Derive Ethereum addresses from public keys
- Store keys securely using the keystore format with encryption
- Understand the relationship between private keys, public keys, and addresses

## Project Structure

```
03-keys-addresses/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── keysaddresses/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go        # Complete solution
│       └── solution_no_err.reference.go # Error-free variant
└── README.md               # This file
```

## Quick Start

### 1. Implement the Exercise

Open `internal/keysaddresses/exercise.go` and implement the `Run` function.

### 2. Run Tests

```bash
go test -v ./...
```

### 3. Test with CLI

```bash
# Generate a key with default settings
go run ./cmd/app/main.go

# Generate a key with custom output directory
go run ./cmd/app/main.go --output ./my-keys

# Generate a key with custom password
go run ./cmd/app/main.go --password "mySecurePassword123"

# Specify both
go run ./cmd/app/main.go --output ./my-keys --password "mySecurePassword123"
```

### 4. Debug with Dev Harness

```bash
go run ./cmd/dev/main.go
```

## CLI Arguments (cmd/app/main.go)

### Syntax

```bash
go run ./cmd/app/main.go [OPTIONS]
```

### Options

- `--output <dir>` - Output directory for keystore (default: "./keystore-demo")
- `--password <pass>` - Password to encrypt the keystore (default: "changeit")

### Example Commands

```bash
# Use defaults
go run ./cmd/app/main.go

# Custom output directory
go run ./cmd/app/main.go --output ./prod-keys

# Custom password
go run ./cmd/app/main.go --password "MySecurePass123!"

# Both custom options
go run ./cmd/app/main.go --output ./my-wallet --password "SecurePass456"
```

## What the Dev Harness Demonstrates

The `cmd/dev/main.go` automatically demonstrates:

1. **Key Generation** - Creates a new secp256k1 private key
2. **Address Derivation** - Derives Ethereum address from public key
3. **Keystore Creation** - Encrypts and stores the key
4. **Key Recovery** - Unlocks the keystore to prove it works
5. **Raw Key Export** - Shows the private key in hex format

## Key Concepts

### secp256k1 Elliptic Curve

Ethereum uses the secp256k1 elliptic curve for key generation:
- Private key: 256-bit random number
- Public key: Derived from private key using elliptic curve multiplication
- Address: Last 20 bytes of Keccak-256 hash of public key

### Keystore Format

The keystore format (as used by Geth):
- Uses scrypt for key derivation from password
- Encrypts the private key with AES-128-CTR
- Stores the encrypted key in a JSON file
- Includes parameters for decryption and verification

### Security Considerations

- **Never share private keys**: Anyone with your private key can control your funds
- **Use strong passwords**: The keystore is only as secure as your password
- **Backup carefully**: Loss of private key means permanent loss of access
- **Test with testnets**: Always test with testnet tokens first

## Next Steps

After completing this exercise, proceed to:
- **geth/04-accounts-balances** - Query account balances and transactions

## Resources

- [Ethereum Keys and Addresses](https://ethereum.org/en/developers/docs/accounts/)
- [Keystore Format Specification](https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition)
- [secp256k1 Curve](https://en.bitcoin.it/wiki/Secp256k1)
