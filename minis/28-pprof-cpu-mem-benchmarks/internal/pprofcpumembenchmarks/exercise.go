//go:build !solution && !reference

package pprofcpumembenchmarks

// TODO: Import required packages
// You'll need:
// - "fmt" for error formatting and string operations
// - "math" for mathematical calculations
// - "strings" for string builder operations
// - "encoding/json" for JSON marshaling (optional optimization)
//
// import (
//     "fmt"
//     "math"
//     "strings"
// )

// ============================================================================
// PROFILING: Understanding CPU and Memory Performance
// ============================================================================
//
// This exercise teaches you how to identify and fix performance bottlenecks
// using Go's built-in profiling tools.
//
// What is profiling?
// - CPU Profiling: Shows where your program spends execution time
// - Memory Profiling: Shows where your program allocates memory
// - Helps identify "hot spots" (slow code) and memory leaks
//
// How to use Go's profiling tools:
// 1. CPU Profile: go test -cpuprofile=cpu.prof -bench=.
// 2. Memory Profile: go test -memprofile=mem.prof -bench=.
// 3. Analyze: go tool pprof cpu.prof
// 4. Visualize: go tool pprof -http=:8080 cpu.prof
//
// Common commands in pprof:
// - top: Show top functions by CPU/memory usage
// - list FunctionName: Show line-by-line breakdown
// - web: Generate graphical visualization
//
// Benchmarking in Go:
// - Use testing.B for benchmarks
// - Run with: go test -bench=.
// - Iterations: b.N is automatically adjusted
// - Reset timer: b.ResetTimer() after setup
// - Report allocations: b.ReportAllocs()
//
// Key optimization strategies:
// 1. Algorithm complexity: O(n²) → O(n log n) → O(n)
// 2. Reduce allocations: Preallocate slices, reuse buffers
// 3. String operations: Use strings.Builder instead of concatenation
// 4. Avoid repeated calculations: Cache results when possible
//
// ============================================================================

// ============================================================================
// Exercise 1: Optimize Prime Number Finding (CPU Bottleneck)
// ============================================================================

// FindPrimes finds all prime numbers up to n.
func FindPrimes(n int) []int {
	// TODO: Replace this naive O(n^2) implementation with the Sieve of Eratosthenes (O(n log log n)).

	// The Sieve of Eratosthenes is a highly efficient algorithm for finding all prime numbers up to a specified integer.

	// Step 1: Create a boolean slice `isPrime` of size `n + 1`.
	// - `isPrime := make([]bool, n+1)`
	// - Initialize all entries from 2 to `n` as `true`.

	// Step 2: Implement the "sieving" process.
	// - Loop from `p = 2` up to `p*p <= n`.
	// - Inside the loop, if `isPrime[p]` is still `true`, then `p` is a prime.
	// - For each prime `p`, mark all of its multiples as not prime.
	//   - Start from `p*p` and increment by `p`: `for i := p * p; i <= n; i += p`.
	//   - `isPrime[i] = false`.
	// - Why start at `p*p`? All smaller multiples of `p` would have already been marked by smaller primes.

	// Step 3: Collect the results.
	// - Create a new slice to hold the prime numbers. Pre-allocate its capacity for better performance (e.g., `make([]int, 0, n/10)`).
	// - Loop from 2 to `n`. If `isPrime[i]` is `true`, append `i` to your results slice.

	// Step 4: Return the slice of primes.
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
// Exercise 2: Reduce String Allocation (Memory Bottleneck)
// ============================================================================

// BuildReport creates a report string from items.
func BuildReport(items []Item) string {
	// TODO: Replace this inefficient string concatenation with a `strings.Builder`.

	// String concatenation with `+=` is very inefficient in a loop.
	// Each `+=` creates a completely new string and copies the contents of the old and new parts.
	// This leads to a large number of memory allocations and a lot of work for the garbage collector.

	// Step 1: Create a `strings.Builder`.
	// - `var buf strings.Builder`

	// Step 2: (Optional but recommended) Pre-allocate memory.
	// - If you can estimate the final size of the string, `buf.Grow(estimatedSize)` can prevent the builder from having to re-allocate its internal buffer as it grows. This is a significant optimization.

	// Step 3: Write to the builder instead of concatenating.
	// - Use `buf.WriteString()` for constant strings.
	// - Use `fmt.Fprintf(&buf, ...)` for formatted strings. `strings.Builder` implements the `io.Writer` interface, making it very flexible.

	// Step 4: Get the final string.
	// - `return buf.String()`
	// - This performs the final allocation and copy to create the immutable result string.
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
func SearchDocuments(docs []Document, query string) []Document {
	// TODO: Optimize this search by reducing redundant string operations and allocations.

	// The current implementation calls `containsCaseInsensitive` twice, which in turn calls `strings.ToLower` on the query twice for every single document. This is very inefficient.

	// Step 1: Convert the `query` to lowercase *once*, before the loop begins.
	// - `queryLower := strings.ToLower(query)`

	// Step 2: Pre-allocate the `results` slice.
	// - We don't know the exact number of matches, but we can make a reasonable guess. Estimating that 10% of documents will match is a good start.
	// - `results := make([]Document, 0, len(docs)/10)`

	// Step 3: Loop through the documents.
	// - Inside the loop, convert the document's title and content to lowercase.
	// - Use `strings.Contains` to check if the lowercase title or content contains the `queryLower`.
	// - If it does, append the document to the `results` slice.

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

func containsCaseInsensitive(text, substr string) bool {
	// This helper function is the source of the inefficiency.
	// By moving the `strings.ToLower(substr)` call outside the main loop, we can eliminate a huge number of redundant operations.
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}

// ============================================================================
// Exercise 4: Reduce Allocations in Loop
// ============================================================================

// ProcessItems processes a slice of items and returns results.
func ProcessItems(items []Item) []Result {
	// TODO: Preallocate the results slice to avoid multiple reallocations during the loop.

	// When you `append` to a slice that is at full capacity, Go allocates a new, larger underlying array and copies all existing elements to it.
	// Doing this repeatedly in a loop is inefficient.

	// Step 1: Pre-allocate the slice.
	// - Since we know the final number of results will be exactly `len(items)`, we can create the slice with the perfect capacity.
	// - `results := make([]Result, 0, len(items))`
	// - This creates a slice with length 0 but capacity `len(items)`. All subsequent `append` calls will simply use the existing allocated memory and will not cause any re-allocations.

	// Step 2: Loop and append.
	// - The rest of the function remains the same. The `append` operations will now be much more efficient.
	var results []Result

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

func calculateScore(item Item) float64 {
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
//
// Problems with manual JSON building:
// 1. String concatenation (see Exercise 2) - O(n²) copying
// 2. Doesn't handle escaping properly (security issue!)
// 3. Error-prone (easy to create invalid JSON)
// 4. Slower than encoding/json package
//
// Why encoding/json is better:
// - Uses efficient buffer management internally
// - Handles all escaping (quotes, newlines, unicode)
// - Validates JSON structure
// - Can use struct tags for customization
// - Can stream large JSON without loading entire result in memory
//
// Alternative: strings.Builder
// - If you must build JSON manually (custom format)
// - Still much better than concatenation
// - But use encoding/json for production code!
//
// Memory comparison (1000 items):
// - String concatenation: ~1000 allocations, 1.25MB copied
// - strings.Builder: ~2 allocations, 50KB copied
// - encoding/json: ~1 allocation, 50KB copied + proper escaping
//
// TODO: Use encoding/json.MarshalIndent for correctness and efficiency
//
// func FormatItemsAsJSON(items []Item) string {
//     // Create simplified representation for JSON
//     type SimpleItem struct {
//         ID    int     `json:"id"`
//         Name  string  `json:"name"`
//         Value float64 `json:"value"`
//     }
//
//     simpleItems := make([]SimpleItem, len(items))
//     for i, item := range items {
//         simpleItems[i] = SimpleItem{
//             ID:    item.ID,
//             Name:  item.Name,
//             Value: item.Value,
//         }
//     }
//
//     bytes, err := json.MarshalIndent(simpleItems, "", "  ")
//     if err != nil {
//         return ""
//     }
//
//     return string(bytes)
// }

func FormatItemsAsJSON(items []Item) string {
	// TODO: Replace this manual, inefficient, and incorrect JSON building with the standard `encoding/json` package.

	// The current implementation is bad for several reasons:
	// 1. Inefficient: It uses string concatenation in a loop (see Exercise 2).
	// 2. Incorrect/Unsafe: It does not handle string escaping. If an item's name contained a quote (`"`), it would produce invalid JSON. This can be a security risk.

	// Step 1: Define a new struct that represents the desired JSON structure.
	// - It's good practice to create a specific struct for your JSON output rather than exporting the fields of your internal `Item` struct directly.
	// - Use `json` struct tags to control the field names in the output (e.g., `json:"id"`).
	/*
		type SimpleItem struct {
			ID    int     `json:"id"`
			Name  string  `json:"name"`
			Value float64 `json:"value"`
		}
	*/

	// Step 2: Convert your `[]Item` slice to a `[]SimpleItem` slice.

	// Step 3: Use `json.MarshalIndent` to serialize the slice to a formatted JSON byte slice.
	// - `bytes, err := json.MarshalIndent(simpleItems, "", "  ")`
	// - This handles all the complex parts: building the string efficiently, escaping special characters, and formatting the output nicely.

	// Step 4: Convert the byte slice to a string and return it.
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
func ComputeDistances(points [][2]float64) []float64 {
	// TODO: Preallocate the result slice with the correct size to avoid reallocations.

	// The number of unique pairs in a set of `n` points is `n * (n - 1) / 2`.

	// Step 1: Calculate the exact number of distances that will be computed.
	// - `n := len(points)`
	// - `size := n * (n - 1) / 2`

	// Step 2: Pre-allocate the `distances` slice with this capacity.
	// - `distances := make([]float64, 0, size)`

	// Step 3: Perform the calculations and append to the slice. The appends will now be efficient.
	var distances []float64

	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
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
func CountWordFrequency(docs []Document) map[string]int {
	// TODO: Preallocate the map with an estimated size and use a more efficient split function.

	// When a map grows, it has to rehash all its existing elements into a new, larger underlying array. This is expensive.
	// Pre-allocating the map with an estimate of the final size can prevent most of these rehashes.

	// Step 1: Pre-allocate the map.
	// - `wordCount := make(map[string]int, estimatedSize)`
	// - A good estimate might be `len(docs) * 50` if you assume around 50 unique words per document.

	// Step 2: In the loop, use `strings.Fields` instead of `strings.Split(text, " ")`.
	// - `strings.Fields` is generally faster and correctly handles all forms of whitespace (tabs, newlines, multiple spaces).

	// Step 3: (Optional) A more advanced optimization could be to trim punctuation from words.
	// - `word = strings.Trim(word, ".,!?;:")`
	wordCount := make(map[string]int)

	for _, doc := range docs {
		// TODO: Use strings.Fields instead of strings.Split
		text := strings.ToLower(doc.Content)
		words := strings.Split(text, " ")

		for _, word := range words {
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
type SimpleCache struct {
	mu       sync.RWMutex
	data     map[string]interface{}
	capacity int
	order    []string
}

func NewSimpleCache(capacity int) *SimpleCache {
	// TODO: Initialize the cache correctly.
	// - A `SimpleCache` needs a `sync.RWMutex` for thread safety.
	// - The `data` map needs to be initialized. Pre-allocating it with the given `capacity` is a good optimization.
	// - The `order` slice, used for our simple eviction policy, should also be pre-allocated.
	return &SimpleCache{
		data: make(map[string]interface{}),
	}
}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
	// TODO: Implement a thread-safe read.
	// - This is a read-only operation, so use a read lock (`c.mu.RLock()`) for better performance under concurrent reads.
	// - Defer the unlock (`c.mu.RUnlock()`).
	// - Read from the map and return.
	val, ok := c.data[key]
	return val, ok
}

func (c *SimpleCache) Set(key string, value interface{}) {
	// TODO: Implement a thread-safe write with a capacity limit and eviction.

	// Step 1: Acquire a full write lock (`c.mu.Lock()`).
	// Step 2: Defer the unlock.
	// Step 3: Implement the eviction policy.
	// - Check if the key already exists. If it does, we're just updating, so no eviction is needed.
	// - If the key does not exist AND the cache is at capacity (`len(c.data) >= c.capacity`):
	//   - We need to evict the oldest item. In our simple FIFO (First-In, First-Out) implementation, this is the first item in the `c.order` slice.
	//   - Get the key of the oldest item: `oldest := c.order[0]`.
	//   - Remove it from the `c.order` slice: `c.order = c.order[1:]`.
	//   - Delete it from the `c.data` map: `delete(c.data, oldest)`.
	// Step 4: Set the new key-value pair in the map.
	// Step 5: Add the new key to the end of the `c.order` slice to track its age.
	c.data[key] = value
}

// ============================================================================
// Exercise 9: Optimize Slice Operations
// ============================================================================

// FilterAndTransform filters items and transforms them to results.
func FilterAndTransform(items []Item, minValue float64) []Result {
	// TODO: Preallocate the result slice with a reasonable estimated capacity.

	// In this scenario, we don't know the exact final size of the `results` slice because the filter is conditional.
	// However, appending one-by-one is still inefficient if many elements pass the filter.

	// Step 1: Choose an allocation strategy.
	// - No pre-allocation: `var results []Result` (bad if many items pass).
	// - Full pre-allocation: `make([]Result, 0, len(items))` (best for performance, but can waste memory if the filter is very selective).
	// - Estimated pre-allocation: `make([]Result, 0, len(items)/2)` (a good compromise).

	// Step 2: Implement using an estimated pre-allocation.
	results := make([]Result, 0, len(items)/2) // Let's estimate 50% will pass.

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
func Fibonacci(n int) int {
	// TODO: Replace this exponential-time recursive implementation with a linear-time iterative one.

	// The current implementation is O(2^n), which is extremely slow for `n` larger than ~40.
	// It re-calculates the same Fibonacci numbers over and over again.

	// Step 1: Handle the base cases.
	// - If `n <= 1`, the Fibonacci number is `n` itself. Return `n`.

	// Step 2: Initialize variables for the two previous numbers in the sequence.
	// - `prev := 0`
	// - `curr := 1`

	// Step 3: Loop from 2 up to `n`.
	// - In each iteration, calculate the next number in the sequence.
	// - `next := prev + curr`
	// - Then, update `prev` and `curr` for the next iteration. A parallel assignment is great for this:
	// - `prev, curr = curr, next` or `prev, curr = curr, prev + curr`

	// Step 4: After the loop, `curr` will hold the nth Fibonacci number. Return it.
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

// After implementing all optimizations:
// - Run: go test -bench=. -benchmem
// - Compare: BenchmarkNaive vs BenchmarkOptimized
// - Check: Allocations per operation (should be much lower)
// - Profile: go test -cpuprofile=cpu.prof -bench=.
// - Analyze: go tool pprof cpu.prof
// - Compare with solution.go to see detailed implementations
