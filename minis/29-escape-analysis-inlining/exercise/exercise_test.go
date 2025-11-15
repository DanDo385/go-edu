package exercise

import (
	"bytes"
	"strings"
	"testing"
)

// ============================================================================
// TESTS FOR EXERCISE 1: Fix Unnecessary Escapes
// ============================================================================

func TestSumIntsOptimized(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"empty slice", []int{}, 0},
		{"single element", []int{5}, 5},
		{"positive numbers", []int{1, 2, 3, 4, 5}, 15},
		{"mixed numbers", []int{-5, 10, -3, 8, 0}, 10},
		{"all negative", []int{-1, -2, -3}, -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumIntsOptimized(tt.input)
			if result != tt.expected {
				t.Errorf("SumIntsOptimized(%v) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TESTS FOR EXERCISE 2: Enable Inlining
// ============================================================================

func TestCalculateAreaOptimized(t *testing.T) {
	tests := []struct {
		name     string
		width    float64
		height   float64
		expected float64
	}{
		{"positive dimensions", 10.0, 5.0, 50.0},
		{"zero width", 0.0, 5.0, 0.0},
		{"zero height", 10.0, 0.0, 0.0},
		{"decimal dimensions", 3.5, 2.5, 8.75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAreaOptimized(tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("CalculateAreaOptimized(%f, %f) = %f; want %f",
					tt.width, tt.height, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TESTS FOR EXERCISE 3: Optimize String Building
// ============================================================================

func TestJoinStringsOptimized(t *testing.T) {
	tests := []struct {
		name      string
		parts     []string
		separator string
		expected  string
	}{
		{"empty slice", []string{}, ",", ""},
		{"single element", []string{"hello"}, ",", "hello"},
		{"two elements", []string{"hello", "world"}, " ", "hello world"},
		{"comma separated", []string{"a", "b", "c"}, ",", "a,b,c"},
		{"pipe separated", []string{"one", "two", "three"}, "|", "one|two|three"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinStringsOptimized(tt.parts, tt.separator)
			if result != tt.expected {
				t.Errorf("JoinStringsOptimized(%v, %q) = %q; want %q",
					tt.parts, tt.separator, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TESTS FOR EXERCISE 4: Pointer vs Value Receivers
// ============================================================================

func TestAreaValueReceiver(t *testing.T) {
	tests := []struct {
		name     string
		rect     Rectangle
		expected float64
	}{
		{"square", Rectangle{5, 5}, 25.0},
		{"rectangle", Rectangle{10, 3}, 30.0},
		{"zero width", Rectangle{0, 5}, 0.0},
		{"decimals", Rectangle{3.5, 2.5}, 8.75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect.AreaValueReceiver()
			if result != tt.expected {
				t.Errorf("Rectangle{%f, %f}.AreaValueReceiver() = %f; want %f",
					tt.rect.Width, tt.rect.Height, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TESTS FOR EXERCISE 5: Optimize Buffer Reuse
// ============================================================================

func TestProcessItemsOptimized(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		expected []string
	}{
		{"empty", []string{}, []string{}},
		{"single item", []string{"test"}, []string{"processed: test"}},
		{"multiple items", []string{"a", "b", "c"},
			[]string{"processed: a", "processed: b", "processed: c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessItemsOptimized(tt.items)
			if len(result) != len(tt.expected) {
				t.Errorf("ProcessItemsOptimized(%v) returned %d items; want %d",
					tt.items, len(result), len(tt.expected))
				return
			}
			for i, got := range result {
				if string(got) != tt.expected[i] {
					t.Errorf("ProcessItemsOptimized(%v)[%d] = %q; want %q",
						tt.items, i, string(got), tt.expected[i])
				}
			}
		})
	}
}

// ============================================================================
// TESTS FOR EXERCISE 6: Avoid Interface{} Boxing
// ============================================================================

func TestFormatIntOptimized(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		value    int
		expected string
	}{
		{"positive", "value: ", 42, "value: 42"},
		{"negative", "count: ", -10, "count: -10"},
		{"zero", "result: ", 0, "result: 0"},
		{"empty prefix", "", 123, "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatIntOptimized(tt.prefix, tt.value)
			if result != tt.expected {
				t.Errorf("FormatIntOptimized(%q, %d) = %q; want %q",
					tt.prefix, tt.value, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TESTS FOR EXERCISE 7: Pre-Allocate Slices
// ============================================================================

func TestFilterPositiveOptimized(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"empty", []int{}, []int{}},
		{"all positive", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all negative", []int{-1, -2, -3}, []int{}},
		{"mixed", []int{-5, 3, -2, 7, 0, 1}, []int{3, 7, 1}},
		{"with zero", []int{0, 1, 2}, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterPositiveOptimized(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("FilterPositiveOptimized(%v) length = %d; want %d",
					tt.input, len(result), len(tt.expected))
				return
			}
			for i, got := range result {
				if got != tt.expected[i] {
					t.Errorf("FilterPositiveOptimized(%v)[%d] = %d; want %d",
						tt.input, i, got, tt.expected[i])
				}
			}
		})
	}
}

// ============================================================================
// TESTS FOR EXERCISE 8: Escape Analysis Challenge
// ============================================================================

func TestGetConfigOptimized(t *testing.T) {
	cfg := GetConfigOptimized()
	if cfg.Host != "localhost" {
		t.Errorf("GetConfigOptimized().Host = %q; want %q", cfg.Host, "localhost")
	}
	if cfg.Port != 8080 {
		t.Errorf("GetConfigOptimized().Port = %d; want %d", cfg.Port, 8080)
	}
}

// ============================================================================
// BENCHMARKS
// ============================================================================

// Benchmark: Sum operations
func BenchmarkSumIntsNaive(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SumIntsNaive(data)
	}
}

func BenchmarkSumIntsOptimized(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SumIntsOptimized(data)
	}
}

// Benchmark: Area calculations
func BenchmarkCalculateAreaNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CalculateAreaNaive(10.5, 20.3)
	}
}

func BenchmarkCalculateAreaOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CalculateAreaOptimized(10.5, 20.3)
	}
}

// Benchmark: String joining
func BenchmarkJoinStringsNaive(b *testing.B) {
	parts := []string{"one", "two", "three", "four", "five"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JoinStringsNaive(parts, ",")
	}
}

func BenchmarkJoinStringsOptimized(b *testing.B) {
	parts := []string{"one", "two", "three", "four", "five"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JoinStringsOptimized(parts, ",")
	}
}

// Benchmark: Receivers
func BenchmarkAreaPointerReceiver(b *testing.B) {
	rect := &Rectangle{Width: 10.5, Height: 20.3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rect.AreaPointerReceiver()
	}
}

func BenchmarkAreaValueReceiver(b *testing.B) {
	rect := Rectangle{Width: 10.5, Height: 20.3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rect.AreaValueReceiver()
	}
}

// Benchmark: Buffer reuse
func BenchmarkProcessItemsNaive(b *testing.B) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProcessItemsNaive(items)
	}
}

func BenchmarkProcessItemsOptimized(b *testing.B) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProcessItemsOptimized(items)
	}
}

// Benchmark: Integer formatting
func BenchmarkFormatIntNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatIntNaive("value: ", 12345)
	}
}

func BenchmarkFormatIntOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatIntOptimized("value: ", 12345)
	}
}

// Benchmark: Filtering with pre-allocation
func BenchmarkFilterPositiveNaive(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i - 500 // Mix of positive and negative
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterPositiveNaive(numbers)
	}
}

func BenchmarkFilterPositiveOptimized(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i - 500 // Mix of positive and negative
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterPositiveOptimized(numbers)
	}
}

// Benchmark: Escape analysis
func BenchmarkGetConfigNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetConfigNaive()
	}
}

func BenchmarkGetConfigOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetConfigOptimized()
	}
}

// ============================================================================
// ALLOCATION BENCHMARKS
// These specifically measure allocation counts
// Run with: go test -bench=Alloc -benchmem
// ============================================================================

func BenchmarkAllocations(b *testing.B) {
	b.Run("SumNaive", func(b *testing.B) {
		data := make([]int, 100)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = SumIntsNaive(data)
		}
	})

	b.Run("SumOptimized", func(b *testing.B) {
		data := make([]int, 100)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = SumIntsOptimized(data)
		}
	})
}

// ============================================================================
// EXAMPLE BENCHMARKS (showing expected improvements)
// ============================================================================

// Example showing stack vs heap allocation
func BenchmarkStackAllocation(b *testing.B) {
	b.Run("StackValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x := 42
			_ = x
		}
	})

	b.Run("HeapPointer", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x := new(int)
			*x = 42
			_ = x
		}
	})
}

// Example showing string concatenation vs Builder
func BenchmarkStringBuilding(b *testing.B) {
	parts := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	b.Run("Concatenation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := ""
			for _, part := range parts {
				result += part
			}
			_ = result
		}
	})

	b.Run("Builder", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for _, part := range parts {
				builder.WriteString(part)
			}
			_ = builder.String()
		}
	})

	b.Run("Buffer", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			for _, part := range parts {
				buf.WriteString(part)
			}
			_ = buf.String()
		}
	})
}
