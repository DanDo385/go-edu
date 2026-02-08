//go:build reference

package clitodofiles

/*
Reference Solution - CLI Todo List with File Persistence
======================================================

This file demonstrates building a CLI application with persistent JSON storage.
It combines file I/O, JSON serialization, in-memory data structures, and interface
design into a complete, testable application. This is the foundation for many
desktop and CLI tools.

This connects to the broader Go ecosystem by showing:
- How to structure data for JSON serialization with struct tags
- Interface-based design for testability (swap file store for mock)
- File system operations with proper error handling
- CRUD operations on in-memory collections with persistence

The exercise builds understanding of:
- JSON marshaling/unmarshaling and struct tags
- Slice operations: append, copy, iteration patterns
- Interface design: Store interface enables multiple implementations
- File I/O: os.ReadFile, os.WriteFile with error handling
- Data invariants: unique IDs, idempotent toggle operations

Time/Space Complexity:
- Load/Save: O(n) - JSON marshal/unmarshal of all items
- Add: O(n) to find max ID, O(1) to append
- Toggle: O(n) to find item by ID
- List: O(n) to filter items

Teaching notes:
- Memory/ownership: List returns a copy of items to prevent caller mutations
  from affecting the store's internal state. The items slice is our source of truth.
- Invariants: IDs are auto-incrementing and unique. Toggle is idempotent (toggling
  twice returns to original state). These invariants simplify reasoning about correctness.
- Error surfaces: File and JSON operations can fail. We propagate errors so the CLI
  can display meaningful messages to users.
*/

import (
	"encoding/json"
	"fmt"
	"os"
)

type Item struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Store interface {
	Load() error
	Save() error
	Add(text string) Item
	Toggle(id int) (Item, bool)
	List(onlyPending bool) []Item
}

type fileStore struct {
	path  string // File path for persistence
	items []Item // In-memory storage (slice)
}

/*
NewFileStore - Constructor for File-Backed Store

Creates a Store implementation that persists items to a JSON file.
The path parameter specifies where the JSON file lives on disk.

Returns a *fileStore which satisfies the Store interface.
This enables the CLI to use file storage while tests can provide mock implementations.
*/
func NewFileStore(path string) Store {
	// Initialize with empty slice - make([]Item, 0) creates slice with len=0, cap=0
	// Will grow on first append; alternatively could use make([]Item, 0, 16) for initial cap
	return &fileStore{path: path, items: make([]Item, 0)}
}

/*
Load - Deserialize JSON File into Memory

Reads the JSON file from disk and populates the in-memory items slice.
This is called at startup to restore previous session state.

Error handling: File might not exist (first run), might be corrupted,
or JSON might be malformed. We surface all errors for caller to handle.
*/
func (fs *fileStore) Load() error {
	// Read entire file into memory - fine for small todo lists
	// For large files, would use streaming JSON decoder
	data, err := os.ReadFile(fs.path)
	if err != nil {
		// File not found, permission denied, etc.
		return err
	}

	// Parse JSON bytes into []Item slice
	// json.Unmarshal uses reflection to map JSON keys to struct fields
	// Struct tags (`json:"id"`) tell Unmarshal which JSON key maps to which field
	var items []Item
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal todos: %w", err)
	}

	// Replace in-memory state with loaded data
	// This overwrites any existing items - file is source of truth
	fs.items = items
	return nil
}

/*
Save - Serialize In-Memory State to JSON File

Writes the current items slice to the JSON file.
MarshalIndent produces human-readable output (pretty-printed) for manual editing.

File permissions: 0644 = owner read/write, group/others read only
*/
func (fs *fileStore) Save() error {
	// Marshal items to JSON with indentation for readability
	// Second param "" = no prefix, third param "  " = 2-space indent
	data, err := json.MarshalIndent(fs.items, "", "  ")
	if err != nil {
		// Unlikely for simple structs, but custom types might fail
		return fmt.Errorf("marshal todos: %w", err)
	}

	// Write to file - overwrites if exists, creates if not
	// 0644 = rw-r--r-- (owner: read+write, others: read only)
	if err := os.WriteFile(fs.path, data, 0644); err != nil {
		return fmt.Errorf("write todos: %w", err)
	}
	return nil
}

/*
Add - Append New Todo Item with Auto-Incrementing ID

Adds a new item to the list with the next available ID.
We find max ID by linear scan - O(n) but acceptable for typical list sizes.
Alternative: maintain a separate maxID counter, but that adds state to persist.
*/
func (fs *fileStore) Add(text string) Item {
	// Find highest existing ID to ensure uniqueness
	// Start at 0 so first item gets ID 1
	maxID := 0
	for _, item := range fs.items {
		if item.ID > maxID {
			maxID = item.ID
		}
	}

	// Create new item with next ID
	// Done defaults to false (Go zero value for bool)
	newItem := Item{
		ID:   maxID + 1,
		Text: text,
		Done: false,
	}

	// Append to slice - may trigger reallocation if at capacity
	fs.items = append(fs.items, newItem)

	// Return copy of new item for caller to display
	return newItem
}

/*
Toggle - Flip Completion Status by ID

Marks an item as done/not-done. Idempotent: toggling twice returns to original state.
Returns (Item, true) if found and toggled, (zero Item, false) if ID not found.

We use range over indices (for i := range) because we need to mutate fs.items[i].
*/
func (fs *fileStore) Toggle(id int) (Item, bool) {
	for i := range fs.items {
		if fs.items[i].ID == id {
			// Flip the Done boolean: !true = false, !false = true
			fs.items[i].Done = !fs.items[i].Done
			// Return the updated item and success indicator
			return fs.items[i], true
		}
	}
	// ID not found - return zero value and false
	return Item{}, false
}

/*
List - Return Items with Optional Filtering

Returns a COPY (per .cursorrules): out := make([]Item, len(fs.items)); copy(out, fs.items).
Why? If we returned fs.items directly, the caller would get a slice that shares
the SAME backing array as our internal state. Mutating out[0].Done would change
fs.items[0].Done — we'd have multiple owners of mutable data. By copying, we
give the caller an independent slice. Mutations to the returned slice don't
affect our internal items. Think of it as "hand over a clone, not the original."
*/
func (fs *fileStore) List(onlyPending bool) []Item {
	if !onlyPending {
		// Return full list - allocate exact size and copy
		out := make([]Item, len(fs.items))
		copy(out, fs.items)
		return out
	}

	// Filter to pending items only
	// Pre-allocate with capacity for worst case (all pending)
	out := make([]Item, 0, len(fs.items))
	for _, item := range fs.items {
		if !item.Done {
			out = append(out, item)
		}
	}
	return out
}
