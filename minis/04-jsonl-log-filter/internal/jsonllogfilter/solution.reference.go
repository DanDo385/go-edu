//go:build reference

package jsonllogfilter

/*
Reference Solution - JSON Lines Log Processing and Filtering
==========================================================

This file demonstrates structured log processing using JSON Lines (JSONL) format.
JSONL stores each log entry as a separate JSON object on its own line, making
it ideal for streaming log processing, grep-style filtering, and big data tools.

This connects to the broader Go ecosystem by showing:
- encoding/json package for robust JSON parsing with custom unmarshaling
- bufio.Scanner for efficient line-by-line text processing
- time.Time for timestamp handling and sorting
- sort package for custom comparison functions
- How Go's type system enables type-safe log level enums

The exercise builds understanding of:
- Streaming data processing: handle logs as they arrive, not all at once
- Custom JSON unmarshaling: convert string log levels to typed enums
- Error resilience: continue processing when encountering malformed log entries
- Sorting algorithms: implement custom comparison logic for multi-field sorting
- Data filtering: apply business logic to select relevant log entries

Teaching notes:
- Memory/ownership: returned slices are newly allocated, so callers own them.
  This prevents aliasing issues while allowing efficient append operations.
- Invariants: JSON schema is validated during unmarshaling, establishing data
  integrity guarantees that simplify downstream processing.
- Error surfaces: log processing can encounter malformed JSON, I/O errors, and
  invalid data. We surface these explicitly so callers can implement appropriate
  error handling (alerting, retries, data repair, etc.).
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

/*
Level - Typed Log Level Enumeration

In Go, we use typed constants instead of strings for log levels to:
1. Prevent typos (compile-time checking vs runtime string comparison)
2. Enable efficient comparison operations
3. Provide IDE autocompletion and documentation

The `int` underlying type allows Level values to be compared with <, >, >= operators.
This is more efficient than string comparison for filtering operations.
*/
type Level int

/*
Log Level Constants - Ordered Severity Scale

iota generates sequential integer constants starting from 0.
This creates a natural ordering: Debug (0) < Info (1) < Warn (2) < Error (3)

Benefits of this approach:
- Constants are self-documenting and grouped together
- iota ensures they stay in sync if we add/remove levels
- The ordering enables efficient range queries (level >= Info)
*/
const (
	Debug Level = iota // 0 - Most verbose, typically disabled in production
	Info               // 1 - General information about normal operations
	Warn               // 2 - Warning about potential issues
	Error              // 3 - Error conditions that need attention
)

/*
Level Constants - Alternative Naming Convention

Some codebases prefer these explicit constant names.
They provide the same values as above but with different naming.
This demonstrates that multiple constant names can refer to the same value.

In practice, you'd choose one naming convention for consistency.
*/
const (
	LevelDebug = Debug
	LevelInfo  = Info
	LevelWarn  = Warn
	LevelError = Error
)

/*
SortField - Enumeration for Sort Criteria

Similar to Level, this uses an int-based enum for type safety and efficiency.
The sort field determines which Entry field to use as the primary sort key.
*/
type SortField int

/*
Sort Field Constants

Defines the available sorting options for log entries.
- SortByTimestamp: chronological order (most common for log analysis)
- SortByLevel: group by severity level
*/
const (
	SortByTimestamp SortField = iota // Primary: timestamp, Secondary: none
	SortByLevel                      // Primary: level, Secondary: timestamp
)

/*
Entry - Log Entry Structure

This struct represents a single log entry with structured data.
The JSON tags specify how fields map to/from JSON format.

JSON marshaling automatically handles:
- time.Time ↔ RFC3339 timestamp strings
- Custom Level enum ↔ string conversion (via UnmarshalJSON)
- string ↔ JSON string values

This structure enables type-safe log processing while maintaining
compatibility with JSON-based log aggregation systems.
*/
type Entry struct {
	TS    time.Time `json:"ts"`    // Timestamp when the log was created
	Level Level     `json:"level"` // Log severity level (Debug, Info, Warn, Error)
	Msg   string    `json:"msg"`   // The actual log message text
}

/*
UnmarshalJSON - Custom JSON Deserialization for Level

This method implements the json.Unmarshaler interface, allowing Level values
to be parsed from JSON strings like "debug", "info", "warn", "error".

Why custom unmarshaling?
- JSON has no native enum type, only strings and numbers
- We want case-insensitive parsing ("DEBUG" = "debug")
- We want to validate log levels at parse time, not later
- We want to convert strings to efficient integer comparisons

The method signature requires:
- *Level receiver: we modify the Level value being unmarshaled
- []byte parameter: raw JSON bytes for this field
- error return: parsing failures are surfaced as errors

Algorithm steps:
1. Unmarshal JSON bytes into a temporary string variable
2. Convert to lowercase for case-insensitive matching
3. Map string to Level constant using switch statement
4. Return error for unknown level strings
*/
func (l *Level) UnmarshalJSON(data []byte) error {
	// Step 1: Parse JSON bytes as string
	// json.Unmarshal can convert JSON strings to Go strings automatically
	// We use a temporary variable to avoid modifying *l until we know it's valid
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		// If the JSON isn't even a valid string, surface the parsing error
		return err
	}

	// Step 2: Case-insensitive string matching
	// strings.ToLower normalizes the input for consistent comparison
	// This allows "DEBUG", "debug", "Debug" to all map to the same Level
	switch strings.ToLower(raw) {
	case "debug":
		// Assign the Debug constant (0) to the receiver
		// *l dereferences the pointer receiver to modify the actual Level value
		*l = Debug
	case "info":
		*l = Info
	case "warn":
		*l = Warn
	case "error":
		*l = Error
	default:
		// Unknown log level - return descriptive error
		// fmt.Errorf creates a new error with formatted message
		// %q adds quotes around the string for clarity
		return fmt.Errorf("unknown level %q", raw)
	}

	// Step 3: Success - no error to return
	return nil
}

/*
FilterAndSort - Complete Log Processing Pipeline

This function implements the full log processing pipeline: parsing, filtering, sorting.
It demonstrates production-ready data processing with proper error handling and resilience.

Parameters:
- r: io.Reader providing JSONL log data (could be file, network, memory)
- minLevel: minimum log level to include (filters out less severe logs)
- sortBy: primary sort criteria (timestamp or level)

Returns:
- []Entry: filtered and sorted log entries (newly allocated slice)
- error: any processing errors encountered

Algorithm overview:
1. Parse JSONL lines into Entry structs
2. Filter entries by minimum log level
3. Sort entries by specified criteria
4. Return results with error summary

Error handling strategy:
- Malformed JSON lines are skipped (with count)
- I/O errors during reading cause immediate failure
- Invalid log levels are caught during JSON unmarshaling
- Sorting errors are prevented by type safety
*/
func FilterAndSort(r io.Reader, minLevel Level, sortBy SortField) ([]Entry, error) {
	// Step 1: Initialize processing state
	// bufio.NewScanner provides efficient line-by-line reading
	// It handles buffering automatically and works with any io.Reader
	scanner := bufio.NewScanner(r)

	// Pre-allocate slice with initial capacity of 0 (will grow as needed)
	// This slice will store successfully parsed and filtered log entries
	entries := make([]Entry, 0)

	// Track malformed lines for error reporting
	// We continue processing despite errors, but inform the caller
	skipped := 0

	// Step 2: Process each line of JSONL input
	// scanner.Scan() returns true while there are more lines to read
	for scanner.Scan() {
		// Get the current line and remove leading/trailing whitespace
		// JSONL allows whitespace around JSON objects for readability
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines (common in hand-edited log files)
		if line == "" {
			continue
		}

		// Step 2a: Parse JSON line into Entry struct
		// json.Unmarshal converts JSON bytes to Go struct
		// The Entry struct's UnmarshalJSON methods handle custom parsing
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Malformed JSON - count it and continue processing
			// In production, you might log these errors or collect them
			skipped++
			continue
		}

		// Step 2b: Apply level filtering
		// Only include entries at or above the minimum severity level
		// Level enum values allow efficient integer comparison
		if entry.Level >= minLevel {
			// Add entry to results slice - may trigger slice growth
			entries = append(entries, entry)
		}
	}

	// Step 3: Check for I/O errors during scanning
	// scanner.Err() returns any error encountered while reading
	// This is separate from JSON parsing errors above
	if err := scanner.Err(); err != nil {
		// I/O error (file corrupted, network issue, etc.) - can't continue
		return nil, err
	}

	// Step 4: Sort the filtered entries
	// sort.Slice sorts in-place using a custom comparison function
	// The comparison function determines sort order and stability
	sort.Slice(entries, func(i, j int) bool {
		// Switch on sort criteria to implement different sort orders
		switch sortBy {
		case SortByLevel:
			// Primary sort: by level (Debug < Info < Warn < Error)
			// If levels are equal, secondary sort by timestamp
			if entries[i].Level == entries[j].Level {
				// Same level - sort by timestamp (chronological)
				// time.Time.Before() provides proper time comparison
				return entries[i].TS.Before(entries[j].TS)
			}
			// Different levels - sort by level severity
			return entries[i].Level < entries[j].Level

		default: // SortByTimestamp (or any unknown value)
			// Default to chronological sorting
			// This is the most common log sorting requirement
			return entries[i].TS.Before(entries[j].TS)
		}
	})

	// Step 5: Return results with error summary
	// If we skipped any malformed lines, include that in the error
	// This allows callers to decide if the error is acceptable
	if skipped > 0 {
		return entries, fmt.Errorf("skipped %d malformed log lines", skipped)
	}

	// Success - return filtered and sorted entries
	return entries, nil
}

// FilterLogs implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	return FilterAndSort(r, minLevel, SortByTimestamp)
}
