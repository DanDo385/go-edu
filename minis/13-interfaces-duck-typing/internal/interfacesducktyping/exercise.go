//go:build !solution
// +build !solution

package interfacesducktyping

// TODO: Import required packages
// You'll need:
// - "fmt" for formatting strings and type names
//
// import "fmt"

// ============================================================================
// INTERFACES AND DUCK TYPING: Understanding Go's Interface System
// ============================================================================
//
// Go's interface philosophy: "If it walks like a duck and quacks like a duck..."
//
// Key concepts:
// 1. Interfaces are IMPLICIT (no "implements" keyword)
// 2. If a type has all the methods, it satisfies the interface automatically
// 3. Interfaces are implemented by VALUE or POINTER depending on receiver type
// 4. Empty interface (interface{}) accepts ANY type
// 5. nil interfaces vs interfaces containing nil pointers (gotcha!)
//
// Interface internals (runtime representation):
// An interface value contains TWO pointers (16 bytes total on 64-bit):
//   - Type pointer: Points to type metadata (what concrete type is this?)
//   - Data pointer: Points to the actual value
//
// Example:
//   var i interface{} = 42
//   // Type pointer: points to int type metadata
//   // Data pointer: points to memory containing 42
//
// Interface satisfaction rules:
// - Methods with VALUE receiver (func (t T) Method()):
//   * Both T and *T satisfy the interface
//   * Reason: Go can automatically create pointers from values
//
// - Methods with POINTER receiver (func (t *T) Method()):
//   * Only *T satisfies the interface
//   * Reason: Can't automatically create value from pointer (might be nil)
//
// The nil gotcha:
//   var p *Person = nil
//   var i interface{} = p
//   i == nil  // FALSE! (Type is *Person, Value is nil)
//   p == nil  // TRUE
//
// Type assertions:
// - Single return: value := i.(Type)  // Panics if wrong type
// - Comma-ok: value, ok := i.(Type)  // Returns false if wrong type
// - Type switch: switch v := i.(type) { case int: ... }
//
// Empty interface (interface{} or any in Go 1.18+):
// - Accepts any value (all types satisfy it)
// - Use for generic containers, unmarshaling, reflection
// - Downside: Loses type safety, requires type assertions
//
// ============================================================================

// ============================================================================
// Exercise 1: Implementing Interfaces (Implicit Satisfaction)
// ============================================================================

// TODO: Implement the String() method for Person to satisfy the Stringer interface.
//
// The Stringer interface (defined in types.go):
//   type Stringer interface {
//       String() string
//   }
//
// Any type with a String() string method automatically satisfies Stringer.
// This is Go's "duck typing" - no explicit declaration needed!
//
// Requirements:
// - Return a string in the format: "Name (Age years old)"
// - Example: Person{Name: "Alice", Age: 30} → "Alice (30 years old)"
//
// Steps to implement:
//
// 1. Use fmt.Sprintf to format the string
//    - Use: fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
//    - %s formats the string (Name)
//    - %d formats the integer (Age)
//    - Spaces and parentheses are literal text
//
// 2. Return the formatted string
//    - The return satisfies the String() string signature
//
// Key Go concepts:
// - No "implements Stringer" declaration needed
// - Just implement the method with matching signature
// - Receiver type matters: (p Person) vs (p *Person)
//
// Method receiver choice:
// - Use VALUE receiver (p Person) because:
//   * Person is small (16 bytes: string header + int)
//   * Method doesn't modify Person
//   * Allows both Person and *Person to satisfy Stringer
//
// Memory behavior:
// - Value receiver: Makes a copy of Person (16 bytes)
// - Pointer receiver: Just copies pointer (8 bytes)
// - For small structs, value receiver is fine
//
// Example usage:
//   p := Person{Name: "Bob", Age: 25}
//   fmt.Println(p.String())  // "Bob (25 years old)"
//   fmt.Println(p)           // Also calls String() automatically!
//
// TODO: Implement the String method below
// func (p Person) String() string {
//     return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
// }

// ============================================================================
// Exercise 2: Type Assertions (Extracting Concrete Types)
// ============================================================================

// GetAge extracts the age from a Stringer if it's a Person.
// TODO: Implement type assertion with comma-ok idiom.
//
// Type assertions convert interface → concrete type:
//   value, ok := interfaceVar.(ConcreteType)
//   - If interfaceVar contains ConcreteType: value = the value, ok = true
//   - If interfaceVar contains different type: value = zero value, ok = false
//
// Requirements:
// - If s contains a Person, return the person's age and true
// - If s contains anything else, return 0 and false
//
// Steps to implement:
//
// 1. Perform type assertion to Person
//    - Use: p, ok := s.(Person)
//    - This checks if the concrete type inside s is Person
//    - p gets the Person value (or zero value if wrong type)
//    - ok indicates success or failure
//
// 2. Check if assertion succeeded
//    - Use: if !ok { return 0, false }
//    - If wrong type, return default values
//    - This prevents using an invalid p value
//
// 3. Return the age if successful
//    - Use: return p.Age, true
//    - p is guaranteed to be a valid Person
//    - true indicates we found a Person
//
// Key Go concepts:
// - Type assertion WITHOUT comma-ok panics on wrong type
// - Type assertion WITH comma-ok returns false on wrong type (safer)
// - Always use comma-ok for untrusted interface values
//
// Why this is safe:
//   p, ok := s.(Person)  // Never panics
//   vs
//   p := s.(Person)      // Panics if s isn't Person!
//
// Common pattern:
//   if p, ok := s.(Person); ok {
//       // Use p here (it's a Person)
//       return p.Age, true
//   }
//   return 0, false
//
// TODO: Implement the GetAge function below
// func GetAge(s Stringer) (int, bool) {
//     p, ok := s.(Person)
//     if !ok {
//         return 0, false
//     }
//     return p.Age, true
// }

// ============================================================================
// Exercise 3: Type Switches (Multiple Type Checks)
// ============================================================================

// DescribeType returns a description of the type of the value.
// TODO: Implement type switch for runtime type inspection.
//
// Type switch syntax:
//   switch v := i.(type) {
//   case int:
//       // v has type int here
//   case string:
//       // v has type string here
//   default:
//       // v has type interface{}
//   }
//
// Requirements:
// - For int: return "Integer: <value>"
// - For string: return "String: <value>"
// - For bool: return "Boolean: <value>"
// - For Person: return "Person: <name>"
// - For nil: return "Nil"
// - For any other type: return "Unknown"
//
// Steps to implement:
//
// 1. Use type switch to check the type
//    - Use: switch v := i.(type) { ... }
//    - v takes on the type of each case
//    - Special syntax: .(type) only works in switch
//
// 2. Handle each type case
//    - case int: v is an int, format with %d
//    - case string: v is a string, format with %s
//    - case bool: v is a bool, format with %t
//    - case Person: v is a Person, access v.Name
//
// 3. Handle nil case
//    - case nil: Checks for nil interface
//    - Return "Nil"
//
// 4. Handle unknown types
//    - default: Catches everything else
//    - Return "Unknown"
//
// Key Go concepts:
// - Type switch is the idiomatic way to handle multiple types
// - Each case has v with the appropriate type (type-safe)
// - nil is a valid case (checks if interface is nil)
// - Order doesn't matter (unlike value switch where it might)
//
// Memory behavior:
// - Type switch doesn't copy the value
// - Just inspects the type pointer in the interface
// - Fast operation (pointer comparison)
//
// Example usage:
//   DescribeType(42)                    // "Integer: 42"
//   DescribeType("hello")               // "String: hello"
//   DescribeType(true)                  // "Boolean: true"
//   DescribeType(Person{Name: "Alice"}) // "Person: Alice"
//   DescribeType(nil)                   // "Nil"
//   DescribeType(3.14)                  // "Unknown"
//
// TODO: Implement the DescribeType function below
// func DescribeType(i interface{}) string {
//     switch v := i.(type) {
//     case int:
//         return fmt.Sprintf("Integer: %d", v)
//     case string:
//         return fmt.Sprintf("String: %s", v)
//     case bool:
//         return fmt.Sprintf("Boolean: %t", v)
//     case Person:
//         return fmt.Sprintf("Person: %s", v.Name)
//     case nil:
//         return "Nil"
//     default:
//         return "Unknown"
//     }
// }

// ============================================================================
// Exercise 4: The Two-Part Nil Problem
// ============================================================================

// IsValidEmail checks if a Validator is valid, handling nil correctly.
// TODO: Implement nil-safe interface checking.
//
// The two-part nil problem:
// 1. nil interface: Both type and value are nil
//    var i interface{} = nil  // i == nil is TRUE
//
// 2. Interface containing nil pointer: Type is set, value is nil
//    var p *Person = nil
//    var i interface{} = p    // i == nil is FALSE!
//    // i's type is *Person, value is nil
//
// This is CRITICAL to understand!
//
// Requirements:
// - If v is nil (true nil interface), return false
// - If v contains a nil *Email pointer, return false
// - Otherwise, call v.IsValid() and return its result
//
// Steps to implement:
//
// 1. Check if interface itself is nil
//    - Use: if v == nil { return false }
//    - This checks both type and value are nil
//    - Safe to do before any other checks
//
// 2. Type assert to *Email
//    - Use: e, ok := v.(*Email)
//    - Check if concrete type is *Email
//    - If not *Email, it's some other Validator
//
// 3. Handle non-Email Validators
//    - Use: if !ok { return v.IsValid() }
//    - It's a different type that implements Validator
//    - Trust its IsValid() implementation
//
// 4. Check if Email pointer is nil
//    - Use: if e == nil { return false }
//    - At this point, v != nil but e might be nil
//    - This is the "interface containing nil pointer" case
//
// 5. Call IsValid on non-nil Email
//    - Use: return e.IsValid()
//    - Now we know e is a valid *Email
//
// Key Go concepts:
// - An interface is nil ONLY if both type and value are nil
// - An interface containing a nil pointer is NOT nil
// - Always check both the interface and the concrete pointer
//
// The gotcha in code:
//   var e *Email = nil
//   var v Validator = e
//   v == nil  // FALSE! (type is *Email, value is nil)
//   e == nil  // TRUE
//
// Why this happens:
// - Setting interface = pointer stores the TYPE (*Email)
// - Even if the pointer value is nil
// - Interface is only nil if you assign nil directly (or zero value)
//
// Common pattern to avoid this:
//   var e *Email
//   if condition {
//       e = &Email{...}
//   }
//   if e != nil {
//       var v Validator = e  // Only assign if not nil
//   }
//
// TODO: Implement the IsValidEmail function below
// func IsValidEmail(v Validator) bool {
//     if v == nil {
//         return false
//     }
//     e, ok := v.(*Email)
//     if !ok {
//         return v.IsValid()
//     }
//     if e == nil {
//         return false
//     }
//     return e.IsValid()
// }

// ============================================================================
// Exercise 5: Implementing Multiple Interfaces
// ============================================================================

// TODO: Implement Read() for Buffer to satisfy the Reader interface.
//
// The Reader interface (defined in types.go):
//   type Reader interface {
//       Read() string
//   }
//
// Requirements:
// - Return the current data in the buffer
//
// Steps to implement:
//
// 1. Return b.data
//    - Use: return b.data
//    - Simple accessor method
//    - No modification needed
//
// Why pointer receiver:
// - Consistency with Write (which needs pointer receiver)
// - Even though Read doesn't modify, using pointer keeps API consistent
// - Rule: If one method uses pointer, all should
//
// TODO: Implement the Read method below
// func (b *Buffer) Read() string {
//     return b.data
// }

// TODO: Implement Write() for Buffer to satisfy the Writer interface.
//
// The Writer interface (defined in types.go):
//   type Writer interface {
//       Write(data string) error
//   }
//
// Requirements:
// - Append the data to the buffer's existing data
// - Return nil (no errors for in-memory buffer)
//
// Steps to implement:
//
// 1. Append data to existing buffer
//    - Use: b.data += data
//    - String concatenation (creates new string)
//    - Modifies the buffer
//
// 2. Return nil error
//    - Use: return nil
//    - In-memory buffers don't fail
//    - Real I/O might return errors (disk full, etc.)
//
// Why pointer receiver:
// - MUST use pointer to modify b.data
// - Value receiver would modify a copy (changes lost)
//
// Key Go concepts:
// - Buffer now implements BOTH Reader and Writer
// - Therefore also implements ReadWriter (embedded interfaces)
// - All automatic - no declaration needed!
//
// TODO: Implement the Write method below
// func (b *Buffer) Write(data string) error {
//     b.data += data
//     return nil
// }

// ============================================================================
// Exercise 6: Interface Composition (Embedded Interfaces)
// ============================================================================

// IsReadWriter checks if an interface value implements ReadWriter.
// TODO: Implement interface composition checking.
//
// The ReadWriter interface (defined in types.go):
//   type ReadWriter interface {
//       Reader  // Embeds Reader interface
//       Writer  // Embeds Writer interface
//   }
//
// This means ReadWriter requires:
// - Read() string  (from Reader)
// - Write(string) error  (from Writer)
//
// Requirements:
// - Return true if i implements both Read() and Write()
// - Return false otherwise
//
// Steps to implement:
//
// 1. Type assert to ReadWriter
//    - Use: _, ok := i.(ReadWriter)
//    - Checks if i has both Read and Write methods
//    - Underscore _ discards the value (we only need ok)
//
// 2. Return the result
//    - Use: return ok
//    - true if i implements ReadWriter
//    - false otherwise
//
// Key Go concepts:
// - Interface embedding creates composite interfaces
// - ReadWriter = Reader + Writer (all methods required)
// - Type assertion checks if ALL methods are present
//
// Example:
//   buf := &Buffer{}
//   IsReadWriter(buf)  // true (has Read and Write)
//
//   type ReadOnly struct{}
//   func (r ReadOnly) Read() string { return "" }
//   IsReadWriter(ReadOnly{})  // false (missing Write)
//
// TODO: Implement the IsReadWriter function below
// func IsReadWriter(i interface{}) bool {
//     _, ok := i.(ReadWriter)
//     return ok
// }

// ============================================================================
// Exercise 7: Method Sets and Pointer Receivers
// ============================================================================

// TODO: Implement Increment() with a POINTER receiver.
//
// The Incrementer interface (defined in types.go):
//   type Incrementer interface {
//       Increment()
//   }
//
// Requirements:
// - Increment the Value field by 1
// - Use a pointer receiver so the change persists
//
// Steps to implement:
//
// 1. Increment the counter
//    - Use: c.Value++
//    - Increments the Value field
//    - Changes persist because c is a pointer
//
// Why pointer receiver:
// - MUST use pointer to modify c.Value
// - Value receiver would modify a copy
// - Only *Counter implements Incrementer, not Counter
//
// Method set rules:
// - Type T: Methods with receiver T
// - Type *T: Methods with receiver T or *T (both!)
// - Interface requirement: Must have exact receiver type
//
// This means:
//   var c Counter
//   c.Increment()  // OK (Go automatically converts to &c)
//
//   var i Incrementer = c  // ERROR! Counter doesn't implement Incrementer
//   var i Incrementer = &c // OK! *Counter implements Incrementer
//
// TODO: Implement the Increment method below
// func (c *Counter) Increment() {
//     c.Value++
// }

// CanIncrement checks if a value can be used as an Incrementer.
// TODO: Implement interface type checking.
//
// Requirements:
// - Return true if i implements the Incrementer interface
// - Return false otherwise
//
// Steps to implement:
//
// 1. Type assert to Incrementer
//    - Use: _, ok := i.(Incrementer)
//    - Checks if i has Increment() method
//
// 2. Return the result
//    - Use: return ok
//
// Remember:
// - Counter doesn't satisfy Incrementer (value type, method has pointer receiver)
// - *Counter does satisfy Incrementer (pointer type)
//
// TODO: Implement the CanIncrement function below
// func CanIncrement(i interface{}) bool {
//     _, ok := i.(Incrementer)
//     return ok
// }

// ============================================================================
// Exercise 8: Working with Empty Interface
// ============================================================================

// CountTypes counts how many values of each type are in the slice.
// TODO: Implement type counting using reflection.
//
// The empty interface (interface{}) accepts ANY type.
// Use fmt.Sprintf("%T", v) to get type name as string.
//
// Requirements:
// - Return a map where keys are type names and values are counts
// - Example: []interface{}{1, 2, "hello"} → map[string]int{"int": 2, "string": 1}
//
// Steps to implement:
//
// 1. Create a map to store counts
//    - Use: counts := make(map[string]int)
//    - Keys: type names (string)
//    - Values: how many times we've seen that type (int)
//
// 2. Loop through all values
//    - Use: for _, v := range values { ... }
//    - Each v is interface{} (could be any type)
//
// 3. Get the type name for each value
//    - Use: typeName := fmt.Sprintf("%T", v)
//    - %T format verb prints the type
//    - Examples: "int", "string", "exercise.Person", "*int"
//
// 4. Increment the count for this type
//    - Use: counts[typeName]++
//    - If key doesn't exist, starts at 0 (zero value)
//    - Then increments to 1
//
// 5. Return the counts map
//    - Use: return counts
//
// Key Go concepts:
// - Empty interface can hold any value
// - %T format verb gets type name
// - Map zero value (0) makes counting easy
//
// Example usage:
//   values := []interface{}{1, 2, "hello", "world", 1, true}
//   counts := CountTypes(values)
//   // counts = map[string]int{"int": 3, "string": 2, "bool": 1}
//
// TODO: Implement the CountTypes function below
// func CountTypes(values []interface{}) map[string]int {
//     counts := make(map[string]int)
//     for _, v := range values {
//         typeName := fmt.Sprintf("%T", v)
//         counts[typeName]++
//     }
//     return counts
// }

// ============================================================================
// Exercise 9: The Error Interface
// ============================================================================

// TODO: Implement Error() for ValidationError to satisfy the error interface.
//
// The error interface (built into Go):
//   type error interface {
//       Error() string
//   }
//
// ANY type with Error() string method implements error automatically!
//
// Requirements:
// - Return a string in the format: "validation error on <field>: <message>"
// - Example: ValidationError{Field: "email", Message: "invalid format"}
//   → "validation error on email: invalid format"
//
// Steps to implement:
//
// 1. Format the error message
//    - Use: fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
//    - %s for both strings
//    - Include field name and message
//
// 2. Return the formatted string
//    - This satisfies the Error() string signature
//
// Why value receiver:
// - ValidationError is small (two string headers = 32 bytes)
// - Doesn't need to modify itself
// - Allows both ValidationError and *ValidationError to be used as error
//
// Key Go concepts:
// - error is just an interface (nothing special)
// - Any type with Error() string is an error
// - Custom error types add context (field name, codes, etc.)
//
// Usage:
//   err := ValidationError{Field: "username", Message: "too short"}
//   fmt.Println(err.Error())  // "validation error on username: too short"
//   fmt.Println(err)          // Also calls Error() automatically!
//
// TODO: Implement the Error method below
// func (e ValidationError) Error() string {
//     return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
// }

// ============================================================================
// Exercise 10: Polymorphism Through Interfaces
// ============================================================================

// TODO: Implement Area() for Rectangle.
//
// The Shape interface (defined in types.go):
//   type Shape interface {
//       Area() float64
//   }
//
// Requirements:
// - Return width * height
//
// Steps to implement:
//
// 1. Calculate area
//    - Use: return r.Width * r.Height
//    - Simple multiplication
//
// Why value receiver:
// - Rectangle is small (16 bytes: two float64s)
// - Read-only operation
// - Allows both Rectangle and *Rectangle to satisfy Shape
//
// TODO: Implement the Area method for Rectangle below
// func (r Rectangle) Area() float64 {
//     return r.Width * r.Height
// }

// TODO: Implement Area() for Circle.
//
// Requirements:
// - Return π * radius²
// - Use 3.14159 for π
//
// Steps to implement:
//
// 1. Calculate area
//    - Use: return 3.14159 * c.Radius * c.Radius
//    - π * r²
//
// Why value receiver:
// - Circle is small (8 bytes: one float64)
// - Read-only operation
//
// TODO: Implement the Area method for Circle below
// func (c Circle) Area() float64 {
//     return 3.14159 * c.Radius * c.Radius
// }

// TotalArea calculates the total area of all shapes.
// TODO: Implement polymorphic shape processing.
//
// This demonstrates polymorphism:
// - Works with ANY type that implements Shape
// - Doesn't care if it's Rectangle, Circle, or future types
// - Calls the appropriate Area() method for each shape
//
// Requirements:
// - Sum the areas of all shapes in the slice
// - Return the total
//
// Steps to implement:
//
// 1. Initialize total
//    - Use: total := 0.0
//    - Will accumulate the sum
//
// 2. Loop through shapes
//    - Use: for _, shape := range shapes { ... }
//    - Each shape could be different concrete type
//
// 3. Add each shape's area
//    - Use: total += shape.Area()
//    - Dynamic dispatch: Calls Rectangle.Area() or Circle.Area()
//    - Decided at runtime based on concrete type
//
// 4. Return the total
//    - Use: return total
//
// Key Go concepts:
// - Polymorphism: Same code works with different types
// - Dynamic dispatch: Method call determined at runtime
// - Interface = contract (Shape promises Area() method)
//
// Example usage:
//   shapes := []Shape{
//       Rectangle{Width: 10, Height: 5},   // Area = 50
//       Circle{Radius: 3},                  // Area = 28.27...
//       Rectangle{Width: 2, Height: 3},     // Area = 6
//   }
//   total := TotalArea(shapes)  // 84.27...
//
// TODO: Implement the TotalArea function below
// func TotalArea(shapes []Shape) float64 {
//     total := 0.0
//     for _, shape := range shapes {
//         total += shape.Area()
//     }
//     return total
// }

// After implementing all functions:
// - Run: go test ./minis/13-interfaces-duck-typing/exercise/...
// - Check: go test -v for verbose output
// - Experiment with different types and interfaces
// - Compare with solution.go to see detailed explanations
// - Remember: Interfaces enable polymorphism without inheritance!
