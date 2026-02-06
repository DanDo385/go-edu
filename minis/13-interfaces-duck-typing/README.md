# 13: Interfaces Duck Typing

## Core Concepts

- The concrete problem in Interfaces Duck Typing and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Interfaces Duck Typing patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for interfaces duck typing.

At this point in the arc:
Lesson 13 introduces a sharper systems concern so later modules can assume this mental model is stable.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Define the smallest valid behavior and reject invalid input or impossible state early.

### Step 2: Why This Approach
Pick a direct design that keeps control flow and data flow visible for debugging and testing.

### Step 3: Memory / Pointer Impact
Call out where data is copied versus aliased, and where mutable shared state needs synchronization.

### Step 4: What Changed
Produce a stable result shape and explicit error behavior that downstream code can rely on.

## Pointer and Indirection

- Explain * and & in this module when they appear in code or docs.
- Show memory-before and memory-after when data ownership changes.
- Clarify common misconceptions: Go stays pass-by-value even when pointer values are copied.
- Primer link: docs/MEMORY_POINTERS_PRIMER.md

## Verify


a) learner path


go test -v ./...


b) reference path


go test -tags=reference -v ./...


This project explores interfaces, arguably the most powerful feature of Go. You'll learn how interfaces enable "duck typing" to create flexible, decoupled, and testable code. This concept is the cornerstone of Go's architectural philosophy.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.


## What You'll Learn

- How to define **interfaces** as a set of method signatures.
- How Go's **implicit satisfaction** of interfaces works.
- How to write **polymorphic** functions that accept multiple types.
- How to use the **empty interface (`any`)** and **type switches**.

## The Challenge: Programming to Contracts

Imagine building a notification system. You start with email, then need to add SMS, then push notifications. If you write specific functions for each (`SendEmail`, `SendSMS`), your code becomes a tangled mess of `if/else` statements.

Interfaces solve this by defining a **contract**. Instead of depending on a *concrete* type (`Email`, `SMS`), your code depends on an *abstract* `Notifier` interface.

```go
type Notifier interface {
    Send(message string) error
}

// This function doesn't care what the notifier IS, only what it CAN DO.
func SendNotification(n Notifier, msg string) {
    n.Send(msg)
}
```

## Core Concepts

### Implicit Satisfaction & Duck Typing
In Go, a type satisfies an interface **automatically** and **implicitly** if it implements all the methods defined in the interface. There is no `implements` keyword.

This is often called **duck typing**: "If it walks like a duck and it quacks like a duck, then it must be a duck."
- The **interface** is the "duck".
- The **methods** are the "walking" and "quacking".

Any type that has the right methods is a "duck" and can be used where that interface is required. The Go compiler checks this for you, giving you the flexibility of duck typing with the safety of a statically-typed language.

### The Empty Interface: `any`
An interface with zero methods is called the empty interface, written as `any`. Since every type has zero or more methods, **every type satisfies the empty interface**. This allows you to write a function that accepts a value of any type.

### Type Assertions and Switches
What do you do with a value of type `any` or another interface? To access its underlying concrete value, you use a **type switch**.
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

## Your Task

This is another conceptual lesson. Your task is to **read, run, and understand** the code in `cmd/app/main.go`.

The `main.go` file defines several types (`Bird`, `Plane`, `Car`) and interfaces (`Flyer`, `Mover`). It contains functions that demonstrate:
- How different types can satisfy the same interface.
- How a polymorphic function (`MakeItFly`) can accept any type that satisfies the `Flyer` interface.
- How to use a type switch to inspect the concrete type behind an interface.

## How to Verify Your Work

1.  **Read the code:** Open `cmd/app/main.go` and read through the source code and its comments.
2.  **Run the program:**
    ```bash
    go run ./cmd/app/main.go
    ```
3.  **Analyze the output:** Trace the program's output back to the code that produced it. Make sure you understand why `MakeItFly(b)` and `MakeItFly(p)` both work, and how the type switch is able to differentiate between the different types in the `things` slice.

Once you understand each demonstration, you have completed the lesson.

## Related Lessons
- Previous: `minis/12-pointers-zero-values-nil-gotchas`
- Next: `minis/14-methods-value-vs-pointer-receivers`
