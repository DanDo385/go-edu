package main

import (
	"fmt"
	"log"
	"os"

	"github.com/example/go-10x-minis/minis/03-csv-stats/exercise"
)

func main() {
	// Open the testdata CSV file
	file, err := os.Open("minis/03-csv-stats/testdata/transactions.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Compute statistics
	stats, err := exercise.SummarizeCSV(file)
	if err != nil {
		log.Fatalf("Failed to process CSV: %v", err)
	}

	// Display results
	fmt.Println("=== Transaction Statistics by Category ===\n")
	for category, stat := range stats {
		fmt.Printf("Category: %s\n", category)
		fmt.Printf("  Count:   %d\n", stat.Count)
		fmt.Printf("  Sum:     $%.2f\n", stat.Sum)
		fmt.Printf("  Average: $%.2f\n", stat.Avg)
		fmt.Println()
	}
}
