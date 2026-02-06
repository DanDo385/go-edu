//go:build reference

package hellostrings

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps error paths
explicit when an operation can fail, so callers decide how to handle failure.

Read this alongside exercise.go and the tests to understand the intended data
flow, boundary checks, and invariants.
*/

// TitleCase capitalizes the first rune in each whitespace-delimited word.
func TitleCase(s string) string {
	return coreTitleCase(s)
}

// Reverse reverses a string by rune, preserving UTF-8 characters.
func Reverse(s string) string {
	return coreReverse(s)
}

// RuneLen returns the number of runes, not bytes.
func RuneLen(s string) int {
	return coreRuneLen(s)
}

// coreTitleCase demonstrates core logic of capitalizing first letters
// with explicit error handling where failures can occur
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
// with explicit error handling where failures can occur
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

// coreRuneLen demonstrates counting Unicode characters with explicit error handling where failures can occur
//
// Algorithm steps:
// 1. Count the number of Unicode characters in the string
// 2. Return the count (note: different from len(s) for multi-byte characters)
func coreRuneLen(s string) int {
	// Step 1 - Count Unicode runes in string using utf8.RuneCountInString
	// Step 2 - Return the rune count
	return utf8.RuneCountInString(s)
}
