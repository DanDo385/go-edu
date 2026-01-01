//go:build !solution && !reference

package clitodofiles

import (
	"encoding/json"
	"fmt"
	"os"
)

func NewFileStore(path string) Store {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func (fs *fileStore) Load() error {
	// TODO: Implement this function
	panic("not implemented")
}

func (fs *fileStore) Save() error {
	// TODO: Implement this function
	panic("not implemented")
}

func (fs *fileStore) Add(text string) Item {
	// TODO: Implement this function
	panic("not implemented")
}

func (fs *fileStore) Toggle(id int) (Item, bool) {
	// TODO: Implement this function
	panic("not implemented")
}

func (fs *fileStore) List(onlyPending bool) []Item {
	// TODO: Implement this function
	panic("not implemented")
}
