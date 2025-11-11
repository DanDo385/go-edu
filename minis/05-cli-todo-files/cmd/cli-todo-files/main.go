package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/example/go-10x-minis/minis/05-cli-todo-files/exercise"
)

func main() {
	// Define CLI flags
	addText := flag.String("add", "", "Add a new todo item")
	toggleID := flag.Int("toggle", 0, "Toggle done status of item by ID")
	listCmd := flag.Bool("list", false, "List all todo items")
	showAll := flag.Bool("all", false, "Show completed items (use with -list)")
	filePath := flag.String("file", "todos.json", "Path to todo file")

	flag.Parse()

	// Create store
	store := exercise.NewFileStore(*filePath)

	// Load existing data
	if err := store.Load(); err != nil {
		// File not existing is OK (first run), other errors are fatal
		if !os.IsNotExist(err) {
			log.Fatalf("Failed to load todos: %v", err)
		}
	}

	// Handle commands
	switch {
	case *addText != "":
		item := store.Add(*addText)
		fmt.Printf("Added: [%d] %s\n", item.ID, item.Text)
		if err := store.Save(); err != nil {
			log.Fatalf("Failed to save: %v", err)
		}

	case *toggleID > 0:
		item, found := store.Toggle(*toggleID)
		if !found {
			log.Fatalf("Item %d not found", *toggleID)
		}
		status := "not done"
		if item.Done {
			status = "done"
		}
		fmt.Printf("Toggled: [%d] %s (%s)\n", item.ID, item.Text, status)
		if err := store.Save(); err != nil {
			log.Fatalf("Failed to save: %v", err)
		}

	case *listCmd:
		items := store.List(!*showAll)
		if len(items) == 0 {
			fmt.Println("No items to show")
			return
		}
		for _, item := range items {
			status := " "
			if item.Done {
				status = "✓"
			}
			fmt.Printf("[%d] [%s] %s\n", item.ID, status, item.Text)
		}

	default:
		flag.Usage()
		os.Exit(1)
	}
}
