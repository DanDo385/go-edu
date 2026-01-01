//go:build reference

/*
Problem: Understanding method receivers in Go (value vs pointer)

Requirements:
1. Choose correct receiver type for mutation
2. Understand interface satisfaction rules
3. Handle nil receivers safely
4. Optimize for performance (large structs)
5. Maintain API consistency

Data Structure:
- Value receiver: Operates on copy (8-64 bytes typical)
- Pointer receiver: Operates on original (8 bytes pointer)
- Method set: T has methods with receiver T, *T has both T and *T

Algorithm: Receiver Selection
- Mutation needed: Use pointer receiver
- Large struct (>64 bytes): Use pointer receiver
- Small immutable value: Use value receiver
- Interface satisfaction: Consider both T and *T

Why receiver type matters:
- Value receiver: Safe copy, can't modify original
- Pointer receiver: Can modify, more efficient for large types
- Mixed receivers cause interface satisfaction issues
*/

package methodsvaluevspointerreceivers

// DepositSolution adds money to account.
// BREAKPOINT: Set breakpoint here to trace mutation
// DEBUG: Watch 'b.balance' before and after
func (b *BankAccount) DepositSolution(amount int) {
	// BREAKPOINT: Increment balance
	// DEBUG: Watch 'b.balance += amount'
	// DEBUG: Pointer receiver required for mutation
	b.balance += amount
}

// BalanceSolution returns current balance.
// BREAKPOINT: Set breakpoint here to trace read
// DEBUG: Watch 'b.balance' value
func (b *BankAccount) BalanceSolution() int {
	// DEBUG: Pointer receiver for consistency
	return b.balance
}

// WithdrawSolution subtracts money from account.
// BREAKPOINT: Set breakpoint here to trace mutation
// DEBUG: Watch 'b.balance' decrease
func (b *BankAccount) WithdrawSolution(amount int) {
	// BREAKPOINT: Decrement balance
	// DEBUG: Watch 'b.balance -= amount'
	b.balance -= amount
}

// AreaSolution calculates rectangle area.
// BREAKPOINT: Set breakpoint here to trace value receiver
// DEBUG: Watch 'r' Rectangle copy
// DEBUG: Watch calculation with r.Width and r.Height
func (r Rectangle) AreaSolution() float64 {
	// BREAKPOINT: Calculate area
	// DEBUG: Value receiver - both Rectangle and *Rectangle satisfy Shape
	return r.Width * r.Height
}

// AreaSolution calculates circle area.
// BREAKPOINT: Set breakpoint here to trace pointer receiver
// DEBUG: Watch 'c' Circle pointer
// DEBUG: Only *Circle satisfies Shape (not Circle)
func (c *Circle) AreaSolution() float64 {
	// BREAKPOINT: Calculate area
	// DEBUG: Watch 'c.Radius' access
	return 3.14159 * c.Radius * c.Radius
}

// TotalAreaSolution sums all shape areas.
// BREAKPOINT: Set breakpoint here to trace iteration
// DEBUG: Watch 'shapes' slice
// DEBUG: Watch 'total' accumulate
func TotalAreaSolution(shapes []Shape) float64 {
	total := 0.0
	// BREAKPOINT: Loop through shapes
	for _, shape := range shapes {
		// DEBUG: Watch shape.Area() dynamic dispatch
		total += shape.Area()
	}
	return total
}

// AppendSolution adds value to list.
// BREAKPOINT: Set breakpoint here to trace nil handling
// DEBUG: Watch 'l' pointer (may be nil)
// DEBUG: Watch recursive traversal
func (l *StringList) AppendSolution(value string) *StringList {
	// BREAKPOINT: Check nil receiver
	// DEBUG: Watch 'l == nil' condition
	if l == nil {
		return &StringList{value: value, next: nil}
	}

	// BREAKPOINT: Check if tail
	if l.next == nil {
		l.next = &StringList{value: value, next: nil}
		return l
	}

	// BREAKPOINT: Recurse
	// DEBUG: Watch recursion to next node
	l.next = l.next.AppendSolution(value)
	return l
}

// ContainsSolution checks if value exists.
// BREAKPOINT: Set breakpoint here to trace search
// DEBUG: Watch 'l' traverse list
func (l *StringList) ContainsSolution(value string) bool {
	// BREAKPOINT: Check nil
	if l == nil {
		return false
	}

	// BREAKPOINT: Check current
	// DEBUG: Watch 'l.value == value'
	if l.value == value {
		return true
	}

	// BREAKPOINT: Recurse
	return l.next.ContainsSolution(value)
}

// FirstSolution returns first element.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch nil check
func (l *StringList) FirstSolution() string {
	if l == nil {
		return ""
	}
	return l.value
}

// ValidateSolution checks small config.
// BREAKPOINT: Set breakpoint here
// DEBUG: Value receiver for small struct (<64 bytes)
func (c SmallConfig) ValidateSolution() bool {
	return c.ID > 0 && c.Name != ""
}

// SumSolution sums large config data.
// BREAKPOINT: Set breakpoint here
// DEBUG: Pointer receiver for large struct (8000 bytes)
// DEBUG: Avoids copying entire array
func (l *LargeConfig) SumSolution() int {
	total := 0
	for _, v := range l.Data {
		total += v
	}
	return total
}

// SetNameSolution updates user name.
// BREAKPOINT: Set breakpoint here
// DEBUG: Pointer receiver for mutation
func (u *User) SetNameSolution(name string) {
	u.Name = name
}

// SetEmailSolution updates user email.
// BREAKPOINT: Set breakpoint here
// DEBUG: Pointer receiver for consistency
func (u *User) SetEmailSolution(email string) {
	u.Email = email
}

// GetNameSolution returns user name.
// BREAKPOINT: Set breakpoint here
// DEBUG: Pointer receiver for consistency
func (u *User) GetNameSolution() string {
	return u.Name
}

// IsAdultSolution checks if user is 18+.
// BREAKPOINT: Set breakpoint here
// DEBUG: Pointer receiver for consistency
func (u *User) IsAdultSolution() bool {
	return u.Age >= 18
}

// EqualsSolution compares points.
// BREAKPOINT: Set breakpoint here
// DEBUG: Value receiver allows both Point and *Point
func (p Point) EqualsSolution(other Comparable) bool {
	// BREAKPOINT: Type assert
	otherPoint, ok := other.(Point)
	if !ok {
		otherPointPtr, ok := other.(*Point)
		if !ok {
			return false
		}
		otherPoint = *otherPointPtr
	}
	return p.X == otherPoint.X && p.Y == otherPoint.Y
}

// NewSafeCounterMapSolution creates counter map.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch map initialization
func NewSafeCounterMapSolution() SafeCounterMap {
	return SafeCounterMap{
		counters: make(map[string]int),
	}
}

// IncrementSolution increments counter.
// BREAKPOINT: Set breakpoint here
// DEBUG: Pointer receiver to modify map
func (m *SafeCounterMap) IncrementSolution(key string) {
	m.counters[key]++
}

// GetSolution returns counter value.
// BREAKPOINT: Set breakpoint here
// DEBUG: Pointer receiver for consistency
func (m *SafeCounterMap) GetSolution(key string) int {
	return m.counters[key]
}

// AppendIterative is iterative version of append.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch iteration (no recursion)
func (l *StringList) AppendIterative(value string) *StringList {
	newNode := &StringList{value: value, next: nil}
	if l == nil {
		return newNode
	}
	current := l
	for current.next != nil {
		current = current.next
	}
	current.next = newNode
	return l
}
