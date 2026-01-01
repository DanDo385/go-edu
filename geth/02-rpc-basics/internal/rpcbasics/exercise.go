//go:build !solution && !reference

package rpcbasics

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum/core/types"
)

func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Retrieve Latest Block Number
	// TODO: Retrieve Network ID
	// TODO: Retrieve Full Block with Retry Logic
	// TODO: Construct Result with Defensive Copying
	panic("not implemented")
}
