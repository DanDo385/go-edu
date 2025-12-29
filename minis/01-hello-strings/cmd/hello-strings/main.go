package main

import (
	"flag"   // Command-line argument parsing
	"fmt"    // Formatted I/O
	"os"     // Operating system interface

	// Import the exercise package from parent directory
	// This allows main.go to call functions from exercise.go and solution.go
	"github.com/example/go-10x-minis/minis/01-hello-strings/exercise"
)

/*
===================================================================
String Manipulation Demonstration
===================================================================

This program demonstrates fundamental string operations in Go by
calling functions from the exercise package.

USAGE:
    go run ./cmd/hello-strings/main.go
    go run ./cmd/hello-strings/main.go -input="hello world"
    go run ./cmd/hello-strings/main.go -input="café" -demo=reverse
    go run ./cmd/hello-strings/main.go -input="test" -solution=true

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see string transformations
- See /RUN_DEBUG.md for comprehensive debugging guide

PACKAGE STRUCTURE:
- This file (cmd/hello-strings/main.go) is package main
- It imports the exercise package from the parent directory
- exercise.go, solution.go are in package exercise
- Functions must be exported (capitalized) to be called from main
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect default values before parsing
	// DEBUG: Before flag.Parse(), all flags have their default values
	// DEBUG: In Variables panel, you'll see: input="hello world", demo="all", useSolution=false

	// flag.String creates a string flag variable
	// Syntax: flag.String(name, default, usage)
	// Returns a pointer to the string value
	input := flag.String("input", "hello world", "Input string to process")

	// Which demo to run: all, titlecase, reverse, runelen
	demo := flag.String("demo", "all", "Which demo to run: all, titlecase, reverse, runelen")

	// Use solution.go instead of exercise.go
	// Note: solution.go requires build tag -tags=solution
	useSolution := flag.Bool("solution", false, "Use solution.go instead of exercise.go (requires -tags=solution)")

	// flag.Parse() reads command-line arguments and updates flag variables
	flag.Parse()

	// BREAKPOINT: Set breakpoint here to inspect parsed values in Variables panel
	// DEBUG: After flag.Parse(), variables now hold actual command-line values
	// DEBUG: Watch 'input', 'demo', 'useSolution' - they may have changed from defaults

	// ============================================================================
	// STEP 2: Example 1 - TitleCase (Simple)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through simple example
	// DEBUG: This is the simplest use case - capitalizing first letter of each word
	// DEBUG: Step Over (F10) to execute line by line
	// DEBUG: Watch Variables panel to see how 'result' changes

	if *demo == "all" || *demo == "titlecase" {
		fmt.Println("=== TitleCase Demo ===")
		fmt.Println("Capitalizes the first letter of each word")
		fmt.Println()

		testCases := []string{
			"hello world",
			"the quick brown fox",
			"Hello 👋 World",
			"café résumé",
			"日本語",
			*input,  // User-provided text
		}

		// BREAKPOINT: Set breakpoint here before loop
		// DEBUG: Expand 'testCases' in Variables panel to see all test strings
		// DEBUG: Notice the len and cap of the slice

		for _, s := range testCases {
			// BREAKPOINT: Set breakpoint here inside loop
			// DEBUG: Watch 's' variable change with each iteration
			// DEBUG: Step Into (F11) on exercise.TitleCase() to see implementation

			// Call the TitleCase function from exercise package
			// BREAKPOINT: Set breakpoint here before function call
			// DEBUG: Step Into (F11) to enter exercise.TitleCase function
			// DEBUG: You'll see the implementation in exercise.go or solution.go
			result := exercise.TitleCase(s)

			// BREAKPOINT: Set breakpoint here after function call
			// DEBUG: Inspect 'result' variable - see how string was transformed
			// DEBUG: Compare 's' (input) with 'result' (output)

			fmt.Printf("%-25s → %s\n", s, result)
		}
		fmt.Println()
	}

	// ============================================================================
	// STEP 3: Example 2 - Reverse (Intermediate)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through intermediate example
	// DEBUG: This shows more advanced features - working with runes
	// DEBUG: Step Into (F11) on exercise.Reverse() to see rune manipulation

	if *demo == "all" || *demo == "reverse" {
		fmt.Println("=== Reverse Demo ===")
		fmt.Println("Reverses strings (handles Unicode correctly)")
		fmt.Println()

		testCases := []string{
			"hello",
			"Hello 👋 World",
			"café",
			"日本語",
			*input,
		}

		for _, s := range testCases {
			// BREAKPOINT: Set breakpoint here
			// DEBUG: Step Into (F11) to see reverse algorithm
			result := exercise.Reverse(s)

			// BREAKPOINT: Set breakpoint here after function call
			// DEBUG: Compare 's' and 'result' - string should be reversed
			fmt.Printf("%-25s → %s\n", s, result)
		}
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 3 - RuneLen (Advanced)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through advanced example
	// DEBUG: This demonstrates byte vs rune counting
	// DEBUG: Use Call Stack panel to see function call hierarchy

	if *demo == "all" || *demo == "runelen" {
		fmt.Println("=== RuneLen Demo ===")
		fmt.Println("Shows difference between byte length and rune (character) length")
		fmt.Println()

		testCases := []string{
			"hello",           // ASCII only
			"café",            // Contains é (2 bytes)
			"Hello 👋 World",  // Contains emoji (4 bytes)
			"日本語",           // Japanese (3 characters, 9 bytes)
			*input,
		}

		for _, s := range testCases {
			// BREAKPOINT: Set breakpoint here
			// DEBUG: Step Into (F11) to see rune counting implementation
			byteLen := len(s)               // len() returns byte length
			runeLen := exercise.RuneLen(s)  // Count actual characters

			// BREAKPOINT: Set breakpoint here
			// DEBUG: Compare byteLen vs runeLen
			// DEBUG: For ASCII, they're equal. For Unicode, runeLen < byteLen

			fmt.Printf("%-25s → bytes: %2d, runes: %2d\n", s, byteLen, runeLen)
		}
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Usage Note about Solution
	// ============================================================================
	// If -solution flag was used, remind user about build tags
	if *useSolution {
		fmt.Println("NOTE: To use solution.go, build/run with -tags=solution:")
		fmt.Println("  go run -tags=solution ./cmd/hello-strings/main.go -solution=true")
		fmt.Println("  go test -tags=solution")
		os.Exit(0)
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise package functions
	// 5. Watch Variables panel to see how data changes
	// 6. Add watch expressions: len(testCases), *input, result
	// 7. Use Call Stack panel to see function hierarchy
	//
	// STEPPING INTO EXERCISE PACKAGE:
	// - Set breakpoint at exercise.TitleCase(s) call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.go
	// - You can now see the implementation
	// - Step through the implementation with F10
	// - Press Shift+F11 (Step Out) to return to main.go
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}