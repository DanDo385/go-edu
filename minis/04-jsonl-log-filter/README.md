# Project 04: jsonl-log-filter - Understanding JSON and Log Processing

## What Is This Project About?

Imagine you're a developer debugging a production system. Your server has been running for days, generating hundreds of thousands of log entries. You need to find all ERROR-level messages from the last hour, sorted by time. How do you process this efficiently without loading a 10GB log file into Excel?

This project teaches you how to:
1. Parse JSONL (JSON Lines) format - one JSON object per line
2. Convert string log levels ("error", "warn") to typed enums
3. Filter and sort log entries by severity and timestamp
4. Handle malformed data gracefully (skip bad lines, report count)

## First Principles: What Is JSON and Why JSONL?

### Understanding JSON

JSON (JavaScript Object Notation) is a text format for structured data. It uses familiar punctuation:

```json
{
  "ts": "2024-01-01T12:00:00Z",
  "level": "error",
  "msg": "Database connection failed"
}
```

- `{}` curly braces = object (like a map/dictionary)
- `""` quotes = strings
- `:` colon = maps key to value
- `,` comma = separates items

Think of JSON as a way to write Go structs as text that any program can read.

### The Problem with JSON Arrays

You might store logs as a JSON array:

```json
[
  {"ts": "...", "level": "info", "msg": "Started"},
  {"ts": "...", "level": "error", "msg": "Failed"},
  ...1 million more entries...
]
```

**Problems**:
1. Must read **entire file** before parsing (can't stream)
2. One malformed entry breaks the **entire** array
3. Can't append new entries without rewriting the file
4. Opening `[` and closing `]` must be valid

### The JSONL Solution

JSONL (JSON Lines) puts **one JSON object per line**:

```
{"ts":"2024-01-01T12:00:00Z","level":"info","msg":"Server started"}
{"ts":"2024-01-01T12:00:05Z","level":"error","msg":"Database failed"}
{"ts":"2024-01-01T12:00:10Z","level":"warn","msg":"High memory"}
```

**Advantages**:
1. **Stream processing**: Read line-by-line (constant memory)
2. **Append-friendly**: Just add new line to file end
3. **Fault-tolerant**: Skip malformed lines, process the rest
4. **Standard format**: Used by Elasticsearch, MongoDB, many tools

## The Problem We're Solving

**Input**: JSONL file with log entries
**Filter**: Keep only entries >= minimum severity level
**Sort**: Order by timestamp (oldest first)
**Output**: Filtered, sorted entries + error if any lines were skipped

Example:
```
Input:
{"ts":"2024-01-01T12:00:00Z","level":"info","msg":"A"}
{this is malformed
{"ts":"2024-01-01T12:00:05Z","level":"error","msg":"B"}
{"ts":"2024-01-01T12:00:02Z","level":"debug","msg":"C"}

Filter: level >= warn (excludes debug, info)
Result: [{"ts":"...05Z","level":"error","msg":"B"}]
Error: "skipped 3 malformed lines"
```

## Breaking Down the Solution (Step by Step)

### Step 1: Understanding Log Levels as an Enum

In many languages, log levels are strings: "debug", "info", "warn", "error".

**Problem**: Strings are hard to compare. Is "warn" > "info"? You'd need to check spelling, handle case sensitivity, etc.

**Solution**: Use an **enum** (enumeration) - assign each level an integer:

```go
const (
    Debug Level = 0  // Lowest
    Info  Level = 1
    Warn  Level = 2
    Error Level = 3  // Highest
)
```

Now we can compare: `entry.Level >= Warn` is simple integer comparison!

### Step 2: Custom JSON Unmarshaling

JSON has no native enum type. Log levels come as strings:
```json
{"level": "error"}
```

We need to teach Go: "When you see the string 'error', convert it to the number 3 (Error enum value)".

This is done by implementing the `UnmarshalJSON` interface:

```go
func (l *Level) UnmarshalJSON(data []byte) error {
    var s string
    json.Unmarshal(data, &s)  // Get the string

    switch strings.ToLower(s) {
    case "debug": *l = Debug
    case "info": *l = Info
    case "warn": *l = Warn
    case "error": *l = Error
    default: return fmt.Errorf("invalid level: %q", s)
    }
    return nil
}
```

**How it works**:
1. JSON package calls our `UnmarshalJSON` method when it sees a `Level` field
2. We receive the raw JSON bytes: `"error"` (with quotes!)
3. We unmarshal to string, convert to lowercase, map to enum value
4. We modify the Level pointer (`*l = Error`) to set the value

### Step 3: Processing JSONL Line by Line

Think of JSONL processing like reading a book line by line:

```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {              // Read next line
    line := scanner.Text()        // Get the line as string

    var entry Entry
    err := json.Unmarshal([]byte(line), &entry)
    if err != nil {
        skippedCount++            // Track malformed lines
        continue                  // Skip, don't crash
    }

    // Process valid entry...
}
```

**Key insight**: We process one line at a time, so a 10GB log file uses only a few KB of memory!

### Step 4: Filtering by Level

Once we have log levels as integers, filtering is trivial:

```go
if entry.Level >= minLevel {
    entries = append(entries, entry)
}
```

Examples:
- `minLevel = Warn` (2): Keeps Warn (2) and Error (3), filters Debug (0) and Info (1)
- `minLevel = Debug` (0): Keeps everything
- `minLevel = Error` (3): Only errors

### Step 5: Sorting by Timestamp

Go's `time.Time` type represents a moment in time. We can compare them:

```go
t1 := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
t2 := time.Parse(time.RFC3339, "2024-01-01T12:00:05Z")
t1.Before(t2)  // true (t1 is earlier)
```

To sort a slice of entries, we use `sort.Slice` with a comparison function:

```go
sort.Slice(entries, func(i, j int) bool {
    return entries[i].TS.Before(entries[j].TS)
})
```

**How it works**:
- The function says "return true if i should come before j"
- `sort.Slice` uses this to sort the slice in-place
- We're saying "entry i comes before j if its timestamp is earlier"

## The Complete Solution (Explained)

```go
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
    var entries []Entry
    var skippedCount int

    // Create scanner for line-by-line reading
    scanner := bufio.NewScanner(r)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()

        // Skip blank lines gracefully
        if strings.TrimSpace(line) == "" {
            continue
        }

        // Try to parse JSON
        var entry Entry
        if err := json.Unmarshal([]byte(line), &entry); err != nil {
            // Malformed JSON - skip but count
            skippedCount++
            continue
        }

        // Filter by level
        if entry.Level >= minLevel {
            entries = append(entries, entry)
        }
    }

    // Check for I/O errors (not EOF)
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("reading input: %w", err)
    }

    // Sort by timestamp (oldest first)
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].TS.Before(entries[j].TS)
    })

    // Return partial success with error if lines were skipped
    var err error
    if skippedCount > 0 {
        err = fmt.Errorf("skipped %d malformed lines", skippedCount)
    }

    return entries, err
}
```

## Key Concepts Explained

### Why `time.Time` Instead of String?

We could store timestamps as strings:
```go
type Entry struct {
    TS string `json:"ts"`
}
```

**Problems**:
- Can't compare directly: Is "12:00:05" > "12:00:02"? (Lexicographically, yes. Chronologically... also yes, but only if same day!)
- Can't handle timezones: "12:00:00Z" (UTC) vs "13:00:00+01:00" (Berlin, same time!)
- Can't do math: "What's 5 minutes after this time?"

`time.Time` handles all of this:
```go
ts.Before(other)           // Comparison
ts.Add(5 * time.Minute)    // Math
ts.Format("2006-01-02")    // Formatting
```

### Custom Unmarshaling Pattern

The pattern for custom JSON unmarshaling:

```go
type MyType int

func (m *MyType) UnmarshalJSON(data []byte) error {
    // 1. Unmarshal to a simple type
    var temp string
    if err := json.Unmarshal(data, &temp); err != nil {
        return err
    }

    // 2. Convert and validate
    switch temp {
    case "value1": *m = 1
    case "value2": *m = 2
    default: return fmt.Errorf("invalid: %q", temp)
    }

    return nil
}
```

### Error Accumulation vs Fail-Fast

Two philosophies:

**Fail-fast** (Stop on first error):
```go
if err != nil {
    return nil, err
}
```
Pros: Catches problems immediately
Cons: Partial results lost

**Error accumulation** (Continue, report issues):
```go
if err != nil {
    skippedCount++
    continue
}
// At end:
if skippedCount > 0 {
    return results, fmt.Errorf("skipped %d items", skippedCount)
}
```
Pros: Get partial results, know total failure count
Cons: May mask systemic issues

Our solution uses accumulation—better for log processing (one bad line shouldn't abort analysis).

## Common Patterns You're Learning

### Pattern 1: JSONL Processing
```go
scanner := bufio.NewScanner(r)
for scanner.Scan() {
    line := scanner.Text()
    var obj MyType
    if err := json.Unmarshal([]byte(line), &obj); err != nil {
        continue  // Skip malformed
    }
    // Process obj...
}
```

### Pattern 2: Enum-Style Constants
```go
type Level int
const (
    Debug Level = iota  // 0
    Info                // 1
    Warn                // 2
    Error               // 3
)
```

### Pattern 3: Sorting with Custom Comparator
```go
sort.Slice(items, func(i, j int) bool {
    return items[i].Field < items[j].Field
})
```

## Real-World Applications

1. **Log Analysis**: Filtering production logs by severity (Splunk, ELK stack)
2. **Metrics Processing**: Time-series data from monitoring systems
3. **Event Streaming**: Processing Kafka/Kinesis messages
4. **Data Pipelines**: ETL jobs that transform JSONL data
5. **Debugging**: Finding errors in large log dumps

## How to Run

```bash
# Prepare the project
cd minis/04-jsonl-log-filter/exercise
mv solution.go solution.go.reference

# Look at the test data
cat ../testdata/logs.jsonl

# Run tests
go test -v

# Implement your solution in exercise.go
# Then test again
go test

# Run the demo program
cd ../..
make run P=04-jsonl-log-filter
```

## Testing with Time

The tests use `time.Parse` to create timestamps:

```go
ts, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
```

`RFC3339` is the standard format: `YYYY-MM-DDTHH:MM:SSZ`

## Common Mistakes to Avoid

1. **Trying to unmarshal entire file as JSON array**: JSONL is NOT a JSON array!
2. **Not handling blank lines**: Always `strings.TrimSpace()` and check `== ""`
3. **Failing on one bad line**: Use `continue`, track `skippedCount`
4. **String level comparison**: Convert to enum for reliable comparison
5. **Forgetting to sort**: Results should be chronological, not input order
6. **Not handling timezones**: `time.Time` handles this, string timestamps don't

## Stretch Goals

1. **Add time range filtering**: `--after "2024-01-01" --before "2024-01-02"`
   - Use `time.After()` and `time.Before()` methods

2. **Full-text search**: Filter by regex on message field
   - Use `regexp.MatchString(pattern, entry.Msg)`

3. **Colorized output**: Red for errors, yellow for warnings
   - Use ANSI escape codes: `\033[31m` for red

4. **Structured output**: Write results as JSONL to a file
   - Use `json.Marshal()` + write line by line

5. **Performance benchmark**: Compare parsing 10,000 logs vs 100,000
   - Use `testing.B` to measure throughput
