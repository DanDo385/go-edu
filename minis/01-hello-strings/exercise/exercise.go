//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "strings" for string manipulation functions
// - "unicode" for character case conversion
// - "unicode/utf8" for UTF-8 utilities
//
// Uncomment these imports when you implement the functions:
/*
import (
	"strings"
	"unicode"
	"unicode/utf8"
)
*/

// TitleCase converts the first letter of each word to uppercase.
// Words are separated by whitespace.
// Example: "hello world" → "Hello World"
func TitleCase(s string) string {
	// TODO: Implement the function.

	// Step 1: Split the string into words.
	// - Use `strings.Fields(s)` to split the input string `s` into a slice of words.
	// - `strings.Fields` is smart and handles multiple spaces and other whitespace characters.
	// - Memory: `strings.Fields` allocates a new slice `[]string` to store the words. The underlying string data, however, is shared until a word is modified.

	// Step 2: Iterate over the slice of words.
	// - Use a `for i, word := range words` loop.
	// - In Go, `range` on a slice provides the index `i` and a *copy* of the element `word`.
	// - To modify the slice, you must use the index, like `words[i] = ...`.

	// Step 3: For each word, capitalize the first letter.
	// - A `string` in Go is a read-only slice of bytes, and it's UTF-8 encoded. To modify characters, you must convert the string to a slice of runes.
	// - `runes := []rune(word)`
	// - Memory: This conversion allocates a new slice `[]rune` on the heap. Each rune is an `int32`, so this can use more memory than the original string.
	// - Check if the word is not empty (`len(runes) > 0`).
	// - `runes[0] = unicode.ToUpper(runes[0])`. `unicode.ToUpper` correctly handles all Unicode characters, not just ASCII.

	// Step 4: Convert the modified rune slice back to a string.
	// - `words[i] = string(runes)`
	// - Memory: This conversion allocates *another* new string. The Go runtime has to encode the UTF-8 characters from the rune slice into bytes.

	// Step 5: Join the words back into a single string.
	// - Use `strings.Join(words, " ")`.
	// - Memory: This is the most efficient way to join strings. It calculates the final size, allocates one new string, and copies all the words into it. This avoids creating intermediate strings that you would get with repeated `+` concatenation.
	// - Return the result.

	// --- Memory & Language Comparison ---
	//
	// - Go: Strings are immutable. Modifications require allocations. The garbage collector handles cleanup. This is a balance between safety and performance.
	// - C: You would manually manage memory with `malloc` and `free`. You'd likely modify the string in-place (if it's not a string literal). This is fast but error-prone (buffer overflows, memory leaks).
	// - Rust: The "borrow checker" would enforce memory safety at compile time. You'd use methods on `&str` or `String` to manipulate the string, and the compiler would ensure you don't have memory bugs. This provides C-like performance with Go-like safety.
	// - Python: Strings are also immutable. The process is similar to Go (e.g., `s.title()`), but the memory details are more hidden. Python's garbage collector handles all memory management.

	return ""
}

// Reverse returns the string reversed character-by-character (UTF-8 aware).
// This correctly handles multi-byte characters like emoji.
// Example: "Hello 👋" → "👋 olleH"
func Reverse(s string) string {
	// TODO: Implement the function.

	// Step 1: Convert the string to a slice of runes.
	// - `runes := []rune(s)`
	// - Why runes? A string is a sequence of bytes. Reversing the bytes of a multi-byte character (like an emoji) would corrupt it. A slice of runes represents the string as a sequence of Unicode code points, which can be safely reordered.
	// - Memory: This conversion allocates a new `[]rune` slice on the heap.

	// Step 2: Reverse the slice of runes in-place.
	// - Use the "two-pointer" technique.
	// - `for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 { ... }`
	//   - `i` starts at the beginning of the slice and moves forward.
	//   - `j` starts at the end of the slice and moves backward.
	//   - The loop continues until `i` and `j` meet or cross.
	// - Inside the loop, swap the elements: `runes[i], runes[j] = runes[j], runes[i]`
	// - This is a feature of Go called "tuple assignment" or "parallel assignment". It's a clean way to swap two values without a temporary variable.
	// - Memory: This step is "in-place". It modifies the `[]rune` slice's underlying memory directly. No new memory is allocated for the slice itself.

	// Step 3: Convert the reversed slice of runes back to a string.
	// - `return string(runes)`
	// - Memory: This conversion allocates a new string on the heap to store the UTF-8 encoded bytes of the runes. The `[]rune` slice can now be garbage collected.

	// --- Memory & Language Comparison ---
	//
	// - Go: The approach is to convert to a mutable representation (`[]rune`), modify it, then convert back to an immutable `string`. This involves a couple of allocations but is safe and clear.
	// - C: You would likely have a `char*` that points to a mutable buffer. You could reverse the string in-place. If dealing with UTF-8, you'd need a library to handle multi-byte characters correctly, and you would be responsible for all memory management.
	// - Rust: You would collect the `chars()` into a `Vec<char>`, reverse the `Vec`, and then `collect()` it back into a `String`. The memory model is similar to Go's, but the borrow checker provides compile-time guarantees against memory errors.
	// - Python: You can reverse a string with a simple slice: `s[::-1]`. This is very concise, but it also creates a new string object in memory, similar to the Go approach.

	return ""
}


// RuneLen returns the number of UTF-8 runes (characters) in the string,
// not the byte count. This is important for strings with non-ASCII characters.
// Example: "café" has 4 runes but 5 bytes (é is 2 bytes in UTF-8)
func RuneLen(s string) int {
	// TODO: Implement the function.

	// In Go, there's a critical distinction between the number of bytes and the number of characters (runes) in a string.
	// `len(s)` gives you the number of bytes, which is often not what you want for international text.

	// **Option 1 (The Best Way): Use the `unicode/utf8` package.**
	// - `return utf8.RuneCountInString(s)`
	// - Memory: This is the most efficient method. It performs its work without allocating any new memory on the heap. It simply iterates over the bytes of the string and counts the runes.
	// - Performance: It's highly optimized and will be the fastest option.

	// **Option 2 (The Idiomatic Loop Way): Use a `for...range` loop.**
	// - A `for...range` loop over a string in Go iterates over runes, not bytes.
	//   ```
	//   count := 0
	//   for range s {
	//       count++
	//   }
	//   return count
	//   ```
	// - Memory: This approach is also very memory-efficient. Like `RuneCountInString`, it doesn't require any heap allocations.

	// **Option 3 (The Naive Way): Convert to a rune slice.**
	// - `return len([]rune(s))`
	// - Memory: This is the least memory-efficient approach. It allocates a new slice on the heap (`[]rune`) to hold all the runes, and then takes the length of that slice. For a large string, this can be a significant amount of memory.

	// --- Memory & Language Comparison ---
	//
	// - Go: Go is explicit. `len()` gives you bytes. If you want characters, you must ask for them, and you have several ways to do it, with different performance trade-offs. This forces you to be aware of how strings work.
	// - C: A "string" is just a sequence of bytes ending in a null terminator. There is no standard concept of a character count. You would need a library like ICU to correctly count characters in a UTF-8 string, and you'd be managing all the memory yourself.
	// - Rust: Similar to Go, Rust distinguishes between bytes and characters. `s.len()` is the byte length, and `s.chars().count()` is the character count. Rust's philosophy is also about being explicit and safe.
	// - Python: In Python 3, `len(s)` gives you the number of characters, which is usually what you want. The complexity of the underlying encoding is hidden from you. This is more convenient but less explicit about what's happening with memory.

	return 0
}

// After implementing all functions:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Build: go build (should compile without errors)
// - Compare with solution.go to see detailed explanations
