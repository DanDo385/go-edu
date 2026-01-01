package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Project: 12-pointers-zero-values-nil-gotchas")
	fmt.Println("To run tests: go test ./...")
	
	if len(os.Args) > 1 {
		fmt.Printf("Arguments provided: %v\n", os.Args[1:])
	} else {
		fmt.Println("No arguments provided. Usage: go run ./cmd/app/main.go [args...]")
	}
	
	fmt.Println("\nSee internal/pointerszerovaluesnilgotchas/exercise.go for the implementation.")
}
