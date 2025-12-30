//go:build solution
// +build solution

/*
Problem: Compute per-category statistics from a CSV of financial transactions

Given a CSV with columns (id, category, amount), we need to:
1. Parse the CSV line-by-line (streaming for memory efficiency)
2. Group transactions by category
3. Compute count, sum, and average for each category
4. Handle malformed data gracefully

Constraints:
- CSV has a header row that must be validated
- Amounts are decimal numbers (use float64)
- Missing or invalid amounts should cause an error (fail-fast)
- Empty categories should be treated as an error

Time/Space Complexity:
- Time: O(n) where n = number of rows (single pass)
- Space: O(c) where c = number of unique categories (map storage)

Why Go is well-suited:
- `encoding/csv` in stdlib: No external dependencies for CSV parsing
- Streaming I/O: Process line-by-line for constant memory usage
- Strong typing: Compile-time detection of struct field mismatches
- Explicit error handling: No silent data corruption

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical CSV parsing points
2. Watching struct field transformations in the Variables panel
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting map aggregation patterns
5. Using the Debug Console to evaluate expressions
6. Understanding streaming I/O and memory usage
*/

package exercise

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// Stat holds aggregated statistics for a category.
type Stat struct {
	Count int
	Sum   float64
	Avg   float64
}

// SummarizeCSV reads a CSV with headers (id,category,amount) and returns
// per-category statistics.
//
// Go Concepts Demonstrated:
// - encoding/csv: Standard library CSV parser with streaming support
// - Structs: Define custom data types with named fields
// - Error wrapping: Use fmt.Errorf with %w to preserve error chains
// - strconv: Convert strings to numbers with error handling
// - Map aggregation: Group-by pattern using maps
//
// Three-Input Iteration Table:
//
// Input 1: Valid CSV (happy path)
//   Row 1: "1,groceries,12.50" → groceries: {Count:1, Sum:12.50}
//   Row 2: "2,groceries,7.50"  → groceries: {Count:2, Sum:20.00}
//   Row 3: "3,books,10.00"     → books: {Count:1, Sum:10.00}
//   Post-process → groceries: {Count:2, Sum:20.00, Avg:10.00}, books: {Count:1, Sum:10.00, Avg:10.00}
//
// Input 2: Empty CSV (edge case)
//   Only header row: "id,category,amount"
//   No data rows
//   Result: empty map (valid)
//
// Input 3: Malformed amount (failure case)
//   Row 1: "1,groceries,12.50" → groceries: {Count:1, Sum:12.50}
//   Row 2: "2,books,invalid"   → Error: "row 3: invalid amount"
//   Result: nil, error (fail-fast)
//
// DEBUGGING WORKFLOW:
// ===================
// 1. Set breakpoint at function entry (line 100)
// 2. Call with test CSV containing multiple categories
// 3. Step through CSV parsing and aggregation
// 4. Watch stats map build category by category
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	// BREAKPOINT 1: Set here to inspect function entry
	// DEBUG: In Variables panel, expand 'r' to see:
	//   - The io.Reader interface value
	//   - For strings.Reader, you can see the underlying CSV data
	// DEBUG: In Debug Console, type: r
	// This shows what type of reader was passed in
	// ============================================================================
	// CSV READER: Struct with internal state
	// ============================================================================
	// csv.Reader automatically handles:
	// - Line breaks within quoted fields
	// - Escaped quotes (double quotes: "He said ""hello""")
	// - Different line endings (\n vs \r\n)
	//
	// Memory semantics:
	// - csv.NewReader returns a csv.Reader STRUCT (not a pointer!)
	// - But the struct contains pointers to: buffer, io.Reader, etc.
	// - So csvReader is a value type, but shares internal state
	//
	// BREAKPOINT 2: Set here BEFORE csv.NewReader call
	// BREAKPOINT 3: Set here AFTER csv.NewReader call
	// DEBUG: In Variables panel, expand 'csvReader' to see:
	//   - Internal buffer state
	//   - Reader configuration
	// DEBUG: In Debug Console, type: csvReader
	csvReader := csv.NewReader(r)

	// ============================================================================
	// HEADER READING: Validate CSV structure
	// ============================================================================
	// Read the header row
	// This validates the CSV structure and allows us to check column names
	//
	// BREAKPOINT 4: Set here BEFORE reading headers
	// Step Into (F11) to see csv.Reader.Read() implementation
	// BREAKPOINT 5: Set here AFTER reading headers
	// DEBUG: In Variables panel, expand 'headers' to see:
	//   - headers.len: should be 3
	//   - headers[0], headers[1], headers[2]: column names
	// DEBUG: In Debug Console, type: headers
	// DEBUG: In Debug Console, type: len(headers)
	headers, err := csvReader.Read()
	if err != nil {
		// BREAKPOINT 6: Set here in error path
		// DEBUG: In Variables panel, expand 'err'
		// DEBUG: In Debug Console, type: err == io.EOF
		// Distinguish between empty file vs. I/O error
		if err == io.EOF {
			return nil, fmt.Errorf("empty CSV file (no header)")
		}
		return nil, fmt.Errorf("reading header: %w", err)
	}

	// ====================================================================
	// HEADER VALIDATION: Schema checking
	// ====================================================================
	// Validate header format
	// We expect exactly 3 columns: id, category, amount
	// This is a defensive check to fail fast on schema mismatches
	//
	// BREAKPOINT 7: Set here BEFORE validation
	// DEBUG: In Debug Console, type: len(headers) != 3
	// DEBUG: In Debug Console, type: headers[0] != "id"
	// DEBUG: In Debug Console, type: headers[1] != "category"
	// DEBUG: In Debug Console, type: headers[2] != "amount"
	// See which condition (if any) will trigger the error
	if len(headers) != 3 || headers[0] != "id" || headers[1] != "category" || headers[2] != "amount" {
		// BREAKPOINT 8: Set here if validation fails
		return nil, fmt.Errorf("invalid header: expected [id,category,amount], got %v", headers)
	}

	// ============================================================================
	// STATS MAP: Initialize aggregation structure
	// ============================================================================
	// Initialize the aggregation map
	// Key: category name
	// Value: running totals (count and sum; average computed later)
	//
	// BREAKPOINT 9: Set here BEFORE map creation
	// BREAKPOINT 10: Set here AFTER map creation
	// DEBUG: In Variables panel, expand 'stats' to see:
	//   - Empty map: map[]
	// DEBUG: In Debug Console, type: len(stats)
	// Should be 0 initially
	stats := make(map[string]Stat)

	// Track row number for error messages (starting at 2 since row 1 is header)
	// BREAKPOINT 11: Set here AFTER rowNum initialization
	// DEBUG: In Variables panel, watch 'rowNum'
	// It should be 2 (row 1 was the header)
	rowNum := 2

	// ============================================================================
	// MAIN LOOP: Read and process each CSV row
	// ============================================================================
	// Read records line-by-line (streaming)
	// csv.Reader.Read() returns []string for each row
	// It returns io.EOF when the file is exhausted (not an error!)
	//
	// BREAKPOINT 12: Set here at loop start
	// Step Over (F10) to iterate through each row
	// Watch stats map grow with each category
	for {
		// BREAKPOINT 13: Set here at start of loop body
		// DEBUG: In Variables panel, watch:
		//   - rowNum: current row being processed
		//   - stats: see it grow with each iteration
		// DEBUG: In Debug Console, type: rowNum
		// DEBUG: In Debug Console, type: len(stats)

		// ====================================================================
		// ROW READING: Get next CSV record
		// ====================================================================
		// BREAKPOINT 14: Set here BEFORE reading record
		// Step Into (F11) to see csv.Reader.Read() parse the CSV
		// BREAKPOINT 15: Set here AFTER reading record
		// DEBUG: In Variables panel, expand 'record' to see:
		//   - record.len: should be 3 (id, category, amount)
		//   - record[0], record[1], record[2]: field values
		// DEBUG: In Debug Console, type: record
		record, err := csvReader.Read()
		if err == io.EOF {
			// BREAKPOINT 16: Set here at EOF
			// This is the normal loop exit condition
			// DEBUG: In Variables panel, review 'stats' map
			// All rows should be processed at this point
			// End of file is the normal exit condition
			break
		}
		if err != nil {
			// BREAKPOINT 17: Set here at error
			// DEBUG: In Variables panel, expand 'err'
			// DEBUG: In Debug Console, type: err.Error()
			// Unexpected error (e.g., malformed CSV, I/O failure)
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}

		// ====================================================================
		// FIELD COUNT VALIDATION: Ensure correct number of columns
		// ====================================================================
		// Expect exactly 3 fields per row
		// csv.Reader ensures this by default but we check defensively
		//
		// BREAKPOINT 18: Set here BEFORE field count check
		// DEBUG: In Debug Console, type: len(record)
		// DEBUG: In Debug Console, type: len(record) != 3
		if len(record) != 3 {
			// BREAKPOINT 19: Set here if wrong field count
			return nil, fmt.Errorf("row %d: expected 3 fields, got %d", rowNum, len(record))
		}

		// ====================================================================
		// FIELD EXTRACTION: Get id, category, amount
		// ====================================================================
		// Extract fields
		// We don't use the id field in this analysis, but could validate it's numeric
		//
		// BREAKPOINT 20: Set here BEFORE field extraction
		// BREAKPOINT 21: Set here AFTER field extraction
		// DEBUG: In Variables panel, watch:
		//   - category: extracted category name
		//   - amountStr: amount as string (needs parsing)
		// DEBUG: In Debug Console, type: category
		// DEBUG: In Debug Console, type: amountStr
		// id := record[0]  // unused
		category := record[1]
		amountStr := record[2]

		// ====================================================================
		// CATEGORY VALIDATION: Ensure not empty
		// ====================================================================
		// Validate category is not empty
		//
		// BREAKPOINT 22: Set here BEFORE empty check
		// DEBUG: In Debug Console, type: category == ""
		if category == "" {
			// BREAKPOINT 23: Set here if category empty
			return nil, fmt.Errorf("row %d: empty category", rowNum)
		}

		// ====================================================================
		// AMOUNT PARSING: Convert string to float64
		// ====================================================================
		// Parse amount as float64
		// strconv.ParseFloat returns an error if the string is not a valid number
		// The second argument (64) specifies float64 precision
		//
		// BREAKPOINT 24: Set here BEFORE parsing amount
		// Step Into (F11) to see strconv.ParseFloat implementation
		// BREAKPOINT 25: Set here AFTER parsing amount
		// DEBUG: In Variables panel, watch:
		//   - amountStr: original string value
		//   - amount: parsed float64 value
		// DEBUG: In Debug Console, type: amountStr
		// DEBUG: In Debug Console, type: amount
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			// BREAKPOINT 26: Set here if parsing fails
			// DEBUG: In Variables panel, expand 'err'
			// DEBUG: In Debug Console, type: amountStr
			// See what invalid amount caused the error
			return nil, fmt.Errorf("row %d: invalid amount %q: %w", rowNum, amountStr, err)
		}

		// ====================================================================
		// MAP UPDATE: Read-Modify-Write pattern for struct values
		// ====================================================================
		// Map lookup returns the zero value (Stat{Count:0, Sum:0.0}) if key doesn't exist
		// This is perfect for aggregation: we can read, modify, and write back
		//
		// CRITICAL CONCEPT: Maps store VALUES, not pointers
		// - stats[category] returns a COPY of the Stat struct
		// - s is a local variable (lives on stack) containing the copy
		// - Modifying s does NOT modify the map
		// - We MUST write s back to update the map
		//
		// Why not stats[category].Count++?
		// - Go prohibits this! Map elements are not addressable
		// - You can't take address of stats[category] because the hash table
		//   might resize and move elements to different memory locations
		//
		// BREAKPOINT 27: Set here BEFORE map lookup
		// DEBUG: In Variables panel, expand 'stats' map
		// See if category already exists in map
		// DEBUG: In Debug Console, type: stats[category]
		// If category doesn't exist, returns Stat{Count:0, Sum:0.0, Avg:0.0}
		// BREAKPOINT 28: Set here AFTER map lookup
		// DEBUG: In Variables panel, expand 's' to see:
		//   - s.Count: current count (0 if first occurrence)
		//   - s.Sum: current sum (0.0 if first occurrence)
		// DEBUG: In Debug Console, type: s.Count
		// DEBUG: In Debug Console, type: s.Sum
		s := stats[category]

		// BREAKPOINT 29: Set here BEFORE incrementing count
		// DEBUG: In Debug Console, type: s.Count
		// See current count before increment
		s.Count++
		// BREAKPOINT 30: Set here AFTER incrementing count
		// DEBUG: In Debug Console, type: s.Count
		// Verify count increased by 1

		// BREAKPOINT 31: Set here BEFORE adding amount
		// DEBUG: In Debug Console, type: s.Sum
		// DEBUG: In Debug Console, type: amount
		// See current sum and amount to add
		s.Sum += amount
		// BREAKPOINT 32: Set here AFTER adding amount
		// DEBUG: In Debug Console, type: s.Sum
		// Verify sum increased by amount

		// BREAKPOINT 33: Set here BEFORE writing back to map
		// DEBUG: In Variables panel, compare:
		//   - s: modified struct (on stack)
		//   - stats[category]: old value in map (if exists)
		// BREAKPOINT 34: Set here AFTER writing back to map
		// DEBUG: In Variables panel, expand 'stats' map
		// See updated value for this category
		// DEBUG: In Debug Console, type: stats[category]
		// Verify it matches 's'
		stats[category] = s

		// BREAKPOINT 35: Set here BEFORE rowNum increment
		// DEBUG: In Debug Console, type: rowNum
		rowNum++
		// BREAKPOINT 36: Set here AFTER rowNum increment
		// DEBUG: In Debug Console, type: rowNum
		// Verify it increased by 1
		// Step Over (F10) to continue to next row
	}

	// BREAKPOINT 37: Set here AFTER loop completes
	// DEBUG: In Variables panel, expand 'stats' map
	// All categories should have Count and Sum, but Avg is still 0
	// DEBUG: In Debug Console, type: stats
	// DEBUG: In Debug Console, type: len(stats)

	// ============================================================================
	// AVERAGE COMPUTATION: Second pass to calculate averages
	// ============================================================================
	// Compute averages
	// We do this in a separate pass to avoid redundant calculations
	// (Average is Sum/Count, so we only need to compute it once per category)
	//
	// BREAKPOINT 38: Set here BEFORE average loop
	// Step Over (F10) to iterate through each category
	for category, s := range stats {
		// BREAKPOINT 39: Set here at start of loop body
		// DEBUG: In Variables panel, watch:
		//   - category: current category name
		//   - s: struct with Count and Sum (Avg still 0)
		// DEBUG: In Debug Console, type: category
		// DEBUG: In Debug Console, type: s.Count
		// DEBUG: In Debug Console, type: s.Sum

		if s.Count > 0 {
			// BREAKPOINT 40: Set here BEFORE computing average
			// DEBUG: In Debug Console, type: s.Sum
			// DEBUG: In Debug Console, type: float64(s.Count)
			// DEBUG: In Debug Console, type: s.Sum / float64(s.Count)
			// See the computed average
			s.Avg = s.Sum / float64(s.Count)

			// BREAKPOINT 41: Set here AFTER computing average
			// DEBUG: In Variables panel, watch 's.Avg'
			// DEBUG: In Debug Console, type: s.Avg
			// Verify average is correct (Sum / Count)

			// BREAKPOINT 42: Set here BEFORE writing back to map
			stats[category] = s
			// BREAKPOINT 43: Set here AFTER writing back to map
			// DEBUG: In Variables panel, expand 'stats' map
			// Verify this category now has all three fields populated
			// DEBUG: In Debug Console, type: stats[category].Avg
		}
		// If Count is 0 (shouldn't happen in this logic), Avg remains 0.0

		// BREAKPOINT 44: Set here at end of loop iteration
		// Step Over (F10) to continue to next category
	}

	// BREAKPOINT 45: Set here AFTER average loop completes
	// DEBUG: In Variables panel, expand 'stats' map
	// All categories should have Count, Sum, AND Avg now
	// DEBUG: In Debug Console, type: stats
	// Verify all statistics are complete

	// ============================================================================
	// RETURN: Final result
	// ============================================================================
	// BREAKPOINT 46: Set here at return statement
	// DEBUG: In Variables panel, review 'stats' map one final time
	// All categories should have complete statistics
	// DEBUG: In Debug Console, type: len(stats)
	// DEBUG: In Debug Console, type: stats
	// Final verification before returning
	return stats, nil
}

/*
ADVANCED DEBUGGING TECHNIQUES:
===============================

1. CONDITIONAL BREAKPOINTS:
   Right-click on any breakpoint and add conditions
   Examples:
   - In main loop: category == "groceries" (break only for specific category)
   - In main loop: amount > 100.0 (break only for large amounts)
   - In average loop: s.Count > 5 (break only for categories with many transactions)

2. LOGPOINTS:
   Right-click on line number → Add Logpoint
   Examples:
   - In main loop: Log "Row {rowNum}: {category} = ${amount}"
   - In map update: Log "Category {category}: Count={s.Count}, Sum={s.Sum}"
   - No need to modify code or add print statements!

3. DEBUG CONSOLE EXPRESSIONS:
   During debugging, try these in the Debug Console:
   - Type: stats to see current statistics map
   - Type: len(stats) to see number of categories
   - Type: stats["groceries"] to check stats for specific category
   - Type: amount to see current parsed amount
   - Type: s.Sum / float64(s.Count) to manually compute average

4. WATCH EXPRESSIONS:
   Add these to the Watch panel for real-time monitoring:
   - len(stats) - number of categories
   - stats[category] - stats for current category
   - s.Count - transaction count
   - s.Sum - running sum
   - s.Avg - computed average

5. CALL STACK NAVIGATION:
   In the Call Stack panel:
   - Click different frames to see state at each level
   - Observe how variables change across stack frames
   - Use "Step Out" (Shift+F11) to return to caller

6. MEMORY INSPECTION:
   To see memory allocations:
   - Watch how stats map grows with each category
   - Notice struct values vs pointer values
   - Compare memory addresses when structs are copied

7. STEP COMMANDS:
   - F10 (Step Over): Execute line, don't enter functions
   - F11 (Step Into): Enter function calls (like ParseFloat) to see internals
   - Shift+F11 (Step Out): Return to caller
   - F5 (Continue): Run until next breakpoint

8. DATA BREAKPOINTS:
   Watch for when specific variables change:
   - Right-click variable → Break When Value Changes
   - Useful for tracking when stats map is modified
   - Useful for tracking when rowNum increments

Alternatives & Trade-offs:
==========================

1. Pointer values in map:
   stats := make(map[string]*Stat)
   s := stats[category]
   if s == nil { s = &Stat{}; stats[category] = s }
   s.Count++
   s.Sum += amount
   Pros: Modify in-place (no write-back needed)
   Cons: Extra allocations; nil checks required

2. Struct with embedded mutex for thread-safety:
   type Stat struct {
     sync.Mutex
     Count int
     Sum   float64
   }
   s.Lock(); s.Count++; s.Unlock()
   Pros: Safe for concurrent access
   Cons: Overkill for single-goroutine code; performance overhead

3. Accumulate errors instead of failing fast:
   var errs []error
   // ... on parse error: errs = append(errs, err)
   if len(errs) > 0 { return nil, errors.Join(errs...) }
   Pros: Process entire file even with some bad rows
   Cons: More complex error handling; may hide systemic issues

4. Use integer cents instead of float64:
   Floating-point arithmetic has rounding errors:
     0.1 + 0.2 = 0.30000000000000004 (in binary!)
   For financial data, multiply by 100 and use int64 cents:
     amount, _ := strconv.ParseInt(amountStr, 10, 64)
     amountCents := int64(amount * 100)
   Then divide by 100.0 for display.
   Pros: Exact arithmetic; no rounding errors
   Cons: More code; still need float64 for display

DEBUGGING EXERCISES:
====================

Exercise 1: Trace SummarizeCSV execution
- Input: CSV with 3 categories, 2 transactions each
- Set breakpoints at: 94, 214, 326, 380
- Watch: stats, rowNum, category
- Question: How many times does the main loop iterate?
- Question: What are the final values in the stats map?

Exercise 2: Understand struct read-modify-write pattern
- Input: CSV with "groceries" appearing twice
- Set breakpoints at: 326, 354
- In Debug Console at 326: stats["groceries"]
- Question: What does stats["groceries"] return the first time vs second time?
- Question: Why do we need to write back to the map?

Exercise 3: Watch floating-point arithmetic
- Input: CSV with amounts: 10.50, 20.25, 30.10
- Set breakpoints at: 340, 395
- Watch: s.Sum in Variables panel
- Question: Is the sum exact or does it have rounding errors?
- Question: How is the average computed?

Exercise 4: Empty CSV edge case
- Input: CSV with only header row
- Set breakpoints at: 94, 214, 365, 430
- Question: Does the main loop execute? How many times?
- Question: What does the function return?

Exercise 5: Malformed CSV handling
- Input: CSV with invalid amount in row 2
- Set breakpoints at: 289, 291
- Watch: amountStr, err in Variables panel
- Question: When does the error occur?
- Question: What information is in the error message?

Exercise 6: Category validation
- Input: CSV with empty category in row 3
- Set breakpoints at: 269, 271
- Watch: category in Variables panel
- Question: What happens when category is empty?
- Question: At what row number does the error occur?

Exercise 7: Average computation
- Input: CSV with "books" category: 10.00, 20.00, 30.00
- Set breakpoints at: 365, 395, 420
- Watch: stats["books"] at each breakpoint
- Question: When is Avg set to 0? When is it computed?
- Question: What is the final Avg value?
*/
