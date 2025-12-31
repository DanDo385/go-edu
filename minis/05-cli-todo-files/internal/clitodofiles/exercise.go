//go:build !solution
// +build !solution

package clitodofiles

// TODO: Import required packages
// You'll need:
// - "encoding/json" for JSON marshaling/unmarshaling
// - "fmt" for error formatting
// - "os" for file I/O operations
//
// Uncomment these imports when you implement the functions:
/*
// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )
*/

// Item represents a single TODO item
type Item struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Store defines operations for managing TODO items
//
// Key Design Pattern: Interface-based abstraction allows swapping
// storage implementations without changing client code
type Store interface {
	Load() error
	Save() error
	Add(text string) Item
	Toggle(id int) (Item, bool)
	List(onlyPending bool) []Item
}

// fileStore persists items to a JSON file
type fileStore struct {
	path  string  // File path for JSON storage
	items []Item  // In-memory cache of all items
}

// NewFileStore creates a Store backed by a JSON file at path
// TODO: Implement this constructor function
//
// Algorithm:
// 1. Create a fileStore struct with given path
// 2. Initialize empty items slice
// 3. Return the fileStore as a Store interface
func NewFileStore(path string) Store {
	// TODO: Implement NewFileStore
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Load reads items from disk
// TODO: Implement this function
//
// Algorithm:
// 1. Read entire file using os.ReadFile
// 2. Parse JSON array into fs.items slice
// 3. Return error if file read/parse fails
//
// CS Concepts:
// - Persistence: Data survives program restarts via disk storage
// - Serialization: Converting in-memory objects to portable format
// - Idempotence: Multiple loads should work correctly (replace, not append)
func (fs *fileStore) Load() error {
	// TODO: Step 1 - Read file using os.ReadFile

	// TODO: Step 2 - Unmarshal JSON into fs.items

	// TODO: Step 3 - Return error or nil

	return nil
}

// Save writes items to disk
// TODO: Implement this function
//
// Algorithm:
// 1. Marshal fs.items to JSON with indentation
// 2. Write JSON bytes to file with proper permissions
// 3. Return error if marshal or write fails
//
// CS Concepts:
// - Serialization: Converting objects to portable text format
// - File Permissions: Controlling read/write access (0644 = rw-r--r--)
// - Atomicity: Ideally should be atomic (not shown here for simplicity)
func (fs *fileStore) Save() error {
	// TODO: Step 1 - Marshal items to indented JSON

	// TODO: Step 2 - Write JSON to file with permissions

	// TODO: Step 3 - Return error or nil

	return nil
}

// Add inserts a new item with an auto-incremented ID
// TODO: Implement this function
//
// Algorithm:
// 1. Find the maximum ID in current items
// 2. Create new item with ID = maxID + 1
// 3. Append to items slice
// 4. Return the new item
//
// CS Concepts:
// - ID Generation: Simple auto-increment for uniqueness
// - Append Operation: Adding to slice may trigger reallocation
// - Linear Search: O(n) scan to find maximum (acceptable for small lists)
func (fs *fileStore) Add(text string) Item {
	// TODO: Step 1 - Find maximum existing ID

	// TODO: Step 2 - Create new item with nextID and given text

	// TODO: Step 3 - Append to items slice

	// TODO: Step 4 - Return new item

	return Item{}
}

// Toggle flips the Done status for the given ID
// TODO: Implement this function
//
// Algorithm:
// 1. Search items slice for matching ID
// 2. If found, flip Done status and return (item, true)
// 3. If not found, return (zero value, false)
//
// CS Concepts:
// - Linear Search: O(n) time, simple implementation
// - Value Semantics: Modifying struct in slice requires careful handling
// - Optional Pattern: Returning (value, ok) instead of error for "not found"
func (fs *fileStore) Toggle(id int) (Item, bool) {
	// TODO: Step 1 - Loop through items slice

	// TODO: Step 2 - Check if ID matches

	// TODO: Step 3 - Flip Done status if found

	// TODO: Step 4 - Return (item, true) if found

	// TODO: Step 5 - Return (zero value, false) if not found

	return Item{}, false
}

// List returns all items, optionally filtering completed ones
// TODO: Implement this function
//
// Algorithm:
// 1. If onlyPending is false, return all items
// 2. If onlyPending is true, filter to items where Done == false
// 3. Return the resulting slice
//
// CS Concepts:
// - Filtering: Selecting subset based on predicate
// - Memory Sharing: Returned slice may share backing array
// - Optional Filtering: Different return types based on parameter
func (fs *fileStore) List(onlyPending bool) []Item {
	// TODO: Step 1 - Check onlyPending flag

	// TODO: Step 2 - Return all items if not filtering

	// TODO: Step 3 - Build filtered slice if filtering pending only

	// TODO: Step 4 - Return filtered result

	return nil
}

// After implementing:
// - Run: go test -v ./...
// - Build CLI: cd ../cmd/cli-todo-files && go build
// - Test manually:
//   $ ./cli-todo-files -add "Learn Go"
//   $ ./cli-todo-files -list
// - Compare with solution.go to see detailed explanations
