package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/example/go-10x-minis/minis/49-state-machine-pattern/exercise"
)

func main() {
	fmt.Println("=== State Machine Pattern Demonstrations ===\n")

	// Demo 1: Order Processing State Machine
	fmt.Println("--- Demo 1: Order Processing State Machine ---")
	demoOrderProcessing()

	fmt.Println("\n" + separator())

	// Demo 2: Authentication Flow State Machine
	fmt.Println("\n--- Demo 2: Authentication Flow State Machine ---")
	demoAuthenticationFlow()

	fmt.Println("\n" + separator())

	// Demo 3: Invalid Transitions
	fmt.Println("\n--- Demo 3: Invalid Transitions (Guards in Action) ---")
	demoInvalidTransitions()

	fmt.Println("\n" + separator())

	// Demo 4: State History Tracking
	fmt.Println("\n--- Demo 4: State History Tracking ---")
	demoStateHistory()

	fmt.Println("\n" + separator())

	// Demo 5: Concurrent State Machines
	fmt.Println("\n--- Demo 5: Concurrent State Machines ---")
	demoConcurrentStateMachines()
}

// demoOrderProcessing demonstrates a complete order lifecycle
func demoOrderProcessing() {
	fmt.Println("Creating a new order...")

	order := &exercise.Order{
		ID:            "ORDER-001",
		CustomerEmail: "customer@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	sm := exercise.NewOrderStateMachine(order)
	ctx := context.Background()

	fmt.Printf("Initial state: %s\n\n", sm.Current())

	// Step 1: Process payment
	fmt.Println("Step 1: Processing payment...")
	if err := sm.Transition(ctx, exercise.Event(exercise.EventPay)); err != nil {
		log.Fatalf("Payment failed: %v", err)
	}
	fmt.Printf("Current state: %s\n\n", sm.Current())

	// Step 2: Ship the order
	fmt.Println("Step 2: Shipping order...")
	if err := sm.Transition(ctx, exercise.Event(exercise.EventShip)); err != nil {
		log.Fatalf("Shipping failed: %v", err)
	}
	fmt.Printf("Current state: %s\n", sm.Current())
	fmt.Printf("Tracking number: %s\n\n", order.TrackingNumber)

	// Step 3: Deliver the order
	fmt.Println("Step 3: Delivering order...")
	if err := sm.Transition(ctx, exercise.Event(exercise.EventDeliver)); err != nil {
		log.Fatalf("Delivery failed: %v", err)
	}
	fmt.Printf("Current state: %s\n", sm.Current())
	fmt.Printf("Delivered at: %s\n", order.DeliveredAt.Format(time.RFC3339))
}

// demoAuthenticationFlow demonstrates user authentication with MFA
func demoAuthenticationFlow() {
	fmt.Println("Scenario 1: User with MFA enabled\n")

	user := &exercise.User{
		ID:         "USER-001",
		Email:      "user@example.com",
		MFAEnabled: true,
		MFASecret:  "secret123",
	}

	sm := exercise.NewAuthStateMachine(user)
	ctx := context.Background()

	fmt.Printf("Initial state: %s\n\n", sm.Current())

	// Attempt login with MFA enabled
	fmt.Println("Step 1: User attempts login...")
	if err := sm.Transition(ctx, exercise.Event(exercise.EventLogin)); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("Current state: %s (MFA required)\n\n", sm.Current())

	// Provide invalid MFA code
	fmt.Println("Step 2: User provides invalid MFA code...")
	user.MFACode = "wrong_code"
	if err := sm.Transition(ctx, exercise.Event(exercise.EventMFASuccess)); err != nil {
		fmt.Printf("MFA failed as expected: %v\n", err)
		fmt.Printf("Current state: %s (still pending MFA)\n\n", sm.Current())
	}

	// Provide correct MFA code
	fmt.Println("Step 3: User provides correct MFA code...")
	user.MFACode = "secret123" // Simplified: code matches secret
	if err := sm.Transition(ctx, exercise.Event(exercise.EventMFASuccess)); err != nil {
		log.Fatalf("MFA verification failed: %v", err)
	}
	fmt.Printf("Current state: %s\n\n", sm.Current())

	// Logout
	fmt.Println("Step 4: User logs out...")
	if err := sm.Transition(ctx, exercise.Event(exercise.EventLogout)); err != nil {
		log.Fatalf("Logout failed: %v", err)
	}
	fmt.Printf("Current state: %s\n\n", sm.Current())

	// Scenario 2: User without MFA
	fmt.Println("\nScenario 2: User without MFA enabled\n")

	user2 := &exercise.User{
		ID:         "USER-002",
		Email:      "user2@example.com",
		MFAEnabled: false,
	}

	sm2 := exercise.NewAuthStateMachine(user2)

	fmt.Printf("Initial state: %s\n\n", sm2.Current())

	// Login directly without MFA
	fmt.Println("Step 1: User logs in (no MFA required)...")
	if err := sm2.Transition(ctx, exercise.Event(exercise.EventLogin)); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("Current state: %s (logged in directly)\n", sm2.Current())
}

// demoInvalidTransitions shows how guards prevent invalid state changes
func demoInvalidTransitions() {
	fmt.Println("Attempt 1: Pay for order with $0 amount\n")

	order := &exercise.Order{
		ID:            "ORDER-002",
		CustomerEmail: "customer@example.com",
		Amount:        0, // Invalid amount
		PaymentMethod: "credit_card",
	}

	sm := exercise.NewOrderStateMachine(order)
	ctx := context.Background()

	fmt.Printf("Order amount: $%.2f\n", order.Amount)
	fmt.Printf("Current state: %s\n\n", sm.Current())

	if err := sm.Transition(ctx, exercise.Event(exercise.EventPay)); err != nil {
		fmt.Printf("Payment rejected: %v\n", err)
		fmt.Printf("Current state: %s (unchanged)\n\n", sm.Current())
	}

	fmt.Println("Attempt 2: Pay with missing payment method\n")

	order.Amount = 99.99
	order.PaymentMethod = "" // Missing payment method

	fmt.Printf("Order amount: $%.2f\n", order.Amount)
	fmt.Printf("Payment method: '%s'\n", order.PaymentMethod)
	fmt.Printf("Current state: %s\n\n", sm.Current())

	if err := sm.Transition(ctx, exercise.Event(exercise.EventPay)); err != nil {
		fmt.Printf("Payment rejected: %v\n", err)
		fmt.Printf("Current state: %s (unchanged)\n\n", sm.Current())
	}

	fmt.Println("Attempt 3: Cancel already delivered order\n")

	order.PaymentMethod = "credit_card"

	// Process order to delivered state
	sm.Transition(ctx, exercise.Event(exercise.EventPay))
	sm.Transition(ctx, exercise.Event(exercise.EventShip))
	sm.Transition(ctx, exercise.Event(exercise.EventDeliver))

	fmt.Printf("Current state: %s\n\n", sm.Current())

	// Try to cancel
	if err := sm.Transition(ctx, exercise.Event(exercise.EventCancel)); err != nil {
		fmt.Printf("Cancellation rejected: %v\n", err)
		fmt.Printf("Current state: %s (cannot cancel delivered order)\n", sm.Current())
	}
}

// demoStateHistory shows complete history tracking
func demoStateHistory() {
	order := &exercise.Order{
		ID:            "ORDER-003",
		CustomerEmail: "customer@example.com",
		Amount:        149.99,
		PaymentMethod: "paypal",
	}

	sm := exercise.NewOrderStateMachine(order)
	ctx := context.Background()

	// Process order through lifecycle
	fmt.Println("Processing order through complete lifecycle...\n")

	sm.Transition(ctx, exercise.Event(exercise.EventPay))
	time.Sleep(100 * time.Millisecond) // Small delay to show distinct timestamps

	sm.Transition(ctx, exercise.Event(exercise.EventShip))
	time.Sleep(100 * time.Millisecond)

	sm.Transition(ctx, exercise.Event(exercise.EventDeliver))

	// Display history
	fmt.Println("Order State History:")
	fmt.Println("-------------------")

	history := sm.History()
	for i, entry := range history {
		fmt.Printf("%d. %s -> %s (event: %s) at %s\n",
			i+1,
			entry.From,
			entry.To,
			entry.Event,
			entry.Timestamp.Format("15:04:05.000"),
		)
	}

	fmt.Printf("\nTotal transitions: %d\n", len(history))
	fmt.Printf("Final state: %s\n", sm.Current())
}

// demoConcurrentStateMachines demonstrates multiple orders being processed
func demoConcurrentStateMachines() {
	ctx := context.Background()

	// Create multiple orders
	orders := []*exercise.Order{
		{
			ID:            "ORDER-101",
			CustomerEmail: "customer1@example.com",
			Amount:        29.99,
			PaymentMethod: "credit_card",
		},
		{
			ID:            "ORDER-102",
			CustomerEmail: "customer2@example.com",
			Amount:        59.99,
			PaymentMethod: "paypal",
		},
		{
			ID:            "ORDER-103",
			CustomerEmail: "customer3@example.com",
			Amount:        99.99,
			PaymentMethod: "credit_card",
		},
	}

	// Create state machines for each order
	stateMachines := make([]*exercise.StateMachine, len(orders))
	for i, order := range orders {
		stateMachines[i] = exercise.NewOrderStateMachine(order)
	}

	fmt.Println("Processing multiple orders concurrently...\n")

	// Process orders with different lifecycles
	processOrder := func(sm *exercise.StateMachine, order *exercise.Order, scenario string) {
		fmt.Printf("Order %s (%s): Starting\n", order.ID, scenario)

		switch scenario {
		case "complete":
			sm.Transition(ctx, exercise.Event(exercise.EventPay))
			sm.Transition(ctx, exercise.Event(exercise.EventShip))
			sm.Transition(ctx, exercise.Event(exercise.EventDeliver))
			fmt.Printf("Order %s: %s -> Delivered\n", order.ID, exercise.OrderPending)

		case "cancelled":
			sm.Transition(ctx, exercise.Event(exercise.EventCancel))
			fmt.Printf("Order %s: %s -> Cancelled\n", order.ID, exercise.OrderPending)

		case "in_progress":
			sm.Transition(ctx, exercise.Event(exercise.EventPay))
			sm.Transition(ctx, exercise.Event(exercise.EventShip))
			fmt.Printf("Order %s: %s -> Shipped\n", order.ID, exercise.OrderPending)
		}
	}

	// Process orders with different scenarios
	processOrder(stateMachines[0], orders[0], "complete")
	processOrder(stateMachines[1], orders[1], "cancelled")
	processOrder(stateMachines[2], orders[2], "in_progress")

	// Display final states
	fmt.Println("\nFinal States:")
	fmt.Println("------------")
	for i, sm := range stateMachines {
		fmt.Printf("Order %s: %s\n", orders[i].ID, sm.Current())
	}
}

// Additional demonstration functions

// demoCanCheck shows how to check if transition is possible
func demoCanCheck() {
	fmt.Println("--- Demo: Checking Possible Transitions ---\n")

	order := &exercise.Order{
		ID:            "ORDER-004",
		CustomerEmail: "customer@example.com",
		Amount:        99.99,
		PaymentMethod: "credit_card",
	}

	sm := exercise.NewOrderStateMachine(order)

	fmt.Printf("Current state: %s\n\n", sm.Current())

	events := []exercise.OrderEvent{
		exercise.EventPay,
		exercise.EventShip,
		exercise.EventDeliver,
		exercise.EventCancel,
	}

	fmt.Println("Checking which transitions are possible:")
	for _, event := range events {
		canTransition := sm.Can(exercise.Event(event))
		status := "✗"
		if canTransition {
			status = "✓"
		}
		fmt.Printf("  %s Event '%s'\n", status, event)
	}

	// Now transition to paid state
	fmt.Println("\nTransitioning to 'paid' state...")
	sm.Transition(context.Background(), exercise.Event(exercise.EventPay))
	fmt.Printf("Current state: %s\n\n", sm.Current())

	fmt.Println("Checking which transitions are now possible:")
	for _, event := range events {
		canTransition := sm.Can(exercise.Event(event))
		status := "✗"
		if canTransition {
			status = "✓"
		}
		fmt.Printf("  %s Event '%s'\n", status, event)
	}
}

// demoActionsDemo shows entry/exit actions in detail
func demoActionsDemo() {
	fmt.Println("--- Demo: Entry and Exit Actions ---\n")

	order := &exercise.Order{
		ID:            "ORDER-005",
		CustomerEmail: "customer@example.com",
		Amount:        199.99,
		PaymentMethod: "credit_card",
	}

	sm := exercise.NewOrderStateMachine(order)
	ctx := context.Background()

	fmt.Println("Watch for entry and exit actions as we transition...\n")

	events := []exercise.OrderEvent{
		exercise.EventPay,
		exercise.EventShip,
		exercise.EventDeliver,
	}

	for _, event := range events {
		fmt.Printf("Triggering event: %s\n", event)
		if err := sm.Transition(ctx, exercise.Event(event)); err != nil {
			log.Printf("Error: %v\n", err)
		}
		fmt.Printf("New state: %s\n\n", sm.Current())
		time.Sleep(500 * time.Millisecond)
	}
}

// separator returns a visual separator
func separator() string {
	return "========================================================"
}

// Visual state diagram printer
func printStateDiagram() {
	fmt.Println("\nOrder State Machine Diagram:")
	fmt.Println("============================\n")
	fmt.Println("                 pay")
	fmt.Println("    ┌──────────────────────────────────┐")
	fmt.Println("    │                                  ▼")
	fmt.Println("┌─────────┐                        ┌──────┐   ship         ┌──────────┐")
	fmt.Println("│ Pending │                        │ Paid │─────────────→ │ Shipped  │")
	fmt.Println("└─────────┘                        └──────┘                └──────────┘")
	fmt.Println("    │                                                           │")
	fmt.Println("    │ cancel                                                    │ deliver")
	fmt.Println("    ▼                                                           ▼")
	fmt.Println("┌───────────┐                                               ┌────────────┐")
	fmt.Println("│ Cancelled │                                               │ Delivered  │")
	fmt.Println("└───────────┘                                               └────────────┘")
	fmt.Println("\nAuthentication State Machine Diagram:")
	fmt.Println("=====================================\n")
	fmt.Println("                login (no MFA)")
	fmt.Println("    ┌────────────────────────────────┐")
	fmt.Println("    │                                ▼")
	fmt.Println("┌─────────────┐                  ┌────────────┐")
	fmt.Println("│ Logged Out  │                  │ Logged In  │")
	fmt.Println("└─────────────┘                  └────────────┘")
	fmt.Println("    │                                ▲")
	fmt.Println("    │ login (MFA enabled)            │ logout")
	fmt.Println("    ▼                                │")
	fmt.Println("┌──────────────┐  mfa_success       │")
	fmt.Println("│ MFA Pending  │────────────────────┘")
	fmt.Println("└──────────────┘")
	fmt.Println()
}
