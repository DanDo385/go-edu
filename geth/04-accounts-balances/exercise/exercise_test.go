package exercise

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type mockAccountClient struct {
	balances map[common.Address]*big.Int
	codes    map[common.Address][]byte
	errBal   error
	errCode  error
	lastReq  []requestRecord
}

type requestRecord struct {
	addr common.Address
	kind string
}

func (m *mockAccountClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	m.lastReq = append(m.lastReq, requestRecord{addr: account, kind: "balance"})
	if m.errBal != nil {
		return nil, m.errBal
	}
	return m.balances[account], nil
}

func (m *mockAccountClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	m.lastReq = append(m.lastReq, requestRecord{addr: account, kind: "code"})
	if m.errCode != nil {
		return nil, m.errCode
	}
	return m.codes[account], nil
}

func addr(hexStr string) common.Address {
	return common.HexToAddress(hexStr)
}

func TestRunClassifiesAccounts(t *testing.T) {
	a1 := addr("0x0000000000000000000000000000000000000001")
	a2 := addr("0x0000000000000000000000000000000000000002")

	client := &mockAccountClient{
		balances: map[common.Address]*big.Int{
			a1: big.NewInt(10),
			a2: big.NewInt(20),
		},
		codes: map[common.Address][]byte{
			a1: nil,
			a2: {0x60, 0x60},
		},
	}

	res, err := Run(context.Background(), client, Config{
		Addresses: []common.Address{a1, a2},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if len(res.Accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(res.Accounts))
	}
	if res.Accounts[0].Type != AccountTypeEOA {
		t.Fatalf("expected first account to be EOA")
	}
	if res.Accounts[1].Type != AccountTypeContract {
		t.Fatalf("expected second account to be Contract")
	}
	if res.Accounts[0].Balance.Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("unexpected balance copy, got %s", res.Accounts[0].Balance)
	}
	res.Accounts[0].Balance.SetUint64(0)
	if client.balances[a1].Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("original balance should not have been mutated")
	}
	res.Accounts[1].Code[0] = 0x00
	if client.codes[a2][0] != 0x60 {
		t.Fatalf("code slice should be copied")
	}
}

func TestRunErrors(t *testing.T) {
	a1 := addr("0x0000000000000000000000000000000000000001")
	client := &mockAccountClient{
		errBal: errors.New("boom"),
	}
	if _, err := Run(context.Background(), client, Config{Addresses: []common.Address{a1}}); err == nil {
		t.Fatalf("expected balance error")
	}

	client = &mockAccountClient{
		balances: map[common.Address]*big.Int{
			a1: big.NewInt(0),
		},
		errCode: errors.New("fail"),
	}
	if _, err := Run(context.Background(), client, Config{Addresses: []common.Address{a1}}); err == nil {
		t.Fatalf("expected code error")
	}

	if _, err := Run(context.Background(), client, Config{}); err == nil {
		t.Fatalf("expected empty address error")
	}
}
