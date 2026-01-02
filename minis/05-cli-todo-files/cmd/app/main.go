package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/example/go-10x-minis/minis/05-cli-todo-files/internal/clitodofiles"
)

/*
TODO List CLI

Usage:
  go run ./cmd/app/main.go add <text>
  go run ./cmd/app/main.go list [--pending]
  go run ./cmd/app/main.go toggle <id>

Commands:
  add       Add a new todo item
  list      List all todos (--pending for pending only)
  toggle    Toggle completion status by ID

Examples:
  go run ./cmd/app/main.go add "Buy groceries"
  go run ./cmd/app/main.go list
  go run ./cmd/app/main.go list --pending
  go run ./cmd/app/main.go toggle 1
*/

const dataFile = "todos.json"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	store := clitodofiles.NewFileStore(dataFile)
	if err := store.Load(); err != nil {
		// File might not exist yet, that's OK
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: missing text for todo")
			printUsage()
			os.Exit(1)
		}
		text := os.Args[2]
		item := store.Add(text)
		if err := store.Save(); err != nil {
			fmt.Printf("Error saving: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added todo #%d: %s\n", item.ID, item.Text)

	case "list":
		onlyPending := len(os.Args) > 2 && os.Args[2] == "--pending"
		items := store.List(onlyPending)
		if len(items) == 0 {
			if onlyPending {
				fmt.Println("No pending todos!")
			} else {
				fmt.Println("No todos yet. Add one with: go run ./cmd/app/main.go add \"task\"")
			}
			return
		}
		fmt.Println("=== TODO List ===")
		for _, item := range items {
			status := "[ ]"
			if item.Done {
				status = "[✓]"
			}
			fmt.Printf("%s #%d: %s\n", status, item.ID, item.Text)
		}

	case "toggle":
		if len(os.Args) < 3 {
			fmt.Println("Error: missing ID")
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		item, found := store.Toggle(id)
		if !found {
			fmt.Printf("Todo #%d not found\n", id)
			os.Exit(1)
		}
		if err := store.Save(); err != nil {
			fmt.Printf("Error saving: %v\n", err)
			os.Exit(1)
		}
		status := "pending"
		if item.Done {
			status = "done"
		}
		fmt.Printf("Todo #%d is now %s: %s\n", item.ID, status, item.Text)

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("TODO List CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/app/main.go add <text>")
	fmt.Println("  go run ./cmd/app/main.go list [--pending]")
	fmt.Println("  go run ./cmd/app/main.go toggle <id>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add       Add a new todo item")
	fmt.Println("  list      List all todos (--pending for pending only)")
	fmt.Println("  toggle    Toggle completion status by ID")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/app/main.go add \"Buy groceries\"")
	fmt.Println("  go run ./cmd/app/main.go list")
	fmt.Println("  go run ./cmd/app/main.go toggle 1")
}
