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

// These are the log levels, created using iota.
const (
	Debug Level = iota // 0
	Info               // 1
	Warn               // 2
	Error              // 3
)

// Backward-compatible constant names used by cmd/app.
const (
	LevelDebug = Debug
	LevelInfo  = Info
	LevelWarn  = Warn
	LevelError = Error
)

// SortField defines the field to sort by.
type SortField int

const (
	SortByTimestamp SortField = iota
	SortByLevel
)

type Entry struct {
	TS    time.Time `json:"ts"`
	Level Level     `json:"level"`
	Msg   string    `json:"msg"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Level type.
// This method is a pointer receiver because it needs to modify the Level value.
func (l *Level) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch strings.ToLower(raw) {
	case "debug":
		*l = Debug
	case "info":
		*l = Info
	case "warn":
		*l = Warn
	case "error":
		*l = Error
	default:
		return fmt.Errorf("unknown level %q", raw)
	}
	return nil
}

// FilterAndSort is the main function you need to implement.
func FilterAndSort(r io.Reader, minLevel Level, sortBy SortField) ([]Entry, error) {
	scanner := bufio.NewScanner(r)
	entries := make([]Entry, 0)
	skipped := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			skipped++
			continue
		}
		if entry.Level >= minLevel {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		switch sortBy {
		case SortByLevel:
			if entries[i].Level == entries[j].Level {
				return entries[i].TS.Before(entries[j].TS)
			}
			return entries[i].Level < entries[j].Level
		case SortByTimestamp:
			fallthrough
		default:
			return entries[i].TS.Before(entries[j].TS)
		}
	})

	if skipped > 0 {
		return entries, fmt.Errorf("skipped %d malformed log lines", skipped)
	}
	return entries, nil
}

// FilterLogs preserves the older API used by tests and cmd/app.
func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	return FilterAndSort(r, minLevel, SortByTimestamp)
}
