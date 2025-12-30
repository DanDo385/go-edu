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
*/

// coreTitleCase demonstrates core logic of capitalizing first letters
// without error handling
//
// Algorithm steps:
// 1. Split the string into words
// 2. For each word, convert to runes and capitalize first character
// 3. Join words back together with spaces
func coreTitleCase(s string) string {
	// Step 1 - Split input string into words using strings.Fields
	words := strings.Fields(s)

	// Step 2 - Iterate through words, convert each to runes
	for i, word := range words {
		runes := []rune(word)
		// Step 2a - For non-empty words, capitalize first rune using unicode.ToUpper
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		// Step 2b - Convert runes back to string
		words[i] = string(runes)
	}

	// Step 3 - Join all processed words back together with space delimiter
	return strings.Join(words, " ")
}

// coreReverse demonstrates core logic of reversing strings
// without error handling
//
// Algorithm steps:
// 1. Convert string to runes (for proper Unicode handling)
// 2. Use two-pointer technique to swap characters from both ends
// 3. Convert runes back to string
func coreReverse(s string) string {
	// Step 1 - Convert input string to slice of runes
	runes := []rune(s)

	// Step 2 - Use two-pointer technique (i starts at 0, j at len-1)
	//          Swap elements and move pointers until they meet
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Step 3 - Convert runes slice back to string
	return string(runes)
}

// coreRuneLen demonstrates counting Unicode characters without error handling
//
// Algorithm steps:
// 1. Count the number of Unicode characters in the string
// 2. Return the count (note: different from len(s) for multi-byte characters)
func coreRuneLen(s string) int {
	// Step 1 - Count Unicode runes in string using utf8.RuneCountInString
	// Step 2 - Return the rune count
	return utf8.RuneCountInString(s)
}
