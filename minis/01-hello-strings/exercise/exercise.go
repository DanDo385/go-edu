//go:build !solution
// +build !solution

package exercise

import (
	"strings"
	"unicode"
)

// TitleCase converts the first letter of each word to uppercase.
// Words are separated by whitespace.
// Example: "hello world" → "Hello World"
func TitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// Reverse returns the string reversed character-by-character (UTF-8 aware).
// This correctly handles multi-byte characters like emoji.
// Example: "Hello 👋" → "👋 olleH"
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// RuneLen returns the number of UTF-8 runes (characters) in the string,
// not the byte count. This is important for strings with non-ASCII characters.
// Example: "café" has 4 runes but 5 bytes (é is 2 bytes in UTF-8)
func RuneLen(s string) int {
	return len([]rune(s))
}