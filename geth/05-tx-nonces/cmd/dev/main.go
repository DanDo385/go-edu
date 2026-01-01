package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/05-tx-nonces/internal/txnonces"
)

func mustKey(hex string) *ecdsa.PrivateKey {
	hex = strings.TrimPrefix(hex, "0x")
	k, err := crypto.HexToECDSA(hex)
	if err != nil {
		panic(err)
	}
	return k
}

func main() {
	// Deterministic harness: NoSend=true so this is safe.
	//
	// BREAKPOINT: set a throwaway private key for local devnets.
	rpcURL := "https://eth.llamarpc.com"
	privateKeyHex := "0x" // set me
	to := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	amountWei := big.NewInt(1)

	if privateKeyHex == "0x" {
		fmt.Println("Set privateKeyHex in cmd/dev before running (use a devnet key).")
		return
	}

	pk := mustKey(privateKeyHex)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := txnonces.Run(ctx, client, txnonces.Config{
		PrivateKey: pk,
		To:         to,
		AmountWei:  amountWei,
		GasLimit:   21000,
		NoSend:     true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("From:", out.FromAddress.Hex())
	fmt.Println("Nonce:", out.Nonce)
	fmt.Println("TxHash:", out.Tx.Hash().Hex())
}
