package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/01-hello-strings/internal/hellostrings"
)

/*
Debug Harness for Hello Strings Module

This file runs all functions with sample inputs for debugging.
Set breakpoints in exercise.go and press F5 in VS Code.
*/
func main() {
	fmt.Println("=== Hello Strings Debug Harness ===")
	fmt.Println()

	// Demo 1: TitleCase
	fmt.Println("--- TitleCase Demo ---")
	inputs := []string{
		"hello world",
		"the quick BROWN fox",
		"café résumé naïve",
		"hello 世界",
	}
	for _, input := range inputs {
		result := hellostrings.TitleCase(input)
		fmt.Printf("  TitleCase(%q) = %q\n", input, result)
	}
	fmt.Println()

	// Demo 2: Reverse
	fmt.Println("--- Reverse Demo ---")
	reverseInputs := []string{
		"hello",
		"Hello World",
		"Hi👋",
		"café",
		"世界你好",
	}
	for _, input := range reverseInputs {
		result := hellostrings.Reverse(input)
		fmt.Printf("  Reverse(%q) = %q\n", input, result)
	}
	fmt.Println()

	// Demo 3: RuneLen
	fmt.Println("--- RuneLen Demo ---")
	lenInputs := []string{
		"hello",
		"café",
		"👋😀🎉",
		"世界",
		"",
	}
	for _, input := range lenInputs {
		result := hellostrings.RuneLen(input)
		fmt.Printf("  RuneLen(%q) = %d (bytes: %d)\n", input, result, len(input))
	}
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Runes (int32) represent Unicode code points")
	fmt.Println("2. len(s) returns bytes, not characters")
	fmt.Println("3. utf8.RuneCountInString counts characters without allocation")
	fmt.Println("4. []rune(s) converts to rune slice for manipulation")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/02-arrays-maps-basics")
}
