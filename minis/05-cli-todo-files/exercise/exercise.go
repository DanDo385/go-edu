//go:build !solution
// +build !solution

package exercise

// Item represents a single TODO item.
type Item struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Store defines operations for managing TODO items.
type Store interface {
	// Load reads items from persistent storage
	Load() error

	// Save writes items to persistent storage
	Save() error

	// Add creates a new item and returns it
	Add(text string) Item

	// Toggle marks an item as done/not done by ID
	// Returns the updated item and true if found, or zero Item and false if not found
	Toggle(id int) (Item, bool)

	// List returns all items
	// If onlyPending is true, excludes completed items
	List(onlyPending bool) []Item
}

// NewFileStore creates a Store backed by a JSON file at path.
func NewFileStore(path string) Store {
	// TODO: implement
	return nil
}
