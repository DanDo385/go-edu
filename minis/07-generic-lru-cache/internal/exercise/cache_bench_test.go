package exercise

import (
	"fmt"
	"testing"
)

// BenchmarkCache_Set measures Set performance
func BenchmarkCache_Set(b *testing.B) {
	cache := New[int, int](1000, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i%1000, i)
	}
}

// BenchmarkCache_Get measures Get performance
func BenchmarkCache_Get(b *testing.B) {
	cache := New[int, int](1000, 0)
	// Pre-fill
	for i := 0; i < 1000; i++ {
		cache.Set(i, i*10)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % 1000)
	}
}

// BenchmarkCache_Mixed measures realistic workload (70% Get, 30% Set)
func BenchmarkCache_Mixed(b *testing.B) {
	cache := New[int, int](1000, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%10 < 7 {
			cache.Get(i % 1000)
		} else {
			cache.Set(i%1000, i)
		}
	}
}

// BenchmarkCache_SetWithEviction measures performance with constant eviction
func BenchmarkCache_SetWithEviction(b *testing.B) {
	cache := New[int, int](100, 0) // Small capacity
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i, i) // Every insert after 100 causes eviction
	}
}

// BenchmarkMap_Baseline compares with plain map (no LRU, no locking)
func BenchmarkMap_Baseline(b *testing.B) {
	m := make(map[int]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%10 < 7 {
			_ = m[i%1000]
		} else {
			m[i%1000] = i
		}
	}
}

// Benchmark different cache sizes
func BenchmarkCache_Sizes(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			cache := New[int, int](size, 0)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if i%2 == 0 {
					cache.Set(i%size, i)
				} else {
					cache.Get(i % size)
				}
			}
		})
	}
}
