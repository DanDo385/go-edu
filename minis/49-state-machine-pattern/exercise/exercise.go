//go:build !solution

package exercise

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

// TODO: Implement New
// New creates a new state machine with an initial state and user data
// The data parameter will be passed to all Guard and Action functions
func New(initial State, data interface{}) *StateMachine {
	// EXERCISE: Initialize and return a new StateMachine
	// - Set current to initial
	// - Initialize all maps
	// - Initialize history as empty slice
	// - Set data
	panic("TODO: implement New")
}

// TODO: Implement AddTransition
// AddTransition adds a valid transition to the state machine
func (sm *StateMachine) AddTransition(t Transition) {
	// EXERCISE: Add a transition to the state machine
	// - Lock the mutex (write lock)
	// - Initialize the nested map if needed
	// - Store the transition
	panic("TODO: implement AddTransition")
}

// TODO: Implement OnEnter
// OnEnter registers an action to execute when entering a state
func (sm *StateMachine) OnEnter(state State, action Action) {
	// EXERCISE: Register an entry action
	// - Lock the mutex (write lock)
	// - Append action to onEnter[state]
	panic("TODO: implement OnEnter")
}

// TODO: Implement OnExit
// OnExit registers an action to execute when exiting a state
func (sm *StateMachine) OnExit(state State, action Action) {
	// EXERCISE: Register an exit action
	// - Lock the mutex (write lock)
	// - Append action to onExit[state]
	panic("TODO: implement OnExit")
}

// TODO: Implement Transition
// Transition attempts to transition from current state using the given event
// Returns error if:
// - No transition exists for current state and event
// - Guard condition fails
// - Any action fails
func (sm *StateMachine) Transition(ctx context.Context, event Event) error {
	// EXERCISE: Implement the complete transition logic
	// 1. Lock the mutex (write lock)
	// 2. Find the transition for current state and event
	// 3. If guard exists, check it (return error if fails)
	// 4. Execute exit actions for current state
	// 5. Execute transition action (if exists)
	// 6. Update current state
	// 7. Execute entry actions for new state
	// 8. Record in history
	//
	// IMPORTANT: If any action fails, consider rollback strategy
	panic("TODO: implement Transition")
}

// TODO: Implement Current
// Current returns the current state of the state machine
func (sm *StateMachine) Current() State {
	// EXERCISE: Return current state (with read lock)
	panic("TODO: implement Current")
}

// TODO: Implement Can
// Can checks if a transition is possible from the current state
// It checks both transition existence and guard condition
func (sm *StateMachine) Can(event Event) bool {
	// EXERCISE: Check if transition is possible
	// - Lock with read lock
	// - Check if transition exists
	// - Check guard condition (if exists)
	panic("TODO: implement Can")
}

// TODO: Implement History
// History returns a copy of all state transitions
func (sm *StateMachine) History() []HistoryEntry {
	// EXERCISE: Return a copy of history (with read lock)
	// - Lock with read lock
	// - Create a copy of history slice
	// - Return the copy (to prevent external modification)
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

// TODO: Implement NewOrderStateMachine
// NewOrderStateMachine creates a state machine for order processing
func NewOrderStateMachine(order *Order) *StateMachine {
	// EXERCISE: Create and configure an order state machine
	// 1. Create new state machine with OrderPending initial state
	// 2. Add transitions:
	//    - Pending -> Paid (event: pay, guard: amount > 0 && payment method exists)
	//    - Paid -> Shipped (event: ship, action: set tracking number)
	//    - Shipped -> Delivered (event: deliver, action: set delivered time)
	//    - Pending -> Cancelled (event: cancel)
	// 3. Add entry/exit actions:
	//    - OnEnter Paid: log payment confirmation
	//    - OnEnter Shipped: log shipping notification
	//    - OnEnter Delivered: log delivery confirmation
	//
	// HINT: Use type assertions to access order in guards/actions
	// Example: order := data.(*Order)
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

// TODO: Implement NewAuthStateMachine
// NewAuthStateMachine creates a state machine for user authentication
func NewAuthStateMachine(user *User) *StateMachine {
	// EXERCISE: Create and configure an authentication state machine
	// 1. Create new state machine with AuthLoggedOut initial state
	// 2. Add transitions:
	//    - LoggedOut -> MFAPending (event: login, guard: MFA enabled)
	//    - LoggedOut -> LoggedIn (event: login, guard: MFA NOT enabled)
	//    - MFAPending -> LoggedIn (event: mfa_success, guard: valid MFA code)
	//    - LoggedIn -> LoggedOut (event: logout)
	// 3. Add actions:
	//    - OnEnter LoggedIn: generate session ID, log successful login
	//    - OnEnter MFAPending: log MFA challenge sent
	//    - OnExit LoggedIn: clear session, log logout
	//
	// HINT: For simplified MFA validation, check if code == secret
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
