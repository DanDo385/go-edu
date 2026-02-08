//go:build reference

package toolbox

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var errNilClient = errors.New("nil toolbox client")

type statusOutput struct {
	ChainID     *big.Int
	NetworkID   *big.Int
	LatestBlock uint64
	LatestHash  string
	IsSyncing   bool
	PeerCount   uint64
}

type blockOutput struct {
	Number   uint64
	Hash     string
	Parent   string
	TxCount  int
	GasUsed  uint64
	GasLimit uint64
}

type txOutput struct {
	Hash    string
	To      string
	Nonce   uint64
	Value   *big.Int
	Pending bool
	Gas     uint64
}

/*
Reference Solution

Structure:
- Dispatch by command (`status`, `block`, `tx`).
- Reuse focused handler functions for each operation.
- Return command-specific payloads through `Result.Output`.
*/
func Run(ctx context.Context, client ToolboxClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if strings.TrimSpace(cfg.Command) == "" {
		return nil, errors.New("command is required")
	}

	switch strings.ToLower(cfg.Command) {
	case "status":
		return handleStatus(ctx, client)
	case "block":
		return handleBlock(ctx, client, cfg.Args)
	case "tx":
		return handleTx(ctx, client, cfg.Args)
	default:
		return nil, fmt.Errorf("unknown command %q", cfg.Command)
	}
}

func handleStatus(ctx context.Context, client ToolboxClient) (*Result, error) {
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id: nil response")
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}
	if networkID == nil {
		return nil, errors.New("network id: nil response")
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("latest header: %w", err)
	}
	if header == nil {
		return nil, errors.New("latest header: nil response")
	}

	progress, err := client.SyncProgress(ctx)
	if err != nil {
		return nil, fmt.Errorf("sync progress: %w", err)
	}

	peers, err := client.PeerCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("peer count: %w", err)
	}

	out := statusOutput{
		ChainID:     new(big.Int).Set(chainID),
		NetworkID:   new(big.Int).Set(networkID),
		LatestBlock: header.Number.Uint64(),
		LatestHash:  header.Hash().Hex(),
		IsSyncing:   progress != nil,
		PeerCount:   peers,
	}

	return &Result{Command: "status", Output: out, Status: "ok"}, nil
}

func handleBlock(ctx context.Context, client ToolboxClient, args []string) (*Result, error) {
	if len(args) < 1 {
		return nil, errors.New("block command requires block number argument")
	}

	n, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse block number: %w", err)
	}

	block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(n))
	if err != nil {
		return nil, fmt.Errorf("block by number: %w", err)
	}
	if block == nil {
		return nil, errors.New("block by number: nil response")
	}

	out := blockOutput{
		Number:   block.NumberU64(),
		Hash:     block.Hash().Hex(),
		Parent:   block.ParentHash().Hex(),
		TxCount:  len(block.Transactions()),
		GasUsed:  block.GasUsed(),
		GasLimit: block.GasLimit(),
	}

	return &Result{Command: "block", Output: out, Status: "ok"}, nil
}

func handleTx(ctx context.Context, client ToolboxClient, args []string) (*Result, error) {
	if len(args) < 1 {
		return nil, errors.New("tx command requires transaction hash argument")
	}

	tx, pending, err := client.TransactionByHash(ctx, common.HexToHash(args[0]))
	if err != nil {
		return nil, fmt.Errorf("transaction by hash: %w", err)
	}
	if tx == nil {
		return nil, errors.New("transaction by hash: nil response")
	}

	to := ""
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	out := txOutput{
		Hash:    tx.Hash().Hex(),
		To:      to,
		Nonce:   tx.Nonce(),
		Value:   new(big.Int).Set(tx.Value()),
		Pending: pending,
		Gas:     tx.Gas(),
	}

	return &Result{Command: "tx", Output: out, Status: "ok"}, nil
}
