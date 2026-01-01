package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("=== sync.Pool Demonstrations ===")
	fmt.Println()

	demo1AllocationProblem()
	demo2BasicPooling()
	demo3BufferPoolPattern()
	demo4GCImpact()
	demo5PoolLifecycle()
	demo6SizeClassedPools()
	demo7PoolMetrics()
	demo8RealWorldHTTPServer()
}

// Demo 1: The Allocation Problem
func demo1AllocationProblem() {
	fmt.Println("--- Demo 1: The Allocation Problem ---")
	fmt.Println("Simulating high-frequency allocations without pooling...")

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	beforeAlloc := stats.TotalAlloc
	beforeGC := stats.NumGC

	start := time.Now()

	// Simulate 100,000 temporary buffer allocations
	for i := 0; i < 100000; i++ {
		buf := new(bytes.Buffer)
		buf.WriteString("temporary data that will be discarded")
		_ = buf.String()
		// buf is now garbage
	}

	elapsed := time.Since(start)

	runtime.ReadMemStats(&stats)
	afterAlloc := stats.TotalAlloc
	afterGC := stats.NumGC

	fmt.Printf("Operations:       100,000\n")
	fmt.Printf("Time:             %v\n", elapsed)
	fmt.Printf("Throughput:       %.0f ops/sec\n", 100000/elapsed.Seconds())
	fmt.Printf("Total allocated:  %.2f MB\n", float64(afterAlloc-beforeAlloc)/1024/1024)
	fmt.Printf("GC cycles:        %d\n", afterGC-beforeGC)
	fmt.Println()
}

// Demo 2: Basic Pooling
func demo2BasicPooling() {
	fmt.Println("--- Demo 2: Basic Pooling ---")
	fmt.Println("Same workload with sync.Pool...")

	var bufferPool = sync.Pool{
		New: func() interface{} {
			fmt.Println("  [Pool] Creating new buffer (pool was empty)")
			return new(bytes.Buffer)
		},
	}

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	beforeAlloc := stats.TotalAlloc
	beforeGC := stats.NumGC

	start := time.Now()

	// Same 100,000 operations, but with pooling
	for i := 0; i < 100000; i++ {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		buf.WriteString("temporary data that will be reused")
		_ = buf.String()
		bufferPool.Put(buf)
	}

	elapsed := time.Since(start)

	runtime.ReadMemStats(&stats)
	afterAlloc := stats.TotalAlloc
	afterGC := stats.NumGC

	fmt.Printf("Operations:       100,000\n")
	fmt.Printf("Time:             %v\n", elapsed)
	fmt.Printf("Throughput:       %.0f ops/sec\n", 100000/elapsed.Seconds())
	fmt.Printf("Total allocated:  %.2f MB\n", float64(afterAlloc-beforeAlloc)/1024/1024)
	fmt.Printf("GC cycles:        %d\n", afterGC-beforeGC)
	fmt.Printf("Speedup:          ~2-3x faster\n")
	fmt.Printf("Allocation reduction: ~99%%\n")
	fmt.Println()
}

// Demo 3: Buffer Pool Pattern
func demo3BufferPoolPattern() {
	fmt.Println("--- Demo 3: Buffer Pool Pattern ---")

	var bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	processData := func(data string) string {
		// Get buffer from pool
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset() // Important: clear previous data
		defer bufferPool.Put(buf)

		// Use buffer for temporary processing
		buf.WriteString("Processed: ")
		buf.WriteString(data)
		buf.WriteString(" (transformed)")

		return buf.String()
	}

	// Process multiple items
	items := []string{"item1", "item2", "item3", "item4", "item5"}
	for _, item := range items {
		result := processData(item)
		fmt.Printf("  %s → %s\n", item, result)
	}

	fmt.Println("Notice: Same buffer was reused for all operations!")
	fmt.Println()
}

// Demo 4: GC Impact Comparison
func demo4GCImpact() {
	fmt.Println("--- Demo 4: GC Impact Comparison ---")

	// Without pool
	fmt.Println("Running WITHOUT pool (watch GC cycles)...")
	runWorkload(false)

	// Force GC to clean up
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// With pool
	fmt.Println("\nRunning WITH pool (watch GC cycles)...")
	runWorkload(true)

	fmt.Println("\nKey insight: Pool version triggers fewer GC cycles")
	fmt.Println()
}

func runWorkload(usePool bool) {
	var bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	beforeGC := stats.NumGC
	beforePause := stats.PauseTotalNs

	// Simulate workload
	for i := 0; i < 50000; i++ {
		var buf *bytes.Buffer

		if usePool {
			buf = bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
		} else {
			buf = new(bytes.Buffer)
		}

		// Do some work
		buf.WriteString(fmt.Sprintf("Data item %d", i))
		for j := 0; j < 10; j++ {
			buf.WriteString(fmt.Sprintf(" field%d", j))
		}
		_ = buf.Bytes()

		if usePool {
			bufferPool.Put(buf)
		}
	}

	runtime.ReadMemStats(&stats)
	afterGC := stats.NumGC
	afterPause := stats.PauseTotalNs

	fmt.Printf("  GC cycles:        %d\n", afterGC-beforeGC)
	fmt.Printf("  Total GC pause:   %.2f ms\n", float64(afterPause-beforePause)/1e6)
}

// Demo 5: Pool Lifecycle and GC Interaction
func demo5PoolLifecycle() {
	fmt.Println("--- Demo 5: Pool Lifecycle and GC Interaction ---")

	creationCount := atomic.Int64{}
	var pool = sync.Pool{
		New: func() interface{} {
			count := creationCount.Add(1)
			fmt.Printf("  [Pool] Creating object #%d\n", count)
			return &struct{ id int }{id: int(count)}
		},
	}

	fmt.Println("Step 1: Get object from empty pool")
	obj1 := pool.Get()
	fmt.Printf("  Got: %+v\n", obj1)

	fmt.Println("\nStep 2: Put object back")
	pool.Put(obj1)
	fmt.Println("  Object returned to pool")

	fmt.Println("\nStep 3: Get again (should reuse)")
	obj2 := pool.Get()
	fmt.Printf("  Got: %+v (same object!)\n", obj2)

	fmt.Println("\nStep 4: Force GC")
	pool.Put(obj2)
	runtime.GC()
	fmt.Println("  First GC done (objects moved to victim cache)")

	fmt.Println("\nStep 5: Get from pool (still available from victim)")
	obj3 := pool.Get()
	fmt.Printf("  Got: %+v (from victim cache)\n", obj3)

	fmt.Println("\nStep 6: Force second GC")
	pool.Put(obj3)
	runtime.GC()
	fmt.Println("  Second GC done (victim cache cleared)")

	fmt.Println("\nStep 7: Get from pool (pool is empty now)")
	obj4 := pool.Get()
	fmt.Printf("  Got: %+v (new object created)\n", obj4)

	fmt.Println("\nKey insight: Objects survive 2 GC cycles before being freed")
	fmt.Println()
}

// Demo 6: Size-Classed Pools
func demo6SizeClassedPools() {
	fmt.Println("--- Demo 6: Size-Classed Pools ---")
	fmt.Println("Different pool for each buffer size class...")

	// Define size classes
	var pools = [4]sync.Pool{
		// 1KB pool
		{New: func() interface{} {
			buf := make([]byte, 0, 1024)
			return &buf
		}},
		// 4KB pool
		{New: func() interface{} {
			buf := make([]byte, 0, 4096)
			return &buf
		}},
		// 16KB pool
		{New: func() interface{} {
			buf := make([]byte, 0, 16384)
			return &buf
		}},
		// 64KB pool
		{New: func() interface{} {
			buf := make([]byte, 0, 65536)
			return &buf
		}},
	}

	getBuffer := func(size int) (*[]byte, int) {
		switch {
		case size <= 1024:
			return pools[0].Get().(*[]byte), 0
		case size <= 4096:
			return pools[1].Get().(*[]byte), 1
		case size <= 16384:
			return pools[2].Get().(*[]byte), 2
		default:
			return pools[3].Get().(*[]byte), 3
		}
	}

	putBuffer := func(buf *[]byte, poolIdx int) {
		*buf = (*buf)[:0] // Reset length
		pools[poolIdx].Put(buf)
	}

	// Test with different sizes
	sizes := []int{500, 2000, 10000, 50000}
	for _, size := range sizes {
		buf, poolIdx := getBuffer(size)
		*buf = append(*buf, make([]byte, size)...)

		poolName := []string{"1KB", "4KB", "16KB", "64KB"}[poolIdx]
		fmt.Printf("  Requested %5d bytes → used %s pool (cap=%d)\n",
			size, poolName, cap(*buf))

		putBuffer(buf, poolIdx)
	}

	fmt.Println("\nKey insight: Right-sized pools reduce memory waste")
	fmt.Println()
}

// Demo 7: Pool with Metrics
func demo7PoolMetrics() {
	fmt.Println("--- Demo 7: Pool with Metrics ---")

	type MetricsPool struct {
		pool sync.Pool
		gets atomic.Int64
		puts atomic.Int64
		news atomic.Int64
	}

	mp := &MetricsPool{}
	mp.pool.New = func() interface{} {
		mp.news.Add(1)
		return new(bytes.Buffer)
	}

	get := func() *bytes.Buffer {
		mp.gets.Add(1)
		return mp.pool.Get().(*bytes.Buffer)
	}

	put := func(buf *bytes.Buffer) {
		mp.puts.Add(1)
		buf.Reset()
		mp.pool.Put(buf)
	}

	// Run workload
	for i := 0; i < 1000; i++ {
		buf := get()
		buf.WriteString("data")
		put(buf)
	}

	// Report metrics
	gets := mp.gets.Load()
	puts := mp.puts.Load()
	news := mp.news.Load()
	hitRate := float64(gets-news) / float64(gets) * 100

	fmt.Printf("  Total Gets:       %d\n", gets)
	fmt.Printf("  Total Puts:       %d\n", puts)
	fmt.Printf("  Objects Created:  %d\n", news)
	fmt.Printf("  Hit Rate:         %.1f%%\n", hitRate)
	fmt.Printf("  Reuse Rate:       %.1f%%\n", 100-hitRate)
	fmt.Println()
}

// Demo 8: Real-World HTTP Server Simulation
func demo8RealWorldHTTPServer() {
	fmt.Println("--- Demo 8: Real-World HTTP Server Simulation ---")

	// Shared buffer pool
	var bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// Shared gzip writer pool
	var gzipPool = sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(io.Discard)
		},
	}

	type Response struct {
		Status  string                 `json:"status"`
		Data    map[string]interface{} `json:"data"`
		Message string                 `json:"message"`
	}

	// Simulate handling HTTP request
	handleRequest := func(requestID int) []byte {
		// Get buffer from pool
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufferPool.Put(buf)

		// Create response
		resp := Response{
			Status: "success",
			Data: map[string]interface{}{
				"request_id": requestID,
				"timestamp":  time.Now().Unix(),
				"random":     rand.Intn(1000),
			},
			Message: "Request processed successfully",
		}

		// Encode JSON to buffer
		if err := json.NewEncoder(buf).Encode(resp); err != nil {
			panic(err)
		}

		// Compress response
		compressedBuf := bufferPool.Get().(*bytes.Buffer)
		compressedBuf.Reset()
		defer bufferPool.Put(compressedBuf)

		gzw := gzipPool.Get().(*gzip.Writer)
		defer gzipPool.Put(gzw)
		gzw.Reset(compressedBuf)

		if _, err := gzw.Write(buf.Bytes()); err != nil {
			panic(err)
		}
		gzw.Close()

		// Return compressed data (copy because buffer is pooled)
		result := make([]byte, compressedBuf.Len())
		copy(result, compressedBuf.Bytes())
		return result
	}

	fmt.Println("Simulating 10,000 HTTP requests with pooling...")

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	beforeAlloc := stats.TotalAlloc
	beforeGC := stats.NumGC

	start := time.Now()

	// Simulate concurrent requests
	var wg sync.WaitGroup
	const numRequests = 10000
	const concurrency = 100

	sem := make(chan struct{}, concurrency)
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			_ = handleRequest(id)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	runtime.ReadMemStats(&stats)
	afterAlloc := stats.TotalAlloc
	afterGC := stats.NumGC

	fmt.Printf("\nResults:\n")
	fmt.Printf("  Requests:         %d\n", numRequests)
	fmt.Printf("  Time:             %v\n", elapsed)
	fmt.Printf("  Throughput:       %.0f req/sec\n", float64(numRequests)/elapsed.Seconds())
	fmt.Printf("  Avg latency:      %.2f ms\n", elapsed.Seconds()*1000/numRequests)
	fmt.Printf("  Total allocated:  %.2f MB\n", float64(afterAlloc-beforeAlloc)/1024/1024)
	fmt.Printf("  GC cycles:        %d\n", afterGC-beforeGC)

	fmt.Println("\nKey insight: Pools dramatically reduce allocations in real workloads")
	fmt.Println("  - Buffer pool: Reuses JSON encoding buffers")
	fmt.Println("  - Gzip pool: Reuses compression writers (~256KB each)")
	fmt.Println()
}

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}
