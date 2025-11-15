//go:build !solution
// +build !solution

package exercise

import (
	"fmt"
	"math"
	"strings"
)

// ============================================================================
// Exercise 1: Optimize Prime Number Finding
// ============================================================================

// FindPrimes finds all prime numbers up to n.
// TODO: This uses a naive O(n²) algorithm. Optimize it using the Sieve of Eratosthenes.
// The current implementation will be slow for large n and will show up in CPU profiles.
//
// Hint: The Sieve of Eratosthenes:
// 1. Create a boolean array of size n+1
// 2. Mark all multiples of each prime starting from 2
// 3. Collect all unmarked numbers as primes
// Time complexity: O(n log log n) vs current O(n²)
func FindPrimes(n int) []int {
	// TODO: Implement efficient sieve algorithm
	// Currently using naive algorithm (SLOW!)
	var primes []int
	for i := 2; i <= n; i++ {
		isPrime := true
		for j := 2; j < i; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, i)
		}
	}
	return primes
}

// ============================================================================
// Exercise 2: Reduce String Allocation
// ============================================================================

// BuildReport creates a report string from items.
// TODO: This uses string concatenation which creates many intermediate allocations.
// Optimize using strings.Builder to reduce memory allocations.
//
// Hint:
// 1. Create a strings.Builder
// 2. Use buf.Grow() to preallocate capacity (estimate: 100 bytes per item)
// 3. Use buf.WriteString() instead of +=
// 4. Return buf.String()
func BuildReport(items []Item) string {
	// TODO: Optimize to use strings.Builder
	// Current implementation allocates a new string on each concatenation
	report := "=== Report ===\n"
	for _, item := range items {
		report += fmt.Sprintf("ID: %d, Name: %s, Value: %.2f\n",
			item.ID, item.Name, item.Value)
	}
	report += "=== End Report ===\n"
	return report
}

// ============================================================================
// Exercise 3: Optimize Search Algorithm
// ============================================================================

// SearchDocuments searches for documents containing the query string.
// TODO: This uses a naive linear search with case-insensitive comparison.
// Optimize by:
// 1. Converting query to lowercase once (not in loop)
// 2. Precomputing lowercased content if possible
// 3. Using more efficient string search (strings.Contains vs loops)
//
// For even better performance, consider building an inverted index.
func SearchDocuments(docs []Document, query string) []Document {
	// TODO: Optimize this search
	var results []Document
	for _, doc := range docs {
		// Inefficient: converts to lowercase on every check
		if containsCaseInsensitive(doc.Title, query) ||
			containsCaseInsensitive(doc.Content, query) {
			results = append(results, doc)
		}
	}
	return results
}

// containsCaseInsensitive is intentionally inefficient for the exercise
func containsCaseInsensitive(text, substr string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}

// ============================================================================
// Exercise 4: Reduce Allocations in Loop
// ============================================================================

// ProcessItems processes a slice of items and returns results.
// TODO: This allocates many intermediate values. Optimize by:
// 1. Preallocating the results slice with correct capacity
// 2. Reusing the category string computation
// 3. Avoiding repeated calculations
func ProcessItems(items []Item) []Result {
	// TODO: Preallocate results slice
	var results []Result

	for _, item := range items {
		// TODO: This recalculates category every time - can we optimize?
		category := determineCategory(item.Value)

		result := Result{
			ItemID:    item.ID,
			Score:     calculateScore(item),
			Category:  category,
			Processed: item.Timestamp,
		}
		results = append(results, result)
	}

	return results
}

// determineCategory is a helper function (intentionally simple)
func determineCategory(value float64) string {
	switch {
	case value < 10:
		return "low"
	case value < 100:
		return "medium"
	default:
		return "high"
	}
}

// calculateScore computes a score from an item
func calculateScore(item Item) float64 {
	// Simple scoring logic
	score := item.Value * float64(len(item.Tags))
	if item.Name != "" {
		score *= 1.1
	}
	return score
}

// ============================================================================
// Exercise 5: Optimize JSON Processing
// ============================================================================

// FormatItemsAsJSON converts items to JSON-like strings.
// TODO: This builds JSON manually with string concatenation (very inefficient!).
// Optimize by:
// 1. Using encoding/json.Marshal for correctness
// 2. OR use strings.Builder if manual formatting is required
// 3. Preallocate buffer capacity
func FormatItemsAsJSON(items []Item) string {
	// TODO: Optimize this JSON formatting
	result := "[\n"
	for i, item := range items {
		result += "  {\n"
		result += fmt.Sprintf("    \"id\": %d,\n", item.ID)
		result += fmt.Sprintf("    \"name\": \"%s\",\n", item.Name)
		result += fmt.Sprintf("    \"value\": %.2f\n", item.Value)
		result += "  }"
		if i < len(items)-1 {
			result += ","
		}
		result += "\n"
	}
	result += "]\n"
	return result
}

// ============================================================================
// Exercise 6: Optimize Mathematical Computation
// ============================================================================

// ComputeDistances calculates Euclidean distances between all pairs of points.
// TODO: This has O(n²) complexity which is unavoidable, but you can optimize:
// 1. Reduce allocations by preallocating result slice
// 2. Avoid redundant calculations (distance(a,b) == distance(b,a))
// 3. Use integer math where possible
// 4. Consider parallel processing for large datasets
func ComputeDistances(points [][2]float64) []float64 {
	// TODO: Preallocate result slice with correct size
	var distances []float64

	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			// TODO: This calculation is expensive - profile shows it's a hotspot
			dx := points[i][0] - points[j][0]
			dy := points[i][1] - points[j][1]
			dist := math.Sqrt(dx*dx + dy*dy)
			distances = append(distances, dist)
		}
	}

	return distances
}

// ============================================================================
// Exercise 7: Optimize Map Operations
// ============================================================================

// CountWordFrequency counts word occurrences in documents.
// TODO: This has several inefficiencies:
// 1. Repeated string splitting and lowercase conversion
// 2. Map is not preallocated
// 3. Could use strings.Fields instead of strings.Split
func CountWordFrequency(docs []Document) map[string]int {
	// TODO: Preallocate map with estimated size
	wordCount := make(map[string]int)

	for _, doc := range docs {
		// TODO: Optimize string processing
		text := strings.ToLower(doc.Content)
		words := strings.Split(text, " ")

		for _, word := range words {
			// TODO: Trim and validate words
			word = strings.TrimSpace(word)
			if word != "" {
				wordCount[word]++
			}
		}
	}

	return wordCount
}

// ============================================================================
// Exercise 8: Implement Efficient Cache
// ============================================================================

// SimpleCache is a basic cache implementation.
// TODO: This cache is missing important optimizations:
// 1. No capacity limit (unbounded growth - memory leak!)
// 2. Uses sync.Mutex for all reads (could use sync.RWMutex)
// 3. No eviction policy (LRU, TTL, etc.)
// 4. Could use sync.Map for specific access patterns
type SimpleCache struct {
	// TODO: Add proper fields (mutex, map, capacity limit, etc.)
	data map[string]interface{}
}

// NewSimpleCache creates a new cache.
// TODO: Initialize with proper fields and capacity limit
func NewSimpleCache(capacity int) *SimpleCache {
	// TODO: Implement proper initialization
	return &SimpleCache{
		data: make(map[string]interface{}),
	}
}

// Get retrieves a value from cache.
// TODO: Implement with proper locking (RWMutex for reads)
func (c *SimpleCache) Get(key string) (interface{}, bool) {
	// TODO: Implement with RLock/RUnlock
	val, ok := c.data[key]
	return val, ok
}

// Set stores a value in cache.
// TODO: Implement with:
// 1. Proper locking (Lock/Unlock)
// 2. Capacity checking and eviction
// 3. LRU tracking if implementing LRU eviction
func (c *SimpleCache) Set(key string, value interface{}) {
	// TODO: Implement with capacity limit and eviction
	c.data[key] = value
}

// ============================================================================
// Exercise 9: Optimize Slice Operations
// ============================================================================

// FilterAndTransform filters items and transforms them to results.
// TODO: This has several allocation inefficiencies:
// 1. Doesn't preallocate result slice
// 2. Could avoid intermediate allocations
// 3. Could process in-place if original slice isn't needed
func FilterAndTransform(items []Item, minValue float64) []Result {
	// TODO: Preallocate with estimated capacity
	var results []Result

	for _, item := range items {
		if item.Value >= minValue {
			result := Result{
				ItemID:   item.ID,
				Score:    item.Value * 2,
				Category: "filtered",
			}
			results = append(results, result)
		}
	}

	return results
}

// ============================================================================
// Exercise 10: Optimize Recursive Algorithm
// ============================================================================

// Fibonacci computes the nth Fibonacci number.
// TODO: This recursive implementation is extremely inefficient for large n.
// Optimize using:
// 1. Iterative approach (O(n) time, O(1) space)
// 2. Memoization (O(n) time, O(n) space)
// 3. Matrix exponentiation (O(log n) time)
func Fibonacci(n int) int {
	// TODO: Replace with iterative or memoized version
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2) // Exponential time!
}

// ============================================================================
// Helper Functions (for testing)
// ============================================================================

// GenerateTestItems creates test data
func GenerateTestItems(count int) []Item {
	items := make([]Item, count)
	for i := 0; i < count; i++ {
		items[i] = Item{
			ID:        i,
			Name:      fmt.Sprintf("Item-%d", i),
			Value:     float64(i % 100),
			Timestamp: int64(i),
			Tags:      []string{"tag1", "tag2"},
			Metadata: map[string]string{
				"key": "value",
			},
		}
	}
	return items
}

// GenerateTestDocuments creates test documents
func GenerateTestDocuments(count int) []Document {
	docs := make([]Document, count)
	content := []string{
		"The quick brown fox jumps over the lazy dog",
		"Go is a statically typed compiled programming language",
		"Profiling helps identify performance bottlenecks",
		"Memory allocation can be optimized using sync.Pool",
		"CPU profiling shows hot paths in your code",
	}

	for i := 0; i < count; i++ {
		docs[i] = Document{
			ID:      i,
			Title:   fmt.Sprintf("Document %d", i),
			Content: content[i%len(content)],
			Author:  "Test Author",
			Tags:    []string{"test", "document"},
		}
	}
	return docs
}

// GenerateTestPoints creates test points for distance calculation
func GenerateTestPoints(count int) [][2]float64 {
	points := make([][2]float64, count)
	for i := 0; i < count; i++ {
		points[i] = [2]float64{float64(i), float64(i * 2)}
	}
	return points
}
