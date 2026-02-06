//go:build reference

package ethcall

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
	errNilClient     = errors.New("nil call client")
	errShortABIValue = errors.New("abi value too short")

	selectorName        = selector("name()")
	selectorSymbol      = selector("symbol()")
	selectorDecimals    = selector("decimals()")
	selectorTotalSupply = selector("totalSupply()")
)

/*
Reference Solution

Structure:
- Build low-level selectors manually.
- Use eth_call for read-only contract queries.
- Decode ABI return payloads explicitly.

Invariant:
- All calls must target the configured contract address.
*/
func Run(ctx context.Context, client CallClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}

	call := func(sel []byte) ([]byte, error) {
		out, err := client.CallContract(ctx, ethereum.CallMsg{To: &cfg.Contract, Data: sel}, cfg.BlockNumber)
		if err != nil {
			return nil, err
		}
		return out, nil
	}

	nameRaw, err := call(selectorName)
	if err != nil {
		return nil, fmt.Errorf("call name: %w", err)
	}
	name, err := decodeString(nameRaw)
	if err != nil {
		return nil, fmt.Errorf("decode name: %w", err)
	}

	symbolRaw, err := call(selectorSymbol)
	if err != nil {
		return nil, fmt.Errorf("call symbol: %w", err)
	}
	symbol, err := decodeString(symbolRaw)
	if err != nil {
		return nil, fmt.Errorf("decode symbol: %w", err)
	}

	decimalsRaw, err := call(selectorDecimals)
	if err != nil {
		return nil, fmt.Errorf("call decimals: %w", err)
	}
	decimals, err := decodeUint8(decimalsRaw)
	if err != nil {
		return nil, fmt.Errorf("decode decimals: %w", err)
	}

	totalSupplyRaw, err := call(selectorTotalSupply)
	if err != nil {
		return nil, fmt.Errorf("call totalSupply: %w", err)
	}
	totalSupply, err := decodeUint256(totalSupplyRaw)
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
