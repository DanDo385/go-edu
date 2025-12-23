//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "encoding/json" for JSON marshaling/unmarshaling
// - "fmt" for error formatting
// - "os" for file I/O operations
//
import (
    "encoding/json"
    "fmt"
    "os"
)

// ============================================================================
// CLI TODO APPLICATION: Building a Persistent Task Manager
// ============================================================================
//
// This project demonstrates:
// 1. JSON persistence (serialization/deserialization)
// 2. File I/O operations (reading/writing)
// 3. Interface-based design (abstraction)
// 4. CLI flag parsing (command-line interface)
// 5. CRUD operations (Create, Read, Update via Toggle, List)
//
// Key Concepts:
// - Persistence: Data survives program restarts by saving to disk
// - JSON Format: Human-readable structured data format
// - Interface: Contract defining behavior without implementation details
// - File Atomicity: Preventing partial writes during crashes
//
// Memory Management:
// - Slices share underlying arrays (be careful with modifications)
// - JSON marshal/unmarshal allocates new memory
// - File I/O reads entire file into memory (not streaming)
// - Maps use pointer semantics internally (modifications visible)
//
// ============================================================================

// ============================================================================
// Exercise 1: Define the Data Model
// ============================================================================

// Item represents a single TODO item.
// TODO: Define the Item struct with proper JSON tags.
type Item struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`		
}

// Memory considerations:
// - struct is a value type (copied when assigned)
// - String fields internally contain pointer to data (8 bytes) + len/cap
// - int field is typically 8 bytes on 64-bit systems
// - bool field is 1 byte (but padded for alignment)
// - Total size: ~32 bytes per item (with padding)
//
// JSON struct tags format: `json:"fieldname"`
// - Maps Go field names to JSON field names
// - Enables lowercase JSON keys (JSON convention)
// - Controls marshaling/unmarshaling behavior
//
// type Item struct {
//     ID   int    `json:"id"`    // Unique identifier
//     Text string `json:"text"`  // Task description
//     Done bool   `json:"done"`  // Completion status
// }

// ============================================================================
// Exercise 2: Define the Store Interface
// ============================================================================

// Store defines operations for managing TODO items.
// TODO: Define the Store interface with all required methods.
//
// Why use an interface?
// - Abstraction: Hide implementation details from users
// - Testability: Can create mock implementations for testing
// - Flexibility: Swap file storage for database without changing callers
// - Documentation: Interface = contract of what a Store must do
//
// Methods to implement:
// - Load() error: Read items from persistent storage
// - Save() error: Write items to persistent storage
// - Add(text string) Item: Create a new item and return it
// - Toggle(id int) (Item, bool): Mark item as done/undone, return (item, found)
// - List(onlyPending bool) []Item: Return items (optionally filter completed)
//
type Store interface {
    Load() error
    Save() error
    Add(text string) Item
    Toggle(id int) (Item, bool)
    List(onlyPending bool) []Item
}

// ============================================================================
// Exercise 3: Implement the FileStore Type
// ============================================================================

// fileStore persists items to a JSON file.
// TODO: Define the fileStore struct with necessary fields.
//
// Memory considerations:
// - path is a string (pointer + length internally, ~16 bytes)
// - items is a slice (pointer + len + cap, 24 bytes header)
// - The slice header points to an array on the heap
// - Each Item in the array is ~32 bytes
// - Total: 40 bytes struct + (32 * n) bytes for items
//
// Design pattern:
// - Keep data in memory (items slice)
// - Only touch disk on Load() and Save()
// - This is fast for small datasets (<10k items)
// - For large datasets, consider streaming or database
//
type fileStore struct {
    path  string  // File path for JSON storage
    items []Item  // In-memory cache of all items
}

// ============================================================================
// Exercise 4: Constructor Function
// ============================================================================

// NewFileStore creates a Store backed by a JSON file at path.
// TODO: Implement the constructor function.
//
// Memory allocation:
// - Creates a fileStore on the heap (return pointer)
// - Initializes empty slice (no underlying array yet)
// - Returns interface type (hides implementation)
//
// Why return *fileStore as Store interface?
// - Caller doesn't know it's a fileStore (abstraction)
// - Easy to swap implementations later
// - Interface = behavior, not structure
//
func NewFileStore(path string) Store {
    return &fileStore{
        path:  path,
        items: []Item{},  // Empty slice (not nil)
    }
}

// ============================================================================
// Exercise 5: Load from JSON
// ============================================================================

// Load reads items from disk.
// TODO: Implement thread-safe file reading and JSON unmarshaling.
//
// Steps:
// 1. Read entire file as []byte using os.ReadFile()
// 2. Unmarshal JSON array into fs.items slice
// 3. Return error if file doesn't exist or JSON is invalid
//
// Memory operations:
// - os.ReadFile allocates a new []byte (size = file size)
// - json.Unmarshal allocates new Item structs and slice
// - Old fs.items is garbage collected if no other references
// - For 1000 items: ~32KB for items + file size in bytes
//
// Error handling:
// - File doesn't exist: Return error (caller decides if that's OK)
// - Malformed JSON: Wrap error with context using fmt.Errorf("%w")
// - Empty file: Unmarshal into empty slice (valid JSON: [])
//
func (fs *fileStore) Load() error {
    // Read entire file into memory
    data, err := os.ReadFile(fs.path)
    if err != nil {
        return err  // Caller can check os.IsNotExist(err)
    }

    // Parse JSON array into slice
    if err := json.Unmarshal(data, &fs.items); err != nil {
        return fmt.Errorf("parsing JSON: %w", err)
    }

    return nil
}

// ============================================================================
// Exercise 6: Save to JSON
// ============================================================================

// Save writes items to disk.
// TODO: Implement JSON marshaling and file writing.
//
// Steps:
// 1. Marshal fs.items to JSON with indentation (pretty-print)
// 2. Write JSON bytes to file with proper permissions
// 3. Return error if marshal or write fails
//
// Memory operations:
// - json.MarshalIndent allocates new []byte (size proportional to data)
// - For 1000 items: ~50KB JSON (more than binary due to formatting)
// - os.WriteFile writes all at once (not streaming)
//
// File permissions (0644):
// - Owner: read (4) + write (2) = 6
// - Group: read (4)
// - Others: read (4)
// - Standard for data files (not executables)
//
// Atomicity consideration:
// - os.WriteFile is not atomic (crash mid-write = corrupted file)
// - Production: Write to temp file, then rename (atomic on POSIX)
// - For this app: Acceptable risk
//
func (fs *fileStore) Save() error {
    // Marshal with indentation for readability
    data, err := json.MarshalIndent(fs.items, "", "  ")
    if err != nil {
        return fmt.Errorf("marshaling JSON: %w", err)
    }

    // Write to file
    if err := os.WriteFile(fs.path, data, 0644); err != nil {
        return fmt.Errorf("writing file: %w", err)
    }

    return nil
}

// ============================================================================
// Exercise 7: Add New Items
// ============================================================================

// Add inserts a new item with an auto-incremented ID.
// TODO: Implement ID generation and item creation.
//
// ID generation strategy:
// - Find maximum existing ID in slice
// - New ID = max + 1
// - Ensures uniqueness within this file
//
// Memory operations:
// - Loop through all items: O(n) time
// - Create new Item struct: 32 bytes allocation
// - append() may grow underlying array:
//   * If len < cap: Just increment length (O(1))
//   * If len == cap: Allocate new array (2x capacity), copy all items (O(n))
//   * Amortized O(1) for append operations
//
// Limitations:
// - Not thread-safe (multiple goroutines = race condition)
// - IDs not reused after deletion
// - Max ID grows forever (overflow after ~2 billion items)
//
// Alternatives:
// - UUID: Globally unique, works in distributed systems
// - Database auto-increment: Let DB handle IDs
// - Timestamp + random: Time-sortable unique IDs
//
func (fs *fileStore) Add(text string) Item {
// Find max ID
    maxID := 0
    for _, item := range fs.items {
        if item.ID > maxID {
            maxID = item.ID
        }
    }
//
// Create item with next ID
    newItem := Item{
        ID:   maxID + 1,
        Text: text,
        Done: false,
    }
//
//     // Append to slice
    fs.items = append(fs.items, newItem)

    return newItem
// }

// ============================================================================
// Exercise 8: Toggle Item Completion
// ============================================================================

// Toggle flips Done for the given ID.
// TODO: Implement find-and-toggle logic.
//
// Returns: (updated item, true) if found, (zero value, false) if not found
//
// Why (Item, bool) instead of error?
// - "Not found" is expected behavior (not exceptional)
// - Matches Go idiom: value, ok := map[key]
// - Caller can check ok without error handling ceremony
//
// Memory operations:
// - Linear search through slice: O(n) time
// - Modify item in-place (no allocation)
// - Return copy of Item (value type, 32 bytes copied)
//
// Performance consideration:
// - For large lists (>10k items), use map[int]*Item for O(1) lookup
// - Trade-off: More memory (map overhead) vs. faster lookups
//
func (fs *fileStore) Toggle(id int) (Item, bool) {
	// Linear search for the item
	// For large lists, consider a map[int]*Item for O(1) lookup
	for i := range fs.items {
		if fs.items[i].ID == id {
			// Toggle the done status
			fs.items[i].Done = !fs.items[i].Done
			return fs.items[i], true
		}
	}

	// Not found: return zero value and false
	return Item{}, false
}

// ============================================================================
// Exercise 9: List Items
// ============================================================================

// List returns all items, optionally filtering out completed ones.
// TODO: Implement filtering logic.
//
// Parameters:
// - onlyPending: If true, exclude items where Done == true
//
// Memory operations:
// - If onlyPending == false: Return slice as-is (shares backing array)
// - If onlyPending == true: Build new slice (allocates new array)
// - Filtered slice: O(n) to iterate, O(k) space where k = pending count
//
// Slice sharing caution:
// - Returned slice shares backing array with fs.items
// - Caller should NOT modify returned items (will affect stored data)
// - For defensive copy: return append([]Item(nil), fs.items...)
// - Trade-off: Performance vs. safety
//
// Alternative: Return []Item copy always
// - Pros: Caller can safely modify
// - Cons: Extra allocation and copying (32 * n bytes)
//
func (fs *fileStore) List(onlyPending bool) []Item {
    // Return all if not filtering
    if !onlyPending {
        return fs.items
    }

    // Filter pending items
    var filtered []Item
    for _, item := range fs.items {
        if !item.Done {
            filtered = append(filtered, item)
        }
    }

    return filtered
}

// ============================================================================
// Testing Strategy
// ============================================================================
//
// 1. Unit tests for each method:
//    - Add: Verify ID generation, check item in list
//    - Toggle: Test found/not found, verify Done flips
//    - List: Test all vs. pending filtering
//
// 2. Integration tests:
//    - Load → Add → Save → Load again (persistence)
//    - Test with non-existent file (Load should fail gracefully)
//    - Test with malformed JSON (Load should return error)
//
// 3. Use t.TempDir() for test files:
//    - Automatically cleaned up after test
//    - Isolated from other tests (no shared state)
//
// 4. Table-driven tests for edge cases:
//    - Empty list, single item, many items
//    - Toggle non-existent ID
//    - List with all done, all pending, mixed
//
// ============================================================================
// Performance Characteristics
// ============================================================================
//
// Time Complexity:
// - Load: O(n) to unmarshal JSON
// - Save: O(n) to marshal JSON
// - Add: O(n) to find max ID + O(1) amortized for append
// - Toggle: O(n) to search by ID
// - List: O(n) to filter items
//
// Space Complexity:
// - O(n) for items slice
// - O(n) for JSON bytes during Load/Save
// - O(k) for filtered list where k = pending count
//
// Optimization opportunities:
// 1. Index by ID: map[int]*Item for O(1) Toggle
// 2. Track max ID: Avoid O(n) scan in Add
// 3. Lazy load: Only Load on first access (if file is large)
// 4. Write-through cache: Save only changed items (requires diffing)
//
// ============================================================================
// Production Considerations
// ============================================================================
//
// 1. Atomic writes:
//    - Current: os.WriteFile can partial-write on crash
//    - Better: Write to temp file, then os.Rename (atomic on POSIX)
//
// 2. Concurrent access:
//    - Current: Not thread-safe (race conditions)
//    - Fix: Use sync.Mutex or sync.RWMutex
//
// 3. File locking:
//    - Current: Multiple processes can corrupt file
//    - Fix: Use flock() or similar OS-level file locking
//
// 4. Data corruption:
//    - Current: No validation after Load
//    - Fix: Validate IDs are unique, no negative IDs, etc.
//
// 5. Large datasets:
//    - Current: Load entire file into memory
//    - Fix: Use database or streaming JSON parsing
//
// 6. Backup/recovery:
//    - Current: No backups
//    - Fix: Keep versioned backups, implement undo
//
// ============================================================================
// After implementing all functions:
// - Run: go test -v ./...
// - Build CLI: cd ../cmd/cli-todo-files && go build
// - Test manually:
//   $ ./cli-todo-files -add "Learn Go"
//   $ ./cli-todo-files -list
//   $ ./cli-todo-files -toggle 1
//   $ ./cli-todo-files -list -all
// - Inspect JSON: cat todos.json
// - Compare with solution: diff exercise.go solution.go
// ============================================================================
