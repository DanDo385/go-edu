package main

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/example/go-10x-minis/minis/46-generics-map-reduce/exercise"
)

func main() {
	fmt.Println("=== Go Generics & Map-Reduce Demo ===")
	fmt.Println()

	demo1_BasicGenerics()
	demo2_MapOperations()
	demo3_FilterOperations()
	demo4_ReduceOperations()
	demo5_FlatMapOperations()
	demo6_GenericDataStructures()
	demo7_ParallelMapReduce()
	demo8_RealWorldExample()
}

// Demo 1: Basic Generic Functions
func demo1_BasicGenerics() {
	fmt.Println("--- Demo 1: Basic Generic Functions ---")

	// Generic Identity function
	intVal := exercise.Identity(42)
	strVal := exercise.Identity("hello")
	floatVal := exercise.Identity(3.14)

	fmt.Printf("Identity[int](42) = %v\n", intVal)
	fmt.Printf("Identity[string](\"hello\") = %v\n", strVal)
	fmt.Printf("Identity[float64](3.14) = %v\n", floatVal)

	// Generic Contains
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("\nContains(numbers, 3) = %v\n", exercise.Contains(numbers, 3))
	fmt.Printf("Contains(numbers, 10) = %v\n", exercise.Contains(numbers, 10))

	// Generic Reverse
	reversed := exercise.Reverse(numbers)
	fmt.Printf("\nReverse([1,2,3,4,5]) = %v\n", reversed)

	words := []string{"hello", "world", "from", "Go"}
	reversedWords := exercise.Reverse(words)
	fmt.Printf("Reverse(%v) = %v\n", words, reversedWords)

	fmt.Println()
}

// Demo 2: Map Operations
func demo2_MapOperations() {
	fmt.Println("--- Demo 2: Map Operations ---")

	numbers := []int{1, 2, 3, 4, 5}

	// Map: int -> int (double)
	doubled := exercise.Map(numbers, func(x int) int {
		return x * 2
	})
	fmt.Printf("Map(double): %v -> %v\n", numbers, doubled)

	// Map: int -> string (transform type)
	asStrings := exercise.Map(numbers, func(x int) string {
		return fmt.Sprintf("#%d", x)
	})
	fmt.Printf("Map(to string): %v -> %v\n", numbers, asStrings)

	// Map: int -> float64 (square root)
	sqrts := exercise.Map(numbers, func(x int) float64 {
		return math.Sqrt(float64(x))
	})
	fmt.Printf("Map(sqrt): %v -> %.2f\n", numbers, sqrts)

	// Chaining maps
	result := exercise.Map(
		exercise.Map(numbers, func(x int) int { return x * 2 }),
		func(x int) int { return x + 1 },
	)
	fmt.Printf("Map(double then +1): %v -> %v\n", numbers, result)

	fmt.Println()
}

// Demo 3: Filter Operations
func demo3_FilterOperations() {
	fmt.Println("--- Demo 3: Filter Operations ---")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter evens
	evens := exercise.Filter(numbers, func(x int) bool {
		return x%2 == 0
	})
	fmt.Printf("Filter(even): %v -> %v\n", numbers, evens)

	// Filter greater than 5
	greaterThan5 := exercise.Filter(numbers, func(x int) bool {
		return x > 5
	})
	fmt.Printf("Filter(>5): %v -> %v\n", numbers, greaterThan5)

	// Filter strings by length
	words := []string{"a", "ab", "abc", "abcd", "abcde"}
	longWords := exercise.Filter(words, func(s string) bool {
		return len(s) >= 3
	})
	fmt.Printf("Filter(len>=3): %v -> %v\n", words, longWords)

	// Combining filter and map
	evenSquares := exercise.Map(
		exercise.Filter(numbers, func(x int) bool { return x%2 == 0 }),
		func(x int) int { return x * x },
	)
	fmt.Printf("Filter(even) then Map(square): %v -> %v\n", numbers, evenSquares)

	fmt.Println()
}

// Demo 4: Reduce Operations
func demo4_ReduceOperations() {
	fmt.Println("--- Demo 4: Reduce Operations ---")

	numbers := []int{1, 2, 3, 4, 5}

	// Sum
	sum := exercise.Reduce(numbers, 0, func(acc, x int) int {
		return acc + x
	})
	fmt.Printf("Reduce(sum): %v -> %d\n", numbers, sum)

	// Product
	product := exercise.Reduce(numbers, 1, func(acc, x int) int {
		return acc * x
	})
	fmt.Printf("Reduce(product): %v -> %d\n", numbers, product)

	// Max
	max := exercise.Reduce(numbers, numbers[0], func(acc, x int) int {
		if x > acc {
			return x
		}
		return acc
	})
	fmt.Printf("Reduce(max): %v -> %d\n", numbers, max)

	// Concatenate strings
	words := []string{"Hello", "world", "from", "Go"}
	sentence := exercise.Reduce(words, "", func(acc, word string) string {
		if acc == "" {
			return word
		}
		return acc + " " + word
	})
	fmt.Printf("Reduce(concat): %v -> \"%s\"\n", words, sentence)

	// Count occurrences (reduce to map)
	items := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	counts := exercise.Reduce(items, make(map[string]int), func(acc map[string]int, item string) map[string]int {
		acc[item]++
		return acc
	})
	fmt.Printf("Reduce(count): %v -> %v\n", items, counts)

	fmt.Println()
}

// Demo 5: FlatMap Operations
func demo5_FlatMapOperations() {
	fmt.Println("--- Demo 5: FlatMap Operations ---")

	// Split strings into words
	sentences := []string{"hello world", "go is great"}
	words := exercise.FlatMap(sentences, func(s string) []string {
		return strings.Fields(s)
	})
	fmt.Printf("FlatMap(split): %v -> %v\n", sentences, words)

	// Generate ranges
	lengths := []int{2, 3, 1}
	ranges := exercise.FlatMap(lengths, func(n int) []int {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = i
		}
		return result
	})
	fmt.Printf("FlatMap(range): %v -> %v\n", lengths, ranges)

	// Duplicate elements
	numbers := []int{1, 2, 3}
	duplicated := exercise.FlatMap(numbers, func(x int) []int {
		return []int{x, x}
	})
	fmt.Printf("FlatMap(duplicate): %v -> %v\n", numbers, duplicated)

	fmt.Println()
}

// Demo 6: Generic Data Structures
func demo6_GenericDataStructures() {
	fmt.Println("--- Demo 6: Generic Data Structures ---")

	// Optional[T]
	fmt.Println("Optional[int]:")
	opt1 := exercise.Some(42)
	opt2 := exercise.None[int]()

	if val, ok := opt1.Get(); ok {
		fmt.Printf("  Some(42).Get() = %d\n", val)
	}

	if _, ok := opt2.Get(); !ok {
		fmt.Println("  None().Get() = (no value)")
	}

	fmt.Printf("  Some(42).OrElse(0) = %d\n", opt1.OrElse(0))
	fmt.Printf("  None().OrElse(99) = %d\n", opt2.OrElse(99))

	// Result[T, E]
	fmt.Println("\nResult[int, string]:")
	ok := exercise.Ok[int, string](42)
	err := exercise.Err[int, string]("something went wrong")

	if val, _, isOk := ok.Unwrap(); isOk {
		fmt.Printf("  Ok(42).Unwrap() = %d\n", val)
	}

	if _, errMsg, isOk := err.Unwrap(); !isOk {
		fmt.Printf("  Err(\"...\").Unwrap() = error: %s\n", errMsg)
	}

	// Pair[A, B]
	fmt.Println("\nPair[string, int]:")
	pair := exercise.MakePair("answer", 42)
	fmt.Printf("  MakePair(\"answer\", 42) = {%v, %v}\n", pair.First, pair.Second)

	swapped := pair.Swap()
	fmt.Printf("  Swap() = {%v, %v}\n", swapped.First, swapped.Second)

	// Stack[T]
	fmt.Println("\nStack[int]:")
	stack := exercise.NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	fmt.Println("  Push(1, 2, 3)")

	if val, ok := stack.Peek(); ok {
		fmt.Printf("  Peek() = %d\n", val)
	}

	if val, ok := stack.Pop(); ok {
		fmt.Printf("  Pop() = %d\n", val)
	}

	if val, ok := stack.Pop(); ok {
		fmt.Printf("  Pop() = %d\n", val)
	}

	fmt.Println()
}

// Demo 7: Parallel Map-Reduce
func demo7_ParallelMapReduce() {
	fmt.Println("--- Demo 7: Parallel Map-Reduce ---")

	numCPU := runtime.NumCPU()
	fmt.Printf("Available CPUs: %d\n\n", numCPU)

	// Create large dataset
	n := 1000000
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = i + 1
	}

	// CPU-intensive operation (simulate work)
	expensiveFunc := func(x int) int {
		result := x
		for i := 0; i < 100; i++ {
			result = (result*result + x) % 1000000
		}
		return result
	}

	// Sequential Map
	fmt.Println("Sequential Map:")
	start := time.Now()
	seqResult := exercise.Map(data[:10000], expensiveFunc) // Smaller dataset for demo
	seqDuration := time.Since(start)
	fmt.Printf("  Processed 10,000 items in %v\n", seqDuration)
	fmt.Printf("  First 5 results: %v\n", seqResult[:5])

	// Parallel Map
	fmt.Println("\nParallel Map:")
	start = time.Now()
	parResult := exercise.ParallelMap(data[:10000], expensiveFunc, numCPU)
	parDuration := time.Since(start)
	fmt.Printf("  Processed 10,000 items in %v\n", parDuration)
	fmt.Printf("  First 5 results: %v\n", parResult[:5])

	if seqDuration > parDuration {
		speedup := float64(seqDuration) / float64(parDuration)
		fmt.Printf("  Speedup: %.2fx\n", speedup)
	}

	// Parallel Reduce (sum)
	fmt.Println("\nParallel Reduce (Sum):")
	numbers := make([]int, 1000000)
	for i := 0; i < len(numbers); i++ {
		numbers[i] = i + 1
	}

	start = time.Now()
	seqSum := exercise.Reduce(numbers, 0, func(acc, x int) int { return acc + x })
	seqDuration = time.Since(start)

	start = time.Now()
	parSum := exercise.ParallelReduce(numbers, 0, func(acc, x int) int { return acc + x }, numCPU)
	parDuration = time.Since(start)

	fmt.Printf("  Sequential sum: %d (took %v)\n", seqSum, seqDuration)
	fmt.Printf("  Parallel sum:   %d (took %v)\n", parSum, parDuration)

	if seqSum == parSum {
		fmt.Println("  ✓ Results match!")
	}

	fmt.Println()
}

// Demo 8: Real-World Example (Log Processing)
func demo8_RealWorldExample() {
	fmt.Println("--- Demo 8: Real-World Example (Log Analysis) ---")

	// Simulated log entries
	type LogEntry struct {
		Timestamp time.Time
		Level     string
		Message   string
	}

	logs := []LogEntry{
		{time.Now(), "INFO", "Server started on port 8080"},
		{time.Now(), "ERROR", "Failed to connect to database connection timeout"},
		{time.Now(), "INFO", "Request received from 192.168.1.1"},
		{time.Now(), "ERROR", "Database query failed connection lost"},
		{time.Now(), "WARN", "Slow query detected took 2s"},
		{time.Now(), "ERROR", "Failed to write to disk disk full"},
		{time.Now(), "INFO", "Request completed in 45ms"},
		{time.Now(), "ERROR", "Connection timeout to external service"},
	}

	fmt.Printf("Total logs: %d\n\n", len(logs))

	// Pipeline 1: Count error logs
	errorCount := exercise.Reduce(
		exercise.Filter(logs, func(log LogEntry) bool {
			return log.Level == "ERROR"
		}),
		0,
		func(count int, log LogEntry) int {
			return count + 1
		},
	)
	fmt.Printf("Error logs: %d\n", errorCount)

	// Pipeline 2: Extract error messages
	errorMessages := exercise.Map(
		exercise.Filter(logs, func(log LogEntry) bool {
			return log.Level == "ERROR"
		}),
		func(log LogEntry) string {
			return log.Message
		},
	)
	fmt.Println("\nError messages:")
	for _, msg := range errorMessages {
		fmt.Printf("  - %s\n", msg)
	}

	// Pipeline 3: Word frequency in error messages
	words := exercise.FlatMap(errorMessages, func(msg string) []string {
		return strings.Fields(strings.ToLower(msg))
	})

	wordCounts := exercise.Reduce(words, make(map[string]int), func(acc map[string]int, word string) map[string]int {
		acc[word]++
		return acc
	})

	fmt.Println("\nWord frequency in errors:")
	for word, count := range wordCounts {
		if count > 1 {
			fmt.Printf("  %s: %d\n", word, count)
		}
	}

	// Pipeline 4: Group by level
	logsByLevel := exercise.Reduce(logs, make(map[string][]LogEntry), func(acc map[string][]LogEntry, log LogEntry) map[string][]LogEntry {
		acc[log.Level] = append(acc[log.Level], log)
		return acc
	})

	fmt.Println("\nLogs by level:")
	for level, entries := range logsByLevel {
		fmt.Printf("  %s: %d logs\n", level, len(entries))
	}

	fmt.Println()
}
