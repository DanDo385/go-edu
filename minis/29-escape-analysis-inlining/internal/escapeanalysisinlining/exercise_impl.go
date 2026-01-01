package escapeanalysisinlining

import "fmt"

// SumIntsNaive is a baseline implementation for benchmarks.
func SumIntsNaive(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}

func SumIntsOptimized(values []int) int { return SumIntsOptimizedSolution(values) }

// CalculateAreaNaive is a baseline implementation for benchmarks.
func CalculateAreaNaive(width, height float64) float64 {
	return width * height
}

func CalculateAreaOptimized(width, height float64) float64 {
	return CalculateAreaOptimizedSolution(width, height)
}

// JoinStringsNaive is a baseline implementation for benchmarks.
func JoinStringsNaive(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += separator
		result += parts[i]
	}
	return result
}

func JoinStringsOptimized(parts []string, separator string) string {
	return JoinStringsOptimizedSolution(parts, separator)
}

// AreaPointerReceiver is a baseline method for benchmarking receiver choices.
func (r *Rectangle) AreaPointerReceiver() float64 { return r.Width * r.Height }

func (r Rectangle) AreaValueReceiver() float64 { return r.AreaValueReceiverSolution() }

// ProcessItemsNaive is a baseline implementation for benchmarks.
func ProcessItemsNaive(items []string) [][]byte {
	if len(items) == 0 {
		return [][]byte{}
	}

	results := make([][]byte, 0, len(items))
	for _, item := range items {
		results = append(results, []byte("processed: "+item))
	}
	return results
}

func ProcessItemsOptimized(items []string) [][]byte { return ProcessItemsOptimizedSolution(items) }

// FormatIntNaive is a baseline implementation for benchmarks.
func FormatIntNaive(prefix string, value int) string {
	return fmt.Sprintf("%s%d", prefix, value)
}

func FormatIntOptimized(prefix string, value int) string {
	return FormatIntOptimizedSolution(prefix, value)
}

func FilterPositiveOptimized(numbers []int) []int {
	return FilterPositiveOptimizedSolution(numbers)
}

// FilterPositiveNaive is a baseline implementation (no pre-allocation).
func FilterPositiveNaive(numbers []int) []int {
	var result []int
	for _, n := range numbers {
		if n > 0 {
			result = append(result, n)
		}
	}
	return result
}

func GetConfigOptimized() Config { return GetConfigOptimizedSolution() }

// GetConfigNaive is a baseline implementation used in allocation benchmarks.
func GetConfigNaive() Config {
	cfg := &Config{
		Host: "localhost",
		Port: 8080,
	}
	return *cfg
}
