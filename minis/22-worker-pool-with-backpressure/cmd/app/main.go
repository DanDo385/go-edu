package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Project: 22-worker-pool-with-backpressure")
	fmt.Println("To run tests: go test ./...")
	
	if len(os.Args) > 1 {
		fmt.Printf("Arguments provided: %v\n", os.Args[1:])
	} else {
		fmt.Println("No arguments provided. Usage: go run ./cmd/app/main.go [args...]")
	}
	
	fmt.Println("\nSee internal/workerpoolwithbackpressure/exercise.go for the implementation.")
}
