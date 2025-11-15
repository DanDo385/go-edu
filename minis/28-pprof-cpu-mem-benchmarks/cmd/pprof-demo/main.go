package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof" // Registers /debug/pprof/* handlers
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== pprof Profiling Demonstration ===")
	fmt.Println()
	fmt.Println("Starting HTTP server on :6060")
	fmt.Println("Profile endpoints available at:")
	fmt.Println("  http://localhost:6060/debug/pprof/")
	fmt.Println()
	fmt.Println("Available profiles:")
	fmt.Println("  - CPU:       /debug/pprof/profile?seconds=30")
	fmt.Println("  - Heap:      /debug/pprof/heap")
	fmt.Println("  - Goroutine: /debug/pprof/goroutine")
	fmt.Println("  - Block:     /debug/pprof/block")
	fmt.Println("  - Mutex:     /debug/pprof/mutex")
	fmt.Println("  - Allocs:    /debug/pprof/allocs")
	fmt.Println()
	fmt.Println("Usage examples:")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/heap")
	fmt.Println("  go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=10")
	fmt.Println()

	// Register custom HTTP handlers for workload demos
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/cpu-demo", cpuDemoHandler)
	http.HandleFunc("/mem-demo", memDemoHandler)
	http.HandleFunc("/goroutine-demo", goroutineDemoHandler)
	http.HandleFunc("/stats", statsHandler)

	// Start background workload generators
	fmt.Println("Starting background workload generators...")
	go cpuWorkloadGenerator()
	go memoryWorkloadGenerator()
	go goroutineWorkloadGenerator()

	fmt.Println()
	fmt.Println("Server ready! Press Ctrl+C to stop.")
	fmt.Println()

	// Start HTTP server
	if err := http.ListenAndServe(":6060", nil); err != nil {
		log.Fatal(err)
	}
}

// ============================================================================
// HTTP Handlers
// ============================================================================

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>pprof Demo Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        .section { margin: 20px 0; padding: 15px; background: #f5f5f5; border-radius: 5px; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>pprof Profiling Demo Server</h1>

    <div class="section">
        <h2>🔍 Profile Endpoints</h2>
        <ul>
            <li><a href="/debug/pprof/">/debug/pprof/</a> - Profile index</li>
            <li><a href="/debug/pprof/heap">/debug/pprof/heap</a> - Heap memory profile</li>
            <li><a href="/debug/pprof/goroutine">/debug/pprof/goroutine</a> - Goroutine stack traces</li>
            <li><a href="/debug/pprof/profile?seconds=10">/debug/pprof/profile?seconds=10</a> - 10-second CPU profile</li>
            <li><a href="/debug/pprof/block">/debug/pprof/block</a> - Blocking profile</li>
            <li><a href="/debug/pprof/mutex">/debug/pprof/mutex</a> - Mutex contention</li>
        </ul>
    </div>

    <div class="section">
        <h2>🚀 Demo Workloads</h2>
        <ul>
            <li><a href="/cpu-demo">/cpu-demo</a> - Trigger CPU-intensive work</li>
            <li><a href="/mem-demo">/mem-demo</a> - Trigger memory allocations</li>
            <li><a href="/goroutine-demo">/goroutine-demo</a> - Spawn goroutines</li>
            <li><a href="/stats">/stats</a> - Runtime statistics (JSON)</li>
        </ul>
    </div>

    <div class="section">
        <h2>📖 CLI Usage</h2>
        <p>Collect and analyze profiles using <code>go tool pprof</code>:</p>
        <pre>
# Interactive CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Interactive memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Web UI with flamegraphs
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=10

# Save profile to file
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof
        </pre>
    </div>

    <div class="section">
        <h2>💡 Common Commands</h2>
        <p>Inside <code>pprof</code> interactive mode:</p>
        <ul>
            <li><code>top</code> - Show top functions by flat time</li>
            <li><code>top -cum</code> - Show top functions by cumulative time</li>
            <li><code>list funcName</code> - Show annotated source code</li>
            <li><code>web</code> - Generate call graph (requires graphviz)</li>
            <li><code>pdf > profile.pdf</code> - Export to PDF</li>
            <li><code>help</code> - Show all commands</li>
        </ul>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func cpuDemoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Running CPU-intensive workload...\n\n")

	start := time.Now()

	// Compute prime numbers (inefficient algorithm for profiling demo)
	primes := findPrimesNaive(10000)
	fmt.Fprintf(w, "Found %d prime numbers\n", len(primes))

	// Hash computation
	data := make([]byte, 1024*1024) // 1MB
	rand.Read(data)
	hash := sha256.Sum256(data)
	fmt.Fprintf(w, "SHA256 hash: %x\n", hash[:8])

	// String operations
	result := buildLargeString(1000)
	fmt.Fprintf(w, "Built string of length: %d\n", len(result))

	elapsed := time.Since(start)
	fmt.Fprintf(w, "\nCompleted in %v\n", elapsed)
	fmt.Fprintf(w, "\nNow try: go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10\n")
}

func memDemoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Running memory-intensive workload...\n\n")

	start := time.Now()

	// Allocate large data structures
	slices := allocateLargeSlices(100)
	fmt.Fprintf(w, "Allocated %d slices\n", len(slices))

	// JSON marshaling (allocates a lot)
	data := generateComplexStructure(1000)
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "JSON size: %d bytes\n", len(jsonBytes))

	// String concatenation (inefficient, lots of allocations)
	concat := concatenateStringsInefficient(500)
	fmt.Fprintf(w, "Concatenated string length: %d\n", len(concat))

	elapsed := time.Since(start)
	fmt.Fprintf(w, "\nCompleted in %v\n", elapsed)
	fmt.Fprintf(w, "\nNow try: go tool pprof http://localhost:6060/debug/pprof/heap\n")
}

func goroutineDemoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Spawning goroutines...\n\n")

	before := runtime.NumGoroutine()
	fmt.Fprintf(w, "Goroutines before: %d\n", before)

	// Spawn temporary goroutines
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}(i)
	}

	after := runtime.NumGoroutine()
	fmt.Fprintf(w, "Goroutines after spawn: %d\n", after)

	fmt.Fprintf(w, "\nWaiting for goroutines to complete...\n")
	wg.Wait()

	final := runtime.NumGoroutine()
	fmt.Fprintf(w, "Goroutines after completion: %d\n", final)

	fmt.Fprintf(w, "\nNow try: go tool pprof http://localhost:6060/debug/pprof/goroutine\n")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := map[string]interface{}{
		"goroutines":     runtime.NumGoroutine(),
		"cpus":           runtime.NumCPU(),
		"alloc_mb":       m.Alloc / 1024 / 1024,
		"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
		"sys_mb":         m.Sys / 1024 / 1024,
		"num_gc":         m.NumGC,
		"gc_pause_ns":    m.PauseNs[(m.NumGC+255)%256],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ============================================================================
// Background Workload Generators
// ============================================================================

func cpuWorkloadGenerator() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Simulate CPU-intensive work
		_ = findPrimesNaive(5000)
		_ = sha256.Sum256([]byte(time.Now().String()))
	}
}

func memoryWorkloadGenerator() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Simulate memory allocations
		_ = allocateLargeSlices(10)
		_ = generateComplexStructure(100)

		// Force GC occasionally to see GC in profiles
		if rand.Intn(10) == 0 {
			runtime.GC()
		}
	}
}

func goroutineWorkloadGenerator() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Spawn and wait for goroutines
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			}()
		}
		wg.Wait()
	}
}

// ============================================================================
// CPU-Intensive Functions (for profiling demonstrations)
// ============================================================================

// findPrimesNaive uses a naive O(n²) algorithm - intentionally inefficient for profiling demo
func findPrimesNaive(n int) []int {
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

// buildLargeString inefficiently concatenates strings (allocates a lot)
func buildLargeString(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += fmt.Sprintf("Line %d: some data here\n", i)
	}
	return result
}

// computeFibonacci computes Fibonacci recursively (inefficient)
func computeFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return computeFibonacci(n-1) + computeFibonacci(n-2)
}

// ============================================================================
// Memory-Intensive Functions (for profiling demonstrations)
// ============================================================================

// allocateLargeSlices creates multiple large slices
func allocateLargeSlices(count int) [][]byte {
	slices := make([][]byte, count)
	for i := 0; i < count; i++ {
		slices[i] = make([]byte, 1024*1024) // 1MB each
		rand.Read(slices[i])
	}
	return slices
}

// generateComplexStructure creates nested data structures
func generateComplexStructure(size int) []map[string]interface{} {
	data := make([]map[string]interface{}, size)
	for i := 0; i < size; i++ {
		data[i] = map[string]interface{}{
			"id":        i,
			"name":      fmt.Sprintf("Item-%d", i),
			"timestamp": time.Now().Unix(),
			"data":      make([]byte, 1024), // 1KB per item
			"metadata": map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		}
	}
	return data
}

// concatenateStringsInefficient uses + operator (allocates intermediate strings)
func concatenateStringsInefficient(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += "This is a line of text that will be concatenated. "
	}
	return result
}

// concatenateStringsEfficient uses strings.Builder (efficient)
func concatenateStringsEfficient(n int) string {
	var buf strings.Builder
	buf.Grow(n * 50) // Preallocate
	for i := 0; i < n; i++ {
		buf.WriteString("This is a line of text that will be concatenated. ")
	}
	return buf.String()
}

// ============================================================================
// Demonstration: Before/After Optimization Examples
// ============================================================================

// ExampleBefore shows inefficient prime finding
func ExampleBefore_FindPrimes() {
	primes := findPrimesNaive(1000)
	fmt.Printf("Found %d primes\n", len(primes))
}

// ExampleAfter shows efficient prime finding (Sieve of Eratosthenes)
func ExampleAfter_FindPrimes() {
	primes := findPrimesSieve(1000)
	fmt.Printf("Found %d primes\n", len(primes))
}

// findPrimesSieve uses Sieve of Eratosthenes (O(n log log n))
func findPrimesSieve(n int) []int {
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

	primes := make([]int, 0, n/10) // Approximate prime density
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

// ExampleBefore shows inefficient string building
func ExampleBefore_BuildString() {
	result := ""
	for i := 0; i < 1000; i++ {
		result += fmt.Sprintf("Line %d\n", i)
	}
	fmt.Printf("Built string of length %d\n", len(result))
}

// ExampleAfter shows efficient string building
func ExampleAfter_BuildString() {
	var buf strings.Builder
	buf.Grow(1000 * 20) // Preallocate
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(&buf, "Line %d\n", i)
	}
	result := buf.String()
	fmt.Printf("Built string of length %d\n", len(result))
}

// ExampleBefore shows inefficient byte processing
func ExampleBefore_ProcessBytes() {
	data := make([]byte, 1000000)
	rand.Read(data)

	// Inefficient: Creates many intermediate buffers
	var results [][]byte
	for i := 0; i < len(data); i += 1000 {
		end := i + 1000
		if end > len(data) {
			end = len(data)
		}
		chunk := make([]byte, end-i)
		copy(chunk, data[i:end])
		results = append(results, chunk)
	}

	fmt.Printf("Processed %d chunks\n", len(results))
}

// ExampleAfter shows efficient byte processing with sync.Pool
var byteBufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func ExampleAfter_ProcessBytes() {
	data := make([]byte, 1000000)
	rand.Read(data)

	// Efficient: Reuse buffers
	var results [][]byte
	for i := 0; i < len(data); i += 1000 {
		end := i + 1000
		if end > len(data) {
			end = len(data)
		}

		buf := byteBufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		buf.Write(data[i:end])

		chunk := make([]byte, buf.Len())
		copy(chunk, buf.Bytes())
		results = append(results, chunk)

		byteBufferPool.Put(buf)
	}

	fmt.Printf("Processed %d chunks\n", len(results))
}

// ============================================================================
// Cache Example (for profiling cache hit/miss patterns)
// ============================================================================

type Cache struct {
	mu    sync.RWMutex
	data  map[string][]byte
	hits  int
	misses int
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.data[key]; ok {
		c.hits++
		return val, true
	}
	c.misses++
	return nil, false
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Stats() (hits, misses int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses
}

// ExampleCache demonstrates cache usage under concurrent load
func ExampleCache() {
	cache := NewCache()

	// Populate cache
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := make([]byte, 1024)
		rand.Read(value)
		cache.Set(key, value)
	}

	// Concurrent reads (mix of hits and misses)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", rand.Intn(150)) // Some keys miss
			_, _ = cache.Get(key)
		}(i)
	}
	wg.Wait()

	hits, misses := cache.Stats()
	fmt.Printf("Cache stats: %d hits, %d misses (%.1f%% hit rate)\n",
		hits, misses, float64(hits)*100/float64(hits+misses))
}
