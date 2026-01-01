//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package genericsmapreduce

import "sync"
// TODO: implement Identity.
func Identity[T any](val T) T { panic("TODO: implement") }
// TODO: implement Contains.
func Contains[T comparable](slice []T, val T) bool { panic("TODO: implement") }
// TODO: implement Reverse.
func Reverse[T any](slice []T) []T { panic("TODO: implement") }
// TODO: implement Map.
func Map[T, U any](data []T, fn func(T) U) []U { panic("TODO: implement") }
// TODO: implement Filter.
func Filter[T any](data []T, predicate func(T) bool) []T { panic("TODO: implement") }
// TODO: implement Reduce.
func Reduce[T, U any](data []T, initial U, fn func(U, T) U) U { panic("TODO: implement") }
// TODO: implement FlatMap.
func FlatMap[T, U any](data []T, fn func(T) []U) []U { panic("TODO: implement") }
// TODO: implement ParallelMap.
func ParallelMap[T, U any](data []T, fn func(T) U, numWorkers int) []U { panic("TODO: implement") }
// TODO: implement ParallelReduce.
func ParallelReduce[T, U any](data []T, initial U, fn func(U, T) U, numWorkers int) U {
	panic("TODO: implement")
}

type Optional[T any] struct {
	value T
	valid bool
}
// TODO: implement Some.
func Some[T any](val T) Optional[T] { panic("TODO: implement") }
// TODO: implement None.
func None[T any]() Optional[T] { panic("TODO: implement") }
// TODO: implement Get.
func (o Optional[T]) Get() (T, bool) { panic("TODO: implement") }
// TODO: implement OrElse.
func (o Optional[T]) OrElse(defaultVal T) T { panic("TODO: implement") }

type Result[T, E any] struct {
	value T
	err   E
	ok    bool
}
// TODO: implement Ok.
func Ok[T, E any](val T) Result[T, E] { panic("TODO: implement") }
// TODO: implement Err.
func Err[T, E any](err E) Result[T, E] { panic("TODO: implement") }
// TODO: implement Unwrap.
func (r Result[T, E]) Unwrap() (T, E, bool) { panic("TODO: implement") }

type Pair[A, B any] struct {
	First  A
	Second B
}
// TODO: implement MakePair.
func MakePair[A, B any](a A, b B) Pair[A, B] { panic("TODO: implement") }
// TODO: implement Swap.
func (p Pair[A, B]) Swap() Pair[B, A] { panic("TODO: implement") }

type Stack[T any] struct {
	items []T
	mu    sync.Mutex
}
// TODO: implement NewStack.
func NewStack[T any]() *Stack[T] { panic("TODO: implement") }
// TODO: implement Push.
func (s *Stack[T]) Push(item T) { panic("TODO: implement") }
// TODO: implement Pop.
func (s *Stack[T]) Pop() (T, bool) { panic("TODO: implement") }
// TODO: implement Peek.
func (s *Stack[T]) Peek() (T, bool) { panic("TODO: implement") }
// TODO: implement Len.
func (s *Stack[T]) Len() int { panic("TODO: implement") }
