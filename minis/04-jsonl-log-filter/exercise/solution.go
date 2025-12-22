//go:build solution
// +build solution

/*
Problem: Parse and filter JSONL (JSON Lines) log entries by severity level

Given a file with one JSON object per line, we need to:
1. Parse each line as a structured log entry
2. Filter by minimum severity level (debug < info < warn < error)
3. Sort results by timestamp
4. Handle malformed lines gracefully (skip but report count)

Constraints:
- JSONL format: one JSON object per line (not a JSON array!)
- Timestamps are RFC3339 format (e.g., "2024-01-01T12:00:00Z")
- Level is a string ("debug", "info", "warn", "error") that must map to an enum
- Malformed lines should be skipped, not cause total failure

Time/Space Complexity:
- Time: O(n log n) where n = number of valid entries (O(n) parse + O(n log n) sort)
- Space: O(n) to store filtered entries

Why Go is well-suited:
- `encoding/json`: Robust JSON parser with struct mapping via reflection
- Custom unmarshalers: `UnmarshalJSON` interface for enum-like types
- `time.Time`: First-class time support with zone awareness
- `sort.Slice`: Inline sorting with custom comparators
- Error accumulation: Handle partial failures gracefully without exceptions

Performance Characteristics:
- Go's JSON parsing uses reflection, which is slower than Rust's compile-time codegen
  but faster than Python's pure interpretation
- Memory allocation is handled by Go's garbage collector, eliminating manual management
  required in C while being more predictable than Python's reference counting
- String operations are UTF-8 aware and bounds-checked, unlike C's pointer arithmetic
*/

package exercise

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Level represents log severity as an integer enum.
//
// Memory Layout:
// - Level is stored as an int (32 or 64 bits depending on platform)
// - This is much more efficient than storing strings ("debug", "info", etc.)
// - Comparison operations (>=, <) are single CPU instructions versus string comparisons
//
// Speed Comparison:
// - C: Would use enum, identical performance (single integer comparison)
// - Rust: Would use enum with #[repr(u8)], potentially 1 byte vs Go's 4/8 bytes
// - Python: Would use Enum or just strings, significantly slower (hash lookups)
// - Solidity: Would use uint8 enum, identical approach but in EVM context
type Level int

const (
	Debug Level = iota // iota generates: Debug=0, Info=1, Warn=2, Error=3
	Info
	Warn
	Error
)

// Entry represents a single log entry.
//
// Memory Layout:
// - This struct is laid out in memory contiguously
// - TS: 24 bytes (time.Time contains wall clock, monotonic clock, location pointer)
// - Level: 4 or 8 bytes (platform-dependent int)
// - Msg: 16 bytes (string header: pointer + length, data stored separately on heap)
// - Total: ~48 bytes per entry (plus heap allocation for Msg string data)
//
// The `json:"..."` tags are struct tags that the json package reads via reflection
// at runtime to map JSON field names to struct fields. This is zero-cost at runtime
// after the initial reflection setup.
//
// Speed Comparison:
// - C: Would use struct with char* for message, manual parsing required
// - Rust: serde derives would generate specialized parsing code at compile time (faster)
// - Python: Would use dict or dataclass, stored as PyObject* with overhead (~200+ bytes)
// - Solidity: N/A (no native JSON support, would need external oracle for off-chain data)
type Entry struct {
	TS    time.Time `json:"ts"`    // RFC3339 timestamp, parsed by time package
	Level Level     `json:"level"` // Uses our custom UnmarshalJSON below
	Msg   string    `json:"msg"`   // String data lives on heap, referenced by 16-byte header
}

// UnmarshalJSON implements custom JSON unmarshaling for Level.
//
// Go Concepts Demonstrated:
// - Implementing the json.Unmarshaler interface (method set on pointer receiver)
// - Converting string to enum-like values via switch statement
// - Error handling in custom unmarshalers (returning error stops unmarshaling)
//
// How This Works Under the Hood:
// 1. When json.Unmarshal encounters a Level field, it checks if *Level implements json.Unmarshaler
// 2. If yes, it calls this method instead of default unmarshaling behavior
// 3. We receive raw JSON bytes (e.g., []byte(`"info"`)) including the quotes
// 4. We unmarshal to string first (to strip quotes and handle escaping)
// 5. Then map the string to our integer enum value
// 6. The *l assignment modifies the actual Level value in the parent Entry struct
//
// Why Pointer Receiver?
// We need *Level not Level because we're modifying the value. Go passes by value,
// so we need a pointer to modify the original memory location in the Entry struct.
//
// Speed Comparison:
// - C: Would use strcmp in switch, similar speed but manual error handling
// - Rust: serde would generate this at compile time (no runtime reflection), faster
// - Python: Would use dict lookup, slower due to hash computation and PyObject overhead
// - Solidity: N/A, but similar concept when parsing bytes to enum in assembly
func (l *Level) UnmarshalJSON(data []byte) error {
	// First, unmarshal to a string to handle JSON escaping and quote removal.
	// The json package handles all edge cases: escaped quotes, unicode, etc.
	// Memory: Allocates a new string on the heap (data is copied, not referenced)
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("level must be a string: %w", err)
	}

	// Convert string to Level enum using switch statement.
	// This compiles to a jump table (similar to C switch) for O(1) lookup.
	// strings.ToLower allocates a new string but only if the string contains uppercase.
	//
	// Performance Note:
	// - This switch is compiled to efficient machine code (jump table or binary search)
	// - In hot paths, we could optimize by avoiding ToLower (just accept lowercase)
	// - Or use a map[string]Level for initialization, then direct lookup (O(1) hash)
	switch strings.ToLower(s) {
	case "debug":
		*l = Debug // Assigns 0 to the memory location pointed to by l
	case "info":
		*l = Info // Assigns 1
	case "warn", "warning": // Multiple case values (comma-separated)
		*l = Warn // Assigns 2
	case "error":
		*l = Error // Assigns 3
	default:
		// Return error if string doesn't match any valid level.
		// The fmt.Errorf allocation happens on the heap (escape analysis determines this).
		return fmt.Errorf("invalid level: %q (must be debug/info/warn/error)", s)
	}

	return nil
}

// FilterLogs reads JSONL from r, filters entries >= minLevel, and sorts by timestamp.
//
// Function Signature Breakdown:
// - io.Reader: Interface for any readable stream (files, network, memory buffers)
//   This demonstrates Go's composition over inheritance philosophy
// - []Entry: Slice (dynamic array) returned on the heap
// - error: Go's explicit error handling (no exceptions like Python)
//
// Memory Management Throughout:
// 1. Scanner allocates a buffer (~4KB default) for reading lines
// 2. Each Entry struct is appended to the slice (grows by doubling when capacity exceeded)
// 3. Slice backing array lives on heap (escape analysis determines this)
// 4. After function returns, all allocations are eligible for garbage collection
//
// Speed Comparison:
// - C: Faster raw parsing (no reflection) but requires manual memory management (malloc/free)
// - Rust: Similar speed to C with zero-cost abstractions, compile-time guarantees
// - Python: Much slower (interpreted, dynamic typing, reference counting overhead)
// - Solidity: N/A (cannot read files, but similar filtering logic would cost gas per iteration)
//
// Three-Input Iteration Table:
//
// Input 1: Valid JSONL (happy path)
//   Line 1: {"ts":"2024-01-01T10:00:00Z","level":"info","msg":"A"} → filtered out (info < warn)
//   Line 2: {"ts":"2024-01-01T11:00:00Z","level":"error","msg":"B"} → kept (error >= warn)
//   Line 3: {"ts":"2024-01-01T09:00:00Z","level":"warn","msg":"C"} → kept (warn >= warn)
//   After sort → [warn(09:00), error(11:00)] (sorted by timestamp, oldest first)
//
// Input 2: Empty input (edge case)
//   No lines → scanner.Scan() returns false immediately
//   Result: []Entry{} (empty slice), nil (no error)
//
// Input 3: Malformed JSON (partial failure)
//   Line 1: {"ts":"...","level":"info","msg":"A"} → filtered out (info < warn)
//   Line 2: {invalid json} → json.Unmarshal fails, skippedCount++, continue
//   Line 3: {"ts":"...","level":"error","msg":"C"} → kept
//   Result: []Entry{error entry}, error("skipped 2 lines") - partial success pattern
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	// Initialize slice with zero capacity. Go will allocate backing array dynamically.
	// Memory: nil slice (no allocation yet), capacity grows as we append: 0→1→2→4→8→16...
	// This doubling strategy gives amortized O(1) append, similar to C++ vector or Python list.
	var entries []Entry
	
	// Track number of lines we skip due to parse errors.
	// Memory: Single integer on the stack (typically 8 bytes)
	var skippedCount int

	// Create a buffered scanner for line-by-line reading.
	// Memory: Allocates internal buffer (~4KB default) on heap for reading chunks.
	// This is more efficient than reading byte-by-byte (fewer syscalls).
	//
	// Why Scanner vs Manual Reading?
	// - Scanner handles line splitting, buffer management, and edge cases (long lines)
	// - Alternative: bufio.Reader.ReadString('\n') requires manual loop
	//
	// Speed Comparison:
	// - C: fgets() with fixed buffer, similar performance, manual overflow handling
	// - Rust: BufRead::lines() iterator, equivalent performance with safety guarantees
	// - Python: for line in file, similar but slower due to interpreter overhead
	scanner := bufio.NewScanner(r)

	// Track line number for debugging (currently unused but useful for error messages).
	// Memory: Stack-allocated integer, incremented each iteration
	lineNum := 0

	// Main parsing loop: read each line until EOF or error.
	// scanner.Scan() returns false on EOF or error (check scanner.Err() to distinguish).
	//
	// Loop Mechanics:
	// 1. scanner.Scan() reads next line into internal buffer (stops at \n)
	// 2. scanner.Text() returns the line as a string (excludes \n)
	// 3. Loop body processes the line
	// 4. Repeat until Scan() returns false
	//
	// Performance: This is a streaming approach (O(1) memory per line) vs loading
	// entire file into memory (O(file_size)). Critical for large log files.
	for scanner.Scan() {
		lineNum++
		// Get the current line as a string.
		// Memory: scanner.Text() returns a string that references the scanner's internal buffer.
		// This is efficient (no copy) but the string is invalidated on next Scan() call.
		// Assignment to `line` makes a copy, so it persists across iterations.
		line := scanner.Text()

		// Skip empty lines (whitespace-only lines).
		// strings.TrimSpace allocates a new string if trimming is needed.
		// Alternative optimization: Check if line is empty directly (if len(line) == 0).
		// We use TrimSpace to also handle lines with only spaces/tabs.
		if strings.TrimSpace(line) == "" {
			continue // Jump back to for loop condition (scanner.Scan())
		}

		// Parse the JSON line into an Entry struct.
		// Memory: Entry is allocated on the stack initially (48 bytes).
		// If it escapes (added to slice), it will be moved to heap by escape analysis.
		var entry Entry
		
		// json.Unmarshal performs the following steps:
		// 1. Parse JSON bytes into intermediate representation (using reflection)
		// 2. Map JSON fields to struct fields using `json:"..."` tags
		// 3. Call UnmarshalJSON for types that implement json.Unmarshaler (our Level type)
		// 4. Validate types (string for string, number for int, etc.)
		//
		// Performance Cost:
		// - Reflection overhead: ~2-3x slower than hand-written parsing
		// - Memory allocations: JSON parser allocates intermediate objects
		// - But: Vastly more convenient than manual parsing, good enough for most use cases
		//
		// Speed Comparison:
		// - C: Manual parsing with strtok/sscanf, faster but error-prone (buffer overflows)
		// - Rust: serde_json generates specialized code at compile time, ~2x faster than Go
		// - Python: json.loads() is implemented in C, comparable speed but higher memory overhead
		// - Solidity: No equivalent, but abi.decode() for binary data is similar concept
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Malformed JSON: increment skip counter and continue.
			// This is "fail-soft" behavior: one bad line doesn't abort the entire process.
			// In production, you might log the error with line number for debugging.
			//
			// Alternative Approaches:
			// 1. Fail-fast: return nil, fmt.Errorf("line %d: %w", lineNum, err)
			//    Pros: Catches data quality issues early
			//    Cons: One bad line aborts everything
			// 2. Return detailed errors: []ParseError with line numbers
			//    Pros: Caller can inspect each failure
			//    Cons: More complex API
			skippedCount++
			continue
		}

		// Filter by level using integer comparison.
		// Since Level is an integer enum (debug=0, info=1, warn=2, error=3),
		// this compiles to a single CPU comparison instruction (CMP, JL in x86).
		//
		// Performance: O(1) constant time, single CPU cycle
		// Compare to string comparison: O(min_length) character-by-character scan
		//
		// Speed Comparison:
		// - C: Identical (enum comparison)
		// - Rust: Identical (enum with PartialOrd derived)
		// - Python: Much slower (string comparison or Enum value lookup)
		// - Solidity: Identical (uint8 comparison costs 3 gas)
		if entry.Level >= minLevel {
			// Append entry to slice.
			// Memory: If slice capacity is exhausted, Go allocates new backing array
			// (double the size), copies old elements, then appends new one.
			// Amortized O(1) time complexity, but occasional O(n) copy on resize.
			//
			// Under the Hood:
			// 1. Check if len(entries) < cap(entries)
			// 2. If yes: entries[len] = entry, len++
			// 3. If no: Allocate new array (cap*2), copy old elements, append new one
			//
			// Speed Comparison:
			// - C: realloc() does similar doubling, manual bounds checking required
			// - Rust: Vec::push() identical algorithm with compile-time bounds checking
			// - Python: list.append() same algorithm but higher per-element overhead
			entries = append(entries, entry)
		}
	}

	// Check for scanner errors (I/O failures like permission denied, disk error).
	// scanner.Err() returns nil if we reached EOF normally.
	// This distinguishes between normal EOF and error conditions.
	//
	// Why Check This?
	// scanner.Scan() returns false for both EOF and errors. We need to distinguish:
	// - EOF: Normal completion, proceed to sorting
	// - Error: I/O failure, return error immediately (don't return partial results)
	if err := scanner.Err(); err != nil {
		// Wrap error with context using %w verb for error chain.
		// The original error is preserved in the error chain for debugging.
		return nil, fmt.Errorf("reading input: %w", err)
	}

	// Sort entries by timestamp in ascending order (oldest first).
	// sort.Slice performs in-place quicksort with O(n log n) time complexity.
	//
	// How sort.Slice Works:
	// 1. Takes a slice and a "less" function (comparator)
	// 2. Uses quicksort algorithm (Hoare partition, 3-way partitioning for duplicates)
	// 3. Calls less function O(n log n) times to compare elements
	// 4. Swaps elements in-place (no extra memory allocation)
	//
	// The comparator is a closure that captures i and j from the sort function.
	// Memory: The closure itself is small (~16-24 bytes), allocated on stack or heap.
	//
	// Performance:
	// - Time: O(n log n) comparisons
	// - Space: O(log n) for recursion stack (quicksort is divide-and-conquer)
	// - Stable sort: Not stable (equal elements may swap). Use sort.SliceStable if needed.
	//
	// Speed Comparison:
	// - C: qsort() with function pointer, similar performance but type-unsafe (void*)
	// - Rust: slice.sort_by() with closure, identical algorithm, zero-cost abstraction
	// - Python: list.sort() uses Timsort (O(n log n), stable), slower due to Python overhead
	// - Solidity: No native sort, manual implementation required (gas-expensive for large arrays)
	sort.Slice(entries, func(i, j int) bool {
		// Compare timestamps using time.Time.Before method.
		// This handles timezone-aware comparison correctly (converts to UTC internally).
		//
		// Why time.Time is Superior:
		// - C: time_t (seconds since epoch), no timezone support, error-prone
		// - Rust: chrono::DateTime is similar to Go, proper timezone handling
		// - Python: datetime.datetime comparison is equivalent
		// - Solidity: block.timestamp (uint256 seconds), no timezone concept (always UTC)
		//
		// The Before method compares both wall clock and monotonic clock components.
		// This prevents issues with system time adjustments (NTP corrections, DST).
		return entries[i].TS.Before(entries[j].TS)
	})

	// Return results with potential error for skipped lines (partial success pattern).
	// This allows caller to:
	// 1. Use the valid entries we successfully parsed
	// 2. Be aware that some lines were skipped (data quality issue)
	//
	// Go Error Handling Philosophy:
	// - Explicit: Errors are values, returned explicitly (not thrown)
	// - Caller decides: Caller chooses to check error or ignore it
	// - No hidden control flow: Unlike exceptions, errors don't unwind the stack
	//
	// Speed Comparison:
	// - C: Return error codes, manual checking, similar performance
	// - Rust: Result<T, E> enum, explicit handling with ? operator, zero overhead
	// - Python: Exceptions have significant overhead (stack unwinding, frame inspection)
	// - Solidity: revert() unwinds transaction, all state changes rolled back (expensive)
	var err error
	if skippedCount > 0 {
		// Create error with skip count.
		// Memory: fmt.Errorf allocates error struct on heap (typically ~32 bytes).
		// This happens only if we skipped lines, not on the happy path.
		err = fmt.Errorf("skipped %d malformed lines", skippedCount)
	}

	// Return both results and error (partial success).
	// Memory: entries slice is returned (pointer to backing array + len + cap).
	// The backing array remains allocated on heap, managed by garbage collector.
	// Caller receives ownership of the slice (Go has no explicit ownership model).
	return entries, err
}

/*
Key Performance Insights:

1. JSON Parsing Overhead:
   - Go's reflection-based parsing is ~2-3x slower than Rust's compile-time codegen
   - But vastly more convenient than C's manual parsing (which is error-prone)
   - Python's json.loads() is similar speed (implemented in C) but higher memory overhead

2. Memory Allocations:
   - Scanner buffer: ~4KB, reused across lines (efficient)
   - Slice growth: Amortized O(1) via doubling strategy
   - JSON parsing: Intermediate allocations (garbage collected)
   - Total memory: O(n) where n = number of kept entries

3. Comparison Operations:
   - Integer enum comparison (Level): Single CPU instruction (~1 cycle)
   - String comparison would be O(length), much slower
   - time.Time comparison is efficient (64-bit integer comparison internally)

4. Sorting Performance:
   - O(n log n) quicksort with good cache locality (in-place)
   - Comparator function calls have tiny overhead (inlined by compiler)
   - Much faster than Python's Timsort for primitive types
   - Rust's slice.sort() is equivalent performance

5. Error Handling Cost:
   - Go's explicit error returns have zero overhead on happy path
   - Python's exceptions have significant cost when raised (stack unwinding)
   - Rust's Result enum is zero-cost abstraction (optimized to register return)
   - C's error codes are equivalent to Go's approach

Trade-offs Summary:
- Go prioritizes developer productivity (reflection, garbage collection) over raw speed
- Rust prioritizes zero-cost abstractions with compile-time guarantees
- C prioritizes raw speed but requires manual memory management
- Python prioritizes ease of use but sacrifices performance
- For this log parsing task, Go provides the best balance of speed, safety, and ergonomics
*/