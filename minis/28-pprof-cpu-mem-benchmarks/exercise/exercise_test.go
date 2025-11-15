package exercise

import (
	"testing"
)

// ============================================================================
// Tests for Exercise 1: Prime Number Finding
// ============================================================================

func TestFindPrimes(t *testing.T) {
	tests := []struct {
		n        int
		expected []int
	}{
		{10, []int{2, 3, 5, 7}},
		{20, []int{2, 3, 5, 7, 11, 13, 17, 19}},
		{2, []int{2}},
		{1, []int{}},
	}

	for _, tt := range tests {
		result := FindPrimes(tt.n)
		if !equalIntSlices(result, tt.expected) {
			t.Errorf("FindPrimes(%d) = %v, want %v", tt.n, result, tt.expected)
		}
	}
}

func TestFindPrimesOptimized(t *testing.T) {
	tests := []struct {
		n        int
		expected []int
	}{
		{10, []int{2, 3, 5, 7}},
		{20, []int{2, 3, 5, 7, 11, 13, 17, 19}},
		{2, []int{2}},
		{1, []int{}},
	}

	for _, tt := range tests {
		result := FindPrimesOptimized(tt.n)
		if !equalIntSlices(result, tt.expected) {
			t.Errorf("FindPrimesOptimized(%d) = %v, want %v", tt.n, result, tt.expected)
		}
	}
}

// Benchmark naive vs optimized prime finding
func BenchmarkFindPrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimes(1000)
	}
}

func BenchmarkFindPrimesOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimesOptimized(1000)
	}
}

// ============================================================================
// Tests for Exercise 2: String Building
// ============================================================================

func TestBuildReport(t *testing.T) {
	items := []Item{
		{ID: 1, Name: "Item1", Value: 10.5},
		{ID: 2, Name: "Item2", Value: 20.3},
	}

	result := BuildReport(items)

	// Check that result contains expected content
	if result == "" {
		t.Error("BuildReport returned empty string")
	}
	if len(result) < 50 {
		t.Errorf("BuildReport result too short: %d chars", len(result))
	}
}

func TestBuildReportOptimized(t *testing.T) {
	items := []Item{
		{ID: 1, Name: "Item1", Value: 10.5},
		{ID: 2, Name: "Item2", Value: 20.3},
	}

	result := BuildReportOptimized(items)

	if result == "" {
		t.Error("BuildReportOptimized returned empty string")
	}
	if len(result) < 50 {
		t.Errorf("BuildReportOptimized result too short: %d chars", len(result))
	}
}

// Benchmark string building
func BenchmarkBuildReport(b *testing.B) {
	items := GenerateTestItems(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = BuildReport(items)
	}
}

func BenchmarkBuildReportOptimized(b *testing.B) {
	items := GenerateTestItems(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = BuildReportOptimized(items)
	}
}

func BenchmarkBuildReportWithPool(b *testing.B) {
	items := GenerateTestItems(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = BuildReportWithPool(items)
	}
}

// ============================================================================
// Tests for Exercise 3: Document Search
// ============================================================================

func TestSearchDocuments(t *testing.T) {
	docs := GenerateTestDocuments(5)
	results := SearchDocuments(docs, "Go")

	if len(results) == 0 {
		t.Error("SearchDocuments found no results for 'Go'")
	}
}

func TestSearchDocumentsOptimized(t *testing.T) {
	docs := GenerateTestDocuments(5)
	results := SearchDocumentsOptimized(docs, "Go")

	if len(results) == 0 {
		t.Error("SearchDocumentsOptimized found no results for 'Go'")
	}
}

// Benchmark search
func BenchmarkSearchDocuments(b *testing.B) {
	docs := GenerateTestDocuments(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = SearchDocuments(docs, "profiling")
	}
}

func BenchmarkSearchDocumentsOptimized(b *testing.B) {
	docs := GenerateTestDocuments(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = SearchDocumentsOptimized(docs, "profiling")
	}
}

// ============================================================================
// Tests for Exercise 4: Process Items
// ============================================================================

func TestProcessItems(t *testing.T) {
	items := GenerateTestItems(10)
	results := ProcessItems(items)

	if len(results) != len(items) {
		t.Errorf("ProcessItems returned %d results, expected %d", len(results), len(items))
	}

	for i, result := range results {
		if result.ItemID != items[i].ID {
			t.Errorf("Result %d has wrong ID: got %d, want %d", i, result.ItemID, items[i].ID)
		}
	}
}

func TestProcessItemsOptimized(t *testing.T) {
	items := GenerateTestItems(10)
	results := ProcessItemsOptimized(items)

	if len(results) != len(items) {
		t.Errorf("ProcessItemsOptimized returned %d results, expected %d", len(results), len(items))
	}
}

// Benchmark item processing
func BenchmarkProcessItems(b *testing.B) {
	items := GenerateTestItems(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ProcessItems(items)
	}
}

func BenchmarkProcessItemsOptimized(b *testing.B) {
	items := GenerateTestItems(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ProcessItemsOptimized(items)
	}
}

// ============================================================================
// Tests for Exercise 5: JSON Formatting
// ============================================================================

func TestFormatItemsAsJSON(t *testing.T) {
	items := GenerateTestItems(3)
	result := FormatItemsAsJSON(items)

	if result == "" {
		t.Error("FormatItemsAsJSON returned empty string")
	}
	if len(result) < 20 {
		t.Errorf("FormatItemsAsJSON result too short: %d chars", len(result))
	}
}

func TestFormatItemsAsJSONOptimized(t *testing.T) {
	items := GenerateTestItems(3)
	result := FormatItemsAsJSONOptimized(items)

	if result == "" {
		t.Error("FormatItemsAsJSONOptimized returned empty string")
	}
}

// Benchmark JSON formatting
func BenchmarkFormatItemsAsJSON(b *testing.B) {
	items := GenerateTestItems(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FormatItemsAsJSON(items)
	}
}

func BenchmarkFormatItemsAsJSONOptimized(b *testing.B) {
	items := GenerateTestItems(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FormatItemsAsJSONOptimized(items)
	}
}

func BenchmarkFormatItemsAsJSONManual(b *testing.B) {
	items := GenerateTestItems(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FormatItemsAsJSONManual(items)
	}
}

// ============================================================================
// Tests for Exercise 6: Distance Calculation
// ============================================================================

func TestComputeDistances(t *testing.T) {
	points := GenerateTestPoints(5)
	distances := ComputeDistances(points)

	expectedCount := 5 * 4 / 2 // n*(n-1)/2
	if len(distances) != expectedCount {
		t.Errorf("ComputeDistances returned %d distances, expected %d", len(distances), expectedCount)
	}
}

func TestComputeDistancesOptimized(t *testing.T) {
	points := GenerateTestPoints(5)
	distances := ComputeDistancesOptimized(points)

	expectedCount := 5 * 4 / 2
	if len(distances) != expectedCount {
		t.Errorf("ComputeDistancesOptimized returned %d distances, expected %d", len(distances), expectedCount)
	}
}

// Benchmark distance calculations
func BenchmarkComputeDistances(b *testing.B) {
	points := GenerateTestPoints(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ComputeDistances(points)
	}
}

func BenchmarkComputeDistancesOptimized(b *testing.B) {
	points := GenerateTestPoints(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ComputeDistancesOptimized(points)
	}
}

// ============================================================================
// Tests for Exercise 7: Word Frequency
// ============================================================================

func TestCountWordFrequency(t *testing.T) {
	docs := GenerateTestDocuments(5)
	freq := CountWordFrequency(docs)

	if len(freq) == 0 {
		t.Error("CountWordFrequency returned empty map")
	}

	// Check that common words are present
	if freq["the"] == 0 {
		t.Error("Expected word 'the' to be present")
	}
}

func TestCountWordFrequencyOptimized(t *testing.T) {
	docs := GenerateTestDocuments(5)
	freq := CountWordFrequencyOptimized(docs)

	if len(freq) == 0 {
		t.Error("CountWordFrequencyOptimized returned empty map")
	}
}

// Benchmark word frequency
func BenchmarkCountWordFrequency(b *testing.B) {
	docs := GenerateTestDocuments(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = CountWordFrequency(docs)
	}
}

func BenchmarkCountWordFrequencyOptimized(b *testing.B) {
	docs := GenerateTestDocuments(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = CountWordFrequencyOptimized(docs)
	}
}

// ============================================================================
// Tests for Exercise 8: Cache
// ============================================================================

func TestSimpleCache(t *testing.T) {
	cache := NewSimpleCache(100)

	// Test Set and Get
	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")

	if !ok {
		t.Error("Expected key1 to be present in cache")
	}

	if val != "value1" {
		t.Errorf("Got value %v, expected 'value1'", val)
	}

	// Test missing key
	_, ok = cache.Get("missing")
	if ok {
		t.Error("Expected missing key to not be found")
	}
}

func TestOptimizedCache(t *testing.T) {
	cache := NewOptimizedCache(100)

	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")

	if !ok {
		t.Error("Expected key1 to be present in cache")
	}

	if val != "value1" {
		t.Errorf("Got value %v, expected 'value1'", val)
	}
}

// Benchmark cache operations
func BenchmarkCacheSet(b *testing.B) {
	cache := NewSimpleCache(10000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key", i)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewSimpleCache(10000)
	cache.Set("key", "value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("key")
	}
}

func BenchmarkOptimizedCacheSet(b *testing.B) {
	cache := NewOptimizedCache(10000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key", i)
	}
}

func BenchmarkOptimizedCacheGet(b *testing.B) {
	cache := NewOptimizedCache(10000)
	cache.Set("key", "value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("key")
	}
}

// ============================================================================
// Tests for Exercise 9: Filter and Transform
// ============================================================================

func TestFilterAndTransform(t *testing.T) {
	items := GenerateTestItems(100)
	results := FilterAndTransform(items, 50.0)

	// Check that all results have value >= 50
	for _, result := range results {
		if result.Score < 100.0 { // Score is value * 2
			t.Errorf("Result has score %.2f, expected >= 100.0", result.Score)
		}
	}
}

func TestFilterAndTransformOptimized(t *testing.T) {
	items := GenerateTestItems(100)
	results := FilterAndTransformOptimized(items, 50.0)

	if len(results) == 0 {
		t.Error("FilterAndTransformOptimized returned no results")
	}
}

// Benchmark filter operations
func BenchmarkFilterAndTransform(b *testing.B) {
	items := GenerateTestItems(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FilterAndTransform(items, 50.0)
	}
}

func BenchmarkFilterAndTransformOptimized(b *testing.B) {
	items := GenerateTestItems(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FilterAndTransformOptimized(items, 50.0)
	}
}

// ============================================================================
// Tests for Exercise 10: Fibonacci
// ============================================================================

func TestFibonacci(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
		{10, 55},
	}

	for _, tt := range tests {
		// Test naive (only for small values)
		if tt.n <= 10 {
			result := Fibonacci(tt.n)
			if result != tt.expected {
				t.Errorf("Fibonacci(%d) = %d, want %d", tt.n, result, tt.expected)
			}
		}

		// Test optimized versions
		result := FibonacciIterative(tt.n)
		if result != tt.expected {
			t.Errorf("FibonacciIterative(%d) = %d, want %d", tt.n, result, tt.expected)
		}

		result = FibonacciMemoized(tt.n)
		if result != tt.expected {
			t.Errorf("FibonacciMemoized(%d) = %d, want %d", tt.n, result, tt.expected)
		}

		result = FibonacciMatrix(tt.n)
		if result != tt.expected {
			t.Errorf("FibonacciMatrix(%d) = %d, want %d", tt.n, result, tt.expected)
		}
	}
}

// Benchmark Fibonacci implementations
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20) // Limited to 20 due to exponential complexity
	}
}

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(20)
	}
}

func BenchmarkFibonacciMemoized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciMemoized(20)
	}
}

func BenchmarkFibonacciMatrix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciMatrix(20)
	}
}

// ============================================================================
// Benchmark Variations: Testing Different Sizes
// ============================================================================

func BenchmarkFindPrimes_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimes(100)
	}
}

func BenchmarkFindPrimes_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimes(1000)
	}
}

func BenchmarkFindPrimes_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimes(10000)
	}
}

func BenchmarkFindPrimesOptimized_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimesOptimized(100)
	}
}

func BenchmarkFindPrimesOptimized_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimesOptimized(1000)
	}
}

func BenchmarkFindPrimesOptimized_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPrimesOptimized(10000)
	}
}

// ============================================================================
// Memory Allocation Benchmarks
// ============================================================================

func BenchmarkAllocations_BuildReport(b *testing.B) {
	items := GenerateTestItems(100)
	b.ReportAllocs() // Show allocation stats
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = BuildReport(items)
	}
}

func BenchmarkAllocations_BuildReportOptimized(b *testing.B) {
	items := GenerateTestItems(100)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = BuildReportOptimized(items)
	}
}

func BenchmarkAllocations_ProcessItems(b *testing.B) {
	items := GenerateTestItems(1000)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ProcessItems(items)
	}
}

func BenchmarkAllocations_ProcessItemsOptimized(b *testing.B) {
	items := GenerateTestItems(1000)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ProcessItemsOptimized(items)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ============================================================================
// Example: How to Run Benchmarks with Profiling
// ============================================================================

// To run benchmarks with CPU profiling:
//   go test -bench=. -cpuprofile=cpu.prof
//   go tool pprof cpu.prof
//
// To run benchmarks with memory profiling:
//   go test -bench=. -memprofile=mem.prof
//   go tool pprof mem.prof
//
// To run benchmarks with allocation stats:
//   go test -bench=. -benchmem
//
// To compare before/after optimizations:
//   go test -bench=BenchmarkFindPrimes -count=10 > old.txt
//   (make changes)
//   go test -bench=BenchmarkFindPrimesOptimized -count=10 > new.txt
//   benchstat old.txt new.txt
//
// To run with race detector:
//   go test -race -bench=.
//
// To profile specific benchmarks:
//   go test -bench=BenchmarkFindPrimesOptimized -cpuprofile=cpu.prof
//   go tool pprof -http=:8080 cpu.prof
