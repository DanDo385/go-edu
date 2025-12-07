package exercise

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type mockRPC struct {
	chainID   *big.Int
	networkID *big.Int
	header    *types.Header
	chainErr  error
	netErr    error
	headerErr error
	lastArgs  []*big.Int
}

func (m *mockRPC) ChainID(ctx context.Context) (*big.Int, error) {
	return m.chainID, m.chainErr
}

func (m *mockRPC) NetworkID(ctx context.Context) (*big.Int, error) {
	return m.networkID, m.netErr
}

func (m *mockRPC) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	m.lastArgs = append(m.lastArgs, number)
	return m.header, m.headerErr
}

func TestRunSuccess(t *testing.T) {
	header := &types.Header{
		Number:     big.NewInt(123),
		ParentHash: common.HexToHash("0xbeef"),
		GasUsed:    42,
	}
	mock := &mockRPC{
		chainID:   big.NewInt(1),
		networkID: big.NewInt(1),
		header:    header,
	}

	res, err := Run(context.Background(), mock, Config{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.ChainID.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("unexpected chain id %s", res.ChainID)
	}
	if res.NetworkID.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("unexpected network id %s", res.NetworkID)
	}
	if res.Header.Number.Uint64() != 123 {
		t.Fatalf("unexpected header number %d", res.Header.Number.Uint64())
	}

	// Mutating the result should not mutate the mock responses (defensive copy)
	res.ChainID.SetUint64(99)
	if mock.chainID.Uint64() != 1 {
		t.Fatalf("chain id was not copied")
	}
	res.Header.Number.SetUint64(0)
	if header.Number.Uint64() != 123 {
		t.Fatalf("header was not copied")
	}

	if len(mock.lastArgs) != 1 || mock.lastArgs[0] != nil {
		t.Fatalf("expected HeaderByNumber to be called with nil (latest), got %v", mock.lastArgs)
	}
}

func TestRunErrors(t *testing.T) {
	if _, err := Run(context.Background(), nil, Config{}); err == nil {
		t.Fatalf("expected error for nil client")
	}

	mock := &mockRPC{chainErr: errors.New("boom")}
	if _, err := Run(context.Background(), mock, Config{}); err == nil {
		t.Fatalf("expected chain error")
	}

	mock = &mockRPC{
		chainID: big.NewInt(1),
		netErr:  errors.New("fail"),
	}
	if _, err := Run(context.Background(), mock, Config{}); err == nil {
		t.Fatalf("expected network error")
	}

	mock = &mockRPC{
		chainID:   big.NewInt(1),
		networkID: big.NewInt(1),
		headerErr: errors.New("nope"),
	}
	if _, err := Run(context.Background(), mock, Config{}); err == nil {
		t.Fatalf("expected header error")
	}
}
