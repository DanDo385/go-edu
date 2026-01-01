//go:build !solution && !reference

package pprofcpumembenchmarks

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// Exercise 1: Prime Number Finding (naive vs optimized)
// ============================================================================

// FindPrimes is a simple, intentionally-naive implementation.
func FindPrimes(n int) []int {
	if n < 2 {
		return []int{}
	}

	out := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime(i) {
			out = append(out, i)
		}
	}
	return out
}

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	if x == 2 {
		return true
	}
	if x%2 == 0 {
		return false
	}
	for d := 3; d*d <= x; d += 2 {
		if x%d == 0 {
			return false
		}
	}
	return true
}

// FindPrimesOptimized uses the Sieve of Eratosthenes.
func FindPrimesOptimized(n int) []int {
	if n < 2 {
		return nil
	}

	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	primes := make([]int, 0, n/10)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

// ============================================================================
// Exercise 2: String Building (naive vs optimized)
// ============================================================================

func BuildReport(items []Item) string {
	// Intentionally naive concatenation.
	s := "=== Report ===\n"
	for _, item := range items {
		s += fmt.Sprintf("ID: %d, Name: %s, Value: %.2f\n", item.ID, item.Name, item.Value)
	}
	s += "=== End Report ===\n"
	return s
}

func BuildReportOptimized(items []Item) string {
	var buf strings.Builder
	buf.Grow(len(items)*100 + 100)

	buf.WriteString("=== Report ===\n")
	for _, item := range items {
		fmt.Fprintf(&buf, "ID: %d, Name: %s, Value: %.2f\n", item.ID, item.Name, item.Value)
	}
	buf.WriteString("=== End Report ===\n")
	return buf.String()
}

var StringBuilderPool = sync.Pool{
	New: func() any { return &strings.Builder{} },
}

func BuildReportWithPool(items []Item) string {
	buf := StringBuilderPool.Get().(*strings.Builder)
	buf.Reset()
	defer StringBuilderPool.Put(buf)

	buf.Grow(len(items)*100 + 100)
	buf.WriteString("=== Report ===\n")
	for _, item := range items {
		fmt.Fprintf(buf, "ID: %d, Name: %s, Value: %.2f\n", item.ID, item.Name, item.Value)
	}
	buf.WriteString("=== End Report ===\n")
	return buf.String()
}

// ============================================================================
// Exercise 3: Document Search (naive vs optimized)
// ============================================================================

func SearchDocuments(docs []Document, query string) []Document {
	// Intentionally naive: repeated lowercasing and string work.
	var out []Document
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(doc.Content), strings.ToLower(query)) {
			out = append(out, doc)
		}
	}
	return out
}

func SearchDocumentsOptimized(docs []Document, query string) []Document {
	queryLower := strings.ToLower(query)
	results := make([]Document, 0, len(docs)/10)

	for _, doc := range docs {
		titleLower := strings.ToLower(doc.Title)
		contentLower := strings.ToLower(doc.Content)

		if strings.Contains(titleLower, queryLower) || strings.Contains(contentLower, queryLower) {
			results = append(results, doc)
		}
	}
	return results
}

// ============================================================================
// Exercise 4: Process Items (naive vs optimized)
// ============================================================================

func ProcessItems(items []Item) []Result {
	var out []Result
	for _, item := range items {
		out = append(out, Result{
			ItemID:    item.ID,
			Score:     calculateScore(item),
			Category:  determineCategory(item.Value),
			Processed: item.Timestamp,
		})
	}
	return out
}

func ProcessItemsOptimized(items []Item) []Result {
	results := make([]Result, 0, len(items))
	for _, item := range items {
		category := determineCategory(item.Value)
		results = append(results, Result{
			ItemID:    item.ID,
			Score:     calculateScore(item),
			Category:  category,
			Processed: item.Timestamp,
		})
	}
	return results
}

func determineCategory(v float64) string {
	switch {
	case v < 33:
		return "low"
	case v < 66:
		return "medium"
	default:
		return "high"
	}
}

func calculateScore(item Item) float64 {
	// Not performance-critical for tests; kept simple and deterministic.
	return item.Value * 1.5
}

// ============================================================================
// Exercise 5: JSON Formatting (naive + optimized)
// ============================================================================

func FormatItemsAsJSON(items []Item) string {
	// Very naive manual formatting (not robust JSON, but good enough for the tests).
	var b strings.Builder
	b.WriteString("[\n")
	for i, it := range items {
		fmt.Fprintf(&b, "  {\"id\": %d, \"name\": \"%s\", \"value\": %.2f}", it.ID, it.Name, it.Value)
		if i < len(items)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}
	b.WriteString("]\n")
	return b.String()
}

func FormatItemsAsJSONOptimized(items []Item) string {
	type SimpleItem struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}

	simpleItems := make([]SimpleItem, len(items))
	for i, item := range items {
		simpleItems[i] = SimpleItem{ID: item.ID, Name: item.Name, Value: item.Value}
	}

	bytes, err := json.MarshalIndent(simpleItems, "", "  ")
	if err != nil {
		return ""
	}
	return string(bytes)
}

func FormatItemsAsJSONManual(items []Item) string {
	var buf strings.Builder
	buf.Grow(len(items) * 100)

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
// Exercise 6: Distance Calculation (naive + optimized)
// ============================================================================

func ComputeDistances(points [][2]float64) []float64 {
	var out []float64
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dx := points[i][0] - points[j][0]
			dy := points[i][1] - points[j][1]
			out = append(out, math.Sqrt(dx*dx+dy*dy))
		}
	}
	return out
}

func ComputeDistancesOptimized(points [][2]float64) []float64 {
	n := len(points)
	size := n * (n - 1) / 2
	distances := make([]float64, 0, size)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := points[i][0] - points[j][0]
			dy := points[i][1] - points[j][1]
			distances = append(distances, math.Sqrt(dx*dx+dy*dy))
		}
	}
	return distances
}

// ============================================================================
// Exercise 7: Word Frequency (naive + optimized)
// ============================================================================

func CountWordFrequency(docs []Document) map[string]int {
	out := make(map[string]int)
	for _, doc := range docs {
		words := strings.Split(strings.ToLower(doc.Content), " ")
		for _, w := range words {
			w = strings.Trim(w, ".,!?;:")
			if w == "" {
				continue
			}
			out[w]++
		}
	}
	return out
}

func CountWordFrequencyOptimized(docs []Document) map[string]int {
	wordCount := make(map[string]int, 1000)
	for _, doc := range docs {
		words := strings.Fields(strings.ToLower(doc.Content))
		for _, word := range words {
			word = strings.Trim(word, ".,!?;:")
			if word != "" {
				wordCount[word]++
			}
		}
	}
	return wordCount
}

// ============================================================================
// Exercise 8: Cache (simple + optimized)
// ============================================================================

type SimpleCache struct {
	mu       sync.RWMutex
	data     map[string]interface{}
	capacity int
}

func NewSimpleCache(capacity int) *SimpleCache {
	if capacity <= 0 {
		capacity = 1
	}
	return &SimpleCache{
		data:     make(map[string]interface{}, capacity),
		capacity: capacity,
	}
}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *SimpleCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Very simple eviction: if we're at capacity and this is a new key, delete an arbitrary key.
	if len(c.data) >= c.capacity {
		if _, exists := c.data[key]; !exists {
			for k := range c.data {
				delete(c.data, k)
				break
			}
		}
	}
	c.data[key] = value
}

// OptimizedCache is used by the tests; simple implementation with RWMutex + FIFO eviction.
type OptimizedCache struct {
	mu       sync.RWMutex
	data     map[string]*cacheEntry
	capacity int
	order    []string
}

type cacheEntry struct {
	value     interface{}
	createdAt int64
}

func NewOptimizedCache(capacity int) *OptimizedCache {
	if capacity <= 0 {
		capacity = 1
	}
	return &OptimizedCache{
		data:     make(map[string]*cacheEntry, capacity),
		capacity: capacity,
		order:    make([]string, 0, capacity),
	}
}

func (c *OptimizedCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if entry, ok := c.data[key]; ok {
		return entry.value, true
	}
	return nil, false
}

func (c *OptimizedCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.data) >= c.capacity && c.data[key] == nil {
		if len(c.order) > 0 {
			oldest := c.order[0]
			delete(c.data, oldest)
			c.order = c.order[1:]
		}
	}

	c.data[key] = &cacheEntry{value: value, createdAt: 0}
	c.order = append(c.order, key)
}

func (c *OptimizedCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// ============================================================================
// Exercise 9: Filter + Transform
// ============================================================================

func FilterAndTransform(items []Item, minValue float64) []Result {
	var out []Result
	for _, item := range items {
		if item.Value >= minValue {
			out = append(out, Result{
				ItemID:   item.ID,
				Score:    item.Value * 2,
				Category: "filtered",
			})
		}
	}
	return out
}

func FilterAndTransformOptimized(items []Item, minValue float64) []Result {
	results := make([]Result, 0, len(items)/2)
	for _, item := range items {
		if item.Value >= minValue {
			results = append(results, Result{
				ItemID:   item.ID,
				Score:    item.Value * 2,
				Category: "filtered",
			})
		}
	}
	return results
}

// ============================================================================
// Exercise 10: Fibonacci (naive + optimized)
// ============================================================================

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

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

func FibonacciMemoized(n int) int {
	memo := make(map[int]int)
	return fibMemo(n, memo)
}

func fibMemo(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	if v, ok := memo[n]; ok {
		return v
	}
	v := fibMemo(n-1, memo) + fibMemo(n-2, memo)
	memo[n] = v
	return v
}

func FibonacciMatrix(n int) int {
	if n <= 1 {
		return n
	}
	// Matrix: [[1,1],[1,0]]^(n-1) gives Fibonacci(n) at [0][0]
	m := matrixPower([][]int{{1, 1}, {1, 0}}, n-1)
	return m[0][0]
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
// Test data generators
// ============================================================================

func GenerateTestItems(n int) []Item {
	items := make([]Item, n)
	for i := 0; i < n; i++ {
		items[i] = Item{
			ID:        i + 1,
			Name:      fmt.Sprintf("Item%d", i+1),
			Value:     float64((i * 7) % 100),
			Timestamp: int64(i),
			Tags:      []string{"test", "profiling"},
			Metadata:  map[string]string{"k": "v"},
		}
	}
	return items
}

func GenerateTestDocuments(n int) []Document {
	docs := make([]Document, n)
	for i := 0; i < n; i++ {
		docs[i] = Document{
			ID:      i,
			Title:   fmt.Sprintf("Go Document %d", i),
			Content: "the quick brown fox jumps over the lazy dog profiling Go performance",
			Author:  "tester",
			Tags:    []string{"go", "profiling"},
		}
	}
	return docs
}

func GenerateTestPoints(n int) [][2]float64 {
	r := rand.New(rand.NewSource(1))
	points := make([][2]float64, n)
	for i := 0; i < n; i++ {
		points[i] = [2]float64{r.Float64() * 100, r.Float64() * 100}
	}
	return points
}

func init() {
	// Ensure deterministic generators that rely on time.
	rand.Seed(1)
	_ = time.Now() // keep import used if expanded later
}
