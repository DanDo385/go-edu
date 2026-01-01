//go:build !solution && !reference


package statemachinepattern

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// State represents a discrete condition in a state machine
type State string

// Event represents a trigger that causes state transitions
type Event string

// Guard is a conditional function that must return true for a transition to occur
// It receives the context and the user-provided data
type Guard func(context.Context, interface{}) bool

// Action is a side effect function that executes during transitions
// It can perform operations like logging, sending notifications, etc.
type Action func(context.Context, interface{}) error

// Transition defines a valid state change in the state machine
type Transition struct {
	From   State  // Source state
	Event  Event  // Triggering event
	To     State  // Destination state
	Guard  Guard  // Optional condition for transition
	Action Action // Optional action to execute
}

// HistoryEntry records a state change
type HistoryEntry struct {
	From      State     // Previous state
	Event     Event     // Event that triggered the transition
	To        State     // New state
	Timestamp time.Time // When the transition occurred
}

// StateMachine manages state transitions
type StateMachine struct {
	mu          sync.RWMutex
	current     State
	transitions map[State]map[Event][]*Transition
	onEnter     map[State][]Action
	onExit      map[State][]Action
	history     []HistoryEntry
	data        interface{} // User-provided context data
}

// New creates a new state machine with an initial state and user data
func New(initial State, data interface{}) *StateMachine {
	// EXERCISE: Initialize and return a new StateMachine.
	// - Set the `current` state to `initial`.
	// - Initialize the `transitions`, `onEnter`, and `onExit` maps so they are ready to be used.
	// - The `history` slice should be an empty slice.
	// - Store the user-provided `data`.
	return &StateMachine{
		current:     initial,
		transitions: make(map[State]map[Event][]*Transition),
		onEnter:     make(map[State][]Action),
		onExit:      make(map[State][]Action),
		history:     []HistoryEntry{},
		data:        data,
	}
}

// AddTransition adds a valid transition to the state machine
func (sm *StateMachine) AddTransition(t Transition) {
	// EXERCISE: Add a transition to the state machine's configuration.
	// - Use a write lock (`sm.mu.Lock()`) to safely modify the `transitions` map.
	// - Check if the map for the `t.From` state has been initialized. If not, create it.
	// - Append the new transition `t` to the slice for the given state and event.
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.transitions[t.From] == nil {
		sm.transitions[t.From] = make(map[Event][]*Transition)
	}

	sm.transitions[t.From][t.Event] = append(sm.transitions[t.From][t.Event], &t)
}

// OnEnter registers an action to execute when entering a state
func (sm *StateMachine) OnEnter(state State, action Action) {
	// EXERCISE: Register an action to be called when a state is entered.
	// - Use a write lock.
	// - Append the `action` to the slice of actions for the given `state` in the `onEnter` map.
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.onEnter[state] = append(sm.onEnter[state], action)
}

// OnExit registers an action to execute when exiting a state
func (sm *StateMachine) OnExit(state State, action Action) {
	// EXERCISE: Register an action to be called when a state is exited.
	// - Use a write lock.
	// - Append the `action` to the slice of actions for the given `state` in the `onExit` map.
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.onExit[state] = append(sm.onExit[state], action)
}

// Transition attempts to transition from current state using the given event
func (sm *StateMachine) Transition(ctx context.Context, event Event) error {
	// EXERCISE: Implement the complete transition logic.
	// This is the heart of the state machine.

	// Step 1: Acquire a write lock, as this method modifies the machine's state.
	// Step 2: Find all possible transitions for the current state and the given event.
	// Step 3: Iterate through the possible transitions to find the first one that is allowed.
	//   - A transition is allowed if it has no `Guard` function, or if its `Guard` function returns `true`.
	// Step 4: If no allowed transition is found, return an error.
	// Step 5: If an allowed transition is found:
	//   a. Execute all `OnExit` actions for the *current* state. If any action returns an error, stop and return that error.
	//   b. Execute the `Action` associated with the transition itself (if it exists). Handle any error.
	//   c. Change the state: `sm.current = transition.To`.
	//   d. Execute all `OnEnter` actions for the *new* state. If any action fails, you might need to consider a rollback strategy (for this exercise, just returning the error is fine).
	//   e. Record the successful transition in the `sm.history`.
	// Step 6: If all actions were successful, return `nil`.
	panic("TODO: implement Transition")
}

// Current returns the current state of the state machine
func (sm *StateMachine) Current() State {
	// EXERCISE: Return the current state.
	// - Use a read lock (`sm.mu.RLock()`) for thread-safe access.
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current
}

// Can checks if a transition is possible from the current state
func (sm *StateMachine) Can(event Event) bool {
	// EXERCISE: Check if a transition is possible.
	// - Use a read lock.
	// - Find the transitions for the current state and event.
	// - Return `true` if there is at least one transition for which the `Guard` function passes (or there is no guard).
	panic("TODO: implement Can")
}

// History returns a copy of all state transitions
func (sm *StateMachine) History() []HistoryEntry {
	// EXERCISE: Return a copy of the history.
	// - Use a read lock.
	// - It's important to return a *copy* of the slice, not the slice itself, to prevent the caller from modifying the internal history.
	//   `result := make([]HistoryEntry, len(sm.history)); copy(result, sm.history); return result`
	panic("TODO: implement History")
}

// ============================================================================
// ORDER PROCESSING STATE MACHINE
// ============================================================================

// OrderState represents the state of an order
type OrderState string

const (
	OrderPending   OrderState = "pending"
	OrderPaid      OrderState = "paid"
	OrderShipped   OrderState = "shipped"
	OrderDelivered OrderState = "delivered"
	OrderCancelled OrderState = "cancelled"
)

// OrderEvent represents events that trigger order state changes
type OrderEvent string

const (
	EventPay     OrderEvent = "pay"
	EventShip    OrderEvent = "ship"
	EventDeliver OrderEvent = "deliver"
	EventCancel  OrderEvent = "cancel"
)

// Order represents an e-commerce order
type Order struct {
	ID             string
	CustomerEmail  string
	Amount         float64
	PaymentMethod  string
	TrackingNumber string
	DeliveredAt    time.Time
}

// NewOrderStateMachine creates a state machine for order processing
func NewOrderStateMachine(order *Order) *StateMachine {
	// EXERCISE: Create and configure a state machine for a specific use case: an e-commerce order.

	// Step 1: Create a new state machine instance.
	// - The initial state is `OrderPending`.
	// - The user-provided `data` is the `order` itself.
	//   `sm := New(State(OrderPending), order)`

	// Step 2: Add all the valid transitions for the order lifecycle.
	// - Use `sm.AddTransition()` for each one.
	// - Example for Pending -> Paid:
	/*
		sm.AddTransition(Transition{
			From:  State(OrderPending),
			Event: Event(EventPay),
			To:    State(OrderPaid),
			Guard: func(ctx context.Context, data interface{}) bool {
				// Access the order via type assertion
				order := data.(*Order)
				// The payment can only succeed if the amount is valid.
				return order.Amount > 0 && order.PaymentMethod != ""
			},
			Action: func(ctx context.Context, data interface{}) error {
				// Simulate processing the payment.
				fmt.Printf("Processing payment for order %s\n", data.(*Order).ID)
				return nil
			},
		})
	*/
	// - Implement the other transitions: Paid -> Shipped, Shipped -> Delivered, Pending -> Cancelled.

	// Step 3: Add entry/exit actions for specific states.
	// - Use `sm.OnEnter()` to define actions that happen as soon as a state is entered.
	// - Example for OnEnter Paid:
	/*
		sm.OnEnter(State(OrderPaid), func(ctx context.Context, data interface{}) error {
			fmt.Printf("Sending payment confirmation email for order %s\n", data.(*Order).ID)
			return nil
		})
	*/
	// - Implement other entry actions for Shipped and Delivered.

	// Step 4: Return the configured state machine.
	panic("TODO: implement NewOrderStateMachine")
}

// ============================================================================
// AUTHENTICATION STATE MACHINE
// ============================================================================

// AuthState represents authentication states
type AuthState string

const (
	AuthLoggedOut  AuthState = "logged_out"
	AuthLoggedIn   AuthState = "logged_in"
	AuthMFAPending AuthState = "mfa_pending"
)

// AuthEvent represents authentication events
type AuthEvent string

const (
	EventLogin      AuthEvent = "login"
	EventMFASuccess AuthEvent = "mfa_success"
	EventLogout     AuthEvent = "logout"
)

// User represents a user account
type User struct {
	ID         string
	Email      string
	MFAEnabled bool
	MFASecret  string
	MFACode    string // Code provided by user
	SessionID  string
}

// NewAuthStateMachine creates a state machine for user authentication
func NewAuthStateMachine(user *User) *StateMachine {
	// EXERCISE: Create and configure an authentication state machine.

	// This demonstrates how multiple transitions can exist for the same event from the same state,
	// with `Guard` functions determining which path is taken.

	// Step 1: Create a new state machine with the initial state `AuthLoggedOut`.
	// Step 2: Add transitions for the `login` event.
	//   - One transition from `LoggedOut` to `MFAPending` for the `login` event, with a `Guard` that checks if `user.MFAEnabled` is true.
	//   - A second transition from `LoggedOut` to `LoggedIn` for the same `login` event, with a `Guard` that checks if `user.MFAEnabled` is false.
	// Step 3: Add the transition for MFA success.
	//   - `MFAPending` -> `LoggedIn` on the `mfa_success` event. Include a `Guard` that checks if the provided `user.MFACode` is valid.
	// Step 4: Add the logout transition.
	//   - `LoggedIn` -> `LoggedOut` on the `logout` event.
	// Step 5: Add entry/exit actions.
	//   - `OnEnter(AuthLoggedIn, ...)`: a function that generates and sets the `user.SessionID`.
	//   - `OnExit(AuthLoggedIn, ...)`: a function that clears the `user.SessionID`.
	// Step 6: Return the configured state machine.
	panic("TODO: implement NewAuthStateMachine")
}

// ============================================================================
// HELPER FUNCTIONS (you can use these)
// ============================================================================

// generateSessionID generates a simple session ID
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// generateTrackingNumber generates a tracking number for shipments
func generateTrackingNumber(orderID string) string {
	return fmt.Sprintf("TRACK-%s-%d", orderID, time.Now().Unix())
}
