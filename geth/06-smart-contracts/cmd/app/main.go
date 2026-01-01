package main

import (
	"fmt"
	"os"
)

/*
geth/06-smart-contracts: cmd/app

This module is primarily a console tutorial teaching smart contract interaction concepts.
The main learning happens in the Geth JavaScript console (see README.md).

This Go program serves as an optional reference showing how console concepts
map to Go code. For the actual learning experience, follow the tutorial in README.md.

Usage:
  go run ./cmd/app/main.go

This will print instructions for using the Geth console to interact with smart contracts.
For actual contract interaction in Go, see modules 07 (eth-call), 08 (abigen), and 09 (events).
*/

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         geth/06-smart-contracts: Console Tutorial              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("This module teaches smart contract interaction using the Geth console.")
	fmt.Println("The main learning happens in the JavaScript console, not Go code.")
	fmt.Println()
	fmt.Println("📖 To start the tutorial:")
	fmt.Println("   1. Read README.md in this directory")
	fmt.Println("   2. Start Geth: geth --dev --http --http.api eth,net,web3,personal")
	fmt.Println("   3. Attach console: geth attach")
	fmt.Println("   4. Follow the step-by-step tutorial in README.md")
	fmt.Println()
	fmt.Println("🎯 What you'll learn:")
	fmt.Println("   • Call vs Transaction (read-only vs state-changing)")
	fmt.Println("   • Contract addresses and ABIs")
	fmt.Println("   • Making calls and sending transactions")
	fmt.Println("   • Decoding events and logs")
	fmt.Println("   • Common pitfalls and debugging")
	fmt.Println()
	fmt.Println("🚀 After completing this module:")
	fmt.Println("   • geth/07-eth-call: Same concepts in Go")
	fmt.Println("   • geth/08-abigen: Type-safe Go bindings")
	fmt.Println("   • geth/09-events: Event subscriptions in Go")
	fmt.Println()
	fmt.Println("For Go code examples of these concepts, see modules 07-09.")
	fmt.Println()

	// BREAKPOINT: Set a breakpoint here if you want to step through the instructions
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("Additional resources:")
		fmt.Println("  • Geth Console Docs: https://geth.ethereum.org/docs/interacting-with-geth/javascript-console")
		fmt.Println("  • ABI Specification: https://docs.soliditylang.org/en/latest/abi-spec.html")
		fmt.Println("  • ERC20 Standard: https://eips.ethereum.org/EIPS/eip-20")
	}
}
