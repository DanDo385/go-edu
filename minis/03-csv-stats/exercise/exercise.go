package exercise

import "io"

// Stat holds aggregated statistics for a category.
type Stat struct {
	Count int     // Number of transactions
	Sum   float64 // Total amount
	Avg   float64 // Average amount (Sum / Count)
}

// SummarizeCSV reads a CSV with headers (id,category,amount) from r and returns
// per-category statistics.
//
// CSV format:
//   id,category,amount
//   1,groceries,12.50
//   2,books,10.00
//
// Returns:
//   map[string]Stat{
//     "groceries": {Count: 1, Sum: 12.50, Avg: 12.50},
//     "books": {Count: 1, Sum: 10.00, Avg: 10.00},
//   }
//
// Errors:
//   - Malformed CSV (wrong number of columns)
//   - Invalid amount (not a number)
//   - Missing headers
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	// TODO: implement
	return nil, nil
}
