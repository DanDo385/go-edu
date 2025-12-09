//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "strings" for string manipulation functions
// - "unicode" for character case conversion
//
// import (
//     "strings"
//     "unicode"
// )

// TitleCase converts the first letter of each word to uppercase.
// Words are separated by whitespace.
// Example: "hello world" → "Hello World"
//
// TODO: Implement TitleCase function
// Function signature: func TitleCase(s string) string
//
// Steps to implement:
// 1. Split the string into words using strings.Fields(s)
//    - Why Fields? It automatically handles multiple spaces and all whitespace types
//    - Returns a []string (slice of strings) - this is a reference type in Go
//
// 2. Loop through each word with range
//    - range returns (index, value) - both are COPIES of the slice elements
//    - Modifying 'value' won't change the original slice, you must use index
//
// 3. Convert each word to []rune (slice of Unicode code points)
//    - Why runes? Strings in Go are UTF-8 byte sequences
//    - A rune is an int32 representing a Unicode code point
//    - Converting to []rune lets you manipulate individual characters safely
//    - This creates a NEW slice (allocation happens here)
//
// 4. Capitalize the first rune using unicode.ToUpper()
//    - unicode.ToUpper() takes a rune and returns a rune (passed by value)
//    - Modify runes[0] directly to change the first character
//
// 5. Convert back to string and update the words slice
//    - string(runes) creates a NEW string from the rune slice
//    - Strings in Go are immutable - you can't modify them in place
//    - Use words[i] = string(runes) to update the slice element
//
// 6. Join words with spaces using strings.Join(words, " ")
//    - This creates a NEW string by concatenating all words
//    - Returns the final result
//
// Key Go concepts:
// - Strings are immutable (can't change bytes directly)
// - Slices are references to underlying arrays (but slice elements can be updated)
// - Converting string → []rune allocates memory (creates a copy)
// - range makes copies of values (use index to modify slice elements)

// TODO: Implement the TitleCase function below
// func TitleCase(s string) string {
//     return ""
// }

// Reverse returns the string reversed character-by-character (UTF-8 aware).
// This correctly handles multi-byte characters like emoji.
// Example: "Hello 👋" → "👋 olleH"
//
// TODO: Implement Reverse function
// Function signature: func Reverse(s string) string
//
// Steps to implement:
// 1. Convert string to []rune
//    - Why? Reversing bytes directly would corrupt multi-byte UTF-8 characters
//    - Example: emoji "👋" is 4 bytes [0xF0, 0x9F, 0x91, 0x8B]
//    - Reversing those bytes = invalid UTF-8! But as a rune it's one value
//
// 2. Use two-pointer technique to reverse the slice in place
//    - Initialize: i = 0 (start), j = len(runes) - 1 (end)
//    - Loop while i < j
//    - Swap: runes[i], runes[j] = runes[j], runes[i]
//      * This is simultaneous assignment (Go feature)
//      * Both sides are evaluated before assignment
//      * No temporary variable needed!
//    - Move pointers: i++ (increment i by 1), j-- (decrement j by 1)
//
// 3. Convert []rune back to string
//    - string(runes) allocates a new string with UTF-8 encoding
//    - The rune slice can be garbage collected after this
//
// Key Go concepts:
// - []rune is a slice (reference type pointing to underlying array)
// - Swapping modifies the slice in place (no copy of whole slice)
// - Simultaneous assignment (a, b = b, a) is safe and idiomatic
// - Memory: one allocation for []rune, one for final string

// TODO: Implement the Reverse function below
// func Reverse(s string) string {
//     return ""
// }

// RuneLen returns the number of UTF-8 runes (characters) in the string,
// not the byte count. This is important for strings with non-ASCII characters.
// Example: "café" has 4 runes but 5 bytes (é is 2 bytes in UTF-8)
//
// TODO: Implement RuneLen function
// Function signature: func RuneLen(s string) int
//
// Steps to implement:
// Option 1 (Simple but allocates memory):
// 1. Convert string to []rune
// 2. Return len([]rune(s))
//    - This works but allocates a slice on the heap
//    - Fine for small strings, wasteful for large ones
//
// Option 2 (Efficient, no allocation):
// 1. Use range loop: for range s
//    - range over a string iterates RUNES, not bytes!
//    - Each iteration decodes one UTF-8 character
//    - The value returned is the rune, the index is the byte position
// 2. Count iterations
// 3. Return count
//
// Key Go concepts:
// - len(s) returns BYTE length, not character count
// - range over string yields runes (automatically decodes UTF-8)
// - []rune(s) allocates memory, range does not
// - For production: use utf8.RuneCountInString(s) from "unicode/utf8"
//   (it's optimized and doesn't allocate)

// TODO: Implement the RuneLen function below
// func RuneLen(s string) int {
//     return 0
// }

// After implementing all functions:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Build: go build (should compile without errors)
// - Compare with solution.go to see detailed explanations
