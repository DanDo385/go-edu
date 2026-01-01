package smartcontracts

import (
	"testing"
)

/*
geth/06-smart-contracts: Tests

This module is primarily console-tutorial based.
Tests here are minimal since the main learning happens in the Geth JavaScript console.

To complete this module:
1. Follow the tutorial in README.md
2. Complete exercises in the Geth console
3. Use geth/07-eth-call for Go-based contract interaction testing
*/

func TestModuleUnderstanding(t *testing.T) {
	t.Log("This module teaches smart contract interaction via Geth console")
	t.Log("Key concepts:")
	t.Log("  • Call vs Transaction (read-only vs state-changing)")
	t.Log("  • Contract addresses and ABIs")
	t.Log("  • Making calls: eth_call (free, instant, read-only)")
	t.Log("  • Sending transactions: eth_sendTransaction (signed, costs gas, state-changing)")
	t.Log("  • Event decoding: Understanding logs and topics")
	t.Log("")
	t.Log("To complete exercises:")
	t.Log("  1. Read README.md")
	t.Log("  2. Start Geth: geth --dev --http --http.api eth,net,web3,personal")
	t.Log("  3. Attach: geth attach")
	t.Log("  4. Follow step-by-step console tutorial")
	t.Log("")
	t.Log("Next modules:")
	t.Log("  • geth/07-eth-call: Go implementation of eth_call")
	t.Log("  • geth/08-abigen: Type-safe Go contract bindings")
	t.Log("  • geth/09-events: Event subscriptions in Go")
}

func TestConceptualUnderstanding(t *testing.T) {
	tests := []struct {
		name     string
		concept  string
		question string
		answer   string
	}{
		{
			name:     "Call vs Transaction",
			concept:  "Read-only operations",
			question: "What method do you use for view/pure functions?",
			answer:   "eth_call - free, instant, no signature required",
		},
		{
			name:     "Call vs Transaction",
			concept:  "State-changing operations",
			question: "What method do you use for non-view functions?",
			answer:   "eth_sendTransaction - signed, costs gas, creates receipt",
		},
		{
			name:     "ABI",
			concept:  "Contract interface",
			question: "Why do you need an ABI?",
			answer:   "To know what functions exist, their parameters, and return types",
		},
		{
			name:     "Events",
			concept:  "Contract communication",
			question: "How do contracts return data from transactions?",
			answer:   "Via events in the transaction receipt logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Concept: %s", tt.concept)
			t.Logf("Question: %s", tt.question)
			t.Logf("Answer: %s", tt.answer)
		})
	}
}

func TestPrerequisiteCheck(t *testing.T) {
	t.Log("Prerequisites for this module:")
	t.Log("  ✓ Completion of geth/01-stack through geth/05-tx-nonces")
	t.Log("  ✓ Understanding of Ethereum basics (addresses, transactions, gas)")
	t.Log("  ✓ Geth installed: run 'geth version' to verify")
	t.Log("")
	t.Log("If prerequisites are met, proceed to README.md tutorial")
}
