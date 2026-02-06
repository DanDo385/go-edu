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
Reference Solution

Structure:
- Build four read-only call payloads using function selectors.
- Execute eth_call for each selector.
- Decode ABI-encoded return values.

Invariant:
- Contract must be non-zero address and client must be non-nil.
*/
func Run(ctx context.Context, client ContractCaller, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}

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

func selector(sig string) []byte {
	h := crypto.Keccak256([]byte(sig))
	return append([]byte(nil), h[:4]...)
}

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

func decodeUint256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errShortABIValue
	}
	return new(big.Int).SetBytes(data[:32]), nil
}
