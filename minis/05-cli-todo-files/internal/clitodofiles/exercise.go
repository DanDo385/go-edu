//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package clitodofiles

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
// TODO: implement NewFileStore.
func NewFileStore(path string) Store { panic("TODO: implement") }
// TODO: implement Load.
func (fs *fileStore) Load() error { panic("TODO: implement") }
// TODO: implement Save.
func (fs *fileStore) Save() error { panic("TODO: implement") }
// TODO: implement Add.
func (fs *fileStore) Add(text string) Item { panic("TODO: implement") }
// TODO: implement Toggle.
func (fs *fileStore) Toggle(id int) (Item, bool) { panic("TODO: implement") }
// TODO: implement List.
func (fs *fileStore) List(onlyPending bool) []Item { panic("TODO: implement") }
