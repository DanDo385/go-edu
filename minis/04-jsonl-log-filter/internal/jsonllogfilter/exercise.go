//go:build !solution && !reference

package jsonllogfilter

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

type Entry struct {
	TS    time.Time `json:"ts"`    // RFC3339 timestamp, parsed by time package
	Level Level     `json:"level"` // Uses our custom UnmarshalJSON below
	Msg   string    `json:"msg"`   // String data lives on heap, referenced by 16-byte header
}

// UnmarshalJSON implements the exercise.
//
// TODO: Implement this function
func (l *Level) UnmarshalJSON(data []byte) error {
	// TODO: Implement
	return nil
}

// FilterLogs implements the exercise.
//
// TODO: Implement this function
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	// TODO: Implement
	return nil, nil
}
