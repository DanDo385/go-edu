//go:build !solution
// +build !solution

package exercise

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
Core Logic Without Error Handling
==================================

This file demonstrates the core string manipulation algorithms without
error handling complexity. It's designed to help you understand the
fundamental logic flow.

KEY CONCEPTS:
- Working with Unicode strings and runes
- String immutability in Go
- Converting between strings and runes
- UTF-8 encoding awareness

DEBUGGING:
- Set breakpoints at function entry points (marked with // BREAKPOINT)
- Step through the pure logic without error handling distractions
- Focus on understanding the algorithm
- Watch Variables panel to see string transformations

See /RUN_DEBUG.md for comprehensive debugging guide.
*/

// CoreTitleCase demonstrates the core logic of capitalizing first letters
// BREAKPOINT: Set breakpoint here to step through algorithm
// DEBUG: Watch Variables panel to see how words are processed
func CoreTitleCase(s string) string {
	// DEBUG: Inspect 's' parameter - see input string

	// Split into words
	words := strings.Fields(s)
	// DEBUG: Expand 'words' - see slice of individual words

	// Process each word
	for i, word := range words {
		// Convert to runes to work with Unicode characters
		runes := []rune(word)
		// DEBUG: Expand 'runes' - see Unicode code points

		if len(runes) > 0 {
			// Capitalize first rune
			runes[0] = unicode.ToUpper(runes[0])
			// DEBUG: Inspect runes[0] - now uppercase

			// Convert back to string
			words[i] = string(runes)
		}
	}

	// Join words back together
	result := strings.Join(words, " ")
	// DEBUG: Inspect 'result' - see final capitalized string

	return result
}

// CoreReverse demonstrates the core logic of reversing strings
// BREAKPOINT: Set breakpoint here to step through algorithm
// DEBUG: Watch how runes are swapped in Variables panel
func CoreReverse(s string) string {
	// DEBUG: Inspect 's' parameter - see input string

	// Convert to runes to handle Unicode correctly
	runes := []rune(s)
	// DEBUG: Expand 'runes' - see all characters as code points

	// Reverse using two-pointer technique
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		// DEBUG: Watch 'i' increase and 'j' decrease
		// DEBUG: After swap, inspect runes[i] and runes[j]
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert back to string
	result := string(runes)
	// DEBUG: Compare 's' (input) with 'result' (reversed)

	return result
}

// CoreRuneLen demonstrates counting Unicode characters
// BREAKPOINT: Set breakpoint here to step through algorithm
// DEBUG: Compare byte count vs rune count in Watch panel
func CoreRuneLen(s string) int {
	// DEBUG: Add watch expression: len(s) to see byte count
	// DEBUG: Compare with return value (rune count)

	// Count Unicode characters (not bytes)
	count := utf8.RuneCountInString(s)
	// DEBUG: For ASCII: count == len(s)
	// DEBUG: For Unicode: count < len(s)

	return count
}
