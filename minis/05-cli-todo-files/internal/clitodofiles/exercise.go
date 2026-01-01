//go:build !solution && !reference

package clitodofiles

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

// NewFileStore implements the exercise.
//
// TODO: Implement this function
func NewFileStore(path string) Store {
	// TODO: Implement
	return Store{}
}

// Load implements the exercise.
//
// TODO: Implement this function
func (fs *fileStore) Load() error {
	// TODO: Implement
	return nil
}

// Save implements the exercise.
//
// TODO: Implement this function
func (fs *fileStore) Save() error {
	// TODO: Implement
	return nil
}

// Add implements the exercise.
//
// TODO: Implement this function
func (fs *fileStore) Add(text string) Item {
	// TODO: Implement
	return Item{}
}

// Toggle implements the exercise.
//
// TODO: Implement this function
func (fs *fileStore) Toggle(id int) (Item, bool) {
	// TODO: Implement
	return Item{}, false
}

// List implements the exercise.
//
// TODO: Implement this function
func (fs *fileStore) List(onlyPending bool) []Item {
	// TODO: Implement
	return nil
}
