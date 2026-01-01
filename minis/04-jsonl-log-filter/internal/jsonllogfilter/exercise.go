//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package jsonllogfilter

import (
	"io"

	"time"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Entry struct {
	TS    time.Time `json:"ts"`    // RFC3339 timestamp, parsed by time package
	Level Level     `json:"level"` // Uses our custom UnmarshalJSON below
	Msg   string    `json:"msg"`   // String data lives on heap, referenced by 16-byte header
}
// TODO: implement UnmarshalJSON.
func (l *Level) UnmarshalJSON(data []byte) error { panic("TODO: implement") }
// TODO: implement FilterLogs.
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) { panic("TODO: implement") }
