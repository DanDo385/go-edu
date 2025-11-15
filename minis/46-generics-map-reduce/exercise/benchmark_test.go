package exercise

import (
	"runtime"
	"testing"
)

// ============================================================================
// Sequential vs Parallel Benchmarks
// ============================================================================

// Expensive computation function for realistic benchmarks
func expensiveComputation(x int) int {
	result := x
	for i := 0; i < 100; i++ {
		result = (result*result + x) % 1000000
	}
	return result
}

// BenchmarkMapSequential measures sequential map performance
func BenchmarkMapSequential(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(string(rune(size)), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Map(data, expensiveComputation)
			}
		})
	}
}

// BenchmarkMapParallel measures parallel map performance
func BenchmarkMapParallel(b *testing.B) {
	sizes := []int{100, 1000, 10000}
	workers := runtime.NumCPU()

	for _, size := range sizes {
		b.Run(string(rune(size)), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ParallelMap(data, expensiveComputation, workers)
			}
		})
	}
}

// BenchmarkReduceSequential measures sequential reduce performance
func BenchmarkReduceSequential(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(string(rune(size)), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Reduce(data, 0, func(acc, x int) int { return acc + x })
			}
		})
	}
}

// BenchmarkReduceParallel measures parallel reduce performance
func BenchmarkReduceParallel(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}
	workers := runtime.NumCPU()

	for _, size := range sizes {
		b.Run(string(rune(size)), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ParallelReduce(data, 0, func(acc, x int) int { return acc + x }, workers)
			}
		})
	}
}

// ============================================================================
// Generic vs Non-Generic Benchmarks
// ============================================================================

// Non-generic map for comparison
func mapInt(data []int, fn func(int) int) []int {
	result := make([]int, len(data))
	for i, item := range data {
		result[i] = fn(item)
	}
	return result
}

// BenchmarkGenericMap measures generic map performance
func BenchmarkGenericMap(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(data, func(x int) int { return x * 2 })
	}
}

// BenchmarkNonGenericMap measures non-generic map performance
func BenchmarkNonGenericMap(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapInt(data, func(x int) int { return x * 2 })
	}
}

// ============================================================================
// Data Structure Benchmarks
// ============================================================================

// BenchmarkStack measures stack operations
func BenchmarkStack(b *testing.B) {
	b.Run("Push", func(b *testing.B) {
		stack := NewStack[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stack.Push(i)
		}
	})

	b.Run("Pop", func(b *testing.B) {
		stack := NewStack[int]()
		for i := 0; i < b.N; i++ {
			stack.Push(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stack.Pop()
		}
	})

	b.Run("Peek", func(b *testing.B) {
		stack := NewStack[int]()
		stack.Push(42)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stack.Peek()
		}
	})
}

// BenchmarkOptional measures Optional operations
func BenchmarkOptional(b *testing.B) {
	b.Run("Some", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Some(42)
		}
	})

	b.Run("Get", func(b *testing.B) {
		opt := Some(42)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			opt.Get()
		}
	})

	b.Run("OrElse", func(b *testing.B) {
		opt := Some(42)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			opt.OrElse(0)
		}
	})
}

// ============================================================================
// Filter Benchmarks
// ============================================================================

func BenchmarkFilter(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(string(rune(size)), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Filter(data, func(x int) bool { return x%2 == 0 })
			}
		})
	}
}

// ============================================================================
// FlatMap Benchmarks
// ============================================================================

func BenchmarkFlatMap(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FlatMap(data, func(x int) []int {
			return []int{x, x * 2}
		})
	}
}

// ============================================================================
// Worker Count Comparison
// ============================================================================

func BenchmarkParallelMapWorkers(b *testing.B) {
	data := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	workerCounts := []int{1, 2, 4, 8, 16}
	for _, workers := range workerCounts {
		b.Run(string(rune(workers)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ParallelMap(data, expensiveComputation, workers)
			}
		})
	}
}

// ============================================================================
// Memory Allocation Benchmarks
// ============================================================================

func BenchmarkMapAllocations(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(data, func(x int) int { return x * 2 })
	}
}

func BenchmarkFilterAllocations(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Filter(data, func(x int) bool { return x%2 == 0 })
	}
}

func BenchmarkReduceAllocations(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Reduce(data, 0, func(acc, x int) int { return acc + x })
	}
}
