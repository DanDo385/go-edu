package main

import (
	"fmt"

	"minis/01-hello-strings/internal/hellostrings"
)

/*
Debug Harness for Hello Strings

This file automatically demonstrates the project's capabilities by running through
different scenarios with pre-configured inputs. No CLI arguments needed!

How to use:
  1. Set breakpoints at "// BREAKPOINT:" comments
  2. Press F5 in VS Code
  3. Select "Debug cmd/dev (Debug Harness)"
  4. Step through with F10 (Step Over) and F11 (Step Into)
*/

func main() {
	fmt.Println("=== Hello Strings - Auto Demo ===\n")

	// Demo 1: TitleCase
	fmt.Println("Demo 1: Title Case Conversion")
	fmt.Println("----------------------------------------")
	testCases := []string{
		"hello world",
		"HELLO WORLD",
		"hello    world",
		"café résumé",
		"hello 👋 world",
	}

	for _, input := range testCases {
		// BREAKPOINT: Step into TitleCase to see how UTF-8 is handled
		result := hellostrings.TitleCase(input)
		fmt.Printf("Input:  %-20q\n", input)
		fmt.Printf("Output: %-20q\n\n", result)
	}

	// Demo 2: Reverse
	fmt.Println("\nDemo 2: String Reversal (UTF-8 Aware)")
	fmt.Println("----------------------------------------")
	reverseCases := []string{
		"hello",
		"Hello World",
		"café",
		"Hello 👋 World",
		"🚀 Go is awesome!",
	}

	for _, input := range reverseCases {
		// BREAKPOINT: Step into Reverse to see rune-based reversal
		result := hellostrings.Reverse(input)
		fmt.Printf("Input:  %-20q\n", input)
		fmt.Printf("Output: %-20q\n\n", result)
	}

	// Demo 3: RuneLen
	fmt.Println("\nDemo 3: Rune Length vs Byte Length")
	fmt.Println("----------------------------------------")
	lenCases := []string{
		"hello",
		"café",
		"👋",
		"Hello 👋 World",
		"🚀 Go is awesome!",
	}

	for _, input := range lenCases {
		// BREAKPOINT: Step into RuneLen to see rune counting
		runeLen := hellostrings.RuneLen(input)
		byteLen := len(input)
		fmt.Printf("Input:    %-20q\n", input)
		fmt.Printf("Runes:    %d\n", runeLen)
		fmt.Printf("Bytes:    %d\n", byteLen)
		fmt.Printf("Diff:     %d bytes (multi-byte characters)\n\n", byteLen-runeLen)
	}

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Go strings are UTF-8 encoded byte sequences")
	fmt.Println("2. len(string) counts bytes, not characters")
	fmt.Println("3. Use utf8 package or range to iterate over runes")
	fmt.Println("4. Title case must handle UTF-8 correctly (accented chars, emoji)")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/02-arrays-maps-basics")
}
