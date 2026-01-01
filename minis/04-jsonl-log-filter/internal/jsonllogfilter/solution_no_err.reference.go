//go:build reference

package jsonllogfilter

import (
	"bufio"
	"encoding/json"
	"io"
	"sort"
	"strings"
)

/*
Core Logic Without Error Handling
==================================

This file demonstrates the core JSONL log filtering algorithm without
error handling complexity. It's designed to help you understand the
fundamental logic flow.

KEY CONCEPTS:
- Parsing JSON data from text
- Filtering data based on criteria
- Sorting collections by keys
- Working with custom types and enums
*/

// CoreFilterLogs demonstrates log filtering without error handling
//
// Algorithm steps:
// 1. Create a slice to accumulate filtered entries
// 2. Create scanner for reading lines
// 3. For each line, parse JSON into Entry struct
// 4. Filter entries by minimum severity level
// 5. Append matching entries to results
// 6. Sort entries by timestamp
// 7. Return sorted filtered entries
func CoreFilterLogs(r io.Reader, minLevel Level) []Entry {
	// TODO: Step 1 - Create entries slice

	// TODO: Step 2 - Create scanner for reading

	// TODO: Step 3 - Loop through each line

	// TODO: Step 4 - Skip empty lines

	// TODO: Step 5 - Parse JSON into Entry

	// TODO: Step 6 - Filter by minimum level

	// TODO: Step 7 - Append matching entries

	// TODO: Step 8 - Sort entries by timestamp

	// TODO: Step 9 - Return filtered and sorted entries

	return nil
}
