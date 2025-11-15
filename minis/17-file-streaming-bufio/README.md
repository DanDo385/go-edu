# Project 17: File Streaming with bufio - Buffered I/O

## What Is This Project About?

This project teaches you how to efficiently read and write files in Go without loading entire files into memory. You'll learn:

1. **File I/O fundamentals** (os.File, io.Reader, io.Writer interfaces)
2. **Buffered I/O with bufio** (Reader, Writer, Scanner)
3. **Streaming large files** (line-by-line processing)
4. **Memory-efficient patterns** (avoiding OOM errors)
5. **Performance trade-offs** (buffering strategies)
6. **Real-world file processing** (logs, CSVs, large datasets)

By the end, you'll understand why buffering matters, how to process files of any size efficiently, and when to use different bufio tools.

---

## The Fundamental Problem: Files Can Be Huge

### Real-World Scenario

Imagine you're processing server logs. A single log file might be:
- **10GB** on a production server
- **Millions of lines** (one per request)
- **Still growing** (actively being written to)

**Bad approach (loads entire file into memory):**
```go
// WRONG: Reads entire 10GB file into memory!
data, err := os.ReadFile("huge.log")  // OOM crash on 4GB machine!
lines := strings.Split(string(data), "\n")
for _, line := range lines {
    process(line)
}
```

**Good approach (streams line-by-line):**
```go
// RIGHT: Reads one line at a time (constant memory usage)
file, _ := os.Open("huge.log")
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    process(line)  // Only one line in memory at a time!
}
```

**Memory usage:**
- Bad approach: **10GB** (entire file)
- Good approach: **~4KB** (one line + buffer)

This is the power of **streaming** with **bufio**.

---

## First Principles: Understanding File I/O in Go

### What Is a File?

A **file** is a sequence of bytes stored on disk. In Go, files are represented by the `os.File` type:

```go
file, err := os.Open("data.txt")  // Open for reading
if err != nil {
    log.Fatal(err)
}
defer file.Close()  // ALWAYS close files!
```

**Key insight:** `os.File` implements the `io.Reader` and `io.Writer` interfaces:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

This means you can read/write files using generic I/O functions!

### The Problem with Direct os.File I/O

When you call `file.Read()`, it makes a **system call** to the operating system:

```
Your Go Program → System Call → OS Kernel → Disk → OS Kernel → Your Program
                  (expensive!)
```

**System calls are slow:**
- Context switch (user space → kernel space)
- Disk I/O (mechanical disks: ~10ms latency, SSDs: ~0.1ms)
- Each `Read()` call typically reads only what you ask for (e.g., 1 byte at a time = disaster!)

**Example of terrible performance:**
```go
// HORRIBLE: 1 million system calls for a 1MB file!
file, _ := os.Open("data.txt")
var b [1]byte
for {
    n, err := file.Read(b[:])
    if err == io.EOF {
        break
    }
    process(b[0])  // Process one byte
}
```

**Solution:** Use **buffering** to minimize system calls.

---

## What Is Buffering?

**Buffering** means reading a **large chunk** of data from disk into memory (the buffer), then serving small reads from that buffer (no system calls needed).

### Analogy: Grocery Shopping

**No buffering** (direct I/O):
- Drive to the store every time you need one ingredient
- Need sugar? Drive to store.
- Need flour? Drive to store.
- Need eggs? Drive to store.
- Result: 100 trips for 100 items!

**With buffering**:
- Drive to store ONCE
- Buy a cart full of groceries (fill the buffer)
- Use items from your pantry as needed
- Result: 1 trip for 100 items!

### Buffered Reading in Go

```go
file, _ := os.Open("data.txt")
bufferedReader := bufio.NewReader(file)

// Now reads come from the buffer (fast!), not disk (slow)
byte, err := bufferedReader.ReadByte()  // No system call (usually)
```

**How it works:**

```
┌─────────────────────────────────┐
│  Disk File (e.g., 1GB)         │
└─────────────────────────────────┘
         ↓ (one system call reads 4KB)
    ┌────────────────┐
    │  Buffer (4KB)  │  ← bufio.Reader
    └────────────────┘
         ↓ (many reads from buffer, no syscalls)
    Your program reads bytes/lines one at a time
```

**Default buffer size:** 4096 bytes (4KB)

---

## bufio.Reader: Buffered Reading

### Creating a Buffered Reader

```go
file, _ := os.Open("input.txt")
defer file.Close()

reader := bufio.NewReader(file)  // Default 4KB buffer
// OR
reader := bufio.NewReaderSize(file, 64*1024)  // Custom 64KB buffer
```

### Reading Methods

#### 1. ReadByte() - Read One Byte
```go
b, err := reader.ReadByte()
if err == io.EOF {
    // End of file
}
```

#### 2. ReadLine() - Read Until '\n' (Low-Level)
```go
line, isPrefix, err := reader.ReadLine()
// isPrefix = true if line was longer than buffer
```

**Gotcha:** `ReadLine()` can return partial lines! Usually avoid it.

#### 3. ReadString(delim) - Read Until Delimiter
```go
line, err := reader.ReadString('\n')  // Includes '\n' at end
if err == io.EOF {
    // Last line might not have '\n'
}
```

**Use case:** Reading line-by-line when you need to keep the delimiter.

#### 4. ReadBytes(delim) - Like ReadString but Returns []byte
```go
line, err := reader.ReadBytes('\n')  // Returns []byte instead of string
```

**Use case:** When you need to modify the line (avoid string → []byte conversion).

---

## bufio.Scanner: The Best Way to Read Lines

`bufio.Scanner` is a **higher-level** abstraction that makes line-by-line reading trivial:

```go
file, _ := os.Open("data.txt")
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()  // Current line as string
    fmt.Println(line)
}

if err := scanner.Err(); err != nil {
    log.Fatal(err)  // Check for errors after scanning
}
```

### Why Use Scanner?

**Advantages:**
- **Simple API**: No need to check `io.EOF` manually
- **Handles edge cases**: Last line without '\n', long lines, etc.
- **Flexible splitting**: Can scan by words, bytes, custom delimiters

**Example: Scan by Words**
```go
scanner := bufio.NewScanner(file)
scanner.Split(bufio.ScanWords)  // Split on whitespace

for scanner.Scan() {
    word := scanner.Text()
    fmt.Println(word)
}
```

### Scanner Split Functions

Built-in split functions:
- `bufio.ScanLines` (default): Split on '\n'
- `bufio.ScanWords`: Split on whitespace
- `bufio.ScanBytes`: One byte at a time
- `bufio.ScanRunes`: One UTF-8 character at a time

**Custom split function:**
```go
// Split on commas (simple CSV)
scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    for i := 0; i < len(data); i++ {
        if data[i] == ',' {
            return i + 1, data[:i], nil
        }
    }
    if atEOF && len(data) > 0 {
        return len(data), data, nil
    }
    return 0, nil, nil
})
```

### Scanner Limitations

**Maximum token size:** 64KB by default (can increase with `scanner.Buffer()`).

If you need to process lines > 64KB:
```go
scanner := bufio.NewScanner(file)
buf := make([]byte, 0, 1024*1024)  // 1MB initial buffer
scanner.Buffer(buf, 10*1024*1024)  // 10MB max token size

for scanner.Scan() {
    line := scanner.Text()  // Can now handle huge lines
}
```

---

## bufio.Writer: Buffered Writing

Writing data one byte at a time is just as slow as reading. `bufio.Writer` buffers writes and flushes to disk in chunks.

### Creating a Buffered Writer

```go
file, _ := os.Create("output.txt")
defer file.Close()

writer := bufio.NewWriter(file)
defer writer.Flush()  // CRITICAL: Flush before closing!

writer.WriteString("Hello, world!\n")
```

**Important:** Always call `Flush()` before closing the file, or buffered data will be lost!

### Writing Methods

#### 1. Write([]byte) - Write Byte Slice
```go
data := []byte("Hello\n")
n, err := writer.Write(data)
```

#### 2. WriteString(string) - Write String
```go
n, err := writer.WriteString("Hello, world!\n")
```

#### 3. WriteByte(byte) - Write Single Byte
```go
err := writer.WriteByte('A')
```

#### 4. WriteRune(rune) - Write Single Character
```go
n, err := writer.WriteRune('日')  // Multi-byte UTF-8 character
```

### Flush Behavior

**When does bufio.Writer flush to disk?**
1. Buffer is full (default 4KB)
2. You call `Flush()` manually
3. (Never automatically on `Close()`—you MUST call `Flush()`)

**Example: Ensuring data is written**
```go
writer := bufio.NewWriter(file)

writer.WriteString("Important data\n")
writer.Flush()  // Ensure it's on disk NOW

// vs.

writer.WriteString("Data 1\n")
writer.WriteString("Data 2\n")
writer.WriteString("Data 3\n")
writer.Flush()  // Flush all three at once (more efficient)
```

---

## Performance Comparison: Buffered vs Unbuffered

### Benchmark: Reading a 100MB File

**Unbuffered (direct os.File):**
```go
file, _ := os.Open("100mb.txt")
defer file.Close()

buf := make([]byte, 1)
for {
    _, err := file.Read(buf)
    if err == io.EOF {
        break
    }
}
// Time: ~30 seconds (100M system calls!)
```

**Buffered (bufio.Reader):**
```go
file, _ := os.Open("100mb.txt")
defer file.Close()

reader := bufio.NewReader(file)
for {
    _, err := reader.ReadByte()
    if err == io.EOF {
        break
    }
}
// Time: ~0.5 seconds (only ~25,000 system calls)
```

**Speed improvement:** **60x faster!**

### Benchmark: Writing 1 Million Lines

**Unbuffered:**
```go
file, _ := os.Create("output.txt")
defer file.Close()

for i := 0; i < 1000000; i++ {
    file.WriteString(fmt.Sprintf("Line %d\n", i))
}
// Time: ~20 seconds
```

**Buffered:**
```go
file, _ := os.Create("output.txt")
defer file.Close()

writer := bufio.NewWriter(file)
defer writer.Flush()

for i := 0; i < 1000000; i++ {
    writer.WriteString(fmt.Sprintf("Line %d\n", i))
}
// Time: ~1 second
```

**Speed improvement:** **20x faster!**

---

## Streaming Large Files: Practical Patterns

### Pattern 1: Line-by-Line Processing (Logs, CSVs)

```go
func processLogFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()

        // Process line (e.g., parse, filter, aggregate)
        if strings.Contains(line, "ERROR") {
            fmt.Printf("Line %d: %s\n", lineNum, line)
        }
    }

    return scanner.Err()  // Return any scan errors
}
```

**Memory usage:** Constant (~4KB), regardless of file size.

### Pattern 2: Filtering and Transforming Files

```go
func filterFile(inputPath, outputPath string, predicate func(string) bool) error {
    // Open input file
    input, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer input.Close()

    // Create output file
    output, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer output.Close()

    // Buffered reader and writer
    scanner := bufio.NewScanner(input)
    writer := bufio.NewWriter(output)
    defer writer.Flush()

    // Stream line-by-line
    for scanner.Scan() {
        line := scanner.Text()
        if predicate(line) {
            writer.WriteString(line + "\n")
        }
    }

    return scanner.Err()
}

// Usage: Extract only ERROR lines
filterFile("app.log", "errors.log", func(line string) bool {
    return strings.Contains(line, "ERROR")
})
```

**Why this works well:**
- Input and output are both streamed (constant memory)
- No temporary storage of entire file
- Can process files larger than available RAM

### Pattern 3: Aggregating Statistics

```go
func wordFrequency(filename string) (map[string]int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    counts := make(map[string]int)
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)  // Split on words, not lines

    for scanner.Scan() {
        word := strings.ToLower(scanner.Text())
        counts[word]++
    }

    return counts, scanner.Err()
}
```

**Memory usage:** O(unique words), not O(file size).

### Pattern 4: Chunked Processing with Custom Buffer Sizes

```go
func processInChunks(filename string, chunkSize int) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := bufio.NewReaderSize(file, chunkSize)
    buffer := make([]byte, chunkSize)

    for {
        n, err := reader.Read(buffer)
        if n > 0 {
            processChunk(buffer[:n])
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }

    return nil
}
```

**Use case:** Processing binary files, network protocols, etc.

---

## Common Pitfalls and How to Avoid Them

### Pitfall 1: Forgetting to Flush

```go
// WRONG: Data might not be written!
writer := bufio.NewWriter(file)
writer.WriteString("Important data\n")
file.Close()  // BUG: Buffer not flushed!
```

**Fix:**
```go
writer := bufio.NewWriter(file)
defer writer.Flush()  // Ensure flush happens
writer.WriteString("Important data\n")
```

### Pitfall 2: Not Checking Scanner.Err()

```go
// WRONG: Ignores scan errors
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    process(scanner.Text())
}
// What if scanner encountered an error?
```

**Fix:**
```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    process(scanner.Text())
}
if err := scanner.Err(); err != nil {
    return fmt.Errorf("scan error: %w", err)
}
```

### Pitfall 3: Reusing Scanner.Text() or Scanner.Bytes()

```go
// WRONG: scanner.Text() is invalidated on next Scan()!
var lines []string
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    lines = append(lines, scanner.Text())  // BUG: All pointers point to same buffer!
}
```

**Fix:** Copy the string if you need to store it:
```go
var lines []string
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    lines = append(lines, string(scanner.Text()))  // Explicit copy
}
```

Actually, `scanner.Text()` already returns a copy as a string, but `scanner.Bytes()` does NOT:
```go
// WRONG with Bytes():
var lines [][]byte
for scanner.Scan() {
    lines = append(lines, scanner.Bytes())  // BUG: All slices share same buffer!
}

// FIX:
var lines [][]byte
for scanner.Scan() {
    line := make([]byte, len(scanner.Bytes()))
    copy(line, scanner.Bytes())
    lines = append(lines, line)
}
```

### Pitfall 4: Token Too Long Error

```go
// File has a line > 64KB
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    // ...
}
// scanner.Err() == bufio.ErrTooLong
```

**Fix:** Increase buffer size:
```go
scanner := bufio.NewScanner(file)
buf := make([]byte, 1024*1024)  // 1MB buffer
scanner.Buffer(buf, 10*1024*1024)  // Allow up to 10MB tokens
```

### Pitfall 5: Not Closing Files

```go
// WRONG: File descriptor leak!
func readFile(path string) string {
    file, _ := os.Open(path)
    // Missing: defer file.Close()
    scanner := bufio.NewScanner(file)
    // ...
}
```

**Fix:** Always defer Close():
```go
file, err := os.Open(path)
if err != nil {
    return "", err
}
defer file.Close()  // Guaranteed to run
```

---

## Real-World Applications

### 1. Log Analysis (DevOps)

**Problem:** Analyze 100GB of server logs to find error patterns.

```go
func analyzeErrors(logPath string) (map[string]int, error) {
    file, _ := os.Open(logPath)
    defer file.Close()

    errorCounts := make(map[string]int)
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "ERROR") {
            // Extract error type and count it
            errorType := extractErrorType(line)
            errorCounts[errorType]++
        }
    }

    return errorCounts, scanner.Err()
}
```

**Used by:** Splunk, Datadog, ELK stack internals

### 2. ETL Pipelines (Data Engineering)

**Problem:** Transform CSV files with millions of rows.

```go
func transformCSV(inputPath, outputPath string) error {
    in, _ := os.Open(inputPath)
    defer in.Close()
    out, _ := os.Create(outputPath)
    defer out.Close()

    scanner := bufio.NewScanner(in)
    writer := bufio.NewWriter(out)
    defer writer.Flush()

    for scanner.Scan() {
        row := scanner.Text()
        transformedRow := transform(row)  // Apply transformation
        writer.WriteString(transformedRow + "\n")
    }

    return scanner.Err()
}
```

**Used by:** Apache Spark, AWS Glue, data warehouses

### 3. Chat Server (Real-Time Communication)

**Problem:** Read/write messages over network connections.

```go
func handleClient(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            return
        }

        response := processMessage(message)
        writer.WriteString(response + "\n")
        writer.Flush()  // Send immediately
    }
}
```

**Used by:** Slack, Discord, IRC servers

### 4. Database Import (Batch Processing)

**Problem:** Import millions of records from a text file.

```go
func importRecords(db *sql.DB, filepath string) error {
    file, _ := os.Open(filepath)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    tx, _ := db.Begin()  // Use a transaction for speed

    for scanner.Scan() {
        record := parseRecord(scanner.Text())
        tx.Exec("INSERT INTO table VALUES (?)", record)
    }

    return tx.Commit()
}
```

**Used by:** PostgreSQL COPY, MySQL LOAD DATA

---

## Choosing the Right Tool

| Task | Tool | Reason |
|------|------|--------|
| Read line-by-line (text files) | `bufio.Scanner` | Simplest API, handles edge cases |
| Read with custom delimiter | `bufio.Reader.ReadString()` | More control than Scanner |
| Read binary data | `bufio.Reader.Read()` | Byte-oriented, no line assumptions |
| Write many small writes | `bufio.Writer` | Batches writes, reduces syscalls |
| Read huge lines (>64KB) | `bufio.Scanner` + custom buffer | Increase token size limit |
| Process fixed-size chunks | `bufio.Reader` + custom buffer | Full control over chunk size |

---

## How to Run

```bash
# Run the demonstration program
cd minis/17-file-streaming-bufio
go run cmd/file-stream/main.go

# Run tests
cd exercise
go test -v

# Run benchmarks
go test -bench=. -benchmem
```

---

## Expected Output (Demo Program)

```
=== Creating Test Files ===
Created small.txt (10 lines, 50 bytes)
Created large.txt (1,000,000 lines, ~20MB)

=== Example 1: Reading Lines with Scanner ===
Line 1: This is line 1 of the test file
Line 2: This is line 2 of the test file
...

=== Example 2: Word Count ===
Total words: 2,000,000

=== Example 3: Buffered Writing ===
Wrote 100,000 lines to output.txt in 45ms

=== Example 4: Performance Comparison ===
Unbuffered read: 1.2s
Buffered read:    0.08s
Speedup:          15x faster
```

---

## Key Takeaways

1. **Always use buffered I/O** for files (unless you have a specific reason not to)
2. **Scanner is best for line-oriented text** (logs, CSVs, config files)
3. **Always flush Writers** before closing files
4. **Check scanner.Err()** after the scan loop
5. **Stream large files** line-by-line to avoid memory exhaustion
6. **Pre-allocate buffers** if you know the typical line/token size
7. **Close files with defer** to prevent leaks

---

## Connections to Other Projects

- **Project 03 (csv-stats)**: Could use Scanner for more efficient CSV reading
- **Project 04 (jsonl-log-filter)**: Streaming JSONL with Scanner
- **Project 06 (worker-pool-wordcount)**: Combine streaming with concurrency
- **Project 11 (slices-internals)**: Understanding buffer growth and allocation
- **Project 28 (pprof-benchmarks)**: Profiling buffered vs unbuffered I/O

---

## Stretch Goals

1. **Build a `tail -f` clone** that follows a growing log file in real-time
   - Hint: Use `file.Seek()` and poll for new data

2. **Implement a line-oriented database** with indexing
   - Hint: Store byte offsets for fast random access

3. **Create a CSV parser** that handles quoted fields with commas
   - Hint: Write a custom `SplitFunc` for Scanner

4. **Benchmark different buffer sizes** (1KB, 4KB, 64KB, 1MB)
   - Hint: Use `testing.B` and vary `bufio.NewReaderSize()`

5. **Build a log rotator** that splits files when they exceed a size limit
   - Hint: Track bytes written and create new files as needed
