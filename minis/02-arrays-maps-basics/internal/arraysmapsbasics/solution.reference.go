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

This file is the canonical reference for this exercise. It keeps error paths
explicit when an operation can fail, so callers decide how to handle failure.

Read this alongside exercise.go and the tests to understand the intended data
flow, boundary checks, and invariants.
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
