package exercise

import (
	"bytes"
	"strconv"
	"strings"
)

// ============================================================================
// SOLUTION 1: Fix Unnecessary Escapes
// ============================================================================

// SumIntsOptimizedSolution calculates the sum without escaping.
// IMPROVEMENT: No slice allocation needed - just iterate directly.
func SumIntsOptimizedSolution(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
	// Escape analysis: sum stays on stack (it's just an int)
	// No new slice allocation, so nothing escapes
}

// ============================================================================
// SOLUTION 2: Enable Inlining
// ============================================================================

// CalculateAreaOptimizedSolution is simple enough to be inlined.
// IMPROVEMENT: Removed all unnecessary conditionals and complexity.
func CalculateAreaOptimizedSolution(width, height float64) float64 {
	return width * height
	// This is simple enough that the compiler will inline it
	// Run: go build -gcflags='-m' to verify "can inline CalculateAreaOptimizedSolution"
}

// ============================================================================
// SOLUTION 3: Optimize String Building
// ============================================================================

// JoinStringsOptimizedSolution uses strings.Builder for efficiency.
// IMPROVEMENT: Single allocation for the final string, no intermediate strings.
func JoinStringsOptimizedSolution(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}

	// Calculate total size needed
	totalLen := len(parts[0])
	for i := 1; i < len(parts); i++ {
		totalLen += len(separator) + len(parts[i])
	}

	// Pre-allocate the builder's internal buffer
	var builder strings.Builder
	builder.Grow(totalLen)

	// Build the string
	builder.WriteString(parts[0])
	for i := 1; i < len(parts); i++ {
		builder.WriteString(separator)
		builder.WriteString(parts[i])
	}

	return builder.String()
	// strings.Builder is optimized to minimize allocations
	// Only 1-2 allocations total vs N allocations with concatenation
}

// ============================================================================
// SOLUTION 4: Pointer vs Value Receivers
// ============================================================================

// AreaValueReceiverSolution uses value receiver for better performance.
// IMPROVEMENT: For small structs (16 bytes), value receivers are faster.
func (r Rectangle) AreaValueReceiverSolution() float64 {
	return r.Width * r.Height
	// Value receiver: struct is copied (16 bytes)
	// This is faster than pointer dereference for small structs
	// Also allows the compiler to optimize better (can inline)
}

// ============================================================================
// SOLUTION 5: Optimize Buffer Reuse
// ============================================================================

// ProcessItemsOptimizedSolution reuses a single buffer.
// IMPROVEMENT: One buffer allocation vs N buffer allocations.
func ProcessItemsOptimizedSolution(items []string) [][]byte {
	if len(items) == 0 {
		return [][]byte{}
	}

	results := make([][]byte, 0, len(items))
	buf := &bytes.Buffer{}

	for _, item := range items {
		buf.Reset() // Clear the buffer for reuse
		buf.WriteString("processed: ")
		buf.WriteString(item)

		// Important: Copy the bytes, don't reference buf.Bytes() directly
		// buf.Bytes() returns a slice that will be invalidated on next Reset()
		data := make([]byte, buf.Len())
		copy(data, buf.Bytes())

		results = append(results, data)
	}

	return results
	// Allocation count: 1 buffer + len(items) result slices
	// vs naive: len(items) buffers + len(items) result slices
}

// ============================================================================
// SOLUTION 6: Avoid Interface{} Boxing
// ============================================================================

// FormatIntOptimizedSolution formats without interface{} escape.
// IMPROVEMENT: Direct string building, no fmt package, no escaping.
func FormatIntOptimizedSolution(prefix string, value int) string {
	var builder strings.Builder
	builder.Grow(len(prefix) + 20) // Estimate size (20 chars for int)
	builder.WriteString(prefix)
	builder.WriteString(strconv.Itoa(value))
	return builder.String()
	// No interface{} means no escape
	// strconv.Itoa is much faster than fmt.Sprintf for this use case
}

// Alternative using manual conversion (even faster for small ints):
func FormatIntOptimizedManual(prefix string, value int) string {
	// For educational purposes - shows how to avoid strconv too
	var builder strings.Builder
	builder.Grow(len(prefix) + 20)
	builder.WriteString(prefix)

	if value == 0 {
		builder.WriteRune('0')
		return builder.String()
	}

	negative := value < 0
	if negative {
		value = -value
		builder.WriteRune('-')
	}

	// Convert digits (in reverse order)
	digits := make([]byte, 0, 10)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}

	// Write digits in correct order
	for i := len(digits) - 1; i >= 0; i-- {
		builder.WriteByte(digits[i])
	}

	return builder.String()
}

// ============================================================================
// SOLUTION 7: Pre-Allocate Slices
// ============================================================================

// FilterPositiveOptimizedSolution pre-allocates to avoid reallocations.
// IMPROVEMENT: Single allocation vs multiple grow operations.
func FilterPositiveOptimizedSolution(numbers []int) []int {
	// Pre-allocate for worst case (all positive)
	result := make([]int, 0, len(numbers))

	for _, n := range numbers {
		if n > 0 {
			result = append(result, n)
		}
	}

	return result
	// Worst case: 1 allocation for capacity len(numbers)
	// vs naive: log2(N) allocations as slice grows
	// Even if we over-allocate, it's a small price for no reallocations
}

// Alternative: If you know the approximate hit rate, you can optimize further:
func FilterPositiveOptimizedEstimate(numbers []int) []int {
	// If you know ~50% are positive, allocate len/2
	result := make([]int, 0, len(numbers)/2)

	for _, n := range numbers {
		if n > 0 {
			result = append(result, n)
		}
	}

	return result
	// Trades off potential reallocation vs over-allocation
}

// ============================================================================
// SOLUTION 8: Escape Analysis Challenge
// ============================================================================

// GetConfigOptimizedSolution returns config by value.
// IMPROVEMENT: Config stays on stack, no heap allocation.
func GetConfigOptimizedSolution() Config {
	return Config{
		Host: "localhost",
		Port: 8080,
	}
	// Return by value: struct copied to caller's stack
	// No pointer, so no escape to heap
	// For small structs like Config (~24 bytes), this is very fast
}

// Note: If Config were very large (>1KB), returning a pointer might be better
// to avoid the copy cost. But for typical config structs, value is optimal.

// ============================================================================
// ADVANCED SOLUTIONS
// ============================================================================

// Example: Zero-allocation string building for known size
func BuildStringNoAlloc(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}

	// Calculate exact size
	totalLen := len(parts[0])
	for i := 1; i < len(parts); i++ {
		totalLen += len(separator) + len(parts[i])
	}

	// Pre-allocate exact buffer
	buf := make([]byte, 0, totalLen)
	buf = append(buf, parts[0]...)

	for i := 1; i < len(parts); i++ {
		buf = append(buf, separator...)
		buf = append(buf, parts[i]...)
	}

	return string(buf)
	// Exactly 1 allocation (the buffer)
	// No reallocations, no waste
}

// Example: Escape elimination through inlining
type Point struct {
	X, Y float64
}

// This can be inlined, eliminating allocation
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

// Usage: p := NewPoint(1, 2)
// After inlining: p := Point{X: 1, Y: 2}
// Compiler sees through the function call

// Example: Bounds check elimination
func SumArrayNoBoundsCheck(arr *[1000]int) int {
	sum := 0
	n := len(arr)
	for i := 0; i < n; i++ {
		sum += arr[i] // Bounds check can be eliminated
	}
	return sum
	// Compiler proves: i < n == len(arr), so arr[i] is always valid
}

// Example: Using sync.Pool for buffer reuse (advanced)
/*
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func ProcessWithPool(item string) []byte {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteString("processed: ")
	buf.WriteString(item)
	result := buf.Bytes()
	bufferPool.Put(buf) // Return to pool
	return result
}
*/

// ============================================================================
// PERFORMANCE COMPARISON HELPERS
// ============================================================================

// These functions demonstrate the performance difference

func EscapingAllocation() *int {
	x := 42
	return &x // Escapes to heap
}

func NonEscapingAllocation() int {
	x := 42
	return x // Stays on stack
}

// Benchmark these to see ~50x performance difference

func InlinableFunction(a, b int) int {
	return a + b
}

func NonInlinableFunction(a, b, c, d, e, f int) int {
	// Too complex to inline
	result := 0
	if a > b {
		result += c * d
	} else {
		result += e * f
	}
	return result
}

// ============================================================================
// ESCAPE ANALYSIS PATTERNS
// ============================================================================

// Pattern 1: Return value, not pointer (if struct is small)
type SmallStruct struct{ A, B int }

func Good() SmallStruct     { return SmallStruct{1, 2} } // Stack
func Bad() *SmallStruct     { s := SmallStruct{1, 2}; return &s } // Heap

// Pattern 2: Use local slices when size is known and small
func GoodLocal() int {
	s := [5]int{1, 2, 3, 4, 5} // Stack array
	return s[0]
}

func BadLocal() []int {
	s := []int{1, 2, 3, 4, 5} // Slice escapes if returned
	return s
}

// Pattern 3: Avoid capturing variables in closures if possible
func GoodClosure() func() int {
	x := 42
	return func() int { return x } // x must escape (captured)
}

func BetterNoClosure(x int) int {
	return x // No closure, no escape
}

// Pattern 4: Use concrete types instead of interfaces in hot paths
func GoodConcrete(x int) int     { return x * 2 }
func BadInterface(x interface{}) interface{} { return x } // x escapes
