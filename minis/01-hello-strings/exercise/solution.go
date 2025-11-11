//go:build solution
// +build solution

/*
Problem: Implement UTF-8-aware string utilities in Go

We need three functions that correctly handle Unicode text:
1. TitleCase - Capitalize the first letter of each word
2. Reverse - Reverse a string character-by-character (not byte-by-byte!)
3. RuneLen - Count characters (runes), not bytes

Constraints:
- Must handle multi-byte UTF-8 characters (emoji, accented letters, CJK)
- Preserve all characters without corruption
- Use only the Go standard library

Time/Space Complexity:
- TitleCase: O(n) time, O(n) space (allocates new string)
- Reverse: O(n) time, O(n) space (allocates rune slice + result string)
- RuneLen: O(n) time, O(1) space (just counting)

Why Go is well-suited:
- Built-in UTF-8 support: strings are UTF-8 byte sequences by default
- Clear byte/rune distinction: prevents subtle encoding bugs
- Excellent stdlib: `unicode/utf8` and `strings` cover most needs
- Fast: no string copying overhead (immutable strings are shared internally)

Compared to other languages:
- Python: Easier (strings are always Unicode), but slower and less explicit
- JavaScript: Similar ease, but UTF-16 internals cause edge cases (surrogate pairs)
- Rust: More control, but `String`/`&str`/`char` has a steeper learning curve
- C/C++: Requires third-party libraries (ICU), prone to encoding errors
*/

package exercise

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TitleCase converts the first letter of each word to uppercase.
//
// Go Concepts Demonstrated:
// - strings.Fields(): splits on any Unicode whitespace (spaces, tabs, newlines)
// - []rune conversion: allows character-level manipulation
// - unicode.ToUpper(): handles all Unicode uppercase rules (not just ASCII)
// - strings.Builder: efficient string concatenation (reduces allocations)
func TitleCase(s string) string {
	// Split the input string into words using whitespace as delimiter
	// strings.Fields() is preferred over strings.Split() because it:
	// 1. Handles all Unicode whitespace (spaces, tabs, newlines, non-breaking spaces)
	// 2. Automatically trims leading/trailing whitespace
	// 3. Collapses multiple consecutive spaces
	words := strings.Fields(s)

	// Process each word: capitalize first rune, lowercase the rest
	for i, word := range words {
		// Convert string to []rune to work with characters (not bytes)
		// This is critical for UTF-8 correctness:
		// - The string "café" as bytes is [99 97 102 195 169]
		// - As runes it's [99 97 102 233] where 233 is the single character 'é'
		runes := []rune(word)
		if len(runes) > 0 {
			// Capitalize the first rune using unicode.ToUpper()
			// This handles all Unicode case rules (e.g., German ß → SS)
			runes[0] = unicode.ToUpper(runes[0])
			// Note: We don't lowercase the rest to preserve existing capitalization
			// If you wanted "hELLO" → "Hello", add: runes[1:] = unicode.ToLower(...)
		}
		// Convert back to string and update the slice
		words[i] = string(runes)
	}

	// Join words back together with single spaces
	// Alternative: Use strings.Builder for large strings (avoids intermediate allocations)
	return strings.Join(words, " ")
}

// Reverse returns the string reversed character-by-character (UTF-8 aware).
//
// Go Concepts Demonstrated:
// - Rune slices for character-level operations
// - Two-pointer reversal algorithm (in-place swap)
// - String immutability: we must allocate a new string
//
// Three-Input Iteration Table:
//
// Input 1: "Hello" (ASCII-only, happy path)
//   []rune conversion → [H, e, l, l, o]
//   After swap loop  → [o, l, l, e, H]
//   Result: "olleH"
//
// Input 2: "" (empty string, edge case)
//   []rune conversion → []
//   Loop never executes (i=0, j=-1, condition false)
//   Result: ""
//
// Input 3: "Hi👋" (multi-byte emoji)
//   []rune conversion → [H, i, 👋] (3 runes, but 6 bytes!)
//   After swap loop  → [👋, i, H]
//   Result: "👋iH" (correctly preserves emoji)
func Reverse(s string) string {
	// Convert string to slice of runes (Unicode code points)
	// Why? Because reversing bytes would corrupt multi-byte characters:
	// - "👋" in UTF-8 is 4 bytes: [0xF0, 0x9F, 0x91, 0x8B]
	// - Reversing those bytes would produce invalid UTF-8!
	// - But as a rune, it's a single value we can safely move
	runes := []rune(s)

	// Two-pointer reversal: swap elements from both ends moving inward
	// This is the classic O(n/2) in-place reversal algorithm
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		// Simultaneous assignment (Go feature): swap without temp variable
		// Equivalent to: temp := runes[i]; runes[i] = runes[j]; runes[j] = temp
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert []rune back to string
	// This allocates a new string with the UTF-8 encoding of the runes
	return string(runes)
}

// RuneLen returns the number of UTF-8 runes (not bytes) in the string.
//
// Go Concepts Demonstrated:
// - utf8.RuneCountInString(): efficient rune counting
// - Byte vs rune distinction (the most common gotcha for new Go devs!)
//
// Why not just len(s)?
// - len(s) returns the byte count, not character count
// - Example: "café" has len()=5 bytes but 4 characters
// - The 'é' is encoded as 2 bytes in UTF-8: [0xC3, 0xA9]
//
// Three-Input Iteration Table:
//
// Input 1: "hello" (ASCII, happy path)
//   Each ASCII char is 1 byte
//   utf8.RuneCountInString internally iterates: h(1), e(1), l(1), l(1), o(1)
//   Result: 5 runes
//
// Input 2: "" (empty, edge case)
//   No iterations
//   Result: 0 runes
//
// Input 3: "👋😀" (emoji, 2 characters but 8 bytes)
//   👋 = 4 bytes [0xF0, 0x9F, 0x91, 0x8B]
//   😀 = 4 bytes [0xF0, 0x9F, 0x98, 0x80]
//   utf8.RuneCountInString recognizes multi-byte sequences
//   Result: 2 runes (not 8!)
func RuneLen(s string) int {
	// Use the standard library function for correctness and performance
	// This is implemented in assembly on many platforms for speed
	return utf8.RuneCountInString(s)

	// Alternative implementation (educational, less efficient):
	// count := 0
	// for range s {
	//     count++
	// }
	// return count
	// The `for range` loop over a string iterates runes, not bytes!
}

/*
Alternatives & Trade-offs:

1. TitleCase: We could use strings.Title() from the stdlib, but it's deprecated
   as of Go 1.18 because it doesn't follow Unicode word-breaking rules properly.
   For production, use golang.org/x/text/cases.Title() with language-specific rules.

2. Reverse: For ASCII-only strings, byte-level reversal is faster:
   b := []byte(s)
   for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
       b[i], b[j] = b[j], b[i]
   }
   return string(b)
   But this corrupts multi-byte characters. Profile first!

3. RuneLen: Converting to []rune and taking len() works but allocates:
   return len([]rune(s))  // Allocates a slice!
   utf8.RuneCountInString() is O(n) time but O(1) space.

Go vs X:

Go vs Python:
- Python: `s.title()` just works, but encoding surprises can occur (Python 2 vs 3)
- Go: More verbose, but encoding is always explicit and predictable

Go vs JavaScript:
- JS: `s.split('').reverse().join('')` works, but can split surrogate pairs (emoji)
- Go: Rune handling prevents these edge cases

Go vs Rust:
- Rust: `s.chars().rev().collect()` is similar to our []rune approach
- Go: Simpler syntax, but Rust's zero-cost abstractions are faster for large strings
- Both are memory-safe, but Go's GC vs Rust's ownership is a philosophical choice

Go vs C:
- C: Requires ICU library for proper Unicode, easy to introduce buffer overflows
- Go: Standard library + memory safety = fewer production bugs
*/
