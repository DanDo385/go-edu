//go:build !solution && !reference

package trace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

func Run(ctx context.Context, client TraceClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Trace Transaction
	// TODO: Defensive Copy of JSON Data
	// TODO: Return Result
	panic("not implemented")
}
