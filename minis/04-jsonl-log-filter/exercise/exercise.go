package exercise

import (
	"io"
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
	// TODO: implement
	return nil, nil
}
