package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/example/go-10x-minis/geth/13-trace/internal/trace"
)

func rpcURLFromArgs(defaultURL string) string {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		return os.Args[1]
	}
	if v := os.Getenv("RPC_URL"); v != "" {
		return v
	}
	return defaultURL
}

type rpcTraceClient struct {
	rpc *rpc.Client
}

func (c rpcTraceClient) TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error) {
	var out json.RawMessage
	err := c.rpc.CallContext(ctx, &out, "debug_traceTransaction", txHash)
	return out, err
}

func main() {
	// geth/13-trace
	//
	// NOTE: Many public RPC providers disable debug tracing.
	// This program is best-effort: it will explain failures clearly.
	//
	// Usage:
	//   go run ./geth/13-trace/cmd/app <RPC_URL> [tx_hash]
	//
	// BREAKPOINT: parse args
	rpcURL := rpcURLFromArgs("https://eth.llamarpc.com")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var txHash common.Hash
	if len(os.Args) >= 3 {
		txHash = common.HexToHash(os.Args[2])
	} else {
		// Pick a recent tx hash (first tx in latest block), if available.
		c, err := ethclient.DialContext(ctx, rpcURL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dial:", err)
			os.Exit(1)
		}
		defer c.Close()

		b, err := c.BlockByNumber(ctx, nil)
		if err == nil && b != nil && len(b.Transactions()) > 0 {
			txHash = b.Transactions()[0].Hash()
		} else {
			fmt.Fprintln(os.Stderr, "provide a tx_hash: could not auto-select from latest block")
			os.Exit(2)
		}
	}

	// BREAKPOINT: dial raw RPC for debug_* methods
	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial rpc:", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	out, err := trace.Run(ctx, rpcTraceClient{rpc: rpcClient}, trace.Config{TxHash: txHash})
	if err != nil {
		fmt.Fprintln(os.Stderr, "trace failed:", err)
		fmt.Fprintln(os.Stderr, "Hint: use a local geth node with --http.api debug,eth,net,web3 (and often an archive node for older txs).")
		os.Exit(1)
	}

	fmt.Println("TxHash:", out.TxHash.Hex())
	fmt.Printf("TraceBytes: %d\n", len(out.Trace))
}
