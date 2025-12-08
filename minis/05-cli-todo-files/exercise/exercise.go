//go:build !solution
// +build !solution

package exercise

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

// fileStore persists items to a JSON file.
type fileStore struct {
	path  string
	items []Item
}

// NewFileStore creates a Store backed by a JSON file at path.
func NewFileStore(path string) Store {
	return &fileStore{
		path:  path,
		items: []Item{},
	}
}

// Load reads items from disk.
func (fs *fileStore) Load() error {
	data, err := os.ReadFile(fs.path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &fs.items); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}
	return nil
}

// Save writes items to disk.
func (fs *fileStore) Save() error {
	data, err := json.MarshalIndent(fs.items, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(fs.path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// Add inserts a new item with an auto-incremented ID.
func (fs *fileStore) Add(text string) Item {
	maxID := 0
	for _, item := range fs.items {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	newItem := Item{ID: maxID + 1, Text: text}
	fs.items = append(fs.items, newItem)
	return newItem
}

// Toggle flips Done for the given ID.
func (fs *fileStore) Toggle(id int) (Item, bool) {
	for i := range fs.items {
		if fs.items[i].ID == id {
			fs.items[i].Done = !fs.items[i].Done
			return fs.items[i], true
		}
	}
	return Item{}, false
}

// List returns all items, optionally filtering out completed ones.
func (fs *fileStore) List(onlyPending bool) []Item {
	if !onlyPending {
		return append([]Item(nil), fs.items...)
	}
	var filtered []Item
	for _, item := range fs.items {
		if !item.Done {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
