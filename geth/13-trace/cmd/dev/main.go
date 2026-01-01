package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/example/go-10x-minis/geth/13-trace/internal/trace"
)

type rpcTraceClient struct {
	rpc *rpc.Client
}

func (c rpcTraceClient) TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error) {
	var out json.RawMessage
	err := c.rpc.CallContext(ctx, &out, "debug_traceTransaction", txHash)
	return out, err
}

func main() {
	// BREAKPOINT: deterministic inputs (requires debug-enabled endpoint)
	rpcURL := "https://eth.llamarpc.com"
	txHash := common.Hash{} // set me

	if txHash == (common.Hash{}) {
		fmt.Println("Set txHash in cmd/dev to a real tx hash (and use a debug-enabled node).")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer rpcClient.Close()

	out, err := trace.Run(ctx, rpcTraceClient{rpc: rpcClient}, trace.Config{TxHash: txHash})
	if err != nil {
		panic(err)
	}

	fmt.Println("TraceBytes:", len(out.Trace))
}
