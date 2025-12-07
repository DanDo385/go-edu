//go:build !solution
// +build !solution

package exercise

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type Stat struct {
	Count int     // Number of transactions
	Sum   float64 // Total amount
	Avg   float64 // Average amount (Sum / Count)
}

func SummarizeCSV(r io.Reader) (map[string]Stat, error) {

	csvReader := csv.NewReader(r)

	headers, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty CSV file")
		}
		return nil, fmt.Errorf("reader header: %w", err)
	}
	
	if len(headers) != 3 || headers[0] != "id" || headers[1] != "category" || headers[2] != "amount" {
		return nil, fmt.Errorf("invalid header: expected [id,category,amount], got %v", headers)
	}

	stats := make(map[string]Stat)
	rowNum := 2

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}
		if len(record) !=3 {
			return nil, fmt.Errorf("row %d: expected 3 fields, got %d", rowNum, len(record))
		}

		category := record[1]
		amountStr := record[2]

		if category == "" {
			return nil, fmt.Errorf("row %d: empty category", rowNum)
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid amount %q: %w", rowNum, amountStr, err)
		}

		s := stats[category]
		s.Count++
		s.Sum += amount
		stats[category] = s
		rowNum++
	}

	for category, s := range stats {
		if s.Count > 0 {
			s.Avg = s.Sum / float64(s.Count)
			stats[category] = s
		}
	}

	return stats, nil

}