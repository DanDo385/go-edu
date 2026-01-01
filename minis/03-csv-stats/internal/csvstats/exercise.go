//go:build !solution && !reference

package csvstats

/*
Problem: Compute per-category statistics from a CSV of financial transactions
Constraints:
- CSV has a header row that must be validated
- Amounts are decimal numbers (use float64)
- Missing or invalid amounts should cause an error (fail-fast)
- Empty categories should be treated as an error
Time/Space Complexity:
- Time: O(n) where n = number of rows (single pass)
- Space: O(c) where c = number of unique categories (map storage)
*/

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type Stat struct {
	Count int
	Sum   float64
	Avg   float64
}

// SummarizeCSV - TODO: implement this function
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

