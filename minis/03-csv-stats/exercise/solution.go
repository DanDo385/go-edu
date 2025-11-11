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

Compared to other languages:
- Python (pandas): df.groupby('category')['amount'].agg(['count','sum','mean'])
  Pros: One-liner, powerful analytics
  Cons: Loads entire file into memory; slower for large files
- JavaScript: Requires external CSV library; async I/O complicates streaming
- Rust: `csv` crate is excellent; zero-copy parsing is faster than Go
- SQL: Natural fit (SELECT category, COUNT(*), SUM(amount), AVG(amount) GROUP BY category)
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
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	// Create a CSV reader
	// csv.Reader automatically handles:
	// - Line breaks within quoted fields
	// - Escaped quotes (double quotes: "He said ""hello""")
	// - Different line endings (\n vs \r\n)
	csvReader := csv.NewReader(r)

	// Read the header row
	// This validates the CSV structure and allows us to check column names
	headers, err := csvReader.Read()
	if err != nil {
		// Distinguish between empty file vs. I/O error
		if err == io.EOF {
			return nil, fmt.Errorf("empty CSV file (no header)")
		}
		return nil, fmt.Errorf("reading header: %w", err)
	}

	// Validate header format
	// We expect exactly 3 columns: id, category, amount
	// This is a defensive check to fail fast on schema mismatches
	if len(headers) != 3 || headers[0] != "id" || headers[1] != "category" || headers[2] != "amount" {
		return nil, fmt.Errorf("invalid header: expected [id,category,amount], got %v", headers)
	}

	// Initialize the aggregation map
	// Key: category name
	// Value: running totals (count and sum; average computed later)
	stats := make(map[string]Stat)

	// Track row number for error messages (starting at 2 since row 1 is header)
	rowNum := 2

	// Read records line-by-line (streaming)
	// csv.Reader.Read() returns []string for each row
	// It returns io.EOF when the file is exhausted (not an error!)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			// End of file is the normal exit condition
			break
		}
		if err != nil {
			// Unexpected error (e.g., malformed CSV, I/O failure)
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}

		// Expect exactly 3 fields per row
		// csv.Reader ensures this by default (FieldsPerRecord = 0 means "match header")
		// but we check defensively
		if len(record) != 3 {
			return nil, fmt.Errorf("row %d: expected 3 fields, got %d", rowNum, len(record))
		}

		// Extract fields
		// We don't use the id field in this analysis, but could validate it's numeric
		// id := record[0]  // unused
		category := record[1]
		amountStr := record[2]

		// Validate category is not empty
		if category == "" {
			return nil, fmt.Errorf("row %d: empty category", rowNum)
		}

		// Parse amount as float64
		// strconv.ParseFloat returns an error if the string is not a valid number
		// The second argument (64) specifies float64 precision
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid amount %q: %w", rowNum, amountStr, err)
		}

		// Update statistics for this category
		// Map lookup returns the zero value (Stat{Count:0, Sum:0.0}) if key doesn't exist
		// This is perfect for aggregation: we can read, modify, and write back
		s := stats[category]
		s.Count++
		s.Sum += amount
		stats[category] = s
		// Note: We must write back to the map because structs are value types!
		// Alternatively: Use *Stat (pointer) as the map value type to modify in-place

		rowNum++
	}

	// Compute averages
	// We do this in a separate pass to avoid redundant calculations
	// (Average is Sum/Count, so we only need to compute it once per category)
	for category, s := range stats {
		if s.Count > 0 {
			s.Avg = s.Sum / float64(s.Count)
			stats[category] = s
		}
		// If Count is 0 (shouldn't happen in this logic), Avg remains 0.0
	}

	return stats, nil
}

/*
Alternatives & Trade-offs:

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

Go vs X:

Go vs Python (pandas):
  import pandas as pd
  df = pd.read_csv('transactions.csv')
  stats = df.groupby('category')['amount'].agg(['count','sum','mean'])
  Pros: Concise; powerful analytics (median, std dev, etc.)
  Cons: Loads entire file into memory; 100MB file = 500MB+ RAM
        Slower for large files (pandas overhead)
        Dynamic typing hides schema errors until runtime
  Go: Constant memory usage; catches type errors at compile time

Go vs SQL:
  SELECT category,
         COUNT(*) as count,
         SUM(amount) as sum,
         AVG(amount) as avg
  FROM transactions
  GROUP BY category;
  Pros: Declarative; optimized by query planner; scales to billions of rows
  Cons: Requires database setup; less portable than CSV file
  Go: Simpler deployment (single binary); good for <10M rows

Go vs Rust:
  use csv::ReaderBuilder;
  let mut stats: HashMap<String, (usize, f64)> = HashMap::new();
  for result in rdr.records() {
      let record = result?;
      let (count, sum) = stats.entry(record[1].to_string()).or_insert((0, 0.0));
      *count += 1;
      *sum += record[2].parse::<f64>()?;
  }
  Pros: Zero-copy parsing; faster execution; no GC pauses
  Cons: More complex ownership rules; steeper learning curve
  Go: Simpler code; "fast enough" for most use cases

Go vs JavaScript (Node.js):
  const csv = require('csv-parser');
  const stats = {};
  fs.createReadStream('transactions.csv')
    .pipe(csv())
    .on('data', (row) => {
      if (!stats[row.category]) stats[row.category] = {count:0, sum:0};
      stats[row.category].count++;
      stats[row.category].sum += parseFloat(row.amount);
    });
  Pros: Streaming; widely known syntax
  Cons: Async complexity (callbacks/promises); requires external library
        Dynamic typing misses schema errors
  Go: Synchronous code is simpler; built-in CSV parser
*/
