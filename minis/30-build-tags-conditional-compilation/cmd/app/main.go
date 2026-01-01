package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Project: 30-build-tags-conditional-compilation")
	fmt.Println("To run tests: go test ./...")
	
	if len(os.Args) > 1 {
		fmt.Printf("Arguments provided: %v\n", os.Args[1:])
	} else {
		fmt.Println("No arguments provided. Usage: go run ./cmd/app/main.go [args...]")
	}
	
	fmt.Println("\nSee internal/buildtagsconditionalcompilation/exercise.go for the implementation.")
}
