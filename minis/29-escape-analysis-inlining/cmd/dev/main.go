package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// ESCAPE ANALYSIS EXAMPLES
// Run with: go build -gcflags='-m' cmd/escape-demo/main.go
// ============================================================================

// Example 1: Stays on stack (value returned)
func createValue() int {
	x := 42
	return x // x stays on stack, value copied to caller
}

// Example 2: Escapes to heap (pointer returned)
func createPointer() *int {
	x := 42
	return &x // x escapes: pointer returned to caller
}

// Example 3: Stays on stack (used locally)
func sumLocal() int {
	nums := []int{1, 2, 3, 4, 5}
	total := 0
	for _, n := range nums {
		total += n
	}
	return total // nums used only locally, can be stack-allocated
}

// Example 4: Escapes to heap (returned slice)
func createSlice() []int {
	s := []int{1, 2, 3, 4, 5}
	return s // s's backing array escapes to heap
}

// Example 5: Escapes via interface{}
func escapeViaInterface(x int) {
	fmt.Println(x) // x escapes because fmt.Println takes interface{}
}

// Example 6: Large allocation escapes
func largeAllocation() {
	// Large arrays typically escape to heap
	var big [10000]int
	big[0] = 1
	_ = big
}

// Example 7: Closure captures by reference
func makeClosure() func() int {
	count := 0
	return func() int {
		count++ // count escapes: captured by closure
		return count
	}
}

// ============================================================================
// INLINING EXAMPLES
// Run with: go build -gcflags='-m' cmd/escape-demo/main.go
// ============================================================================

// Simple function: WILL be inlined
func add(a, b int) int {
	return a + b
}

// Simple function: WILL be inlined
func multiply(a, b int) int {
	return a * b
}

// Complex function: Will NOT be inlined
func complexCalculation(a, b, c, d int) int {
	if a > b {
		for i := 0; i < c; i++ {
			a += b * d
		}
	} else {
		for i := 0; i < d; i++ {
			b += a * c
		}
	}
	return a + b
}

// ============================================================================
// OPTIMIZATION EXAMPLES
// ============================================================================

// Token represents a parsed token
type Token struct {
	Type  string
	Value string
}

// parseTokensNaive: Inefficient version (many allocations)
func parseTokensNaive(input string) []*Token {
	var tokens []*Token // No pre-allocation
	for _, word := range strings.Fields(input) {
		token := &Token{ // Each token escapes to heap
			Type:  "WORD",
			Value: word,
		}
		tokens = append(tokens, token)
	}
	return tokens
}

// parseTokensOptimized: Efficient version (fewer allocations)
func parseTokensOptimized(input string) []Token {
	words := strings.Fields(input)
	tokens := make([]Token, 0, len(words)) // Pre-allocated

	for _, word := range words {
		tokens = append(tokens, Token{ // Values, not pointers
			Type:  "WORD",
			Value: word,
		})
	}
	return tokens
}

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// DistanceValueReceiver: Efficient for small structs
func (p Point) DistanceValueReceiver() float64 {
	return p.X*p.X + p.Y*p.Y
}

// DistancePointerReceiver: Less efficient for small structs
func (p *Point) DistancePointerReceiver() float64 {
	return p.X*p.X + p.Y*p.Y
}

// ============================================================================
// DEMONSTRATION FUNCTIONS
// ============================================================================

func demonstrateEscapeAnalysis() {
	fmt.Println("=== Escape Analysis Demonstrations ===\n")

	// Example 1: Value vs Pointer
	fmt.Println("1. Value (stack) vs Pointer (heap):")
	val := createValue()
	ptr := createPointer()
	fmt.Printf("   Value: %d (allocated on stack)\n", val)
	fmt.Printf("   Pointer: %d (allocated on heap)\n", *ptr)

	// Example 2: Local slice
	fmt.Println("\n2. Local slice (may stay on stack):")
	sum := sumLocal()
	fmt.Printf("   Sum: %d\n", sum)

	// Example 3: Returned slice
	fmt.Println("\n3. Returned slice (escapes to heap):")
	s := createSlice()
	fmt.Printf("   Slice: %v\n", s)

	// Example 4: Closure
	fmt.Println("\n4. Closure capturing variable:")
	counter := makeClosure()
	fmt.Printf("   Count 1: %d\n", counter())
	fmt.Printf("   Count 2: %d\n", counter())
	fmt.Printf("   Count 3: %d\n", counter())

	fmt.Println()
}

func demonstrateInlining() {
	fmt.Println("=== Inlining Demonstrations ===\n")

	// Simple functions (will be inlined)
	fmt.Println("1. Simple functions (inlined):")
	a, b := 10, 20
	sum := add(a, b)
	product := multiply(a, b)
	fmt.Printf("   %d + %d = %d\n", a, b, sum)
	fmt.Printf("   %d * %d = %d\n", a, b, product)

	// Complex function (will NOT be inlined)
	fmt.Println("\n2. Complex function (not inlined):")
	result := complexCalculation(5, 10, 3, 2)
	fmt.Printf("   Complex result: %d\n", result)

	fmt.Println()
}

func demonstrateOptimizations() {
	fmt.Println("=== Optimization Demonstrations ===\n")

	input := "the quick brown fox jumps over the lazy dog"

	// Naive parsing
	fmt.Println("1. Naive parsing (many allocations):")
	tokensNaive := parseTokensNaive(input)
	fmt.Printf("   Parsed %d tokens (pointer-based)\n", len(tokensNaive))
	fmt.Printf("   First token: %s '%s'\n", tokensNaive[0].Type, tokensNaive[0].Value)

	// Optimized parsing
	fmt.Println("\n2. Optimized parsing (fewer allocations):")
	tokensOpt := parseTokensOptimized(input)
	fmt.Printf("   Parsed %d tokens (value-based)\n", len(tokensOpt))
	fmt.Printf("   First token: %s '%s'\n", tokensOpt[0].Type, tokensOpt[0].Value)

	// Value vs pointer receivers
	fmt.Println("\n3. Value vs pointer receivers:")
	p := Point{X: 3.0, Y: 4.0}
	distVal := p.DistanceValueReceiver()
	distPtr := (&p).DistancePointerReceiver()
	fmt.Printf("   Value receiver: %.2f\n", distVal)
	fmt.Printf("   Pointer receiver: %.2f\n", distPtr)

	fmt.Println()
}

func demonstrateCompilerFlags() {
	fmt.Println("=== Compiler Flags Guide ===\n")

	fmt.Println("To see escape analysis decisions:")
	fmt.Println("  go build -gcflags='-m' cmd/escape-demo/main.go")
	fmt.Println()

	fmt.Println("To see more verbose escape analysis:")
	fmt.Println("  go build -gcflags='-m -m' cmd/escape-demo/main.go")
	fmt.Println()

	fmt.Println("To see inlining decisions:")
	fmt.Println("  go build -gcflags='-m' cmd/escape-demo/main.go 2>&1 | grep inline")
	fmt.Println()

	fmt.Println("To disable inlining and compare:")
	fmt.Println("  go build -gcflags='-l -m' cmd/escape-demo/main.go")
	fmt.Println()

	fmt.Println("To see generated assembly:")
	fmt.Println("  go build -gcflags='-S' cmd/escape-demo/main.go 2>&1 | less")
	fmt.Println()

	fmt.Println("To check bounds check elimination:")
	fmt.Println("  go build -gcflags='-d=ssa/check_bce/debug=1' cmd/escape-demo/main.go")
	fmt.Println()
}

func demonstrateMemoryBehavior() {
	fmt.Println("=== Memory Allocation Behavior ===\n")

	// Stack allocation example
	fmt.Println("1. Stack allocation (fast):")
	fmt.Println("   func stackExample() int {")
	fmt.Println("       x := 42")
	fmt.Println("       return x  // x stays on stack")
	fmt.Println("   }")
	fmt.Println("   → Allocation: ~0.5 ns (just stack pointer move)")
	fmt.Println()

	// Heap allocation example
	fmt.Println("2. Heap allocation (slower):")
	fmt.Println("   func heapExample() *int {")
	fmt.Println("       x := 42")
	fmt.Println("       return &x  // x escapes to heap")
	fmt.Println("   }")
	fmt.Println("   → Allocation: ~25 ns (memory allocator + GC tracking)")
	fmt.Println()

	// Performance comparison
	fmt.Println("3. Performance impact:")
	fmt.Println("   - Stack allocation: ~50x faster")
	fmt.Println("   - No GC pressure from stack allocations")
	fmt.Println("   - Better CPU cache locality on stack")
	fmt.Println()
}

func demonstrateCommonPatterns() {
	fmt.Println("=== Common Optimization Patterns ===\n")

	fmt.Println("1. Avoid unnecessary pointers:")
	fmt.Println("   BAD:  func process() *Result { r := Result{}; return &r }")
	fmt.Println("   GOOD: func process() Result { r := Result{}; return r }")
	fmt.Println()

	fmt.Println("2. Pre-allocate slices:")
	fmt.Println("   BAD:  var s []int; for ... { s = append(s, x) }")
	fmt.Println("   GOOD: s := make([]int, 0, knownSize)")
	fmt.Println()

	fmt.Println("3. Avoid interface{} in hot paths:")
	fmt.Println("   BAD:  fmt.Println(x)  // x escapes to interface{}")
	fmt.Println("   GOOD: Use custom logging with concrete types")
	fmt.Println()

	fmt.Println("4. Use value receivers for small structs:")
	fmt.Println("   BAD:  func (p *Point) Distance() { ... }  // for 16-byte struct")
	fmt.Println("   GOOD: func (p Point) Distance() { ... }   // pass by value")
	fmt.Println()

	fmt.Println("5. Reuse buffers:")
	fmt.Println("   BAD:  for ... { buf := make([]byte, 1024); ... }")
	fmt.Println("   GOOD: buf := make([]byte, 1024); for ... { clear(buf); ... }")
	fmt.Println()
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  Go Escape Analysis & Inlining Demonstration              ║")
	fmt.Println("║  Project 29: Compiler Optimizations Deep Dive             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	demonstrateEscapeAnalysis()
	demonstrateInlining()
	demonstrateOptimizations()
	demonstrateMemoryBehavior()
	demonstrateCommonPatterns()
	demonstrateCompilerFlags()

	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("Next Steps:")
	fmt.Println("1. Run: go build -gcflags='-m' cmd/escape-demo/main.go")
	fmt.Println("2. Complete exercises in exercise/exercise.go")
	fmt.Println("3. Run benchmarks: cd exercise && go test -bench=. -benchmem")
	fmt.Println("════════════════════════════════════════════════════════════")
}
