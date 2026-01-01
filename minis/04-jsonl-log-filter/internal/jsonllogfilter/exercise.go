//go:build !solution && !reference

package jsonllogfilter

/*
Problem: Parse and filter JSONL (JSON Lines) log entries by severity level
Constraints:
- JSONL format: one JSON object per line (not a JSON array!)
- Timestamps are RFC3339 format (e.g., "2024-01-01T12:00:00Z")
- Level is a string ("debug", "info", "warn", "error") that must map to an enum
- Malformed lines should be skipped, not cause total failure
Time/Space Complexity:
- Time: O(n log n) where n = number of valid entries (O(n) parse + O(n log n) sort)
- Space: O(n) to store filtered entries
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

type Level int

const (
type Entry struct {
	TS    time.Time `json:"ts"`    // RFC3339 timestamp, parsed by time package
	Level Level     `json:"level"` // Uses our custom UnmarshalJSON below
	Msg   string    `json:"msg"`   // String data lives on heap, referenced by 16-byte header
}

// UnmarshalJSON - TODO: implement this function
func (l *Level) UnmarshalJSON(data []byte) error {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// FilterLogs - TODO: implement this function
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

