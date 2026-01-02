# 49-state-machine-pattern

**State Machine Pattern**

Implement a state machine for workflow management.

## What You'll Learn

- State machine design
- Transition validation
- Event handling
- State persistence

## Functions to Implement

| Function | Description |
|----------|-------------|
| Define states | State enum |
| Transition | Validate and execute |
| Handle events | Trigger transitions |

## Project Structure

```
49-state-machine-pattern/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/statemachinepattern/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/49-state-machine-pattern

# Run state machine demo
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run demo
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **States**: Distinct conditions
2. **Transitions**: Valid state changes
3. **Events**: Triggers for transitions
4. **Guards**: Conditional transitions

## Next Steps

After completing this exercise, proceed to `minis/50-mini-service-all-features`.
