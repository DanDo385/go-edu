package exercise

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStore_AddAndList(t *testing.T) {
	// Create a temporary directory for test files
	// t.TempDir() automatically cleans up after the test
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test.json")

	store := NewFileStore(storePath)

	// Add items
	item1 := store.Add("Learn Go")
	item2 := store.Add("Write tests")

	if item1.ID == 0 {
		t.Error("Item1 ID should not be 0")
	}
	if item2.ID <= item1.ID {
		t.Error("Item2 ID should be greater than Item1 ID")
	}

	// List all items
	items := store.List(false)
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestFileStore_Toggle(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test.json")

	store := NewFileStore(storePath)
	item := store.Add("Test task")

	// Toggle to done
	toggled, found := store.Toggle(item.ID)
	if !found {
		t.Fatal("Item not found after adding")
	}
	if !toggled.Done {
		t.Error("Item should be marked as done")
	}

	// Toggle back to not done
	toggled, found = store.Toggle(item.ID)
	if !found {
		t.Fatal("Item not found on second toggle")
	}
	if toggled.Done {
		t.Error("Item should be marked as not done")
	}

	// Try to toggle non-existent item
	_, found = store.Toggle(9999)
	if found {
		t.Error("Should not find non-existent item")
	}
}

func TestFileStore_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test.json")

	// Create store and add items
	store1 := NewFileStore(storePath)
	store1.Add("Item 1")
	store1.Add("Item 2")
	store1.Toggle(1) // Mark first item as done

	// Save to file
	if err := store1.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Create new store and load
	store2 := NewFileStore(storePath)
	if err := store2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify items
	items := store2.List(false)
	if len(items) != 2 {
		t.Errorf("Expected 2 items after load, got %d", len(items))
	}

	if !items[0].Done {
		t.Error("First item should be marked as done")
	}
	if items[1].Done {
		t.Error("Second item should not be marked as done")
	}
}

func TestFileStore_ListFiltering(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test.json")

	store := NewFileStore(storePath)
	store.Add("Task 1")
	store.Add("Task 2")
	store.Add("Task 3")
	store.Toggle(2) // Mark second item as done

	// List only pending
	pending := store.List(true)
	if len(pending) != 2 {
		t.Errorf("Expected 2 pending items, got %d", len(pending))
	}

	// List all
	all := store.List(false)
	if len(all) != 3 {
		t.Errorf("Expected 3 total items, got %d", len(all))
	}
}

func TestFileStore_LoadNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "nonexistent.json")

	store := NewFileStore(storePath)
	err := store.Load()

	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}

	if !os.IsNotExist(err) {
		t.Errorf("Expected IsNotExist error, got: %v", err)
	}
}

func TestFileStore_LoadMalformedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "malformed.json")

	// Write malformed JSON
	os.WriteFile(storePath, []byte("{not valid json"), 0644)

	store := NewFileStore(storePath)
	err := store.Load()

	if err == nil {
		t.Error("Expected error when loading malformed JSON")
	}
}
