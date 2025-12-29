//go:build ignore
// +build ignore

package main

import (
	"flag"   // Command-line argument parsing
	"fmt"    // Formatted I/O
	"strings" // String manipulation
	"unicode" // Unicode character utilities
	"unicode/utf8" // UTF-8 encoding/decoding
)

/*
===================================================================
String Manipulation Demonstration
===================================================================

This program demonstrates fundamental string operations in Go:
- Working with Unicode strings
- Converting between strings and runes
- String reversal and case manipulation

USAGE:
    go run main.go
    go run main.go -text="hello world"
    go run main.go -text="café" -demo=reverse

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging
- Watch Variables panel to see string transformations
- See /RUN_DEBUG.md for comprehensive debugging guide

NOTE: This file has build tag 'ignore' so it doesn't conflict with
the exercise package in the same directory. Run it directly with:
    go run main.go
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect default values before parsing
	// DEBUG: Before flag.Parse(), all flags have their default values
	// DEBUG: In Variables panel, you'll see: text="hello world", demo="all"

	// flag.String creates a string flag variable
	// Syntax: flag.String(name, default, usage)
	text := flag.String("text", "hello world", "Text to demonstrate with")
	demo := flag.String("demo", "all", "Which demo to run: all, titlecase, reverse, runelen")

	flag.Parse()  // Parse command-line arguments

	// BREAKPOINT: Set breakpoint here to inspect parsed values in Variables panel
	// DEBUG: After flag.Parse(), variables now hold actual command-line values
	// DEBUG: Watch 'text' and 'demo' variables - they may have changed from defaults

	// ============================================================================
	// STEP 2: Example 1 - TitleCase (Simple)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through simple example
	// DEBUG: This is the simplest use case - capitalizing first letter of each word
	// DEBUG: Step Over (F10) to execute line by line
	// DEBUG: Watch Variables panel to see how 'result' changes

	if *demo == "all" || *demo == "titlecase" {
		fmt.Println("=== TitleCase Demo ===")
		fmt.Println("Capitalizes the first letter of each word")
		fmt.Println()

		testCases := []string{
			"hello world",
			"the quick brown fox",
			"Hello 👋 World",
			"café résumé",
			"日本語",
			*text,  // User-provided text
		}

		// BREAKPOINT: Set breakpoint here before loop
		// DEBUG: Expand 'testCases' in Variables panel to see all test strings
		// DEBUG: Notice the len and cap of the slice

		for _, s := range testCases {
			// BREAKPOINT: Set breakpoint here inside loop
			// DEBUG: Watch 's' variable change with each iteration
			// DEBUG: Step Into (F11) on titleCase() to see function internals

			result := titleCase(s)

			// BREAKPOINT: Set breakpoint here after function call
			// DEBUG: Inspect 'result' variable - see how string was transformed
			// DEBUG: Compare 's' (input) with 'result' (output)

			fmt.Printf("%-25s → %s\n", s, result)
		}
		fmt.Println()
	}

	// ============================================================================
	// STEP 3: Example 2 - Reverse (Intermediate)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through intermediate example
	// DEBUG: This shows more advanced features - working with runes
	// DEBUG: Step Into (F11) on reverse() to see how runes are processed

	if *demo == "all" || *demo == "reverse" {
		fmt.Println("=== Reverse Demo ===")
		fmt.Println("Reverses strings (handles Unicode correctly)")
		fmt.Println()

		testCases := []string{
			"hello",
			"Hello 👋 World",
			"café",
			"日本語",
			*text,
		}

		for _, s := range testCases {
			result := reverse(s)
			fmt.Printf("%-25s → %s\n", s, result)
		}
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 3 - RuneLen (Advanced)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through advanced example
	// DEBUG: This demonstrates byte vs rune counting
	// DEBUG: Use Call Stack panel to see function call hierarchy

	if *demo == "all" || *demo == "runelen" {
		fmt.Println("=== RuneLen Demo ===")
		fmt.Println("Shows difference between byte length and rune (character) length")
		fmt.Println()

		testCases := []string{
			"hello",           // ASCII only
			"café",            // Contains é (2 bytes)
			"Hello 👋 World",  // Contains emoji (4 bytes)
			"日本語",           // Japanese (3 characters, 9 bytes)
			*text,
		}

		for _, s := range testCases {
			byteLen := len(s)  // len() returns byte length
			runeLen := runeLen(s)  // Count actual characters

			// BREAKPOINT: Set breakpoint here
			// DEBUG: Compare byteLen vs runeLen
			// DEBUG: For ASCII, they're equal. For Unicode, runeLen < byteLen

			fmt.Printf("%-25s → bytes: %2d, runes: %2d\n", s, byteLen, runeLen)
		}
	}
}

// ============================================================================
// Helper Functions with Detailed Comments
// ============================================================================

// titleCase converts the first letter of each word to uppercase
// BREAKPOINT: Set breakpoint at function entry to inspect parameters
func titleCase(s string) string {
	// DEBUG: Inspect 's' parameter in Variables panel

	// Step 1: Split string into words
	// strings.Fields splits on whitespace (spaces, tabs, newlines)
	// Memory: Allocates a new slice []string to store words
	words := strings.Fields(s)

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Expand 'words' in Variables panel
	// DEBUG: See: words ([]string) with len and individual elements

	// Step 2: Process each word
	for i, word := range words {
		// BREAKPOINT: Set breakpoint here inside loop
		// DEBUG: Watch 'i' (index) and 'word' change with each iteration

		// Convert string to []rune to work with Unicode characters
		// A rune is an int32 representing a Unicode code point
		// Memory: Allocates new slice []rune on heap
		runes := []rune(word)

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Expand 'runes' in Variables panel
		// DEBUG: See each rune value (Unicode code points)

		if len(runes) > 0 {
			// Capitalize first rune
			// unicode.ToUpper works correctly for all Unicode, not just ASCII
			runes[0] = unicode.ToUpper(runes[0])

			// BREAKPOINT: Set breakpoint here
			// DEBUG: Inspect runes[0] - it should now be uppercase

			// Convert []rune back to string
			// Memory: Allocates new string, encodes UTF-8
			words[i] = string(runes)
		}
	}

	// Step 3: Join words back together
	// strings.Join concatenates with separator " "
	// Memory: Allocates new string for result
	result := strings.Join(words, " ")

	// BREAKPOINT: Set breakpoint here before return
	// DEBUG: Inspect 'result' - see final transformed string

	return result
}

// reverse reverses a string (Unicode-aware)
// BREAKPOINT: Set breakpoint at function entry
func reverse(s string) string {
	// DEBUG: Inspect 's' parameter

	// Convert to []rune to handle Unicode correctly
	// If we used bytes, we'd break multi-byte UTF-8 characters
	runes := []rune(s)

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Expand 'runes' to see all characters as code points

	// Reverse the slice in place
	// This is the classic "two-pointer" algorithm
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		// BREAKPOINT: Set breakpoint here inside loop
		// DEBUG: Watch 'i' increase and 'j' decrease
		// DEBUG: After swap, inspect runes[i] and runes[j]

		// Swap elements
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert back to string
	result := string(runes)

	// BREAKPOINT: Set breakpoint here before return
	// DEBUG: Compare 's' with 'result' - string should be reversed

	return result
}

// runeLen returns the number of runes (characters) in a string
// BREAKPOINT: Set breakpoint at function entry
func runeLen(s string) int {
	// DEBUG: Inspect 's' parameter

	// utf8.RuneCountInString counts Unicode code points
	// This is more accurate than len(s) for non-ASCII strings
	count := utf8.RuneCountInString(s)

	// BREAKPOINT: Set breakpoint here before return
	// DEBUG: Compare count with len(s) in Watch panel
	// DEBUG: For ASCII: count == len(s)
	// DEBUG: For Unicode: count < len(s)

	return count
}
