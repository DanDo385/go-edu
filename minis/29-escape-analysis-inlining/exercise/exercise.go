//go:build !solution
// +build !solution

package exercise

import (
	"bytes"
)

// ============================================================================
// EXERCISE 1: Fix Unnecessary Escapes
// ============================================================================

// SumIntsNaive calculates the sum of integers.
// PROBLEM: The slice escapes to heap unnecessarily.
// TODO: Modify to keep the slice on the stack.
// HINT: The slice doesn't need to be returned, only the sum.
func SumIntsNaive(values []int) int {
	// Current implementation (causes escape)
	result := make([]int, len(values))
	copy(result, values)

	sum := 0
	for _, v := range result {
		sum += v
	}
	return sum
}

// SumIntsOptimized should do the same thing but avoid the escape.
// TODO: Implement this function without creating a new slice.
func SumIntsOptimized(values []int) int {
	// TODO: Implement me!
	// Hint: Do you really need to copy the slice?
	return 0
}

// ============================================================================
// EXERCISE 2: Enable Inlining
// ============================================================================

// CalculateAreaNaive calculates the area of a rectangle.
// PROBLEM: Too complex to be inlined.
// TODO: Simplify this function so the compiler can inline it.
// HINT: Remove unnecessary complexity.
func CalculateAreaNaive(width, height float64) float64 {
	// Unnecessary complexity prevents inlining
	var result float64
	if width > 0 && height > 0 {
		temp := width * height
		if temp > 0 {
			result = temp
		} else {
			result = 0
		}
	} else {
		result = 0
	}
	return result
}

// CalculateAreaOptimized should be simple enough to inline.
// TODO: Implement a simple version that can be inlined.
// HINT: Just multiply width * height directly.
func CalculateAreaOptimized(width, height float64) float64 {
	// TODO: Implement me!
	return 0
}

// ============================================================================
// EXERCISE 3: Optimize String Building
// ============================================================================

// JoinStringsNaive concatenates strings with a separator.
// PROBLEM: Uses string concatenation in a loop (many allocations).
// TODO: This works but creates many intermediate strings.
func JoinStringsNaive(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result = result + separator + parts[i] // Each + creates a new string
	}
	return result
}

// JoinStringsOptimized should use a more efficient approach.
// TODO: Implement using strings.Builder or bytes.Buffer.
// HINT: strings.Builder is designed for this use case.
func JoinStringsOptimized(parts []string, separator string) string {
	// TODO: Implement me!
	// Hint: Use strings.Builder
	return ""
}

// ============================================================================
// EXERCISE 4: Pointer vs Value Receivers
// ============================================================================

// Rectangle represents a rectangle with width and height.
type Rectangle struct {
	Width  float64
	Height float64
}

// AreaPointerReceiver calculates area using a pointer receiver.
// PROBLEM: Pointer receiver for a small struct (16 bytes).
// This is less efficient than necessary.
func (r *Rectangle) AreaPointerReceiver() float64 {
	return r.Width * r.Height
}

// AreaValueReceiver should use a value receiver for better performance.
// TODO: Implement this method with a value receiver.
// HINT: For small structs (<=32 bytes), value receivers are often faster.
func (r Rectangle) AreaValueReceiver() float64 {
	// TODO: Implement me!
	return 0
}

// ============================================================================
// EXERCISE 5: Optimize Buffer Reuse
// ============================================================================

// ProcessItemsNaive processes items, creating a new buffer each time.
// PROBLEM: Allocates a new buffer for every item.
func ProcessItemsNaive(items []string) [][]byte {
	results := make([][]byte, 0, len(items))

	for _, item := range items {
		buf := &bytes.Buffer{} // New allocation each iteration
		buf.WriteString("processed: ")
		buf.WriteString(item)
		results = append(results, buf.Bytes())
	}

	return results
}

// ProcessItemsOptimized should reuse a single buffer.
// TODO: Implement this function reusing a single buffer.
// HINT: Reset the buffer between iterations.
func ProcessItemsOptimized(items []string) [][]byte {
	// TODO: Implement me!
	// Hint: Create one buffer outside the loop and reset it each iteration
	return nil
}

// ============================================================================
// EXERCISE 6: Avoid Interface{} Boxing
// ============================================================================

// FormatIntNaive formats an integer with a prefix.
// PROBLEM: Uses fmt.Sprintf which causes the int to escape via interface{}.
// This is fine for occasional use, but terrible in hot loops.
func FormatIntNaive(prefix string, value int) string {
	// Using fmt.Sprintf causes value to escape
	_ = prefix // Avoid unused warning for exercise
	_ = value
	// In real implementation: return fmt.Sprintf("%s%d", prefix, value)
	return "" // Placeholder to avoid import issues
}

// FormatIntOptimized should format without using interface{}.
// TODO: Implement without fmt.Sprintf.
// HINT: Use strings.Builder and manual int-to-string conversion.
func FormatIntOptimized(prefix string, value int) string {
	// TODO: Implement me!
	// Hint: Use strings.Builder and strconv.Itoa or manual conversion
	return ""
}

// ============================================================================
// EXERCISE 7: Pre-Allocate Slices
// ============================================================================

// FilterPositiveNaive filters positive numbers from a slice.
// PROBLEM: Doesn't pre-allocate, causing multiple reallocations.
func FilterPositiveNaive(numbers []int) []int {
	var result []int // No pre-allocation
	for _, n := range numbers {
		if n > 0 {
			result = append(result, n)
		}
	}
	return result
}

// FilterPositiveOptimized should pre-allocate to avoid reallocations.
// TODO: Implement with pre-allocation.
// HINT: In worst case, all numbers are positive, so pre-allocate len(numbers).
func FilterPositiveOptimized(numbers []int) []int {
	// TODO: Implement me!
	// Hint: make([]int, 0, len(numbers))
	return nil
}

// ============================================================================
// EXERCISE 8: Escape Analysis Challenge
// ============================================================================

// Config represents configuration data.
type Config struct {
	Host string
	Port int
}

// GetConfigNaive returns a pointer to a config.
// PROBLEM: Config escapes to heap.
func GetConfigNaive() *Config {
	cfg := Config{
		Host: "localhost",
		Port: 8080,
	}
	return &cfg // cfg escapes
}

// GetConfigOptimized should return config without escaping.
// TODO: Implement this to avoid heap allocation.
// HINT: Return by value, not by pointer.
func GetConfigOptimized() Config {
	// TODO: Implement me!
	return Config{}
}

// ============================================================================
// HELPER FUNCTIONS (for testing)
// ============================================================================

// You can add helper functions here if needed for your implementations.
