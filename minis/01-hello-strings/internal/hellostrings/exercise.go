//go:build !solution && !reference

package hellostrings

/*
Problem: Implement UTF-8-aware string utilities in Go
Constraints:
- Must handle multi-byte UTF-8 characters (emoji, accented letters, CJK)
- Preserve all characters without corruption
- Use only the Go standard library
Time/Space Complexity:
- TitleCase: O(n) time, O(n) space (allocates new string)
- Reverse: O(n) time, O(n) space (allocates rune slice + result string)
- RuneLen: O(n) time, O(1) space (just counting)
*/

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TitleCase - TODO: implement this function
func TitleCase(s string) string {
	// split input string into words
	words := strings.Fields(s)
	// iterate through words, converting to runes
	for i, word := range words {
		runes := []rune(word)
		// capitalize non-empty words
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		// convert runes back to a string
		words[i] = string(runes)
	}
	
	return strings.Join(words, " ")
}

// Reverse - TODO: implement this function
func Reverse(s string) string {
	runes := []rune(s) // convert string to runes
	// iterate through runes in reverse order
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		// swap runes
		runes[i], runes[j] = runes[j], runes[i]
	}
	// convert runes back to a string
	return string(runes)
}	

// RuneLen - TODO: implement this function
func RuneLen(s string) int {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return utf8.RuneCountInString(s)
}

