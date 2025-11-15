package exercise

// Item represents a data item to be processed
type Item struct {
	ID        int
	Name      string
	Value     float64
	Timestamp int64
	Tags      []string
	Metadata  map[string]string
}

// Result represents the output of processing an item
type Result struct {
	ItemID    int
	Score     float64
	Category  string
	Processed int64
}

// Document represents a text document for search/indexing
type Document struct {
	ID      int
	Title   string
	Content string
	Author  string
	Tags    []string
}

// CacheEntry represents a cached value with metadata
type CacheEntry struct {
	Key       string
	Value     interface{}
	CreatedAt int64
	Hits      int
}
