//go:build !solution
// +build !solution

package exercise

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
// TODO: This uses a naive O(n²) algorithm. Optimize it using the Sieve of Eratosthenes.
//
// Current implementation analysis:
// - Time complexity: O(n²) - checks every number against every previous number
// - For n=10,000: ~50 million operations
// - For n=100,000: ~5 billion operations (very slow!)
//
// Why is this slow?
// - For each number i, we check divisibility against all numbers from 2 to i-1
// - Most of these checks are unnecessary
// - Example: For 97, we check 2,3,4,5...96 (95 checks, but only need to check primes up to √97)
//
// Optimization strategy - Sieve of Eratosthenes:
// 1. Create a boolean array of size n+1 (all initially true = prime)
// 2. Start with 2 (first prime)
// 3. Mark all multiples of 2 as composite (4, 6, 8, 10...)
// 4. Move to next unmarked number (3), mark its multiples (6, 9, 12...)
// 5. Continue until √n
// 6. Collect all unmarked numbers as primes
//
// Time complexity improvement: O(n log log n) vs O(n²)
// - For n=100,000: ~150,000 operations vs 5 billion
// - This is ~33,000x faster!
//
// Memory consideration:
// - Sieve uses O(n) space (boolean array)
// - Trade-off: More memory for much faster execution
// - For n=1,000,000: ~1MB of memory (acceptable)
//
// TODO: Implement the Sieve of Eratosthenes algorithm
//
// func FindPrimes(n int) []int {
//     if n < 2 {
//         return nil
//     }
//
//     // Create boolean array - true means "is prime"
//     isPrime := make([]bool, n+1)
//     for i := 2; i <= n; i++ {
//         isPrime[i] = true
//     }
//
//     // Sieve algorithm
//     for i := 2; i*i <= n; i++ {
//         if isPrime[i] {
//             // Mark all multiples of i as composite
//             for j := i * i; j <= n; j += i {
//                 isPrime[j] = false
//             }
//         }
//     }
//
//     // Collect primes
//     primes := make([]int, 0, n/10) // Approximate: ~n/ln(n) primes
//     for i := 2; i <= n; i++ {
//         if isPrime[i] {
//             primes = append(primes, i)
//         }
//     }
//
//     return primes
// }

func FindPrimes(n int) []int {
	// TODO: Replace this naive implementation with the Sieve of Eratosthenes
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
// TODO: This uses string concatenation which creates many intermediate allocations.
//
// Why is string concatenation slow?
// - Strings in Go are IMMUTABLE (cannot be changed after creation)
// - Each += creates a NEW string and copies all previous content
// - Example with 1000 items:
//   * Iteration 1: Copy 50 bytes
//   * Iteration 2: Copy 100 bytes (50 old + 50 new)
//   * Iteration 3: Copy 150 bytes
//   * Total: 50 + 100 + 150 + ... + 50,000 = ~1.25 million bytes copied!
//
// Memory allocations:
// - Each concatenation allocates a new string on the heap
// - For n items: n allocations + n string copies
// - Garbage collector must clean up n-1 temporary strings
// - This creates GC pressure and slows down the program
//
// Optimization: strings.Builder
// - Maintains a single growing buffer ([]byte)
// - No copying until you call String()
// - Can preallocate capacity with Grow()
// - Only 1-2 allocations total!
//
// How strings.Builder works:
// 1. Internally uses a []byte slice
// 2. Appends to slice (amortized O(1) if capacity is sufficient)
// 3. Grows slice by doubling when needed (like slice append)
// 4. Final String() converts []byte to string (1 copy)
//
// Memory improvement:
// - Before: n allocations, ~n²/2 bytes copied
// - After: 1-2 allocations, ~n bytes copied
// - For 1000 items: 50KB vs 1.25MB = 25x less memory traffic
//
// TODO: Implement using strings.Builder
//
// func BuildReport(items []Item) string {
//     var buf strings.Builder
//
//     // Preallocate capacity (estimate: 100 bytes per item + header/footer)
//     buf.Grow(len(items)*100 + 100)
//
//     buf.WriteString("=== Report ===\n")
//
//     for _, item := range items {
//         // fmt.Fprintf writes to io.Writer (strings.Builder implements this)
//         fmt.Fprintf(&buf, "ID: %d, Name: %s, Value: %.2f\n",
//             item.ID, item.Name, item.Value)
//     }
//
//     buf.WriteString("=== End Report ===\n")
//
//     return buf.String()
// }

func BuildReport(items []Item) string {
	// TODO: Replace string concatenation with strings.Builder
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
// TODO: This has multiple inefficiencies. Optimize the string processing.
//
// Current inefficiencies:
// 1. strings.ToLower() called multiple times for the same string
//    - Each call allocates a new string
//    - For 1000 docs with query in loop: 2000 allocations
//
// 2. Case-insensitive comparison is expensive
//    - strings.Contains(strings.ToLower(text), strings.ToLower(query))
//    - Allocates 2 lowercase strings per comparison
//
// Optimization strategies:
// 1. Convert query to lowercase ONCE (outside loop)
//    - Reduces allocations from 2n to n+1
//
// 2. Consider precomputing lowercase content
//    - If searching repeatedly, store lowercase version
//    - Trade-off: More memory, but much faster searches
//
// 3. For better performance with many searches: Build an inverted index
//    - Preprocessing: O(total words) time
//    - Search: O(query words) time vs O(documents) time
//    - Used by search engines (Elasticsearch, etc.)
//
// Memory analysis:
// - Before: 2 allocations per document per search
// - After: 1 allocation per document per search
// - With index: 0 allocations per search (after preprocessing)
//
// TODO: Optimize by reducing repeated string operations
//
// func SearchDocuments(docs []Document, query string) []Document {
//     // Convert query to lowercase ONCE
//     queryLower := strings.ToLower(query)
//
//     // Preallocate result slice (estimate: ~10% of docs match)
//     results := make([]Document, 0, len(docs)/10)
//
//     for _, doc := range docs {
//         // Convert document text to lowercase
//         titleLower := strings.ToLower(doc.Title)
//         contentLower := strings.ToLower(doc.Content)
//
//         // Check if query is in title or content
//         if strings.Contains(titleLower, queryLower) ||
//            strings.Contains(contentLower, queryLower) {
//             results = append(results, doc)
//         }
//     }
//
//     return results
// }

func SearchDocuments(docs []Document, query string) []Document {
	// TODO: Optimize this search by reducing string allocations
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
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}

// ============================================================================
// Exercise 4: Reduce Allocations in Loop
// ============================================================================

// ProcessItems processes a slice of items and returns results.
// TODO: This allocates inefficiently. Preallocate the results slice.
//
// Slice growth mechanics in Go:
// - When append() exceeds capacity, a new array is allocated
// - New capacity = old capacity * 2 (approximately)
// - Old contents are copied to new array
// - Old array is garbage collected
//
// Growth sequence for 1000 items:
// - Start: cap=0
// - After 1 append: cap=1 (allocate, copy 0 → 1 items)
// - After 2 appends: cap=2 (allocate, copy 1 → 2 items)
// - After 3 appends: cap=4 (allocate, copy 2 → 4 items)
// - After 5 appends: cap=8 (allocate, copy 4 → 8 items)
// - ...
// - After 1000 appends: cap=1024 (log₂(1000) ≈ 10 allocations)
//
// Total cost without preallocation:
// - ~10 allocations
// - ~2000 items copied (1 + 2 + 4 + 8 + ... + 512)
//
// With preallocation:
// - 1 allocation
// - 0 items copied
//
// How to preallocate:
// - If exact size known: make([]Type, length)
// - If size unknown but bounded: make([]Type, 0, capacity)
//
// Memory consideration:
// - Overallocation is cheap compared to reallocation
// - Wasting 10% capacity is better than 10 reallocations
//
// TODO: Preallocate the results slice with correct capacity
//
// func ProcessItems(items []Item) []Result {
//     // Preallocate with exact capacity (we know final size)
//     results := make([]Result, 0, len(items))
//
//     for _, item := range items {
//         category := determineCategory(item.Value)
//
//         result := Result{
//             ItemID:    item.ID,
//             Score:     calculateScore(item),
//             Category:  category,
//             Processed: item.Timestamp,
//         }
//         results = append(results, result)
//     }
//
//     return results
// }

func ProcessItems(items []Item) []Result {
	// TODO: Preallocate results slice to avoid multiple reallocations
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
	// TODO: Replace manual string concatenation with encoding/json or strings.Builder
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
// TODO: Preallocate result slice to reduce allocations.
//
// Algorithm analysis:
// - Computes distance for all pairs: n*(n-1)/2 pairs
// - For n=1000: 499,500 distances
// - Time complexity: O(n²) - unavoidable for all-pairs
//
// Current inefficiency:
// - Results slice grows dynamically
// - log₂(n²/2) reallocations
// - For n=1000: ~19 reallocations
//
// Optimization: Preallocate
// - Calculate result size: n*(n-1)/2
// - Allocate once: make([]float64, 0, size)
// - No reallocations during computation
//
// Memory consideration:
// - For n=1000: 499,500 float64 values = ~3.8 MB
// - Preallocation saves ~50KB of intermediate arrays
// - More importantly: Saves CPU time (no copying)
//
// Advanced optimization (not required):
// - Parallel processing: Divide pairs among goroutines
// - Trade-off: Coordination overhead vs parallelism benefit
// - Beneficial for n > 10,000
//
// TODO: Preallocate the result slice with correct capacity
//
// func ComputeDistances(points [][2]float64) []float64 {
//     n := len(points)
//     size := n * (n - 1) / 2  // Number of unique pairs
//
//     // Preallocate result slice
//     distances := make([]float64, 0, size)
//
//     for i := 0; i < n; i++ {
//         for j := i + 1; j < n; j++ {
//             dx := points[i][0] - points[j][0]
//             dy := points[i][1] - points[j][1]
//             dist := math.Sqrt(dx*dx + dy*dy)
//             distances = append(distances, dist)
//         }
//     }
//
//     return distances
// }

func ComputeDistances(points [][2]float64) []float64 {
	// TODO: Preallocate result slice with correct size: n*(n-1)/2
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
// TODO: This has several inefficiencies in string processing.
//
// Current inefficiencies:
// 1. Map not preallocated
//    - Default capacity is very small
//    - Grows dynamically (rehashing on each growth)
//    - Rehashing copies all entries to new backing array
//
// 2. strings.Split vs strings.Fields
//    - Split(" ") only splits on space (not tab, newline, etc.)
//    - Fields() splits on any whitespace (more flexible)
//    - Fields() is slightly faster (optimized)
//
// 3. Repeated string.ToLower
//    - Done per document (good!)
//    - But could cache if documents are reused
//
// Map growth mechanics:
// - Maps in Go use hash tables
// - When load factor exceeds threshold (~6.5), map doubles in size
// - Doubling requires rehashing ALL entries
// - For 10,000 unique words: ~14 rehashes
//
// Optimization: Preallocate map
// - Estimate final size
// - make(map[string]int, estimatedSize)
// - Avoids most rehashes
//
// Memory consideration:
// - Overallocation wastes some memory
// - But rehashing wastes CPU time
// - Better to overallocate by 20% than to rehash
//
// TODO: Optimize with map preallocation and strings.Fields
//
// func CountWordFrequency(docs []Document) map[string]int {
//     // Preallocate map with estimated size
//     // Estimate: 100 unique words per document
//     wordCount := make(map[string]int, len(docs)*100)
//
//     for _, doc := range docs {
//         // Use strings.Fields (more efficient than Split)
//         words := strings.Fields(strings.ToLower(doc.Content))
//
//         for _, word := range words {
//             // Trim punctuation if needed
//             word = strings.Trim(word, ".,!?;:")
//             if word != "" {
//                 wordCount[word]++
//             }
//         }
//     }
//
//     return wordCount
// }

func CountWordFrequency(docs []Document) map[string]int {
	// TODO: Preallocate map with estimated size
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
// TODO: This cache has serious problems. Fix them!
//
// Current problems:
// 1. Uses sync.Mutex for all operations (even reads!)
//    - Reads don't modify data, can use RLock
//    - Multiple readers can proceed concurrently with RWMutex
//
// 2. No capacity limit
//    - Cache grows indefinitely = memory leak!
//    - Real caches need eviction policy
//
// 3. No eviction policy
//    - LRU (Least Recently Used): Track access time, evict oldest
//    - TTL (Time To Live): Expire entries after timeout
//    - FIFO (First In First Out): Evict in insertion order
//
// 4. Not thread-safe (missing mutex!)
//    - Concurrent map access causes panic
//
// Why RWMutex is better for caches:
// - Caches are read-heavy (90% reads, 10% writes)
// - RWMutex allows multiple concurrent readers
// - Only one writer at a time (blocks readers)
// - Much better throughput for read-heavy workloads
//
// Cache eviction:
// - LRU is most common (best hit rate)
// - Implementation: map + doubly-linked list
// - Get: Move item to front
// - Set: Add to front, evict from back if over capacity
//
// Memory consideration:
// - LRU needs extra memory for list pointers
// - Trade-off: ~16 bytes per entry vs unbounded growth
//
// TODO: Fix the cache implementation
//
// type SimpleCache struct {
//     mu       sync.RWMutex
//     data     map[string]interface{}
//     capacity int
//     order    []string  // For FIFO eviction (simple)
// }
//
// func NewSimpleCache(capacity int) *SimpleCache {
//     return &SimpleCache{
//         data:     make(map[string]interface{}, capacity),
//         capacity: capacity,
//         order:    make([]string, 0, capacity),
//     }
// }
//
// func (c *SimpleCache) Get(key string) (interface{}, bool) {
//     c.mu.RLock()
//     defer c.mu.RUnlock()
//
//     val, ok := c.data[key]
//     return val, ok
// }
//
// func (c *SimpleCache) Set(key string, value interface{}) {
//     c.mu.Lock()
//     defer c.mu.Unlock()
//
//     // Check capacity and evict if needed (FIFO)
//     if len(c.data) >= c.capacity && c.data[key] == nil {
//         // Evict oldest entry
//         if len(c.order) > 0 {
//             oldest := c.order[0]
//             delete(c.data, oldest)
//             c.order = c.order[1:]
//         }
//     }
//
//     c.data[key] = value
//     c.order = append(c.order, key)
// }

type SimpleCache struct {
	// TODO: Add proper fields (RWMutex, map, capacity limit, eviction tracking)
	data map[string]interface{}
}

func NewSimpleCache(capacity int) *SimpleCache {
	// TODO: Initialize with proper fields and capacity limit
	return &SimpleCache{
		data: make(map[string]interface{}),
	}
}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
	// TODO: Use RLock/RUnlock for concurrent reads
	val, ok := c.data[key]
	return val, ok
}

func (c *SimpleCache) Set(key string, value interface{}) {
	// TODO: Use Lock/Unlock and implement capacity checking with eviction
	c.data[key] = value
}

// ============================================================================
// Exercise 9: Optimize Slice Operations
// ============================================================================

// FilterAndTransform filters items and transforms them to results.
// TODO: Preallocate result slice to avoid allocations.
//
// Filtering pattern analysis:
// - Input: n items
// - Output: k items (where k ≤ n)
// - Problem: We don't know k in advance!
//
// Allocation strategies:
// 1. Don't preallocate: var results []Result
//    - Pro: No wasted memory
//    - Con: log₂(k) reallocations
//
// 2. Preallocate for all: make([]Result, 0, n)
//    - Pro: No reallocations
//    - Con: May waste memory if filter is selective
//
// 3. Preallocate with estimate: make([]Result, 0, n/2)
//    - Pro: Balance between allocations and waste
//    - Con: Needs domain knowledge
//
// When to use each strategy:
// - Strategy 1: Output size << input size (e.g., 1% filter)
// - Strategy 2: Output size ≈ input size (e.g., 90% filter)
// - Strategy 3: Output size ≈ 50% (most common)
//
// Memory waste calculation:
// - Strategy 2 with 10% filter: 90% wasted capacity
// - For 10,000 items: 9,000 unused slots × 40 bytes = 360KB wasted
// - But saves ~14 reallocations and copying ~20,000 items
//
// General rule:
// - If you can estimate output size within 2x, preallocate
// - Better to waste some memory than to copy repeatedly
//
// TODO: Preallocate with reasonable estimate
//
// func FilterAndTransform(items []Item, minValue float64) []Result {
//     // Estimate: assume ~50% pass filter (adjust based on your data)
//     results := make([]Result, 0, len(items)/2)
//
//     for _, item := range items {
//         if item.Value >= minValue {
//             result := Result{
//                 ItemID:   item.ID,
//                 Score:    item.Value * 2,
//                 Category: "filtered",
//             }
//             results = append(results, result)
//         }
//     }
//
//     return results
// }

func FilterAndTransform(items []Item, minValue float64) []Result {
	// TODO: Preallocate with estimated capacity (e.g., len(items)/2)
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
//
// Why is naive recursion slow?
// - Fibonacci(n) = Fibonacci(n-1) + Fibonacci(n-2)
// - Recomputes same values many times
// - Example: Fibonacci(5)
//   * Fib(5) calls Fib(4) and Fib(3)
//   * Fib(4) calls Fib(3) and Fib(2)
//   * Fib(3) is called TWICE (and will call Fib(2) twice)
//   * Fib(2) is called THREE times
//
// Time complexity: O(2ⁿ) - exponential!
// - Fibonacci(10): ~177 calls
// - Fibonacci(20): ~21,891 calls
// - Fibonacci(40): ~204,668,309 calls (billions!)
//
// Space complexity: O(n) - call stack depth
//
// Optimization 1: Iterative approach
// - Time: O(n), Space: O(1)
// - Keep only last two values
// - No recursion overhead
//
// func FibonacciIterative(n int) int {
//     if n <= 1 {
//         return n
//     }
//     prev, curr := 0, 1
//     for i := 2; i <= n; i++ {
//         prev, curr = curr, prev+curr
//     }
//     return curr
// }
//
// Optimization 2: Memoization (cache results)
// - Time: O(n), Space: O(n)
// - Store computed values in map
// - Each value computed once
//
// func FibonacciMemoized(n int) int {
//     memo := make(map[int]int)
//     return fibMemo(n, memo)
// }
//
// func fibMemo(n int, memo map[int]int) int {
//     if n <= 1 {
//         return n
//     }
//     if val, ok := memo[n]; ok {
//         return val
//     }
//     result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
//     memo[n] = result
//     return result
// }
//
// Optimization 3: Matrix exponentiation (advanced)
// - Time: O(log n), Space: O(1)
// - Uses mathematical property of Fibonacci
// - Overkill for most use cases
//
// Performance comparison (Fibonacci(40)):
// - Naive recursion: ~2 seconds
// - Memoization: ~0.00001 seconds (200,000x faster!)
// - Iterative: ~0.000001 seconds (2,000,000x faster!)
//
// TODO: Implement iterative version

func Fibonacci(n int) int {
	// TODO: Replace exponential recursion with iterative approach
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
