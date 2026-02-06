# Lesson Blueprint

Use this blueprint for any new lesson or substantial refactor.

## 1. Pre-Coding Lecture (Required)

Start the lesson with three explicit sections:

1. Core concepts involved.
2. CS fundamentals connection.
3. End-state understanding.

Template:

```md
## Core Concepts
- stack vs heap
- values vs pointers
- aliasing and mutability
- goroutines/channels or other relevant runtime behavior

## CS Connection
- memory layout
- addresses vs values
- what is copied vs what is shared
- how synchronization affects safety

## End-State
- what the learner should be able to explain and implement
```

## 2. Implementation Steps (Required)

Do not jump straight to the final solution. Use small steps.

For each step:

1. State the problem this step solves.
2. Explain why this design was chosen.
3. Explain memory impact (copying, aliasing, escape if relevant).
4. Implement code.
5. Summarize what changed.

## 3. `*` and `&` Rule (Highest Priority)

Whenever `*` or `&` appears in lesson text or code, explain:

1. Meaning in this exact context.
2. Memory state before operation.
3. Memory state after operation.
4. Whether value, pointer value, or dereferenced data is being passed.
5. Whether data is copied, aliased, or mutated in place.

Always call out misconceptions:

1. `*` in a type means pointer type, not dereference.
2. `*` in an expression can mean dereference or multiplication.
3. `&` creates an address value; it does not switch Go to pass-by-reference.
4. Pointer receiver methods do not guarantee heap allocation.

## 4. Required Files per Lesson

For non-trivial modules:

1. `internal/<lesson>/exercise.go`
2. `internal/<lesson>/solution.reference.go`
3. `internal/<lesson>/exercise_test.go`
4. `README.md`
5. Optional: `types.go`, `example_test.go`, `cmd/app/main.go`

## 5. Reference Solution Standard

`solution.reference.go` should be:

1. Correct and idiomatic.
2. No extra abstraction beyond what the lesson taught.
3. Clearly structured for readability.
4. Commented only where reasoning would otherwise be ambiguous.

In lesson text, explicitly compare:

1. Intermediate student steps vs reference.
2. Why the final shape is idiomatic Go.
3. Which tradeoffs were resolved.

## 6. Testing Standard

Every lesson should define:

1. `go test` command(s).
2. What each test proves.
3. Edge cases and failure modes.
4. Concurrency risks and race-avoidance strategy when applicable.

Minimum checks:

1. Happy path.
2. Invalid input or boundary behavior.
3. Resource cleanup/cancellation behavior for long-running patterns.
4. Race-safe behavior for concurrent modules (`go test -race` where useful).

Repository-level contract check:

1. Run `make verify-lessons` from repo root to validate required lesson files.
2. Run `make verify-teaching` from repo root to validate required README lecture sections.

## 7. Debugging from First Principles

When diagnosing bugs, explain state transitions explicitly:

1. Which value changed.
2. Which goroutine changed it.
3. Whether change happened through aliasing (`[]T`, `map`, `*T`).
4. Why observed behavior follows from memory and scheduling.

## 8. Review Checklist

Before merging a lesson:

1. Pre-coding lecture sections are present.
2. Stepwise implementation narrative is present.
3. `*` and `&` uses are fully explained.
4. `exercise.go` and `solution.reference.go` both exist and agree on behavior.
5. Tests pass and cover critical edge conditions.
6. `make verify-lessons` passes.
7. `make verify-teaching` passes.
