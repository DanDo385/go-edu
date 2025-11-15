//go:build solution
// +build solution

package exercise

import (
	"sync"
)

// ============================================================================
// Basic Generic Functions
// ============================================================================

// Identity returns the value unchanged.
func Identity[T any](val T) T {
	return val
}

// Contains checks if a slice contains a value.
func Contains[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Reverse returns a new slice with elements in reverse order.
func Reverse[T any](slice []T) []T {
	n := len(slice)
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = slice[n-1-i]
	}
	return result
}

// ============================================================================
// Map, Filter, Reduce, FlatMap
// ============================================================================

// Map applies a function to each element of a slice.
func Map[T, U any](data []T, fn func(T) U) []U {
	result := make([]U, len(data))
	for i, item := range data {
		result[i] = fn(item)
	}
	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](data []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(data))
	for _, item := range data {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce combines all elements into a single value.
func Reduce[T, U any](data []T, initial U, fn func(U, T) U) U {
	acc := initial
	for _, item := range data {
		acc = fn(acc, item)
	}
	return acc
}

// FlatMap applies a function that returns a slice, then flattens the results.
func FlatMap[T, U any](data []T, fn func(T) []U) []U {
	result := make([]U, 0, len(data))
	for _, item := range data {
		result = append(result, fn(item)...)
	}
	return result
}

// ============================================================================
// Parallel Map-Reduce
// ============================================================================

// ParallelMap applies a function to each element in parallel using worker pool.
func ParallelMap[T, U any](data []T, fn func(T) U, numWorkers int) []U {
	n := len(data)
	if n == 0 {
		return []U{}
	}

	// Pre-allocate result slice
	result := make([]U, n)

	// Channel for distributing work (indices)
	jobs := make(chan int, n)

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				result[i] = fn(data[i])
			}
		}()
	}

	// Send all indices to workers
	for i := 0; i < n; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	return result
}

// ParallelReduce reduces data in parallel by splitting into chunks.
// LIMITATION: This only works when the reduce function can combine values of the
// same type (T). The operation must be associative. For example:
// - Sum: works (int + int = int)
// - Product: works (int * int = int)
// - String concatenation: works
// This is a simplified implementation for educational purposes.
func ParallelReduce[T, U any](data []T, initial U, fn func(U, T) U, numWorkers int) U {
	n := len(data)
	if n == 0 {
		return initial
	}

	// For small datasets, use sequential reduce (overhead not worth it)
	if n < 100 {
		return Reduce(data, initial, fn)
	}

	// Calculate chunk size
	chunkSize := (n + numWorkers - 1) / numWorkers

	// Channel for partial results
	partials := make(chan U, numWorkers)

	// Process chunks in parallel
	var wg sync.WaitGroup
	actualWorkers := 0

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		if start >= n {
			break
		}

		actualWorkers++
		wg.Add(1)
		go func(chunk []T) {
			defer wg.Done()
			// Reduce this chunk
			acc := initial
			for _, item := range chunk {
				acc = fn(acc, item)
			}
			partials <- acc
		}(data[start:end])
	}

	// Close channel when all workers are done
	go func() {
		wg.Wait()
		close(partials)
	}()

	// Collect all partial results
	results := make([]U, 0, actualWorkers)
	for partial := range partials {
		results = append(results, partial)
	}

	// Combine partial results sequentially
	// For this to work correctly, we need the operation to be associative
	// and we're relying on the caller to use T=U for the accumulator
	acc := initial
	for _, partial := range results {
		// We have to do an unsafe conversion here because Go's type system
		// can't express the constraint that U must equal T
		// This will work for common cases like sum/product where T=U
		// Cast partial (type U) to interface{} and back to T, then apply fn
		var item T
		if v, ok := any(partial).(T); ok {
			item = v
			acc = fn(acc, item)
		}
	}

	return acc
}

// ============================================================================
// Generic Data Structures
// ============================================================================

// Optional represents a value that may or may not exist.
type Optional[T any] struct {
	value T
	valid bool
}

// Some creates an Optional with a value.
func Some[T any](val T) Optional[T] {
	return Optional[T]{value: val, valid: true}
}

// None creates an empty Optional.
func None[T any]() Optional[T] {
	return Optional[T]{valid: false}
}

// Get returns the value and whether it exists.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.valid
}

// OrElse returns the value if it exists, otherwise returns the default.
func (o Optional[T]) OrElse(defaultVal T) T {
	if o.valid {
		return o.value
	}
	return defaultVal
}

// Note: Map method with type parameters is not supported in Go.
// You cannot add new type parameters to methods.
// Use a standalone MapOptional function instead.

// Result represents a value or an error.
type Result[T, E any] struct {
	value T
	err   E
	ok    bool
}

// Ok creates a successful Result.
func Ok[T, E any](val T) Result[T, E] {
	return Result[T, E]{value: val, ok: true}
}

// Err creates a failed Result.
func Err[T, E any](err E) Result[T, E] {
	return Result[T, E]{err: err, ok: false}
}

// Unwrap returns the value, error, and success flag.
func (r Result[T, E]) Unwrap() (T, E, bool) {
	return r.value, r.err, r.ok
}

// Note: Map method with type parameters is not supported in Go.
// You cannot add new type parameters to methods.
// Use a standalone MapResult function instead.

// Pair represents a tuple of two values.
type Pair[A, B any] struct {
	First  A
	Second B
}

// MakePair creates a new Pair.
func MakePair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{First: a, Second: b}
}

// Swap returns a new Pair with First and Second swapped.
func (p Pair[A, B]) Swap() Pair[B, A] {
	return Pair[B, A]{First: p.Second, Second: p.First}
}

// Stack is a generic LIFO data structure.
type Stack[T any] struct {
	items []T
	mu    sync.Mutex
}

// NewStack creates an empty Stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Pop removes and returns the top item.
func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// Peek returns the top item without removing it.
func (s *Stack[T]) Peek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	return s.items[len(s.items)-1], true
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}
