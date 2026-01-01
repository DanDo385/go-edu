//go:build reference

/*
Problem: Count word frequencies from text input and find the most common word

Given an io.Reader containing one word per line, we need to:
1. Build a frequency map (word → count)
2. Find the word with the highest count
3. Handle errors gracefully (I/O failures, empty input)

Constraints:
- Normalize to lowercase ("Hello" == "hello")
- Ignore blank lines
- For ties, return any of the tied words (arbitrary but deterministic)

Time/Space Complexity:
- Time: O(n) where n = number of words (one pass to build map, one to find max)
- Space: O(u) where u = number of unique words (map storage)

Why Go is well-suited:
- Built-in maps with clean syntax: `map[string]int` and `count++` patterns
- `io.Reader` interface enables testing without real files
- `bufio.Scanner` handles line-by-line reading efficiently
- Zero value semantics: missing map keys return 0 (perfect for counting!)

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical transformation points
2. Watching variable state changes in the Variables panel
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting map data structure transformations
5. Using the Debug Console to evaluate expressions
6. Understanding the call stack at each level
*/

package arraysmapsbasics

import (
	"bufio"
	"io"
	"strings"
)

// FreqFromReader reads words from r (one per line) and returns frequency data.
//
// Go Concepts Demonstrated:
// - io.Reader interface: Accept any readable source (files, strings, network)
// - bufio.Scanner: Efficient line-by-line reading with automatic buffering
// - Maps: Hash table with zero-value semantics (missing keys return 0)
// - Multiple return values: (result, result, error) is the Go convention
// - strings package: ToLower() and TrimSpace() for normalization
//
// Three-Input Iteration Table:
//
// Input 1: "hello\nworld\nhello\n" (happy path)
//   Line 1: "hello" → freq["hello"] = 1, maxWord = "hello", maxCount = 1
//   Line 2: "world" → freq["world"] = 1 (maxWord unchanged)
//   Line 3: "hello" → freq["hello"] = 2, maxWord = "hello", maxCount = 2
//   Result: map[hello:2 world:1], "hello", nil
//
// Input 2: "" (empty input)
//   No lines scanned
//   Result: map[], "", nil (empty map is valid)
//
// Input 3: "Go\n\ngo\nGO\n" (blank line + case variations)
//   Line 1: "Go" → "go" → freq["go"] = 1
//   Line 2: "" (blank) → skipped
//   Line 3: "go" → freq["go"] = 2
//   Line 4: "GO" → "go" → freq["go"] = 3
//   Result: map[go:3], "go", nil
//
// DEBUGGING WORKFLOW:
// ===================
// 1. Set breakpoint at function entry (line 76)
// 2. Call with test input: FreqFromReader(strings.NewReader("hello\nworld\nhello\n"))
// 3. Step through each line of text processing
// 4. Watch map build up word by word
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	// BREAKPOINT 1: Set here to inspect function entry
	// DEBUG: In Variables panel, expand 'r' to see:
	//   - The io.Reader interface value
	//   - For strings.Reader, you can see the underlying string data
	// DEBUG: In Debug Console, type: r
	// This shows what type of reader was passed in

	// ============================================================================
	// MAP CREATION: Hash table with zero-value semantics
	// ============================================================================
	// Maps in Go are reference types (like slices), so we use `make()`
	// Alternative: var freq = map[string]int{} (creates empty map)
	//
	// Memory Layout:
	// - make() allocates a hash table on the heap
	// - freq is a pointer to this hash table (8 bytes on 64-bit systems)
	// - The map starts with a small number of buckets (8 by default)
	// - Grows automatically when load factor > 6.5 (rehashing occurs)
	//
	// BREAKPOINT 2: Set here BEFORE map creation
	// BREAKPOINT 3: Set here AFTER map creation
	// DEBUG: In Variables panel, expand 'freq' to see:
	//   - Empty map: map[]
	//   - Internal structure is hidden but managed by runtime
	// DEBUG: In Debug Console, type: len(freq)
	// This should be 0 (no elements yet)
	// DEBUG: In Debug Console, type: freq
	// See the empty map representation
	freq := make(map[string]int)

	// ============================================================================
	// SCANNER CREATION: Buffered line-by-line reading
	// ============================================================================
	// Create a buffered scanner for line-by-line reading
	// bufio.Scanner is preferred over reading the entire file because:
	// 1. Memory efficient: processes line-by-line (constant memory)
	// 2. Handles different line endings (\n, \r\n) automatically
	// 3. Built-in token splitting (default: ScanLines)
	//
	// Memory Layout:
	// - Scanner contains a buffer (default 64KB)
	// - Reads chunks from io.Reader, splits into tokens
	// - Returns one token per Scan() call
	//
	// BREAKPOINT 4: Set here BEFORE scanner creation
	// BREAKPOINT 5: Set here AFTER scanner creation
	// DEBUG: In Variables panel, expand 'scanner' to see:
	//   - Internal buffer state
	//   - Current position in the input stream
	// DEBUG: In Debug Console, type: scanner
	// See the scanner structure
	scanner := bufio.NewScanner(r)

	// ============================================================================
	// MAIN LOOP: Process each line of input
	// ============================================================================
	// Scan() advances to the next line and returns true if successful
	// It returns false when:
	// - End of input is reached
	// - An error occurs (check scanner.Err())
	//
	// BREAKPOINT 6: Set here at loop start
	// Step Over (F10) to iterate through each line
	// Watch 'freq' map grow with each iteration
	for scanner.Scan() {
		// BREAKPOINT 7: Set here at start of loop body
		// DEBUG: In Variables panel, watch:
		//   - freq: see how it grows
		//   - Current iteration count (implicit in debugger)
		// DEBUG: In Debug Console, type: len(freq)
		// See how many unique words we've seen so far

		// ====================================================================
		// LINE EXTRACTION: Get current line as string
		// ====================================================================
		// Get the current line as a string
		// Text() returns the most recent token (without the newline)
		//
		// Memory: Text() returns a string that references scanner's buffer
		// Assigning to 'line' makes a copy, so it persists after next Scan()
		//
		// BREAKPOINT 8: Set here BEFORE Text() call
		// BREAKPOINT 9: Set here AFTER Text() call
		// DEBUG: In Variables panel, expand 'line' to see:
		//   - line.len: byte length
		//   - line.data: pointer to string data
		// DEBUG: In Debug Console, type: line
		// See the current line content
		line := scanner.Text()

		// ====================================================================
		// NORMALIZATION: Clean and standardize the word
		// ====================================================================
		// Normalize the word: remove whitespace and convert to lowercase
		// TrimSpace() removes leading/trailing whitespace (spaces, tabs, newlines)
		// ToLower() handles Unicode case folding (e.g., Turkish İ → i)
		//
		// Function composition: strings.ToLower(strings.TrimSpace(line))
		// - Inner function (TrimSpace) executes first
		// - Result is passed to outer function (ToLower)
		//
		// BREAKPOINT 10: Set here BEFORE normalization
		// Step Into (F11) to see TrimSpace implementation
		// BREAKPOINT 11: Set here AFTER normalization
		// DEBUG: In Variables panel, compare:
		//   - line: original value
		//   - word: normalized value
		// DEBUG: In Debug Console, type: line
		// DEBUG: In Debug Console, type: word
		// Compare the before and after normalization
		word := strings.ToLower(strings.TrimSpace(line))

		// ====================================================================
		// BLANK LINE HANDLING: Skip empty words
		// ====================================================================
		// Skip blank lines
		// This handles both empty lines and lines with only whitespace
		//
		// BREAKPOINT 12: Set here at conditional check
		// DEBUG: In Variables panel, watch 'word'
		// DEBUG: In Debug Console, type: word == ""
		// See if this line will be skipped
		// DEBUG: In Debug Console, type: len(word)
		// Length of 0 means empty string
		if word == "" {
			// BREAKPOINT 13: Set here inside if block
			// This only executes for blank lines
			// Step Over (F10) to continue to next iteration
			continue
		}

		// ====================================================================
		// MAP INCREMENT: Core counting logic with zero-value semantics
		// ====================================================================
		// Increment the frequency count
		// Go maps have zero-value semantics:
		// - If word doesn't exist, freq[word] returns 0
		// - So freq[word]++ sets it to 1 (first occurrence)
		// - Subsequent increments work as expected
		// This is cleaner than: if _, ok := freq[word]; !ok { freq[word] = 0 }
		//
		// Under the Hood:
		// 1. Hash the word string to find bucket
		// 2. Look up word in the bucket (linked list or similar)
		// 3. If not found, return 0 (zero value for int)
		// 4. Increment: 0++ = 1 or existingCount++ = existingCount+1
		// 5. Store back in the map
		//
		// BREAKPOINT 14: Set here BEFORE increment
		// DEBUG: In Variables panel, expand 'freq' map
		// Look for the current word in the map
		// DEBUG: In Debug Console, type: freq[word]
		// See current count (0 if first occurrence)
		// DEBUG: In Debug Console, type: word
		// See which word we're counting
		// BREAKPOINT 15: Set here AFTER increment
		// DEBUG: In Variables panel, expand 'freq' map again
		// See the updated count
		// DEBUG: In Debug Console, type: freq[word]
		// Verify the count increased by 1
		// DEBUG: Watch how the map grows dynamically
		freq[word]++

		// BREAKPOINT 16: Set here at end of loop iteration
		// DEBUG: In Variables panel, review:
		//   - word: current word processed
		//   - freq: entire frequency map so far
		// DEBUG: In Debug Console, type: freq
		// See all word counts accumulated so far
		// Step Over (F10) to continue to next iteration
	}

	// BREAKPOINT 17: Set here AFTER loop completes
	// DEBUG: In Variables panel, expand 'freq' map
	// All words should be counted now
	// DEBUG: In Debug Console, type: freq
	// See the complete frequency map
	// DEBUG: In Debug Console, type: len(freq)
	// See how many unique words were found

	// ============================================================================
	// ERROR CHECKING: Distinguish EOF from I/O errors
	// ============================================================================
	// Check for I/O errors during scanning
	// scanner.Err() returns nil if we reached EOF normally
	// Returns an error if there was a read failure (disk error, closed file, etc.)
	//
	// BREAKPOINT 18: Set here BEFORE error check
	// DEBUG: In Variables panel, watch 'err' variable (if it exists)
	// DEBUG: In Debug Console, type: scanner.Err()
	// See if any error occurred during scanning
	// Step Into (F11) to see scanner.Err() implementation
	if err := scanner.Err(); err != nil {
		// BREAKPOINT 19: Set here inside error block
		// This only executes if there was a scanner error
		// DEBUG: In Variables panel, expand 'err' to see:
		//   - Error message
		//   - Error type
		// DEBUG: In Debug Console, type: err.Error()
		// See the error message string

		// Return what we have so far + the error
		// The caller can decide whether to use partial results
		return freq, "", err
	}

	// ============================================================================
	// FINDING MAXIMUM: Linear scan through map
	// ============================================================================
	// Find the most common word
	// We need to iterate the entire map (no built-in "max" function)
	//
	// Algorithm: Linear scan with running maximum
	// - Initialize maxWord = "" (empty), maxCount = 0
	// - For each (word, count) in map:
	//   - If count > maxCount, update both maxWord and maxCount
	// - After loop, maxWord contains the most frequent word
	//
	// Time Complexity: O(n) where n = number of unique words
	//
	// BREAKPOINT 20: Set here BEFORE variable initialization
	// BREAKPOINT 21: Set here AFTER variable initialization
	// DEBUG: In Variables panel, watch:
	//   - maxWord: should be "" (empty string, zero value)
	//   - maxCount: should be 0 (zero value for int)
	// DEBUG: In Debug Console, type: maxWord
	// DEBUG: In Debug Console, type: maxCount
	// Both should be zero values initially
	var maxWord string
	var maxCount int

	// ============================================================================
	// MAP ITERATION: Randomized order iteration
	// ============================================================================
	// Iterate over the map
	// Important: Map iteration order is RANDOMIZED in Go!
	// This prevents relying on insertion order
	// For deterministic results, we'd need to sort keys first
	//
	// BREAKPOINT 22: Set here at loop start
	// Step Over (F10) to iterate through each map entry
	// Watch maxWord and maxCount update when a higher count is found
	for word, count := range freq {
		// BREAKPOINT 23: Set here at start of loop body
		// DEBUG: In Variables panel, watch:
		//   - word: current word being examined
		//   - count: frequency of current word
		//   - maxWord: best word found so far
		//   - maxCount: highest count found so far
		// DEBUG: In Debug Console, type: word
		// DEBUG: In Debug Console, type: count
		// DEBUG: In Debug Console, type: maxWord
		// DEBUG: In Debug Console, type: maxCount

		// ====================================================================
		// COMPARISON: Check if current count is greater than max
		// ====================================================================
		// BREAKPOINT 24: Set here BEFORE comparison
		// DEBUG: In Debug Console, type: count > maxCount
		// See if this word will become the new maximum
		if count > maxCount {
			// BREAKPOINT 25: Set here inside if block
			// This only executes when we find a new maximum
			// DEBUG: In Variables panel, watch values about to change
			// DEBUG: In Debug Console, type: word
			// This will become the new maxWord

			maxWord = word
			maxCount = count

			// BREAKPOINT 26: Set here AFTER update
			// DEBUG: In Variables panel, verify:
			//   - maxWord now equals word
			//   - maxCount now equals count
			// DEBUG: In Debug Console, type: maxWord == word
			// Should be true
			// DEBUG: In Debug Console, type: maxCount == count
			// Should be true
		}
		// Note: For ties, we keep the first word encountered
		// Due to random iteration, this is arbitrary but deterministic per run

		// BREAKPOINT 27: Set here at end of loop iteration
		// DEBUG: In Variables panel, review current state
		// DEBUG: In Debug Console, type: maxWord
		// DEBUG: In Debug Console, type: maxCount
		// See the best word found so far
		// Step Over (F10) to continue to next iteration
	}

	// BREAKPOINT 28: Set here AFTER loop completes
	// DEBUG: In Variables panel, verify:
	//   - maxWord: the most frequent word
	//   - maxCount: its frequency
	// DEBUG: In Debug Console, type: maxWord
	// DEBUG: In Debug Console, type: maxCount
	// DEBUG: In Debug Console, type: freq[maxWord]
	// Verify that freq[maxWord] == maxCount

	// ============================================================================
	// RETURN: Multi-value return
	// ============================================================================
	// Return results
	// If freq is empty, maxWord will be "" (zero value for string)
	//
	// BREAKPOINT 29: Set here at return statement
	// DEBUG: In Variables panel, review all return values:
	//   - freq: complete frequency map
	//   - maxWord: most common word
	//   - nil: no error occurred
	// DEBUG: In Debug Console, type: freq
	// DEBUG: In Debug Console, type: maxWord
	// DEBUG: In Debug Console, type: len(freq)
	// Final verification before returning
	return freq, maxWord, nil
}

/*
ADVANCED DEBUGGING TECHNIQUES:
===============================

1. CONDITIONAL BREAKPOINTS:
   Right-click on any breakpoint and add conditions
   Examples:
   - In main loop: word == "hello" (break only for specific word)
   - In main loop: len(freq) > 5 (break only after 5 unique words)
   - In max finding loop: count > maxCount (break only when new max found)

2. LOGPOINTS:
   Right-click on line number → Add Logpoint
   Examples:
   - In main loop: Log "Processing word: {word}, count: {freq[word]}"
   - In max finding loop: Log "Comparing {word}({count}) vs {maxWord}({maxCount})"
   - No need to modify code or add print statements!

3. DEBUG CONSOLE EXPRESSIONS:
   During debugging, try these in the Debug Console:
   - Type: freq to see current frequency map
   - Type: len(freq) to see number of unique words
   - Type: freq["hello"] to check count of specific word
   - Type: maxWord to see most frequent word so far
   - Type: maxCount to see its frequency

4. WATCH EXPRESSIONS:
   Add these to the Watch panel for real-time monitoring:
   - len(freq) - number of unique words
   - freq[word] - count of current word
   - maxCount - highest frequency seen
   - count > maxCount - will this become new max?

5. CALL STACK NAVIGATION:
   In the Call Stack panel:
   - Click different frames to see state at each level
   - Observe how variables change across stack frames
   - Use "Step Out" (Shift+F11) to return to caller

6. MEMORY INSPECTION:
   To see memory allocations:
   - Watch how freq map grows
   - Notice when map rehashes (internal capacity increases)
   - Compare memory addresses of different strings

7. STEP COMMANDS:
   - F10 (Step Over): Execute line, don't enter functions
   - F11 (Step Into): Enter function calls to see internals
   - Shift+F11 (Step Out): Return to caller
   - F5 (Continue): Run until next breakpoint

8. DATA BREAKPOINTS:
   Watch for when specific variables change:
   - Right-click variable → Break When Value Changes
   - Useful for tracking when freq map is modified
   - Useful for tracking when maxWord updates

Alternatives & Trade-offs:
==========================

1. Reading the entire file at once:
   data, _ := io.ReadAll(r)
   words := strings.Split(string(data), "\n")
   Pros: Simpler code
   Cons: O(n) memory usage; fails on large files; requires converting bytes to string

2. Using a struct for return values:
   type Result struct { Freq map[string]int; MostCommon string }
   Pros: Named fields (self-documenting)
   Cons: More verbose; Go prefers multiple returns for simple cases

3. Returning top N words instead of just one:
   We'd need to sort the map by value:
   - Convert to []struct{word string; count int}
   - Use sort.Slice() with custom comparator
   - Return top N elements

4. Case sensitivity:
   Currently we normalize to lowercase. For case-sensitive counting,
   remove the ToLower() call. For case-insensitive but preserve original case,
   use a map[string]string to track original → normalized mapping.

5. Using sync.Map for concurrent access:
   var freq sync.Map
   freq.LoadOrStore(word, 0)
   freq.Store(word, count+1)
   Pros: Thread-safe for concurrent reads/writes
   Cons: More complex API; overkill for single-goroutine code

DEBUGGING EXERCISES:
====================

Exercise 1: Trace FreqFromReader execution
- Input: strings.NewReader("hello\nworld\nhello\n")
- Set breakpoints at: 84, 135, 148, 172, 245, 326
- Watch: freq, maxWord, maxCount
- Question: How many times does the main loop iterate?
- Question: What's the value of freq after processing all lines?
- Question: What is maxWord and maxCount at the end?

Exercise 2: Understand map zero-value semantics
- Input: strings.NewReader("go\n")
- Set breakpoints at: 235, 245
- In Debug Console before line 245: freq["go"]
- Question: What does freq["go"] return before we increment it?
- Question: Why don't we need to check if the key exists?

Exercise 3: Watch map growth
- Input: strings.NewReader("a\nb\nc\nd\ne\n")
- Set breakpoints at: 112, 245
- Watch: freq map in Variables panel
- Question: How does the map display change as we add more keys?
- Question: Can you see when the map internally rehashes?

Exercise 4: Empty input edge case
- Input: strings.NewReader("")
- Set breakpoints at: 84, 148, 276, 398
- Question: Does the main loop execute? How many times?
- Question: What are the values of maxWord and maxCount?
- Question: What does the function return?

Exercise 5: Blank lines handling
- Input: strings.NewReader("hello\n\n\nworld\n")
- Set breakpoints at: 172, 208, 213, 245
- Watch: word variable
- Question: How many times does the continue statement execute?
- Question: What is len(freq) after processing all input?

Exercise 6: Finding maximum with ties
- Input: strings.NewReader("a\nb\na\nb\n")
- Set breakpoints at: 326, 344, 351
- Watch: maxWord, maxCount in max finding loop
- Question: Which word becomes maxWord? (Remember: random iteration!)
- Question: What is maxCount at the end?
*/
