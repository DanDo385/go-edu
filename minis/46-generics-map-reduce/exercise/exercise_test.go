package exercise

import (
	"reflect"
	"strings"
	"testing"
)

// ============================================================================
// Basic Generic Functions Tests
// ============================================================================

func TestIdentity(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 42},
		{"string", "hello"},
		{"float", 3.14},
		{"bool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case int:
				if got := Identity(v); got != v {
					t.Errorf("Identity(%v) = %v, want %v", v, got, v)
				}
			case string:
				if got := Identity(v); got != v {
					t.Errorf("Identity(%v) = %v, want %v", v, got, v)
				}
			case float64:
				if got := Identity(v); got != v {
					t.Errorf("Identity(%v) = %v, want %v", v, got, v)
				}
			case bool:
				if got := Identity(v); got != v {
					t.Errorf("Identity(%v) = %v, want %v", v, got, v)
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		val   int
		want  bool
	}{
		{"found", []int{1, 2, 3, 4, 5}, 3, true},
		{"not found", []int{1, 2, 3, 4, 5}, 10, false},
		{"empty slice", []int{}, 1, false},
		{"first element", []int{1, 2, 3}, 1, true},
		{"last element", []int{1, 2, 3}, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.slice, tt.val); got != tt.want {
				t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.val, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		want := []int{5, 4, 3, 2, 1}
		got := Reverse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Reverse(%v) = %v, want %v", input, got, want)
		}
	})

	t.Run("string slice", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		want := []string{"c", "b", "a"}
		got := Reverse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Reverse(%v) = %v, want %v", input, got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		want := []int{}
		got := Reverse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Reverse(%v) = %v, want %v", input, got, want)
		}
	})
}

// ============================================================================
// Map, Filter, Reduce, FlatMap Tests
// ============================================================================

func TestMap(t *testing.T) {
	t.Run("int to int", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		got := Map(input, func(x int) int { return x * 2 })
		want := []int{2, 4, 6, 8, 10}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map(double) = %v, want %v", got, want)
		}
	})

	t.Run("int to string", func(t *testing.T) {
		input := []int{1, 2, 3}
		got := Map(input, func(x int) string { return strings.Repeat("*", x) })
		want := []string{"*", "**", "***"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map(to stars) = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		got := Map(input, func(x int) int { return x * 2 })
		want := []int{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map(empty) = %v, want %v", got, want)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("filter evens", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		got := Filter(input, func(x int) bool { return x%2 == 0 })
		want := []int{2, 4, 6}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Filter(even) = %v, want %v", got, want)
		}
	})

	t.Run("filter none", func(t *testing.T) {
		input := []int{1, 3, 5}
		got := Filter(input, func(x int) bool { return x%2 == 0 })
		want := []int{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Filter(even) = %v, want %v", got, want)
		}
	})

	t.Run("filter all", func(t *testing.T) {
		input := []int{2, 4, 6}
		got := Filter(input, func(x int) bool { return x%2 == 0 })
		want := []int{2, 4, 6}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Filter(even) = %v, want %v", got, want)
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		got := Reduce(input, 0, func(acc, x int) int { return acc + x })
		want := 15
		if got != want {
			t.Errorf("Reduce(sum) = %v, want %v", got, want)
		}
	})

	t.Run("product", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		got := Reduce(input, 1, func(acc, x int) int { return acc * x })
		want := 120
		if got != want {
			t.Errorf("Reduce(product) = %v, want %v", got, want)
		}
	})

	t.Run("concatenate strings", func(t *testing.T) {
		input := []string{"Hello", "World", "Go"}
		got := Reduce(input, "", func(acc, s string) string {
			if acc == "" {
				return s
			}
			return acc + " " + s
		})
		want := "Hello World Go"
		if got != want {
			t.Errorf("Reduce(concat) = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		got := Reduce(input, 100, func(acc, x int) int { return acc + x })
		want := 100
		if got != want {
			t.Errorf("Reduce(empty) = %v, want %v", got, want)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("split strings", func(t *testing.T) {
		input := []string{"hello world", "go lang"}
		got := FlatMap(input, func(s string) []string { return strings.Fields(s) })
		want := []string{"hello", "world", "go", "lang"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("FlatMap(split) = %v, want %v", got, want)
		}
	})

	t.Run("duplicate elements", func(t *testing.T) {
		input := []int{1, 2, 3}
		got := FlatMap(input, func(x int) []int { return []int{x, x} })
		want := []int{1, 1, 2, 2, 3, 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("FlatMap(duplicate) = %v, want %v", got, want)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		input := []int{1, 2, 3}
		got := FlatMap(input, func(x int) []int { return []int{} })
		want := []int{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("FlatMap(empty) = %v, want %v", got, want)
		}
	})
}

// ============================================================================
// Parallel Map-Reduce Tests
// ============================================================================

func TestParallelMap(t *testing.T) {
	t.Run("basic functionality", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		got := ParallelMap(input, func(x int) int { return x * 2 }, 2)
		want := []int{2, 4, 6, 8, 10}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ParallelMap(double) = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		got := ParallelMap(input, func(x int) int { return x * 2 }, 4)
		want := []int{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ParallelMap(empty) = %v, want %v", got, want)
		}
	})

	t.Run("large dataset", func(t *testing.T) {
		n := 10000
		input := make([]int, n)
		for i := 0; i < n; i++ {
			input[i] = i
		}

		got := ParallelMap(input, func(x int) int { return x * 2 }, 4)

		if len(got) != n {
			t.Errorf("ParallelMap length = %v, want %v", len(got), n)
		}

		for i := 0; i < n; i++ {
			if got[i] != i*2 {
				t.Errorf("ParallelMap[%d] = %v, want %v", i, got[i], i*2)
			}
		}
	})
}

func TestParallelReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		got := ParallelReduce(input, 0, func(acc, x int) int { return acc + x }, 4)
		want := 55
		if got != want {
			t.Errorf("ParallelReduce(sum) = %v, want %v", got, want)
		}
	})

	t.Run("product", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		got := ParallelReduce(input, 1, func(acc, x int) int { return acc * x }, 2)
		want := 120
		if got != want {
			t.Errorf("ParallelReduce(product) = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		got := ParallelReduce(input, 100, func(acc, x int) int { return acc + x }, 4)
		want := 100
		if got != want {
			t.Errorf("ParallelReduce(empty) = %v, want %v", got, want)
		}
	})
}

// ============================================================================
// Generic Data Structures Tests
// ============================================================================

func TestOptional(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		opt := Some(42)
		val, ok := opt.Get()
		if !ok || val != 42 {
			t.Errorf("Some(42).Get() = (%v, %v), want (42, true)", val, ok)
		}
	})

	t.Run("None", func(t *testing.T) {
		opt := None[int]()
		_, ok := opt.Get()
		if ok {
			t.Errorf("None().Get() should return false")
		}
	})

	t.Run("OrElse with value", func(t *testing.T) {
		opt := Some(42)
		got := opt.OrElse(99)
		if got != 42 {
			t.Errorf("Some(42).OrElse(99) = %v, want 42", got)
		}
	})

	t.Run("OrElse without value", func(t *testing.T) {
		opt := None[int]()
		got := opt.OrElse(99)
		if got != 99 {
			t.Errorf("None().OrElse(99) = %v, want 99", got)
		}
	})

}

func TestResult(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		result := Ok[int, string](42)
		val, err, ok := result.Unwrap()
		if !ok || val != 42 || err != "" {
			t.Errorf("Ok(42).Unwrap() = (%v, %v, %v), want (42, \"\", true)", val, err, ok)
		}
	})

	t.Run("Err", func(t *testing.T) {
		result := Err[int, string]("error occurred")
		val, err, ok := result.Unwrap()
		if ok || val != 0 || err != "error occurred" {
			t.Errorf("Err(\"error occurred\").Unwrap() = (%v, %v, %v), want (0, \"error occurred\", false)", val, err, ok)
		}
	})

}

func TestPair(t *testing.T) {
	t.Run("MakePair", func(t *testing.T) {
		pair := MakePair("hello", 42)
		if pair.First != "hello" || pair.Second != 42 {
			t.Errorf("MakePair(\"hello\", 42) = {%v, %v}, want {\"hello\", 42}", pair.First, pair.Second)
		}
	})

	t.Run("Swap", func(t *testing.T) {
		pair := MakePair("hello", 42)
		swapped := pair.Swap()
		if swapped.First != 42 || swapped.Second != "hello" {
			t.Errorf("Swap() = {%v, %v}, want {42, \"hello\"}", swapped.First, swapped.Second)
		}
	})
}

func TestStack(t *testing.T) {
	t.Run("Push and Pop", func(t *testing.T) {
		stack := NewStack[int]()
		stack.Push(1)
		stack.Push(2)
		stack.Push(3)

		val, ok := stack.Pop()
		if !ok || val != 3 {
			t.Errorf("Pop() = (%v, %v), want (3, true)", val, ok)
		}

		val, ok = stack.Pop()
		if !ok || val != 2 {
			t.Errorf("Pop() = (%v, %v), want (2, true)", val, ok)
		}
	})

	t.Run("Peek", func(t *testing.T) {
		stack := NewStack[int]()
		stack.Push(42)

		val, ok := stack.Peek()
		if !ok || val != 42 {
			t.Errorf("Peek() = (%v, %v), want (42, true)", val, ok)
		}

		// Peek shouldn't remove the item
		val, ok = stack.Peek()
		if !ok || val != 42 {
			t.Errorf("Peek() again = (%v, %v), want (42, true)", val, ok)
		}
	})

	t.Run("Pop empty", func(t *testing.T) {
		stack := NewStack[int]()
		_, ok := stack.Pop()
		if ok {
			t.Errorf("Pop() on empty stack should return false")
		}
	})

	t.Run("Len", func(t *testing.T) {
		stack := NewStack[int]()
		if stack.Len() != 0 {
			t.Errorf("Len() = %v, want 0", stack.Len())
		}

		stack.Push(1)
		stack.Push(2)
		if stack.Len() != 2 {
			t.Errorf("Len() = %v, want 2", stack.Len())
		}

		stack.Pop()
		if stack.Len() != 1 {
			t.Errorf("Len() = %v, want 1", stack.Len())
		}
	})
}
