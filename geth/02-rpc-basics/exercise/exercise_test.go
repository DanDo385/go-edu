package exercise

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type mockClient struct {
	blockNumber    uint64
	blockNumberErr error
	networkID      *big.Int
	networkErr     error
	block          *types.Block
	blockErrs      []error
	blockCalls     int
}

func (m *mockClient) BlockNumber(ctx context.Context) (uint64, error) {
	return m.blockNumber, m.blockNumberErr
}

func (m *mockClient) NetworkID(ctx context.Context) (*big.Int, error) {
	return m.networkID, m.networkErr
}

func (m *mockClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	call := m.blockCalls
	m.blockCalls++
	if call < len(m.blockErrs) && m.blockErrs[call] != nil {
		return nil, m.blockErrs[call]
	}
	return m.block, nil
}

func TestRunSuccess(t *testing.T) {
	block := types.NewBlockWithHeader(&types.Header{
		Number:     big.NewInt(123),
		ParentHash: common.HexToHash("0xdeadbeef"),
		GasUsed:    42,
	})
	client := &mockClient{
		blockNumber: 123,
		networkID:   big.NewInt(1),
		block:       block,
	}

	res, err := Run(context.Background(), client, Config{Retries: 1})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Block != block {
		t.Fatalf("expected returned block pointer")
	}
	if res.BlockNumber != 123 {
		t.Fatalf("unexpected block number %d", res.BlockNumber)
	}
	if res.NetworkID.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("unexpected network id %s", res.NetworkID)
	}
	if client.blockCalls != 1 {
		t.Fatalf("expected single block fetch, got %d", client.blockCalls)
	}
}

func TestRunRetries(t *testing.T) {
	block := types.NewBlockWithHeader(&types.Header{
		Number:     big.NewInt(100),
		ParentHash: common.HexToHash("0xbeef"),
	})
	client := &mockClient{
		blockNumber: 100,
		networkID:   big.NewInt(1),
		block:       block,
		blockErrs: []error{
			errors.New("temporary"),
			nil,
		},
	}

	start := time.Now()
	res, err := Run(context.Background(), client, Config{Retries: 2})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Block != block {
		t.Fatalf("expected final block")
	}
	if client.blockCalls != 2 {
		t.Fatalf("expected 2 block calls, got %d", client.blockCalls)
	}
	if time.Since(start) < retryDelay {
		t.Fatalf("expected at least one retry delay")
	}
}

func TestRunPropagatesErrors(t *testing.T) {
	client := &mockClient{
		blockNumberErr: errors.New("boom"),
	}
	if _, err := Run(context.Background(), client, Config{}); err == nil {
		t.Fatalf("expected error from BlockNumber")
	}
}
