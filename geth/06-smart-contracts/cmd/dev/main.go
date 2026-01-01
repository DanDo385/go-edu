package main

import (
	"fmt"
)

/*
Debug harness for geth/06-smart-contracts module.

This module is primarily console-based, so this Go debug harness is minimal.
The real learning happens in the Geth JavaScript console following the README.md tutorial.

This file demonstrates how the console concepts map to Go code that you'll learn
in modules 07-09.
*/

func main() {
	fmt.Println("geth/06-smart-contracts: Debug Harness")
	fmt.Println()
	fmt.Println("This module is primarily a console-based tutorial.")
	fmt.Println("The debug harness demonstrates concepts that will be covered in Go in later modules.")
	fmt.Println()
	fmt.Println("Key Concepts (learned via console in this module):")
	fmt.Println("  1. Call vs Transaction distinction")
	fmt.Println("  2. Contract addresses and ABIs")
	fmt.Println("  3. Creating contract objects")
	fmt.Println("  4. Making read-only calls")
	fmt.Println("  5. Sending state-changing transactions")
	fmt.Println("  6. Decoding events and logs")
	fmt.Println()
	fmt.Println("To learn these concepts:")
	fmt.Println("  1. Read README.md for comprehensive console tutorial")
	fmt.Println("  2. Start Geth: geth --dev --http --http.api eth,net,web3,personal")
	fmt.Println("  3. Open console: geth attach")
	fmt.Println("  4. Follow exercises in README.md")
	fmt.Println()
	fmt.Println("Module 07 (eth-call) will teach you how to do the same things in Go.")

	// BREAKPOINT: Set here to inspect this debug harness
	// This demonstrates the structure, but the real learning is in the console
}
