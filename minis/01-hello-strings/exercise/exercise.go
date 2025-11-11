//go:build !solution
// +build !solution

package exercise

// TitleCase converts the first letter of each word to uppercase.
// Words are separated by whitespace.
// Example: "hello world" → "Hello World"
func TitleCase(s string) string {
	// TODO: implement
	return ""
}

// Reverse returns the string reversed character-by-character (UTF-8 aware).
// This correctly handles multi-byte characters like emoji.
// Example: "Hello 👋" → "👋 olleH"
func Reverse(s string) string {
	// TODO: implement
	return ""
}

// RuneLen returns the number of UTF-8 runes (characters) in the string,
// not the byte count. This is important for strings with non-ASCII characters.
// Example: "café" has 4 runes but 5 bytes (é is 2 bytes in UTF-8)
func RuneLen(s string) int {
	// TODO: implement
	return 0
}
