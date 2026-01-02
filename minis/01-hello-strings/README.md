# 01-hello-strings

**UTF-8 String Utilities**

Implement UTF-8-aware string utilities that correctly handle multi-byte characters.

## What You'll Learn

- UTF-8 encoding in Go
- Runes vs bytes
- strings and unicode packages
- Proper string manipulation

## Functions to Implement

| Function | Description |
|----------|-------------|
| `TitleCase(s string) string` | Capitalize first letter of each word |
| `Reverse(s string) string` | Reverse string character-by-character |
| `RuneLen(s string) int` | Count characters (runes), not bytes |

## Project Structure

```
01-hello-strings/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/hellostrings/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/01-hello-strings

# Title case a string
go run ./cmd/app/main.go titlecase "hello world"

# Reverse a string
go run ./cmd/app/main.go reverse "hello 世界"

# Count runes
go run ./cmd/app/main.go runelen "hello 👋"
```

### Run the Debug Harness

```bash
# Auto-runs all functions with sample inputs
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Command | Argument | Description |
|---------|----------|-------------|
| `titlecase` | `TEXT` | Text to title-case |
| `reverse` | `TEXT` | Text to reverse |
| `runelen` | `TEXT` | Text to count runes |

## Quick Copy & Paste

```bash
# Title case
go run ./cmd/app/main.go titlecase "hello world"

# Reverse with emoji
go run ./cmd/app/main.go reverse "Hi👋"

# Count runes in UTF-8 string
go run ./cmd/app/main.go runelen "café 世界"

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Runes**: int32 representing Unicode code points
2. **UTF-8**: Variable-width encoding (1-4 bytes per rune)
3. **len(s)**: Returns bytes, not characters!
4. **utf8.RuneCountInString**: Returns character count

## Next Steps

After completing this exercise, proceed to `minis/02-arrays-maps-basics`.
