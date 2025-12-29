package main

import (
	"flag"   // Command-line argument parsing
	"fmt"    // Formatted I/O
	"log"    // Logging
	"os"     // Operating system interface

	// Import the exercise package from parent directory
	"github.com/example/go-10x-minis/minis/05-cli-todo-files/exercise"
)

/*
===================================================================
CLI Todo List: File-Based Task Manager
===================================================================

This program demonstrates file I/O, JSON persistence, and command-line
interfaces by implementing a persistent todo list stored in a JSON file.

USAGE:
    go run ./cmd/cli-todo-files/main.go -list
    go run ./cmd/cli-todo-files/main.go -add="Buy groceries"
    go run ./cmd/cli-todo-files/main.go -toggle=1
    go run ./cmd/cli-todo-files/main.go -list -all
    go run ./cmd/cli-todo-files/main.go -file=custom-todos.json -list

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see JSON encoding/decoding
- See /RUN_DEBUG.md for comprehensive debugging guide

PACKAGE STRUCTURE:
- This file (cmd/cli-todo-files/main.go) is package main
- It imports the exercise package from the parent directory
- exercise.go, solution.go are in package exercise
- Functions must be exported (capitalized) to be called from main
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect default values before parsing
	// DEBUG: Before flag.Parse(), all flags have their default values
	// DEBUG: In Variables panel, you'll see: file="todos.json", add="", toggle=0

	// Add a new todo item
	addText := flag.String("add", "", "Add a new todo item")

	// Toggle done status of item by ID
	toggleID := flag.Int("toggle", 0, "Toggle done status of item by ID")

	// List all todo items
	listCmd := flag.Bool("list", false, "List all todo items")

	// Show completed items (use with -list)
	showAll := flag.Bool("all", false, "Show completed items (use with -list)")

	// Path to todo file
	filePath := flag.String("file", "todos.json", "Path to todo file")

	// Use solution.go instead of exercise.go
	useSolution := flag.Bool("solution", false, "Use solution.go instead of exercise.go (requires -tags=solution)")

	flag.Parse()

	// BREAKPOINT: Set breakpoint here to inspect parsed values in Variables panel
	// DEBUG: After flag.Parse(), variables now hold actual command-line values
	// DEBUG: Watch 'addText', 'toggleID', 'listCmd', 'showAll', 'filePath', 'useSolution'

	// ============================================================================
	// STEP 2: Create FileStore Instance
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before creating store
	// DEBUG: Watch 'filePath' - this is where todos will be persisted
	// DEBUG: Step Into (F11) on exercise.NewFileStore() to see initialization

	store := exercise.NewFileStore(*filePath)

	// BREAKPOINT: Set breakpoint here after creating store
	// DEBUG: Inspect 'store' variable - see its fields (filepath, items slice)
	// DEBUG: Notice items slice is initially empty (len=0)

	// ============================================================================
	// STEP 3: Load Existing Data from File
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before loading
	// DEBUG: Step Into (F11) on store.Load() to see JSON deserialization
	// DEBUG: Watch how file is read and decoded into Go structs

	if err := store.Load(); err != nil {
		// File not existing is OK (first run), other errors are fatal
		// BREAKPOINT: Set breakpoint here in error handler
		// DEBUG: Check os.IsNotExist(err) - true means file doesn't exist yet
		// DEBUG: This is expected on first run - no todos.json file yet

		if !os.IsNotExist(err) {
			log.Fatalf("Failed to load todos: %v", err)
		}
		// File doesn't exist - that's OK, we'll create it on first save
		fmt.Printf("(No existing todo file found at %s - will create on first save)\n\n", *filePath)
	}

	// BREAKPOINT: Set breakpoint here after loading
	// DEBUG: Inspect 'store' - items slice may now be populated from file
	// DEBUG: Expand 'store' → 'items' to see loaded todo items
	// DEBUG: Each item has: ID (int), Text (string), Done (bool)

	// ============================================================================
	// STEP 4: Handle Commands (switch statement pattern)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before command handling
	// DEBUG: Watch which flag is set - determines which case executes
	// DEBUG: Only one command should be active at a time

	switch {
	// ------------------------------------------------------------------------
	// Command: Add new todo item
	// ------------------------------------------------------------------------
	case *addText != "":
		// BREAKPOINT: Set breakpoint here in add command
		// DEBUG: Watch 'addText' - this is the new todo text from user
		// DEBUG: Step Into (F11) on store.Add() to see how item is created

		fmt.Printf("Adding new todo: %q\n", *addText)

		item := store.Add(*addText)

		// BREAKPOINT: Set breakpoint here after adding
		// DEBUG: Inspect 'item' - see ID, Text, Done fields
		// DEBUG: Notice item.Done is false (new items start as not done)
		// DEBUG: ID is auto-incremented based on existing items

		fmt.Printf("Added: [%d] %s\n", item.ID, item.Text)

		// BREAKPOINT: Set breakpoint here before saving
		// DEBUG: Step Into (F11) on store.Save() to see JSON serialization
		// DEBUG: Watch how Go structs are encoded to JSON and written to file

		if err := store.Save(); err != nil {
			log.Fatalf("Failed to save: %v", err)
		}

		fmt.Println("✓ Saved to file")

	// ------------------------------------------------------------------------
	// Command: Toggle done status
	// ------------------------------------------------------------------------
	case *toggleID > 0:
		// BREAKPOINT: Set breakpoint here in toggle command
		// DEBUG: Watch 'toggleID' - this is the ID of item to toggle
		// DEBUG: Step Into (F11) on store.Toggle() to see lookup and update

		fmt.Printf("Toggling todo item #%d\n", *toggleID)

		item, found := store.Toggle(*toggleID)

		// BREAKPOINT: Set breakpoint here after toggling
		// DEBUG: Inspect 'found' - should be true if item exists
		// DEBUG: Inspect 'item.Done' - flipped from previous state
		// DEBUG: If found=false, item ID doesn't exist in store

		if !found {
			log.Fatalf("Item %d not found", *toggleID)
		}

		status := "not done"
		if item.Done {
			status = "done ✓"
		}

		fmt.Printf("Toggled: [%d] %s → %s\n", item.ID, item.Text, status)

		// BREAKPOINT: Set breakpoint here before saving
		// DEBUG: Changes are only in memory until we call Save()
		// DEBUG: Step Into (F11) to see persistence to disk

		if err := store.Save(); err != nil {
			log.Fatalf("Failed to save: %v", err)
		}

		fmt.Println("✓ Saved to file")

	// ------------------------------------------------------------------------
	// Command: List all items
	// ------------------------------------------------------------------------
	case *listCmd:
		// BREAKPOINT: Set breakpoint here in list command
		// DEBUG: Watch 'showAll' flag - determines if completed items show
		// DEBUG: Step Into (F11) on store.List() to see filtering logic

		items := store.List(!*showAll)

		// BREAKPOINT: Set breakpoint here after getting list
		// DEBUG: Inspect 'items' slice - expand to see all todo items
		// DEBUG: Notice filtering: if showAll=false, only incomplete items
		// DEBUG: Each item is a pointer: *TodoItem

		if len(items) == 0 {
			if *showAll {
				fmt.Println("No items in todo list")
			} else {
				fmt.Println("No incomplete items (use -all to show completed)")
			}
			return
		}

		// Display header
		if *showAll {
			fmt.Println("=== All Todo Items ===")
		} else {
			fmt.Println("=== Incomplete Todo Items ===")
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here before display loop
		// DEBUG: Expand 'items' in Variables panel to see all entries
		// DEBUG: Notice slice of pointers: []*TodoItem

		for _, item := range items {
			// BREAKPOINT: Set breakpoint here inside display loop
			// DEBUG: Watch 'item' change with each iteration
			// DEBUG: See item.ID (int), item.Text (string), item.Done (bool)

			status := " "
			if item.Done {
				status = "✓"
			}

			fmt.Printf("[%d] [%s] %s\n", item.ID, status, item.Text)
		}

		fmt.Println()
		fmt.Printf("Total items: %d\n", len(items))

		// Show count of hidden completed items if applicable
		// BREAKPOINT: Set breakpoint here before counting completed
		// DEBUG: We're about to count how many completed items exist

		if !*showAll {
			allItems := store.List(false) // Get all items
			completedCount := len(allItems) - len(items)

			// BREAKPOINT: Set breakpoint here after counting
			// DEBUG: Inspect 'completedCount' - number of hidden completed items
			// DEBUG: This is len(all) - len(incomplete)

			if completedCount > 0 {
				fmt.Printf("(%d completed item", completedCount)
				if completedCount > 1 {
					fmt.Print("s")
				}
				fmt.Println(" hidden - use -all to show)")
			}
		}

	// ------------------------------------------------------------------------
	// Default: Show usage if no command given
	// ------------------------------------------------------------------------
	default:
		// BREAKPOINT: Set breakpoint here in default case
		// DEBUG: This executes when no valid command flags are provided
		// DEBUG: Helps user understand how to use the program

		fmt.Println("CLI Todo List Manager")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  Add item:    -add=\"Buy groceries\"")
		fmt.Println("  List items:  -list")
		fmt.Println("  Show all:    -list -all")
		fmt.Println("  Toggle done: -toggle=1")
		fmt.Println("  Custom file: -file=custom.json")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	// ============================================================================
	// STEP 5: Usage Note about Solution
	// ============================================================================
	if *useSolution {
		fmt.Println()
		fmt.Println("NOTE: To use solution.go, build/run with -tags=solution:")
		fmt.Println("  go run -tags=solution ./cmd/cli-todo-files/main.go -solution=true")
		fmt.Println("  go test -tags=solution")
		os.Exit(0)
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES FOR FILE PERSISTENCE AND CLI:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise package functions
	// 5. Watch Variables panel to see struct fields and slice contents
	// 6. Add watch expressions: len(items), items[0].Done, store.filepath
	// 7. Use Debug Console to inspect: items[0], store
	//
	// UNDERSTANDING FILE PERSISTENCE:
	// - Step Into (F11) store.Load() to see JSON deserialization
	// - Watch how json.Decoder reads file and creates Go structs
	// - Step Into (F11) store.Save() to see JSON serialization
	// - Watch how Go structs are encoded to JSON and written
	// - Check todos.json file after save to see JSON format
	//
	// UNDERSTANDING POINTERS:
	// - store.List() returns []*TodoItem (slice of pointers)
	// - Pointers allow modification of original data
	// - In Variables panel, expand pointer to see actual struct
	// - Notice &TodoItem{...} creates pointer to struct
	//
	// UNDERSTANDING SWITCH ON BOOLEANS:
	// - switch without expression evaluates each case as boolean
	// - First true case executes (like if-else chain)
	// - Only one case executes - mutual exclusion
	// - Set breakpoint in each case to see which executes
	//
	// UNDERSTANDING AUTO-INCREMENT IDs:
	// - Step Into (F11) store.Add() to see ID generation
	// - Watch how maxID is found by iterating existing items
	// - New ID = maxID + 1 ensures uniqueness
	//
	// TESTING PERSISTENCE:
	// - Run: go run main.go -add="Test item"
	// - Check todos.json file exists and contains JSON
	// - Run: go run main.go -list
	// - See item persisted across program runs
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}
