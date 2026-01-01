//go:build reference

package interfacesducktyping

/*
Core Interface Patterns Without Complexity
===========================================

This file focuses on essential interface patterns in Go.
Mastering these concepts is key to writing idiomatic Go code.

CORE CONCEPTS:
- Implicit interface satisfaction
- Type assertions and type switches
- Nil interface gotchas
- Interface composition
- Polymorphism through interfaces

DEBUGGING WORKFLOW:
1. Set breakpoint at function entry
2. Watch interface values (type + data)
3. Trace dynamic method dispatch
4. Observe type assertions

See /RUN_DEBUG.md for comprehensive debugging guide.
*/

// TODO: Implement String() for Person
// BREAKPOINT: Set breakpoint at method start
// DEBUG: Watch 'p' Person value
// DEBUG: Watch formatted string output
//
// Algorithm:
// 1. Use fmt.Sprintf to format string
// 2. Return "Name (Age years old)" format
//
// func (p Person) String() string {
//     // BREAKPOINT: Format string
//     // DEBUG: Watch fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
//     return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
// }

// TODO: Implement GetAge
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch type assertion succeed/fail
// DEBUG: Watch extracted age value
//
// Algorithm:
// 1. Type assert to Person
// 2. Return 0, false if not Person
// 3. Return age, true if Person
//
// func GetAge(s Stringer) (int, bool) {
//     // BREAKPOINT: Type assert
//     // DEBUG: Watch p, ok := s.(Person)
//     p, ok := s.(Person)
//     if !ok {
//         return 0, false
//     }
//     return p.Age, true
// }

// TODO: Implement DescribeType
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch type switch determine concrete type
// DEBUG: Watch each case branch
//
// Algorithm:
// 1. Use type switch on interface{}
// 2. Return description for each type
// 3. Handle int, string, bool, Person, nil, default
//
// func DescribeType(i interface{}) string {
//     // BREAKPOINT: Type switch
//     // DEBUG: Watch v := i.(type)
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

// TODO: Implement IsValidEmail
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch nil interface check
// DEBUG: Watch nil pointer check (the gotcha!)
//
// Algorithm:
// 1. Check if v == nil (interface nil)
// 2. Type assert to *Email
// 3. Check if *Email pointer is nil
// 4. Call IsValid() if non-nil
//
// func IsValidEmail(v Validator) bool {
//     // BREAKPOINT: Check interface nil
//     // DEBUG: Watch v == nil
//     if v == nil {
//         return false
//     }
//
//     // BREAKPOINT: Type assert
//     // DEBUG: Watch e, ok := v.(*Email)
//     e, ok := v.(*Email)
//     if !ok {
//         return v.IsValid()
//     }
//
//     // BREAKPOINT: Check pointer nil (the gotcha!)
//     // DEBUG: Watch e == nil (pointer check)
//     if e == nil {
//         return false
//     }
//
//     return e.IsValid()
// }

// TODO: Implement Buffer.Read and Buffer.Write
// BREAKPOINT: Set breakpoint in each method
// DEBUG: Watch buffer data transformations
//
// Algorithm for Read:
// 1. Return b.data
//
// Algorithm for Write:
// 1. Append data to b.data
// 2. Return nil
//
// func (b *Buffer) Read() string {
//     return b.data
// }
//
// func (b *Buffer) Write(data string) error {
//     b.data += data
//     return nil
// }

// TODO: Implement IsReadWriter
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch composite interface check
//
// Algorithm:
// 1. Type assert to ReadWriter
// 2. Return ok result
//
// func IsReadWriter(i interface{}) bool {
//     _, ok := i.(ReadWriter)
//     return ok
// }

// TODO: Implement Counter.Increment
// BREAKPOINT: Set breakpoint at method start
// DEBUG: Watch c.Value increment
//
// Algorithm:
// 1. Increment c.Value
//
// func (c *Counter) Increment() {
//     c.Value++
// }

// TODO: Implement CanIncrement
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch Incrementer interface check
//
// Algorithm:
// 1. Type assert to Incrementer
// 2. Return ok result
//
// func CanIncrement(i interface{}) bool {
//     _, ok := i.(Incrementer)
//     return ok
// }

// TODO: Implement CountTypes
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch map build up type counts
// DEBUG: Watch %T format verb extract types
//
// Algorithm:
// 1. Create map[string]int
// 2. Loop through values
// 3. Get type name with fmt.Sprintf("%T", v)
// 4. Increment count for that type
// 5. Return map
//
// func CountTypes(values []interface{}) map[string]int {
//     counts := make(map[string]int)
//     for _, v := range values {
//         typeName := fmt.Sprintf("%T", v)
//         counts[typeName]++
//     }
//     return counts
// }

// TODO: Implement ValidationError.Error
// BREAKPOINT: Set breakpoint at method start
// DEBUG: Watch error message formatting
//
// Algorithm:
// 1. Format error message with field and message
// 2. Return formatted string
//
// func (e ValidationError) Error() string {
//     return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
// }

// TODO: Implement Shape methods and TotalArea
// BREAKPOINT: Set breakpoints in Area methods
// DEBUG: Watch polymorphic dispatch in TotalArea
//
// Algorithm for Rectangle.Area:
// 1. Return r.Width * r.Height
//
// Algorithm for Circle.Area:
// 1. Return 3.14159 * c.Radius * c.Radius
//
// Algorithm for TotalArea:
// 1. Initialize total = 0.0
// 2. Loop through shapes
// 3. Add each shape.Area() to total
// 4. Return total
//
// func (r Rectangle) Area() float64 {
//     return r.Width * r.Height
// }
//
// func (c Circle) Area() float64 {
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
