//go:build !solution
// +build !solution

package exercise

import (
	"sync"
)

// ============================================================================
// Basic Generic Functions
// ============================================================================

// Identity returns the value unchanged.
// This demonstrates the simplest generic function.
func Identity[T any](val T) T {
	// TODO: implement
	return val
}

// Contains checks if a slice contains a value.
// T must be comparable to use == operator.
func Contains[T comparable](slice []T, val T) bool {
	// TODO: implement
	return false
}

// Reverse returns a new slice with elements in reverse order.
func Reverse[T any](slice []T) []T {
	// TODO: implement
	return nil
}

// ============================================================================
// Map, Filter, Reduce, FlatMap
// ============================================================================

// Map applies a function to each element of a slice.
// Can transform from type T to type U.
func Map[T, U any](data []T, fn func(T) U) []U {
	// TODO: implement
	return nil
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](data []T, predicate func(T) bool) []T {
	// TODO: implement
	return nil
}

// Reduce combines all elements into a single value.
// fn takes (accumulator, current element) and returns new accumulator.
func Reduce[T, U any](data []T, initial U, fn func(U, T) U) U {
	// TODO: implement
	return initial
}

// FlatMap applies a function that returns a slice, then flattens the results.
func FlatMap[T, U any](data []T, fn func(T) []U) []U {
	// TODO: implement
	return nil
}

// ============================================================================
// Parallel Map-Reduce
// ============================================================================

// ParallelMap applies a function to each element in parallel using worker pool.
// numWorkers specifies how many goroutines to use.
func ParallelMap[T, U any](data []T, fn func(T) U, numWorkers int) []U {
	// TODO: implement
	// Hints:
	// 1. Create result slice with same length as input
	// 2. Create a channel for distributing work (indices)
	// 3. Start numWorkers goroutines
	// 4. Each worker reads indices from channel and processes data[i]
	// 5. Use sync.WaitGroup to wait for all workers
	return nil
}

// ParallelReduce reduces data in parallel by splitting into chunks.
// Only works correctly for associative operations (e.g., +, *, min, max).
// Non-associative operations (e.g., -) will give incorrect results.
func ParallelReduce[T, U any](data []T, initial U, fn func(U, T) U, numWorkers int) U {
	// TODO: implement
	// Hints:
	// 1. Calculate chunk size based on data length and numWorkers
	// 2. Process each chunk in a goroutine
	// 3. Each goroutine reduces its chunk and sends result to channel
	// 4. Combine partial results sequentially
	return initial
}

// ============================================================================
// Generic Data Structures
// ============================================================================

// Optional represents a value that may or may not exist.
// This is a type-safe alternative to using pointers or zero values.
type Optional[T any] struct {
	// TODO: add fields
}

// Some creates an Optional with a value.
func Some[T any](val T) Optional[T] {
	// TODO: implement
	return Optional[T]{}
}

// None creates an empty Optional.
func None[T any]() Optional[T] {
	// TODO: implement
	return Optional[T]{}
}

// Get returns the value and whether it exists.
func (o Optional[T]) Get() (T, bool) {
	// TODO: implement
	var zero T
	return zero, false
}

// OrElse returns the value if it exists, otherwise returns the default.
func (o Optional[T]) OrElse(defaultVal T) T {
	// TODO: implement
	var zero T
	return zero
}

// Note: Map method with type parameters is not supported in Go.
// You cannot add new type parameters to methods.
// Use a standalone MapOptional function instead.

// Result represents a value or an error.
// Similar to Rust's Result<T, E> type.
type Result[T, E any] struct {
	// TODO: add fields
}

// Ok creates a successful Result.
func Ok[T, E any](val T) Result[T, E] {
	// TODO: implement
	return Result[T, E]{}
}

// Err creates a failed Result.
func Err[T, E any](err E) Result[T, E] {
	// TODO: implement
	return Result[T, E]{}
}

// Unwrap returns the value, error, and success flag.
func (r Result[T, E]) Unwrap() (T, E, bool) {
	// TODO: implement
	var zeroT T
	var zeroE E
	return zeroT, zeroE, false
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
	// TODO: implement
	return Pair[A, B]{}
}

// Swap returns a new Pair with First and Second swapped.
func (p Pair[A, B]) Swap() Pair[B, A] {
	// TODO: implement
	return Pair[B, A]{}
}

// Stack is a generic LIFO data structure.
type Stack[T any] struct {
	// TODO: add fields
}

// NewStack creates an empty Stack.
func NewStack[T any]() *Stack[T] {
	// TODO: implement
	return nil
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	// TODO: implement
}

// Pop removes and returns the top item.
// Returns (item, true) if successful, (zero, false) if stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	// TODO: implement
	var zero T
	return zero, false
}

// Peek returns the top item without removing it.
// Returns (item, true) if successful, (zero, false) if stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	// TODO: implement
	var zero T
	return zero, false
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	// TODO: implement
	return 0
}

// ============================================================================
// Helper Types (used by ParallelMap/Reduce)
// ============================================================================

var _ sync.Mutex // Prevent "imported and not used" error
