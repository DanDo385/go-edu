//go:build !solution
// +build !solution

package exercise

import (
	"bufio"
	"io"
	"strings"
)

// FreqFromReader reads words from r (one per line) and returns:
// - A frequency map (word -> count)
// - The most common word (arbitrary choice if there's a tie)
// - An error if reading fails

func FreqFromReader(r io.Reader) (map[string]int, string, error) {
	freq := make(map[string]int)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		word := strings.ToLower(strings.TrimSpace(line))
		if word == "" {
			continue
		}
		freq[word]++
	}

	if err := scanner.Err(); err != nil {
		return freq, "", err
	}

	var maxWord string
	var maxCount int

	for word, count := range freq {
		if count > maxCount {
			maxWord = word
			maxCount = count
		}
	}
	return freq, maxWord, nil
}
