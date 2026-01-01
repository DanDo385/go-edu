package main

import (
	"fmt"
	
	_ "github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/internal/pluginsystemhotreload"
)

func main() {
	fmt.Println("=== 47-plugin-system-hot-reload ===")
	fmt.Println("This is the application entry point.")
	fmt.Println("Implement your solution in internal/pluginsystemhotreload/exercise.go")
	fmt.Println("Run tests: go test ./internal/pluginsystemhotreload/...")
}
