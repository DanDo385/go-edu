package main

import (
	"fmt"
	"os"
)

/*
This module is primarily a tutorial-based learning experience using the Geth console.
The exercises are designed to be completed in the Geth JavaScript console, not in Go.

After completing this module, you'll understand:
- Call vs Transaction distinction
- Contract addresses and ABIs
- Console-based contract interaction
- Event decoding

Module 07 (eth-call) will teach you how to do the same things in Go.

This Go wrapper demonstrates the concepts learned in the console tutorial.
*/

func main() {
	if len(os.Args) < 2 {
		fmt.Println("geth/06-smart-contracts: Smart Contract Interaction Fundamentals")
		fmt.Println()
		fmt.Println("This module is primarily a console-based tutorial.")
		fmt.Println("Please see README.md for the comprehensive Geth console tutorial.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  This program demonstrates concepts from the console tutorial.")
		fmt.Println("  For the full learning experience, follow the README.md tutorial.")
		fmt.Println()
		fmt.Println("To start learning:")
		fmt.Println("  1. Read README.md for the step-by-step console tutorial")
		fmt.Println("  2. Start Geth: geth --dev --http --http.api eth,net,web3,personal")
		fmt.Println("  3. Open console: geth attach")
		fmt.Println("  4. Follow the exercises in README.md")
		return
	}

	// BREAKPOINT: Set here to inspect command-line arguments
	command := os.Args[1]

	switch command {
	case "concepts":
		fmt.Println("Key Concepts Covered in This Module:")
		fmt.Println("  1. Call vs Transaction (read-only vs state-changing)")
		fmt.Println("  2. Contract addresses and ABIs")
		fmt.Println("  3. Creating contract objects in Geth console")
		fmt.Println("  4. Making read-only calls")
		fmt.Println("  5. Sending state-changing transactions")
		fmt.Println("  6. Decoding events and logs")
		fmt.Println()
		fmt.Println("See README.md for detailed explanations and console examples.")

	case "next":
		fmt.Println("Next Steps:")
		fmt.Println("  Module 07 (eth-call): Doing the same interactions from Go")
		fmt.Println("  Module 08 (abigen): Understanding eth_call vs SendTransaction")
		fmt.Println("  Module 09 (events): Watching pending txs and simulating execution")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: concepts, next")
	}
}
