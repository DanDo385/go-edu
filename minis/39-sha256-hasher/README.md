# minis/39-sha256-hasher

## Problem

Problem: Implementing SHA-256 cryptographic hash operations

We need to:
1. Hash strings and return hex-encoded digests
2. Hash files efficiently (streaming, not loading into memory)
3. Verify file integrity using checksums
4. Demonstrate incremental hashing (streaming multiple inputs)
5. Compare hashes securely (constant-time to prevent timing attacks)

Key Concepts:
- Cryptographic hash: One-way function mapping arbitrary data to fixed-size digest
- SHA-256: 256-bit (32-byte) secure hash algorithm
- Hex encoding: Binary → human-readable (64 hex characters for SHA-256)
- Incremental hashing: Process data in chunks (essential for large files)
- Constant-time comparison: Prevent timing attacks in security-critical code

Time/Space Complexity:
- HashString: O(n) time, O(1) space (where n = string length)
- HashFile: O(n) time, O(1) space (where n = file size, streams in chunks)
- VerifyFile: O(n) time, O(1) space (same as HashFile)
- HashIncremental: O(n) time, O(1) space (where n = total length of all parts)
- CompareHashes: O(1) time (hashes are fixed size), O(1) space

Why Go is well-suited:
- crypto/sha256: Battle-tested, optimized implementation
- io.Copy: Efficient streaming without manual buffer management
- hash.Hash interface: Clean abstraction for incremental hashing
- crypto/subtle: Security primitives (constant-time operations)

Real-world applications:
- File integrity: Verify downloads haven't been corrupted or tampered
- Deduplication: Identify identical files by content (not name)
- Content-addressable storage: Git, IPFS use hashes as identifiers
- Password verification: Store Hash(password), verify by comparing hashes
- Digital signatures: Hash documents before signing (efficiency)

## Quickstart

```bash
cd minis/39-sha256-hasher
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/app -fn "HashFile" -in "Hello, 世界 👋"
go run ./cmd/app -fn "HashString" -in "Hello, 世界 👋"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/sha256hasher/exercise.go`: implement the TODOs here
- `internal/sha256hasher/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
