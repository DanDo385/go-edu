//go:build !solution
// +build !solution

package exercise

import (
	"encoding/csv"
	"io"
	"strconv"
)

/*
Core Logic Without Error Handling
==================================

This file demonstrates the core CSV statistics algorithm without
error handling complexity. It's designed to help you understand the
fundamental logic flow.

KEY CONCEPTS:
- Parsing structured data (CSV records)
- Aggregating values into statistics
- Type conversions and arithmetic operations
- Working with maps of structs (value types)
*/

// CoreSummarizeCSV demonstrates CSV statistics without error handling
//
// Algorithm steps:
// 1. Create CSV reader for parsing records
// 2. Skip header row (assume it's valid)
// 3. Create statistics map to accumulate results
// 4. For each data row:
//    a. Extract category and amount fields
//    b. Parse amount string to float value
//    c. Update statistics for that category
// 5. Calculate average for each category
// 6. Return the accumulated statistics
func CoreSummarizeCSV(r io.Reader) map[string]Stat {
	// TODO: Step 1 - Create CSV reader

	// TODO: Step 2 - Skip header row

	// TODO: Step 3 - Create statistics map

	// TODO: Step 4 - Loop through data rows

	// TODO: Step 5 - Extract and parse category and amount

	// TODO: Step 6 - Update statistics (read-modify-write pattern)

	// TODO: Step 7 - Calculate averages for all categories

	// TODO: Step 8 - Return statistics map

	return nil
}
