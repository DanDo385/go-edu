package main

import (
	"flag"
	"fmt"
	"os"
	
	// Import the internal package
	_ "github.com/example/go-10x-minis/minis/39-sha256-hasher/internal/sha256hasher"
)

func main() {
	// Custom CLI Arguments
	// TODO: Add your custom flags here.
	// Example:
	// var input = flag.String("input", "", "Input value")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  go run ./cmd/app/main.go --help")
	}
	
	flag.Parse()

	fmt.Println("Running 39-sha256-hasher...")
	
	// TODO: Call your internal package function here using the flags
	// Example:
	// sha256hasher.Run(*input)
}
