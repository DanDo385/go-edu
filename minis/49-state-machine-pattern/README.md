# Project 49: State Machine Pattern

## 1. What Is This About?

### Real-World Scenario

You're building an order processing system for an e-commerce platform:

**❌ Without state machines:**
1. Order status scattered across multiple boolean flags: `isPaid`, `isShipped`, `isCancelled`
2. Bug: Order gets both shipped AND refunded (impossible state)
3. Invalid transition: Customer cancels after shipment (should be prevented)
4. No audit trail of state changes
5. Spaghetti code with nested if-else checking all flag combinations
6. Team member adds new status, breaks existing logic

**✅ With state machines:**
1. Clear states: Pending → Paid → Shipped → Delivered
2. Impossible states prevented by design
3. Invalid transitions rejected automatically
4. Complete history of state changes with timestamps
5. Clean, declarative transition rules
6. Easy to add new states and transitions

This project teaches you how to build **robust state machines** using:
- **Finite State Automata (FSA)**: Mathematical model of computation
- **State transitions**: Valid paths between states
- **Guards**: Conditional logic for transitions
- **Actions**: Side effects on state changes
- **Type safety**: Compile-time guarantees of correctness

### What You'll Learn

1. **State machine fundamentals**: States, transitions, events
2. **Finite State Automata theory**: DFA, NFA, accepting states
3. **Guard conditions**: Preventing invalid transitions
4. **State pattern**: Object-oriented implementation
5. **Event-driven design**: Triggering state changes
6. **State persistence**: Saving and restoring machine state

### The Challenge

Build a state machine framework with:
- Type-safe state and event definitions
- Declarative transition configuration
- Guard conditions for conditional transitions
- Action hooks (onEnter, onExit, onTransition)
- State history tracking
- Two complete examples: order processing and authentication

---

## 2. First Principles: What Is a State Machine?

### The Core Problem

**Problem**: Systems have complex behavior that changes based on their current condition.

**Example**: A door can be open or closed. You can only close an open door, and open a closed door.

**Reality**: Many systems have dozens of states with complex rules about valid transitions. Managing this with if-else statements becomes unmaintainable.

**Solution**: Model the system as a state machine with:
- **States**: Discrete conditions the system can be in
- **Events**: Triggers that cause state changes
- **Transitions**: Valid paths from one state to another
- **Guards**: Conditions that must be true for a transition
- **Actions**: Side effects that occur during transitions

### Why State Machines Matter

**Real incidents**:
- **Therac-25 (1985-87)**: Radiation therapy machine killed patients due to race conditions in state management
- **Mars Pathfinder (1997)**: Rover failed due to priority inversion in state transitions
- **Toyota unintended acceleration (2009-2010)**: Software state machine bugs caused crashes

**Use cases**:
1. **Order processing**: Shopping carts, payments, fulfillment
2. **Authentication**: Login, MFA, session management
3. **Game logic**: Character states, turn-based gameplay
4. **Network protocols**: TCP state machine, HTTP/2
5. **UI workflows**: Multi-step forms, wizards
6. **Hardware control**: Embedded systems, IoT devices

### Finite State Automata (FSA)

**Mathematical definition**:

A finite state automaton is a 5-tuple: `M = (Q, Σ, δ, q₀, F)`

- **Q**: Finite set of states
- **Σ**: Finite set of input symbols (alphabet/events)
- **δ**: Transition function `δ: Q × Σ → Q`
- **q₀**: Initial state (`q₀ ∈ Q`)
- **F**: Set of accepting (final) states (`F ⊆ Q`)

**Example: Door FSA**

```
Q = {Closed, Open}
Σ = {open, close}
q₀ = Closed
F = {Closed}  (door should end closed)

δ(Closed, open) = Open
δ(Open, close) = Closed
δ(Closed, close) = Closed  (no-op)
δ(Open, open) = Open      (no-op)
```

**Visualization**:

```
       open
    ┌────────┐
    │        ▼
┌─────────┐   ┌──────┐
│ Closed  │   │ Open │
└─────────┘   └──────┘
    ▲        │
    └────────┘
      close
```

### Types of State Machines

| Type | Description | Example |
|------|-------------|---------|
| **Deterministic FSA (DFA)** | Each state has exactly one transition per input | Door lock |
| **Non-deterministic FSA (NFA)** | Multiple possible transitions for same input | Pattern matching |
| **Mealy Machine** | Output depends on current state AND input | Vending machine |
| **Moore Machine** | Output depends only on current state | Traffic light |
| **Hierarchical FSM** | States can contain nested state machines | Game AI |
| **Concurrent FSM** | Multiple state machines running in parallel | Robot control |

**We'll implement Mealy machines** because they're:
- Most common in software
- Output (actions) triggered by transitions
- More expressive than Moore machines

---

## 3. State Machine Core Concepts

### Concept 1: States

**States** represent distinct conditions or modes of a system.

**Characteristics**:
- Mutually exclusive (can only be in one state at a time)
- Discrete (no "in-between" states)
- Enumerable (finite number of states)

**Example: Order states**

```go
type OrderState string

const (
    OrderPending   OrderState = "pending"    // Created, awaiting payment
    OrderPaid      OrderState = "paid"       // Payment received
    OrderShipped   OrderState = "shipped"    // In transit
    OrderDelivered OrderState = "delivered"  // Completed successfully
    OrderCancelled OrderState = "cancelled"  // Cancelled by user/system
    OrderRefunded  OrderState = "refunded"   // Money returned
)
```

**Bad example: Non-discrete states**

```go
// ❌ WRONG: Continuous state
type OrderProgress float64 // 0.0 to 1.0

// What does 0.73 mean? Is payment done? Is shipping started?
// Impossible to reason about!
```

### Concept 2: Events

**Events** are triggers that cause state transitions.

**Characteristics**:
- External stimuli (user action, API call, timer)
- Atomic (happen instantaneously)
- Can carry data (event payload)

**Example: Order events**

```go
type OrderEvent string

const (
    EventPaymentReceived OrderEvent = "payment_received"
    EventShipmentStarted OrderEvent = "shipment_started"
    EventDelivered       OrderEvent = "delivered"
    EventCancel          OrderEvent = "cancel"
    EventRefund          OrderEvent = "refund"
)
```

**Event with payload**:

```go
type Event struct {
    Type      OrderEvent
    Timestamp time.Time
    Metadata  map[string]interface{}
}

// Example: Payment event with amount
paymentEvent := Event{
    Type:      EventPaymentReceived,
    Timestamp: time.Now(),
    Metadata: map[string]interface{}{
        "amount":   99.99,
        "currency": "USD",
        "method":   "credit_card",
    },
}
```

### Concept 3: Transitions

**Transitions** define valid paths between states.

**Transition anatomy**:

```
From State + Event + [Guard] → To State + [Actions]
```

**Example: Order transition table**

| From State | Event | Guard | To State | Actions |
|-----------|-------|-------|----------|---------|
| Pending | payment_received | amount_valid | Paid | send_confirmation |
| Paid | shipment_started | - | Shipped | update_tracking |
| Shipped | delivered | - | Delivered | send_receipt |
| Pending | cancel | - | Cancelled | release_inventory |
| Paid | refund | - | Refunded | process_refund |

**Visual representation**:

```
                 payment_received
    ┌──────────────────────────────────┐
    │                                  ▼
┌─────────┐                        ┌──────┐   shipment_started   ┌──────────┐
│ Pending │                        │ Paid │────────────────────→ │ Shipped  │
└─────────┘                        └──────┘                      └──────────┘
    │                                  │                              │
    │ cancel                           │ refund                       │ delivered
    ▼                                  ▼                              ▼
┌───────────┐                      ┌──────────┐                  ┌────────────┐
│ Cancelled │                      │ Refunded │                  │ Delivered  │
└───────────┘                      └──────────┘                  └────────────┘
```

**Code representation**:

```go
type Transition struct {
    From   OrderState
    Event  OrderEvent
    To     OrderState
    Guard  func(context.Context, *Order) bool
    Action func(context.Context, *Order) error
}

var orderTransitions = []Transition{
    {
        From:  OrderPending,
        Event: EventPaymentReceived,
        To:    OrderPaid,
        Guard: func(ctx context.Context, order *Order) bool {
            return order.Amount > 0 && order.PaymentMethod != ""
        },
        Action: func(ctx context.Context, order *Order) error {
            return sendConfirmationEmail(order.CustomerEmail)
        },
    },
    // ... more transitions
}
```

### Concept 4: Guards

**Guards** are boolean conditions that must be true for a transition to occur.

**Purpose**:
- Enforce business rules
- Validate preconditions
- Prevent invalid transitions

**Example: Authentication state machine with guards**

```go
type AuthState string

const (
    AuthLoggedOut AuthState = "logged_out"
    AuthLoggedIn  AuthState = "logged_in"
    AuthMFAPending AuthState = "mfa_pending"
)

type AuthEvent string

const (
    EventLogin      AuthEvent = "login"
    EventMFARequired AuthEvent = "mfa_required"
    EventMFASuccess  AuthEvent = "mfa_success"
    EventLogout      AuthEvent = "logout"
)

// Guard: Check if user has MFA enabled
func requiresMFA(ctx context.Context, user *User) bool {
    return user.MFAEnabled
}

// Guard: Verify MFA code
func validMFACode(ctx context.Context, user *User) bool {
    code := user.MFACode
    return verifyTOTP(user.MFASecret, code)
}

var authTransitions = []Transition{
    {
        From:  AuthLoggedOut,
        Event: EventLogin,
        To:    AuthMFAPending,
        Guard: requiresMFA,  // Only if MFA enabled
    },
    {
        From:  AuthLoggedOut,
        Event: EventLogin,
        To:    AuthLoggedIn,
        Guard: func(ctx context.Context, user *User) bool {
            return !requiresMFA(ctx, user)  // Only if MFA disabled
        },
    },
    {
        From:  AuthMFAPending,
        Event: EventMFASuccess,
        To:    AuthLoggedIn,
        Guard: validMFACode,  // Only if code is valid
    },
}
```

**Without guards** (bad):

```go
// ❌ WRONG: State machine allows invalid transition
sm.Transition(EventLogin)

// Now we have to check AFTER transition if it was valid
if sm.CurrentState == AuthLoggedIn && user.MFAEnabled && !user.MFAVerified {
    // Oops! User logged in without MFA!
    // Have to manually revert state
    sm.CurrentState = AuthLoggedOut
}
```

**With guards** (correct):

```go
// ✅ CORRECT: Guard prevents invalid transition
err := sm.Transition(EventLogin)
if err != nil {
    // Transition rejected, state unchanged
    log.Printf("Login failed: %v", err)
}
```

### Concept 5: Actions

**Actions** are side effects that occur during transitions.

**Types of actions**:
1. **Entry actions**: Run when entering a state
2. **Exit actions**: Run when leaving a state
3. **Transition actions**: Run during specific transition

**Example: Order actions**

```go
type StateMachine struct {
    current OrderState
    onEnter map[OrderState][]Action
    onExit  map[OrderState][]Action
}

type Action func(context.Context, *Order) error

// Entry action: Log when order is paid
sm.OnEnter(OrderPaid, func(ctx context.Context, order *Order) error {
    log.Printf("Order %s paid: $%.2f", order.ID, order.Amount)
    return nil
})

// Exit action: Clean up when leaving pending state
sm.OnExit(OrderPending, func(ctx context.Context, order *Order) error {
    return releaseInventoryReservation(order.Items)
})

// Transition action: Send email when shipped
sm.AddTransition(Transition{
    From:  OrderPaid,
    Event: EventShipmentStarted,
    To:    OrderShipped,
    Action: func(ctx context.Context, order *Order) error {
        return sendShippingNotification(order.CustomerEmail, order.TrackingNumber)
    },
})
```

**Action execution order**:

```
1. Run exit actions for current state
2. Execute transition action
3. Update current state
4. Run entry actions for new state
```

**Example execution**:

```go
// Current state: Paid
sm.Transition(EventShipmentStarted)

// Execution:
// 1. onExit[Paid]         → Clean up payment processor session
// 2. transition.Action    → Send shipping notification
// 3. current = Shipped    → Update state
// 4. onEnter[Shipped]     → Start tracking updates
```

---

## 4. Building a State Machine Framework

### Step 1: Define Core Types

```go
package statemachine

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// State represents a discrete condition
type State string

// Event represents a trigger for transition
type Event string

// Guard is a condition that must be true for transition
type Guard func(context.Context, interface{}) bool

// Action is a side effect during transition
type Action func(context.Context, interface{}) error

// Transition defines a valid state change
type Transition struct {
    From   State
    Event  Event
    To     State
    Guard  Guard
    Action Action
}

// StateMachine manages state and transitions
type StateMachine struct {
    mu          sync.RWMutex
    current     State
    transitions map[State]map[Event]*Transition
    onEnter     map[State][]Action
    onExit      map[State][]Action
    history     []HistoryEntry
    data        interface{} // User-provided context
}

// HistoryEntry records a state change
type HistoryEntry struct {
    From      State
    Event     Event
    To        State
    Timestamp time.Time
}
```

### Step 2: Initialize State Machine

```go
func New(initial State, data interface{}) *StateMachine {
    return &StateMachine{
        current:     initial,
        transitions: make(map[State]map[Event]*Transition),
        onEnter:     make(map[State][]Action),
        onExit:      make(map[State][]Action),
        history:     []HistoryEntry{},
        data:        data,
    }
}

func (sm *StateMachine) AddTransition(t Transition) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if sm.transitions[t.From] == nil {
        sm.transitions[t.From] = make(map[Event]*Transition)
    }

    sm.transitions[t.From][t.Event] = &t
}

func (sm *StateMachine) OnEnter(state State, action Action) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sm.onEnter[state] = append(sm.onEnter[state], action)
}

func (sm *StateMachine) OnExit(state State, action Action) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sm.onExit[state] = append(sm.onExit[state], action)
}
```

### Step 3: Implement Transition Logic

```go
func (sm *StateMachine) Transition(ctx context.Context, event Event) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    current := sm.current

    // Find transition for current state and event
    transitions, exists := sm.transitions[current]
    if !exists {
        return fmt.Errorf("no transitions defined for state %s", current)
    }

    transition, exists := transitions[event]
    if !exists {
        return fmt.Errorf("no transition for event %s in state %s", event, current)
    }

    // Check guard condition
    if transition.Guard != nil && !transition.Guard(ctx, sm.data) {
        return fmt.Errorf("guard condition failed for transition %s -> %s",
            current, transition.To)
    }

    // Execute exit actions
    for _, action := range sm.onExit[current] {
        if err := action(ctx, sm.data); err != nil {
            return fmt.Errorf("exit action failed: %w", err)
        }
    }

    // Execute transition action
    if transition.Action != nil {
        if err := transition.Action(ctx, sm.data); err != nil {
            return fmt.Errorf("transition action failed: %w", err)
        }
    }

    // Update state
    sm.current = transition.To

    // Execute entry actions
    for _, action := range sm.onEnter[transition.To] {
        if err := action(ctx, sm.data); err != nil {
            // Rollback state change
            sm.current = current
            return fmt.Errorf("entry action failed: %w", err)
        }
    }

    // Record history
    sm.history = append(sm.history, HistoryEntry{
        From:      current,
        Event:     event,
        To:        transition.To,
        Timestamp: time.Now(),
    })

    return nil
}

func (sm *StateMachine) Current() State {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.current
}

func (sm *StateMachine) Can(event Event) bool {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    transitions, exists := sm.transitions[sm.current]
    if !exists {
        return false
    }

    transition, exists := transitions[event]
    if !exists {
        return false
    }

    // Check guard without executing
    if transition.Guard != nil {
        return transition.Guard(context.Background(), sm.data)
    }

    return true
}

func (sm *StateMachine) History() []HistoryEntry {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    // Return copy to prevent external modification
    result := make([]HistoryEntry, len(sm.history))
    copy(result, sm.history)
    return result
}
```

### Step 4: Example - Order Processing

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
)

type OrderState string

const (
    OrderPending   OrderState = "pending"
    OrderPaid      OrderState = "paid"
    OrderShipped   OrderState = "shipped"
    OrderDelivered OrderState = "delivered"
    OrderCancelled OrderState = "cancelled"
)

type OrderEvent string

const (
    EventPay    OrderEvent = "pay"
    EventShip   OrderEvent = "ship"
    EventDeliver OrderEvent = "deliver"
    EventCancel OrderEvent = "cancel"
)

type Order struct {
    ID             string
    CustomerEmail  string
    Amount         float64
    PaymentMethod  string
    TrackingNumber string
}

func setupOrderStateMachine(order *Order) *StateMachine {
    sm := New(State(OrderPending), order)

    // Define transitions
    sm.AddTransition(Transition{
        From:  State(OrderPending),
        Event: Event(EventPay),
        To:    State(OrderPaid),
        Guard: func(ctx context.Context, data interface{}) bool {
            order := data.(*Order)
            return order.Amount > 0 && order.PaymentMethod != ""
        },
        Action: func(ctx context.Context, data interface{}) error {
            order := data.(*Order)
            log.Printf("Processing payment of $%.2f", order.Amount)
            return nil
        },
    })

    sm.AddTransition(Transition{
        From:  State(OrderPaid),
        Event: Event(EventShip),
        To:    State(OrderShipped),
        Action: func(ctx context.Context, data interface{}) error {
            order := data.(*Order)
            order.TrackingNumber = "TRACK-" + order.ID
            log.Printf("Order shipped with tracking: %s", order.TrackingNumber)
            return nil
        },
    })

    sm.AddTransition(Transition{
        From:  State(OrderShipped),
        Event: Event(EventDeliver),
        To:    State(OrderDelivered),
    })

    sm.AddTransition(Transition{
        From:  State(OrderPending),
        Event: Event(EventCancel),
        To:    State(OrderCancelled),
    })

    // Entry/exit actions
    sm.OnEnter(State(OrderPaid), func(ctx context.Context, data interface{}) error {
        order := data.(*Order)
        log.Printf("Sending confirmation email to %s", order.CustomerEmail)
        return nil
    })

    return sm
}
```

---

## 5. Real-World Applications

### Application 1: TCP Connection State Machine

```
         SYN_SENT
             │
             ▼
        ESTABLISHED ←─────┐
             │            │
             ├─ CLOSE_WAIT │
             │            │
             └─ FIN_WAIT  │
                  │       │
                  └───────┘
                      │
                      ▼
                   CLOSED
```

### Application 2: Game Character States

```go
type CharacterState string

const (
    Idle     CharacterState = "idle"
    Walking  CharacterState = "walking"
    Running  CharacterState = "running"
    Jumping  CharacterState = "jumping"
    Attacking CharacterState = "attacking"
    Dead     CharacterState = "dead"
)

// Transitions
// Idle → (walk) → Walking
// Walking → (run) → Running
// Any → (jump) → Jumping
// Any → (attack) → Attacking
// Any → (die) → Dead
// Dead → (respawn) → Idle
```

### Application 3: Document Approval Workflow

```
Draft → (submit) → Pending Review
Pending Review → (approve) → Approved
Pending Review → (reject) → Rejected
Rejected → (revise) → Draft
```

---

## 6. Common Mistakes to Avoid

### Mistake 1: Not Using Guards

**❌ Wrong**:

```go
sm.Transition(ctx, EventPay)
// Oops! Payment processed even though amount was $0
```

**✅ Correct**:

```go
sm.AddTransition(Transition{
    From:  OrderPending,
    Event: EventPay,
    To:    OrderPaid,
    Guard: func(ctx context.Context, data interface{}) bool {
        order := data.(*Order)
        return order.Amount > 0
    },
})
```

### Mistake 2: Allowing Invalid Transitions

**❌ Wrong**:

```go
// Allow cancellation even after delivery
sm.AddTransition(Transition{
    From:  OrderDelivered,
    Event: EventCancel,
    To:    OrderCancelled,
})
```

**✅ Correct**:

```go
// Only allow cancellation from pending state
sm.AddTransition(Transition{
    From:  OrderPending,
    Event: EventCancel,
    To:    OrderCancelled,
})
```

### Mistake 3: Not Handling Failed Actions

**❌ Wrong**:

```go
Action: func(ctx context.Context, data interface{}) error {
    sendEmail(order.Email) // Ignore error
    return nil
}
```

**✅ Correct**:

```go
Action: func(ctx context.Context, data interface{}) error {
    if err := sendEmail(order.Email); err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }
    return nil
}
```

---

## 7. Advanced Patterns

### Pattern 1: Hierarchical State Machines

```go
// Parent state: Vehicle
type VehicleState struct {
    engineState EngineStateMachine
    doorState   DoorStateMachine
}

// Engine can be running/stopped independently of doors
```

### Pattern 2: State Machine Persistence

```go
func (sm *StateMachine) Save() ([]byte, error) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    snapshot := struct {
        Current State
        History []HistoryEntry
    }{
        Current: sm.current,
        History: sm.history,
    }

    return json.Marshal(snapshot)
}

func (sm *StateMachine) Restore(data []byte) error {
    var snapshot struct {
        Current State
        History []HistoryEntry
    }

    if err := json.Unmarshal(data, &snapshot); err != nil {
        return err
    }

    sm.mu.Lock()
    defer sm.mu.Unlock()

    sm.current = snapshot.Current
    sm.history = snapshot.History

    return nil
}
```

### Pattern 3: Event Queuing

```go
type EventQueue struct {
    sm     *StateMachine
    events chan Event
}

func (eq *EventQueue) Process(ctx context.Context) {
    for event := range eq.events {
        if err := eq.sm.Transition(ctx, event); err != nil {
            log.Printf("Transition failed: %v", err)
        }
    }
}

func (eq *EventQueue) Enqueue(event Event) {
    eq.events <- event
}
```

---

## 8. Stretch Goals

### Goal 1: Parallel States ⭐⭐

Implement multiple concurrent state machines.

```go
type ParallelStateMachine struct {
    machines []*StateMachine
}

func (psm *ParallelStateMachine) Transition(event Event) error {
    for _, sm := range psm.machines {
        if err := sm.Transition(context.Background(), event); err != nil {
            // Handle partial failure
        }
    }
    return nil
}
```

### Goal 2: State Machine Visualization ⭐⭐⭐

Generate DOT graph for visualization.

```go
func (sm *StateMachine) ToDot() string {
    var buf bytes.Buffer
    buf.WriteString("digraph StateMachine {\n")

    for from, events := range sm.transitions {
        for event, t := range events {
            buf.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s\"];\n",
                from, t.To, event))
        }
    }

    buf.WriteString("}\n")
    return buf.String()
}

// Use: dot -Tpng graph.dot -o graph.png
```

### Goal 3: Time-based Transitions ⭐⭐

Automatic transitions after timeout.

```go
type TimedTransition struct {
    Transition
    Timeout time.Duration
}

func (sm *StateMachine) AddTimedTransition(t TimedTransition) {
    sm.AddTransition(t.Transition)

    go func() {
        time.Sleep(t.Timeout)
        sm.Transition(context.Background(), t.Event)
    }()
}
```

---

## How to Run

```bash
# Run the demo
cd /home/user/go-edu/minis/49-state-machine-pattern
go run ./cmd/statemachine-demo

# Run tests
cd exercise
go test -v

# Run with solution
go test -v -tags=solution

# Visualize state machine (requires graphviz)
dot -Tpng state_machine.dot -o state_machine.png
```

---

## Summary

**What you learned**:
- ✅ State machine fundamentals (states, events, transitions)
- ✅ Finite State Automata theory (DFA, Mealy machines)
- ✅ Guard conditions for validation
- ✅ Actions for side effects (onEnter, onExit, transition)
- ✅ Type-safe state machine implementation
- ✅ Real-world applications (orders, auth, games)

**Why this matters**:
State machines provide a rigorous framework for managing complex behavior. They make impossible states unrepresentable and invalid transitions impossible, preventing entire classes of bugs.

**Key takeaway**:
When you have complex conditional logic about what can happen in different situations, you probably need a state machine.

**Companies using this**:
- AWS Step Functions: Orchestrate microservices with state machines
- Uber: Trip lifecycle managed by state machines
- Stripe: Payment processing state machines
- Game engines: Unity, Unreal use state machines for AI
- Network protocols: TCP, HTTP/2, WebSocket

**Next steps**:
Identify complex if-else logic in your codebase. Refactor it into a state machine for clarity and correctness.

Build robust systems!
