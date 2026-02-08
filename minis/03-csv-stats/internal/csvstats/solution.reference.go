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
Reference Solution - CSV Processing and Statistical Aggregation
============================================================

This file demonstrates structured data processing using Go's standard library.
CSV (Comma-Separated Values) is a ubiquitous data format for tabular data exchange.
This exercise shows how to parse CSV, validate data, and compute aggregate statistics.

This connects to the broader Go ecosystem by demonstrating:
- encoding/csv package for robust CSV parsing with proper escaping
- io.Reader interface for composable data sources (files, network, strings)
- strconv package for safe string-to-number conversion
- Map-based aggregation patterns common in data processing

The exercise builds understanding of:
- Streaming data processing: handling data as it arrives rather than loading everything
- Error propagation: surfacing parsing errors so callers can handle them appropriately
- Data validation: checking data integrity before processing
- Aggregation algorithms: building statistics incrementally

Teaching notes:
- Memory/ownership: the returned map is newly allocated, so callers own it completely.
  This avoids the aliasing issues that would occur if we returned internal references.
- Invariants: we validate the CSV header structure upfront, establishing assumptions
  that simplify the per-row processing logic.
- Error surfaces: CSV parsing can fail in many ways (malformed data, encoding issues,
  type conversion errors), so we surface all errors rather than making assumptions.
*/

/*
SummarizeCSV - CSV Statistical Analysis

This function reads CSV data and computes statistics grouped by category.
It demonstrates production-ready data processing with proper error handling.

The CSV format expected:
id,category,amount
1,electronics,99.99
2,books,25.50
3,electronics,149.99

Returns a map where each key is a category and the value contains:
- Count: number of items in that category
- Sum: total amount spent on that category
- Avg: average amount per item in that category

Error handling: Returns detailed errors for malformed CSV, invalid headers,
or type conversion failures. This allows callers to implement appropriate
error recovery (retry, skip bad records, alert users, etc.).
*/
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	// Step 1: Create CSV reader
	// csv.NewReader wraps any io.Reader and provides CSV parsing functionality
	// It handles proper CSV escaping, quoted fields, and different delimiters
	// This abstraction allows us to process CSV from files, network, or memory
	cr := csv.NewReader(r)

	// Step 2: Read and validate header row
	// CSV files typically have a header row that names the columns
	// cr.Read() returns the next row as a []string, or an error
	header, err := cr.Read()
	if err != nil {
		// If we can't even read the header, the file is malformed or empty
		// Return the error so the caller knows what went wrong
		return nil, err
	}

	// Validate header structure and field names
	// We expect exactly 3 columns: id, category, amount
	// len(header) checks the number of columns
	// header[0], header[1], header[2] check the specific field names
	if len(header) != 3 || header[0] != "id" || header[1] != "category" || header[2] != "amount" {
		// If the header doesn't match our expectations, return an error
		// io.ErrUnexpectedEOF is a bit misleading here, but it's a standard error
		// In production code, we'd define a custom error type for this
		return nil, io.ErrUnexpectedEOF
	}

	// Step 3: Initialize statistics map
	// We'll accumulate statistics by category in this map
	// Key: category name (string), Value: aggregated statistics (Stat struct)
	// The map starts empty and grows as we encounter new categories
	stats := make(map[string]Stat)

	// Step 4: Process data rows one by one
	// CSV processing is typically done in a loop until EOF
	// This streaming approach handles large files without loading everything into memory
	for {
		// Read the next row from the CSV
		record, err := cr.Read()

		// Check for end of file
		if err == io.EOF {
			// Normal termination - we've read all the data
			break
		}

		// Check for other errors (malformed CSV, I/O issues, etc.)
		if err != nil {
			// Surface the error so caller can handle it appropriately
			return nil, err
		}

		// Validate data row structure
		// Each data row should have exactly 3 fields
		// The category field (index 1) should not be empty
		if len(record) != 3 || record[1] == "" {
			// Malformed data row - return error
			return nil, io.ErrUnexpectedEOF
		}

		// Step 4a: Parse the amount field from string to float64
		// CSV data comes as strings, but we need numeric values for calculations
		// strconv.ParseFloat handles various number formats and decimal points
		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			// Invalid number format - surface the parsing error
			return nil, err
		}

		// Step 4b: Update statistics for this category
		// Get current statistics for this category (creates zero value if new)
		// In Go, accessing a missing map key returns the zero value for the value type
		// So cur will be Stat{Count: 0, Sum: 0.0, Avg: 0.0} for new categories
		cur := stats[record[1]]

		// Increment the count for this category
		cur.Count++

		// Add this amount to the running sum
		cur.Sum += amount

		// Store the updated statistics back in the map
		// This overwrites the previous value for this category
		stats[record[1]] = cur
	}

	// Step 5: Calculate averages for each category
	// We deferred average calculation until all data was processed
	// This avoids recalculating averages on every addition
	for k, v := range stats {
		// Calculate average: sum divided by count
		// float64(v.Count) converts int to float64 for division
		v.Avg = v.Sum / float64(v.Count)

		// Store the updated statistics (now including average) back in the map
		stats[k] = v
	}

	// Step 6: Return results
	// Return the completed statistics map and nil error (success)
	// The map contains all categories with their final count, sum, and average
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
