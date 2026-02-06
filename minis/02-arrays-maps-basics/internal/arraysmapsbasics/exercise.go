//go:build !solution && !reference

package arraysmapsbasics

/*
Problem: Count word frequencies from text input and find the most common word
Constraints:
- Normalize to lowercase ("Hello" == "hello")
- Ignore blank lines
- For ties, return any of the tied words (arbitrary but deterministic)
Time/Space Complexity:
- Time: O(n) where n = number of words (one pass to build map, one to find max)
- Space: O(u) where u = number of unique words (map storage)
*/

import (
	"bufio"
	"io"
	"strings"
)

// FreqFromReader - TODO: implement this function
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	scanner := bufio.NewScanner(r)
	freq := make(map[string]int)
	for scanner.Scan() {
		word := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if word == "" {
			continue
		}
		freq[word]++
	}
	if err := scanner.Err(); err != nil {
		return nil, "", err
	}

	mostCommon := ""
	maxCount := 0
	for word, count := range freq {
		if count > maxCount || (count == maxCount && (mostCommon == "" || word < mostCommon)) {
			mostCommon = word
			maxCount = count
		}
	}

	return freq, mostCommon, nil
}
