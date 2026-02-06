package main

import (
	"fmt"
	"os"

	"github.com/example/go-10x-minis/minis/02-arrays-maps-basics/internal/arraysmapsbasics"
)

/*
Word Frequency Counter CLI

Usage:

	go run ./cmd/app/main.go <file>
	echo "text" | go run ./cmd/app/main.go -

Arguments:

	file   Path to text file, or "-" for stdin

Examples:

	go run ./cmd/app/main.go testdata/input.txt
	echo "hello world hello" | go run ./cmd/app/main.go -
*/
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var reader *os.File
	var err error

	if os.Args[1] == "-" {
		reader = os.Stdin
		fmt.Println("Reading from stdin (press Ctrl+D when done)...")
	} else {
		reader, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer reader.Close()
		fmt.Printf("Reading from: %s\n", os.Args[1])
	}
	fmt.Println()

	freq, mostCommon, err := arraysmapsbasics.FreqFromReader(reader)
	if err != nil {
		fmt.Printf("Error processing input: %v\n", err)
		os.Exit(1)
	}

	// Show results
	fmt.Println("=== Word Frequency Results ===")
	fmt.Printf("Total unique words: %d\n", len(freq))
	fmt.Printf("Most common word:   %q (count: %d)\n", mostCommon, freq[mostCommon])
	fmt.Println()

	// Show top 10 words
	fmt.Println("Top words:")
	count := 0
	for word, cnt := range freq {
		if count >= 10 {
			fmt.Println("  ...")
			break
		}
		fmt.Printf("  %q: %d\n", word, cnt)
		count++
	}
}

func printUsage() {
	fmt.Println("Word Frequency Counter CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/app/main.go <file>")
	fmt.Println("  echo \"text\" | go run ./cmd/app/main.go -")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  file   Path to text file, or \"-\" for stdin")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/app/main.go testdata/input.txt")
	fmt.Println("  echo \"hello world hello\" | go run ./cmd/app/main.go -")
}
