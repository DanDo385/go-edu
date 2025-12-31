// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/clitodofiles.go
// 2. Open this file (cmd/debug/main.go)
// 3. Use "Debug Main Program (Current Package)" configuration
// 4. Press F5 - that's it! The debugger will stop at your breakpoints
//
// Usage:
//   go run ./cmd/debug
//   # Or just press F5 in VS Code

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/go-10x-minis/minis/05-cli-todo-files/internal/clitodofiles"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	tmpFile := filepath.Join(os.TempDir(), "debug-todos.json")
	defer os.Remove(tmpFile) // Clean up

	fmt.Println("=== Debugging FileStore ===")
	fmt.Printf("Using file: %s\n\n", tmpFile)

	// Set breakpoint in clitodofiles.go at NewFileStore function
	store := clitodofiles.NewFileStore(tmpFile)

	// Add default test todos
	item1 := store.Add("Buy groceries")
	fmt.Printf("Added: [%d] %s\n", item1.ID, item1.Text)

	item2 := store.Add("Finish homework")
	fmt.Printf("Added: [%d] %s\n", item2.ID, item2.Text)

	// Save to file
	if err := store.Save(); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		return
	}
	fmt.Println("✓ Saved to file")

	// List todos
	items := store.List(false) // false = only incomplete
	fmt.Printf("\nIncomplete todos: %+v\n", items)

	// Toggle one item
	if len(items) > 0 {
		item, found := store.Toggle(items[0].ID)
		if found {
			fmt.Printf("Toggled: [%d] %s (done: %v)\n", item.ID, item.Text, item.Done)
		}
	}
}

