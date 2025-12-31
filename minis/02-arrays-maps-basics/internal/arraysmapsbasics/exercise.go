//go:build !solution
// +build !solution

package arraysmapsbasics

// TODO: Import required packages
// You'll need:
// - "bufio" for buffered I/O and line-by-line reading
// - "io" for the Reader interface
// - "strings" for string manipulation (ToLower, TrimSpace)
//
// Uncomment these imports when you implement the functions:
/*
// import (
// 	"bufio"
// 	"io"
// 	"strings"
// )
*/

// FreqFromReader reads words from r (one per line) and returns:
// - A frequency map (word -> count)
// - The most common word (arbitrary choice if there's a tie)
// - An error if reading fails
//
// Example: If input has "hello\nworld\nhello", returns {"hello": 2, "world": 1}, "hello", nil
// TODO: Implement this function
//
// Algorithm:
// 1. Create an empty frequency map to track word counts
// 2. Create a scanner to read lines efficiently from the reader
// 3. For each line, normalize (trim whitespace, convert to lowercase)
// 4. Skip empty lines and increment the count for each word
// 5. Check for scanner errors after reading completes
// 6. Find the word with the highest count by iterating through the map
// 7. Return the frequency map, most common word, and any errors
//
// CS Concepts:
// - Hash Tables: Maps implement hash tables for O(1) average lookup
// - Iteration Order: Maps have random iteration order (not guaranteed sequence)
// - Reference Semantics: Maps are reference types; modifications are visible to caller
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	// TODO: Implement FreqFromReader
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// After implementing:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Compare with solution.go to see detailed explanations
