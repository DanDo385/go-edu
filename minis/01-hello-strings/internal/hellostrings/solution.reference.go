//go:build reference

package hellostrings

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
Reference Solution - Strings, Runes, and UTF-8
==============================================

First principles (per .cursorrules):

Strings in Go are immutable: the bytes cannot be changed. Any "modification"
creates a new string. Think of a string like a fixed placard — you can read it,
but to "change" it you make a new placard.

Runes vs bytes: A rune is a Unicode code point (int32). " café" has 4 runes
but 5 bytes — é is 2 bytes in UTF-8. len(s) counts bytes; RuneCount counts
logical characters. We use []rune(s) to get a slice of code points for
character-level manipulation.

[]rune(s): Converts string to slice of runes. This COPIES the data — the slice
has its own backing array. Mutating runes[i] does not affect s. string(runes)
creates a new string from the (possibly mutated) slice.
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
//
// Memory and mutability (per .cursorrules):
// s is immutable — we cannot modify s[i]. []rune(s) creates a NEW slice with
// a backing array that holds a copy of the runes. We mutate runes[i], runes[j]
// in place — that's our copy, not the original. string(runes) allocates a new
// string from the mutated slice. The original s is never touched.
func coreReverse(s string) string {
	runes := []rune(s) // Copy: runes has its own backing array

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

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
