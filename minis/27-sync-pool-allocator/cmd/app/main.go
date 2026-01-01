package main

import (
	"bytes"
	"fmt"
	
	"github.com/example/go-10x-minis/minis/27-sync-pool-allocator/internal/syncpoolallocator"
)

func main() {
	fmt.Println("=== sync.Pool Allocator Examples ===")
	fmt.Println()

	// Example 1: Basic Buffer Pool
	fmt.Println("--- Example 1: Basic Buffer Pool ---")
	bufPool := syncpoolallocator.NewBufferPool()
	if bufPool != nil {
		buf := bufPool.Get()
		if buf != nil {
			buf.WriteString("Hello from buffer pool!")
			fmt.Printf("Buffer content: %s\n", buf.String())
			bufPool.Put(buf)
		}
	}
	fmt.Println()

	// Example 2: Slice Pool
	fmt.Println("--- Example 2: Slice Pool ---")
	slicePool := syncpoolallocator.NewSlicePool(1024)
	if slicePool != nil {
		slice := slicePool.Get()
		if slice != nil {
			*slice = append(*slice, []byte("Data in slice pool")...)
			fmt.Printf("Slice content: %s\n", string(*slice))
			slicePool.Put(slice)
		}
	}
	fmt.Println()

	// Example 3: Generic Pool
	fmt.Println("--- Example 3: Generic Pool ---")
	type Data struct {
		Value string
	}
	genericPool := syncpoolallocator.NewPool(
		func() *Data { return &Data{} },
		func(d *Data) { d.Value = "" },
	)
	if genericPool != nil {
		data := genericPool.Get()
		if data != nil {
			data.Value = "Generic pool data"
			fmt.Printf("Generic data: %s\n", data.Value)
			genericPool.Put(data)
		}
	}
	fmt.Println()

	// Example 4: Metrics Pool
	fmt.Println("--- Example 4: Metrics Pool ---")
	metricsPool := syncpoolallocator.NewMetricsPool(func() interface{} {
		return new(bytes.Buffer)
	})
	if metricsPool != nil {
		for i := 0; i < 10; i++ {
			obj := metricsPool.Get()
			metricsPool.Put(obj)
		}
		stats := metricsPool.Stats()
		fmt.Printf("Pool Stats: Gets=%d, Puts=%d, News=%d, HitRate=%.1f%%\n",
			stats.Gets, stats.Puts, stats.News, stats.HitRate)
	}
	fmt.Println()

	// Example 5: Size-Classed Pool
	fmt.Println("--- Example 5: Size-Classed Pool ---")
	sizePool := syncpoolallocator.NewSizeClassedPool()
	if sizePool != nil {
		sizes := []int{500, 2000, 10000, 50000}
		for _, size := range sizes {
			buf := sizePool.Get(size)
			if buf != nil {
				fmt.Printf("Requested %5d bytes → got buffer with cap=%d\n", size, cap(*buf))
				sizePool.Put(buf)
			}
		}
	}
	fmt.Println()

	fmt.Println("=== All examples completed ===")
}
