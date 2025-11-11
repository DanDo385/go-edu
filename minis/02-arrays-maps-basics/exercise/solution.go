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

Compared to other languages:
- Python: `Counter(words).most_common(1)` is shorter but less explicit
- JavaScript: Requires manual map checking (`map.has()` or `||` patterns)
- Rust: HashMap requires `.entry()` API for insert-or-update
- C: Requires choosing a hash table library; Go's is built-in
*/

package exercise

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
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	// Create frequency map
	// Maps in Go are reference types (like slices), so we use `make()`
	// Alternative: var freq = map[string]int{} (creates empty map)
	freq := make(map[string]int)

	// Create a buffered scanner for line-by-line reading
	// bufio.Scanner is preferred over reading the entire file because:
	// 1. Memory efficient: processes line-by-line (constant memory)
	// 2. Handles different line endings (\n, \r\n) automatically
	// 3. Built-in token splitting (default: ScanLines)
	scanner := bufio.NewScanner(r)

	// Scan() advances to the next line and returns true if successful
	// It returns false when:
	// - End of input is reached
	// - An error occurs (check scanner.Err())
	for scanner.Scan() {
		// Get the current line as a string
		// Text() returns the most recent token (without the newline)
		line := scanner.Text()

		// Normalize the word: remove whitespace and convert to lowercase
		// TrimSpace() removes leading/trailing whitespace (spaces, tabs, newlines)
		// ToLower() handles Unicode case folding (e.g., Turkish İ → i)
		word := strings.ToLower(strings.TrimSpace(line))

		// Skip blank lines
		// This handles both empty lines and lines with only whitespace
		if word == "" {
			continue
		}

		// Increment the frequency count
		// Go maps have zero-value semantics:
		// - If word doesn't exist, freq[word] returns 0
		// - So freq[word]++ sets it to 1 (first occurrence)
		// - Subsequent increments work as expected
		// This is cleaner than: if _, ok := freq[word]; !ok { freq[word] = 0 }
		freq[word]++
	}

	// Check for I/O errors during scanning
	// scanner.Err() returns nil if we reached EOF normally
	// Returns an error if there was a read failure (disk error, closed file, etc.)
	if err := scanner.Err(); err != nil {
		// Return what we have so far + the error
		// The caller can decide whether to use partial results
		return freq, "", err
	}

	// Find the most common word
	// We need to iterate the entire map (no built-in "max" function)
	var maxWord string
	var maxCount int

	// Iterate over the map
	// Important: Map iteration order is RANDOMIZED in Go!
	// This prevents relying on insertion order (unlike Python 3.7+)
	// For deterministic results, we'd need to sort keys first
	for word, count := range freq {
		if count > maxCount {
			maxWord = word
			maxCount = count
		}
		// Note: For ties, we keep the first word encountered
		// Due to random iteration, this is arbitrary but deterministic per run
	}

	// Return results
	// If freq is empty, maxWord will be "" (zero value for string)
	return freq, maxWord, nil
}

/*
Alternatives & Trade-offs:

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
   (See stretch goals!)

4. Case sensitivity:
   Currently we normalize to lowercase. For case-sensitive counting,
   remove the ToLower() call. For case-insensitive but preserve original case,
   use a map[string]string to track original → normalized mapping.

Go vs X:

Go vs Python:
- Python:
    from collections import Counter
    words = [line.strip().lower() for line in f if line.strip()]
    freq = Counter(words)
    most_common = freq.most_common(1)[0][0]
  Pros: One-liner with Counter
  Cons: Requires import; less explicit error handling
  Go's version is more verbose but clearer about what's happening

Go vs JavaScript (Node.js):
- JS:
    const freq = {};
    for (const word of words) {
      freq[word] = (freq[word] || 0) + 1;
    }
  Pros: Similar brevity
  Cons: No zero-value semantics (needs `|| 0` pattern); async file I/O
  Go's synchronous I/O is simpler for line-by-line processing

Go vs Rust:
- Rust:
    let mut freq = HashMap::new();
    *freq.entry(word).or_insert(0) += 1;
  Pros: Zero-cost abstractions; no GC
  Cons: More complex API (.entry().or_insert()); steeper learning curve
  Go's `freq[word]++` is more intuitive for beginners

Go vs Java:
- Java:
    Map<String, Integer> freq = new HashMap<>();
    freq.put(word, freq.getOrDefault(word, 0) + 1);
  Pros: Similar approach
  Cons: More verbose (generic syntax, no operator overloading)
  Go's built-in map literal syntax is cleaner
*/
