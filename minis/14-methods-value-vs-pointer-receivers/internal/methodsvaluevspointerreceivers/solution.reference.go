//go:build reference

package methodsvaluevspointerreceivers

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

func (b *BankAccount) DepositSolution(amount int) {
	if b == nil {
		return
	}
	b.balance += amount
}

// BalanceSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (b *BankAccount) BalanceSolution() int {
	if b == nil {
		return 0
	}
	return b.balance
}

// WithdrawSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (b *BankAccount) WithdrawSolution(amount int) {
	if b == nil {
		return
	}
	b.balance -= amount
}

// AreaSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (r Rectangle) AreaSolution() float64 {
	return r.Width * r.Height
}

// AreaSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Circle) AreaSolution() float64 {
	if c == nil {
		return 0
	}
	return 3.14159 * c.Radius * c.Radius
}

// TotalAreaSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func TotalAreaSolution(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// AppendSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (l *StringList) AppendSolution(value string) *StringList {
	if l == nil {
		return &StringList{value: value}
	}
	if l.next == nil {
		l.next = &StringList{value: value}
		return l
	}
	l.next = l.next.AppendSolution(value)
	return l
}

// ContainsSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (l *StringList) ContainsSolution(value string) bool {
	if l == nil {
		return false
	}
	if l.value == value {
		return true
	}
	return l.next.ContainsSolution(value)
}

// FirstSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (l *StringList) FirstSolution() string {
	if l == nil {
		return ""
	}
	return l.value
}

// ValidateSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c SmallConfig) ValidateSolution() bool {
	return c.ID > 0 && c.Name != ""
}

// SumSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (l *LargeConfig) SumSolution() int {
	if l == nil {
		return 0
	}
	total := 0
	for _, v := range l.Data {
		total += v
	}
	return total
}

// SetNameSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (u *User) SetNameSolution(name string) {
	if u == nil {
		return
	}
	u.Name = name
}

// SetEmailSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (u *User) SetEmailSolution(email string) {
	if u == nil {
		return
	}
	u.Email = email
}

// GetNameSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (u *User) GetNameSolution() string {
	if u == nil {
		return ""
	}
	return u.Name
}

// IsAdultSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (u *User) IsAdultSolution() bool {
	if u == nil {
		return false
	}
	return u.Age >= 18
}

// EqualsSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (p Point) EqualsSolution(other Comparable) bool {
	switch v := other.(type) {
	case Point:
		return p.X == v.X && p.Y == v.Y
	case *Point:
		if v == nil {
			return false
		}
		return p.X == v.X && p.Y == v.Y
	default:
		return false
	}
}

// NewSafeCounterMapSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewSafeCounterMapSolution() SafeCounterMap {
	return SafeCounterMap{counters: make(map[string]int)}
}

// IncrementSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m *SafeCounterMap) IncrementSolution(key string) {
	if m == nil {
		return
	}
	if m.counters == nil {
		m.counters = make(map[string]int)
	}
	m.counters[key]++
}

// GetSolution implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m *SafeCounterMap) GetSolution(key string) int {
	if m == nil || m.counters == nil {
		return 0
	}
	return m.counters[key]
}

// AppendIterative implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (l *StringList) AppendIterative(value string) *StringList {
	if l == nil {
		return &StringList{value: value}
	}
	cur := l
	for cur.next != nil {
		cur = cur.next
	}
	cur.next = &StringList{value: value}
	return l
}
