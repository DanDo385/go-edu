//go:build !solution
// +build !solution

package csvstats

// TODO: Import required packages
// You'll need:
// - "encoding/csv" for CSV parsing
// - "io" for the Reader interface
// - "strconv" for string to float conversion
//
// Uncomment these imports when you implement the functions:
/*
// import (
// 	"encoding/csv"
// 	"io"
// 	"strconv"
// )
*/

// Stat holds aggregated statistics for a category
type Stat struct {
	Count int     // Number of transactions in this category
	Sum   float64 // Total amount for this category
	Avg   float64 // Average amount (Sum / Count)
}

// SummarizeCSV reads a CSV file from r and returns category statistics
//
// Expected CSV format:
//   id,category,amount
//   1,groceries,12.50
//   2,books,10.00
//
// Returns:
// - map[string]Stat: Statistics for each category
// - error: Any validation or parsing errors
//
// Example: CSV with "groceries,12.50" and "groceries,7.50" returns
//   {"groceries": {Count: 2, Sum: 20.0, Avg: 10.0}}
// TODO: Implement this function
//
// Algorithm:
// 1. Create a CSV reader from the input reader
// 2. Read and validate the header row (expect: id, category, amount)
// 3. Create a statistics map to store results
// 4. Loop through each data row:
//    a. Parse the amount string to float
//    b. Get or create the category's Stat entry
//    c. Increment count and add to sum
// 5. Calculate average for each category (Sum / Count)
// 6. Return the statistics map and any errors encountered
//
// CS Concepts:
// - Data Aggregation: Grouping data by key and computing aggregate functions (count, sum)
// - Parsing: Converting string representations to typed values (CSV to struct)
// - Struct Modification: Value types require read-modify-write pattern when stored in maps
// - Type Conversion: Implicit conversions (int to float) for arithmetic operations
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	// TODO: Implement SummarizeCSV
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// After implementing:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Try with invalid CSV to test error handling
// - Compare with solution.go to see detailed explanations
