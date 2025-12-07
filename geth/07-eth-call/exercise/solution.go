//go:build solution
// +build solution

package exercise

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
	selectorName        = selector("name()")
	selectorSymbol      = selector("symbol()")
	selectorDecimals    = selector("decimals()")
	selectorTotalSupply = selector("totalSupply()")
)

// Run contains the reference solution for module 07-eth-call.
func Run(ctx context.Context, client CallClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address required")
	}

	call := func(selector []byte) ([]byte, error) {
		msg := ethereum.CallMsg{
			To:   &cfg.Contract,
			Data: selector,
		}
		return client.CallContract(ctx, msg, cfg.BlockNumber)
	}

	nameBytes, err := call(selectorName)
	if err != nil {
		return nil, fmt.Errorf("call name(): %w", err)
	}
	name, err := decodeString(nameBytes)
	if err != nil {
		return nil, fmt.Errorf("decode name(): %w", err)
	}

	symbolBytes, err := call(selectorSymbol)
	if err != nil {
		return nil, fmt.Errorf("call symbol(): %w", err)
	}
	symbol, err := decodeString(symbolBytes)
	if err != nil {
		return nil, fmt.Errorf("decode symbol(): %w", err)
	}

	decBytes, err := call(selectorDecimals)
	if err != nil {
		return nil, fmt.Errorf("call decimals(): %w", err)
	}
	decimals, err := decodeUint8(decBytes)
	if err != nil {
		return nil, fmt.Errorf("decode decimals(): %w", err)
	}

	supplyBytes, err := call(selectorTotalSupply)
	if err != nil {
		return nil, fmt.Errorf("call totalSupply(): %w", err)
	}
	totalSupply, err := decodeUint256(supplyBytes)
	if err != nil {
		return nil, fmt.Errorf("decode totalSupply(): %w", err)
	}

	return &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}, nil
}

func selector(sig string) []byte {
	hash := crypto.Keccak256([]byte(sig))
	return hash[:4]
}

func decodeString(data []byte) (string, error) {
	if len(data) < 64 {
		return "", errors.New("data too short for string")
	}
	offset := new(big.Int).SetBytes(data[:32]).Int64()
	if offset < 0 || offset+32 > int64(len(data)) {
		return "", errors.New("invalid offset")
	}
	lengthStart := int(offset)
	lengthEnd := lengthStart + 32
	if lengthEnd > len(data) {
		return "", errors.New("invalid length data")
	}
	length := new(big.Int).SetBytes(data[lengthStart:lengthEnd]).Int64()
	if length < 0 {
		return "", errors.New("negative length")
	}
	dataStart := lengthEnd
	dataEnd := dataStart + int(length)
	if dataEnd > len(data) {
		return "", errors.New("string exceeds data bounds")
	}
	return string(data[dataStart:dataEnd]), nil
}

func decodeUint8(data []byte) (uint8, error) {
	if len(data) < 32 {
		return 0, errors.New("data too short for uint8")
	}
	return uint8(data[len(data)-1]), nil
}

func decodeUint256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errors.New("data too short for uint256")
	}
	out := new(big.Int).SetBytes(data[len(data)-32:])
	return out, nil
}
