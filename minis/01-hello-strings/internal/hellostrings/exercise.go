//go:build !solution
// +build !solution

package hellostrings

// TODO: Import required packages
// You'll need:
// - "strings" for string manipulation functions
// - "unicode" for character case conversion
// - "unicode/utf8" for UTF-8 utilities
//
// Uncomment these imports when you implement the functions:

// import (
// 	"strings"
// 	"unicode"
// 	"unicode/utf8"
// )


// TitleCase converts the first letter of each word to uppercase.
// Words are separated by whitespace.
// Example: "hello world" → "Hello World"
// TODO: Implement this function
//
// Algorithm:
// 1. Split the string into words using strings.Fields()
// 2. Iterate over each word and capitalize the first character
// 3. Convert modified words back to strings
// 4. Join all words back into a single string
//
// CS Concept: String immutability - since Go strings are immutable, you must convert to runes (mutable representation), modify, then convert back. This connects to the concept of data types and their properties.
func TitleCase(s string) string {
	// TODO: Implement TitleCase
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Reverse returns the string reversed character-by-character (UTF-8 aware).
// This correctly handles multi-byte characters like emoji.
// Example: "Hello 👋" → "👋 olleH"
// TODO: Implement this function
//
// Algorithm:
// 1. Convert string to a slice of runes (handles multi-byte UTF-8 characters)
// 2. Reverse the rune slice in-place using the two-pointer technique
// 3. Convert the reversed rune slice back to a string
//
// CS Concept: Character encoding - UTF-8 uses variable-length byte sequences. Runes represent Unicode code points. Reversing at the byte level would corrupt multi-byte characters, so we must work at the rune level.
func Reverse(s string) string {
	// TODO: Implement Reverse
	// See solution.reference.go for reference implementation
	panic("not implemented")
}



// RuneLen returns the number of UTF-8 runes (characters) in the string,
// not the byte count. This is important for strings with non-ASCII characters.
// Example: "café" has 4 runes but 5 bytes (é is 2 bytes in UTF-8)
// TODO: Implement this function
//
// Algorithm:
// Count the number of Unicode code points (runes) in the string, not bytes.
// (Note: len(s) returns byte count; range over s iterates by runes)
//
// CS Concept: Encoding and character representation - The distinction between bytes (the physical storage unit) and runes (the logical character unit) is fundamental to understanding how text is represented in modern programming languages.
func RuneLen(s string) int {
	// TODO: Implement RuneLen
	// See solution.reference.go for reference implementation
	panic("not implemented")
}



// After implementing all functions:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Build: go build (should compile without errors)
// - Compare with solution.go to see detailed explanations
