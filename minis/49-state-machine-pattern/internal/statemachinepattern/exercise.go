//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package statemachinepattern

import (
	"context"

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
// TODO: implement New.
func New(initial State, data interface{}) *StateMachine { panic("TODO: implement") }
// TODO: implement AddTransition.
func (sm *StateMachine) AddTransition(t Transition) { panic("TODO: implement") }
// TODO: implement OnEnter.
func (sm *StateMachine) OnEnter(state State, action Action) { panic("TODO: implement") }
// TODO: implement OnExit.
func (sm *StateMachine) OnExit(state State, action Action) { panic("TODO: implement") }
// TODO: implement Transition.
func (sm *StateMachine) Transition(ctx context.Context, event Event) error { panic("TODO: implement") }
// TODO: implement Current.
func (sm *StateMachine) Current() State { panic("TODO: implement") }
// TODO: implement Can.
func (sm *StateMachine) Can(event Event) bool { panic("TODO: implement") }
// TODO: implement History.
func (sm *StateMachine) History() []HistoryEntry { panic("TODO: implement") }

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
	EventPay     OrderEvent = "pay"
	EventShip    OrderEvent = "ship"
	EventDeliver OrderEvent = "deliver"
	EventCancel  OrderEvent = "cancel"
)

type Order struct {
	ID             string
	CustomerEmail  string
	Amount         float64
	PaymentMethod  string
	TrackingNumber string
	DeliveredAt    time.Time
}
// TODO: implement NewOrderStateMachine.
func NewOrderStateMachine(order *Order) *StateMachine { panic("TODO: implement") }

type AuthState string

const (
	AuthLoggedOut  AuthState = "logged_out"
	AuthLoggedIn   AuthState = "logged_in"
	AuthMFAPending AuthState = "mfa_pending"
)

type AuthEvent string

const (
	EventLogin      AuthEvent = "login"
	EventMFASuccess AuthEvent = "mfa_success"
	EventLogout     AuthEvent = "logout"
)

type User struct {
	ID         string
	Email      string
	MFAEnabled bool
	MFASecret  string
	MFACode    string
	SessionID  string
}
// TODO: implement NewAuthStateMachine.
func NewAuthStateMachine(user *User) *StateMachine { panic("TODO: implement") }
// TODO: implement generateSessionID.
func generateSessionID() string { panic("TODO: implement") }
// TODO: implement generateTrackingNumber.
func generateTrackingNumber(orderID string) string { panic("TODO: implement") }
