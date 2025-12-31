//go:build !solution
// +build !solution

package methodsvaluevspointerreceivers

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
func (b *BankAccount) Deposit(amount int) {
	// TODO: Implement the function.

	// Why a pointer receiver (`*BankAccount`)?
	// - This method MODIFIES the `balance` field of the `BankAccount` struct.
	// - If you use a value receiver (`b BankAccount`), the method operates on a COPY of the `BankAccount`.
	// - The change to `balance` would be lost when the method returns.
	// - By using a pointer receiver, you are working with a pointer to the original `BankAccount` instance, so modifications are persistent.
	//
	// Memory & Language Comparison:
	// - Go: The choice between value and pointer receivers is a fundamental part of Go's design, giving developers control over memory and mutation.
	// - C++: This is similar to the difference between passing an object by value and passing a pointer or reference to an object.
	// - Java: All objects are passed by reference (technically, the reference itself is passed by value), so you don't explicitly make this choice. Modifications to an object's fields within a method will always persist.
	// - Rust: The concept of "borrowing" (`&` for a shared reference, `&mut` for a mutable reference) provides similar control but with compile-time safety guarantees. A method that modifies the struct would take `&mut self`.
}

// Balance returns the current balance.
//
// REQUIREMENTS:
// - Return the current balance
// - Use the appropriate receiver type (read-only operation)
func (b *BankAccount) Balance() int {
	// TODO: Implement the function.

	// Why a pointer receiver here?
	// - This method is "read-only"; it doesn't modify the `BankAccount`.
	// - Technically, a value receiver (`b BankAccount`) would also work and might even be slightly more efficient for very small structs because it avoids a pointer dereference.
	// - However, the Go community strongly recommends consistency. If any method on a type needs a pointer receiver (like `Deposit` and `Withdraw`), all methods on that type should have a pointer receiver.
	// - This makes the type's method set predictable and avoids subtle issues with interface satisfaction.
	return 0
}

// Withdraw subtracts money from the account.
//
// REQUIREMENTS:
// - Use the appropriate receiver type
// - Subtract the amount from the balance
// - For this exercise, don't worry about negative balances
func (b *BankAccount) Withdraw(amount int) {
	// TODO: Implement the function.

	// Why a pointer receiver (`*BankAccount`)?
	// - Same reason as `Deposit`. This method MODIFIES the `balance`.
	// - A pointer receiver is required for the change to be saved to the original `BankAccount` instance.
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
func (r Rectangle) Area() float64 {
	// TODO: Implement the function.

	// Why a value receiver (`r Rectangle`)?
	// - The Go spec says: "A type T satisfies an interface if it implements all the methods of the interface."
	// - If a method has a value receiver, both the value type `T` and the pointer type `*T` can call it.
	//   - For `T`, the method is called directly.
	//   - For `*T`, Go automatically dereferences the pointer to get the value, then calls the method.
	// - Therefore, to have both `Rectangle` and `*Rectangle` satisfy the `Shape` interface, the `Area` method must have a value receiver.
	//
	// Memory & Language Comparison:
	// - This is a unique aspect of Go's type system. It provides flexibility.
	// - In C++, you would have to explicitly define how pointers and values relate to a base class or interface.
	// - In Java/C#, since everything is a reference, this distinction doesn't really exist in the same way.

	return 0.0
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
func (c *Circle) Area() float64 {
	// TODO: Implement the function.

	// Why a pointer receiver (`*Circle`)?
	// - The Go spec also says that if a method has a pointer receiver, only the pointer type `*T` can call it directly.
	// - A value type `T` *cannot* call a method with a pointer receiver because Go doesn't know how to automatically create a pointer to an arbitrary value (e.g., what if it's a temporary value?).
	// - Therefore, by giving `Area` a pointer receiver, you are saying that only a `*Circle` can be considered a `Shape`. A `Circle` value will not satisfy the interface.
	return 0.0
}

// TotalArea calculates the total area of multiple shapes.
//
// REQUIREMENTS:
// - Sum the areas of all shapes in the slice
// - This tests your understanding of which types satisfy the Shape interface
func TotalArea(shapes []Shape) float64 {
	// TODO: Implement the function.

	// This function will work correctly if you've implemented the `Area` methods with the correct receiver types.
	// You can pass a `Rectangle`, a `*Rectangle`, and a `*Circle` into a `[]Shape` slice.
	// You *cannot* pass a `Circle` value.
	//
	// The loop will call the `Area()` method on each element. This is an example of polymorphism in Go.
	// The Go runtime knows the concrete type of each element in the `shapes` slice and calls the correct `Area` method.
	return 0.0
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
func (l *StringList) Append(value string) *StringList {
	// TODO: Implement the function.

	// This exercise demonstrates the power of "nil receivers".
	// In Go, you can call methods on a nil pointer without causing a panic, as long as the method itself handles the nil case.

	// Step 1: Handle the nil receiver case.
	// - `if l == nil { ... }`
	// - This is the key to making the API "fluent". If the list doesn't exist yet, this method will create it.
	// - Return a new `*StringList` containing the value.

	// Step 2: Handle the recursive case.
	// - If `l.next` is nil, you've reached the end of the list. Create a new node and assign it to `l.next`.
	// - Otherwise, recursively call `Append` on the rest of the list: `l.next = l.next.Append(value)`.
	// - Return the original receiver `l`.

	// Memory & Language Comparison:
	// - Go: The ability to call methods on nil pointers is a unique feature that enables this kind of safe, fluent API design.
	// - C/C++: Calling a method on a null pointer is undefined behavior and will almost certainly crash your program (segmentation fault). You must always check for null *before* calling a function or method.
	// - Java/C#: Calling a method on a `null` reference will throw a `NullPointerException`. You are expected to check for `null` before method calls.
	// - Rust: The concept of `null` doesn't exist in safe Rust. Instead, you use the `Option` type (`Some(value)` or `None`), which forces you to handle the "no value" case at compile time, making null-related runtime errors impossible.
	return nil
}

// Contains checks if a value exists in the list.
//
// REQUIREMENTS:
// - Return true if the value is found, false otherwise
// - Handle nil receivers safely (nil list doesn't contain anything)
func (l *StringList) Contains(value string) bool {
	// TODO: Implement the function.

	// Step 1: Handle the nil receiver case.
	// - A nil list contains nothing, so `if l == nil { return false }`.

	// Step 2: Check the current node.
	// - `if l.value == value { return true }`.

	// Step 3: Recursively check the rest of the list.
	// - `return l.next.Contains(value)`.
	// - This works because the recursive call will handle the case where `l.next` is nil.
	return false
}

// First returns the first element in the list.
//
// REQUIREMENTS:
// - If the list is nil or empty, return ""
// - Otherwise, return the first value
func (l *StringList) First() string {
	// TODO: Implement the function.

	// Step 1: Handle the nil receiver case.
	// - A nil list has no first element, so `if l == nil { return "" }`.
	// - This prevents a panic from trying to access `l.value`.

	// Step 2: Return the value of the current node.
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
	// TODO: Implement the function.

	// Why a value receiver (`c SmallConfig`)?
	// - This struct is very small (an `int` and a `string` header total about 24 bytes on a 64-bit system).
	// - Copying 24 bytes is extremely fast and happens on the stack, which is cheaper than creating a pointer and dereferencing it (which involves heap memory).
	// - The general rule in Go: for small, read-only structs, prefer value receivers. It can be more performant and makes the "read-only" nature of the operation clear.
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
func (l *LargeConfig) Sum() int {
	// TODO: Implement the function.

	// Why a pointer receiver (`l *LargeConfig`)?
	// - This struct is very large (an array of 1000 integers is 8000 bytes on a 64-bit system).
	// - If you used a value receiver, all 8000 bytes would be copied every time this method is called. This is a significant performance penalty.
	// - By using a pointer receiver, only the pointer (8 bytes) is copied. The method then operates on the original data.
	// - The general rule in Go: for large structs, always use a pointer receiver, even for read-only methods, to avoid expensive copies.
	//
	// Memory & Language Comparison:
	// - Go: This is a manual performance optimization that Go developers are expected to know.
	// - C++: The same principle applies. You would pass a large struct by `const` reference (`const LargeConfig&`) to avoid a copy.
	// - Java/C#: Objects are always passed by reference, so you don't need to worry about this. The runtime handles it for you.
	// - Rust: The compiler is often smart enough to optimize away copies in some cases, but the standard practice is still to pass large structs by reference (`&self`) for read-only operations.
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
func (u *User) SetName(name string) {
	// TODO: Implement the function.
	// This method must use a pointer receiver because it modifies the `User` struct.
}

// SetEmail updates the user's email.
//
// REQUIREMENTS:
// - Set the Email field to the provided value
// - Use the SAME receiver type as SetName (consistency!)
func (u *User) SetEmail(email string) {
	// TODO: Implement the function.
	// This method also mutates the struct, so it also needs a pointer receiver.
}

// GetName returns the user's name.
//
// REQUIREMENTS:
// - Return the Name field
// - Use the SAME receiver type as other methods (consistency!)
func (u *User) GetName() string {
	// TODO: Implement the function.

	// Why a pointer receiver here, even though it's read-only?
	// - **Consistency.** When you have a type where some methods need to be pointers (because they modify the struct), it's a strong convention in Go to make *all* methods on that type have pointer receivers.
	// - This avoids a mixed method set. A mixed set can be confusing and can lead to situations where `T` and `*T` satisfy different interfaces, which is usually not what you want.
	// - For a type like `User`, which represents an "entity" that can be changed, you should always interact with it via a pointer.
	return ""
}

// IsAdult checks if the user is 18 or older.
//
// REQUIREMENTS:
// - Return true if Age >= 18
// - Use the SAME receiver type (consistency!)
func (u *User) IsAdult() bool {
	// TODO: Implement the function.

	// Again, use a pointer receiver for consistency with `SetName` and `SetEmail`.
	// This makes the API for `*User` predictable. All methods are available on the pointer.
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
func (p Point) Equals(other Comparable) bool {
	// TODO: Implement the function.

	// This method needs to compare `p` with `other`.
	// But `other` is an interface type (`Comparable`). To get the X and Y values, you must get the concrete `Point` type out of the interface.

	// Step 1: Type Assert the interface.
	// - `otherPoint, ok := other.(Point)`
	// - This is a "type assertion". It checks if the concrete type stored inside the `other` interface is a `Point`.
	// - If it is, `ok` will be `true` and `otherPoint` will be the `Point` value.
	// - If it's not, `ok` will be `false`.

	// Step 2: Handle the case where the assertion fails.
	// - It's possible the `other` value is a `*Point`. Because this method has a value receiver, a `*Point` also satisfies the `Comparable` interface.
	// - So, if the first assertion fails, you should try another one: `otherPointPtr, ok := other.(*Point)`.
	// - If this second assertion succeeds, you have a pointer. You need the value to compare, so you must dereference it: `otherPoint = *otherPointPtr`.
	// - If both assertions fail, then `other` is not a `Point` or a `*Point`, so they can't be equal. Return `false`.

	// Step 3: Compare the fields.
	// - `return p.X == otherPoint.X && p.Y == otherPoint.Y`

	// Memory & Language Comparison:
	// - Go: Type assertions are Go's way of getting a concrete type from an interface. The `value, ok` form is the safe way to do it, preventing a panic if the type is wrong.
	// - C++: This is similar to using `dynamic_cast`, which can be used to safely downcast a base class pointer to a derived class pointer.
	// - Java: This is like using `instanceof` to check the type, and then casting the object.
	// - Rust: The `match` statement on an `enum` or a `Box<dyn Trait>` provides a more structured and compile-time safe way to handle different concrete types.
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
	// TODO: Implement the function.

	// Why initialize the map?
	// - The zero value for a map is `nil`.
	// - If you try to write to a `nil` map, your program will panic.
	// - By initializing the map with `make(map[string]int)`, you ensure that it's ready to be used.
	return SafeCounterMap{}
}

// Increment increments the counter for a given key.
//
// REQUIREMENTS:
// - If the key doesn't exist, initialize it to 0 then increment to 1
// - If the key exists, increment its value
// - Use a POINTER receiver (we're modifying the map)
func (m *SafeCounterMap) Increment(key string) {
	// TODO: Implement the function.

	// Why a pointer receiver?
	// - Although a map is a reference type (it's a pointer to a hash table header), the `SafeCounterMap` struct itself is a value type.
	// - If this method had a value receiver (`m SafeCounterMap`), `m` would be a copy.
	// - While you could modify the map through that copy (since the map itself is a reference), it's not idiomatic.
	// - The `Increment` method is conceptually a mutation of the `SafeCounterMap`, so you should use a pointer receiver to make that clear.

	// How does `m.counters[key]++` work?
	// - If `key` doesn't exist in the map, Go's map implementation will return the zero value for the value type (which is `0` for `int`).
	// - The expression `m.counters[key]++` will then behave like `0++`, resulting in `1` being stored for that key.
	// - This is a convenient feature of Go maps.
}

// Get returns the counter value for a given key.
//
// REQUIREMENTS:
// - Return the value for the key
// - If the key doesn't exist, return 0
func (m *SafeCounterMap) Get(key string) int {
	// TODO: Implement the function.
	// - Use a pointer receiver for consistency with `Increment`.
	// - Reading from a map is safe. If the key doesn't exist, the zero value (`0` for `int`) is returned.
	return 0
}
