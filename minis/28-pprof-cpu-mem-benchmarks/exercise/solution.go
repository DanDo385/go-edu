package exercise

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
)

// ============================================================================
// Solution 1: Optimized Prime Number Finding
// ============================================================================

// FindPrimesOptimized uses the Sieve of Eratosthenes algorithm.
// Time complexity: O(n log log n) vs O(n²) for naive approach
// Space complexity: O(n) for the boolean array
func FindPrimesOptimized(n int) []int {
	if n < 2 {
		return nil
	}

	// Create boolean array - false means prime
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	// Sieve: mark all multiples as composite
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			// Mark all multiples of i as composite
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	// Collect primes
	primes := make([]int, 0, n/10) // Approximate prime density
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

// ============================================================================
// Solution 2: Optimized String Building
// ============================================================================

// BuildReportOptimized uses strings.Builder to avoid allocations.
// Reduces allocations from O(n) to O(1) by using a single growing buffer.
func BuildReportOptimized(items []Item) string {
	var buf strings.Builder

	// Preallocate capacity (estimate: 100 bytes per item + header/footer)
	buf.Grow(len(items)*100 + 100)

	buf.WriteString("=== Report ===\n")

	for _, item := range items {
		fmt.Fprintf(&buf, "ID: %d, Name: %s, Value: %.2f\n",
			item.ID, item.Name, item.Value)
	}

	buf.WriteString("=== End Report ===\n")

	return buf.String()
}

// ============================================================================
// Solution 3: Optimized Document Search
// ============================================================================

// SearchDocumentsOptimized improves search performance.
// Optimizations:
// 1. Convert query to lowercase once
// 2. Use strings.Contains directly (optimized implementation)
func SearchDocumentsOptimized(docs []Document, query string) []Document {
	queryLower := strings.ToLower(query)

	// Preallocate with estimated size
	results := make([]Document, 0, len(docs)/10)

	for _, doc := range docs {
		titleLower := strings.ToLower(doc.Title)
		contentLower := strings.ToLower(doc.Content)

		if strings.Contains(titleLower, queryLower) ||
			strings.Contains(contentLower, queryLower) {
			results = append(results, doc)
		}
	}

	return results
}

// SearchDocumentsWithIndex uses an inverted index for even better performance.
// This is more complex but much faster for multiple searches.
type DocumentIndex struct {
	docs  []Document
	index map[string][]int // word -> document IDs
}

// BuildIndex creates an inverted index
func BuildIndex(docs []Document) *DocumentIndex {
	idx := &DocumentIndex{
		docs:  docs,
		index: make(map[string][]int),
	}

	for _, doc := range docs {
		words := extractWords(doc.Title + " " + doc.Content)
		for _, word := range words {
			idx.index[word] = append(idx.index[word], doc.ID)
		}
	}

	return idx
}

// Search uses the inverted index
func (idx *DocumentIndex) Search(query string) []Document {
	words := extractWords(query)
	if len(words) == 0 {
		return nil
	}

	// Find documents containing all words
	docIDs := make(map[int]int)
	for _, word := range words {
		for _, docID := range idx.index[word] {
			docIDs[docID]++
		}
	}

	// Collect documents that match
	var results []Document
	for docID, count := range docIDs {
		if count == len(words) { // All words present
			results = append(results, idx.docs[docID])
		}
	}

	return results
}

func extractWords(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text) // More efficient than Split
	return words
}

// ============================================================================
// Solution 4: Optimized Item Processing
// ============================================================================

// ProcessItemsOptimized preallocates and reduces redundant calculations.
func ProcessItemsOptimized(items []Item) []Result {
	// Preallocate result slice with exact capacity
	results := make([]Result, 0, len(items))

	for _, item := range items {
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

// ============================================================================
// Solution 5: Optimized JSON Formatting
// ============================================================================

// FormatItemsAsJSONOptimized uses encoding/json for correctness and efficiency.
func FormatItemsAsJSONOptimized(items []Item) string {
	// Create simplified representation for JSON
	type SimpleItem struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}

	simpleItems := make([]SimpleItem, len(items))
	for i, item := range items {
		simpleItems[i] = SimpleItem{
			ID:    item.ID,
			Name:  item.Name,
			Value: item.Value,
		}
	}

	bytes, err := json.MarshalIndent(simpleItems, "", "  ")
	if err != nil {
		return ""
	}

	return string(bytes)
}

// FormatItemsAsJSONManual shows manual optimization with strings.Builder
func FormatItemsAsJSONManual(items []Item) string {
	var buf strings.Builder
	buf.Grow(len(items) * 100) // Preallocate

	buf.WriteString("[\n")
	for i, item := range items {
		buf.WriteString("  {\n")
		fmt.Fprintf(&buf, "    \"id\": %d,\n", item.ID)
		fmt.Fprintf(&buf, "    \"name\": \"%s\",\n", item.Name)
		fmt.Fprintf(&buf, "    \"value\": %.2f\n", item.Value)
		buf.WriteString("  }")
		if i < len(items)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}
	buf.WriteString("]\n")

	return buf.String()
}

// ============================================================================
// Solution 6: Optimized Distance Calculation
// ============================================================================

// ComputeDistancesOptimized preallocates and could be parallelized.
func ComputeDistancesOptimized(points [][2]float64) []float64 {
	n := len(points)
	size := n * (n - 1) / 2 // Number of unique pairs

	// Preallocate result slice
	distances := make([]float64, 0, size)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := points[i][0] - points[j][0]
			dy := points[i][1] - points[j][1]
			dist := math.Sqrt(dx*dx + dy*dy)
			distances = append(distances, dist)
		}
	}

	return distances
}

// ComputeDistancesParallel uses goroutines for large datasets.
func ComputeDistancesParallel(points [][2]float64) []float64 {
	n := len(points)
	size := n * (n - 1) / 2

	distances := make([]float64, size)
	var wg sync.WaitGroup

	// Divide work among workers
	numWorkers := 4
	chunkSize := n / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		start := w * chunkSize
		end := start + chunkSize
		if w == numWorkers-1 {
			end = n
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for j := i + 1; j < n; j++ {
					dx := points[i][0] - points[j][0]
					dy := points[i][1] - points[j][1]
					dist := math.Sqrt(dx*dx + dy*dy)

					// Calculate index in result array
					idx := i*n - (i*(i+1))/2 + (j - i - 1)
					distances[idx] = dist
				}
			}
		}(start, end)
	}

	wg.Wait()
	return distances
}

// ============================================================================
// Solution 7: Optimized Word Frequency
// ============================================================================

// CountWordFrequencyOptimized uses efficient string processing.
func CountWordFrequencyOptimized(docs []Document) map[string]int {
	// Preallocate map with estimated size
	wordCount := make(map[string]int, 1000)

	for _, doc := range docs {
		// Use strings.Fields (more efficient than Split)
		words := strings.Fields(strings.ToLower(doc.Content))

		for _, word := range words {
			// Trim punctuation if needed
			word = strings.Trim(word, ".,!?;:")
			if word != "" {
				wordCount[word]++
			}
		}
	}

	return wordCount
}

// ============================================================================
// Solution 8: Optimized Cache Implementation
// ============================================================================

// OptimizedCache uses RWMutex and has capacity limits.
type OptimizedCache struct {
	mu       sync.RWMutex
	data     map[string]*cacheEntry
	capacity int
	order    []string // For LRU eviction
}

type cacheEntry struct {
	value     interface{}
	createdAt int64
}

// NewOptimizedCache creates a cache with capacity limit.
func NewOptimizedCache(capacity int) *OptimizedCache {
	return &OptimizedCache{
		data:     make(map[string]*cacheEntry, capacity),
		capacity: capacity,
		order:    make([]string, 0, capacity),
	}
}

// Get retrieves a value using read lock.
func (c *OptimizedCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, ok := c.data[key]; ok {
		return entry.value, true
	}
	return nil, false
}

// Set stores a value with capacity management.
func (c *OptimizedCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check capacity and evict if needed
	if len(c.data) >= c.capacity && c.data[key] == nil {
		// Evict oldest entry (simple FIFO, could be LRU)
		if len(c.order) > 0 {
			oldest := c.order[0]
			delete(c.data, oldest)
			c.order = c.order[1:]
		}
	}

	// Add new entry
	c.data[key] = &cacheEntry{
		value:     value,
		createdAt: 0, // Could use time.Now().Unix()
	}
	c.order = append(c.order, key)
}

// Len returns the current size.
func (c *OptimizedCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// ============================================================================
// Solution 9: Optimized Filter and Transform
// ============================================================================

// FilterAndTransformOptimized preallocates and reduces allocations.
func FilterAndTransformOptimized(items []Item, minValue float64) []Result {
	// Estimate result size (assume ~50% pass filter)
	results := make([]Result, 0, len(items)/2)

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

// FilterAndTransformInPlace avoids allocation if possible.
func FilterAndTransformInPlace(items []Item, minValue float64) []Result {
	results := make([]Result, len(items))
	n := 0

	for _, item := range items {
		if item.Value >= minValue {
			results[n] = Result{
				ItemID:   item.ID,
				Score:    item.Value * 2,
				Category: "filtered",
			}
			n++
		}
	}

	return results[:n] // Return only filled portion
}

// ============================================================================
// Solution 10: Optimized Fibonacci
// ============================================================================

// FibonacciIterative uses O(n) time and O(1) space.
func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}

	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}

	return curr
}

// FibonacciMemoized uses memoization for recursive calls.
func FibonacciMemoized(n int) int {
	memo := make(map[int]int)
	return fibMemo(n, memo)
}

func fibMemo(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}

	if val, ok := memo[n]; ok {
		return val
	}

	result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
	memo[n] = result
	return result
}

// FibonacciMatrix uses matrix exponentiation - O(log n) time.
func FibonacciMatrix(n int) int {
	if n <= 1 {
		return n
	}

	// Matrix: [[1,1],[1,0]]^n gives Fibonacci numbers
	result := matrixPower([][]int{{1, 1}, {1, 0}}, n-1)
	return result[0][0]
}

func matrixPower(m [][]int, n int) [][]int {
	if n == 1 {
		return m
	}

	if n%2 == 0 {
		half := matrixPower(m, n/2)
		return matrixMultiply(half, half)
	}

	return matrixMultiply(m, matrixPower(m, n-1))
}

func matrixMultiply(a, b [][]int) [][]int {
	return [][]int{
		{a[0][0]*b[0][0] + a[0][1]*b[1][0], a[0][0]*b[0][1] + a[0][1]*b[1][1]},
		{a[1][0]*b[0][0] + a[1][1]*b[1][0], a[1][0]*b[0][1] + a[1][1]*b[1][1]},
	}
}

// ============================================================================
// Additional Optimizations: sync.Pool Example
// ============================================================================

// StringBuilderPool demonstrates sync.Pool for buffer reuse
var StringBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// BuildReportWithPool uses sync.Pool to reuse buffers.
func BuildReportWithPool(items []Item) string {
	buf := StringBuilderPool.Get().(*strings.Builder)
	buf.Reset()
	defer StringBuilderPool.Put(buf)

	buf.Grow(len(items)*100 + 100)
	buf.WriteString("=== Report ===\n")

	for _, item := range items {
		fmt.Fprintf(buf, "ID: %d, Name: %s, Value: %.2f\n",
			item.ID, item.Name, item.Value)
	}

	buf.WriteString("=== End Report ===\n")

	return buf.String()
}
