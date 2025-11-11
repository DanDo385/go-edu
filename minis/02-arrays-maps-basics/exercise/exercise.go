package exercise

import "io"

// FreqFromReader reads words from r (one per line) and returns:
// - A frequency map (word -> count)
// - The most common word (arbitrary choice if there's a tie)
// - An error if reading fails
//
// Words are normalized to lowercase. Blank lines are ignored.
//
// Example:
//   Input:  "hello\nworld\nhello\n"
//   Output: map[string]int{"hello": 2, "world": 1}, "hello", nil
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	// TODO: implement
	return nil, "", nil
}
