# 17-file-streaming-bufio

**File Streaming with bufio**

Efficiently read large files using buffered I/O.

## What You'll Learn

- bufio.Scanner for line reading
- bufio.Reader for buffered reading
- Memory-efficient file processing
- Custom split functions

## Functions to Implement

| Function | Description |
|----------|-------------|
| Stream and process file | Line-by-line processing |

## Project Structure

```
17-file-streaming-bufio/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/filestreamingbufio/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/17-file-streaming-bufio

# Process file
go run ./cmd/app/main.go input.txt

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Process file
go run ./cmd/app/main.go input.txt

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **bufio.Scanner**: Line-by-line reading
2. **bufio.Reader**: Buffered byte reading
3. **O(1) Memory**: Process files larger than RAM
4. **Custom Split**: ScanLines, ScanWords, custom

## Next Steps

After completing this exercise, proceed to `minis/18-goroutines-1M-demo`.
