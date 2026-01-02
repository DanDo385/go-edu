package main

import (
	"fmt"
	"strings"

	"github.com/example/go-10x-minis/minis/02-arrays-maps-basics/internal/arraysmapsbasics"
)

/*
Debug Harness for Arrays Maps Basics Module

This file runs the word frequency counter with sample inputs.
Set breakpoints in exercise.go and press F5 in VS Code.
*/
func main() {
	fmt.Println("=== Arrays Maps Basics Debug Harness ===")
	fmt.Println()

	// Demo 1: Simple text
	fmt.Println("--- Demo 1: Simple Text ---")
	text1 := `the quick brown fox jumps over the lazy dog
the fox was quick and the dog was lazy`

	fmt.Printf("Input:\n%s\n\n", text1)
	reader1 := strings.NewReader(text1)
	freq1, mostCommon1, err := arraysmapsbasics.FreqFromReader(reader1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Most common word: %q (count: %d)\n", mostCommon1, freq1[mostCommon1])
	fmt.Printf("Word frequencies: %v\n\n", freq1)

	// Demo 2: Mixed case
	fmt.Println("--- Demo 2: Case Normalization ---")
	text2 := "Hello hello HELLO World world"
	fmt.Printf("Input: %q\n", text2)
	reader2 := strings.NewReader(text2)
	freq2, mostCommon2, _ := arraysmapsbasics.FreqFromReader(reader2)
	fmt.Printf("Most common: %q (count: %d)\n", mostCommon2, freq2[mostCommon2])
	fmt.Printf("Frequencies: %v\n\n", freq2)

	// Demo 3: Single word
	fmt.Println("--- Demo 3: Single Word ---")
	text3 := "unique"
	fmt.Printf("Input: %q\n", text3)
	reader3 := strings.NewReader(text3)
	freq3, mostCommon3, _ := arraysmapsbasics.FreqFromReader(reader3)
	fmt.Printf("Most common: %q (count: %d)\n", mostCommon3, freq3[mostCommon3])
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Maps are Go's hash table implementation")
	fmt.Println("2. strings.Fields splits on whitespace")
	fmt.Println("3. strings.ToLower normalizes case")
	fmt.Println("4. bufio.Scanner efficiently reads line by line")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/03-csv-stats")
}
