package exercise

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type mockTXClient struct {
	nonce         uint64
	gasPrice      *big.Int
	chainID       *big.Int
	sendErr       error
	nonceErr      error
	gasErr        error
	chainErr      error
	sentTx        *types.Transaction
	pendingCalled bool
	suggestCalled bool
}

func (m *mockTXClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	m.pendingCalled = true
	if m.nonceErr != nil {
		return 0, m.nonceErr
	}
	return m.nonce, nil
}

func (m *mockTXClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	m.suggestCalled = true
	if m.gasErr != nil {
		return nil, m.gasErr
	}
	return m.gasPrice, nil
}

func (m *mockTXClient) ChainID(ctx context.Context) (*big.Int, error) {
	if m.chainErr != nil {
		return nil, m.chainErr
	}
	return m.chainID, nil
}

func (m *mockTXClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.sentTx = tx
	return nil
}

func mustKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := crypto.HexToECDSA("4c0883a69102937d6231471b5dbb6204fe51296170827906d7a3c0480b813b20")
	if err != nil {
		t.Fatalf("hex to ecdsa: %v", err)
	}
	return key
}

func TestRunSuccessDefaults(t *testing.T) {
	client := &mockTXClient{
		nonce:    7,
		gasPrice: big.NewInt(1_000_000_000),
		chainID:  big.NewInt(1),
	}

	key := mustKey(t)
	to := common.HexToAddress("0x1111111111111111111111111111111111111111")

	res, err := Run(context.Background(), client, Config{
		PrivateKey: key,
		To:         to,
		AmountWei:  big.NewInt(123),
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !client.pendingCalled {
		t.Fatalf("expected PendingNonceAt to be called")
	}
	if !client.suggestCalled {
		t.Fatalf("expected SuggestGasPrice to be called")
	}
	if res.Nonce != 7 {
		t.Fatalf("unexpected nonce %d", res.Nonce)
	}
	if client.sentTx == nil {
		t.Fatalf("expected SendTransaction to be invoked")
	}
	if client.sentTx.To() == nil || *client.sentTx.To() != to {
		t.Fatalf("tx to mismatch")
	}
}

func TestRunOverridesAndNoSend(t *testing.T) {
	client := &mockTXClient{
		gasPrice: big.NewInt(99),
		chainID:  big.NewInt(5),
	}
	nonce := uint64(42)
	key := mustKey(t)
	res, err := Run(context.Background(), client, Config{
		PrivateKey: key,
		To:         common.HexToAddress("0x2222222222222222222222222222222222222222"),
		AmountWei:  big.NewInt(0),
		Nonce:      &nonce,
		GasPrice:   big.NewInt(77),
		NoSend:     true,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if client.pendingCalled {
		t.Fatalf("pending nonce should not be called when override provided")
	}
	if client.suggestCalled {
		t.Fatalf("suggest gas should not be called when override provided")
	}
	if res.Tx.GasPrice().Cmp(big.NewInt(77)) != 0 {
		t.Fatalf("expected override gas price, got %s", res.Tx.GasPrice())
	}
	if client.sentTx != nil {
		t.Fatalf("should not send when NoSend is true")
	}
}

func TestRunErrors(t *testing.T) {
	key := mustKey(t)
	to := common.HexToAddress("0x123")
	client := &mockTXClient{chainID: big.NewInt(1), gasPrice: big.NewInt(1), nonceErr: errors.New("boom")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected nonce error")
	}
	client = &mockTXClient{chainErr: errors.New("fail")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected chain error")
	}
	client = &mockTXClient{chainID: big.NewInt(1), gasErr: errors.New("nope")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected gas error")
	}
	client = &mockTXClient{chainID: big.NewInt(1), gasPrice: big.NewInt(1), sendErr: errors.New("bad send")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected send error")
	}
}
