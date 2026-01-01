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

const defaultMaxHeads = 5
const defaultPollInterval = time.Second

// Run contains the reference solution for module 10-filters.
func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	panic("unimplemented")
}

func subscribeHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	panic("unimplemented")
}

func pollHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	panic("unimplemented")
}
