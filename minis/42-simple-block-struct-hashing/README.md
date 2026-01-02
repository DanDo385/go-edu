# minis/42-simple-block-struct-hashing

## Problem

Problem: Implement a simple blockchain with blocks, hashing, and validation

Requirements:
1. Block structure with header and transactions
2. Cryptographic hash linking between blocks
3. Deterministic serialization for hashing
4. Merkle root for transaction integrity
5. Chain validation to detect tampering

Data Structure:
- Block: Header + Transactions + Hash
- BlockHeader: Index, Timestamp, PrevHash, MerkleRoot, Nonce
- Chain: Linked list of blocks via hash pointers

Time/Space Complexity:
- NewBlock: O(n) where n = number of transactions (hash each)
- ValidateChain: O(m*n) where m = chain length, n = avg transactions
- Space: O(m*n) for storing entire chain

Why Go is well-suited:
- crypto/sha256: Built-in cryptographic hashing
- encoding/hex: Easy hex encoding for hashes
- bytes.Buffer: Efficient serialization
- encoding/binary: Deterministic binary encoding
- Strong typing: Prevents serialization mistakes

Compared to other languages:
- Python: hashlib for hashing, but slower and less type-safe
- Rust: Similar primitives, but more complex lifetimes
- JavaScript: Web Crypto API, but less suitable for backend
- C++: Fast but error-prone manual memory management

Blockchain concepts:
- Hash chain: Each block hashes previous block's hash
- Tamper evidence: Changing any block breaks chain
- Genesis block: First block with no predecessor
- Merkle root: Commit to all transactions with one hash
- Serialization: Must be deterministic for consensus

## Quickstart

```bash
cd minis/42-simple-block-struct-hashing
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
go run ./cmd/app -fn "NewGenesisBlock"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/simpleblockstructhashing/exercise.go`: implement the TODOs here
- `internal/simpleblockstructhashing/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
