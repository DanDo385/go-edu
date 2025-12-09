//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const erc20ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

/*
Problem: Use BoundContract for type-safe contract calls with automatic ABI encoding/decoding.

This module teaches you how to use go-ethereum's BoundContract pattern for cleaner,
safer contract interactions. Instead of manually encoding/decoding like module 07,
you'll use the abi package to handle encoding/decoding automatically.

Computer science principles highlighted:
  - Adapter pattern: BoundContract wraps low-level RPC with high-level interface
  - Type safety: ABI definitions provide compile-time checks
  - Code reuse: Helper functions eliminate boilerplate
  - Separation of concerns: ABI encoding is separate from business logic
*/
func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Same Pattern as Module 07
	// ============================================================================
	// Why validate inputs? This is a library function that will be called by other
	// code. We can't trust callers to always pass valid inputs.
	//
	// Context handling: This pattern repeats from modules 01, 06, and 07. If ctx
	// is nil, we provide context.Background() as a safe default. By module 08, this
	// should be muscle memory—always check for nil context and provide default.
	//
	// Building on previous concepts: Every module starts with this validation pattern.
	// The repetition is intentional—it reinforces defensive programming habits.
	if ctx == nil {
		ctx = context.Background()
	}

	// Backend validation: ContractCaller is the interface for making contract calls.
	// It's a subset of ethclient.Client that only includes CallContract method.
	// This is interface segregation principle—depend only on what you need.
	if backend == nil {
		return nil, errors.New("backend is nil")
	}

	// Contract address validation: We need a valid contract address to call.
	// Zero address is invalid for contract interactions.
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address required")
	}

	// ============================================================================
	// STEP 2: Parse ABI JSON - Understanding Contract Interface
	// ============================================================================
	// ABI (Application Binary Interface) is like an interface definition for contracts.
	// It tells us:
	//   - What functions exist
	//   - What parameters they take
	//   - What types they return
	//   - Whether they're view/pure (read-only) or state-changing
	//
	// Why do we need ABI? Without it, we can't encode function calls or decode
	// return values. The ABI is the "contract" between our Go code and the smart
	// contract.
	//
	// Default ABI: We provide erc20ABI as a default so users don't have to specify
	// it for standard ERC20 tokens. This is the "sensible defaults" pattern—make
	// common cases easy, advanced cases possible.
	abiJSON := cfg.ABI
	if strings.TrimSpace(abiJSON) == "" {
		abiJSON = erc20ABI
	}

	// Parse ABI: The abi.JSON function parses the JSON string into an abi.ABI object.
	// This object contains method definitions, event definitions, and type information
	// needed for encoding/decoding.
	//
	// Error handling: ABI parsing can fail if:
	//   - JSON is malformed
	//   - Function signatures are invalid
	//   - Type definitions are incorrect
	// We wrap the error with context for debugging.
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("parse ABI: %w", err)
	}

	// ============================================================================
	// STEP 3: Create BoundContract - The Adapter Pattern
	// ============================================================================
	// BoundContract is the adapter pattern in action. It adapts the low-level
	// RPC interface (CallContract) to a high-level interface that understands
	// ABI encoding/decoding.
	//
	// What BoundContract provides:
	//   - Call(): For view/pure functions (read operations)
	//   - Transact(): For state-changing functions (write operations)
	//   - FilterLogs(): For querying events
	//
	// Parameters explained:
	//   - cfg.Contract: The contract address we want to interact with
	//   - parsedABI: The ABI definition (tells us how to encode/decode)
	//   - backend: The RPC client for read operations (ContractCaller)
	//   - nil: Transactor for write operations (not needed here—we're only reading)
	//   - nil: Filterer for event queries (not needed in this module)
	//
	// Why three separate interfaces? Separation of concerns. A read-only client
	// doesn't need transaction signing capabilities. This follows the interface
	// segregation principle.
	//
	// Compare to module 07: In module 07, we manually built CallMsg and decoded
	// responses. BoundContract does all of this for us automatically based on the ABI.
	contract := bind.NewBoundContract(cfg.Contract, parsedABI, backend, nil, nil)

	// ============================================================================
	// STEP 4: Create CallOpts - Configuring Contract Calls
	// ============================================================================
	// CallOpts is like a request context for contract calls. It contains:
	//   - Context: For cancellation and timeouts
	//   - BlockNumber: Which block to query (nil = latest, number = historical)
	//   - From: Optional sender address (for view functions that check msg.sender)
	//
	// Why separate CallOpts from contract? Because different calls might need
	// different options. For example, you might query the same contract at different
	// block numbers. This is separation of concerns—call configuration is separate
	// from contract definition.
	callOpts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: cfg.BlockNumber,
	}

	// From field: Some view functions check msg.sender even though they don't
	// modify state. For example, allowance(owner, spender) might check msg.sender
	// for access control. We only set From if the caller provided a holder address.
	//
	// Why pointer? cfg.Holder is *common.Address to distinguish between:
	//   - nil: No holder provided (don't set From)
	//   - &address{}: Zero address provided (set From to zero address)
	// This is the "nil means absent" pattern in Go.
	if cfg.Holder != nil {
		callOpts.From = *cfg.Holder
	}

	// ============================================================================
	// STEP 5: Call name() - First BoundContract Usage
	// ============================================================================
	// name() is a view function that returns a string.
	//
	// Compare to module 07:
	//   Module 07: Manually compute selector, build CallMsg, decode string
	//   Module 08: Just call callString(contract, opts, "name")
	//
	// This is the power of abstraction. BoundContract + helpers hide all the
	// encoding/decoding complexity. The code is simpler, safer, and easier to maintain.
	//
	// The callString helper (implemented below) does:
	//   1. Calls contract.Call() with method name "name"
	//   2. Gets back []interface{} of return values
	//   3. Converts first return value to string type
	//   4. Returns the string or error
	//
	// Type safety: The helper does runtime type conversion. In a production system,
	// you'd use abigen to generate typed bindings with compile-time type safety.
	// This module shows you how it works under the hood.
	name, err := callString(contract, callOpts, "name")
	if err != nil {
		return nil, err
	}

	// ============================================================================
	// STEP 6: Call symbol() - Pattern Repetition
	// ============================================================================
	// symbol() follows the exact same pattern as name(). This demonstrates how
	// consistent patterns make code predictable and easy to understand.
	//
	// Notice: Same helper (callString), different method name. The helper is
	// reusable because all string-returning functions have the same structure.
	// This is the benefit of abstraction and code reuse.
	//
	// Building on previous modules: You've seen this "call → check error" pattern
	// in every module (01-stack, 06-eip1559, 07-eth-call). By module 08, it should
	// feel automatic.
	symbol, err := callString(contract, callOpts, "symbol")
	if err != nil {
		return nil, err
	}

	// ============================================================================
	// STEP 7: Call decimals() - Different Return Type, Same Pattern
	// ============================================================================
	// decimals() returns uint8 instead of string.
	//
	// Compare to name/symbol: Different return type, so we use callUint8 instead
	// of callString. But the pattern is identical:
	//   1. Call helper with (contract, opts, method name)
	//   2. Check error
	//   3. Use the value
	//
	// This demonstrates polymorphism through function naming. We can't use Go
	// generics here (pre-Go 1.18 code), so we use separate functions for each
	// return type. The pattern is consistent even though the types differ.
	decimals, err := callUint8(contract, callOpts, "decimals")
	if err != nil {
		return nil, err
	}

	// ============================================================================
	// STEP 8: Call totalSupply() - Understanding big.Int Returns
	// ============================================================================
	// totalSupply() returns uint256, which maps to *big.Int in Go.
	//
	// Why *big.Int? Solidity uint256 can hold values up to 2^256-1, which is far
	// larger than Go's uint64 (max ~2^64). Go's math/big package provides arbitrary
	// precision integers that can handle Solidity's uint256 range.
	//
	// Real-world importance: Token supplies can be enormous. For example:
	//   - Shiba Inu: 1 quadrillion tokens (10^15)
	//   - Many tokens: Supply in the trillions
	// These values exceed uint64's max, so we must use *big.Int.
	//
	// Pattern repetition: Again, same pattern—call helper, check error, use value.
	// The only difference is the return type (callUint256 returns *big.Int).
	totalSupply, err := callUint256(contract, callOpts, "totalSupply")
	if err != nil {
		return nil, err
	}

	// ============================================================================
	// STEP 9: Optionally Call balanceOf(address) - Functions with Parameters
	// ============================================================================
	// balanceOf(address) is different from the previous functions because it takes
	// a parameter. This demonstrates how BoundContract handles function parameters.
	//
	// Conditional logic: We only call balanceOf if cfg.Holder is provided. This is
	// the "optional feature" pattern—some fields in Result are optional (nil if
	// not requested).
	//
	// Why optional? Not all callers need balance information. By making it optional,
	// we avoid unnecessary RPC calls and make the function more flexible.
	var balance *big.Int
	if cfg.Holder != nil {
		// Variadic parameters: The callUint256 helper uses variadic params
		// (params ...interface{}) to accept any number of function arguments.
		// Here we pass *cfg.Holder as the address parameter for balanceOf.
		//
		// Type handling: *cfg.Holder is common.Address, which BoundContract
		// automatically encodes according to the ABI (addresses are 20 bytes,
		// padded to 32 bytes for ABI encoding).
		//
		// Compare to module 07: In module 07, you'd manually encode the address
		// parameter. Here, BoundContract does it automatically based on the ABI.
		balance, err = callUint256(contract, callOpts, "balanceOf", *cfg.Holder)
		if err != nil {
			return nil, err
		}
	}

	// ============================================================================
	// STEP 10: Construct and Return Result - No Defensive Copying Needed
	// ============================================================================
	// Why no defensive copying? The values we're returning are already independent
	// copies created by our helper functions' type conversions.
	//
	// Ownership analysis:
	//   - name, symbol: Strings are immutable in Go (safe to share)
	//   - decimals: uint8 is a value type (copied, not shared)
	//   - totalSupply, balance: *big.Int pointers from type conversion (fresh allocations)
	//
	// Compare to module 01: In module 01, we needed defensive copying because
	// RPC client might return shared pointers. Here, our helpers create new values,
	// so they're already independent.
	//
	// This demonstrates understanding of ownership:
	//   - Know when to copy (shared mutable data)
	//   - Know when NOT to copy (already independent data)
	//
	// Building on previous concepts:
	//   - Module 01: Learned defensive copying for RPC client returns
	//   - Module 07: Learned when copying is NOT needed (decoder output)
	//   - Module 08: Reinforces the "decoder output is safe" pattern
	return &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Balance:     balance,
	}, nil
}

// callString is a helper that calls a contract method returning string.
//
// How it works:
//   1. Calls contract.Call() with method name and parameters
//   2. contract.Call encodes parameters using ABI
//   3. Executes eth_call via backend
//   4. Decodes return value using ABI
//   5. Returns []interface{} of return values
//   6. We convert first value to string type
//
// Why []interface{}? Solidity functions can return multiple values. The ABI
// decoder returns a slice of interface{} values, one per return value. We
// extract the first one and convert to string.
//
// Error handling: contract.Call can fail if:
//   - Network error (RPC call failed)
//   - Contract reverted (require/revert statement)
//   - Decoding failed (return value doesn't match ABI)
func callString(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (string, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return "", fmt.Errorf("call %s: %w", method, err)
	}
	if len(out) == 0 {
		return "", fmt.Errorf("call %s: empty result", method)
	}
	// abi.ConvertType does type conversion from interface{} to concrete type.
	// It handles the low-level type assertion and conversion, returning a pointer
	// to the requested type. We dereference to get the actual string value.
	return *abi.ConvertType(out[0], new(string)).(*string), nil
}

// callUint8 is a helper that calls a contract method returning uint8.
//
// Same pattern as callString, but converts to uint8 instead of string.
// See callString comments for detailed explanation of how this works.
func callUint8(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (uint8, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return 0, fmt.Errorf("call %s: %w", method, err)
	}
	if len(out) == 0 {
		return 0, fmt.Errorf("call %s: empty result", method)
	}
	return *abi.ConvertType(out[0], new(uint8)).(*uint8), nil
}

// callUint256 is a helper that calls a contract method returning uint256.
//
// Same pattern as callString/callUint8, but converts to *big.Int.
// uint256 in Solidity maps to *big.Int in Go because uint256 can represent
// values larger than Go's native integer types (up to 2^256-1).
//
// See callString comments for detailed explanation of how this works.
func callUint256(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (*big.Int, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return nil, fmt.Errorf("call %s: %w", method, err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("call %s: empty result", method)
	}
	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int), nil
}
