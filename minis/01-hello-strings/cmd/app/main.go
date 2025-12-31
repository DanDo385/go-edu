package main

import (
	"fmt"    // Formatted I/O
	"os"     // Operating system interface

	// Import the exercise package from parent directory
	// This allows main.go to call functions from hellostrings.go and solution.go
	"github.com/example/go-10x-minis/minis/01-hello-strings/internal/hellostrings"
)

/*
===================================================================
String Manipulation Demonstration
===================================================================

This program demonstrates fundamental string operations in Go by
calling functions from the exercise package.

USAGE:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go "hello world"
    go run ./cmd/app/main.go "café" reverse
    go run ./cmd/app/main.go "test" titlecase

Arguments:
    [input]     - Input string to process (default: "hello world")
    [function]  - Function to run: all, titlecase, reverse, runelen (default: "all")

Examples:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go "hello world"
    go run ./cmd/app/main.go "café" reverse
    go run ./cmd/app/main.go "Hello 👋 World" runelen

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see string transformations
- See /RUN_DEBUG.md for comprehensive debugging guide

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the exercise package from internal/exercise
- hellostrings.go, solution.go are in package exercise
- Functions must be exported (capitalized) to be called from main
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the first argument (if provided)
	// DEBUG: In Variables panel, expand 'os.Args' to see all arguments

	// Get input string from command line arguments
	// os.Args is a slice of strings containing command-line arguments
	// Args[0] is the program name, Args[1] is first argument, etc.
	input := "hello world"
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	// Get function to run from command line arguments
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to inspect parsed values
	// DEBUG: Watch 'input' and 'demo' variables in Variables panel
	// DEBUG: These now hold the actual command-line values

	// ============================================================================
	// STEP 2: Example 1 - TitleCase (Simple)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through simple example
	// DEBUG: This is the simplest use case - capitalizing first letter of each word
	// DEBUG: Step Over (F10) to execute line by line
	// DEBUG: Watch Variables panel to see how 'result' changes

	if demo == "all" || demo == "titlecase" {
		fmt.Println("=== TitleCase Demo ===")
		fmt.Println("Capitalizes the first letter of each word")
		fmt.Println()

		testCases := []string{
			"hello world",
			"the quick brown fox",
			"Hello 👋 World",
			"café résumé",
			"日本語",
			input,  // User-provided text
		}

		// BREAKPOINT: Set breakpoint here before loop
		// DEBUG: Expand 'testCases' in Variables panel to see all test strings
		// DEBUG: Notice the len and cap of the slice

		for _, s := range testCases {
			// BREAKPOINT: Set breakpoint here inside loop
			// DEBUG: Watch 's' variable change with each iteration
			// DEBUG: Step Into (F11) on hellostrings.TitleCase() to see implementation

			// Call the TitleCase function from exercise package
			// BREAKPOINT: Set breakpoint here before function call
			// DEBUG: Step Into (F11) to enter hellostrings.TitleCase function
			// DEBUG: You'll see the implementation in hellostrings.go or solution.go
			result := hellostrings.TitleCase(s)

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
	// DEBUG: Step Into (F11) on hellostrings.Reverse() to see rune manipulation

	if demo == "all" || demo == "reverse" {
		fmt.Println("=== Reverse Demo ===")
		fmt.Println("Reverses strings (handles Unicode correctly)")
		fmt.Println()

		testCases := []string{
			"hello",
			"Hello 👋 World",
			"café",
			"日本語",
			input,
		}

		for _, s := range testCases {
			// BREAKPOINT: Set breakpoint here
			// DEBUG: Step Into (F11) to see reverse algorithm
			result := hellostrings.Reverse(s)

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

	if demo == "all" || demo == "runelen" {
		fmt.Println("=== RuneLen Demo ===")
		fmt.Println("Shows difference between byte length and rune (character) length")
		fmt.Println()

		testCases := []string{
			"hello",           // ASCII only
			"café",            // Contains é (2 bytes)
			"Hello 👋 World",  // Contains emoji (4 bytes)
			"日本語",           // Japanese (3 characters, 9 bytes)
			input,
		}

		for _, s := range testCases {
			// BREAKPOINT: Set breakpoint here
			// DEBUG: Step Into (F11) to see rune counting implementation
			byteLen := len(s)               // len() returns byte length
			runeLen := hellostrings.RuneLen(s)  // Count actual characters

			// BREAKPOINT: Set breakpoint here
			// DEBUG: Compare byteLen vs runeLen
			// DEBUG: For ASCII, they're equal. For Unicode, runeLen < byteLen

			fmt.Printf("%-25s → bytes: %2d, runes: %2d\n", s, byteLen, runeLen)
		}
		fmt.Println()
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
	// - Set breakpoint at hellostrings.TitleCase(s) call
	// - Press F11 (Step Into)
	// - Debugger will jump to hellostrings.go or solution.go
	// - You can now see the implementation
	// - Step through the implementation with F10
	// - Press Shift+F11 (Step Out) to return to main.go
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}