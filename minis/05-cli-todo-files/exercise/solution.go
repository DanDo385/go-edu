/*
Problem: Build a persistent TODO list with JSON file storage

We need to implement:
1. CRUD operations (Create, Read, Update, Delete-ish with Toggle)
2. JSON serialization/deserialization for persistence
3. CLI interface with flag parsing
4. Atomic file writes (no partial corruption)

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

Why Go is well-suited:
- `flag` package for CLI parsing (built-in, type-safe)
- JSON marshal/unmarshal with struct tags (no external dependencies)
- Interfaces enable testing without real files
- Pointer receivers for mutable state (clear semantics)

Compared to other languages:
- Python: `argparse` + `json.dump()` similar, but dynamic typing
- Node.js: Requires external CLI library; async file I/O more complex
- Rust: More boilerplate for serialization; ownership rules are steeper
*/

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
	return &fileStore{
		path:  path,
		items: []Item{}, // Initialize empty slice
	}
}

// Load reads items from the JSON file.
//
// Go Concepts Demonstrated:
// - os.ReadFile: Read entire file into memory (simple but not streaming)
// - json.Unmarshal: Deserialize JSON into Go structs
// - Error handling: Return error for caller to handle
// - Pointer receiver: Method can access/modify struct fields
//
// Three-Input Iteration Table:
//
// Input 1: File exists with valid JSON (happy path)
//   os.ReadFile → []byte of JSON array
//   json.Unmarshal → populates fs.items
//   Result: nil error
//
// Input 2: File doesn't exist (edge case)
//   os.ReadFile → error (os.IsNotExist)
//   Result: return error (caller should handle gracefully)
//
// Input 3: File exists but malformed JSON (failure)
//   os.ReadFile → []byte
//   json.Unmarshal → error
//   Result: return error with context
func (fs *fileStore) Load() error {
	// Read entire file into memory
	// For very large files (>100MB), consider streaming with json.Decoder
	// But for a TODO app, this is fine
	data, err := os.ReadFile(fs.path)
	if err != nil {
		// Return error as-is (caller can check os.IsNotExist)
		return err
	}

	// Unmarshal JSON array into slice
	// The json package uses reflection to map JSON fields to struct tags
	// Tags like `json:"id"` specify the JSON field name
	if err := json.Unmarshal(data, &fs.items); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	return nil
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
	// Marshal to JSON with indentation for human readability
	// MarshalIndent(value, prefix, indent)
	// - prefix: string to prepend to each line (usually "")
	// - indent: string for each indentation level (usually "  " or "\t")
	data, err := json.MarshalIndent(fs.items, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	// Write to file atomically
	// Permissions 0644 = owner read/write, group read, others read
	// For production, consider writing to temp file then renaming:
	//   os.WriteFile(tmpPath, data, 0644)
	//   os.Rename(tmpPath, fs.path)
	// This prevents partial writes if process is killed
	if err := os.WriteFile(fs.path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// Add creates a new item and returns it.
//
// Go Concepts Demonstrated:
// - Slice append (grows capacity automatically)
// - ID generation (find max + 1)
// - Value return (Item is small, copying is cheap)
func (fs *fileStore) Add(text string) Item {
	// Generate new ID by finding max existing ID
	// This ensures uniqueness but has limitations:
	// - Not thread-safe (race conditions if multiple processes)
	// - IDs are not reused after deletion
	// For production: use database auto-increment or UUIDs
	maxID := 0
	for _, item := range fs.items {
		if item.ID > maxID {
			maxID = item.ID
		}
	}

	// Create new item
	newItem := Item{
		ID:   maxID + 1,
		Text: text,
		Done: false,
	}

	// Append to slice
	// append() may allocate a new backing array if capacity is exceeded
	// This is O(1) amortized (occasional O(n) when reallocation occurs)
	fs.items = append(fs.items, newItem)

	return newItem
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

// List returns all items, optionally filtering out completed ones.
//
// Go Concepts Demonstrated:
// - Slice filtering (build new slice)
// - Conditional logic with boolean parameter
// - Return slice (no defensive copying needed; slices share backing array)
func (fs *fileStore) List(onlyPending bool) []Item {
	// If showing all items, return directly
	if !onlyPending {
		return fs.items
	}

	// Filter out completed items
	var result []Item
	for _, item := range fs.items {
		if !item.Done {
			result = append(result, item)
		}
	}

	return result
}

/*
Alternatives & Trade-offs:

1. Use map instead of slice:
   items map[int]Item
   Pros: O(1) lookup by ID
   Cons: No ordering; JSON marshal requires converting to slice anyway

2. Pointer items in slice:
   items []*Item
   Pros: Modify in-place without index; less copying
   Cons: More allocations; nil pointer checks

3. Atomic file writes:
   tmpFile, _ := os.CreateTemp(filepath.Dir(fs.path), "todos-*.tmp")
   tmpFile.Write(data)
   tmpFile.Close()
   os.Rename(tmpFile.Name(), fs.path)
   Pros: Prevents corruption if process is killed mid-write
   Cons: More code; overkill for simple TODO app

4. Use database (SQLite):
   Pros: ACID guarantees; SQL queries; scales better
   Cons: Requires CGO (or pure Go SQL library); more complexity

5. Streaming JSON for large datasets:
   enc := json.NewEncoder(file)
   for _, item := range fs.items {
       enc.Encode(item)  // One JSON object per line (JSONL)
   }
   Pros: Constant memory usage
   Cons: More complex; loses pretty-printing

Go vs X:

Go vs Python:
  import json
  class FileStore:
      def __init__(self, path):
          self.path = path
          self.items = []
      def load(self):
          with open(self.path) as f:
              self.items = json.load(f)
      def save(self):
          with open(self.path, 'w') as f:
              json.dump(self.items, f, indent=2)
  Pros: Less code; dynamic typing is flexible
  Cons: No compile-time safety (typos in field names are runtime errors)
        No interface/protocol checking
  Go: Type safety catches errors early; clearer structure

Go vs Node.js:
  const fs = require('fs').promises;
  class FileStore {
      async load() {
          const data = await fs.readFile(this.path, 'utf8');
          this.items = JSON.parse(data);
      }
      async save() {
          await fs.writeFile(this.path, JSON.stringify(this.items, null, 2));
      }
  }
  Pros: Similar brevity; async/await is clean
  Cons: Async complexity (Promises, error handling)
        Requires external CLI library (commander, yargs)
  Go: Synchronous I/O is simpler; built-in flag package

Go vs Rust:
  use serde::{Deserialize, Serialize};
  #[derive(Serialize, Deserialize)]
  struct Item { id: u32, text: String, done: bool }
  impl FileStore {
      fn load(&mut self) -> Result<(), Error> {
          let data = std::fs::read_to_string(&self.path)?;
          self.items = serde_json::from_str(&data)?;
          Ok(())
      }
      fn save(&self) -> Result<(), Error> {
          let data = serde_json::to_string_pretty(&self.items)?;
          std::fs::write(&self.path, data)?;
          Ok(())
      }
  }
  Pros: Zero-cost abstractions; compile-time guarantees
  Cons: Ownership/borrowing complexity (lifetimes, &mut self)
        More boilerplate (derive macros, Result types)
  Go: Simpler mental model; faster iteration

Go vs Java:
  class FileStore {
      private List<Item> items = new ArrayList<>();
      public void load() throws IOException {
          String json = Files.readString(Path.of(path));
          items = objectMapper.readValue(json, new TypeReference<List<Item>>() {});
      }
      public void save() throws IOException {
          String json = objectMapper.writerWithDefaultPrettyPrinter().writeValueAsString(items);
          Files.writeString(Path.of(path), json);
      }
  }
  Pros: Similar structure; Jackson is robust
  Cons: Much more verbose (getters/setters, checked exceptions, generics)
        Requires external library (Jackson/Gson)
  Go: Cleaner syntax; built-in JSON support
*/
