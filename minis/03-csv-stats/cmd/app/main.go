package main

import (
	"fmt"
	"os"

	"minis/03-csv-stats/internal/csvstats"
)

/*
CSV Stats CLI

This application computes per-category statistics from a CSV of financial transactions.
The CSV must have a header row with columns for category and amount.

Usage:

	go run ./cmd/app/main.go <CSV_FILE>

Arguments:

	CSV_FILE - Path to CSV file with header row (category, amount)

Examples:

	# Process transactions CSV
	go run ./cmd/app/main.go testdata/transactions.csv

	# Process from stdin
	go run ./cmd/app/main.go < testdata/transactions.csv

Copy & Paste Examples:

	go run ./cmd/app/main.go testdata/transactions.csv
*/

func main() {
	var input *os.File
	var err error

	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/app/main.go <CSV_FILE>")
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  CSV_FILE - Path to CSV file with header row (category, amount)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run ./cmd/app/main.go testdata/transactions.csv")
		os.Exit(1)
	}

	filePath := os.Args[1]
	input, err = os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer input.Close()

	fmt.Printf("Processing CSV file: %s\n\n", filePath)

	// Process CSV
	stats, err := csvstats.SummarizeCSV(input)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	// Display results
	fmt.Println("=== Category Statistics ===")
	if len(stats) == 0 {
		fmt.Println("(no data found)")
	} else {
		for category, stat := range stats {
			fmt.Printf("\nCategory: %s\n", category)
			fmt.Printf("  Count: %d\n", stat.Count)
			fmt.Printf("  Sum:   %.2f\n", stat.Sum)
			fmt.Printf("  Avg:   %.2f\n", stat.Avg)
		}
	}
}
