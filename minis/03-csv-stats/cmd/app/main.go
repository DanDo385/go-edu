package main

import (
	"fmt"
	"os"

	"github.com/example/go-10x-minis/minis/03-csv-stats/internal/csvstats"
)

/*
CSV Statistics CLI

Usage:
  go run ./cmd/app/main.go <file>

Arguments:
  file   Path to CSV file

Examples:
  go run ./cmd/app/main.go testdata/transactions.csv
*/
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("Processing: %s\n\n", os.Args[1])

	stats, err := csvstats.SummarizeCSV(file)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Category Statistics ===")
	for category, stat := range stats {
		fmt.Printf("\n%s:\n", category)
		fmt.Printf("  Count:   %d\n", stat.Count)
		fmt.Printf("  Sum:     %.2f\n", stat.Sum)
		fmt.Printf("  Average: %.2f\n", stat.Avg)
	}
}

func printUsage() {
	fmt.Println("CSV Statistics CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/app/main.go <file>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  file   Path to CSV file")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/app/main.go testdata/transactions.csv")
}
