package main

import (
	"fmt"    // Formatted I/O
	"log"    // Logging
	"os"     // Operating system interface
	"strings" // String manipulation

	// Import the exercise package from parent directory
	"github.com/example/go-10x-minis/minis/02-arrays-maps-basics/internal/exercise"
)

/*
===================================================================
Arrays and Maps: Word Frequency Counter
===================================================================

This program demonstrates fundamental Go arrays, slices, and maps by
counting word frequencies from text using the exercise package.

USAGE:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go "go is awesome"
    go run ./cmd/app/main.go "testdata/input.txt" file
    go run ./cmd/app/main.go "hello world hello" string

Arguments:
    [input]  - Input text or file path (default: "go is awesome go rocks go")
    [demo]   - Demo to run: all, file, string, custom (default: "all")

Examples:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go "go is awesome"
    go run ./cmd/app/main.go "testdata/input.txt" file
    go run ./cmd/app/main.go "hello world hello" string

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see map and slice transformations
- See /RUN_DEBUG.md for comprehensive debugging guide

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the exercise package from internal/exercise
- exercise.go, solution.go are in package exercise
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

	// Get arguments from command line
	// Default values
	input := "go is awesome go rocks go"
	demo := "all"
	filePath := "testdata/input.txt"

	// Override defaults if arguments provided
	if len(os.Args) > 1 {
		input = os.Args[1]
	}
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}
	// If input looks like a file path, use it as filePath
	if len(os.Args) > 1 && (strings.Contains(os.Args[1], "/") || strings.HasSuffix(os.Args[1], ".txt")) {
		filePath = os.Args[1]
	}

	// BREAKPOINT: Set breakpoint here to inspect parsed values in Variables panel
	// DEBUG: After parsing, variables now hold actual command-line values
	// DEBUG: Watch 'demo', 'input', 'filePath'

	// ============================================================================
	// STEP 2: Example 1 - Read from File (Intermediate)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through file reading
	// DEBUG: This demonstrates io.Reader interface with os.File
	// DEBUG: Step Into (F11) on exercise.FreqFromReader() to see buffered reading

	if demo == "all" || demo == "file" {
		fmt.Println("=== Demo 1: Reading from File ===")
		fmt.Printf("File: %s\n", filePath)
		fmt.Println()

		// BREAKPOINT: Set breakpoint here before opening file
		// DEBUG: Watch 'filePath' variable - this is the path we'll open
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Could not open file %q: %v\n", filePath, err)
			log.Println("(This is OK for demo - continuing with other examples)")
			fmt.Println()
		} else {
			defer file.Close()

			// BREAKPOINT: Set breakpoint here before function call
			// DEBUG: Step Into (F11) to see how FreqFromReader processes the file
			// DEBUG: Watch how it uses bufio.Scanner to read line by line
			freq, mostCommon, err := exercise.FreqFromReader(file)

			// BREAKPOINT: Set breakpoint here after function call
			// DEBUG: Inspect 'freq' map - see word → count mappings
			// DEBUG: Inspect 'mostCommon' string - the word that appears most
			// DEBUG: Expand 'freq' in Variables panel to see all entries

			if err != nil {
				log.Fatalf("Error reading file: %v", err)
			}

			fmt.Printf("Frequency map: %v\n", freq)
			fmt.Printf("Most common word: %q (appears %d times)\n", mostCommon, freq[mostCommon])
			fmt.Printf("Total unique words: %d\n", len(freq))

			// BREAKPOINT: Set breakpoint here to inspect final results
			// DEBUG: Compare len(freq) with file contents
			// DEBUG: Verify mostCommon is actually the highest frequency

			fmt.Println()
		}
	}

	// ============================================================================
	// STEP 3: Example 2 - Read from String (Simple)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through simple example
	// DEBUG: This is the simplest use case - strings.NewReader creates io.Reader
	// DEBUG: Step Over (F10) to execute line by line
	// DEBUG: Watch Variables panel to see how 'freq' map is built

	if demo == "all" || demo == "string" {
		fmt.Println("=== Demo 2: Reading from String ===")
		fmt.Println("This shows io.Reader flexibility - same function, different source")
		fmt.Println()

		// Test with predefined string
		testInput := `go
is
awesome
go
rocks
go`

		// BREAKPOINT: Set breakpoint here before creating reader
		// DEBUG: Watch 'testInput' variable - notice the newlines
		// DEBUG: strings.NewReader converts string to io.Reader interface
		reader := strings.NewReader(testInput)

		// BREAKPOINT: Set breakpoint here before function call
		// DEBUG: Step Into (F11) on exercise.FreqFromReader() to see implementation
		// DEBUG: Notice it uses the same code for both file and string!
		freq, mostCommon, err := exercise.FreqFromReader(reader)

		// BREAKPOINT: Set breakpoint here after function call
		// DEBUG: Inspect 'freq' map - should show: go=3, is=1, awesome=1, rocks=1
		// DEBUG: Inspect 'mostCommon' - should be "go"

		if err != nil {
			log.Fatalf("Error reading string: %v", err)
		}

		fmt.Printf("Input text:\n%s\n\n", testInput)
		fmt.Printf("Frequency map: %v\n", freq)
		fmt.Printf("Most common word: %q (appears %d times)\n", mostCommon, freq[mostCommon])
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 3 - Custom Input (Interactive)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through custom input
	// DEBUG: This demonstrates using CLI flag for input
	// DEBUG: User can run: go run main.go -demo=custom -input="hello world hello"

	if demo == "all" || demo == "custom" {
		fmt.Println("=== Demo 3: Custom Input from CLI ===")
		fmt.Printf("Processing input: %q\n", input)
		fmt.Println()

		// BREAKPOINT: Set breakpoint here before creating reader
		// DEBUG: Watch 'input' variable - this came from CLI argument
		reader := strings.NewReader(input)

		// BREAKPOINT: Set breakpoint here before function call
		// DEBUG: Step Into (F11) to see the same algorithm process different input
		freq, mostCommon, err := exercise.FreqFromReader(reader)

		// BREAKPOINT: Set breakpoint here after function call
		// DEBUG: Inspect 'freq' map - see counts for your custom input
		// DEBUG: Notice how map size equals number of unique words

		if err != nil {
			log.Fatalf("Error reading custom input: %v", err)
		}

		// Display results in table format
		fmt.Println("Word Frequency Table:")
		fmt.Println("---------------------")

		// BREAKPOINT: Set breakpoint here before loop
		// DEBUG: We're about to iterate over the map
		// DEBUG: Maps in Go are unordered - notice order may vary between runs

		for word, count := range freq {
			// BREAKPOINT: Set breakpoint here inside loop
			// DEBUG: Watch 'word' and 'count' change with each iteration
			// DEBUG: Notice map iteration order is non-deterministic

			marker := ""
			if word == mostCommon {
				marker = " ← most common"
			}
			fmt.Printf("  %-15s : %d%s\n", word, count, marker)
		}

		fmt.Println()
		fmt.Printf("Total unique words: %d\n", len(freq))
		fmt.Printf("Most common word: %q (appears %d times)\n", mostCommon, freq[mostCommon])
		fmt.Println()
	}


	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES FOR MAPS AND SLICES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise.FreqFromReader function
	// 5. Watch Variables panel to see map contents
	// 6. Add watch expressions: len(freq), freq["go"], mostCommon
	// 7. Use Debug Console to inspect: freq["word"], len(freq)
	//
	// UNDERSTANDING MAPS:
	// - In Variables panel, expand 'freq' to see all key-value pairs
	// - Map iteration order is non-deterministic (changes between runs)
	// - Watch how map grows as new words are encountered
	// - Notice len(map) equals number of unique keys
	//
	// UNDERSTANDING IO.READER:
	// - Same FreqFromReader function works with files and strings
	// - This is the power of Go's io.Reader interface
	// - Step Into (F11) to see how bufio.Scanner processes input
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}
