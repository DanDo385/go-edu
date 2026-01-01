package main

import (
	"fmt"
)

/*
geth/06-smart-contracts: cmd/dev (Debug Harness)

This module is primarily a console tutorial teaching smart contract interaction concepts.
The main learning happens in the Geth JavaScript console (see README.md).

This debug harness serves as a reference point, but the actual exercises
are meant to be completed in the Geth console.

To use this module properly:
1. Read README.md
2. Follow the console tutorial
3. Use modules 07-09 for Go implementations

BREAKPOINT: This file is minimal because the learning happens in the Geth console.
*/

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  geth/06-smart-contracts: Debug Harness")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// BREAKPOINT: Set a breakpoint here to understand the module's purpose
	fmt.Println("⚠️  This module is console-tutorial based.")
	fmt.Println()
	fmt.Println("The exercises are completed in the Geth JavaScript console,")
	fmt.Println("not in Go code. This teaches you the foundational concepts")
	fmt.Println("before moving to Go implementations.")
	fmt.Println()
	fmt.Println("Steps:")
	fmt.Println("  1. Start Geth dev chain:")
	fmt.Println("     geth --dev --http --http.api eth,net,web3,personal")
	fmt.Println()
	fmt.Println("  2. Attach to console:")
	fmt.Println("     geth attach")
	fmt.Println()
	fmt.Println("  3. Follow README.md tutorial:")
	fmt.Println("     - Load ABI")
	fmt.Println("     - Create contract object")
	fmt.Println("     - Make read-only calls")
	fmt.Println("     - Send state-changing transactions")
	fmt.Println("     - Decode events")
	fmt.Println()
	fmt.Println("  4. After mastering console concepts, proceed to:")
	fmt.Println("     • geth/07-eth-call (Go implementation)")
	fmt.Println("     • geth/08-abigen (Type-safe bindings)")
	fmt.Println("     • geth/09-events (Event subscriptions)")
	fmt.Println()

	// BREAKPOINT: Demonstrate the concept of console-first learning
	conceptualFlow := []string{
		"Console → Understand concepts visually",
		"Go      → Implement programmatically",
		"Scale   → Build production systems",
	}

	fmt.Println("Learning flow:")
	for i, step := range conceptualFlow {
		fmt.Printf("  %d. %s\n", i+1, step)
	}
	fmt.Println()
	fmt.Println("See README.md for the complete tutorial.")
}
