# 01-hello-strings: UTF-8 String Operations

## Overview

Learn UTF-8-aware string operations in Go including title case conversion, string reversal, and rune counting. Foundation for text processing.

## Learning Objectives

- Understand UTF-8 encoding in Go strings
- Work with runes vs bytes
- Handle multi-byte characters (emoji, CJK, accented letters)
- Implement safe string manipulation

## Quick Start

```bash
# Run tests
go test -v ./...

# Run CLI application
go run ./cmd/app/main.go "hello world"

# Debug with fixed inputs
go run ./cmd/dev/main.go
```

## CLI Arguments

```bash
# Title case a string
go run ./cmd/app/main.go titlecase "hello world"
# Output: Hello World

# Reverse a string (UTF-8 aware)
go run ./cmd/app/main.go reverse "Hello 👋 World"
# Output: dlroW 👋 olleH

# Count runes (not bytes)
go run ./cmd/app/main.go runelen "Hello 世界"
# Output: 8
```

## What the Dev Harness Demonstrates

1. **TitleCase** - Converts first letter of each word to uppercase
2. **Reverse** - Reverses string preserving multi-byte characters
3. **RuneLen** - Counts Unicode code points, not bytes
4. **UTF-8 Handling** - Works correctly with emoji and international text

## Key Concepts

### Strings are Byte Slices
- Go strings are immutable byte slices
- May contain multi-byte UTF-8 sequences
- Direct indexing gives bytes, not characters

### Runes
- `rune` is an alias for `int32`
- Represents a Unicode code point
- Use `[]rune()` to convert string to runes for character-level operations

## Next Steps

Proceed to **minis/02-arrays-maps-basics**
