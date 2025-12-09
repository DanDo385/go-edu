//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "bufio" for line-by-line reading (efficient buffered I/O)
// - "encoding/json" for parsing JSON strings into structs
// - "fmt" for error formatting
// - "io" for the Reader interface
// - "sort" for sorting slices with custom comparators
// - "strings" for string manipulation (trimming whitespace)
// - "time" for timestamp handling
//
// import (
//     "bufio"
//     "encoding/json"
//     "fmt"
//     "io"
//     "sort"
//     "strings"
//     "time"
// )

// Level represents log severity (enum-like type in Go)
// This is an INTEGER type, not a string, for efficient comparisons
//
// TODO: Define the Level type
// type Level int
//
// Key Go concepts:
// - Type alias: Level is a distinct type (not just an int)
// - You can define methods on it (like UnmarshalJSON below)
// - Enum pattern: Use iota to auto-generate sequential integer constants

// TODO: Define severity level constants using iota
// const (
//     Debug Level = iota  // 0 - Most verbose, includes everything
//     Info                // 1 - Informational messages
//     Warn                // 2 - Warning messages
//     Error               // 3 - Error messages
// )
//
// Key Go concepts:
// - iota is a constant generator that starts at 0 and increments
// - First const gets iota (0), subsequent consts auto-increment (1, 2, 3, ...)
// - This creates an ordered enum where Debug < Info < Warn < Error
// - Integer comparison (>=) is fast and safe

// Entry represents a single log entry
// This is a STRUCT TYPE (value type, not reference type)
//
// TODO: Define the Entry struct with JSON tags
// type Entry struct {
//     TS    time.Time `json:"ts"`     // Timestamp in RFC3339 format
//     Level Level     `json:"level"`  // Severity level (will use custom unmarshaler)
//     Msg   string    `json:"msg"`    // Log message text
// }
//
// Key Go concepts for struct tags:
// - Struct tags are metadata attached to fields
// - `json:"ts"` tells encoding/json to map JSON field "ts" to TS field
// - Tags are accessed via reflection (encoding/json uses them internally)
// - time.Time is a struct (value type) that handles timezones correctly
// - When Entry is passed by value, all fields are copied
// - time.Time is ~24 bytes, Level is ~8 bytes, string is ~16 bytes (ptr+len)

// UnmarshalJSON implements custom JSON unmarshaling for Level
// This method allows us to convert JSON strings ("debug", "info") to Level enum values
//
// TODO: Implement UnmarshalJSON method for Level type
// Function signature: func (l *Level) UnmarshalJSON(data []byte) error
//
// WHY POINTER RECEIVER (*Level, not Level)?
// - This method MUST modify the Level value that's being unmarshaled
// - If we used (l Level), we'd get a COPY - changes wouldn't affect original
// - Pointer receiver allows us to write: *l = Debug (dereference and assign)
// - json.Unmarshal will call this method with a pointer to the Level field
//
// Steps to implement:
//
// 1. Declare a string variable to hold the unmarshaled JSON value
//    - Use: var s string
//    - This will receive the string value from JSON (e.g., "info")
//
// 2. Unmarshal the raw JSON bytes into the string
//    - Use: if err := json.Unmarshal(data, &s); err != nil { return err }
//    - data is []byte containing the JSON value (e.g., `"info"` with quotes)
//    - json.Unmarshal needs a pointer (&s) to write the result
//    - This handles unquoting: `"info"` becomes just `info`
//    - If data is not a valid JSON string, return error with fmt.Errorf
//
// 3. Convert string to Level enum using a switch statement
//    - Use strings.ToLower(s) for case-insensitive matching
//    - Match cases: "debug" → Debug, "info" → Info, "warn" → Warn, "error" → Error
//    - Also accept "warning" as an alias for Warn (multiple values in one case)
//    - For each match: *l = <Level constant> (dereference pointer and assign)
//    - Default case: return fmt.Errorf("invalid level: %q", s)
//
// 4. Return nil on success
//
// Key Go concepts:
// - Pointer receiver: Method can modify the receiver value
// - Dereferencing: *l = value writes to the memory location l points to
// - Interface satisfaction: This method makes Level implement json.Unmarshaler
// - json.Unmarshal automatically calls this when unmarshaling into a Level field
// - Error handling: Return descriptive errors for invalid input
//
// Common mistakes:
// - Using value receiver (l Level) instead of pointer (*Level)
// - Forgetting to dereference: l = Debug instead of *l = Debug
// - Not handling the JSON quotes (use json.Unmarshal, not manual parsing)

// TODO: Implement the UnmarshalJSON method below
// func (l *Level) UnmarshalJSON(data []byte) error {
//     return nil
// }

// FilterLogs reads JSONL (JSON Lines) from r, filters entries by minimum level, and sorts by timestamp
//
// JSONL format: One JSON object per line (NOT a JSON array!)
//   {"ts":"2024-01-01T12:00:00Z","level":"info","msg":"Server started"}
//   {"ts":"2024-01-01T12:00:05Z","level":"error","msg":"Database failed"}
//
// Returns:
// - []Entry: All log entries where entry.Level >= minLevel, sorted by timestamp
// - error: Non-nil if any lines were skipped (malformed JSON), includes count
//
// TODO: Implement FilterLogs function
// Function signature: func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error)
//
// PARAMETER PASSING SEMANTICS:
// - r is an io.Reader interface (passed by value, but interfaces are small)
//   * Interface is: struct { type *typeInfo; value *actualData }
//   * Copying the interface copies the pointers, not the underlying data
//   * This is why you can read from r - it points to the actual reader
//
// - minLevel is a Level (int) - passed by value (copied)
//   * Level is just an int (~8 bytes), so copying is cheap
//   * We get our own copy, caller's minLevel is unchanged
//
// RETURN VALUE SEMANTICS:
// - []Entry is a SLICE (reference type)
//   * Slice is: struct { ptr *Entry; len, cap int }
//   * Returning a slice copies the slice header (~24 bytes), not the data
//   * Caller gets a slice pointing to the same underlying array
//   * Both caller and function can see changes to array elements
//
// - error is an interface
//   * nil or a pointer to an error struct
//   * Copying the interface is cheap
//
// Steps to implement:
//
// 1. Initialize variables
//    - var entries []Entry (starts as nil slice - no allocation yet)
//    - var skipped int (counter for malformed lines)
//    - Slices in Go:
//      * nil slice has len=0, cap=0, ptr=nil
//      * append() will allocate backing array when first element is added
//      * Capacity grows exponentially (roughly doubles each time)
//
// 2. Create a buffered scanner for line-by-line reading
//    - Use: scanner := bufio.NewScanner(r)
//    - bufio.Scanner is a STRUCT (value type) but contains pointers internally
//    - Scanner maintains internal buffer (default 64KB) for efficiency
//    - Reads from io.Reader and splits by newlines automatically
//    - Memory: Scanner buffer is reused for each line (no allocations per line)
//
// 3. Loop through lines
//    - Use: for scanner.Scan() { ... }
//    - scanner.Scan() advances to next line, returns false at EOF or error
//    - Get line text with: line := scanner.Text()
//      * Returns a string (points to scanner's internal buffer)
//      * IMPORTANT: This string is only valid until next Scan() call
//      * If you need to keep it, make a copy: line = string([]byte(line))
//
// 4. Skip empty lines
//    - Use: if strings.TrimSpace(line) == "" { continue }
//    - TrimSpace removes leading/trailing whitespace
//    - Empty lines are common in log files (blank lines between sections)
//
// 5. Parse JSON line into Entry struct
//    - Declare: var e Entry (zero value: TS is zero time, Level is 0, Msg is "")
//    - Use: err := json.Unmarshal([]byte(line), &e)
//      * Convert string to []byte (makes a copy of string data)
//      * Pass &e (pointer) so Unmarshal can write to the struct
//      * Unmarshal uses reflection to match JSON fields to struct fields
//      * For Level field, it calls our custom UnmarshalJSON method
//    - Handle error (malformed JSON):
//      * if err != nil { skipped++; continue }
//      * This is "fail-soft" approach: skip bad lines but keep processing
//      * Alternative: "fail-fast" would return error immediately
//
// 6. Filter by level
//    - Check: if e.Level >= minLevel { ... }
//    - Since Level is an int enum, >= comparison is simple and fast
//    - Only keep entries that meet or exceed the minimum severity
//
// 7. Append matching entries
//    - Use: entries = append(entries, e)
//    - append() copies e (the Entry struct) into the slice
//    - MEMORY ALLOCATION:
//      * If len == cap, append allocates new backing array (larger capacity)
//      * Copies all existing elements to new array
//      * Returns new slice header pointing to new array
//      * Old array can be garbage collected if no other references
//    - MUST assign result: entries = append(...) because slice header may change
//
// 8. Check for scanner errors after loop
//    - Use: if err := scanner.Err(); err != nil { return nil, err }
//    - scanner.Err() returns nil if we reached EOF normally
//    - Returns non-nil for I/O errors (disk failure, network timeout, etc.)
//    - This separates I/O errors from malformed JSON (which we handled above)
//
// 9. Sort entries by timestamp
//    - Use: sort.Slice(entries, func(i, j int) bool { ... })
//    - sort.Slice takes a "less" function (comparator)
//    - Function signature: func(i, j int) bool
//      * i, j are indices into the entries slice
//      * Return true if entries[i] should come before entries[j]
//    - Comparison: return entries[i].TS.Before(entries[j].TS)
//      * time.Time.Before() compares timestamps correctly (handles timezones)
//      * Alternative: entries[i].TS.Compare(entries[j].TS) < 0 (Go 1.20+)
//    - MEMORY:
//      * sort.Slice sorts IN-PLACE (modifies entries slice)
//      * Uses quicksort internally (O(n log n) comparisons, O(log n) stack)
//      * No allocations (just swaps elements in existing array)
//
// 10. Return results
//     - If skipped > 0:
//       * return entries, fmt.Errorf("skipped %d malformed lines", skipped)
//       * This is "partial success": caller gets valid entries AND an error
//       * Caller can decide whether to use data or treat error as fatal
//     - If skipped == 0:
//       * return entries, nil
//       * Complete success
//
// Key Go concepts:
// - Slices are reference types (share backing array)
// - append() may reallocate (must assign result)
// - Structs are copied when appended to slice
// - bufio.Scanner for efficient line-by-line reading
// - json.Unmarshal uses reflection and struct tags
// - Custom unmarshalers (UnmarshalJSON method)
// - sort.Slice with inline comparator function
// - Partial success pattern (return both data and error)
// - Fail-soft vs fail-fast error handling
//
// Common mistakes:
// - Forgetting to assign append result: append(entries, e) without entries =
// - Using scanner.Text() string after next Scan() call (buffer reuse)
// - Not checking scanner.Err() after loop (missing I/O errors)
// - Comparing times with < instead of time.Before() (doesn't work!)
// - Not handling empty input (should return [], nil - not an error)

// TODO: Implement the FilterLogs function below
// func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
//     return nil, nil
// }

// After implementing all functions:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Try with malformed JSON to test error handling
// - Compare with solution.go to see detailed explanations
