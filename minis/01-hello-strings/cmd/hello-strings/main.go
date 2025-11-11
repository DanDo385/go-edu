package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/01-hello-strings/exercise"
)

func main() {
	// Demonstrate string utilities with various inputs
	testCases := []string{
		"hello world",
		"the quick brown fox",
		"Hello 👋 World",
		"café résumé",
		"日本語",
	}

	fmt.Println("=== TitleCase Demo ===")
	for _, s := range testCases {
		fmt.Printf("%-25s → %s\n", s, exercise.TitleCase(s))
	}

	fmt.Println("\n=== Reverse Demo ===")
	for _, s := range testCases {
		fmt.Printf("%-25s → %s\n", s, exercise.Reverse(s))
	}

	fmt.Println("\n=== RuneLen Demo ===")
	for _, s := range testCases {
		byteLen := len(s)
		runeLen := exercise.RuneLen(s)
		fmt.Printf("%-25s → bytes: %2d, runes: %2d\n", s, byteLen, runeLen)
	}
}
