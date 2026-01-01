//go:build !solution && !reference

package clitodofiles



import (
	"encoding/json"
	"fmt"
	"os"
)

// Item represents a single TODO item.
type Item struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Store defines operations for managing TODO items.
// This interface enables:
// - Testing with mock implementations (no real file I/O)
// - Swapping backends (file → database → API) without changing CLI code
// - Clear API contract (documentation through types)
type Store interface {
	Load() error
	Save() error
	Add(text string) Item
	Toggle(id int) (Item, bool)
	List(onlyPending bool) []Item
}

// fileStore is the concrete implementation backed by a JSON file.
// Go Concepts Demonstrated:
// - Struct fields (state)
// - Pointer receivers (methods that mutate state)
// - Interface implementation (no explicit "implements" keyword!)
type fileStore struct {
	path  string // File path for persistence
	items []Item // In-memory storage (slice)
}

// NewFileStore creates a Store backed by a JSON file.
//
// Go Concepts Demonstrated:
// - Constructor pattern (factory function)
// - Returning interface type (hides implementation)
// - Pointer allocation (struct fields are zero-initialized)
func NewFileStore(path string) Store {
	// TODO: Implement this function
	panic("unimplemented")
}


func (fs *fileStore) Load() error {
	// TODO: Implement this function
	panic("unimplemented")
}

// Save writes items to the JSON file.
//
// Go Concepts Demonstrated:
// - json.MarshalIndent: Serialize with pretty-printing (readable files)
// - os.WriteFile: Atomic write with specified permissions
// - Error wrapping with fmt.Errorf and %w
//
// Why not streaming?
// For small datasets (<10k items), marshaling to memory then writing is simple.
// For large datasets, use json.Encoder for streaming writes.
func (fs *fileStore) Save() error {
	// TODO: Implement this function
	panic("unimplemented")
}

// Add creates a new item and returns it.
//
// Go Concepts Demonstrated:
// - Slice append (grows capacity automatically)
// - ID generation (find max + 1)
// - Value return (Item is small, copying is cheap)
func (fs *fileStore) Add(text string) Item {
	// TODO: Implement this function
	panic("unimplemented")
}

// Toggle marks an item as done/not done by ID.
//
// Go Concepts Demonstrated:
// - Multiple return values (Item, bool) for "found" pattern
// - Slice iteration with index (to modify in-place)
// - Early return on success
//
// Why not return error?
// The "not found" case isn't exceptional—it's expected user behavior.
// Returning (Item, bool) is more idiomatic for optional results.
// Compare to map lookups: value, ok := m[key]
func (fs *fileStore) Toggle(id int) (Item, bool) {
	// TODO: Implement this function
	panic("unimplemented")
}

// List returns all items, optionally filtering out completed ones.
//
// Go Concepts Demonstrated:
// - Slice filtering (build new slice)
// - Conditional logic with boolean parameter
// - Return slice (no defensive copying needed; slices share backing array)
func (fs *fileStore) List(onlyPending bool) []Item {
	// TODO: Implement this function
	panic("unimplemented")
}


