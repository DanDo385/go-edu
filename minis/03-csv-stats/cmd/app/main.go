package main

import (
	"fmt"    // Formatted I/O
	"log"    // Logging
	"os"     // Operating system interface
	"sort"   // Sorting utilities

	// Import the exercise package from parent directory
	"github.com/example/go-10x-minis/minis/03-csv-stats/internal/csvstats"
)

/*
===================================================================
CSV Statistics: Transaction Analysis
===================================================================

This program demonstrates CSV file parsing and statistics computation
by analyzing transaction data grouped by category.

USAGE:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go testdata/transactions.csv
    go run ./cmd/app/main.go custom.csv sum
    go run ./cmd/app/main.go testdata/transactions.csv count reverse

Arguments:
    [file]    - Path to CSV transactions file (default: "testdata/transactions.csv")
    [sort]    - Sort by: none, count, sum, avg, category (default: "category")
    [reverse] - Reverse sort order if "reverse" provided (default: false)

Examples:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go testdata/transactions.csv
    go run ./cmd/app/main.go testdata/transactions.csv sum
    go run ./cmd/app/main.go testdata/transactions.csv count reverse

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter csvstats.SummarizeCSV function
- Watch Variables panel to see map[category]Stats transformations
- See /RUN_DEBUG.md for comprehensive debugging guide

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the exercise package from internal/exercise
- csvstats.go, solution.go are in package exercise
- Functions must be exported (capitalized) to be called from main
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the first argument (if provided)
	// DEBUG: In Variables panel, expand 'os.Args' to see all arguments

	// Default values
	filePath := "testdata/transactions.csv"
	sortBy := "category"
	reverse := false
	verbose := false

	// Override defaults if arguments provided
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	if len(os.Args) > 2 {
		sortBy = os.Args[2]
	}
	if len(os.Args) > 3 && os.Args[3] == "reverse" {
		reverse = true
	}
	if len(os.Args) > 3 && os.Args[3] == "verbose" {
		verbose = true
	}

	// BREAKPOINT: Set breakpoint here to inspect parsed values in Variables panel
	// DEBUG: After parsing, variables now hold actual command-line values
	// DEBUG: Watch 'filePath', 'sortBy', 'reverse', 'verbose'

	// ============================================================================
	// STEP 2: Open and Read CSV File
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before opening file
	// DEBUG: Watch 'filePath' variable - this is the file we'll process
	// DEBUG: Step Over (F10) to execute file opening

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file %q: %v", filePath, err)
	}
	defer file.Close()

	// BREAKPOINT: Set breakpoint here before parsing CSV
	// DEBUG: File is now open and ready to read
	// DEBUG: Step Into (F11) on csvstats.SummarizeCSV() to see CSV parsing

	fmt.Printf("Processing CSV file: %s\n", filePath)
	fmt.Println()

	// ============================================================================
	// STEP 3: Compute Statistics
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before calling SummarizeCSV
	// DEBUG: Step Into (F11) to see how CSV is parsed line by line
	// DEBUG: Watch how map[string]Stats is built incrementally

	stats, err := csvstats.SummarizeCSV(file)

	// BREAKPOINT: Set breakpoint here after computing statistics
	// DEBUG: Inspect 'stats' map - expand to see all categories
	// DEBUG: Each category has Count, Sum, and Avg fields
	// DEBUG: Notice map structure: map[string]Stats where Stats is a struct

	if err != nil {
		log.Fatalf("Failed to process CSV: %v", err)
	}

	// ============================================================================
	// STEP 4: Sort Results (Optional)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before sorting
	// DEBUG: Watch 'sortBy' variable to see which sort mode is active
	// DEBUG: Maps are unordered - we need to convert to slice for sorting

	// Convert map to slice for sorting
	type CategoryStats struct {
		Category string
		Stats    csvstats.Stats
	}

	// BREAKPOINT: Set breakpoint here before building slice
	// DEBUG: We're converting map to slice so we can sort it
	// DEBUG: Watch 'categories' slice grow as we iterate the map

	categories := make([]CategoryStats, 0, len(stats))
	for category, stat := range stats {
		// BREAKPOINT: Set breakpoint here inside loop
		// DEBUG: Watch 'category' and 'stat' change with each iteration
		// DEBUG: See how we append to 'categories' slice

		categories = append(categories, CategoryStats{
			Category: category,
			Stats:    stat,
		})
	}

	// BREAKPOINT: Set breakpoint here before sorting
	// DEBUG: 'categories' slice now contains all data
	// DEBUG: Step Over (F10) to execute sort and watch order change

	// Sort based on flag
	if sortBy != "none" {
		sort.Slice(categories, func(i, j int) bool {
			// BREAKPOINT: Set breakpoint here inside sort comparator
			// DEBUG: This function is called repeatedly during sorting
			// DEBUG: Watch 'i' and 'j' indices as sort algorithm progresses

			less := false
			switch sortBy {
			case "count":
				less = categories[i].Stats.Count < categories[j].Stats.Count
			case "sum":
				less = categories[i].Stats.Sum < categories[j].Stats.Sum
			case "avg":
				less = categories[i].Stats.Avg < categories[j].Stats.Avg
			case "category":
				less = categories[i].Category < categories[j].Category
			default:
				less = categories[i].Category < categories[j].Category
			}

			if reverse {
				return !less
			}
			return less
		})
	}

	// ============================================================================
	// STEP 5: Display Results
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before displaying results
	// DEBUG: 'categories' slice is now sorted according to flags
	// DEBUG: Step Over (F10) to see results printed

	fmt.Println("=== Transaction Statistics by Category ===")
	fmt.Printf("Sort: %s", sortBy)
	if reverse {
		fmt.Print(" (reverse)")
	}
	fmt.Println()
	fmt.Println()

	// BREAKPOINT: Set breakpoint here before loop
	// DEBUG: Expand 'categories' in Variables panel to see all entries
	// DEBUG: Notice the sort order based on your flags

	totalCount := 0
	totalSum := 0.0

	for _, cs := range categories {
		// BREAKPOINT: Set breakpoint here inside display loop
		// DEBUG: Watch 'cs.Category' and 'cs.Stats' change with each iteration
		// DEBUG: See how we accumulate totals in 'totalCount' and 'totalSum'

		if verbose {
			fmt.Printf("Category: %s\n", cs.Category)
			fmt.Printf("  Count:   %d transactions\n", cs.Stats.Count)
			fmt.Printf("  Sum:     $%.2f\n", cs.Stats.Sum)
			fmt.Printf("  Average: $%.2f\n", cs.Stats.Avg)
			fmt.Println()
		} else {
			fmt.Printf("%-20s | Count: %4d | Sum: $%10.2f | Avg: $%8.2f\n",
				cs.Category, cs.Stats.Count, cs.Stats.Sum, cs.Stats.Avg)
		}

		// Accumulate totals
		totalCount += cs.Stats.Count
		totalSum += cs.Stats.Sum
	}

	// ============================================================================
	// STEP 6: Display Summary
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before summary
	// DEBUG: Inspect 'totalCount' and 'totalSum' - accumulated from all categories
	// DEBUG: Verify totals match sum of individual categories

	fmt.Println()
	fmt.Printf("Total Categories: %d\n", len(categories))
	fmt.Printf("Total Transactions: %d\n", totalCount)
	fmt.Printf("Total Amount: $%.2f\n", totalSum)
	if totalCount > 0 {
		fmt.Printf("Overall Average: $%.2f\n", totalSum/float64(totalCount))
	}

	// ============================================================================
	// STEP 7: Usage Note about Solution
	// ============================================================================
	if *useSolution {
		fmt.Println()
		fmt.Println("NOTE: To use solution.go, build/run with -tags=solution:")
		fmt.Println("  go run -tags=solution ./cmd/csv-stats/main.go -solution=true")
		fmt.Println("  go test -tags=solution")
		os.Exit(0)
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES FOR CSV PARSING AND STATISTICS:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter csvstats.SummarizeCSV function
	// 5. Watch Variables panel to see map[string]Stats build up
	// 6. Add watch expressions: len(stats), stats["Food"], totalSum
	// 7. Use Debug Console to inspect: stats["category"].Sum
	//
	// UNDERSTANDING MAP[STRING]STRUCT:
	// - In Variables panel, expand 'stats' to see all categories
	// - Each key is a category name (string)
	// - Each value is a Stats struct with Count, Sum, Avg fields
	// - Watch how map grows as CSV is processed line by line
	//
	// UNDERSTANDING SORTING:
	// - Maps are unordered - we convert to slice for sorting
	// - sort.Slice uses a comparison function (less function)
	// - Set breakpoint inside comparator to see sorting algorithm
	// - Watch indices 'i' and 'j' as QuickSort proceeds
	//
	// CSV PARSING:
	// - Step Into (F11) csvstats.SummarizeCSV to see csv.Reader
	// - Watch how each row is parsed into fields
	// - See how amounts are converted from string to float64
	// - Notice error handling for malformed CSV
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}
