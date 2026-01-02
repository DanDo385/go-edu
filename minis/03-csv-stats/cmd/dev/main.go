package main

import (
	"fmt"
	"os"
	"strings"

	"minis/03-csv-stats/internal/csvstats"
)

/*
Debug Harness for CSV Stats

This file automatically demonstrates the project's capabilities by running through
different scenarios with pre-configured inputs. No CLI arguments needed!

How to use:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5 in VS Code
  3. Select "Debug cmd/dev (Debug Harness)"
  4. Step through with F10 (Step Over) and F11 (Step Into)
*/

func main() {
	fmt.Println("=== CSV Stats - Auto Demo ===\n")

	// Demo 1: Process testdata file
	fmt.Println("Demo 1: Processing testdata/transactions.csv")
	fmt.Println("----------------------------------------")
	testdataPath := "testdata/transactions.csv"
	if _, err := os.Stat(testdataPath); err == nil {
		file, err := os.Open(testdataPath)
		if err == nil {
			fmt.Printf("Reading from: %s\n\n", testdataPath)
			
			// BREAKPOINT: Step into SummarizeCSV to see CSV parsing
			stats, err := csvstats.SummarizeCSV(file)
			file.Close()
			
			if err != nil {
				fmt.Printf("ERROR: %v\n\n", err)
			} else {
				fmt.Println("Results:")
				for category, stat := range stats {
					fmt.Printf("  %s: Count=%d, Sum=%.2f, Avg=%.2f\n", 
						category, stat.Count, stat.Sum, stat.Avg)
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Printf("Testdata file not found: %s\n\n", testdataPath)
	}

	// Demo 2: Inline CSV data
	fmt.Println("Demo 2: Processing inline CSV data")
	fmt.Println("----------------------------------------")
	csvData := `category,amount
Food,10.50
Transport,5.25
Food,15.75
Entertainment,20.00
Transport,3.50
Food,8.00`
	
	fmt.Println("CSV Data:")
	fmt.Println(csvData)
	fmt.Println()
	
	// BREAKPOINT: See how CSV parsing handles different categories
	stats, err := csvstats.SummarizeCSV(strings.NewReader(csvData))
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
	} else {
		fmt.Println("Results:")
		for category, stat := range stats {
			fmt.Printf("  %s:\n", category)
			fmt.Printf("    Count: %d\n", stat.Count)
			fmt.Printf("    Sum:   %.2f\n", stat.Sum)
			fmt.Printf("    Avg:   %.2f\n\n", stat.Avg)
		}
	}

	// Demo 3: Single category
	fmt.Println("Demo 3: Single category")
	fmt.Println("----------------------------------------")
	singleCatCSV := `category,amount
Food,10.00
Food,20.00
Food,30.00`
	
	fmt.Println("CSV Data:")
	fmt.Println(singleCatCSV)
	fmt.Println()
	
	// BREAKPOINT: See how average is calculated
	stats, err = csvstats.SummarizeCSV(strings.NewReader(singleCatCSV))
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
	} else {
		fmt.Println("Results:")
		for category, stat := range stats {
			fmt.Printf("  %s: Count=%d, Sum=%.2f, Avg=%.2f\n", 
				category, stat.Count, stat.Sum, stat.Avg)
		}
		fmt.Println()
	}

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. CSV parsing requires header validation")
	fmt.Println("2. Amount parsing must handle decimal numbers")
	fmt.Println("3. Statistics aggregation uses maps for grouping")
	fmt.Println("4. Error handling is critical for malformed data")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/04-jsonl-log-filter")
}
