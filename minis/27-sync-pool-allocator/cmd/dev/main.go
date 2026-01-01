package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/27-sync-pool-allocator/internal/syncpoolallocator"
)

func main() {
	fmt.Println("=== Debug Harness: sync.Pool Allocator ===")
	fmt.Println("Set breakpoints in internal/syncpoolallocator/ to debug")
	fmt.Println()

	// Test 1: Buffer Pool
	fmt.Println("Testing BufferPool...")
	bufPool := syncpoolallocator.NewBufferPool()
	if bufPool != nil {
		buf := bufPool.Get()
		if buf != nil {
			buf.WriteString("test data")
			fmt.Printf("✓ BufferPool works: %s\n", buf.String())
			bufPool.Put(buf)
		} else {
			fmt.Println("✗ BufferPool.Get() returned nil - needs implementation")
		}
	} else {
		fmt.Println("✗ NewBufferPool() returned nil - needs implementation")
	}
	fmt.Println()

	// Test 2: Slice Pool
	fmt.Println("Testing SlicePool...")
	slicePool := syncpoolallocator.NewSlicePool(1024)
	if slicePool != nil {
		slice := slicePool.Get()
		if slice != nil {
			fmt.Printf("✓ SlicePool works: capacity=%d\n", cap(*slice))
			slicePool.Put(slice)
		} else {
			fmt.Println("✗ SlicePool.Get() returned nil - needs implementation")
		}
	} else {
		fmt.Println("✗ NewSlicePool() returned nil - needs implementation")
	}
	fmt.Println()

	// Test 3: Metrics Pool
	fmt.Println("Testing MetricsPool...")
	metricsPool := syncpoolallocator.NewMetricsPool(func() interface{} {
		return new(struct{})
	})
	if metricsPool != nil {
		for i := 0; i < 5; i++ {
			obj := metricsPool.Get()
			if obj != nil {
				metricsPool.Put(obj)
			}
		}
		stats := metricsPool.Stats()
		fmt.Printf("✓ MetricsPool works: Gets=%d, HitRate=%.1f%%\n", stats.Gets, stats.HitRate)
	} else {
		fmt.Println("✗ NewMetricsPool() returned nil - needs implementation")
	}

	fmt.Println()
	fmt.Println("=== Debug session complete ===")
}
