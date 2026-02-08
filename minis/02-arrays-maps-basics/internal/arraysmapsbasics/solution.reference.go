//go:build reference

package arraysmapsbasics

import (
	"bufio"
	"io"
	"strings"
)

/*
Reference Solution - Arrays, Maps, and Basic Data Structures
==========================================================

This file demonstrates fundamental Go data structures: arrays and maps.
Arrays provide fixed-size, contiguous memory for ordered data access.
Maps provide key-value storage with hash-based lookups.

This connects to the broader Go ecosystem by showing:
- How Go's built-in types (slices, maps) eliminate manual memory management
- Why interfaces like io.Reader enable composable, testable code
- How bufio.Scanner abstracts streaming I/O into line-by-line processing

The exercise builds understanding of:
- Value semantics: how data copying affects performance and correctness
- Hash table behavior: O(1) average lookups vs O(n) worst-case scenarios
- Memory layout: how arrays provide cache-friendly sequential access

Teaching notes (per .cursorrules):
- Memory/ownership: A map value is a reference type — it holds a pointer to the
  hash table. When we return freq, we return that descriptor. The caller gets
  the SAME underlying hash table. Mutations (freq[word]++) affect the shared
  data. This is NOT a copy of the map contents — it's a copy of the map header
  that points to the same table.
- Invariants: establish data structure assumptions early, then build algorithms
  that can rely on those guarantees.
- Error surfaces: I/O operations can fail, so we design APIs that surface errors
  explicitly rather than panicking, allowing callers to implement retry logic.
*/

/*
FreqFromReader - Public API Boundary

This function serves as the public interface for word frequency counting.
It wraps the core logic while providing an error return for API consistency.
In Go, it's common to have thin wrapper functions that handle error propagation
while keeping complex logic in focused helper functions.

The signature (map[string]int, string, error) shows Go's multiple return values:
- First return: the frequency map (mutable reference)
- Second return: most common word (immutable string)
- Third return: error (nil means success)

This pattern allows callers to check errors first, then safely use the data.
*/
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	// Call the core implementation - no errors can occur in this simplified version
	// In production code, this would handle scanner errors from the reader
	freq, mostCommon := CoreFreqFromReader(r)
	return freq, mostCommon, nil
}

/*
CoreFreqFromReader - Word Frequency Counting Algorithm

This function demonstrates the complete algorithm for counting word frequencies
from streaming text input. It combines several fundamental concepts:

Algorithm overview:
1. Initialize empty frequency map (hash table for O(1) lookups)
2. Stream text line-by-line using bufio.Scanner (efficient I/O)
3. Normalize each line (trim whitespace, convert to lowercase)
4. Update frequency counts in the map
5. Find the most frequent word through linear scan

Memory considerations:
- Map grows dynamically as new words are encountered
- Scanner buffers input for efficient reading
- No premature optimization - we let Go's runtime handle growth

Error handling: This core function assumes the reader works correctly.
In production, we'd add scanner.Err() checks and return errors.
*/
func CoreFreqFromReader(r io.Reader) (map[string]int, string) {
	// Step 1: Initialize empty frequency map
	// make(map[string]int) creates a hash table that maps strings to integers
	// The map starts empty but will grow as we add key-value pairs
	// This is more efficient than pre-allocating when we don't know the size
	freq := make(map[string]int)

	// Step 2: Create scanner for line-by-line reading
	// bufio.Scanner wraps the io.Reader and provides Scan() method
	// It automatically handles buffering and line splitting (\n, \r\n, \r)
	// This is more efficient than reading the entire file into memory
	scanner := bufio.NewScanner(r)

	// Step 3: Process each line until EOF or error
	// scanner.Scan() returns true if there's more data, false at EOF
	// Each call advances to the next line, accessible via scanner.Text()
	for scanner.Scan() {
		// Get the current line as a string
		line := scanner.Text()

		// Step 3a: Normalize the line
		// strings.TrimSpace removes leading/trailing whitespace
		// strings.ToLower converts to lowercase for case-insensitive counting
		// This ensures "Word" and "word" are counted as the same
		word := strings.ToLower(strings.TrimSpace(line))

		// Step 3b: Skip empty lines
		// After trimming, empty strings indicate blank lines
		// We continue to the next iteration rather than counting empty words
		if word == "" {
			continue
		}

		// Step 4: Increment frequency count
		// freq[word]++ is syntactic sugar for freq[word] = freq[word] + 1
		// If the key doesn't exist, Go automatically initializes it to the zero value (0)
		// This is why maps are perfect for frequency counting - no manual initialization needed
		freq[word]++
	}

	// Step 5: Find the most common word
	// We need to scan the entire map to find the maximum frequency
	// There are no guarantees about map iteration order in Go

	// Initialize tracking variables
	mostCommon := ""  // Empty string initially - we'll update this
	maxCount := 0     // Start with 0 - any word with count > 0 will be better

	// Iterate through all word-frequency pairs
	// range over map returns (key, value) pairs in random order
	for word, count := range freq {
		// Check if this word should become the new "most common"
		// First condition: higher count than current maximum
		// Second condition (tied counts): lexicographically smaller word wins
		// This makes the result deterministic when there are ties
		if count > maxCount || (count == maxCount && (mostCommon == "" || word < mostCommon)) {
			// Update our tracking variables
			mostCommon = word
			maxCount = count
		}
	}

	// Step 6: Return results
	// Return both the complete frequency map and the most common word
	// The map is returned by reference, so the caller can modify it if needed
	// The string is returned by value (copied), so it's immutable to the caller
	return freq, mostCommon
}
