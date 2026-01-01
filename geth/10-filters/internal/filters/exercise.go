//go:build !solution && !reference

package filters

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func subscribeHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func pollHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	panic("not implemented")
}
