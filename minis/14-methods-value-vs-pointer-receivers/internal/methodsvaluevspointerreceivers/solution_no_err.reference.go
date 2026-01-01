//go:build reference

package methodsvaluevspointerreceivers

/*
Core Method Receiver Patterns Without Complexity
=================================================

This file focuses on understanding when to use value vs pointer receivers.
This is one of the most important decisions in Go API design.

CORE CONCEPTS:
- Pointer receiver for mutation
- Value receiver for immutability
- Pointer receiver for large structs
- Consistency across method sets
- Interface satisfaction rules

DEBUGGING WORKFLOW:
1. Set breakpoint at method entry
2. Watch receiver type (value copy vs pointer)
3. Observe mutation behavior
4. Verify interface satisfaction

See /RUN_DEBUG.md for comprehensive debugging guide.
*/

// TODO: Implement BankAccount methods
// BREAKPOINT: Set breakpoints in each method
// DEBUG: Watch balance changes
//
// Use pointer receivers for ALL methods (consistency + mutation)
//
// func (b *BankAccount) Deposit(amount int) {
//     b.balance += amount
// }
//
// func (b *BankAccount) Balance() int {
//     return b.balance
// }
//
// func (b *BankAccount) Withdraw(amount int) {
//     b.balance -= amount
// }

// TODO: Implement Shape methods
// BREAKPOINT: Watch receiver type impact on interface
//
// Rectangle.Area: Use VALUE receiver (allows both Rectangle and *Rectangle to satisfy Shape)
// Circle.Area: Use POINTER receiver (only *Circle satisfies Shape)
//
// func (r Rectangle) Area() float64 {
//     return r.Width * r.Height
// }
//
// func (c *Circle) Area() float64 {
//     return 3.14159 * c.Radius * c.Radius
// }
//
// func TotalArea(shapes []Shape) float64 {
//     total := 0.0
//     for _, shape := range shapes {
//         total += shape.Area()
//     }
//     return total
// }

// TODO: Implement StringList methods
// BREAKPOINT: Watch nil receiver handling
// DEBUG: Watch list traversal
//
// Use pointer receivers for mutation and nil safety
//
// func (l *StringList) Append(value string) *StringList {
//     if l == nil {
//         return &StringList{value: value, next: nil}
//     }
//     if l.next == nil {
//         l.next = &StringList{value: value, next: nil}
//         return l
//     }
//     l.next = l.next.Append(value)
//     return l
// }
//
// func (l *StringList) Contains(value string) bool {
//     if l == nil {
//         return false
//     }
//     if l.value == value {
//         return true
//     }
//     return l.next.Contains(value)
// }
//
// func (l *StringList) First() string {
//     if l == nil {
//         return ""
//     }
//     return l.value
// }

// TODO: Implement Config methods
// BREAKPOINT: Watch performance difference
//
// SmallConfig: VALUE receiver (small struct, cheap to copy)
// LargeConfig: POINTER receiver (large struct, expensive to copy)
//
// func (c SmallConfig) Validate() bool {
//     return c.ID > 0 && c.Name != ""
// }
//
// func (l *LargeConfig) Sum() int {
//     total := 0
//     for _, v := range l.Data {
//         total += v
//     }
//     return total
// }

// TODO: Implement User methods
// BREAKPOINT: Watch consistency pattern
// DEBUG: All methods use pointer receiver
//
// Use pointer receivers for ALL methods (consistency)
//
// func (u *User) SetName(name string) {
//     u.Name = name
// }
//
// func (u *User) SetEmail(email string) {
//     u.Email = email
// }
//
// func (u *User) GetName() string {
//     return u.Name
// }
//
// func (u *User) IsAdult() bool {
//     return u.Age >= 18
// }

// TODO: Implement Point.Equals
// BREAKPOINT: Watch type assertion
// DEBUG: Value receiver allows both Point and *Point to satisfy Comparable
//
// func (p Point) Equals(other Comparable) bool {
//     otherPoint, ok := other.(Point)
//     if !ok {
//         otherPointPtr, ok := other.(*Point)
//         if !ok {
//             return false
//         }
//         otherPoint = *otherPointPtr
//     }
//     return p.X == otherPoint.X && p.Y == otherPoint.Y
// }

// TODO: Implement SafeCounterMap methods
// BREAKPOINT: Watch map operations
//
// func NewSafeCounterMap() SafeCounterMap {
//     return SafeCounterMap{
//         counters: make(map[string]int),
//     }
// }
//
// func (m *SafeCounterMap) Increment(key string) {
//     m.counters[key]++
// }
//
// func (m *SafeCounterMap) Get(key string) int {
//     return m.counters[key]
// }
