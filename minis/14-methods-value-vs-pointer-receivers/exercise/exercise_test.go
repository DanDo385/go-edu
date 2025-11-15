package exercise

import (
	"testing"
)

// ============================================================================
// TESTS FOR EXERCISE 1: Basic Receiver Types
// ============================================================================

func TestBankAccountDeposit(t *testing.T) {
	account := BankAccount{balance: 0}

	account.Deposit(100)
	if account.Balance() != 100 {
		t.Errorf("After depositing 100, expected balance 100, got %d", account.Balance())
	}

	account.Deposit(50)
	if account.Balance() != 150 {
		t.Errorf("After depositing 50 more, expected balance 150, got %d", account.Balance())
	}
}

func TestBankAccountWithdraw(t *testing.T) {
	account := BankAccount{balance: 200}

	account.Withdraw(50)
	if account.Balance() != 150 {
		t.Errorf("After withdrawing 50, expected balance 150, got %d", account.Balance())
	}

	account.Withdraw(100)
	if account.Balance() != 50 {
		t.Errorf("After withdrawing 100 more, expected balance 50, got %d", account.Balance())
	}
}

func TestBankAccountChaining(t *testing.T) {
	account := BankAccount{balance: 0}

	account.Deposit(100)
	account.Withdraw(30)
	account.Deposit(50)
	account.Withdraw(20)

	expected := 100
	if account.Balance() != expected {
		t.Errorf("Expected balance %d, got %d", expected, account.Balance())
	}
}

// ============================================================================
// TESTS FOR EXERCISE 2: Interface Satisfaction
// ============================================================================

func TestRectangleArea(t *testing.T) {
	tests := []struct {
		name     string
		rect     Rectangle
		expected float64
	}{
		{"square", Rectangle{Width: 5, Height: 5}, 25.0},
		{"rectangle", Rectangle{Width: 10, Height: 3}, 30.0},
		{"zero area", Rectangle{Width: 0, Height: 5}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test value receiver
			area := tt.rect.Area()
			if area != tt.expected {
				t.Errorf("Expected area %v, got %v", tt.expected, area)
			}

			// Test that value satisfies interface
			var _ Shape = tt.rect

			// Test that pointer also satisfies interface
			var _ Shape = &tt.rect
		})
	}
}

func TestCircleArea(t *testing.T) {
	c := Circle{Radius: 10}

	expected := 3.14159 * 10 * 10
	area := c.Area()

	if area < expected-0.01 || area > expected+0.01 {
		t.Errorf("Expected area ≈ %v, got %v", expected, area)
	}

	// This should work: *Circle satisfies Shape
	var _ Shape = &c

	// This should NOT compile (if Circle has pointer receiver):
	// Uncomment to verify:
	// var _ Shape = c  // Should be compile error if implementation is correct
}

func TestTotalArea(t *testing.T) {
	shapes := []Shape{
		&Rectangle{Width: 5, Height: 4},   // 20
		&Rectangle{Width: 3, Height: 2},   // 6
		&Circle{Radius: 1},                // π ≈ 3.14159
	}

	total := TotalArea(shapes)
	expected := 20.0 + 6.0 + 3.14159

	if total < expected-0.01 || total > expected+0.01 {
		t.Errorf("Expected total area ≈ %v, got %v", expected, total)
	}
}

// ============================================================================
// TESTS FOR EXERCISE 3: Nil Receiver Safety
// ============================================================================

func TestStringListAppend(t *testing.T) {
	t.Run("append to nil list", func(t *testing.T) {
		var list *StringList
		list = list.Append("first")

		if list == nil {
			t.Fatal("Expected non-nil list after append")
		}
		if list.First() != "first" {
			t.Errorf("Expected first element 'first', got '%s'", list.First())
		}
	})

	t.Run("append multiple elements", func(t *testing.T) {
		var list *StringList
		list = list.Append("a")
		list = list.Append("b")
		list = list.Append("c")

		if !list.Contains("a") || !list.Contains("b") || !list.Contains("c") {
			t.Error("List should contain all appended elements")
		}
	})
}

func TestStringListContains(t *testing.T) {
	var list *StringList
	list = list.Append("apple")
	list = list.Append("banana")
	list = list.Append("cherry")

	tests := []struct {
		value    string
		expected bool
	}{
		{"apple", true},
		{"banana", true},
		{"cherry", true},
		{"grape", false},
		{"", false},
	}

	for _, tt := range tests {
		result := list.Contains(tt.value)
		if result != tt.expected {
			t.Errorf("Contains(%q) = %v, expected %v", tt.value, result, tt.expected)
		}
	}
}

func TestStringListNilSafety(t *testing.T) {
	var list *StringList  // nil

	// These should not panic
	if list.Contains("anything") {
		t.Error("Nil list should not contain anything")
	}

	if list.First() != "" {
		t.Error("Nil list First() should return empty string")
	}
}

func TestStringListFirst(t *testing.T) {
	var list *StringList
	list = list.Append("first")
	list = list.Append("second")

	if list.First() != "first" {
		t.Errorf("Expected 'first', got '%s'", list.First())
	}
}

// ============================================================================
// TESTS FOR EXERCISE 4: Performance-Aware Design
// ============================================================================

func TestSmallConfigValidate(t *testing.T) {
	tests := []struct {
		name     string
		config   SmallConfig
		expected bool
	}{
		{"valid config", SmallConfig{ID: 1, Name: "test"}, true},
		{"zero ID", SmallConfig{ID: 0, Name: "test"}, false},
		{"empty name", SmallConfig{ID: 1, Name: ""}, false},
		{"both invalid", SmallConfig{ID: 0, Name: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.Validate()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLargeConfigSum(t *testing.T) {
	config := LargeConfig{}

	// Fill with known values
	for i := 0; i < 1000; i++ {
		config.Data[i] = i
	}

	expected := 0
	for i := 0; i < 1000; i++ {
		expected += i
	}

	result := config.Sum()
	if result != expected {
		t.Errorf("Expected sum %d, got %d", expected, result)
	}
}

func TestLargeConfigUsesPointerReceiver(t *testing.T) {
	// This is more of a documentation test
	// In real code, you'd use benchmarks to verify performance
	config := &LargeConfig{}
	for i := 0; i < 1000; i++ {
		config.Data[i] = 1
	}

	// Should be efficient (no copying 8000 bytes)
	sum := config.Sum()
	if sum != 1000 {
		t.Errorf("Expected sum 1000, got %d", sum)
	}
}

// ============================================================================
// TESTS FOR EXERCISE 5: Consistency and Best Practices
// ============================================================================

func TestUserSetters(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 25}

	user.SetName("Bob")
	if user.GetName() != "Bob" {
		t.Errorf("Expected name 'Bob', got '%s'", user.GetName())
	}

	user.SetEmail("bob@example.com")
	if user.Email != "bob@example.com" {
		t.Errorf("Expected email 'bob@example.com', got '%s'", user.Email)
	}
}

func TestUserIsAdult(t *testing.T) {
	tests := []struct {
		name     string
		age      int
		expected bool
	}{
		{"adult", 18, true},
		{"older adult", 30, true},
		{"minor", 17, false},
		{"child", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{Age: tt.age}
			result := user.IsAdult()
			if result != tt.expected {
				t.Errorf("Age %d: expected IsAdult() = %v, got %v", tt.age, tt.expected, result)
			}
		})
	}
}

func TestUserConsistency(t *testing.T) {
	// All methods should use the same receiver type
	// This test verifies that mutations persist
	user := User{Name: "Alice", Email: "alice@example.com", Age: 17}

	user.SetName("Bob")
	user.SetEmail("bob@example.com")

	if user.GetName() != "Bob" {
		t.Error("SetName should persist (check receiver type)")
	}

	if user.Email != "bob@example.com" {
		t.Error("SetEmail should persist (check receiver type)")
	}
}

// ============================================================================
// TESTS FOR EXERCISE 6: Method Set Understanding
// ============================================================================

func TestPointEquals(t *testing.T) {
	p1 := Point{X: 3, Y: 4}
	p2 := Point{X: 3, Y: 4}
	p3 := Point{X: 5, Y: 6}

	if !p1.Equals(p2) {
		t.Error("Equal points should return true")
	}

	if p1.Equals(p3) {
		t.Error("Different points should return false")
	}
}

func TestPointImplementsComparable(t *testing.T) {
	p := Point{X: 1, Y: 2}

	// Both value and pointer should satisfy Comparable
	var _ Comparable = p
	var _ Comparable = &p
}

// ============================================================================
// TESTS FOR EXERCISE 7: Map Element Challenge
// ============================================================================

func TestSafeCounterMapIncrement(t *testing.T) {
	m := NewSafeCounterMap()

	m.Increment("a")
	if m.Get("a") != 1 {
		t.Errorf("After one increment, expected 1, got %d", m.Get("a"))
	}

	m.Increment("a")
	m.Increment("a")
	if m.Get("a") != 3 {
		t.Errorf("After three increments, expected 3, got %d", m.Get("a"))
	}
}

func TestSafeCounterMapMultipleKeys(t *testing.T) {
	m := NewSafeCounterMap()

	m.Increment("x")
	m.Increment("y")
	m.Increment("x")
	m.Increment("z")
	m.Increment("x")

	if m.Get("x") != 3 {
		t.Errorf("Expected x=3, got %d", m.Get("x"))
	}
	if m.Get("y") != 1 {
		t.Errorf("Expected y=1, got %d", m.Get("y"))
	}
	if m.Get("z") != 1 {
		t.Errorf("Expected z=1, got %d", m.Get("z"))
	}
}

func TestSafeCounterMapGetNonexistent(t *testing.T) {
	m := NewSafeCounterMap()

	if m.Get("nonexistent") != 0 {
		t.Error("Nonexistent key should return 0")
	}
}

// ============================================================================
// BENCHMARK TESTS (Optional: Run with go test -bench=.)
// ============================================================================

func BenchmarkSmallConfigValue(b *testing.B) {
	config := SmallConfig{ID: 1, Name: "test"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkLargeConfigPointer(b *testing.B) {
	config := &LargeConfig{}
	for i := 0; i < 1000; i++ {
		config.Data[i] = i
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = config.Sum()
	}
}

// This benchmark would be slow if Sum uses value receiver
// func BenchmarkLargeConfigValue(b *testing.B) {
//     config := LargeConfig{}
//     for i := 0; i < 1000; i++ {
//         config.Data[i] = i
//     }
//     b.ResetTimer()
//
//     for i := 0; i < b.N; i++ {
//         _ = config.Sum()  // Would copy 8000 bytes each iteration!
//     }
// }
