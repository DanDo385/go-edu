//go:build reference

package smartcontracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	errNilClient     = errors.New("nil contract caller")
	errShortABIValue = errors.New("abi value too short")
)

/*
Reference Solution - Contract Call via ABI Encoding
===================================================

This file demonstrates read-only contract calls (eth_call) to query ERC20
metadata: name, symbol, decimals, totalSupply. No gas is spent, no transaction
is signed — we simulate execution at a block and read the return value.

This connects to the Ethereum ecosystem by showing:
- Keccak256 of function signature: first 4 bytes = selector ("name()", etc.)
- CallContract: simulates execution, returns hex-encoded bytes
- ABI encoding: dynamic types (string) use offset + length layout
- cfg.BlockNumber: nil = latest, or specific block for historical queries

The exercise builds understanding of:
- Function selectors: deterministic 4-byte identifier for method dispatch
- ABI layout: first 32 bytes = offset to dynamic data; at offset, 32 bytes =
  length, then raw bytes
- &cfg.Contract: CallMsg.To expects *common.Address; nil To = contract creation
- Read-only vs state-changing: eth_call never persists changes

Teaching notes (per .cursorrules):
- Pointer semantics: cfg.BlockNumber can be nil (latest). &cfg.Contract passes
  address of the config's Contract — CallMsg expects a pointer.
- Memory: crypto.Keccak256 returns a slice; we copy first 4 bytes via
  append([]byte(nil), h[:4]...) so callers can't mutate our selector.
- Slice bounds: data[:32], data[off:off+32] — ABI uses 32-byte words.
*/

/*
Run - Query ERC20 Metadata via eth_call

Parameters:
  - ctx: cancellation and timeout for RPC
  - client: ContractCaller (CallContract)
  - cfg: contract address and optional block number

Returns *Result with Name, Symbol, Decimals, TotalSupply; error on RPC or ABI decode failure.

Algorithm:
  1. Validate client and contract address
  2. For each selector (name, symbol, decimals, totalSupply): CallContract, decode
  3. decodeString: ABI dynamic encoding — offset at [0:32], length at [off:off+32], data at [off+32:...]
  4. decodeUint8/256: first 32 bytes as big-endian integer
*/
func Run(ctx context.Context, client ContractCaller, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}

	// call wraps CallContract: To = contract addr, Data = selector + encoded args (none here)
	call := func(sel []byte) ([]byte, error) {
		payload, err := client.CallContract(ctx, ethereum.CallMsg{To: &cfg.Contract, Data: sel}, cfg.BlockNumber)
		if err != nil {
			return nil, err
		}
		return payload, nil
	}

	nameRaw, err := call(selector("name()"))
	if err != nil {
		return nil, fmt.Errorf("call name: %w", err)
	}
	name, err := decodeString(nameRaw)
	if err != nil {
		return nil, fmt.Errorf("decode name: %w", err)
	}

	symbolRaw, err := call(selector("symbol()"))
	if err != nil {
		return nil, fmt.Errorf("call symbol: %w", err)
	}
	symbol, err := decodeString(symbolRaw)
	if err != nil {
		return nil, fmt.Errorf("decode symbol: %w", err)
	}

	decimalsRaw, err := call(selector("decimals()"))
	if err != nil {
		return nil, fmt.Errorf("call decimals: %w", err)
	}
	decimals, err := decodeUint8(decimalsRaw)
	if err != nil {
		return nil, fmt.Errorf("decode decimals: %w", err)
	}

	supplyRaw, err := call(selector("totalSupply()"))
	if err != nil {
		return nil, fmt.Errorf("call totalSupply: %w", err)
	}
	totalSupply, err := decodeUint256(supplyRaw)
	if err != nil {
		return nil, fmt.Errorf("decode totalSupply: %w", err)
	}

	return &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}, nil
}

/*
selector - Compute ABI Function Selector

Keccak256 of "name()" etc., first 4 bytes. The selector uniquely identifies
the function for contract dispatch. We copy via append([]byte(nil), h[:4]...)
so the returned slice is independent of Keccak256's internal buffer.
*/
func selector(sig string) []byte {
	h := crypto.Keccak256([]byte(sig))
	return append([]byte(nil), h[:4]...)
}

/*
decodeString - Decode ABI-Encoded Dynamic String

ABI layout for dynamic types: first 32 bytes = offset to data; at offset:
32 bytes = length, then raw bytes. We validate bounds at each step to avoid
panic on malformed or truncated responses.
*/
func decodeString(data []byte) (string, error) {
	if len(data) < 64 {
		return "", errShortABIValue
	}
	offset := new(big.Int).SetBytes(data[:32])
	if !offset.IsInt64() || offset.Int64() < 0 {
		return "", errors.New("invalid string offset")
	}
	off := int(offset.Int64())
	if off+32 > len(data) {
		return "", errors.New("string offset out of bounds")
	}

	length := new(big.Int).SetBytes(data[off : off+32])
	if !length.IsInt64() || length.Int64() < 0 {
		return "", errors.New("invalid string length")
	}
	n := int(length.Int64())
	start := off + 32
	end := start + n
	if end > len(data) {
		return "", errors.New("string length out of bounds")
	}

	return string(data[start:end]), nil
}

// decodeUint8 decodes first 32 bytes as uint256, then checks 0–255 range.
func decodeUint8(data []byte) (uint8, error) {
	v, err := decodeUint256(data)
	if err != nil {
		return 0, err
	}
	if !v.IsUint64() || v.Uint64() > 255 {
		return 0, errors.New("uint8 out of range")
	}
	return uint8(v.Uint64()), nil
}

// decodeUint256 decodes first 32 bytes as big-endian uint256; returns defensive copy.
func decodeUint256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errShortABIValue
	}
	return new(big.Int).SetBytes(data[:32]), nil
}
