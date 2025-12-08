package exercise

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
)

type mockProofClient struct {
	resp    *gethclient.AccountResult
	err     error
	account common.Address
	slots   []string
	block   *big.Int
}

func (m *mockProofClient) GetProof(ctx context.Context, account common.Address, slots []string, blockNumber *big.Int) (*gethclient.AccountResult, error) {
	m.account = account
	m.slots = append([]string(nil), slots...)
	m.block = blockNumber
	if m.err != nil {
		return nil, m.err
	}
	return m.resp, nil
}

func TestRunSuccess(t *testing.T) {
	account := common.HexToAddress("0x000000000000000000000000000000000000c0de")
	slot := common.HexToHash("0x01")
	resp := &gethclient.AccountResult{
		Balance:      big.NewInt(123),
		Nonce:        7,
		CodeHash:     common.HexToHash("0x02"),
		StorageHash:  common.HexToHash("0x03"),
		AccountProof: []string{"node1", "node2"},
		StorageProof: []gethclient.StorageResult{
			{
				Key:   slot.Hex(),
				Value: big.NewInt(42),
				Proof: []string{"snode1"},
			},
		},
	}

	client := &mockProofClient{resp: resp}

	res, err := Run(context.Background(), client, Config{
		Account: account,
		Slots:   []common.Hash{slot},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Account.Balance.Cmp(big.NewInt(123)) != 0 {
		t.Fatalf("balance mismatch")
	}
	if len(res.Account.Storage) != 1 {
		t.Fatalf("expected storage proof")
	}
	if res.Account.Storage[0].Key != slot {
		t.Fatalf("slot mismatch")
	}
	if client.slots[0] != slot.Hex() {
		t.Fatalf("client slots not hex encoded")
	}
}

func TestRunErrors(t *testing.T) {
	client := &mockProofClient{err: errors.New("boom")}
	if _, err := Run(context.Background(), client, Config{Account: common.HexToAddress("0x1")}); err == nil {
		t.Fatalf("expected error from GetProof")
	}
	if _, err := Run(context.Background(), client, Config{}); err == nil {
		t.Fatalf("expected account validation error")
	}
	if _, err := Run(context.Background(), nil, Config{Account: common.HexToAddress("0x1")}); err == nil {
		t.Fatalf("expected nil client error")
	}
}
