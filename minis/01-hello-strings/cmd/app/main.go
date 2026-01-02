package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/example/go-10x-minis/minis/01-hello-strings/internal/hellostrings"
)

/*
UTF-8 String Utilities CLI

Usage:
  go run ./cmd/app/main.go <command> <text>

Commands:
  titlecase   Capitalize first letter of each word
  reverse     Reverse string character-by-character
  runelen     Count characters (runes), not bytes

Examples:
  go run ./cmd/app/main.go titlecase "hello world"
  go run ./cmd/app/main.go reverse "hello 世界"
  go run ./cmd/app/main.go runelen "café 👋"
*/
func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	text := strings.Join(os.Args[2:], " ")

	switch command {
	case "titlecase":
		result := hellostrings.TitleCase(text)
		fmt.Printf("Input:  %q\n", text)
		fmt.Printf("Output: %q\n", result)

	case "reverse":
		result := hellostrings.Reverse(text)
		fmt.Printf("Input:  %q\n", text)
		fmt.Printf("Output: %q\n", result)

	case "runelen":
		result := hellostrings.RuneLen(text)
		fmt.Printf("Input:      %q\n", text)
		fmt.Printf("Byte count: %d\n", len(text))
		fmt.Printf("Rune count: %d\n", result)

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("UTF-8 String Utilities CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/app/main.go <command> <text>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  titlecase   Capitalize first letter of each word")
	fmt.Println("  reverse     Reverse string character-by-character")
	fmt.Println("  runelen     Count characters (runes), not bytes")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/app/main.go titlecase \"hello world\"")
	fmt.Println("  go run ./cmd/app/main.go reverse \"hello 世界\"")
	fmt.Println("  go run ./cmd/app/main.go runelen \"café 👋\"")
}
