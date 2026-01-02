package main

import (
	"fmt"
	"os"

	"minis/01-hello-strings/internal/hellostrings"
)

/*
Hello Strings CLI

This application demonstrates UTF-8-aware string utilities including title case conversion,
string reversal, and rune counting.

Usage:

	go run ./cmd/app/main.go <command> [input]

Commands:

	titlecase <string>   Convert string to title case
	reverse <string>     Reverse a string (UTF-8 aware)
	runelen <string>     Count runes (characters) in a string

Examples:

	# Title case conversion
	go run ./cmd/app/main.go titlecase "hello world"

	# Reverse a string
	go run ./cmd/app/main.go reverse "Hello 👋 World"

	# Count runes
	go run ./cmd/app/main.go runelen "Hello 👋 World"

Copy & Paste Examples:

	go run ./cmd/app/main.go titlecase "hello world"
	go run ./cmd/app/main.go reverse "café résumé"
	go run ./cmd/app/main.go runelen "Hello 👋 World"
*/

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./cmd/app/main.go <command> <input>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  titlecase <string>   Convert string to title case")
		fmt.Println("  reverse <string>     Reverse a string (UTF-8 aware)")
		fmt.Println("  runelen <string>     Count runes (characters) in a string")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println(`  go run ./cmd/app/main.go titlecase "hello world"`)
		fmt.Println(`  go run ./cmd/app/main.go reverse "Hello 👋 World"`)
		fmt.Println(`  go run ./cmd/app/main.go runelen "Hello 👋 World"`)
		os.Exit(1)
	}

	command := os.Args[1]
	input := os.Args[2]

	switch command {
	case "titlecase":
		result := hellostrings.TitleCase(input)
		fmt.Printf("Input:  %q\n", input)
		fmt.Printf("Output: %q\n", result)

	case "reverse":
		result := hellostrings.Reverse(input)
		fmt.Printf("Input:  %q\n", input)
		fmt.Printf("Output: %q\n", result)

	case "runelen":
		result := hellostrings.RuneLen(input)
		fmt.Printf("Input:  %q\n", input)
		fmt.Printf("Runes:  %d\n", result)
		fmt.Printf("Bytes:  %d\n", len(input))

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Valid commands: titlecase, reverse, runelen")
		os.Exit(1)
	}
}
