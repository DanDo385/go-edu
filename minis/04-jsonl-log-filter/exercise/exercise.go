//go:build !solution
// +build !solution

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

// Level represents log severity (enum-like type).
type Level int

const (
	Debug Level = iota // 0
	Info               // 1
	Warn               // 2
	Error              // 3
)

// Entry represents a single log entry.
type Entry struct {
	TS    time.Time `json:"ts"`    // Timestamp (RFC3339 format)
	Level Level     `json:"level"` // Severity level
	Msg   string    `json:"msg"`   // Log message
}

// UnmarshalJSON allows parsing a Level from a JSON string (e.g., "info").
func (l *Level) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("level must be string: %w", err)
	}
	switch strings.ToLower(s) {
	case "debug":
		*l = Debug
	case "info":
		*l = Info
	case "warn", "warning":
		*l = Warn
	case "error":
		*l = Error
	default:
		return fmt.Errorf("invalid level: %q", s)
	}
	return nil
}

// FilterLogs reads JSONL from r, filters entries >= minLevel, and sorts by timestamp.
//
// JSONL format (one JSON object per line):
//   {"ts":"2024-01-01T12:00:00Z","level":"info","msg":"Server started"}
//   {"ts":"2024-01-01T12:00:05Z","level":"error","msg":"Database failed"}
//
// Behavior:
//   - Unparseable lines are skipped (not fatal)
//   - Returns an error if any lines were skipped (includes count)
//   - Results are sorted by timestamp (oldest first)
//
// Example:
//   entries, err := FilterLogs(r, Warn)
//   // Returns only "warn" and "error" entries
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	var entries []Entry
	var skipped int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			skipped++
			continue
		}

		if e.Level >= minLevel {
			entries = append(entries, e)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read logs: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].TS.Before(entries[j].TS)
	})

	if skipped > 0 {
		return entries, fmt.Errorf("skipped %d malformed lines", skipped)
	}
	return entries, nil
}
