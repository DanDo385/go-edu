package exercise

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// ============================================================================
// CORE STATE MACHINE TESTS
// ============================================================================

func TestStateMachine_New(t *testing.T) {
	data := &Order{ID: "TEST-001"}
	sm := New(State(OrderPending), data)

	if sm == nil {
		t.Fatal("New() returned nil")
	}

	if sm.Current() != State(OrderPending) {
		t.Errorf("Expected initial state %s, got %s", OrderPending, sm.Current())
	}
}

func TestStateMachine_AddTransition(t *testing.T) {
	sm := New(State(OrderPending), nil)

	transition := Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
	}

	sm.AddTransition(transition)

	// Verify transition was added by attempting to use it
	ctx := context.Background()
	err := sm.Transition(ctx, Event(EventPay))

	if err != nil {
		t.Errorf("Transition failed: %v", err)
	}

	if sm.Current() != State(OrderPaid) {
		t.Errorf("Expected state %s, got %s", OrderPaid, sm.Current())
	}
}

func TestStateMachine_Transition_Simple(t *testing.T) {
	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
	})

	ctx := context.Background()

	// Initial state
	if sm.Current() != State(OrderPending) {
		t.Errorf("Expected initial state %s, got %s", OrderPending, sm.Current())
	}

	// Transition
	err := sm.Transition(ctx, Event(EventPay))
	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	// New state
	if sm.Current() != State(OrderPaid) {
		t.Errorf("Expected state %s, got %s", OrderPaid, sm.Current())
	}
}

func TestStateMachine_Transition_InvalidEvent(t *testing.T) {
	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
	})

	ctx := context.Background()

	// Try invalid event
	err := sm.Transition(ctx, Event(EventShip))
	if err == nil {
		t.Error("Expected error for invalid event, got nil")
	}

	// State should remain unchanged
	if sm.Current() != State(OrderPending) {
		t.Errorf("Expected state to remain %s, got %s", OrderPending, sm.Current())
	}
}

func TestStateMachine_Transition_WithGuard(t *testing.T) {
	order := &Order{
		ID:            "TEST-001",
		Amount:        100.0,
		PaymentMethod: "credit_card",
	}

	sm := New(State(OrderPending), order)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
		Guard: func(ctx context.Context, data interface{}) bool {
			order := data.(*Order)
			return order.Amount > 0 && order.PaymentMethod != ""
		},
	})

	ctx := context.Background()

	// Transition should succeed
	err := sm.Transition(ctx, Event(EventPay))
	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	if sm.Current() != State(OrderPaid) {
		t.Errorf("Expected state %s, got %s", OrderPaid, sm.Current())
	}
}

func TestStateMachine_Transition_GuardFails(t *testing.T) {
	order := &Order{
		ID:            "TEST-001",
		Amount:        0, // Invalid amount
		PaymentMethod: "credit_card",
	}

	sm := New(State(OrderPending), order)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
		Guard: func(ctx context.Context, data interface{}) bool {
			order := data.(*Order)
			return order.Amount > 0 && order.PaymentMethod != ""
		},
	})

	ctx := context.Background()

	// Transition should fail due to guard
	err := sm.Transition(ctx, Event(EventPay))
	if err == nil {
		t.Error("Expected error due to guard failure, got nil")
	}

	// State should remain unchanged
	if sm.Current() != State(OrderPending) {
		t.Errorf("Expected state to remain %s, got %s", OrderPending, sm.Current())
	}
}

func TestStateMachine_Transition_WithAction(t *testing.T) {
	actionExecuted := false

	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
		Action: func(ctx context.Context, data interface{}) error {
			actionExecuted = true
			return nil
		},
	})

	ctx := context.Background()
	err := sm.Transition(ctx, Event(EventPay))

	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	if !actionExecuted {
		t.Error("Transition action was not executed")
	}

	if sm.Current() != State(OrderPaid) {
		t.Errorf("Expected state %s, got %s", OrderPaid, sm.Current())
	}
}

func TestStateMachine_Transition_ActionFails(t *testing.T) {
	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
		Action: func(ctx context.Context, data interface{}) error {
			return fmt.Errorf("simulated action failure")
		},
	})

	ctx := context.Background()
	err := sm.Transition(ctx, Event(EventPay))

	if err == nil {
		t.Error("Expected error due to action failure, got nil")
	}

	// State should remain unchanged when action fails
	if sm.Current() != State(OrderPending) {
		t.Errorf("Expected state to remain %s, got %s", OrderPending, sm.Current())
	}
}

func TestStateMachine_OnEnter(t *testing.T) {
	entryExecuted := false

	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
	})

	sm.OnEnter(State(OrderPaid), func(ctx context.Context, data interface{}) error {
		entryExecuted = true
		return nil
	})

	ctx := context.Background()
	err := sm.Transition(ctx, Event(EventPay))

	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	if !entryExecuted {
		t.Error("Entry action was not executed")
	}
}

func TestStateMachine_OnExit(t *testing.T) {
	exitExecuted := false

	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
	})

	sm.OnExit(State(OrderPending), func(ctx context.Context, data interface{}) error {
		exitExecuted = true
		return nil
	})

	ctx := context.Background()
	err := sm.Transition(ctx, Event(EventPay))

	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	if !exitExecuted {
		t.Error("Exit action was not executed")
	}
}

func TestStateMachine_ActionExecutionOrder(t *testing.T) {
	order := []string{}

	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
		Action: func(ctx context.Context, data interface{}) error {
			order = append(order, "transition")
			return nil
		},
	})

	sm.OnExit(State(OrderPending), func(ctx context.Context, data interface{}) error {
		order = append(order, "exit")
		return nil
	})

	sm.OnEnter(State(OrderPaid), func(ctx context.Context, data interface{}) error {
		order = append(order, "enter")
		return nil
	})

	ctx := context.Background()
	err := sm.Transition(ctx, Event(EventPay))

	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	// Verify execution order: exit -> transition -> enter
	expected := []string{"exit", "transition", "enter"}
	if len(order) != len(expected) {
		t.Fatalf("Expected %d actions, got %d", len(expected), len(order))
	}

	for i, action := range expected {
		if order[i] != action {
			t.Errorf("Action %d: expected %s, got %s", i, action, order[i])
		}
	}
}

func TestStateMachine_Can(t *testing.T) {
	order := &Order{
		ID:            "TEST-001",
		Amount:        100.0,
		PaymentMethod: "credit_card",
	}

	sm := New(State(OrderPending), order)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
		Guard: func(ctx context.Context, data interface{}) bool {
			order := data.(*Order)
			return order.Amount > 0 && order.PaymentMethod != ""
		},
	})

	sm.AddTransition(Transition{
		From:  State(OrderPaid),
		Event: Event(EventShip),
		To:    State(OrderShipped),
	})

	// Can pay from pending state (guard passes)
	if !sm.Can(Event(EventPay)) {
		t.Error("Expected Can(EventPay) to return true")
	}

	// Cannot ship from pending state (no transition)
	if sm.Can(Event(EventShip)) {
		t.Error("Expected Can(EventShip) to return false")
	}

	// Make order invalid
	order.Amount = 0

	// Now cannot pay (guard fails)
	if sm.Can(Event(EventPay)) {
		t.Error("Expected Can(EventPay) to return false when guard fails")
	}
}

func TestStateMachine_History(t *testing.T) {
	sm := New(State(OrderPending), nil)

	sm.AddTransition(Transition{
		From:  State(OrderPending),
		Event: Event(EventPay),
		To:    State(OrderPaid),
	})

	sm.AddTransition(Transition{
		From:  State(OrderPaid),
		Event: Event(EventShip),
		To:    State(OrderShipped),
	})

	ctx := context.Background()

	// Initial history should be empty
	history := sm.History()
	if len(history) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(history))
	}

	// Make first transition
	sm.Transition(ctx, Event(EventPay))

	history = sm.History()
	if len(history) != 1 {
		t.Fatalf("Expected 1 history entry, got %d", len(history))
	}

	if history[0].From != State(OrderPending) {
		t.Errorf("Expected From=%s, got %s", OrderPending, history[0].From)
	}

	if history[0].To != State(OrderPaid) {
		t.Errorf("Expected To=%s, got %s", OrderPaid, history[0].To)
	}

	if history[0].Event != Event(EventPay) {
		t.Errorf("Expected Event=%s, got %s", EventPay, history[0].Event)
	}

	// Make second transition
	sm.Transition(ctx, Event(EventShip))

	history = sm.History()
	if len(history) != 2 {
		t.Fatalf("Expected 2 history entries, got %d", len(history))
	}

	if history[1].From != State(OrderPaid) {
		t.Errorf("Expected From=%s, got %s", OrderPaid, history[1].From)
	}

	if history[1].To != State(OrderShipped) {
		t.Errorf("Expected To=%s, got %s", OrderShipped, history[1].To)
	}
}

// ============================================================================
// ORDER STATE MACHINE TESTS
// ============================================================================

func TestOrderStateMachine_HappyPath(t *testing.T) {
	order := &Order{
		ID:            "ORDER-001",
		CustomerEmail: "test@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	sm := NewOrderStateMachine(order)
	ctx := context.Background()

	// Initial state
	if sm.Current() != State(OrderPending) {
		t.Errorf("Expected initial state %s, got %s", OrderPending, sm.Current())
	}

	// Pay
	err := sm.Transition(ctx, Event(EventPay))
	if err != nil {
		t.Fatalf("Payment transition failed: %v", err)
	}

	if sm.Current() != State(OrderPaid) {
		t.Errorf("Expected state %s after payment, got %s", OrderPaid, sm.Current())
	}

	// Ship
	err = sm.Transition(ctx, Event(EventShip))
	if err != nil {
		t.Fatalf("Ship transition failed: %v", err)
	}

	if sm.Current() != State(OrderShipped) {
		t.Errorf("Expected state %s after shipping, got %s", OrderShipped, sm.Current())
	}

	if order.TrackingNumber == "" {
		t.Error("Expected tracking number to be set")
	}

	// Deliver
	err = sm.Transition(ctx, Event(EventDeliver))
	if err != nil {
		t.Fatalf("Deliver transition failed: %v", err)
	}

	if sm.Current() != State(OrderDelivered) {
		t.Errorf("Expected state %s after delivery, got %s", OrderDelivered, sm.Current())
	}

	if order.DeliveredAt.IsZero() {
		t.Error("Expected delivered timestamp to be set")
	}

	// Verify history
	history := sm.History()
	if len(history) != 3 {
		t.Errorf("Expected 3 transitions in history, got %d", len(history))
	}
}

func TestOrderStateMachine_InvalidPayment(t *testing.T) {
	tests := []struct {
		name   string
		order  *Order
		reason string
	}{
		{
			name: "zero amount",
			order: &Order{
				ID:            "ORDER-002",
				CustomerEmail: "test@example.com",
				Amount:        0,
				PaymentMethod: "credit_card",
			},
			reason: "zero amount",
		},
		{
			name: "missing payment method",
			order: &Order{
				ID:            "ORDER-003",
				CustomerEmail: "test@example.com",
				Amount:        99.99,
				PaymentMethod: "",
			},
			reason: "missing payment method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewOrderStateMachine(tt.order)
			ctx := context.Background()

			err := sm.Transition(ctx, Event(EventPay))
			if err == nil {
				t.Errorf("Expected payment to fail for %s, but it succeeded", tt.reason)
			}

			// State should remain pending
			if sm.Current() != State(OrderPending) {
				t.Errorf("Expected state to remain %s, got %s", OrderPending, sm.Current())
			}
		})
	}
}

func TestOrderStateMachine_Cancel(t *testing.T) {
	order := &Order{
		ID:            "ORDER-004",
		CustomerEmail: "test@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	sm := NewOrderStateMachine(order)
	ctx := context.Background()

	// Cancel from pending state
	err := sm.Transition(ctx, Event(EventCancel))
	if err != nil {
		t.Fatalf("Cancel transition failed: %v", err)
	}

	if sm.Current() != State(OrderCancelled) {
		t.Errorf("Expected state %s after cancel, got %s", OrderCancelled, sm.Current())
	}
}

func TestOrderStateMachine_CannotCancelAfterDelivery(t *testing.T) {
	order := &Order{
		ID:            "ORDER-005",
		CustomerEmail: "test@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	sm := NewOrderStateMachine(order)
	ctx := context.Background()

	// Process order to delivered
	sm.Transition(ctx, Event(EventPay))
	sm.Transition(ctx, Event(EventShip))
	sm.Transition(ctx, Event(EventDeliver))

	// Try to cancel
	err := sm.Transition(ctx, Event(EventCancel))
	if err == nil {
		t.Error("Expected cancel to fail after delivery, but it succeeded")
	}

	// State should remain delivered
	if sm.Current() != State(OrderDelivered) {
		t.Errorf("Expected state to remain %s, got %s", OrderDelivered, sm.Current())
	}
}

// ============================================================================
// AUTHENTICATION STATE MACHINE TESTS
// ============================================================================

func TestAuthStateMachine_WithoutMFA(t *testing.T) {
	user := &User{
		ID:         "USER-001",
		Email:      "user@example.com",
		MFAEnabled: false,
	}

	sm := NewAuthStateMachine(user)
	ctx := context.Background()

	// Initial state
	if sm.Current() != State(AuthLoggedOut) {
		t.Errorf("Expected initial state %s, got %s", AuthLoggedOut, sm.Current())
	}

	// Login (should go directly to logged in)
	err := sm.Transition(ctx, Event(EventLogin))
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if sm.Current() != State(AuthLoggedIn) {
		t.Errorf("Expected state %s after login, got %s", AuthLoggedIn, sm.Current())
	}

	if user.SessionID == "" {
		t.Error("Expected session ID to be set")
	}

	// Logout
	err = sm.Transition(ctx, Event(EventLogout))
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	if sm.Current() != State(AuthLoggedOut) {
		t.Errorf("Expected state %s after logout, got %s", AuthLoggedOut, sm.Current())
	}

	if user.SessionID != "" {
		t.Error("Expected session ID to be cleared")
	}
}

func TestAuthStateMachine_WithMFA_Success(t *testing.T) {
	user := &User{
		ID:         "USER-002",
		Email:      "user@example.com",
		MFAEnabled: true,
		MFASecret:  "secret123",
	}

	sm := NewAuthStateMachine(user)
	ctx := context.Background()

	// Login (should go to MFA pending)
	err := sm.Transition(ctx, Event(EventLogin))
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if sm.Current() != State(AuthMFAPending) {
		t.Errorf("Expected state %s after login with MFA, got %s", AuthMFAPending, sm.Current())
	}

	// Provide correct MFA code
	user.MFACode = "secret123"
	err = sm.Transition(ctx, Event(EventMFASuccess))
	if err != nil {
		t.Fatalf("MFA verification failed: %v", err)
	}

	if sm.Current() != State(AuthLoggedIn) {
		t.Errorf("Expected state %s after MFA success, got %s", AuthLoggedIn, sm.Current())
	}

	if user.SessionID == "" {
		t.Error("Expected session ID to be set after successful MFA")
	}
}

func TestAuthStateMachine_WithMFA_Failure(t *testing.T) {
	user := &User{
		ID:         "USER-003",
		Email:      "user@example.com",
		MFAEnabled: true,
		MFASecret:  "secret123",
	}

	sm := NewAuthStateMachine(user)
	ctx := context.Background()

	// Login
	sm.Transition(ctx, Event(EventLogin))

	// Provide incorrect MFA code
	user.MFACode = "wrong_code"
	err := sm.Transition(ctx, Event(EventMFASuccess))
	if err == nil {
		t.Error("Expected MFA verification to fail with wrong code")
	}

	// Should remain in MFA pending state
	if sm.Current() != State(AuthMFAPending) {
		t.Errorf("Expected state to remain %s, got %s", AuthMFAPending, sm.Current())
	}
}

// ============================================================================
// CONCURRENCY TESTS
// ============================================================================

func TestStateMachine_ConcurrentReads(t *testing.T) {
	sm := New(State(OrderPending), nil)

	done := make(chan bool)

	// Spawn multiple goroutines reading current state
	for i := 0; i < 100; i++ {
		go func() {
			_ = sm.Current()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestStateMachine_ConcurrentTransitions(t *testing.T) {
	order := &Order{
		ID:            "ORDER-CONCURRENT",
		CustomerEmail: "test@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	sm := NewOrderStateMachine(order)
	ctx := context.Background()

	// First transition to establish a state
	sm.Transition(ctx, Event(EventPay))

	done := make(chan bool)

	// Try to trigger same transition from multiple goroutines
	// Only one should succeed
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		go func() {
			err := sm.Transition(ctx, Event(EventShip))
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Only one transition should have succeeded
	if successCount != 1 {
		t.Errorf("Expected 1 successful transition, got %d", successCount)
	}

	if sm.Current() != State(OrderShipped) {
		t.Errorf("Expected final state %s, got %s", OrderShipped, sm.Current())
	}
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkStateMachine_Transition(b *testing.B) {
	order := &Order{
		ID:            "BENCH-001",
		CustomerEmail: "bench@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sm := NewOrderStateMachine(order)
		sm.Transition(ctx, Event(EventPay))
	}
}

func BenchmarkStateMachine_Can(b *testing.B) {
	order := &Order{
		ID:            "BENCH-002",
		CustomerEmail: "bench@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	sm := NewOrderStateMachine(order)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sm.Can(Event(EventPay))
	}
}
