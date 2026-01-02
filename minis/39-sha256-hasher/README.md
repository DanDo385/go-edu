# 39-sha256-hasher

**SHA256 Hashing**

Implement cryptographic hashing utilities.

## What You'll Learn

- crypto/sha256 package
- Hash.Write pattern
- Hex encoding
- File hashing

## Functions to Implement

| Function | Description |
|----------|-------------|
| Hash data with SHA256 | String and file hashing |

## Project Structure

```
39-sha256-hasher/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/sha256hasher/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/39-sha256-hasher

# Hash a string
go run ./cmd/app/main.go hash "hello world"

# Hash a file
go run ./cmd/app/main.go file /path/to/file

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Hash a string
go run ./cmd/app/main.go hash "hello world"

# Hash a file
go run ./cmd/app/main.go file README.md

# Verify with sha256sum
echo -n "hello world" | sha256sum

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **sha256.Sum256**: Direct hash
2. **sha256.New**: Streaming hash
3. **hex.EncodeToString**: Human-readable output
4. **io.Copy(hash, reader)**: Efficient file hashing

## Next Steps

After completing this exercise, proceed to `minis/40-merkle-tree-basics`.
