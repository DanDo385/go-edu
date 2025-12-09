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
// - Pass by value vs reference semantics
func TitleCase(s string) string {
	// ============================================================================
	// PARAMETER PASSING: s is passed by VALUE
	// ============================================================================
	// In Go, ALL parameters are passed by value (copy the bytes).
	// However, strings have special semantics:
	// - A string is internally: struct { ptr *byte; len int }
	// - When you pass a string, you copy the pointer + length (cheap!)
	// - You DON'T copy the actual string data (it's shared)
	// - Strings are immutable, so sharing is safe
	//
	// Memory: If s = "hello world", passing it copies ~16 bytes (pointer + len),
	// not 11 bytes of actual data. This is why passing strings is efficient.

	// Split the input string into words using whitespace as delimiter
	// strings.Fields() is preferred over strings.Split() because it:
	// 1. Handles all Unicode whitespace (spaces, tabs, newlines, non-breaking spaces)
	// 2. Automatically trims leading/trailing whitespace
	// 3. Collapses multiple consecutive spaces
	//
	// ============================================================================
	// SLICE CREATION: words is a slice
	// ============================================================================
	// Slices in Go are: struct { ptr *element; len, cap int }
	// The slice 'words' is a REFERENCE to an underlying array
	// - The slice header lives on the stack (if it doesn't escape)
	// - The underlying array lives on the heap (allocated by strings.Fields)
	// - Modifying words[i] changes the underlying array
	// - Passing words to another function passes the slice header BY VALUE
	//   (copies ptr + len + cap, but the ptr still points to same array!)
	words := strings.Fields(s)

	// ============================================================================
	// RANGE LOOP: value semantics
	// ============================================================================
	// range returns (index, value) where:
	// - i is the index (int, copied)
	// - word is a COPY of words[i] (string copied, but data shared as explained above)
	// - Modifying 'word' won't change words[i]
	// - You MUST use words[i] = ... to modify the slice
	for i, word := range words {
		// ========================================================================
		// CONVERSION: string → []rune allocates memory
		// ========================================================================
		// This is critical for UTF-8 correctness:
		// - The string "café" as bytes is [99 97 102 195 169] (5 bytes)
		// - As runes it's [99 97 102 233] (4 runes) where 233 is 'é'
		//
		// Memory allocation:
		// - Go allocates a new []rune on the heap (size = rune count * 4 bytes)
		// - Decodes UTF-8 from string into runes
		// - Returns a slice pointing to this new array
		//
		// The rune slice is: struct { ptr *rune; len, cap int }
		// - 'runes' is the slice header (on stack)
		// - The actual rune array is on the heap
		runes := []rune(word)

		if len(runes) > 0 {
			// ====================================================================
			// INDEXING: Direct modification through slice
			// ====================================================================
			// runes[0] is a rune (int32), which is a VALUE type
			// We're assigning to a slice element, which MODIFIES the underlying array
			//
			// unicode.ToUpper() takes a rune BY VALUE and returns a rune BY VALUE
			// - The input rune is copied (4 bytes, cheap)
			// - A new rune is returned (4 bytes)
			// - No pointers involved, no heap allocation
			runes[0] = unicode.ToUpper(runes[0])
		}

		// ========================================================================
		// CONVERSION: []rune → string allocates memory
		// ========================================================================
		// string(runes) creates a NEW string:
		// - Allocates memory for UTF-8 encoded bytes
		// - Encodes each rune to UTF-8
		// - Returns a new string (immutable)
		//
		// Memory: The old []rune array can now be garbage collected (if not referenced)
		//
		// SLICE MODIFICATION: words[i] = ... updates the underlying array
		// - We're not modifying the slice header
		// - We're changing the data that words[i] points to
		// - This IS visible to anyone holding a reference to the same slice
		words[i] = string(runes)
	}

	// ============================================================================
	// RETURN: strings.Join creates a new string
	// ============================================================================
	// strings.Join(words, " ") concatenates all strings:
	// - Calculates total length needed
	// - Allocates ONE byte array for the result
	// - Copies all strings into it with separators
	// - Returns a new string (immutable)
	//
	// We return this string BY VALUE (copies pointer + len, shares data)
	return strings.Join(words, " ")
}

// Reverse returns the string reversed character-by-character (UTF-8 aware).
//
// Go Concepts Demonstrated:
// - Rune slices for character-level operations
// - Two-pointer reversal algorithm (in-place swap)
// - String immutability: we must allocate a new string
// - Pointer semantics: slices reference underlying arrays
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
	// ============================================================================
	// CONVERSION: string → []rune
	// ============================================================================
	// Why convert? Because reversing bytes would corrupt multi-byte characters:
	// - "👋" in UTF-8 is 4 bytes: [0xF0, 0x9F, 0x91, 0x8B]
	// - Reversing those bytes would produce invalid UTF-8!
	// - But as a rune, it's a single value we can safely move
	//
	// Memory allocation:
	// - Creates a new array on the heap (size = rune count * 4 bytes)
	// - Returns a slice header pointing to it
	// - The slice 'runes' contains: { ptr *rune, len int, cap int }
	runes := []rune(s)

	// ============================================================================
	// TWO-POINTER SWAP: In-place modification
	// ============================================================================
	// This is the classic O(n/2) in-place reversal algorithm
	// Initialize: i = 0 (left pointer), j = len(runes)-1 (right pointer)
	// Loop condition: i < j (stop when pointers meet/cross)
	// Post-iteration: i, j = i+1, j-1 (move pointers inward)
	//
	// Why this works:
	// - We swap elements at positions i and j
	// - After each swap, move i right and j left
	// - When pointers meet, all elements have been swapped
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		// ====================================================================
		// SIMULTANEOUS ASSIGNMENT: Go's tuple assignment
		// ====================================================================
		// This is a Go language feature for swapping without a temp variable
		// How it works:
		// 1. Evaluate RIGHT side: create tuple (runes[j], runes[i])
		// 2. Assign to LEFT side: runes[i] = tuple[0], runes[j] = tuple[1]
		//
		// Why safe? Both values are evaluated BEFORE any assignment
		// Equivalent to:
		//   temp := runes[i]
		//   runes[i] = runes[j]
		//   runes[j] = temp
		//
		// IMPORTANT: This modifies the underlying array!
		// - runes[i] and runes[j] are not pointers
		// - They're direct indexed access to array elements
		// - Assignment changes the array data (which lives on the heap)
		// - The slice 'runes' still points to the same array
		runes[i], runes[j] = runes[j], runes[i]
	}

	// ============================================================================
	// CONVERSION: []rune → string
	// ============================================================================
	// string(runes) creates a NEW string:
	// 1. Allocates a byte array on the heap
	// 2. Encodes each rune to UTF-8 bytes
	// 3. Returns a string struct { ptr *byte, len int }
	//
	// Memory:
	// - The []rune array can now be garbage collected (no longer referenced)
	// - The new string shares no memory with the rune slice
	//
	// Return: We return the string BY VALUE
	// - Copies the string struct (pointer + length, ~16 bytes)
	// - Does NOT copy the actual string data (it's shared immutably)
	return string(runes)
}

// RuneLen returns the number of UTF-8 runes (not bytes) in the string.
//
// Go Concepts Demonstrated:
// - utf8.RuneCountInString(): efficient rune counting
// - Byte vs rune distinction (the most common gotcha for new Go devs!)
// - Pass by value: parameter is copied
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
	// ============================================================================
	// PARAMETER: s is passed by value
	// ============================================================================
	// The string parameter s is passed BY VALUE:
	// - Copies the string struct (pointer to data + length)
	// - Does NOT copy the actual string bytes (they're shared)
	// - This is why passing strings is efficient in Go
	//
	// Example: If caller has str = "hello", calling RuneLen(str):
	// - Copies str's pointer and length (~16 bytes copied)
	// - Both caller's str and our s point to same "hello" data
	// - Safe because strings are immutable (can't be modified)

	// ============================================================================
	// STDLIB FUNCTION: utf8.RuneCountInString
	// ============================================================================
	// Why use this instead of len([]rune(s))?
	// 1. Performance: Doesn't allocate memory (O(n) time, O(1) space)
	// 2. Correctness: Handles malformed UTF-8 gracefully
	// 3. Optimization: Uses assembly on many architectures
	//
	// How it works internally:
	// - Iterates through string bytes
	// - Decodes each UTF-8 sequence
	// - Counts how many runes (code points) exist
	// - Never allocates memory on the heap
	//
	// Alternative implementations (for learning):
	//
	// Method 1: Convert to []rune (allocates!)
	//   return len([]rune(s))
	//   Problem: Allocates a rune slice on the heap (wasteful)
	//
	// Method 2: range loop (no allocation)
	//   count := 0
	//   for range s {
	//       count++
	//   }
	//   return count
	//   Note: range over string iterates RUNES, not bytes!
	//
	// Method 3: Manual decoding (educational)
	//   count := 0
	//   for len(s) > 0 {
	//       _, size := utf8.DecodeRuneInString(s)
	//       count++
	//       s = s[size:]  // Advance by size bytes (creates new string header)
	//   }
	//   return count
	//   Note: s = s[size:] doesn't modify caller's s (we have a copy!)
	//
	// We use utf8.RuneCountInString because it's the most efficient
	// This is implemented in assembly on many platforms for speed
	return utf8.RuneCountInString(s)
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
