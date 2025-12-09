//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "encoding/csv" for CSV parsing
// - "fmt" for error formatting
// - "io" for the Reader interface
// - "strconv" for string to float conversion
//
// import (
//     "encoding/csv"
//     "fmt"
//     "io"
//     "strconv"
// )

// Stat holds aggregated statistics for a category
// This is a STRUCT TYPE (value type, not reference type)
//
// TODO: Define the Stat struct
// type Stat struct {
//     Count int     // Number of transactions in this category
//     Sum   float64 // Total amount for this category
//     Avg   float64 // Average amount (Sum / Count)
// }
//
// Key Go concepts for structs:
// - Structs are VALUE types (copied when assigned or passed)
// - Fields are accessed with dot notation: s.Count, s.Sum, s.Avg
// - Zero value: Stat{Count: 0, Sum: 0.0, Avg: 0.0}
// - When you write freq[category] = s, you COPY the struct into the map
// - To modify a struct in a map, you must:
//   1. Read it out: s := freq[category]
//   2. Modify the copy: s.Count++
//   3. Write it back: freq[category] = s
//   (You CANNOT do freq[category].Count++ directly!)

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
// TODO: Implement SummarizeCSV function
// Function signature: func SummarizeCSV(r io.Reader) (map[string]Stat, error)
//
// Steps to implement:
//
// 1. Create a CSV reader
//    - Use: csvReader := csv.NewReader(r)
//    - csv.Reader is a STRUCT (value type) but contains pointers internally
//    - The Reader maintains internal state (position in file, buffer, etc.)
//
// 2. Read and validate the header row
//    - Use: headers, err := csvReader.Read()
//    - Read() returns []string (slice of column values) and error
//    - Check if err == io.EOF (empty file - should be an error!)
//    - Check if err != nil (other read errors)
//    - Validate headers match exactly: ["id", "category", "amount"]
//    - Use len(headers) to check count
//    - Use headers[0], headers[1], headers[2] to check values
//
// 3. Create the statistics map
//    - Use: stats := make(map[string]Stat)
//    - Map will store category → Stat
//    - Map is a REFERENCE type (stores pointer to hash table)
//
// 4. Track row number for error messages
//    - Use: rowNum := 2 (start at 2: row 1 is header)
//    - Increment after processing each row: rowNum++
//
// 5. Loop through data rows
//    - Use: for { record, err := csvReader.Read(); ... }
//    - Check if err == io.EOF (end of file - break from loop)
//    - Check if err != nil (other errors - return error with row number)
//    - Validate len(record) == 3 (id, category, amount)
//
// 6. Extract and validate fields
//    - category := record[1] (second column, index 1)
//    - amountStr := record[2] (third column, index 2)
//    - Check if category == "" (empty category should error)
//    - Parse amount: amount, err := strconv.ParseFloat(amountStr, 64)
//      * ParseFloat returns (float64, error)
//      * First arg is string to parse
//      * Second arg is bit size (64 for float64)
//      * If parsing fails, return error with row number and invalid value
//
// 7. Update statistics in map
//    - CRITICAL: You cannot modify struct fields directly in a map!
//    - WRONG: stats[category].Count++ // Compile error!
//    - RIGHT:
//      * Read: s := stats[category]
//      * Modify copy: s.Count++, s.Sum += amount
//      * Write back: stats[category] = s
//    - Why? Maps store VALUES, not pointers
//    - Reading stats[category] returns a COPY of the Stat struct
//    - You must write the modified copy back to the map
//    - If category doesn't exist, stats[category] returns zero-value Stat{0, 0.0, 0.0}
//
// 8. Calculate averages
//    - After reading all rows, loop through map again
//    - For each category: for category, s := range stats { ... }
//    - Calculate: s.Avg = s.Sum / float64(s.Count)
//      * Must convert s.Count to float64 for division
//      * Integer division truncates: 5 / 2 = 2 (wrong!)
//      * Float division: 5.0 / 2.0 = 2.5 (correct!)
//    - Write back: stats[category] = s
//
// 9. Return results
//    - Use: return stats, nil
//    - stats is a map (reference type) - passes pointer to caller
//    - Caller can modify the map's contents
//
// Key Go concepts:
// - Structs are value types (copied when assigned)
// - Maps store values, not pointers (can't modify struct fields in place)
// - Must read-modify-write pattern for map entries
// - CSV Reader maintains internal state
// - Error messages should include context (row number, field value)
// - strconv.ParseFloat for string to float conversion
// - Type conversion: int to float64 for division
//
// Common mistakes:
// - Trying stats[key].Field++ (doesn't compile)
// - Forgetting to check err == io.EOF separately
// - Not including row numbers in error messages
// - Integer division instead of float division

// TODO: Implement the SummarizeCSV function below
// func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
//     return nil, nil
// }

// After implementing:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Try with invalid CSV to test error handling
// - Compare with solution.go to see detailed explanations
