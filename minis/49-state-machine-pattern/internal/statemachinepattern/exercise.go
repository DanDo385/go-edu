//go:build !solution && !reference

package statemachinepattern

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type State string

type Event string

type Guard func(context.Context, interface{}) bool

type Action func(context.Context, interface{}) error

type Transition struct {
	From   State
	Event  Event
	To     State
	Guard  Guard
	Action Action
}

type HistoryEntry struct {
	From      State
	Event     Event
	To        State
	Timestamp time.Time
}

type StateMachine struct {
	mu          sync.RWMutex
	current     State
	transitions map[State]map[Event][]*Transition
	onEnter     map[State][]Action
	onExit      map[State][]Action
	history     []HistoryEntry
	data        interface{}
}

type OrderState string

type OrderEvent string

type Order struct {
	ID             string
	CustomerEmail  string
	Amount         float64
	PaymentMethod  string
	TrackingNumber string
	DeliveredAt    time.Time
}

type AuthState string

type AuthEvent string

type User struct {
	ID         string
	Email      string
	MFAEnabled bool
	MFASecret  string
	MFACode    string
	SessionID  string
}

// New implements the exercise.
//
// TODO: Implement this function
func New(initial State, data interface{}) *StateMachine {
	// TODO: Implement
	return nil
}

// AddTransition implements the exercise.
//
// TODO: Implement this function
func (sm *StateMachine) AddTransition(t Transition) {
	// TODO: Implement
}

// OnEnter implements the exercise.
//
// TODO: Implement this function
func (sm *StateMachine) OnEnter(state State, action Action) {
	// TODO: Implement
}

// OnExit implements the exercise.
//
// TODO: Implement this function
func (sm *StateMachine) OnExit(state State, action Action) {
	// TODO: Implement
}

// Transition implements the exercise.
//
// TODO: Implement this function
func (sm *StateMachine) Transition(ctx context.Context, event Event) error {
	// TODO: Implement
	return nil
}

// Current implements the exercise.
//
// TODO: Implement this function
func (sm *StateMachine) Current() State {
	// TODO: Implement
	return State{}
}

// Can implements the exercise.
//
// TODO: Implement this function
func (sm *StateMachine) Can(event Event) bool {
	// TODO: Implement
	return false
}

// History implements the exercise.
//
// TODO: Implement this function
func (sm *StateMachine) History() []HistoryEntry {
	// TODO: Implement
	return nil
}

// NewOrderStateMachine implements the exercise.
//
// TODO: Implement this function
func NewOrderStateMachine(order *Order) *StateMachine {
	// TODO: Implement
	return nil
}

// NewAuthStateMachine implements the exercise.
//
// TODO: Implement this function
func NewAuthStateMachine(user *User) *StateMachine {
	// TODO: Implement
	return nil
}
