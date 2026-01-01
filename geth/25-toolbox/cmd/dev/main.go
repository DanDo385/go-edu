package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	// BREAKPOINT
	rpcURL := "https://eth.llamarpc.com"
	cfg := toolbox.Config{Command: "status"}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := toolbox.Run(ctx, toolboxClient{Client: client}, cfg)
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(out.Output, "", "  ")
	fmt.Println(string(b))
}
