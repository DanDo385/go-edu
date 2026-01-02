package main

import (
	"fmt"

	"minis/01-hello-strings/internal/hellostrings"
)

/*
String Operations Debug Harness

This file provides a deterministic debug environment with fixed inputs.
Use this for stepping through the code with VS Code debugger.

How to use:
  1. Set breakpoints in internal/hellostrings/exercise.go
  2. Press F5 in VS Code
  3. Select "Debug Current Package" or "Debug cmd/dev/main.go"
  4. Step through with F10 (Step Over) and F11 (Step Into)
*/
func main() {
	fmt.Println("=== String Operations Debug Harness ===")
	fmt.Println()

	// ===== Test 1: TitleCase =====
	// BREAKPOINT: Set a breakpoint here to inspect the input
	fmt.Println("--- Test 1: TitleCase ---")
	input1 := "hello world"
	fmt.Printf("Input:  %q\n", input1)

	// BREAKPOINT: Step into TitleCase to see the implementation
	result1 := hellostrings.TitleCase(input1)
	fmt.Printf("Output: %q\n", result1)
	fmt.Printf("Expected: \"Hello World\"\n")
	fmt.Println()

	// ===== Test 2: TitleCase with emoji =====
	fmt.Println("--- Test 2: TitleCase with emoji ---")
	input2 := "hello 👋 world"
	fmt.Printf("Input:  %q\n", input2)

	result2 := hellostrings.TitleCase(input2)
	fmt.Printf("Output: %q\n", result2)
	fmt.Printf("Expected: \"Hello 👋 World\"\n")
	fmt.Println()

	// ===== Test 3: Reverse =====
	// BREAKPOINT: Watch how multi-byte characters are handled
	fmt.Println("--- Test 3: Reverse ---")
	input3 := "Hello World"
	fmt.Printf("Input:    %q\n", input3)

	// BREAKPOINT: Step into Reverse
	result3 := hellostrings.Reverse(input3)
	fmt.Printf("Reversed: %q\n", result3)
	fmt.Printf("Expected: \"dlroW olleH\"\n")
	fmt.Println()

	// ===== Test 4: Reverse with emoji =====
	fmt.Println("--- Test 4: Reverse with emoji ---")
	input4 := "Hello 👋 World"
	fmt.Printf("Input:    %q\n", input4)

	result4 := hellostrings.Reverse(input4)
	fmt.Printf("Reversed: %q\n", result4)
	fmt.Printf("Expected: \"dlroW 👋 olleH\"\n")
	fmt.Printf("Note: The emoji is preserved as a single unit!\n")
	fmt.Println()

	// ===== Test 5: RuneLen =====
	// BREAKPOINT: Compare rune count vs byte count
	fmt.Println("--- Test 5: RuneLen ---")
	input5 := "Hello"
	fmt.Printf("String: %q\n", input5)

	runeCount5 := hellostrings.RuneLen(input5)
	byteCount5 := len(input5)
	fmt.Printf("Rune count: %d\n", runeCount5)
	fmt.Printf("Byte count: %d\n", byteCount5)
	fmt.Printf("Expected: 5 (both should be equal for ASCII)\n")
	fmt.Println()

	// ===== Test 6: RuneLen with multi-byte characters =====
	fmt.Println("--- Test 6: RuneLen with multi-byte characters ---")
	input6 := "Hello 世界 🌍"
	fmt.Printf("String: %q\n", input6)

	// BREAKPOINT: See the difference between rune and byte count
	runeCount6 := hellostrings.RuneLen(input6)
	byteCount6 := len(input6)
	fmt.Printf("Rune count: %d (Unicode code points)\n", runeCount6)
	fmt.Printf("Byte count: %d (UTF-8 encoded bytes)\n", byteCount6)
	fmt.Printf("Expected: 10 runes, 17 bytes\n")
	fmt.Println()

	// ===== Summary =====
	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Go strings are byte slices, not character arrays")
	fmt.Println("2. Multi-byte UTF-8 sequences (emoji, CJK) require rune handling")
	fmt.Println("3. Use []rune() to convert strings for character-level operations")
	fmt.Println("4. len(s) returns bytes, not characters")
	fmt.Println("5. Always use UTF-8 aware functions for international text")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/02-arrays-maps-basics")
}
