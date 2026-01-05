# 14: Methods: Value vs. Pointer Receivers

This project focuses on a fundamental and crucial decision in Go: when to define a method on a value receiver (`func (t T)`) versus a pointer receiver (`func (t *T)`). This choice is not just syntactic sugar; it's a core statement about your type's intended use, its data ownership model, and its performance characteristics.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: A Decision About Mutation and Efficiency](#the-big-picture-a-decision-about-mutation-and-efficiency)
- [First Principles: Recap of Values vs. Pointers](#first-principles-recap-of-values-vs-pointers)
- [Project Structure](#project-structure)
- [The Two Types of Receivers](#the-two-types-of-receivers)
  - [1. Value Receivers: `func (s MyStruct) ...`](#1-value-receivers-func-s-mystruct-)
  - [2. Pointer Receivers: `func (s *MyStruct) ...`](#2-pointer-receivers-func-s-mystruct-)
- [How to Choose: A Simple Guideline](#how-to-choose-a-simple-guideline)
- [A Note on Interfaces and Method Sets](#a-note-on-interfaces-and-method-sets)
- [Progression: Designing Your Own Types](#progression-designing-your-own-types)
- [How to Run](#how-to-run)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Explain the mechanical difference** between a value and a pointer receiver.
-   **Predict whether a method will modify** the original struct based on its receiver type.
-   **Make a clear, reasoned decision** on which receiver type to use for your own methods.
-   **Understand the performance implications** of copying large structs.
-   **Follow the idiomatic Go convention** for receiver consistency.

## The Big Picture: A Decision About Mutation and Efficiency

When you define a type with methods, you are creating a small, self-contained abstraction. The choice of receiver type answers two critical questions about that abstraction:
1.  **Can this object's state be changed by its methods?** (Mutation)
2.  **What is the performance cost of using this object?** (Efficiency)

Getting this right is essential for creating APIs (whether internal or public) that are clear, predictable, and performant.

## First Principles: Recap of Values vs. Pointers

This project builds directly on the concepts from **Project 12**.
-   A **value** is a self-contained piece of data. When you pass a value, you pass a **copy**.
-   A **pointer** is a variable that holds the memory address of a value. When you pass a pointer, you copy the pointer, but both the original and the copy "point to" the **same underlying data**.

The receiver of a method is just a special kind of function argument. The same rules of copying apply.

## Project Structure

```
.
└── cmd/
    └── dev/
        └── main.go       # A program demonstrating the behavior of both receiver types.
```
-   The `main.go` file defines a simple `Counter` type and then implements methods with both value and pointer receivers to make the difference in their behavior explicit and observable.

## The Two Types of Receivers

### 1. Value Receivers: `func (s MyStruct) ...`

A method with a value receiver operates on a **copy** of the original value.

```go
type Counter struct {
    value int
}

// Increment is a VALUE receiver method.
// `c` is a copy of the original Counter.
func (c Counter) Increment() {
    c.value++
    // This modification happens only to the copy.
}
```

**Memory Visualization:**
```
c1 := Counter{value: 10}
c1.Increment() // A copy of c1 is made and passed to Increment.

Original `c1` in main(): [ value: 10 ] at address 0x100
Copy `c` in Increment(): [ value: 11 ] at address 0x200 (a different location)

// After the call, the copy at 0x200 is gone. `c1` is unchanged.
```

**Use a value receiver when:**
-   The method does **not** need to modify the receiver.
-   The method is a simple getter, formatter, or calculation that only needs to read the receiver's state.
-   The receiver is a small, simple `struct` (or a map, slice, or channel type) where the cost of copying is negligible.

### 2. Pointer Receivers: `func (s *MyStruct) ...`

A method with a pointer receiver operates on a **pointer** to the original value. This allows the method to modify the original value.

```go
// Increment is now a POINTER receiver method.
// `c` is a pointer to the original Counter.
func (c *Counter) Increment() {
    c.value++ // This modification affects the original struct.
}
```
*Note: Go permits `c.value` as a convenient shortcut for the more formally correct `(*c).value`.*

**Memory Visualization:**
```
c1 := Counter{value: 10}
(&c1).Increment() // A pointer to c1 is passed.

Original `c1` in main(): [ value: 10 ] at address 0x100
Pointer `c` in Increment(): holds address 0x100

// Inside the method, we follow the pointer and modify the data at 0x100.
Original `c1` becomes [ value: 11 ]

// After the call, `c1` has been permanently changed.
```

**Use a pointer receiver when:**
-   The method **needs to mutate** the receiver's state. This is the most common reason.
-   The `struct` is very large. Using a pointer avoids copying the entire struct on every method call, which can be a significant performance win.

## How to Choose: A Simple Guideline

Follow this decision process, as recommended by the official Go documentation:

1.  **Does the method need to modify the receiver?**
    -   If **YES**, you **must** use a pointer receiver (`*T`).

2.  If NO, you have a choice. Now consider consistency and efficiency.
    -   **Consistency Rule**: If *any* method on the type has a pointer receiver, the other methods should also have pointer receivers to make the type's usage consistent and predictable.
    -   **Efficiency Rule**: Is the struct large? If so, a pointer receiver is more efficient as it only copies a pointer, not the entire struct's data.
    -   If the type is small and simple (like a basic `struct` or a reference type like a slice or map), a value receiver is perfectly fine.

**Conclusion: When in doubt, use a pointer receiver.** It's the most common and often the correct choice.

## A Note on Interfaces and Method Sets

The choice of receiver affects whether a type satisfies an interface. The rule is:
-   A type `T` has a method set containing only methods with a receiver of `T`.
-   A type `*T` has a method set containing methods with a receiver of both `T` and `*T`.

This means a pointer type can satisfy more interfaces than the corresponding value type. However, Go often helps you by automatically taking the address of a value (`T`) to make it a `*T` if needed to satisfy an interface. While this magic is convenient, understanding the underlying rule is important for edge cases.

## Progression: Designing Your Own Types

This project solidifies your ability to design your own types, a skill you started developing in **Project 03 (CSV Stats)** and **Project 05 (CLI App)**. Your choice of receiver type is a fundamental part of a type's API design. It communicates intent: is this type a simple, immutable value, or is it a stateful object that can be modified? Mastering this choice is a key step toward writing clear, idiomatic, and professional Go code.

## How to Run

```bash
go run ./cmd/dev/main.go
```
The program will run two tests: one using a value receiver `Increment` and one using a pointer receiver `IncrementPtr`. Observe how only the pointer receiver version actually changes the original counter's value.

## Key Takeaways

-   **Value receivers (`T`) operate on a copy; the original is unchanged.**
-   **Pointer receivers (`*T`) operate on the original value (via a pointer) and can modify it.**
-   The primary reason to use a pointer receiver is to **mutate the receiver's state**.
-   A secondary reason is to **avoid copying large structs** for efficiency.
-   For consistency, if a type has any pointer receiver methods, all its methods should probably have pointer receivers.
-   **When in doubt, use a pointer receiver.**

## Further Reading

-   [**A Tour of Go: Methods and pointer indirection**](https://go.dev/tour/methods/4)
-   [**Go FAQ: Should I define methods on values or pointers?**](https://go.dev/doc/faq#methods_on_values_or_pointers) - The definitive answer.
-   [**Go by Example: Methods**](https://gobyexample.com/methods)