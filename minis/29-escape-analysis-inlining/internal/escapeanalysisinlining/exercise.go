//go:build !solution && !reference


package escapeanalysisinlining

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
func SumIntsOptimized(values []int) int {
	// TODO: Implement this function without creating an unnecessary copy of the slice.
	// The `SumIntsNaive` function allocates a new slice (`result`) on the heap because the compiler can't be sure if it's needed after the function returns.
	// However, we are only calculating a sum. There is no need to copy the input slice at all.
	// By iterating directly over the `values` slice, you avoid the allocation entirely.
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
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
func CalculateAreaOptimized(width, height float64) float64 {
	// TODO: Implement a simplified version of `CalculateAreaNaive` that the compiler can inline.

	// Inlining is a compiler optimization where the body of a function is inserted directly into the caller's code, avoiding the overhead of a function call.
	// The Go compiler will not inline functions it considers "too complex" (based on the number of AST nodes).
	// The `CalculateAreaNaive` function has unnecessary `if/else` branching that makes it too complex.

	// By simplifying the logic to a single expression, you make it a candidate for inlining.
	// You can check if a function was inlined by running `go build -gcflags='-m'`.
	return width * height
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
func JoinStringsOptimized(parts []string, separator string) string {
	// TODO: Implement this function using `strings.Builder` to avoid numerous small allocations.

	// Step 1: Handle the edge case of zero parts.
	// Step 2: Create a `strings.Builder`.
	// Step 3: (Advanced Optimization) Calculate the total length of the final string and use `builder.Grow()` to pre-allocate the buffer. This will result in only one memory allocation.
	// Step 4: Write the first part to the builder.
	// Step 5: Loop through the rest of the parts, writing the separator and then the part.
	// Step 6: Return `builder.String()`.
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
func (r Rectangle) AreaValueReceiver() float64 {
	// TODO: Implement this method.

	// For small structs (generally a few machine words, e.g., under 64 bytes), passing by value can be more efficient than passing by pointer.
	// - A value receiver copies the struct onto the function's stack frame. For a small struct, this is a very fast operation.
	// - A pointer receiver passes a pointer. Accessing the struct's fields then requires a pointer dereference, which can be slightly slower and can prevent certain compiler optimizations like inlining.
	//
	// In this case, `Rectangle` is only 16 bytes (two float64s). It's a perfect candidate for a value receiver.
	return r.Width * r.Height
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
func ProcessItemsOptimized(items []string) [][]byte {
	// TODO: Implement this function by creating a single buffer and reusing it for each item.

	// The naive version allocates a new `bytes.Buffer` for every single item, leading to a large number of allocations.

	// Step 1: Create one `bytes.Buffer` *before* the loop.
	// Step 2: Pre-allocate the `results` slice.
	// Step 3: Inside the loop:
	//   - Call `buf.Reset()` to clear the buffer from the previous iteration.
	//   - Write the new data to the buffer.
	//   - **CRITICAL**: `buf.Bytes()` returns a slice that points to the buffer's internal memory. This memory will be overwritten in the next iteration. You MUST make a copy of the bytes before appending them to your results.
	//     - `data := make([]byte, buf.Len())`
	//     - `copy(data, buf.Bytes())`
	//     - `results = append(results, data)`
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
func FormatIntOptimized(prefix string, value int) string {
	// TODO: Implement this function without using `fmt.Sprintf` to avoid heap allocation.

	// `fmt.Sprintf("%s%d", ...)` is convenient, but it's designed to work with any type. To do this, it packs its arguments into an `[]interface{}`, which is a slice of empty interfaces.
	// When a value type like `int` is put into an interface, it gets "boxed", which means it's allocated on the heap. This is what causes the escape.

	// Step 1: Use a `strings.Builder` for efficient string construction.
	// Step 2: Use `strconv.Itoa(value)` to convert the integer to a string. `Itoa` is highly optimized for this specific task and does not cause the value to escape to the heap.
	// Step 3: Write the prefix and the converted integer to the builder.
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
func FilterPositiveOptimized(numbers []int) []int {
	// TODO: Implement this function by pre-allocating the result slice.

	// The naive version will re-allocate its `result` slice multiple times as it grows.
	// We can avoid all re-allocations by allocating for the worst-case scenario.

	// Step 1: Create the result slice with a capacity equal to the length of the input slice.
	// - `result := make([]int, 0, len(numbers))`
	// - In the worst case, all numbers are positive, and the final slice will have `len(numbers)` elements. This guarantees that we will never need to re-allocate the underlying array.

	// Step 2: Loop through the numbers and append the positive ones to the `result` slice.
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
func GetConfigOptimized() Config {
	// TODO: Implement this function to avoid having the `Config` struct escape to the heap.

	// In `GetConfigNaive`, a pointer to the `Config` struct is returned (`*Config`).
	// The Go compiler sees that a pointer to a locally defined variable is being returned, so it cannot keep that variable on the function's stack (which is destroyed when the function returns).
	// It "escapes" the variable to the heap, which involves a memory allocation.

	// By changing the return type to `Config` (a value, not a pointer), you are telling the compiler to return a *copy* of the struct.
	// For a small struct like this, making a copy is very cheap and can be done entirely on the stack, avoiding a heap allocation and subsequent garbage collection.
	return Config{
		Host: "localhost",
		Port: 8080,
	}
}

// ============================================================================
// HELPER FUNCTIONS (for testing)
// ============================================================================

// You can add helper functions here if needed for your implementations.
