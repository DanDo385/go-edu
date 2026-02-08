//go:build reference

package csvstats

import (
	"encoding/csv"
	"io"
	"strconv"
)

type Stat struct {
	Count int
	Sum   float64
	Avg   float64
}

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

// SummarizeCSV validates input and returns per-category count/sum/avg statistics.
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	cr := csv.NewReader(r)
	header, err := cr.Read()
	if err != nil {
		return nil, err
	}
	if len(header) != 3 || header[0] != "id" || header[1] != "category" || header[2] != "amount" {
		return nil, io.ErrUnexpectedEOF
	}

	stats := make(map[string]Stat)
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) != 3 || record[1] == "" {
			return nil, io.ErrUnexpectedEOF
		}

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		cur := stats[record[1]]
		cur.Count++
		cur.Sum += amount
		stats[record[1]] = cur
	}

	for k, v := range stats {
		v.Avg = v.Sum / float64(v.Count)
		stats[k] = v
	}
	return stats, nil
}

// CoreSummarizeCSV demonstrates CSV statistics with explicit error handling where failures can occur
//
// Algorithm steps:
//  1. Create CSV reader for parsing records
//  2. Skip header row (assume it's valid)
//  3. Create statistics map to accumulate results
//  4. For each data row:
//     a. Extract category and amount fields
//     b. Parse amount string to float value
//     c. Update statistics for that category
//  5. Calculate average for each category
//  6. Return the accumulated statistics
func CoreSummarizeCSV(r io.Reader) map[string]Stat {
	stats, err := SummarizeCSV(r)
	if err != nil {
		return nil
	}
	return stats
}
