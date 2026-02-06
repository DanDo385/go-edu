//go:build !solution && !reference

package clitodofiles

/*
Problem: Build a persistent TODO list with JSON file storage
Constraints:
- Items have unique IDs (auto-incrementing)
- JSON file stores all items as an array
- Toggle operation is idempotent
- List can filter by completion status
Time/Space Complexity:
- Load/Save: O(n) where n = number of items (JSON marshal/unmarshal)
- Add: O(n) to find max ID, O(1) to append
- Toggle: O(n) to find item by ID
- List: O(n) to filter items
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

// NewFileStore - TODO: implement this function
func NewFileStore(path string) Store {
	return &fileStore{path: path, items: make([]Item, 0)}
}

// Load - TODO: implement this function
func (fs *fileStore) Load() error {
	data, err := os.ReadFile(fs.path)
	if err != nil {
		return err
	}
	var items []Item
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal todos: %w", err)
	}
	fs.items = items
	return nil
}

// Save - TODO: implement this function
func (fs *fileStore) Save() error {
	data, err := json.MarshalIndent(fs.items, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal todos: %w", err)
	}
	if err := os.WriteFile(fs.path, data, 0644); err != nil {
		return fmt.Errorf("write todos: %w", err)
	}
	return nil
}

// Add - TODO: implement this function
func (fs *fileStore) Add(text string) Item {
	maxID := 0
	for _, item := range fs.items {
		if item.ID > maxID {
			maxID = item.ID
		}
	}

	newItem := Item{
		ID:   maxID + 1,
		Text: text,
		Done: false,
	}
	fs.items = append(fs.items, newItem)
	return newItem
}

// Toggle - TODO: implement this function
func (fs *fileStore) Toggle(id int) (Item, bool) {
	for i := range fs.items {
		if fs.items[i].ID == id {
			fs.items[i].Done = !fs.items[i].Done
			return fs.items[i], true
		}
	}
	return Item{}, false
}

// List - TODO: implement this function
func (fs *fileStore) List(onlyPending bool) []Item {
	if !onlyPending {
		out := make([]Item, len(fs.items))
		copy(out, fs.items)
		return out
	}

	out := make([]Item, 0, len(fs.items))
	for _, item := range fs.items {
		if !item.Done {
			out = append(out, item)
		}
	}
	return out
}
