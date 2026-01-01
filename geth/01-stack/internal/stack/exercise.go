//go:build !solution && !reference

package stack

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
)

func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Retrieve Chain ID
	// TODO: Retrieve Network ID
	// TODO: Retrieve Block Header
	// TODO: Construct Result with Defensive Copying
	panic("not implemented")
}
