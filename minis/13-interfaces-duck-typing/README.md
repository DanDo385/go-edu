# 13: Interfaces & Duck Typing

This project explores interfaces, arguably the most powerful and distinctive feature of Go. You'll learn how interfaces enable "duck typing" to create flexible, decoupled, and testable code. This concept is the cornerstone of Go's architectural philosophy.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: Programming to Contracts](#the-big-picture-programming-to-contracts)
- [First Principles: Duck Typing](#first-principles-duck-typing)
- [Project Structure](#project-structure)
- [Key Concepts in Go](#key-concepts-in-go)
  - [1. Defining an Interface](#1-defining-an-interface)
  - [2. Implicit Satisfaction](#2-implicit-satisfaction)
  - [3. The Empty Interface: `any`](#3-the-empty-interface-any)
  - [4. Type Assertions and Type Switches](#4-type-assertions-and-type-switches)
- [Polymorphism in Action](#polymorphism-in-action)
- [Progression: The DNA of Go Libraries](#progression-the-dna-of-go-libraries)
- [How to Run](#how-to-run)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Define interfaces** as a set of method signatures.
-   **Explain implicit satisfaction** and why Go does not have an `implements` keyword.
-   **Describe duck typing** in the context of Go's static, compile-time checks.
-   **Use interfaces to write polymorphic functions** that can accept values of multiple concrete types.
-   **Understand and use the empty interface (`any`)** for handling values of unknown type.
-   **Safely extract concrete types** from interfaces using type assertions and type switches.

## The Big Picture: Programming to Contracts

Imagine you're building a notification system. You might start by sending emails. You write a function `SendEmailNotification(e Email)`. Later, you need to add SMS notifications, so you write `SendSmsNotification(s SMS)`. Now your application is littered with `if/else` statements to decide which function to call. This is rigid and hard to extend.

Interfaces solve this by inverting the dependency. Instead of depending on a *concrete* type (`Email`, `SMS`), your code depends on an *abstract* contract—an interface. You could define a `Notifier` interface:

```go
type Notifier interface {
    Send(message string) error
}
```

Now, you can write a single, simple function: `SendNotification(n Notifier, msg string)`. This function doesn't know or care if it's sending an email, an SMS, or a carrier pigeon. It only knows that it can call the `Send` method. This is the essence of **decoupling**: the sending logic is separate from the notification mechanism.

## First Principles: Duck Typing

Go's approach is often described by the phrase: "If it walks like a duck and it quacks like a duck, then it must be a duck."

-   **The "duck"** is the interface (the contract).
-   **"Walking" and "quacking"** are the methods the interface requires.
-   Any type that has those methods—that "walks" and "quacks"—is considered to be that type of "duck."

Crucially, in Go this is **not** a runtime check as in dynamic languages like Python or Ruby. The Go compiler checks for interface satisfaction at **compile time**. This gives you the flexibility of duck typing with the safety of a statically-typed language.

## Project Structure

```
.
└── cmd/
    └── dev/
        └── main.go       # A program that demonstrates all the interface concepts.
```
-   The `main.go` file defines several types (`Bird`, `Plane`) and interfaces (`Flyer`). It contains functions that demonstrate how these pieces interact, showcasing polymorphism, type assertions, and the empty interface.

## Key Concepts in Go

### 1. Defining an Interface

An interface is a collection of method signatures. It specifies method names, arguments, and return types.

```go
// Any type that wants to be a Flyer MUST have a Fly() method.
type Flyer interface {
    Fly()
}
```

### 2. Implicit Satisfaction

This is the magic of Go. A type satisfies an interface automatically if it implements all the methods defined in the interface. There is no `implements` keyword.

```go
type Bird struct{}

// By implementing this method, Bird implicitly satisfies the Flyer interface.
func (b Bird) Fly() {
    fmt.Println("The bird is flying")
}

type Plane struct{}

// Plane also satisfies the Flyer interface.
func (p Plane) Fly() {
    fmt.Println("The plane is flying")
}
```

Because there's no explicit declaration, interfaces are much easier to create and use. You can define an interface for types that you didn't even write, including types from Go's standard library or third-party packages.

### 3. The Empty Interface: `any`

An interface with zero methods is called the **empty interface**. In modern Go (1.18+), this is written as `any`.

```go
type Any interface {} // or just `any`
```

Since every type has zero or more methods, **every type satisfies the empty interface**. This makes `any` a useful tool for writing functions that need to accept a value of a completely unknown type. However, a variable of type `any` is not very useful on its own; you don't know what you can do with it. To make it useful, you must get the concrete type back out.

### 4. Type Assertions and Type Switches

This is how you "ask" an interface variable what concrete type it's holding.

-   **Type Assertion**: `v.(T)` asserts that the concrete value stored in interface `v` is of type `T`.

    ```go
    var i any = "hello"

    // Simple assertion. Panics if `i` is not a string.
    s := i.(string)

    // "Comma, ok" idiom. Safe, does not panic.
    s, ok := i.(string)
    if ok {
        fmt.Println(s) // "hello"
    }
    ```

-   **Type Switch**: The idiomatic way to handle an interface that could be one of several types.

    ```go
    func doSomething(i any) {
        switch v := i.(type) {
        case string:
            fmt.Printf("It's a string: %s\n", v)
        case int:
            fmt.Printf("It's an int: %d\n", v)
        default:
            fmt.Printf("I don't know this type: %T\n", v)
        }
    }
    ```

## Polymorphism in Action

Here's how interfaces allow for polymorphic behavior. We can write a single function that operates on any `Flyer`.

```go
func MakeItFly(f Flyer) {
    f.Fly()
}

func main() {
    b := Bird{}
    p := Plane{}

    MakeItFly(b) // Prints "The bird is flying"
    MakeItFly(p) // Prints "The plane is flying"
}
```
The `MakeItFly` function is polymorphic. It can take arguments of multiple different types (`Bird`, `Plane`) as long as they satisfy the `Flyer` contract.

## Progression: The DNA of Go Libraries

Interfaces are not just a feature; they are the DNA of Go's standard library and the entire Go ecosystem.
-   The `io.Reader` and `io.Writer` interfaces are the foundation of all I/O in Go, allowing you to read from files, network connections, or in-memory buffers with the same code.
-   The `http.Handler` interface is the basis of Go's web server framework.
-   The `error` interface is how Go handles all errors.

Understanding interfaces moves you from simply using Go to thinking in Go. You'll start designing your own components around small, focused interfaces, leading to cleaner, more modular, and more testable code.

## How to Run

```bash
go run ./cmd/dev/main.go
```
Read the `main.go` source code. Each function demonstrates a different aspect of interfaces. Run the program and trace the output back to the code that produced it.

## Key Takeaways

-   An interface is a contract defined by a set of methods.
-   A type satisfies an interface **implicitly** if it has all the required methods. There is no `implements` keyword.
-   This enables **polymorphism**, allowing you to write functions that operate on multiple concrete types.
-   The empty interface, `any`, can hold a value of any type.
-   Use **type assertions (`v.(T)`)** and **type switches** to inspect the concrete type held by an interface.
-   Programming to interfaces, not implementations, leads to decoupled and flexible systems.

## Further Reading

-   [**A Tour of Go: Interfaces**](https://go.dev/tour/methods/9)
-   [**Go by Example: Interfaces**](https://gobyexample.com/interfaces)
-   [**Effective Go: Interfaces**](https://go.dev/doc/effective_go#interfaces)
-   [**Go Blog: The Laws of Reflection**](https://go.dev/blog/laws-of-reflection) (For a deep understanding of how interfaces work internally).