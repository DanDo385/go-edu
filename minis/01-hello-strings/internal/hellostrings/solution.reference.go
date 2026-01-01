//go:build reference

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

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical transformation points
2. Watching variable state changes in the Variables panel
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting memory allocations and data structure transformations
5. Using the Debug Console to evaluate expressions
6. Understanding the call stack at each level
*/

package hellostrings

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
//
// DEBUGGING WORKFLOW:
// ===================
// 1. Set breakpoint at function entry (line 66)
// 2. Call with test input: TitleCase("hello café 世界")
// 3. Step through each transformation
// 4. Watch memory allocations in real-time
func TitleCase(s string) string {
	// BREAKPOINT 1: Set here to inspect function entry
	// DEBUG: In Variables panel, expand 's' to see:
	//   - s.len: byte length (not character count!)
	//   - s.data: pointer to underlying byte array
	// DEBUG: In Debug Console, type: len(s)
	// DEBUG: In Debug Console, type: utf8.RuneCountInString(s)
	// Compare byte length vs rune count for UTF-8 strings

	// BREAKPOINT 2: Set here BEFORE strings.Fields call
	// Step Into (F11) to see how strings.Fields works internally
	// Watch how it identifies Unicode whitespace characters
	words := strings.Fields(s)

	// BREAKPOINT 3: Set here AFTER strings.Fields call
	// DEBUG: In Variables panel, expand 'words' slice to see:
	//   - words.len: number of words found
	//   - words.cap: capacity of underlying array
	//   - words[0], words[1], etc.: each word as a string
	// DEBUG: In Debug Console, type: len(words)
	// DEBUG: In Debug Console, type: cap(words)
	// DEBUG: In Debug Console, type: words[0]
	// Notice how slice capacity may exceed length (growth strategy)

	// ============================================================================
	// RANGE LOOP: value semantics
	// ============================================================================
	// range returns (index, value) where:
	// - i is the index (int, copied)
	// - word is a COPY of words[i] (string copied, but data shared)
	// - Modifying 'word' won't change words[i]
	// - You MUST use words[i] = ... to modify the slice
	//
	// BREAKPOINT 4: Set here at loop start
	// Step Over (F10) to iterate through each word
	// Watch 'i' and 'word' change with each iteration
	for i, word := range words {
		// BREAKPOINT 5: Set here at start of loop body
		// DEBUG: In Variables panel, watch:
		//   - i: current loop index
		//   - word: current word (a copy of words[i])
		// DEBUG: In Debug Console, type: words[i] == word
		// This should be 'true' - they have the same value

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
		//
		// BREAKPOINT 6: Set here BEFORE []rune conversion
		// DEBUG: In Variables panel, look at 'word':
		//   - word: still a string
		//   - Note its byte representation
		// BREAKPOINT 7: Set here AFTER []rune conversion
		// DEBUG: In Variables panel, expand 'runes' to see:
		//   - runes.len: number of runes (characters)
		//   - runes.cap: capacity of underlying array
		//   - runes[0], runes[1], etc.: each rune as int32
		// DEBUG: In Debug Console, type: len(word)
		// DEBUG: In Debug Console, type: len(runes)
		// Compare byte length of word vs rune count
		// DEBUG: In Debug Console, type: runes[0]
		// See the numeric Unicode code point value
		runes := []rune(word)

		// BREAKPOINT 8: Set here BEFORE length check
		// DEBUG: Watch Variables panel for 'runes.len'
		// Step Over (F10) - if len > 0, we enter the if block
		// If len == 0, we skip the if block (empty word edge case)
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
			//
			// BREAKPOINT 9: Set here BEFORE unicode.ToUpper call
			// DEBUG: In Variables panel, note current value of runes[0]
			// DEBUG: In Debug Console, type: runes[0]
			// DEBUG: In Debug Console, type: string(runes[0])
			// See the character representation
			// Step Into (F11) to see unicode.ToUpper implementation
			// BREAKPOINT 10: Set here AFTER uppercase assignment
			// DEBUG: In Variables panel, see runes[0] is now uppercase
			// DEBUG: In Debug Console, type: runes[0]
			// DEBUG: In Debug Console, type: string(runes[0])
			// Compare before/after values - observe the transformation
			// DEBUG: Watch how lowercase rune value changes to uppercase
			// Example: 'h' (104) becomes 'H' (72)
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
		//
		// BREAKPOINT 11: Set here BEFORE string conversion
		// DEBUG: In Variables panel, expand 'runes' slice
		// See the array of int32 values
		// BREAKPOINT 12: Set here AFTER words[i] assignment
		// DEBUG: In Variables panel, expand 'words' slice
		// DEBUG: Look at words[i] - it now has the capitalized word
		// DEBUG: In Debug Console, type: words[i]
		// DEBUG: In Debug Console, type: string(runes)
		// They should be equal
		// DEBUG: Watch how the []rune array is encoded back to UTF-8 string
		words[i] = string(runes)

		// BREAKPOINT 13: Set here at end of loop iteration
		// DEBUG: In Variables panel, review the transformation:
		//   - Original 'word' variable is unchanged (it was a copy)
		//   - words[i] now contains the capitalized version
		// DEBUG: In Debug Console, type: word
		// DEBUG: In Debug Console, type: words[i]
		// Notice word != words[i] after modification
		// Step Over (F10) to continue to next iteration
	}

	// BREAKPOINT 14: Set here AFTER loop completes
	// DEBUG: In Variables panel, expand 'words' slice
	// All elements should now be title-cased
	// DEBUG: In Debug Console, type: words
	// See the entire transformed slice

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
	//
	// BREAKPOINT 15: Set here BEFORE strings.Join
	// Step Into (F11) to see strings.Join implementation
	// Watch how it calculates buffer size and concatenates
	// BREAKPOINT 16: Set here at return statement
	// DEBUG: In Variables panel, see the final result string
	// DEBUG: In Debug Console, type: strings.Join(words, " ")
	// This is the final output
	// DEBUG: Compare input 's' with output - transformation complete
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
//
// DEBUGGING WORKFLOW:
// ===================
// 1. Set breakpoint at function entry (line 235)
// 2. Call with: Reverse("Hi👋")
// 3. Watch the two-pointer swap algorithm work
// 4. Observe memory transformations step by step
func Reverse(s string) string {
	// BREAKPOINT 17: Set here to inspect function entry
	// DEBUG: In Variables panel, expand 's' to see:
	//   - s.len: byte length
	//   - s.data: pointer to byte array
	// DEBUG: In Debug Console, type: len(s)
	// DEBUG: In Debug Console, type: utf8.RuneCountInString(s)
	// For "Hi👋", len=6 bytes but 3 runes!

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
	//
	// BREAKPOINT 18: Set here BEFORE []rune conversion
	// DEBUG: In Variables panel, look at 's' as bytes
	// For UTF-8 strings, you'll see multi-byte sequences
	// BREAKPOINT 19: Set here AFTER []rune conversion
	// DEBUG: In Variables panel, expand 'runes' slice to see:
	//   - runes.len: number of characters
	//   - runes[0], runes[1], etc.: each character as int32
	// DEBUG: In Debug Console, type: len(runes)
	// DEBUG: In Debug Console, type: cap(runes)
	// DEBUG: For "Hi👋", you should see len=3, cap=3
	// DEBUG: In Debug Console, type: runes[2]
	// This is the emoji as a single int32 value (128075 for 👋)
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
	//
	// BREAKPOINT 20: Set here BEFORE loop initialization
	// DEBUG: In Variables panel, note runes in original order
	// DEBUG: In Debug Console, type: len(runes)
	// Calculate expected iterations: len(runes) / 2
	// BREAKPOINT 21: Set here at loop condition check
	// DEBUG: Watch Variables panel for 'i' and 'j' values
	// On first iteration: i=0, j=len(runes)-1
	// Step Over (F10) to see loop execute
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		// BREAKPOINT 22: Set here at start of loop body
		// DEBUG: In Variables panel, watch:
		//   - i: left pointer (starts at 0, increments)
		//   - j: right pointer (starts at len-1, decrements)
		// DEBUG: In Debug Console, type: runes[i]
		// DEBUG: In Debug Console, type: string(runes[i])
		// DEBUG: In Debug Console, type: runes[j]
		// DEBUG: In Debug Console, type: string(runes[j])
		// See which characters will be swapped

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
		//
		// BREAKPOINT 23: Set here BEFORE swap
		// DEBUG: In Variables panel, expand 'runes' slice
		// Note values at indices i and j
		// DEBUG: In Debug Console, type: runes[i]
		// DEBUG: In Debug Console, type: runes[j]
		// BREAKPOINT 24: Set here AFTER swap
		// DEBUG: In Variables panel, expand 'runes' slice again
		// Values at i and j should be swapped
		// DEBUG: In Debug Console, type: runes[i]
		// DEBUG: In Debug Console, type: runes[j]
		// Compare before/after - they've switched positions
		// DEBUG: Watch how the underlying array is modified in-place
		runes[i], runes[j] = runes[j], runes[i]

		// BREAKPOINT 25: Set here at end of loop iteration
		// DEBUG: In Variables panel, watch 'i' and 'j' update
		// Step Over (F10) - i increments, j decrements
		// DEBUG: In Debug Console, type: i
		// DEBUG: In Debug Console, type: j
		// When i >= j, loop terminates (pointers have met/crossed)
		// DEBUG: In Debug Console, type: runes
		// See progressive reversal with each iteration
	}

	// BREAKPOINT 26: Set here AFTER loop completes
	// DEBUG: In Variables panel, expand 'runes' slice
	// All elements should now be in reverse order
	// DEBUG: In Debug Console, type: runes
	// Compare with original input - completely reversed
	// DEBUG: In Debug Console, type: string(runes)
	// See the reversed string representation

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
	//
	// BREAKPOINT 27: Set here BEFORE string conversion
	// DEBUG: In Variables panel, expand 'runes' - see int32 array
	// BREAKPOINT 28: Set here at return statement
	// DEBUG: In Variables panel, see final string result
	// DEBUG: In Debug Console, type: string(runes)
	// This is the UTF-8 encoded result
	// DEBUG: Watch how []rune is encoded to UTF-8 bytes
	// For emoji, each int32 rune becomes multiple bytes
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
//
// DEBUGGING WORKFLOW:
// ===================
// 1. Set breakpoint at function entry (line 430)
// 2. Call with: RuneLen("👋😀")
// 3. Watch how it counts runes vs bytes
// 4. Step into utf8.RuneCountInString to see internals
func RuneLen(s string) int {
	// BREAKPOINT 29: Set here to inspect function entry
	// DEBUG: In Variables panel, expand 's' to see:
	//   - s.len: byte length
	//   - s.data: pointer to byte array
	// DEBUG: In Debug Console, type: len(s)
	// This is the BYTE count, not character count
	// For "👋😀", len(s) = 8 bytes
	// DEBUG: In Debug Console, type: utf8.RuneCountInString(s)
	// This is the CHARACTER count
	// For "👋😀", RuneCount = 2 characters

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
	//
	// DEBUG: In Variables panel, note that 's' is a local copy
	// Modifying 's' here wouldn't affect the caller's string
	// (though strings are immutable, so we can't modify anyway)

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
	// BREAKPOINT 30: Set here BEFORE utf8.RuneCountInString call
	// Step Into (F11) to see the stdlib implementation
	// Watch how it decodes UTF-8 sequences byte by byte
	// DEBUG: In Variables panel, 's' hasn't been modified
	// BREAKPOINT 31: Set here at return statement
	// DEBUG: In Variables panel, see the return value
	// DEBUG: In Debug Console, type: utf8.RuneCountInString(s)
	// DEBUG: In Debug Console, type: len(s)
	// Compare rune count vs byte count - key difference!
	//
	// Alternative implementations (for learning):
	//
	// Method 1: Convert to []rune (allocates!)
	//   return len([]rune(s))
	//   Problem: Allocates a rune slice on the heap (wasteful)
	//   DEBUG: Try in Debug Console: len([]rune(s))
	//   Same result, but check memory allocation in profiler
	//
	// Method 2: range loop (no allocation)
	//   count := 0
	//   for range s {
	//       count++
	//   }
	//   return count
	//   Note: range over string iterates RUNES, not bytes!
	//   DEBUG: Try this in Debug Console to verify same result
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
	//
	// DEBUGGING TIP: Use Watch Expressions
	// Add these to the Watch panel (right-click in Variables panel):
	// 1. len(s) - byte length
	// 2. utf8.RuneCountInString(s) - character count
	// 3. s - the actual string value
	// Watch these update as you step through the code
	return utf8.RuneCountInString(s)
}

/*
ADVANCED DEBUGGING TECHNIQUES:
===============================

1. CONDITIONAL BREAKPOINTS:
   Right-click on any breakpoint and add conditions
   Examples:
   - In TitleCase loop: i == 2 (break only on 3rd iteration)
   - In Reverse loop: i == j (break when pointers meet)
   - In RuneLen: len(s) > 10 (break only for long strings)

2. LOGPOINTS:
   Right-click on line number → Add Logpoint
   Examples:
   - In TitleCase loop: Log "Processing word {i}: {word}"
   - In Reverse loop: Log "Swapping i={i}, j={j}: {runes[i]} ↔ {runes[j]}"
   - No need to modify code or add print statements!

3. DEBUG CONSOLE EXPRESSIONS:
   During debugging, try these in the Debug Console:
   - Type: s to see current string value
   - Type: []byte(s) to see byte representation
   - Type: []rune(s) to see rune representation
   - Type: utf8.ValidString(s) to check if valid UTF-8
   - Type: fmt.Printf("%q", s) to see escaped representation
   - Type: fmt.Printf("%#v", s) to see Go syntax representation

4. WATCH EXPRESSIONS:
   Add these to the Watch panel for real-time monitoring:
   - len(s) vs utf8.RuneCountInString(s)
   - cap(words) vs len(words)
   - runes[i] for current character value
   - string(runes[i]) for character representation

5. CALL STACK NAVIGATION:
   In the Call Stack panel:
   - Click different frames to see state at each level
   - Observe how variables change across stack frames
   - Use "Step Out" (Shift+F11) to return to caller

6. MEMORY INSPECTION:
   To see memory allocations:
   - Use Go's pprof in production
   - In debugger, watch heap allocations via escape analysis
   - Compare &s (address) in different stack frames
   - Notice when []rune(s) allocates new memory

7. STEP COMMANDS:
   - F10 (Step Over): Execute line, don't enter functions
   - F11 (Step Into): Enter function calls to see internals
   - Shift+F11 (Step Out): Return to caller
   - F5 (Continue): Run until next breakpoint

8. DATA BREAKPOINTS:
   Watch for when specific variables change:
   - Right-click variable → Break When Value Changes
   - Useful for tracking when words[i] gets modified

Alternatives & Trade-offs:
==========================

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

DEBUGGING EXERCISES:
====================

Exercise 1: Trace TitleCase execution
- Input: "hello world"
- Set breakpoints at: 66, 74, 78, 120, 155, 168, 181
- Watch: s, words, i, word, runes
- Question: How many times does the loop iterate?
- Question: What's the value of words[0] before and after the loop?

Exercise 2: Understand UTF-8 vs ASCII
- Input: "café" vs "cafe"
- Set breakpoints at: 66, 120, 155
- In Debug Console: len(s) vs utf8.RuneCountInString(s)
- Question: Why are they different for "café" but same for "cafe"?

Exercise 3: Two-pointer algorithm visualization
- Input: "12345"
- Set breakpoints at: 235, 275, 302, 330
- Watch: runes, i, j
- Step through each iteration
- Question: What are i and j values at each iteration?
- Question: When does the loop terminate?

Exercise 4: Empty string edge case
- Input: ""
- Set breakpoints at: 66, 235, 430
- Question: Does TitleCase loop execute? How many times?
- Question: Does Reverse loop execute? Why not?
- Question: What does RuneLen return?

Exercise 5: Multi-byte emoji handling
- Input: "Hi👋"
- Set breakpoints at: 235, 275, 302, 330, 365
- Watch: s, runes, i, j
- In Debug Console: len(s) vs len(runes)
- Question: How many bytes is the emoji?
- Question: How is it stored as a rune (int32)?
- Question: After reversal, is the emoji still intact?
*/
