package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/25-toolbox/internal/toolbox"
)

type toolboxClient struct {
	*ethclient.Client
}

func (c toolboxClient) TransactionByHash(ctx context.Context, hash string) (*types.Transaction, bool, error) {
	return c.Client.TransactionByHash(ctx, common.HexToHash(hash))
}

func main() {
	// geth/25-toolbox
	//
	// Usage:
	//   go run ./geth/25-toolbox/cmd/app <RPC_URL> <command> [args...]
	//
	// Commands:
	//   status
	//   block <block_number>
	//   tx <tx_hash>
	//
	// BREAKPOINT: parse args
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: <RPC_URL> <command> [args...]")
		os.Exit(2)
	}

	rpcURL := os.Args[1]
	cmd := os.Args[2]
	args := []string{}
	if len(os.Args) > 3 {
		args = os.Args[3:]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// BREAKPOINT: run
	out, err := toolbox.Run(ctx, toolboxClient{Client: client}, toolbox.Config{Command: cmd, Args: args})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(out.Output, "", "  ")
	fmt.Println(string(b))
}
