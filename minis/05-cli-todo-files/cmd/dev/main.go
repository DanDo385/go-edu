package main

import (
	"fmt"
	"os"

	"github.com/example/go-10x-minis/minis/05-cli-todo-files/internal/clitodofiles"
)

/*
Debug Harness for CLI TODO Files Module

This file demonstrates the todo store functionality.
Set breakpoints in exercise.go and press F5 in VS Code.
*/
func main() {
	fmt.Println("=== CLI TODO Files Debug Harness ===")
	fmt.Println()

	// Use a temporary file
	tmpFile := "/tmp/debug_todos.json"
	defer os.Remove(tmpFile)

	store := clitodofiles.NewFileStore(tmpFile)

	// Demo 1: Add items
	fmt.Println("--- Adding Items ---")
	item1 := store.Add("Learn Go basics")
	fmt.Printf("Added: #%d %s\n", item1.ID, item1.Text)

	item2 := store.Add("Build a web server")
	fmt.Printf("Added: #%d %s\n", item2.ID, item2.Text)

	item3 := store.Add("Deploy to production")
	fmt.Printf("Added: #%d %s\n", item3.ID, item3.Text)
	fmt.Println()

	// Demo 2: Save and reload
	fmt.Println("--- Save and Reload ---")
	if err := store.Save(); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		return
	}
	fmt.Println("Saved to disk")

	store2 := clitodofiles.NewFileStore(tmpFile)
	if err := store2.Load(); err != nil {
		fmt.Printf("Error loading: %v\n", err)
		return
	}
	fmt.Println("Loaded from disk")
	fmt.Println()

	// Demo 3: Toggle items
	fmt.Println("--- Toggle Items ---")
	if item, found := store2.Toggle(1); found {
		fmt.Printf("Toggled #%d: done=%v\n", item.ID, item.Done)
	}
	if item, found := store2.Toggle(2); found {
		fmt.Printf("Toggled #%d: done=%v\n", item.ID, item.Done)
	}
	fmt.Println()

	// Demo 4: List items
	fmt.Println("--- All Items ---")
	for _, item := range store2.List(false) {
		status := "[ ]"
		if item.Done {
			status = "[✓]"
		}
		fmt.Printf("%s #%d: %s\n", status, item.ID, item.Text)
	}
	fmt.Println()

	fmt.Println("--- Pending Only ---")
	for _, item := range store2.List(true) {
		fmt.Printf("[ ] #%d: %s\n", item.ID, item.Text)
	}
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. json.Marshal/Unmarshal for serialization")
	fmt.Println("2. os.ReadFile/WriteFile for simple file I/O")
	fmt.Println("3. Interface-based design (Store interface)")
	fmt.Println("4. Auto-incrementing IDs in slices")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/06-worker-pool-wordcount")
}
