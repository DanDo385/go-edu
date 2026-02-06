# Go Edu

Go Edu is a learning-first systems engineering codebase.

The objective is not just "write working Go", but:

1. Build data structures and algorithms from first principles.
2. Understand Go as a systems and infrastructure language.
3. Apply those ideas to concurrency-heavy blockchain-style services.

The code is the textbook. Tests are the correctness contract. Memory models are the core mental model.

## Start Here

1. Read the full learning sequence in `docs/LEARNING_ARC.md`.
2. Read pointer and memory foundations in `docs/MEMORY_POINTERS_PRIMER.md`.
3. Read lesson authoring and review standards in `docs/LESSON_BLUEPRINT.md`.

## Repository Shape

### `minis/`

Small, focused labs. Each lesson centers one concept and usually contains:

- `internal/<lesson>/exercise.go` (student path)
- `internal/<lesson>/solution.reference.go` (reference path)
- `internal/<lesson>/exercise_test.go` (correctness checks)

### `geth/`

Applied Ethereum/client-infrastructure track. Lessons compose concepts from `minis/` into realistic systems patterns:

- RPC data flow
- state and lifecycle management
- concurrent pipelines
- indexing and monitoring

## Three-Stage Learning Arc

The repository progression is explicitly organized as:

1. Data Structures and Algorithms from first principles
2. Go as a systems language (memory, concurrency, networking, process behavior)
3. Infrastructure and blockchain context (long-running services and client patterns)

The detailed lesson order and prerequisites are in `docs/LEARNING_ARC.md`.

## Daily Workflow

From any lesson folder:

```bash
go test -v ./...
go run ./cmd/app/...
```

From repo root:

```bash
make list
make test P=minis/12-pointers-zero-values-nil-gotchas
make test P=geth/16-concurrency
make verify-lessons
make verify-teaching
```

## First Lesson

If you are starting fresh:

```bash
cd minis/01-hello-strings
go test -v ./...
```
