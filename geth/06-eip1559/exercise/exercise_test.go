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

type mockFeeClient struct {
	nonce       uint64
	nonceErr    error
	tip         *big.Int
	tipErr      error
	chainID     *big.Int
	chainErr    error
	header      *types.Header
	headerErr   error
	sendErr     error
	sentTx      *types.Transaction
	nonceCalled bool
	tipCalled   bool
}

func (m *mockFeeClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	m.nonceCalled = true
	if m.nonceErr != nil {
		return 0, m.nonceErr
	}
	return m.nonce, nil
}

func (m *mockFeeClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	m.tipCalled = true
	if m.tipErr != nil {
		return nil, m.tipErr
	}
	return m.tip, nil
}

func (m *mockFeeClient) ChainID(ctx context.Context) (*big.Int, error) {
	if m.chainErr != nil {
		return nil, m.chainErr
	}
	return m.chainID, nil
}

func (m *mockFeeClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if m.headerErr != nil {
		return nil, m.headerErr
	}
	return m.header, nil
}

func (m *mockFeeClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.sentTx = tx
	return nil
}

func mustKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return key
}

func TestRunSuccessDefaults(t *testing.T) {
	header := &types.Header{BaseFee: big.NewInt(30_000_000_000)}
	client := &mockFeeClient{
		nonce:   5,
		tip:     big.NewInt(3_000_000_000),
		chainID: big.NewInt(1),
		header:  header,
	}
	key := mustKey(t)
	to := common.HexToAddress("0x0102030405060708090001020304050607080900")

	res, err := Run(context.Background(), client, Config{
		PrivateKey: key,
		To:         to,
		AmountWei:  big.NewInt(123),
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !client.nonceCalled || !client.tipCalled {
		t.Fatalf("expected nonce/tip suggestions to be invoked")
	}
	if res.BaseFee.Cmp(header.BaseFee) != 0 {
		t.Fatalf("base fee mismatch")
	}
	if client.sentTx == nil {
		t.Fatalf("expected send to be called")
	}
	if client.sentTx.GasFeeCap().Cmp(big.NewInt(63_000_000_000)) != 0 { // base*2 + tip
		t.Fatalf("unexpected fee cap %s", client.sentTx.GasFeeCap())
	}
	if client.sentTx.GasTipCap().Cmp(big.NewInt(3_000_000_000)) != 0 {
		t.Fatalf("unexpected tip %s", client.sentTx.GasTipCap())
	}
}

func TestRunOverridesNoSend(t *testing.T) {
	client := &mockFeeClient{
		chainID: big.NewInt(11155111),
		header:  &types.Header{BaseFee: big.NewInt(1_000_000_000)},
	}
	key := mustKey(t)
	nonce := uint64(99)
	res, err := Run(context.Background(), client, Config{
		PrivateKey:     key,
		To:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
		AmountWei:      big.NewInt(0),
		Nonce:          &nonce,
		MaxPriorityFee: big.NewInt(2),
		MaxFee:         big.NewInt(10),
		NoSend:         true,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if client.sentTx != nil {
		t.Fatalf("should not send when NoSend true")
	}
	if res.Nonce != 99 {
		t.Fatalf("nonce override ignored")
	}
	if client.nonceCalled || client.tipCalled {
		t.Fatalf("should not call nonce/tip suggestions when overrides supplied")
	}
	if res.Tx.GasFeeCap().Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("max fee override ignored")
	}
}

func TestRunErrors(t *testing.T) {
	key := mustKey(t)
	to := common.HexToAddress("0x1")
	client := &mockFeeClient{nonceErr: errors.New("boom")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected nonce error")
	}
	client = &mockFeeClient{chainErr: errors.New("id")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected chain error")
	}
	client = &mockFeeClient{chainID: big.NewInt(1), headerErr: errors.New("hdr")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected header error")
	}
	client = &mockFeeClient{chainID: big.NewInt(1), header: &types.Header{}}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected base fee missing error")
	}
	client = &mockFeeClient{chainID: big.NewInt(1), header: &types.Header{BaseFee: big.NewInt(1)}, tipErr: errors.New("tip")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected tip error")
	}
	client = &mockFeeClient{chainID: big.NewInt(1), header: &types.Header{BaseFee: big.NewInt(1)}, tip: big.NewInt(1), sendErr: errors.New("send")}
	if _, err := Run(context.Background(), client, Config{PrivateKey: key, To: to, AmountWei: big.NewInt(1)}); err == nil {
		t.Fatalf("expected send error")
	}
}
