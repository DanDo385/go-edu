package main

import (
	"fmt"
	"strings"

	"minis/02-arrays-maps-basics/internal/arraysmapsbasics"
)

/*
Debug Harness for Arrays and Maps Basics

This file automatically demonstrates the project's capabilities by running through
different scenarios with pre-configured inputs. No CLI arguments needed!

How to use:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5 in VS Code
  3. Select "Debug cmd/dev (Debug Harness)"
  4. Step through with F10 (Step Over) and F11 (Step Into)
*/

func main() {
	fmt.Println("=== Arrays and Maps Basics - Auto Demo ===\n")

	// Demo 1: Simple word counting
	fmt.Println("Demo 1: Simple word frequency counting")
	fmt.Println("----------------------------------------")
	demo1 := "hello\nworld\nhello\ngo\n"
	fmt.Printf("Input:\n%s\n", demo1)
	
	// BREAKPOINT: Step into FreqFromReader to see how word frequencies are counted
	freq1, mostCommon1, err := arraysmapsbasics.FreqFromReader(strings.NewReader(demo1))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	fmt.Println("Results:")
	for word, count := range freq1 {
		fmt.Printf("  %s: %d\n", word, count)
	}
	fmt.Printf("Most common: %s\n\n", mostCommon1)

	// Demo 2: Case insensitive counting
	fmt.Println("Demo 2: Case insensitive word counting")
	fmt.Println("----------------------------------------")
	demo2 := "Go\ngo\nGO\nRust\nrust\nRUST\n"
	fmt.Printf("Input:\n%s\n", demo2)
	
	// BREAKPOINT: See how case normalization works
	freq2, mostCommon2, err := arraysmapsbasics.FreqFromReader(strings.NewReader(demo2))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	fmt.Println("Results:")
	for word, count := range freq2 {
		fmt.Printf("  %s: %d\n", word, count)
	}
	fmt.Printf("Most common: %s\n\n", mostCommon2)

	// Demo 3: Handling blank lines
	fmt.Println("Demo 3: Handling blank lines and whitespace")
	fmt.Println("----------------------------------------")
	demo3 := "hello\n\nworld\n\nhello\n  go  \n"
	fmt.Printf("Input:\n%s\n", demo3)
	
	// BREAKPOINT: See how blank lines are ignored
	freq3, mostCommon3, err := arraysmapsbasics.FreqFromReader(strings.NewReader(demo3))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	fmt.Println("Results:")
	for word, count := range freq3 {
		fmt.Printf("  %s: %d\n", word, count)
	}
	fmt.Printf("Most common: %s\n\n", mostCommon3)

	// Demo 4: Empty input
	fmt.Println("Demo 4: Empty input handling")
	fmt.Println("----------------------------------------")
	demo4 := ""
	fmt.Printf("Input: (empty)\n\n")
	
	// BREAKPOINT: See how empty input is handled
	freq4, mostCommon4, err := arraysmapsbasics.FreqFromReader(strings.NewReader(demo4))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	fmt.Println("Results:")
	if len(freq4) == 0 {
		fmt.Println("  (no words found)")
	} else {
		for word, count := range freq4 {
			fmt.Printf("  %s: %d\n", word, count)
		}
	}
	fmt.Printf("Most common: %s\n\n", mostCommon4)

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Maps are perfect for counting frequencies")
	fmt.Println("2. Case normalization ensures consistent counting")
	fmt.Println("3. Blank lines should be ignored when processing text")
	fmt.Println("4. Empty input should return empty results, not errors")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/03-csv-stats")
}
