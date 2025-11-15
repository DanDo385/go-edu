//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for understanding method receivers.
//
// LEARNING OBJECTIVES:
// - Implement methods with appropriate receiver types
// - Understand when modifications persist vs when they don't
// - Design types that satisfy interfaces correctly
// - Handle nil receivers safely
// - Make informed decisions about receiver types

package exercise

// TODO: Implement these functions and methods according to the specifications in the tests.

// ============================================================================
// EXERCISE 1: Basic Receiver Types
// ============================================================================

// BankAccount represents a simple bank account.
type BankAccount struct {
	balance int  // Balance in cents
}

// Deposit adds money to the account.
//
// REQUIREMENTS:
// - Use the appropriate receiver type so that the deposit PERSISTS
// - Add the amount to the balance
// - Hint: Does this method modify the receiver? Then use pointer receiver!
func (b BankAccount) Deposit(amount int) {
	// TODO: Implement this
	// Current receiver type is VALUE - is this correct?
}

// Balance returns the current balance.
//
// REQUIREMENTS:
// - Return the current balance
// - Use the appropriate receiver type (read-only operation)
// - Hint: This doesn't modify the receiver, but be consistent!
func (b BankAccount) Balance() int {
	// TODO: Implement this
	return 0
}

// Withdraw subtracts money from the account.
//
// REQUIREMENTS:
// - Use the appropriate receiver type
// - Subtract the amount from the balance
// - For this exercise, don't worry about negative balances
func (b BankAccount) Withdraw(amount int) {
	// TODO: Implement this
}

// ============================================================================
// EXERCISE 2: Interface Satisfaction
// ============================================================================

// Shape is an interface for geometric shapes.
type Shape interface {
	Area() float64
}

// Rectangle represents a rectangle.
type Rectangle struct {
	Width, Height float64
}

// Area calculates the area of the rectangle.
//
// REQUIREMENTS:
// - Implement this method so that BOTH Rectangle and *Rectangle satisfy Shape
// - Hint: What receiver type allows both the value and pointer to satisfy an interface?
func (r Rectangle) Area() float64 {
	// TODO: Implement this
	return 0
}

// Circle represents a circle.
type Circle struct {
	Radius float64
}

// Area calculates the area of the circle.
//
// REQUIREMENTS:
// - Implement this method so that ONLY *Circle satisfies Shape
// - Use π ≈ 3.14159
// - Hint: Use a pointer receiver to restrict interface satisfaction
func (c Circle) Area() float64 {
	// TODO: Implement this (but fix the receiver type first!)
	return 0
}

// TotalArea calculates the total area of multiple shapes.
//
// REQUIREMENTS:
// - Sum the areas of all shapes in the slice
// - This tests your understanding of which types satisfy the Shape interface
func TotalArea(shapes []Shape) float64 {
	// TODO: Implement this
	return 0
}

// ============================================================================
// EXERCISE 3: Nil Receiver Safety
// ============================================================================

// StringList is a linked list of strings.
type StringList struct {
	value string
	next  *StringList
}

// Append adds a value to the end of the list.
//
// REQUIREMENTS:
// - If the list is nil, return a new list with just the value
// - Otherwise, traverse to the end and append
// - Return the head of the list
//
// EXAMPLE:
//   var list *StringList
//   list = list.Append("a")  // Returns: a -> nil
//   list = list.Append("b")  // Returns: a -> b -> nil
//
// HINT: Check if the receiver is nil first!
func (l *StringList) Append(value string) *StringList {
	// TODO: Implement this
	return nil
}

// Contains checks if a value exists in the list.
//
// REQUIREMENTS:
// - Return true if the value is found, false otherwise
// - Handle nil receivers safely (nil list doesn't contain anything)
// - Traverse the list to search
func (l *StringList) Contains(value string) bool {
	// TODO: Implement this
	return false
}

// First returns the first element in the list.
//
// REQUIREMENTS:
// - If the list is nil or empty, return ""
// - Otherwise, return the first value
// - This tests nil safety
func (l *StringList) First() string {
	// TODO: Implement this
	return ""
}

// ============================================================================
// EXERCISE 4: Performance-Aware Design
// ============================================================================

// SmallConfig is a small configuration (use value receiver).
type SmallConfig struct {
	ID   int
	Name string  // Strings are cheap to copy (just 16 bytes)
}

// Validate checks if the configuration is valid.
//
// REQUIREMENTS:
// - Return true if ID > 0 and Name is not empty
// - Use a VALUE receiver (the struct is small, ~24 bytes)
func (c SmallConfig) Validate() bool {
	// TODO: Implement this
	return false
}

// LargeConfig is a large configuration (use pointer receiver).
type LargeConfig struct {
	Data [1000]int  // 8000 bytes!
}

// Sum calculates the sum of all data elements.
//
// REQUIREMENTS:
// - Sum all elements in the Data array
// - Use a POINTER receiver to avoid copying 8000 bytes
func (l LargeConfig) Sum() int {
	// TODO: Implement this (but fix the receiver type first!)
	return 0
}

// ============================================================================
// EXERCISE 5: Consistency and Best Practices
// ============================================================================

// User represents a user with mutable fields.
type User struct {
	Name  string
	Email string
	Age   int
}

// SetName updates the user's name.
//
// REQUIREMENTS:
// - Set the Name field to the provided value
// - Use a POINTER receiver (this mutates the user)
func (u User) SetName(name string) {
	// TODO: Implement this (but fix the receiver type first!)
}

// SetEmail updates the user's email.
//
// REQUIREMENTS:
// - Set the Email field to the provided value
// - Use the SAME receiver type as SetName (consistency!)
func (u User) SetEmail(email string) {
	// TODO: Implement this (but fix the receiver type first!)
}

// GetName returns the user's name.
//
// REQUIREMENTS:
// - Return the Name field
// - Use the SAME receiver type as other methods (consistency!)
// - Even though this is read-only, consistency matters!
func (u User) GetName() string {
	// TODO: Implement this (but fix the receiver type first!)
	return ""
}

// IsAdult checks if the user is 18 or older.
//
// REQUIREMENTS:
// - Return true if Age >= 18
// - Use the SAME receiver type (consistency!)
func (u User) IsAdult() bool {
	// TODO: Implement this (but fix the receiver type first!)
	return false
}

// ============================================================================
// EXERCISE 6: Method Set Understanding
// ============================================================================

// Comparable is an interface for comparable types.
type Comparable interface {
	Equals(other Comparable) bool
}

// Point represents a 2D point.
type Point struct {
	X, Y int
}

// Equals checks if two points are equal.
//
// REQUIREMENTS:
// - Return true if both X and Y match the other point
// - Use a VALUE receiver so that both Point and *Point satisfy Comparable
// - Hint: Type assert 'other' to Point to access X and Y
func (p Point) Equals(other Comparable) bool {
	// TODO: Implement this
	// You'll need to type assert: otherPoint, ok := other.(Point)
	return false
}

// ============================================================================
// EXERCISE 7: Map Element Challenge
// ============================================================================

// SafeCounterMap wraps a map of counters and provides safe increment operations.
type SafeCounterMap struct {
	counters map[string]int
}

// NewSafeCounterMap creates a new SafeCounterMap with an initialized map.
//
// REQUIREMENTS:
// - Return a SafeCounterMap with a non-nil map
// - The map should be empty initially
func NewSafeCounterMap() SafeCounterMap {
	// TODO: Implement this
	return SafeCounterMap{}
}

// Increment increments the counter for a given key.
//
// REQUIREMENTS:
// - If the key doesn't exist, initialize it to 0 then increment to 1
// - If the key exists, increment its value
// - Use a POINTER receiver (we're modifying the map)
//
// This demonstrates the right way to handle map modifications
// (you can't call pointer-receiver methods on map elements directly).
func (m SafeCounterMap) Increment(key string) {
	// TODO: Implement this (but fix the receiver type first!)
}

// Get returns the counter value for a given key.
//
// REQUIREMENTS:
// - Return the value for the key
// - If the key doesn't exist, return 0
func (m SafeCounterMap) Get(key string) int {
	// TODO: Implement this (and use consistent receiver type!)
	return 0
}
