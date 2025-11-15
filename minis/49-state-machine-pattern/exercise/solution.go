//go:build solution

package exercise

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// State represents a discrete condition in a state machine
type State string

// Event represents a trigger that causes state transitions
type Event string

// Guard is a conditional function that must return true for a transition to occur
type Guard func(context.Context, interface{}) bool

// Action is a side effect function that executes during transitions
type Action func(context.Context, interface{}) error

// Transition defines a valid state change in the state machine
type Transition struct {
	From   State
	Event  Event
	To     State
	Guard  Guard
	Action Action
}

// HistoryEntry records a state change
type HistoryEntry struct {
	From      State
	Event     Event
	To        State
	Timestamp time.Time
}

// StateMachine manages state transitions
type StateMachine struct {
	mu          sync.RWMutex
	current     State
	transitions map[State]map[Event][]*Transition
	onEnter     map[State][]Action
	onExit      map[State][]Action
	history     []HistoryEntry
	data        interface{}
}

// New creates a new state machine with an initial state and user data
func New(initial State, data interface{}) *StateMachine {
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
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.transitions[t.From] == nil {
		sm.transitions[t.From] = make(map[Event][]*Transition)
	}

	sm.transitions[t.From][t.Event] = append(sm.transitions[t.From][t.Event], &t)
}

// OnEnter registers an action to execute when entering a state
func (sm *StateMachine) OnEnter(state State, action Action) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.onEnter[state] = append(sm.onEnter[state], action)
}

// OnExit registers an action to execute when exiting a state
func (sm *StateMachine) OnExit(state State, action Action) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.onExit[state] = append(sm.onExit[state], action)
}

// Transition attempts to transition from current state using the given event
func (sm *StateMachine) Transition(ctx context.Context, event Event) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	current := sm.current

	// Find transitions for current state and event
	stateTransitions, exists := sm.transitions[current]
	if !exists {
		return fmt.Errorf("no transitions defined for state %s", current)
	}

	candidateTransitions, exists := stateTransitions[event]
	if !exists {
		return fmt.Errorf("no transition for event %s in state %s", event, current)
	}

	// Find first transition whose guard passes (or has no guard)
	var transition *Transition
	for _, t := range candidateTransitions {
		if t.Guard == nil || t.Guard(ctx, sm.data) {
			transition = t
			break
		}
	}

	if transition == nil {
		return fmt.Errorf("guard condition failed for all transitions from %s on event %s", current, event)
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

// Current returns the current state of the state machine
func (sm *StateMachine) Current() State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current
}

// Can checks if a transition is possible from the current state
func (sm *StateMachine) Can(event Event) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stateTransitions, exists := sm.transitions[sm.current]
	if !exists {
		return false
	}

	candidateTransitions, exists := stateTransitions[event]
	if !exists {
		return false
	}

	// Check if any transition's guard passes (or has no guard)
	for _, transition := range candidateTransitions {
		if transition.Guard == nil || transition.Guard(context.Background(), sm.data) {
			return true
		}
	}

	return false
}

// History returns a copy of all state transitions
func (sm *StateMachine) History() []HistoryEntry {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Return copy to prevent external modification
	result := make([]HistoryEntry, len(sm.history))
	copy(result, sm.history)
	return result
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
	sm := New(State(OrderPending), order)

	// Transition: Pending -> Paid (on payment received)
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
			log.Printf("[Order %s] Processing payment of $%.2f via %s",
				order.ID, order.Amount, order.PaymentMethod)
			return nil
		},
	})

	// Transition: Paid -> Shipped (on shipment started)
	sm.AddTransition(Transition{
		From:  State(OrderPaid),
		Event: Event(EventShip),
		To:    State(OrderShipped),
		Action: func(ctx context.Context, data interface{}) error {
			order := data.(*Order)
			order.TrackingNumber = generateTrackingNumber(order.ID)
			log.Printf("[Order %s] Shipment started with tracking: %s",
				order.ID, order.TrackingNumber)
			return nil
		},
	})

	// Transition: Shipped -> Delivered (on delivery)
	sm.AddTransition(Transition{
		From:  State(OrderShipped),
		Event: Event(EventDeliver),
		To:    State(OrderDelivered),
		Action: func(ctx context.Context, data interface{}) error {
			order := data.(*Order)
			order.DeliveredAt = time.Now()
			log.Printf("[Order %s] Delivered at %s",
				order.ID, order.DeliveredAt.Format(time.RFC3339))
			return nil
		},
	})

	// Transition: Pending -> Cancelled (cancel before payment)
	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventCancel),
		To:    State(OrderCancelled),
		Action: func(ctx context.Context, data interface{}) error {
			order := data.(*Order)
			log.Printf("[Order %s] Cancelled by customer", order.ID)
			return nil
		},
	})

	// Entry action: When order becomes paid
	sm.OnEnter(State(OrderPaid), func(ctx context.Context, data interface{}) error {
		order := data.(*Order)
		log.Printf("[Order %s] Sending confirmation email to %s",
			order.ID, order.CustomerEmail)
		return nil
	})

	// Entry action: When order is shipped
	sm.OnEnter(State(OrderShipped), func(ctx context.Context, data interface{}) error {
		order := data.(*Order)
		log.Printf("[Order %s] Sending shipping notification to %s",
			order.ID, order.CustomerEmail)
		return nil
	})

	// Entry action: When order is delivered
	sm.OnEnter(State(OrderDelivered), func(ctx context.Context, data interface{}) error {
		order := data.(*Order)
		log.Printf("[Order %s] Sending delivery receipt to %s",
			order.ID, order.CustomerEmail)
		return nil
	})

	// Exit action: When leaving pending state
	sm.OnExit(State(OrderPending), func(ctx context.Context, data interface{}) error {
		order := data.(*Order)
		log.Printf("[Order %s] Leaving pending state", order.ID)
		return nil
	})

	return sm
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
	MFACode    string
	SessionID  string
}

// NewAuthStateMachine creates a state machine for user authentication
func NewAuthStateMachine(user *User) *StateMachine {
	sm := New(State(AuthLoggedOut), user)

	// Transition: LoggedOut -> MFAPending (login with MFA enabled)
	sm.AddTransition(Transition{
		From:  State(AuthLoggedOut),
		Event: Event(EventLogin),
		To:    State(AuthMFAPending),
		Guard: func(ctx context.Context, data interface{}) bool {
			user := data.(*User)
			return user.MFAEnabled
		},
		Action: func(ctx context.Context, data interface{}) error {
			user := data.(*User)
			log.Printf("[User %s] MFA required, sending challenge", user.Email)
			return nil
		},
	})

	// Transition: LoggedOut -> LoggedIn (login without MFA)
	sm.AddTransition(Transition{
		From:  State(AuthLoggedOut),
		Event: Event(EventLogin),
		To:    State(AuthLoggedIn),
		Guard: func(ctx context.Context, data interface{}) bool {
			user := data.(*User)
			return !user.MFAEnabled
		},
		Action: func(ctx context.Context, data interface{}) error {
			user := data.(*User)
			log.Printf("[User %s] Login successful (no MFA)", user.Email)
			return nil
		},
	})

	// Transition: MFAPending -> LoggedIn (MFA success)
	sm.AddTransition(Transition{
		From:  State(AuthMFAPending),
		Event: Event(EventMFASuccess),
		To:    State(AuthLoggedIn),
		Guard: func(ctx context.Context, data interface{}) bool {
			user := data.(*User)
			// Simplified: check if code matches secret
			return user.MFACode == user.MFASecret
		},
		Action: func(ctx context.Context, data interface{}) error {
			user := data.(*User)
			log.Printf("[User %s] MFA verification successful", user.Email)
			return nil
		},
	})

	// Transition: LoggedIn -> LoggedOut (logout)
	sm.AddTransition(Transition{
		From:  State(AuthLoggedIn),
		Event: Event(EventLogout),
		To:    State(AuthLoggedOut),
		Action: func(ctx context.Context, data interface{}) error {
			user := data.(*User)
			log.Printf("[User %s] Logged out", user.Email)
			return nil
		},
	})

	// Entry action: When user logs in
	sm.OnEnter(State(AuthLoggedIn), func(ctx context.Context, data interface{}) error {
		user := data.(*User)
		user.SessionID = generateSessionID()
		log.Printf("[User %s] Session created: %s", user.Email, user.SessionID)
		return nil
	})

	// Entry action: When MFA is pending
	sm.OnEnter(State(AuthMFAPending), func(ctx context.Context, data interface{}) error {
		user := data.(*User)
		log.Printf("[User %s] Awaiting MFA code verification", user.Email)
		return nil
	})

	// Exit action: When user logs out
	sm.OnExit(State(AuthLoggedIn), func(ctx context.Context, data interface{}) error {
		user := data.(*User)
		log.Printf("[User %s] Clearing session: %s", user.Email, user.SessionID)
		user.SessionID = ""
		return nil
	})

	return sm
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// generateSessionID generates a simple session ID
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// generateTrackingNumber generates a tracking number for shipments
func generateTrackingNumber(orderID string) string {
	return fmt.Sprintf("TRACK-%s-%d", orderID, time.Now().Unix())
}
