//go:build reference

package arraysmapsbasics

import (
	"bufio"
	"io"
	"strings"
)

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

// FreqFromReader reads words line-by-line, normalizes them, and counts frequency.
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	freq, mostCommon := CoreFreqFromReader(r)
	return freq, mostCommon, nil
}

// CoreFreqFromReader demonstrates frequency counting with explicit error handling where failures can occur
//
// Algorithm steps:
// 1. Create frequency map to track word counts
// 2. Create scanner to read lines from reader
// 3. For each line, normalize it (trim, lowercase)
// 4. Increment count for each word in the frequency map
// 5. Find the most common word by iterating the map
// 6. Return both the frequency map and the most common word
func CoreFreqFromReader(r io.Reader) (map[string]int, string) {
	freq := make(map[string]int)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if word == "" {
			continue
		}
		freq[word]++
	}

	mostCommon := ""
	maxCount := 0
	for word, count := range freq {
		if count > maxCount || (count == maxCount && (mostCommon == "" || word < mostCommon)) {
			mostCommon = word
			maxCount = count
		}
	}

	return freq, mostCommon
}
