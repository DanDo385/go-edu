//go:build !solution && !reference

package abigen

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

func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Parse ABI JSON
	// TODO: Create BoundContract
	// TODO: Create CallOpts
	// TODO: Call name()
	// TODO: Call symbol()
	// TODO: Call decimals()
	// TODO: Call totalSupply()
	// TODO: Optionally Call balanceOf(address)
	// TODO: Construct and Return Result
	panic("not implemented")
}

func callString(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (string, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func callUint8(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (uint8, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func callUint256(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (*big.Int, error) {
	// TODO: Implement this function
	panic("not implemented")
}
