// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/csvstats.go
// 2. Open this file (cmd/debug/main.go)
// 3. Use "Debug Main Program (Current Package)" configuration
// 4. Press F5 - that's it! The debugger will stop at your breakpoints
//
// Usage:
//   go run ./cmd/debug
//   # Or just press F5 in VS Code

package main

import (
	"fmt"
	"strings"

	"github.com/example/go-10x-minis/minis/03-csv-stats/internal/csvstats"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	csvData := `category,amount
groceries,50.00
groceries,30.00
electronics,200.00
electronics,150.00
groceries,25.00`

	fmt.Println("=== Debugging SummarizeCSV ===")
	fmt.Printf("CSV Data:\n%s\n\n", csvData)

	reader := strings.NewReader(csvData)

	// Set breakpoint in csvstats.go at SummarizeCSV function
	stats, err := csvstats.SummarizeCSV(reader)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("CSV Statistics: %+v\n", stats)
}

