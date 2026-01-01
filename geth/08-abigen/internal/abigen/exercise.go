//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package abigen

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"context"
)

const erc20ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
// TODO: implement Run.
func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement callString.
func callString(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (string, error) {
	panic("TODO: implement")
}
// TODO: implement callUint8.
func callUint8(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (uint8, error) {
	panic("TODO: implement")
}
// TODO: implement callUint256.
func callUint256(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (*big.Int, error) {
	panic("TODO: implement")
}
