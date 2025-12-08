//go:build solution
// +build solution

package exercise

// This file contains the solutions to all exercises.
// Try to solve the exercises yourself first before looking at these solutions!

// ============================================================================
// SOLUTION 1: Basic Receiver Types
// ============================================================================

// The key insight: Methods that modify the receiver MUST use pointer receivers.
// All methods should be consistent - if one uses *, all should use *.

// Deposit adds money to the account.
// SOLUTION: Use pointer receiver because we're modifying the balance.
func (b *BankAccount) DepositSolution(amount int) {
	b.balance += amount
}

// Balance returns the current balance.
// SOLUTION: Use pointer receiver for consistency (even though read-only).
func (b *BankAccount) BalanceSolution() int {
	return b.balance
}

// Withdraw subtracts money from the account.
// SOLUTION: Use pointer receiver because we're modifying the balance.
func (b *BankAccount) WithdrawSolution(amount int) {
	b.balance -= amount
}

// ============================================================================
// SOLUTION 2: Interface Satisfaction
// ============================================================================

// Rectangle.Area() should use VALUE receiver.
// EXPLANATION: This allows both Rectangle and *Rectangle to satisfy Shape.
// Small struct (16 bytes), read-only operation, no modification needed.
func (r Rectangle) AreaSolution() float64 {
	return r.Width * r.Height
}

// Circle.Area() should use POINTER receiver.
// EXPLANATION: This restricts interface satisfaction to *Circle only.
// Demonstrates that only pointer types satisfy interfaces with pointer receivers.
func (c *Circle) AreaSolution() float64 {
	return 3.14159 * c.Radius * c.Radius
}

// TotalArea calculates the total area of multiple shapes.
// SOLUTION: Iterate and sum all areas.
func TotalAreaSolution(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// ============================================================================
// SOLUTION 3: Nil Receiver Safety
// ============================================================================

// Append adds a value to the end of the list.
// SOLUTION: Check for nil receiver and handle it specially.
func (l *StringList) AppendSolution(value string) *StringList {
	// If the receiver is nil, create a new list
	if l == nil {
		return &StringList{value: value, next: nil}
	}

	// If this is the last node, append here
	if l.next == nil {
		l.next = &StringList{value: value, next: nil}
		return l
	}

	// Otherwise, recursively append to the next node
	l.next = l.next.AppendSolution(value)
	return l
}

// Contains checks if a value exists in the list.
// SOLUTION: Traverse the list, checking each node.
func (l *StringList) ContainsSolution(value string) bool {
	// Nil list doesn't contain anything
	if l == nil {
		return false
	}

	// Check this node
	if l.value == value {
		return true
	}

	// Check the rest of the list recursively
	return l.next.ContainsSolution(value)
}

// First returns the first element in the list.
// SOLUTION: Handle nil safely.
func (l *StringList) FirstSolution() string {
	if l == nil {
		return ""
	}
	return l.value
}

// ============================================================================
// SOLUTION 4: Performance-Aware Design
// ============================================================================

// SmallConfig.Validate() uses VALUE receiver.
// EXPLANATION: Small struct (~24 bytes), read-only, copying is cheap.
func (c SmallConfig) ValidateSolution() bool {
	return c.ID > 0 && c.Name != ""
}

// LargeConfig.Sum() uses POINTER receiver.
// EXPLANATION: Large struct (8000 bytes), avoid copying overhead.
func (l *LargeConfig) SumSolution() int {
	total := 0
	for _, v := range l.Data {
		total += v
	}
	return total
}

// ============================================================================
// SOLUTION 5: Consistency and Best Practices
// ============================================================================

// All User methods use POINTER receivers for consistency.
// EXPLANATION: Since SetName and SetEmail need to modify, we use * for all methods.

func (u *User) SetNameSolution(name string) {
	u.Name = name
}

func (u *User) SetEmailSolution(email string) {
	u.Email = email
}

func (u *User) GetNameSolution() string {
	return u.Name
}

func (u *User) IsAdultSolution() bool {
	return u.Age >= 18
}

// ============================================================================
// SOLUTION 6: Method Set Understanding
// ============================================================================

// Point.Equals() uses VALUE receiver.
// EXPLANATION: This allows both Point and *Point to satisfy Comparable.
// Also, Point is small (16 bytes) and we don't modify it.
func (p Point) EqualsSolution(other Comparable) bool {
	// Type assert to Point
	otherPoint, ok := other.(Point)
	if !ok {
		// Also try *Point
		otherPointPtr, ok := other.(*Point)
		if !ok {
			return false
		}
		otherPoint = *otherPointPtr
	}

	return p.X == otherPoint.X && p.Y == otherPoint.Y
}

// ============================================================================
// SOLUTION 7: Map Element Challenge
// ============================================================================

// NewSafeCounterMap creates a new SafeCounterMap.
// SOLUTION: Initialize the map so it's ready to use.
func NewSafeCounterMapSolution() SafeCounterMap {
	return SafeCounterMap{
		counters: make(map[string]int),
	}
}

// Increment increments the counter for a given key.
// SOLUTION: Use pointer receiver to modify the map.
func (m *SafeCounterMap) IncrementSolution(key string) {
	// Maps are reference types, so we can modify them through the pointer
	m.counters[key]++
}

// Get returns the counter value for a given key.
// SOLUTION: Read from the map (returns 0 if key doesn't exist).
// Use pointer receiver for consistency.
func (m *SafeCounterMap) GetSolution(key string) int {
	return m.counters[key]
}

// ============================================================================
// ALTERNATIVE SOLUTION: Iterative Append (More Efficient)
// ============================================================================

// AppendIterative is a more efficient version of Append using iteration.
// EXPLANATION: Recursion can cause stack overflow for very long lists.
func (l *StringList) AppendIterative(value string) *StringList {
	newNode := &StringList{value: value, next: nil}

	// If the list is nil, the new node becomes the head
	if l == nil {
		return newNode
	}

	// Find the last node
	current := l
	for current.next != nil {
		current = current.next
	}

	// Append the new node
	current.next = newNode
	return l
}

// ============================================================================
// KEY INSIGHTS FROM SOLUTIONS
// ============================================================================

/*
INSIGHT 1: Pointer Receivers for Mutation
- Any method that modifies the receiver MUST use a pointer receiver
- Otherwise, it modifies a copy that's thrown away

INSIGHT 2: Consistency is Critical
- If one method uses a pointer receiver, ALL methods should
- Makes the API predictable and avoids interface satisfaction issues

INSIGHT 3: Interface Satisfaction Rules
- Type T satisfies interfaces requiring methods with receiver T
- Type *T satisfies interfaces requiring methods with receiver T or *T
- Choose receiver type based on which types you want to satisfy interfaces

INSIGHT 4: Performance Matters for Large Types
- Small types (< 64 bytes): Value receiver is fine
- Large types (> 64 bytes): Use pointer receiver to avoid copying
- The difference can be 10x-100x in performance

INSIGHT 5: Nil Safety
- Methods with pointer receivers can be called on nil
- Always check for nil before dereferencing fields
- This enables powerful nil-safe patterns

INSIGHT 6: Map Elements Aren't Addressable
- You can't call pointer-receiver methods on map elements directly
- Solutions: extract-modify-store, or use pointers in the map

INSIGHT 7: Think About Semantics
- Value types (numbers, IDs): Use value receivers
- Entity types (users, accounts): Use pointer receivers
- The receiver type communicates intent to API users
*/
