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
	cr := csv.NewReader(r)
	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if len(header) != 3 || header[0] != "id" || header[1] != "category" || header[2] != "amount" {
		return nil, fmt.Errorf("invalid header: %v", header)
	}

	stats := make(map[string]Stat)
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		if len(record) != 3 {
			return nil, fmt.Errorf("invalid row width: %d", len(record))
		}

		category := record[1]
		if category == "" {
			return nil, fmt.Errorf("empty category")
		}

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount %q: %w", record[2], err)
		}

		cur := stats[category]
		cur.Count++
		cur.Sum += amount
		stats[category] = cur
	}

	for category, stat := range stats {
		stat.Avg = stat.Sum / float64(stat.Count)
		stats[category] = stat
	}

	return stats, nil
}
