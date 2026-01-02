package main

import (
	"fmt"
	"os"
	"strings"

	"minis/02-arrays-maps-basics/internal/arraysmapsbasics"
)

/*
Arrays and Maps Basics CLI

This application demonstrates word frequency counting from text input.
It reads from stdin or a file and counts word frequencies, finding the most common word.

Usage:

	go run ./cmd/app/main.go [FILE]

Arguments:

	FILE - Optional file path to read from. If not provided, reads from stdin.

Examples:

	# Read from stdin
	echo -e "hello\nworld\nhello" | go run ./cmd/app/main.go

	# Read from file
	go run ./cmd/app/main.go testdata/input.txt

	# Read from file (alternative)
	go run ./cmd/app/main.go < testdata/input.txt

Copy & Paste Examples:

	go run ./cmd/app/main.go testdata/input.txt
	echo -e "go\ngo\nrust\nrust\nrust" | go run ./cmd/app/main.go
*/

func main() {
	var input *os.File
	var err error

	// Parse command line arguments
	if len(os.Args) > 1 {
		// Read from file
		filePath := os.Args[1]
		input, err = os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		defer input.Close()
		fmt.Printf("Reading from file: %s\n\n", filePath)
	} else {
		// Read from stdin
		input = os.Stdin
		fmt.Println("Reading from stdin (type text and press Ctrl+D to finish):\n")
	}

	// Count word frequencies
	freq, mostCommon, err := arraysmapsbasics.FreqFromReader(input)
	if err != nil {
		fmt.Printf("Error processing input: %v\n", err)
		os.Exit(1)
	}

	// Display results
	fmt.Println("=== Word Frequencies ===")
	if len(freq) == 0 {
		fmt.Println("(no words found)")
	} else {
		// Sort by frequency for display (simple approach)
		for word, count := range freq {
			fmt.Printf("%-20s: %d\n", word, count)
		}
	}

	fmt.Println()
	if mostCommon != "" {
		fmt.Printf("Most common word: %s (appears %d times)\n", mostCommon, freq[mostCommon])
	} else {
		fmt.Println("Most common word: (none)")
	}
}
