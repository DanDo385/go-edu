package main

import (
	"fmt"
	"strings"

	"github.com/example/go-10x-minis/minis/03-csv-stats/internal/csvstats"
)

/*
Debug Harness for CSV Stats Module

This file runs the CSV summarizer with sample data.
Set breakpoints in exercise.go and press F5 in VS Code.
*/
func main() {
	fmt.Println("=== CSV Stats Debug Harness ===")
	fmt.Println()

	// Sample CSV data
	csvData := `category,amount
Food,25.50
Transport,15.00
Food,12.75
Entertainment,50.00
Transport,8.25
Food,30.00
Entertainment,25.00
Transport,12.50`

	fmt.Println("--- Sample CSV Data ---")
	fmt.Println(csvData)
	fmt.Println()

	reader := strings.NewReader(csvData)
	stats, err := csvstats.SummarizeCSV(reader)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("--- Results ---")
	for category, stat := range stats {
		fmt.Printf("\n%s:\n", category)
		fmt.Printf("  Count:   %d\n", stat.Count)
		fmt.Printf("  Sum:     %.2f\n", stat.Sum)
		fmt.Printf("  Average: %.2f\n", stat.Avg)
	}
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. encoding/csv parses CSV files")
	fmt.Println("2. Maps aggregate data by key (category)")
	fmt.Println("3. strconv.ParseFloat converts strings to floats")
	fmt.Println("4. Structs organize related statistics")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/04-jsonl-log-filter")
}
