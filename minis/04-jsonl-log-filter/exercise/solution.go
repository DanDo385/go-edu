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
- `encoding/json`: Robust JSON parser with struct mapping
- Custom unmarshalers: `UnmarshalJSON` interface for enum-like types
- `time.Time`: First-class time support with zone awareness
- `sort.Slice`: Inline sorting with custom comparators
- Error accumulation: Handle partial failures gracefully

Compared to other languages:
- Python: `json.loads()` + `datetime.fromisoformat()` is similar
  Cons: Exception-based error handling is less explicit
- JavaScript: `JSON.parse()` but Date parsing is less robust
  Cons: Date timezone handling is notoriously tricky
- Rust: `serde_json` is excellent but more boilerplate for custom types
- Java: Jackson/Gson require annotations; more verbose
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

// Level represents log severity.
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

// Entry represents a single log entry.
type Entry struct {
	TS    time.Time `json:"ts"`
	Level Level     `json:"level"`
	Msg   string    `json:"msg"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Level.
//
// Go Concepts Demonstrated:
// - Implementing the json.Unmarshaler interface
// - Converting string to enum-like values
// - Error handling in custom unmarshalers
//
// Why custom unmarshaling?
// JSON has no native enum type, so levels come as strings: "debug", "info", etc.
// We want to store them as integers (Level type) for efficient comparison (>= checks).
//
// The json package will call this method automatically when unmarshaling into a Level field.
func (l *Level) UnmarshalJSON(data []byte) error {
	// JSON strings are quoted: "info" comes as `"info"` (with quotes)
	// We need to unquote it first
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("level must be a string: %w", err)
	}

	// Convert string to Level enum
	// We use lowercase for case-insensitive matching
	switch strings.ToLower(s) {
	case "debug":
		*l = Debug
	case "info":
		*l = Info
	case "warn", "warning": // Accept both "warn" and "warning"
		*l = Warn
	case "error":
		*l = Error
	default:
		return fmt.Errorf("invalid level: %q (must be debug/info/warn/error)", s)
	}

	return nil
}

// FilterLogs reads JSONL from r, filters entries >= minLevel, and sorts by timestamp.
//
// Go Concepts Demonstrated:
// - bufio.Scanner for line-by-line reading
// - json.Unmarshal for parsing JSON strings
// - Custom unmarshaling (Level.UnmarshalJSON)
// - time.Time for timestamp handling
// - sort.Slice with inline comparator
// - Error accumulation (track skipped lines)
//
// Three-Input Iteration Table:
//
// Input 1: Valid JSONL (happy path)
//   Line 1: {"ts":"...","level":"info","msg":"A"} → filtered (< warn)
//   Line 2: {"ts":"...","level":"error","msg":"B"} → kept
//   Line 3: {"ts":"...","level":"warn","msg":"C"} → kept
//   After sort → [error, warn] (sorted by timestamp)
//
// Input 2: Empty input (edge case)
//   No lines
//   Result: [], nil
//
// Input 3: Malformed JSON (partial failure)
//   Line 1: {"ts":"...","level":"info","msg":"A"} → skipped (< warn)
//   Line 2: {invalid json → skipped, skippedCount++
//   Line 3: {"ts":"...","level":"error","msg":"C"} → kept
//   Result: [error], error("skipped 2 lines")
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	var entries []Entry
	var skippedCount int

	// Create a line-by-line scanner
	// This is memory-efficient for large log files (streaming approach)
	scanner := bufio.NewScanner(r)

	// Track line number for error reporting
	lineNum := 0

	// Read each line
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		// This handles blank lines in the JSONL file gracefully
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse the JSON line into an Entry struct
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Malformed JSON: skip this line but continue processing
			// This is a "fail-soft" approach: we don't abort on one bad line
			skippedCount++
			continue
		}

		// Filter by level
		// Since Level is an integer enum (debug=0, info=1, warn=2, error=3),
		// we can use simple integer comparison
		if entry.Level >= minLevel {
			entries = append(entries, entry)
		}
	}

	// Check for scanner errors (I/O failures)
	// scanner.Err() returns nil if we reached EOF normally
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	// Sort entries by timestamp (oldest first)
	// sort.Slice takes a comparator function (less function)
	// We use time.Before() to compare timestamps
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].TS.Before(entries[j].TS)
	})

	// If we skipped any lines, return an error (but still return valid entries)
	// This is a "partial success" pattern: caller gets both results and an error
	var err error
	if skippedCount > 0 {
		err = fmt.Errorf("skipped %d malformed lines", skippedCount)
	}

	return entries, err
}

/*
Alternatives & Trade-offs:

1. Fail-fast on malformed JSON:
   if err := json.Unmarshal(...); err != nil {
       return nil, fmt.Errorf("line %d: %w", lineNum, err)
   }
   Pros: Catches data quality issues early
   Cons: One bad line aborts entire process; less resilient

2. Use json.Decoder instead of Unmarshal:
   dec := json.NewDecoder(strings.NewReader(line))
   dec.Decode(&entry)
   Pros: Slightly more efficient (avoids []byte allocation)
   Cons: Same result for simple cases; decoder is better for streaming arrays

3. Store Level as string instead of enum:
   type Entry struct { Level string `json:"level"` }
   Pros: Simpler (no custom unmarshaling)
   Cons: String comparisons are slower and error-prone
         entry.Level == "warn" (typo-prone) vs entry.Level >= Warn

4. Use a logging library (log/slog):
   Go 1.21+ has structured logging built-in
   Pros: Standardized format, better performance
   Cons: Different API; this exercise teaches JSON parsing fundamentals

5. Return detailed error information:
   type ParseError struct { Line int; Content string; Err error }
   var errs []ParseError
   Pros: Caller can inspect each failure
   Cons: More complex API; usually simple count is sufficient

Go vs X:

Go vs Python:
  import json
  from datetime import datetime
  entries = []
  for line in f:
      try:
          entry = json.loads(line)
          if level_map[entry['level']] >= min_level:
              entries.append(entry)
      except json.JSONDecodeError:
          skipped += 1
  entries.sort(key=lambda e: e['ts'])
  Pros: Similar logic, fewer lines
  Cons: Exception-based error handling is less explicit
        Dict access (entry['level']) can raise KeyError (no compile-time safety)
  Go: Static typing catches missing fields at compile time

Go vs JavaScript (Node.js):
  const entries = [];
  for (const line of lines) {
      try {
          const entry = JSON.parse(line);
          if (levelMap[entry.level] >= minLevel) {
              entries.push(entry);
          }
      } catch (e) { skipped++; }
  }
  entries.sort((a, b) => new Date(a.ts) - new Date(b.ts));
  Pros: Concise, familiar syntax
  Cons: Date parsing is inconsistent across browsers/Node versions
        No type safety (TypeScript helps but adds complexity)
  Go: time.Time is always RFC3339-aware; no surprises

Go vs Rust:
  use serde::{Deserialize, Serialize};
  #[derive(Deserialize, PartialEq, PartialOrd)]
  enum Level { Debug, Info, Warn, Error }
  let entries: Vec<Entry> = lines
      .filter_map(|line| serde_json::from_str(&line).ok())
      .filter(|e| e.level >= min_level)
      .collect();
  entries.sort_by(|a, b| a.ts.cmp(&b.ts));
  Pros: Zero-cost abstractions; faster execution
  Cons: More boilerplate (derive macros, type annotations)
        serde requires learning a new API
  Go: Simpler reflection-based JSON; easier for beginners

Go vs Java:
  ObjectMapper mapper = new ObjectMapper();
  List<Entry> entries = new ArrayList<>();
  for (String line : lines) {
      try {
          Entry entry = mapper.readValue(line, Entry.class);
          if (entry.getLevel().compareTo(minLevel) >= 0) {
              entries.add(entry);
          }
      } catch (JsonProcessingException e) { skipped++; }
  }
  entries.sort(Comparator.comparing(Entry::getTs));
  Pros: Similar structure; Jackson is robust
  Cons: Much more verbose (getters, class definitions, generics)
        Requires external library (Jackson/Gson)
  Go: Built-in JSON; cleaner syntax with struct tags
*/
