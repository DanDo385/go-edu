package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== Build Tags & Conditional Compilation Demo ===")
	fmt.Println()
	fmt.Printf("GOOS:   %s\n", runtime.GOOS)
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)
	fmt.Println()
	fmt.Println("This project demonstrates build tags and conditional compilation.")
	fmt.Println("Try building with different tags:")
	fmt.Println("  go build -tags=debug ./cmd/app")
	fmt.Println("  go build -tags=prod ./cmd/app")
	fmt.Println("  go build -tags=cloud ./cmd/app")
}
