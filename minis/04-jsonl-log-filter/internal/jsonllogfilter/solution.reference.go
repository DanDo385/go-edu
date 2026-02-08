//go:build reference

package jsonllogfilter

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
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
	Debug Level = iota
	Info
	Warn
	Error
)

const (
	LevelDebug = Debug
	LevelInfo  = Info
	LevelWarn  = Warn
	LevelError = Error
)

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

// UnmarshalJSON implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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

// FilterAndSort implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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
		default:
			return entries[i].TS.Before(entries[j].TS)
		}
	})

	if skipped > 0 {
		return entries, fmt.Errorf("skipped %d malformed log lines", skipped)
	}
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
